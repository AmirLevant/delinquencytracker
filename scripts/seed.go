package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	_ "github.com/lib/pq"

	dt "github.com/amirlevant/delinquencytracker"
)

func main() {
	// Command-line flags for flexibility
	numUsers := flag.Int("users", 5, "number of users to create")
	numLoans := flag.Int("loans", 10, "total number of loans to create")
	teardown := flag.Bool("teardown", false, "delete all existing data before seeding")
	flag.Parse()

	// Connect to database
	db, err := connectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("✅ Connected to database")

	// Optional teardown
	if *teardown {
		log.Println("🗑️  Tearing down existing data...")
		if err := teardownData(db); err != nil {
			log.Fatalf("Teardown failed: %v", err)
		}
		log.Println("✅ Data cleared")
	}

	// Seed the database
	log.Printf("🌱 Starting seed: %d users, %d loans\n", *numUsers, *numLoans)
	if err := seedDatabase(db, *numUsers, *numLoans); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	// Summary
	printSummary(db)
	log.Println("✅ Seeding complete!")
}

// connectDB establishes connection to PostgreSQL database
func connectDB() (*sql.DB, error) {
	config := "host=localhost port=5432 user=postgres password=amir dbname=loan_tracker sslmode=disable"
	db, err := sql.Open("postgres", config)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return db, nil
}

// teardownData deletes all data from tables (respects FK constraints)
func teardownData(db *sql.DB) error {
	queries := []string{
		"DELETE FROM payments",
		"DELETE FROM loans",
		"DELETE FROM users",
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute '%s': %w", query, err)
		}
	}

	return nil
}

// seedDatabase creates users and distributes loans among them
func seedDatabase(db *sql.DB, numUsers, totalLoans int) error {
	// Calculate base loans per user and remainder
	loansPerUser := totalLoans / numUsers
	remainderLoans := totalLoans % numUsers

	userIDs := make([]int64, 0, numUsers)

	// Create users with their base number of loans
	for i := 0; i < numUsers; i++ {
		// Generate fake user data
		name := gofakeit.Name()
		email := gofakeit.Email()
		phone := gofakeit.Phone()

		// Determine how many loans this user gets
		numLoansForUser := loansPerUser
		if i < remainderLoans {
			numLoansForUser++ // Distribute remainder loans to first few users
		}

		log.Printf("Creating user %d/%d: %s (%d loan(s))", i+1, numUsers, name, numLoansForUser)

		// Create user with first loan using InitializeUserWithLoan
		usr, err := createUserWithFirstLoan(db, name, email, phone)
		if err != nil {
			return fmt.Errorf("failed to create user %s: %w", name, err)
		}

		userIDs = append(userIDs, usr.ID)
		log.Printf("  ✓ User ID %d created with loan ID %d", usr.ID, usr.Loans[0].ID)

		// Add additional loans to this user using AddLoanToExistingUser
		for j := 1; j < numLoansForUser; j++ {
			loan, err := addRandomLoan(db, usr.ID)
			if err != nil {
				return fmt.Errorf("failed to add loan %d for user %d: %w", j+1, usr.ID, err)
			}
			log.Printf("  ✓ Added loan ID %d to user %d", loan.ID, usr.ID)
		}
	}

	return nil
}

// createUserWithFirstLoan creates a user with their first loan and payment schedule
// Uses InitializeUserWithLoan from business.go
func createUserWithFirstLoan(db *sql.DB, name, email, phone string) (dt.user, error) {
	// Generate random but realistic loan parameters
	totalAmount := randomLoanAmount()
	interestRate := randomInterestRate()
	termMonths := randomTermMonths()
	dayDue := randomDayDue()
	dateTaken := randomPastDate()

	// Use the robust business function that handles everything
	usr, err := dt.InitializeUserWithLoan(
		db,
		name, email, phone,
		totalAmount, interestRate, termMonths, dayDue, dateTaken,
	)

	if err != nil {
		return dt.User{}, err
	}

	return usr, nil
}

// addRandomLoan adds an additional loan to an existing user
// Uses AddLoanToExistingUser from business.go
func addRandomLoan(db *sql.DB, userID int64) (dt.Loan, error) {
	totalAmount := randomLoanAmount()
	interestRate := randomInterestRate()
	termMonths := randomTermMonths()
	dayDue := randomDayDue()
	dateTaken := randomPastDate()

	// Use the robust business function
	loan, err := dt.AddLoanToExistingUser(
		db,
		userID,
		totalAmount, interestRate, termMonths, dayDue, dateTaken,
	)

	if err != nil {
		return dt.loan{}, err
	}

	return loan, nil
}

// Random data generators for realistic loan parameters

func randomLoanAmount() float64 {
	// Loan amounts between $5,000 and $50,000
	amounts := []float64{
		5000, 7500, 10000, 12500, 15000,
		20000, 25000, 30000, 35000, 40000, 50000,
	}
	return amounts[rand.Intn(len(amounts))]
}

func randomInterestRate() float64 {
	// Interest rates between 3% and 18%
	rates := []float64{
		0.03, 0.045, 0.06, 0.075, 0.09,
		0.105, 0.12, 0.135, 0.15, 0.165, 0.18,
	}
	return rates[rand.Intn(len(rates))]
}

func randomTermMonths() int {
	// Common loan terms: 1, 2, 3, 4, or 5 years
	terms := []int{12, 24, 36, 48, 60}
	return terms[rand.Intn(len(terms))]
}

func randomDayDue() int {
	// Due dates between 1st and 28th (avoids month-length issues)
	return rand.Intn(28) + 1
}

func randomPastDate() time.Time {
	// Random date between 1 and 365 days ago
	daysAgo := rand.Intn(365) + 1
	return time.Now().AddDate(0, 0, -daysAgo).UTC()
}

// printSummary shows final counts in database
func printSummary(db *sql.DB) {
	userCount, _ := dt.CountUsers(db)

	// Count loans
	var loanCount int64
	db.QueryRow("SELECT COUNT(*) FROM loans").Scan(&loanCount)

	// Count payments
	var paymentCount int64
	db.QueryRow("SELECT COUNT(*) FROM payments").Scan(&paymentCount)

	log.Println("\n📊 Database Summary:")
	log.Printf("   Users:    %d", userCount)
	log.Printf("   Loans:    %d", loanCount)
	log.Printf("   Payments: %d", paymentCount)
}
