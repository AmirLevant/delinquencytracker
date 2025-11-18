package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	dt "github.com/amirlevant/delinquencytracker"
)

func main() {

	db := SetupDatabaseConnection()
	defer db.Close()
	CleanDatabaseData(db)

	if err := PopulateTestUsersLoansPayments(db); err != nil {
		log.Fatalf("Failed to populate test data: %v", err)
	}

	fmt.Println("Database seeded successfully!")
}

func SetupDatabaseConnection() *sql.DB {
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "amir"
		dbname   = "loan_tracker"
		sslmode  = "disable"
	)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		panic(err)
	}

	return db
}

func CleanDatabaseData(db *sql.DB) {
	if _, err := db.Exec("DELETE FROM payments"); err != nil {
		log.Printf("Warning: failed to delete payments: %v", err)
	}
	if _, err := db.Exec("DELETE FROM loans"); err != nil {
		log.Printf("Warning: failed to delete loans: %v", err)
	}
	if _, err := db.Exec("DELETE FROM users"); err != nil {
		log.Printf("Warning: failed to delete users: %v", err)
	}
	fmt.Println("Database cleaned successfully")
}

func PopulateTestUsersLoansPayments(db *sql.DB) error {

	// Define static dates for consistent test data
	// Using dates in the past to simulate historical loans
	date1 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC) // Jan 15, 2024
	date2 := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)  // Mar 1, 2024
	date3 := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC) // Jun 10, 2024
	date4 := time.Date(2023, 11, 5, 0, 0, 0, 0, time.UTC) // Nov 5, 2023
	date5 := time.Date(2024, 8, 20, 0, 0, 0, 0, time.UTC) // Aug 20, 2024

	// User 1: John Smith - Short term, low interest car loan
	usr1, err := dt.InitializeUserWithLoan(
		db,
		"John Smith",
		"john.smith@email.com",
		"555-0101",
		15000.00, // $15,000 loan
		0.045,    // 4.5% annual interest rate
		36,       // 36 months (3 years)
		5,        // Payment due on the 5th of each month
		date1,
	)
	if err != nil {
		return fmt.Errorf("failed to create user 1: %w", err)
	}
	fmt.Printf("Created user: %s (ID: %d)\n", usr1.Name, usr1.ID)

	// User 2: Maria Garcia - Mortgage with longer term
	usr2, err := dt.InitializeUserWithLoan(
		db,
		"Maria Garcia",
		"maria.garcia@email.com",
		"555-0102",
		250000.00, // $250,000 loan
		0.035,     // 3.5% annual interest rate
		360,       // 360 months (30 years)
		1,         // Payment due on the 1st of each month
		date2,
	)
	if err != nil {
		return fmt.Errorf("failed to create user 2: %w", err)
	}
	fmt.Printf("Created user: %s (ID: %d)\n", usr2.Name, usr2.ID)

	// User 3: David Lee - Personal loan, medium term
	usr3, err := dt.InitializeUserWithLoan(
		db,
		"David Lee",
		"david.lee@email.com",
		"555-0103",
		8000.00, // $8,000 loan
		0.0899,  // 8.99% annual interest rate
		24,      // 24 months (2 years)
		15,      // Payment due on the 15th of each month
		date3,
	)
	if err != nil {
		return fmt.Errorf("failed to create user 3: %w", err)
	}
	fmt.Printf("Created user: %s (ID: %d)\n", usr3.Name, usr3.ID)

	// User 4: Sarah Johnson - Student loan with 0% interest
	usr4, err := dt.InitializeUserWithLoan(
		db,
		"Sarah Johnson",
		"sarah.johnson@email.com",
		"555-0104",
		35000.00, // $35,000 loan
		0.00,     // 0% interest (special case handled in code)
		120,      // 120 months (10 years)
		28,       // Payment due on the 28th of each month
		date4,
	)
	if err != nil {
		return fmt.Errorf("failed to create user 4: %w", err)
	}
	fmt.Printf("Created user: %s (ID: %d)\n", usr4.Name, usr4.ID)

	// User 5: Robert Chen - Business loan, high amount
	usr5, err := dt.InitializeUserWithLoan(
		db,
		"Robert Chen",
		"robert.chen@email.com",
		"555-0105",
		75000.00, // $75,000 loan
		0.065,    // 6.5% annual interest rate
		60,       // 60 months (5 years)
		10,       // Payment due on the 10th of each month
		date5,
	)
	if err != nil {
		return fmt.Errorf("failed to create user 5: %w", err)
	}
	fmt.Printf("Created user: %s (ID: %d)\n", usr5.Name, usr5.ID)

	// Add a second loan to one user to test multiple loans per user
	loan2ForUsr1, err := dt.AddLoanToExistingUser(
		db,
		usr1.ID,
		5000.00, // $5,000 second loan
		0.0699,  // 6.99% annual interest rate
		12,      // 12 months (1 year)
		5,       // Payment due on the 5th
		time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC),
	)
	if err != nil {
		return fmt.Errorf("failed to add second loan to user 1: %w", err)
	}
	fmt.Printf("Added second loan to %s (Loan ID: %d)\n", usr1.Name, loan2ForUsr1.ID)

	fmt.Println("\n=== Test Data Summary ===")
	fmt.Printf("Total users created: 5\n")
	fmt.Printf("Total loans created: 6\n")

	return nil
}
