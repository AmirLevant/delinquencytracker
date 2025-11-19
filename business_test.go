package delinquencytracker

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestCalculateMonthlyPayment verifies monthly payment calculations for various loan scenarios.
func TestCalculateMonthlyPayment(t *testing.T) {
	tests := []struct {
		name        string
		principal   float64
		annualRate  float64
		months      int
		expected    float64
		description string
	}{
		{
			name:        "Zero interest loan",
			principal:   12000.0,
			annualRate:  0.0,
			months:      12,
			expected:    1000.0,
			description: "$12,000 at 0% for 12 months = $1,000/month",
		},
		{
			name:        "Standard car loan",
			principal:   20000.0,
			annualRate:  0.05,
			months:      60,
			expected:    377.42,
			description: "$20,000 at 5% APR for 60 months",
		},
		{
			name:        "Small personal loan",
			principal:   5000.0,
			annualRate:  0.08,
			months:      24,
			expected:    226.14,
			description: "$5,000 at 8% APR for 24 months",
		},
		{
			name:        "High interest short term",
			principal:   1000.0,
			annualRate:  0.15,
			months:      6,
			expected:    174.03,
			description: "$1,000 at 15% APR for 6 months",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateMonthlyPayment(tt.principal, tt.annualRate, tt.months)

			// Check for NaN first
			if math.IsNaN(result) {
				t.Errorf("%s: result is NaN", tt.description)
				return
			}

			// Check if result is within $0.01 of expected
			if math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("%s: expected $%.2f, got $%.2f (diff: $%.2f)",
					tt.description, tt.expected, result, result-tt.expected)
			}
		})
	}
}

// TestCalculateMonthlyPaymentTotal verifies that total payments exceed principal due to interest.
func TestCalculateMonthlyPaymentTotal(t *testing.T) {
	// Arrange
	principal := 10000.0
	annualRate := 0.06
	months := 12

	// Act
	monthlyPayment := calculateMonthlyPayment(principal, annualRate, months)
	totalPaid := monthlyPayment * float64(months)

	// Assert
	t.Logf("Principal: $%.2f", principal)
	t.Logf("Monthly payment: $%.2f", monthlyPayment)
	t.Logf("Total paid: $%.2f", totalPaid)
	t.Logf("Interest paid: $%.2f", totalPaid-principal)

	// Total paid should be more than principal (because of interest)
	require.Greater(t, totalPaid, principal, "Total paid should be greater than principal")

	// But not unreasonably high
	maxExpected := principal * (1 + annualRate) // Rough upper bound
	require.Less(t, totalPaid, maxExpected, "Total paid seems too high")
}

// TestCalculateDueDate verifies due date calculations including edge cases.
func TestCalculateDueDate(t *testing.T) {
	startDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		paymentNum  int
		dayDue      int
		expected    time.Time
		description string
	}{
		{
			name:        "First payment",
			paymentNum:  1,
			dayDue:      15,
			expected:    time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC),
			description: "First payment due one month after start",
		},
		{
			name:        "Second payment",
			paymentNum:  2,
			dayDue:      15,
			expected:    time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
			description: "Second payment due two months after start",
		},
		{
			name:        "Different day of month",
			paymentNum:  1,
			dayDue:      5,
			expected:    time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC),
			description: "Payment due on different day",
		},
		{
			name:        "End of month edge case",
			paymentNum:  1,
			dayDue:      31,
			expected:    time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
			description: "When dayDue is 31 but month only has 29 days",
		},
		{
			name:        "Year boundary",
			paymentNum:  12,
			dayDue:      15,
			expected:    time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			description: "Payment crosses year boundary",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := calculateDueDate(startDate, tt.paymentNum, tt.dayDue)

			// Assert
			require.Equal(t, tt.expected, result,
				"%s: expected %s, got %s",
				tt.description,
				tt.expected.Format("2006-01-02"),
				result.Format("2006-01-02"))
			t.Logf("✓ %s: %s", tt.description, result.Format("2006-01-02"))
		})
	}
}

