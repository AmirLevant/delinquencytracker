package delinquencytracker

import (
	"database/sql"
	"fmt"
	"math"
	"time"
)

// calculateMonthlyPayment calculates the monthly payment using the amortization formula.
func calculateMonthlyPayment(principal, annualRate float64, months int) float64 {
	var mnthlyPayment float64
	var mnthlyInterestRate float64 = annualRate / 12

	// special case to avoid Nan
	if annualRate == 0 {
		return principal / float64(months)
	}

	numirator := mnthlyInterestRate * math.Pow(1+mnthlyInterestRate, float64(months))
	denominator := (math.Pow(1+mnthlyInterestRate, float64(months)) - 1)

	mnthlyPayment = principal * (numirator / denominator)

	return mnthlyPayment
}

// calculateDueDate calculates the payment due date by adding months to the start date.
func calculateDueDate(startDate time.Time, termMonths, dayDue int) time.Time {
	// Get the target month by adding months to the start date's year and month
	// We need to work with year and month directly to avoid day overflow issues
	year := startDate.Year()
	month := startDate.Month()

	// Add the months
	month += time.Month(termMonths)

	// Normalize year and month (handle overflow)
	for month > 12 {
		month -= 12
		year++
	}

	// Find last day of the target month
	lastDayOfMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()

	// Adjust the day if it exceeds the month's maximum
	actualDay := dayDue
	if dayDue > lastDayOfMonth {
		actualDay = lastDayOfMonth
	}

	// Return the due date in UTC
	return time.Date(year, month, actualDay, 0, 0, 0, 0, time.UTC)
}

// InitializeUserWithLoan creates a new user with a loan and generates the complete payment schedule.
// Use dateTaken to backdate loans for historical data.
func InitializeUserWithLoan(db *sql.DB, name, email, phone string, totalAmount, interestRate float64, termMonths, dayDue int, dateTaken time.Time) (user, error) {

	// Ensure dateTaken is in UTC for consistency
	dateTaken = dateTaken.UTC()

	// Step 1: Create the user
	// This gives us a user ID that we'll need for the loan
	usr, err := CreateUser(db, name, email, phone)
	if err != nil {
		return user{}, fmt.Errorf("failed to create user: %w", err)
	}

	// Step 2: Create the loan
	// The loan uses the provided dateTaken and starts in "active" status
	ln, err := CreateLoan(db, usr.ID, totalAmount, interestRate, termMonths, dayDue, "active", dateTaken)
	if err != nil {
		// NOTE: At this point, the user exists in the DB but the loan failed
		// In Approach 1, we don't clean this up automatically
		// You could add cleanup logic here if desired
		return user{}, fmt.Errorf("failed to create loan for user %d: %w", usr.ID, err)
	}

	// Step 3: Calculate the monthly payment amount
	monthlyPayment := calculateMonthlyPayment(totalAmount, interestRate, termMonths)

	// Step 4: Create all payment records
	// We'll create one payment record for each month of the loan term
	payments := make([]payment, 0, termMonths)

	for i := 1; i <= termMonths; i++ {
		// Calculate when this payment is due
		dueDate := calculateDueDate(dateTaken, i, dayDue)

		// Create the payment record
		// AmountPaid is 0 because it hasn't been paid yet
		// PaidDate is zero time (time.Time{}) because it's unpaid
		pmt, err := CreatePayment(db, ln.ID, int64(i), monthlyPayment, 0, dueDate, time.Time{})
		if err != nil {
			// NOTE: At this point, we have user + loan + some payments in DB
			// The remaining payments failed to create
			// In Approach 1, we don't clean this up automatically
			return user{}, fmt.Errorf("failed to create payment %d for loan %d: %w", i, ln.ID, err)
		}

		payments = append(payments, pmt)
	}

	// Step 5: Assemble the full user object
	// Attach the payments to the loan
	ln.Payments = payments

	// Attach the loan to the user
	usr.Loans = []loan{ln}

	// Step 6: Return the fully populated user
	return usr, nil
}

