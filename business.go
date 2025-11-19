package delinquencytracker

import (
	"database/sql"
	"fmt"
	"math"
	"time"
)

// calculateMonthlyPayment calculates the monthly Payment using the amortization formula.
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

// calculateDueDate calculates the Payment due date by adding months to the start date.
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

// validateLoanParameters validates the input parameters for creating a Loan.
func validateLoanParameters(totalAmount, interestRate float64, termMonths, dayDue int, dateTaken time.Time) error {
	if totalAmount <= 0 {
		return fmt.Errorf("totalAmount must be positive, got %.2f", totalAmount)
	}

	if interestRate < 0 {
		return fmt.Errorf("interestRate cannot be negative, got %.4f", interestRate)
	}

	if termMonths <= 0 {
		return fmt.Errorf("termMonths must be positive, got %d", termMonths)
	}

	if dayDue < 1 || dayDue > 31 {
		return fmt.Errorf("dayDue must be between 1 and 31, got %d", dayDue)
	}

	// Allow dateTaken to be in the past, present, or future
	// Just ensure it's a valid time
	if dateTaken.IsZero() {
		return fmt.Errorf("dateTaken cannot be zero time")
	}

	return nil
}

// createPaymentSchedule generates the complete Payment schedule for a Loan.
// If autoPayPastDue is true, payments with due dates before now will be marked as paid.
// The paidDate for auto-paid payments will be set to the dueDate (assumes on-time payment).
func createPaymentSchedule(db *sql.DB, loanID int64, principal, annualRate float64,
	termMonths, dayDue int, dateTaken time.Time, autoPayPastDue bool) ([]Payment, error) {

	monthlyPayment := calculateMonthlyPayment(principal, annualRate, termMonths)
	payments := make([]Payment, 0, termMonths)
	now := time.Now().UTC()

	for i := 1; i <= termMonths; i++ {
		dueDate := calculateDueDate(dateTaken, i, dayDue)

		// Determine if this payment should be marked as paid
		var amountPaid float64
		var paidDate time.Time

		if autoPayPastDue && dueDate.Before(now) {
			// Payment is in the past - mark as paid with on-time payment
			amountPaid = monthlyPayment
			paidDate = dueDate
		} else {
			// Payment is in the future or we're not auto-paying - leave unpaid
			amountPaid = 0
			paidDate = time.Time{}
		}

		pmt, err := CreatePayment(db, loanID, int64(i), monthlyPayment, amountPaid, dueDate, paidDate)
		if err != nil {
			return nil, fmt.Errorf("failed to create Payment %d: %w", i, err)
		}

		payments = append(payments, pmt)
	}

	return payments, nil
}

// InitializeUserWithLoan creates a new User with a Loan and generates the complete Payment schedule.
// Use dateTaken to backdate loans for historical data.
// If autoPayPastDue is true, payments with due dates before today will be automatically marked as paid.
func InitializeUserWithLoan(db *sql.DB, name, email, phone string, totalAmount, interestRate float64,
	termMonths, dayDue int, dateTaken time.Time, autoPayPastDue bool) (User, error) {

	// Ensure dateTaken is in UTC for consistency
	dateTaken = dateTaken.UTC()

	// Validate input parameters
	if err := validateLoanParameters(totalAmount, interestRate, termMonths, dayDue, dateTaken); err != nil {
		return User{}, fmt.Errorf("invalid loan parameters: %w", err)
	}

	// Step 1: Create the User
	usr, err := CreateUser(db, name, email, phone)
	if err != nil {
		return User{}, fmt.Errorf("failed to create User: %w", err)
	}

	// Step 2: Create the Loan
	ln, err := CreateLoan(db, usr.ID, totalAmount, interestRate, termMonths, dayDue, "active", dateTaken)
	if err != nil {
		return User{}, fmt.Errorf("failed to create Loan for User %d: %w", usr.ID, err)
	}

	// Step 3: Create all Payment records
	payments, err := createPaymentSchedule(db, ln.ID, totalAmount, interestRate, termMonths, dayDue, dateTaken, autoPayPastDue)
	if err != nil {
		return User{}, fmt.Errorf("failed to create payment schedule for Loan %d: %w", ln.ID, err)
	}

	// Step 4: Assemble the full User object
	ln.Payments = payments
	usr.Loans = []Loan{ln}

	return usr, nil
}

// InitializeUserWithLoanNow creates a new User with a Loan starting today.
// All past payments (none in this case) will not be auto-paid since the Loan starts now.
func InitializeUserWithLoanNow(db *sql.DB, name, email, phone string,
	totalAmount, interestRate float64, termMonths, dayDue int) (User, error) {
	// When creating a loan starting now, there are no past payments to auto-pay
	return InitializeUserWithLoan(db, name, email, phone, totalAmount, interestRate,
		termMonths, dayDue, time.Now().UTC(), false)
}