// TestInitializeUserWithLoanUnpaid verifies creation of a user with loan where payments are NOT auto-paid.
func TestInitializeUserWithLoanUnpaid(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange
	dateTaken := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	name := "John Doe"
	email := "john@example.com"
	phone := "555-1234"
	totalAmount := 10000.0
	interestRate := 0.05
	termMonths := 12
	dayDue := 15

	// Act - autoPayPastDue = false
	user, err := InitializeUserWithLoan(db, name, email, phone,
		totalAmount, interestRate, termMonths, dayDue, dateTaken, false)

	// Assert
	require.NoError(t, err, "InitializeUserWithLoan should not return error")
	require.NotEqual(t, int64(0), user.ID, "User should have valid ID")
	require.Equal(t, name, user.Name, "User name should match")
	require.Equal(t, email, user.Email, "User email should match")
	require.Equal(t, phone, user.Phone, "User phone should match")

	// Check loan was created
	require.Len(t, user.Loans, 1, "User should have exactly 1 loan")
	loan := user.Loans[0]
	require.NotEqual(t, int64(0), loan.ID, "Loan should have valid ID")
	require.Equal(t, user.ID, loan.UserID, "Loan should belong to user")
	require.Equal(t, totalAmount, loan.TotalAmount, "Loan amount should match")
	require.Equal(t, interestRate, loan.InterestRate, "Interest rate should match")
	require.Equal(t, termMonths, loan.TermMonths, "Term months should match")
	require.Equal(t, dayDue, loan.DayDue, "Day due should match")
	require.Equal(t, "active", loan.Status, "Loan status should be active")
	require.Equal(t, dateTaken, loan.DateTaken, "Date taken should match")

	// Check payments were created
	require.Len(t, loan.Payments, termMonths, "Should have payment for each month")

	// Verify first payment - should be UNPAID since autoPayPastDue = false
	firstPayment := loan.Payments[0]
	require.Equal(t, int64(1), firstPayment.PaymentNumber, "First payment should be #1")
	require.Equal(t, loan.ID, firstPayment.LoanID, "Payment should belong to loan")
	require.Greater(t, firstPayment.AmountDue, 0.0, "Payment amount should be positive")
	require.Equal(t, 0.0, firstPayment.AmountPaid, "Payment should be unpaid")
	require.True(t, firstPayment.PaidDate.IsZero(), "PaidDate should be zero for unpaid payment")
	expectedFirstDue := time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC)
	require.Equal(t, expectedFirstDue, firstPayment.DueDate, "First payment due date should be correct")

	// Verify last payment
	lastPayment := loan.Payments[termMonths-1]
	require.Equal(t, int64(termMonths), lastPayment.PaymentNumber, "Last payment number should match term")
	expectedLastDue := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	require.Equal(t, expectedLastDue, lastPayment.DueDate, "Last payment due date should be correct")

	// Verify all payments have same amount due
	firstAmount := loan.Payments[0].AmountDue
	for i, pmt := range loan.Payments {
		require.Equal(t, firstAmount, pmt.AmountDue,
			"Payment %d should have same amount as first payment", i+1)
	}

	// Verify ALL payments are unpaid
	for i, pmt := range loan.Payments {
		require.Equal(t, 0.0, pmt.AmountPaid, "Payment %d should be unpaid", i+1)
		require.True(t, pmt.PaidDate.IsZero(), "Payment %d should have zero PaidDate", i+1)
	}

	t.Logf("✓ Successfully created user with loan and %d unpaid payments", termMonths)
}

// TestInitializeUserWithLoanAutoPaid verifies creation of a user with historical loan where past payments are auto-paid.
func TestInitializeUserWithLoanAutoPaid(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange - Create loan that started 6 months ago
	sixMonthsAgo := time.Now().UTC().AddDate(0, -6, 0)
	// Normalize to first of month for predictable testing
	dateTaken := time.Date(sixMonthsAgo.Year(), sixMonthsAgo.Month(), 1, 0, 0, 0, 0, time.UTC)

	name := "Jane Smith"
	email := "jane@example.com"
	phone := "555-5678"
	totalAmount := 10000.0
	interestRate := 0.05
	termMonths := 12
	dayDue := 15

	// Act - autoPayPastDue = true
	user, err := InitializeUserWithLoan(db, name, email, phone,
		totalAmount, interestRate, termMonths, dayDue, dateTaken, true)

	// Assert
	require.NoError(t, err, "InitializeUserWithLoan should not return error")
	require.Len(t, user.Loans, 1, "User should have exactly 1 loan")
	loan := user.Loans[0]
	require.Len(t, loan.Payments, termMonths, "Should have payment for each month")

	// Calculate monthly payment amount
	monthlyPayment := calculateMonthlyPayment(totalAmount, interestRate, termMonths)

	// Count how many payments should be paid (due dates in the past)
	now := time.Now().UTC()
	paidCount := 0
	unpaidCount := 0

	for i, pmt := range loan.Payments {
		if pmt.DueDate.Before(now) {
			// This payment is in the past - should be marked as paid
			require.Equal(t, monthlyPayment, pmt.AmountPaid,
				"Payment %d (due %s) should be paid", i+1, pmt.DueDate.Format("2006-01-02"))
			require.Equal(t, pmt.DueDate, pmt.PaidDate,
				"Payment %d PaidDate should equal DueDate for auto-paid", i+1)
			paidCount++
		} else {
			// This payment is in the future - should be unpaid
			require.Equal(t, 0.0, pmt.AmountPaid,
				"Payment %d (due %s) should be unpaid", i+1, pmt.DueDate.Format("2006-01-02"))
			require.True(t, pmt.PaidDate.IsZero(),
				"Payment %d should have zero PaidDate", i+1)
			unpaidCount++
		}
	}

	t.Logf("✓ Successfully created user with loan: %d paid, %d unpaid payments", paidCount, unpaidCount)
	require.Greater(t, paidCount, 0, "Should have at least some paid payments from 6 months ago")
	require.Greater(t, unpaidCount, 0, "Should have at least some unpaid future payments")
}