// InitializeUserWithLoanNow creates a new user with a loan starting today.
func InitializeUserWithLoanNow(db *sql.DB, name, email, phone string,
	totalAmount, interestRate float64, termMonths, dayDue int) (user, error) {
	return InitializeUserWithLoan(db, name, email, phone, totalAmount, interestRate,
		termMonths, dayDue, time.Now().UTC())
}

// AddLoanToExistingUser adds a new loan with payment schedule to an existing user.
func AddLoanToExistingUser(db *sql.DB, userID int64, totalAmount, interestRate float64,
	termMonths, dayDue int, dateTaken time.Time) (loan, error) {

	// Ensure dateTaken is in UTC for consistency
	dateTaken = dateTaken.UTC()

	// Verify user exists
	_, err := GetUserByID(db, userID)
	if err != nil {
		return loan{}, fmt.Errorf("user %d not found: %w", userID, err)
	}

	// Create the loan
	ln, err := CreateLoan(db, userID, totalAmount, interestRate, termMonths, dayDue, "active", dateTaken)
	if err != nil {
		return loan{}, fmt.Errorf("failed to create loan for user %d: %w", userID, err)
	}

	// Calculate the monthly payment amount
	monthlyPayment := calculateMonthlyPayment(totalAmount, interestRate, termMonths)

	// Create all payment records
	payments := make([]payment, 0, termMonths)

	for i := 1; i <= termMonths; i++ {
		dueDate := calculateDueDate(dateTaken, i, dayDue)

		pmt, err := CreatePayment(db, ln.ID, int64(i), monthlyPayment, 0, dueDate, time.Time{})
		if err != nil {
			return loan{}, fmt.Errorf("failed to create payment %d for loan %d: %w", i, ln.ID, err)
		}

		payments = append(payments, pmt)
	}

	// Attach payments to the loan
	ln.Payments = payments

	return ln, nil
}

// AddLoanToExistingUserNow adds a loan starting today to an existing user.
func AddLoanToExistingUserNow(db *sql.DB, userID int64, totalAmount, interestRate float64,
	termMonths, dayDue int) (loan, error) {
	return AddLoanToExistingUser(db, userID, totalAmount, interestRate,
		termMonths, dayDue, time.Now().UTC())
}

// GetFullUserByID retrieves a user with all their loans and payments.
func GetFullUserByID(db *sql.DB, userID int64) (user, error) {
	// Step 1: Get the basic user information
	usr, err := GetUserByID(db, userID)
	if err != nil {
		return user{}, fmt.Errorf("failed to get user: %w", err)
	}

	// Step 2: Get all loans for this user
	loans, err := GetLoansByUserID(db, userID)
	if err != nil {
		return user{}, fmt.Errorf("failed to get loans for user %d: %w", userID, err)
	}

	// Step 3: For each loan, get all its payments
	for i := range loans {
		payments, err := GetPaymentsByLoanID(db, loans[i].ID)
		if err != nil {
			return user{}, fmt.Errorf("failed to get payments for loan %d: %w", loans[i].ID, err)
		}
		loans[i].Payments = payments
	}

	// Step 4: Attach all loans to the user
	usr.Loans = loans

	return usr, nil
}

// GetFullLoanByID retrieves a loan with all its payment information.
func GetFullLoanByID(db *sql.DB, loanID int64) (loan, error) {
	// Step 1: Get the basic loan information
	ln, err := GetLoanByLoanID(db, loanID)
	if err != nil {
		return loan{}, fmt.Errorf("failed to get loan: %w", err)
	}

	// Step 2: Get all payments for this loan
	payments, err := GetPaymentsByLoanID(db, loanID)
	if err != nil {
		return loan{}, fmt.Errorf("failed to get payments for loan %d: %w", loanID, err)
	}

	// Step 3: Attach payments to the loan
	ln.Payments = payments

	return ln, nil
}
