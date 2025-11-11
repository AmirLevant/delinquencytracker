package delinquencytracker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Test to verify the CheckLate method for the payment struct works apropriately
func TestPaymentCheckLate(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	dateTaken := time.Now().UTC().Truncate(24 * time.Hour)

	// Arrange

	usr, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Creating a loan for the test user
	ln, err := CreateLoan(db, usr.ID, 10000.00, 0.05, 16, 05, "active", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	// Set up payment dates
	dueDate := dateTaken.Add(30 * 24 * time.Hour) // 30 days after loan was taken
	paidDate := dueDate.Add(-2 * 24 * time.Hour)  // paid 2 days before due date

	pyment, err := CreatePayment(db, ln.ID, 1, 1000, 900, dueDate, paidDate)
	if err != nil {
		t.Fatalf("Create Payment failed %v:", err)
	}

	CheckLateExpected := false // payment was paid before due date, hence it is not late
	CheckLateActual := pyment.CheckLate()

	require.Equal(t, CheckLateExpected, CheckLateActual)
}