// TestInitializeUserWithLoanNow verifies creation of a user with loan starting today.
func TestInitializeUserWithLoanNow(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange
	name := "Bob Johnson"
	email := "bob@example.com"
	phone := "555-9999"

	// Act
	user, err := InitializeUserWithLoanNow(db, name, email, phone,
		5000.0, 0.06, 6, 10)

	// Assert
	require.NoError(t, err, "InitializeUserWithLoanNow should not return error")
	require.NotEqual(t, int64(0), user.ID, "User should have valid ID")
	require.Len(t, user.Loans, 1, "User should have exactly 1 loan")
	require.Len(t, user.Loans[0].Payments, 6, "Loan should have 6 payments")

	// Verify loan started recently (within last minute)
	now := time.Now().UTC()
	timeDiff := now.Sub(user.Loans[0].DateTaken)
	require.Less(t, timeDiff, time.Minute, "Loan should have started within last minute")

	// All payments should be unpaid since loan just started
	for i, pmt := range user.Loans[0].Payments {
		require.Equal(t, 0.0, pmt.AmountPaid, "Payment %d should be unpaid for new loan", i+1)
		require.True(t, pmt.PaidDate.IsZero(), "Payment %d should have zero PaidDate", i+1)
	}

	t.Logf("✓ Successfully created user with current-date loan")
}

// TestAddLoanToExistingUser verifies adding a second loan to an existing user.
func TestAddLoanToExistingUser(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange - Create initial user with a loan
	dateTaken1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	user, err := InitializeUserWithLoan(db, "Alice Cooper", "alice@example.com", "555-1111",
		10000.0, 0.05, 12, 15, dateTaken1, false)
	require.NoError(t, err, "Failed to create initial user")

	// Act - Add second loan to same user
	dateTaken2 := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	secondLoan, err := AddLoanToExistingUser(db, user.ID,
		5000.0, 0.055, 24, 20, dateTaken2, false)

	// Assert
	require.NoError(t, err, "AddLoanToExistingUser should not return error")
	require.NotEqual(t, int64(0), secondLoan.ID, "Second loan should have valid ID")
	require.Equal(t, user.ID, secondLoan.UserID, "Second loan should belong to same user")
	require.Equal(t, 5000.0, secondLoan.TotalAmount, "Second loan amount should match")
	require.Equal(t, 0.055, secondLoan.InterestRate, "Second loan rate should match")
	require.Equal(t, 24, secondLoan.TermMonths, "Second loan term should match")
	require.Equal(t, 20, secondLoan.DayDue, "Second loan day due should match")
	require.Len(t, secondLoan.Payments, 24, "Second loan should have 24 payments")

	// Verify user now has 2 loans in database
	fullUser, err := GetFullUserByID(db, user.ID)
	require.NoError(t, err, "Failed to get full user")
	require.Len(t, fullUser.Loans, 2, "User should now have 2 loans")

	t.Logf("✓ Successfully added second loan to existing user")
}

// TestAddLoanToExistingUserNow verifies adding a loan starting today to an existing user.
func TestAddLoanToExistingUserNow(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange - Create initial user
	user, err := InitializeUserWithLoanNow(db, "Charlie Brown", "charlie@example.com", "555-2222",
		8000.0, 0.06, 18, 5)
	require.NoError(t, err, "Failed to create initial user")

	// Act - Add second loan with current date
	secondLoan, err := AddLoanToExistingUserNow(db, user.ID,
		3000.0, 0.07, 12, 10)

	// Assert
	require.NoError(t, err, "AddLoanToExistingUserNow should not return error")
	require.NotEqual(t, int64(0), secondLoan.ID, "Second loan should have valid ID")
	require.Len(t, secondLoan.Payments, 12, "Second loan should have 12 payments")

	// Verify second loan started recently
	now := time.Now().UTC()
	timeDiff := now.Sub(secondLoan.DateTaken)
	require.Less(t, timeDiff, time.Minute, "Second loan should have started within last minute")

	t.Logf("✓ Successfully added second loan with current date")
}