// InitializeUserWithLoanNowAutoPay creates a new User with a Loan starting today.
// This is primarily for testing or special cases where you might want autoPayPastDue enabled.
func InitializeUserWithLoanNowAutoPay(db *sql.DB, name, email, phone string,
	totalAmount, interestRate float64, termMonths, dayDue int, autoPayPastDue bool) (User, error) {
	return InitializeUserWithLoan(db, name, email, phone, totalAmount, interestRate,
		termMonths, dayDue, time.Now().UTC(), autoPayPastDue)
}

// AddLoanToExistingUser adds a new Loan with Payment schedule to an existing User.
// If autoPayPastDue is true, payments with due dates before today will be automatically marked as paid.
func AddLoanToExistingUser(db *sql.DB, userID int64, totalAmount, interestRate float64,
	termMonths, dayDue int, dateTaken time.Time, autoPayPastDue bool) (Loan, error) {

	// Ensure dateTaken is in UTC for consistency
	dateTaken = dateTaken.UTC()

	// Validate input parameters
	if err := validateLoanParameters(totalAmount, interestRate, termMonths, dayDue, dateTaken); err != nil {
		return Loan{}, fmt.Errorf("invalid loan parameters: %w", err)
	}

	// Step 1: Verify User exists
	_, err := GetUserByID(db, userID)
	if err != nil {
		return Loan{}, fmt.Errorf("User %d not found: %w", userID, err)
	}

	// Step 2: Create the Loan
	ln, err := CreateLoan(db, userID, totalAmount, interestRate, termMonths, dayDue, "active", dateTaken)
	if err != nil {
		return Loan{}, fmt.Errorf("failed to create Loan for User %d: %w", userID, err)
	}

	// Step 3: Create all Payment records
	payments, err := createPaymentSchedule(db, ln.ID, totalAmount, interestRate, termMonths, dayDue, dateTaken, autoPayPastDue)
	if err != nil {
		return Loan{}, fmt.Errorf("failed to create payment schedule for Loan %d: %w", ln.ID, err)
	}

	// Step 4: Attach payments to the Loan
	ln.Payments = payments

	return ln, nil
}

// AddLoanToExistingUserNow adds a Loan starting today to an existing User.
// All past payments (none in this case) will not be auto-paid since the Loan starts now.
func AddLoanToExistingUserNow(db *sql.DB, userID int64, totalAmount, interestRate float64,
	termMonths, dayDue int) (Loan, error) {
	// When creating a loan starting now, there are no past payments to auto-pay
	return AddLoanToExistingUser(db, userID, totalAmount, interestRate,
		termMonths, dayDue, time.Now().UTC(), false)
}

// AddLoanToExistingUserNowAutoPay adds a Loan starting today to an existing User.
// This is primarily for testing or special cases where you might want autoPayPastDue enabled.
func AddLoanToExistingUserNowAutoPay(db *sql.DB, userID int64, totalAmount, interestRate float64,
	termMonths, dayDue int, autoPayPastDue bool) (Loan, error) {
	return AddLoanToExistingUser(db, userID, totalAmount, interestRate,
		termMonths, dayDue, time.Now().UTC(), autoPayPastDue)
}

// GetFullUserByID retrieves a User with all their loans and payments.
func GetFullUserByID(db *sql.DB, userID int64) (User, error) {
	// Step 1: Get the basic User information
	usr, err := GetUserByID(db, userID)
	if err != nil {
		return User{}, fmt.Errorf("failed to get User: %w", err)
	}

	// Step 2: Get all loans for this User
	loans, err := GetLoansByUserID(db, userID)
	if err != nil {
		return User{}, fmt.Errorf("failed to get loans for User %d: %w", userID, err)
	}

	// Step 3: For each Loan, get all its payments
	for i := range loans {
		payments, err := GetPaymentsByLoanID(db, loans[i].ID)
		if err != nil {
			return User{}, fmt.Errorf("failed to get payments for Loan %d: %w", loans[i].ID, err)
		}
		loans[i].Payments = payments
	}

	// Step 4: Attach all loans to the User
	usr.Loans = loans

	return usr, nil
}

// GetFullLoanByID retrieves a Loan with all its Payment information.
func GetFullLoanByID(db *sql.DB, loanID int64) (Loan, error) {
	// Step 1: Get the basic Loan information
	ln, err := GetLoanByLoanID(db, loanID)
	if err != nil {
		return Loan{}, fmt.Errorf("failed to get Loan: %w", err)
	}

	// Step 2: Get all payments for this Loan
	payments, err := GetPaymentsByLoanID(db, loanID)
	if err != nil {
		return Loan{}, fmt.Errorf("failed to get payments for Loan %d: %w", loanID, err)
	}

	// Step 3: Attach payments to the Loan
	ln.Payments = payments

	return ln, nil
}