// TestAddLoanToNonexistentUser verifies error handling when adding loan to nonexistent user.
func TestAddLoanToNonexistentUser(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange - Use a user ID that doesn't exist
	nonexistentUserID := int64(99999)
	dateTaken := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// Act
	_, err := AddLoanToExistingUser(db, nonexistentUserID,
		5000.0, 0.05, 12, 15, dateTaken, false)

	// Assert
	require.Error(t, err, "Should return error for nonexistent user")
	require.Contains(t, err.Error(), "not found", "Error should mention user not found")

	t.Logf("✓ Correctly rejected loan for nonexistent user")
}

// TestGetFullUserByID verifies retrieval of user with all loans and payments.
func TestGetFullUserByID(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange - Create user with multiple loans
	dateTaken1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	user, err := InitializeUserWithLoan(db, "Diana Prince", "diana@example.com", "555-3333",
		10000.0, 0.05, 12, 15, dateTaken1, false)
	require.NoError(t, err, "Failed to create user")

	dateTaken2 := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	_, err = AddLoanToExistingUser(db, user.ID, 5000.0, 0.06, 24, 20, dateTaken2, false)
	require.NoError(t, err, "Failed to add second loan")

	// Act
	fullUser, err := GetFullUserByID(db, user.ID)

	// Assert
	require.NoError(t, err, "GetFullUserByID should not return error")
	require.Equal(t, user.ID, fullUser.ID, "User ID should match")
	require.Equal(t, "Diana Prince", fullUser.Name, "User name should match")
	require.Len(t, fullUser.Loans, 2, "User should have 2 loans")

	// Verify first loan
	loan1 := fullUser.Loans[0]
	require.Equal(t, 10000.0, loan1.TotalAmount, "First loan amount should match")
	require.Len(t, loan1.Payments, 12, "First loan should have 12 payments")

	// Verify second loan
	loan2 := fullUser.Loans[1]
	require.Equal(t, 5000.0, loan2.TotalAmount, "Second loan amount should match")
	require.Len(t, loan2.Payments, 24, "Second loan should have 24 payments")

	t.Logf("✓ Successfully retrieved full user with %d loans", len(fullUser.Loans))
}

// TestGetFullLoanByID verifies retrieval of loan with all payment information.
func TestGetFullLoanByID(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange - Create user with loan
	dateTaken := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	user, err := InitializeUserWithLoan(db, "Eve Adams", "eve@example.com", "555-4444",
		15000.0, 0.055, 36, 10, dateTaken, false)
	require.NoError(t, err, "Failed to create user")

	loanID := user.Loans[0].ID

	// Act
	fullLoan, err := GetFullLoanByID(db, loanID)

	// Assert
	require.NoError(t, err, "GetFullLoanByID should not return error")
	require.Equal(t, loanID, fullLoan.ID, "Loan ID should match")
	require.Equal(t, user.ID, fullLoan.UserID, "User ID should match")
	require.Equal(t, 15000.0, fullLoan.TotalAmount, "Loan amount should match")
	require.Equal(t, 0.055, fullLoan.InterestRate, "Interest rate should match")
	require.Equal(t, 36, fullLoan.TermMonths, "Term months should match")
	require.Len(t, fullLoan.Payments, 36, "Loan should have 36 payments")

	// Verify payments are ordered correctly
	for i, pmt := range fullLoan.Payments {
		require.Equal(t, int64(i+1), pmt.PaymentNumber,
			"Payment %d should have correct payment number", i+1)
	}

	t.Logf("✓ Successfully retrieved full loan with %d payments", len(fullLoan.Payments))
}

// TestInitializeUserWithLoanHistoricalDate verifies backdating loans with historical dates.
func TestInitializeUserWithLoanHistoricalDate(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange - Create loan that started 1 year ago on the 15th
	oneYearAgo := time.Date(time.Now().Year()-1, time.Now().Month(), 15, 0, 0, 0, 0, time.UTC)
	dayDue := 1

	// Act
	user, err := InitializeUserWithLoan(db, "Historical User", "history@example.com", "555-5555",
		20000.0, 0.06, 24, dayDue, oneYearAgo, false)

	// Assert
	require.NoError(t, err, "Should create loan with historical date")
	require.Equal(t, oneYearAgo, user.Loans[0].DateTaken, "Date taken should match historical date")

	// Verify first payment is due on dayDue of the next month
	expectedFirstDue := time.Date(oneYearAgo.Year(), oneYearAgo.Month()+1, dayDue, 0, 0, 0, 0, time.UTC)
	actualFirstDue := user.Loans[0].Payments[0].DueDate

	require.Equal(t, expectedFirstDue.Year(), actualFirstDue.Year(), "Year should match")
	require.Equal(t, expectedFirstDue.Month(), actualFirstDue.Month(), "Month should match")
	require.Equal(t, dayDue, actualFirstDue.Day(), "Day should match dayDue parameter")

	t.Logf("✓ Successfully created loan with historical date from %s", oneYearAgo.Format("2006-01-02"))
}

// TestPaymentScheduleIntegrity verifies payment schedule handles month-end edge cases correctly.
func TestPaymentScheduleIntegrity(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange
	dateTaken := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	dayDue := 31

	// Act
	user, err := InitializeUserWithLoan(db, "Edge Case User", "edge@example.com", "555-6666",
		6000.0, 0.05, 6, dayDue, dateTaken, false)

	// Assert
	require.NoError(t, err, "Should create loan with edge case date")
	require.Len(t, user.Loans[0].Payments, 6, "Should have 6 payments")

	// Verify payment schedule handles month-end correctly
	payments := user.Loans[0].Payments

	require.Equal(t, 29, payments[0].DueDate.Day(), "Feb payment should be on 29th")
	require.Equal(t, 31, payments[1].DueDate.Day(), "March payment should be on 31st")
	require.Equal(t, 30, payments[2].DueDate.Day(), "April payment should be on 30th")
	require.Equal(t, 31, payments[3].DueDate.Day(), "May payment should be on 31st")

	t.Logf("✓ Payment schedule correctly handles month-end edge cases")
}

// TestZeroInterestLoan verifies calculation of loans with zero interest rate.
func TestZeroInterestLoan(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange - Create loan with 0% interest
	dateTaken := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	principal := 12000.0
	termMonths := 12

	// Act
	user, err := InitializeUserWithLoan(db, "Zero Interest User", "zero@example.com", "555-7777",
		principal, 0.0, termMonths, 15, dateTaken, false)

	// Assert
	require.NoError(t, err, "Should create zero interest loan")

	// Verify monthly payment is principal divided by months
	expectedMonthlyPayment := principal / float64(termMonths)
	actualMonthlyPayment := user.Loans[0].Payments[0].AmountDue
	require.InDelta(t, expectedMonthlyPayment, actualMonthlyPayment, 0.01,
		"Monthly payment should be principal/months for zero interest")

	// Verify total payments equal principal (no interest)
	totalPayments := 0.0
	for _, pmt := range user.Loans[0].Payments {
		totalPayments += pmt.AmountDue
	}
	require.InDelta(t, principal, totalPayments, 1.0,
		"Total payments should equal principal for zero interest loan")

	t.Logf("✓ Zero interest loan calculated correctly: $%.2f/month", actualMonthlyPayment)
}

// TestValidateLoanParameters verifies input validation.
func TestValidateLoanParameters(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	tests := []struct {
		name         string
		totalAmount  float64
		interestRate float64
		termMonths   int
		dayDue       int
		shouldFail   bool
	}{
		{"Valid parameters", 10000.0, 0.05, 12, 15, false},
		{"Zero total amount", 0.0, 0.05, 12, 15, true},
		{"Negative total amount", -1000.0, 0.05, 12, 15, true},
		{"Negative interest rate", 10000.0, -0.05, 12, 15, true},
		{"Zero term months", 10000.0, 0.05, 0, 15, true},
		{"Negative term months", 10000.0, 0.05, -12, 15, true},
		{"Day due too low", 10000.0, 0.05, 12, 0, true},
		{"Day due too high", 10000.0, 0.05, 12, 32, true},
		{"Zero interest valid", 10000.0, 0.0, 12, 15, false},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dateTaken := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
			// Use unique email for each test case to avoid constraint violation
			email := fmt.Sprintf("test%d@example.com", i)
			_, err := InitializeUserWithLoan(db, "Test User", email, "555-0000",
				tt.totalAmount, tt.interestRate, tt.termMonths, tt.dayDue, dateTaken, false)

			if tt.shouldFail {
				require.Error(t, err, "Should return validation error")
				t.Logf("✓ Correctly rejected: %s", tt.name)
			} else {
				require.NoError(t, err, "Should not return error for valid parameters")
				t.Logf("✓ Correctly accepted: %s", tt.name)
			}
		})
	}
}
