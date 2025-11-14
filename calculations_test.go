package delinquencytracker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// ============================================================================
// CURRENT STATUS TESTS
// ============================================================================

func TestPaymentIsOverdue(t *testing.T) {
	// Test 1: overdue unpaid payment
	payment1 := payment{
		DueDate:    time.Now().UTC().AddDate(0, 0, -5), // 5 days ago
		AmountDue:  100.0,
		AmountPaid: 0.0,
	}
	result1 := payment1.IsOverdue()
	require.True(t, result1, "Payment due 5 days ago with no payment should be overdue")

	// Test 2: overdue partially paid
	payment2 := payment{
		DueDate:    time.Now().UTC().AddDate(0, 0, -3),
		AmountDue:  100.0,
		AmountPaid: 50.0,
	}
	result2 := payment2.IsOverdue()
	require.True(t, result2, "Partially paid overdue payment should still be overdue")

	// Test 3: overdue but fully paid
	payment3 := payment{
		DueDate:    time.Now().UTC().AddDate(0, 0, -10),
		AmountDue:  100.0,
		AmountPaid: 100.0,
	}
	result3 := payment3.IsOverdue()
	require.False(t, result3, "Fully paid payment should not be overdue even if past due date")

	// Test 4: not yet due
	payment4 := payment{
		DueDate:    time.Now().UTC().AddDate(0, 0, 5), // 5 days from now
		AmountDue:  100.0,
		AmountPaid: 0.0,
	}
	result4 := payment4.IsOverdue()
	require.False(t, result4, "Future payment should not be overdue")

	// Test 5: overpaid
	payment5 := payment{
		DueDate:    time.Now().UTC().AddDate(0, 0, -5),
		AmountDue:  100.0,
		AmountPaid: 150.0,
	}
	result5 := payment5.IsOverdue()
	require.False(t, result5, "Overpaid payment should not be overdue")
}

func TestPaymentDaysOverdue(t *testing.T) {
	// Test 1: 5 days overdue
	payment1 := payment{
		DueDate:    time.Now().UTC().AddDate(0, 0, -5),
		AmountDue:  100.0,
		AmountPaid: 0.0,
	}
	result1 := payment1.DaysOverdue()
	require.Equal(t, 5, result1, "Expected 5 days overdue")

	// Test 2: 10 days overdue partially paid
	payment2 := payment{
		DueDate:    time.Now().UTC().AddDate(0, 0, -10),
		AmountDue:  100.0,
		AmountPaid: 30.0,
	}
	result2 := payment2.DaysOverdue()
	require.Equal(t, 10, result2, "Expected 10 days overdue")

	// Test 3: not overdue returns 0
	payment3 := payment{
		DueDate:    time.Now().UTC().AddDate(0, 0, 5),
		AmountDue:  100.0,
		AmountPaid: 0.0,
	}
	result3 := payment3.DaysOverdue()
	require.Equal(t, 0, result3, "Expected 0 days overdue for future payment")

	// Test 4: fully paid returns 0
	payment4 := payment{
		DueDate:    time.Now().UTC().AddDate(0, 0, -10),
		AmountDue:  100.0,
		AmountPaid: 100.0,
	}
	result4 := payment4.DaysOverdue()
	require.Equal(t, 0, result4, "Expected 0 days overdue for fully paid payment")
}

func TestPaymentIsFullyPaid(t *testing.T) {
	// Test 1: fully paid exact amount
	payment1 := payment{
		AmountDue:  100.0,
		AmountPaid: 100.0,
	}
	result1 := payment1.IsFullyPaid()
	require.True(t, result1, "Expected payment to be fully paid")

	// Test 2: overpaid
	payment2 := payment{
		AmountDue:  100.0,
		AmountPaid: 150.0,
	}
	result2 := payment2.IsFullyPaid()
	require.True(t, result2, "Expected overpaid payment to be considered fully paid")

	// Test 3: partially paid
	payment3 := payment{
		AmountDue:  100.0,
		AmountPaid: 50.0,
	}
	result3 := payment3.IsFullyPaid()
	require.False(t, result3, "Expected partially paid payment to not be fully paid")

	// Test 4: unpaid
	payment4 := payment{
		AmountDue:  100.0,
		AmountPaid: 0.0,
	}
	result4 := payment4.IsFullyPaid()
	require.False(t, result4, "Expected unpaid payment to not be fully paid")
}

func TestPaymentRemainingBalance(t *testing.T) {
	// Test 1: no payment made
	payment1 := payment{
		AmountDue:  100.0,
		AmountPaid: 0.0,
	}
	result1 := payment1.RemainingBalance()
	require.Equal(t, 100.0, result1, "Expected remaining balance to be 100.0")

	// Test 2: partial payment
	payment2 := payment{
		AmountDue:  100.0,
		AmountPaid: 30.0,
	}
	result2 := payment2.RemainingBalance()
	require.Equal(t, 70.0, result2, "Expected remaining balance to be 70.0")

	// Test 3: fully paid
	payment3 := payment{
		AmountDue:  100.0,
		AmountPaid: 100.0,
	}
	result3 := payment3.RemainingBalance()
	require.Equal(t, 0.0, result3, "Expected remaining balance to be 0.0")

	// Test 4: overpaid returns 0
	payment4 := payment{
		AmountDue:  100.0,
		AmountPaid: 150.0,
	}
	result4 := payment4.RemainingBalance()
	require.Equal(t, 0.0, result4, "Expected remaining balance to be 0.0 for overpayment")

	// Test 5: almost fully paid
	payment5 := payment{
		AmountDue:  100.0,
		AmountPaid: 99.99,
	}
	result5 := payment5.RemainingBalance()
	require.InDelta(t, 0.01, result5, 0.001, "Expected remaining balance to be approximately 0.01")
}

func TestPaymentIsPartiallyPaid(t *testing.T) {
	// Test 1: partial payment
	payment1 := payment{
		AmountDue:  100.0,
		AmountPaid: 50.0,
	}
	result1 := payment1.IsPartiallyPaid()
	require.True(t, result1, "Expected payment to be partially paid")

	// Test 2: small partial payment
	payment2 := payment{
		AmountDue:  100.0,
		AmountPaid: 0.01,
	}
	result2 := payment2.IsPartiallyPaid()
	require.True(t, result2, "Expected small payment to be considered partially paid")

	// Test 3: fully paid not partial
	payment3 := payment{
		AmountDue:  100.0,
		AmountPaid: 100.0,
	}
	result3 := payment3.IsPartiallyPaid()
	require.False(t, result3, "Expected fully paid payment to not be partially paid")

	// Test 4: unpaid not partial
	payment4 := payment{
		AmountDue:  100.0,
		AmountPaid: 0.0,
	}
	result4 := payment4.IsPartiallyPaid()
	require.False(t, result4, "Expected unpaid payment to not be partially paid")

	// Test 5: overpaid not partial
	payment5 := payment{
		AmountDue:  100.0,
		AmountPaid: 150.0,
	}
	result5 := payment5.IsPartiallyPaid()
	require.False(t, result5, "Expected overpaid payment to not be partially paid")
}

func TestPaymentIsPaid(t *testing.T) {
	// Test 1: fully paid
	payment1 := payment{
		AmountDue:  100.0,
		AmountPaid: 100.0,
	}
	result1 := payment1.IsPaid()
	require.True(t, result1, "Expected fully paid payment to be paid")

	// Test 2: partially paid
	payment2 := payment{
		AmountDue:  100.0,
		AmountPaid: 50.0,
	}
	result2 := payment2.IsPaid()
	require.True(t, result2, "Expected partially paid payment to be paid")

	// Test 3: overpaid
	payment3 := payment{
		AmountDue:  100.0,
		AmountPaid: 150.0,
	}
	result3 := payment3.IsPaid()
	require.True(t, result3, "Expected overpaid payment to be paid")

	// Test 4: unpaid
	payment4 := payment{
		AmountDue:  100.0,
		AmountPaid: 0.0,
	}
	result4 := payment4.IsPaid()
	require.False(t, result4, "Expected unpaid payment to not be paid")
}

// ============================================================================
// HISTORICAL ANALYSIS TESTS
// ============================================================================

func TestPaymentWasPaidLate(t *testing.T) {
	now := time.Now().UTC()

	// Test 1: paid on time
	payment1 := payment{
		DueDate:  now.AddDate(0, 0, -10),
		PaidDate: now.AddDate(0, 0, -11), // Paid 1 day before due
	}
	result1 := payment1.WasPaidLate()
	require.False(t, result1, "Expected payment paid before due date to not be late")

	// Test 2: paid exactly on due date
	payment2 := payment{
		DueDate:  now.AddDate(0, 0, -10),
		PaidDate: now.AddDate(0, 0, -10),
	}
	result2 := payment2.WasPaidLate()
	require.False(t, result2, "Expected payment paid on due date to not be late")

	// Test 3: paid 1 day late
	payment3 := payment{
		DueDate:  now.AddDate(0, 0, -10),
		PaidDate: now.AddDate(0, 0, -9), // Paid 1 day after due
	}
	result3 := payment3.WasPaidLate()
	require.True(t, result3, "Expected payment paid 1 day late to be late")

	// Test 4: paid 30 days late
	payment4 := payment{
		DueDate:  now.AddDate(0, 0, -40),
		PaidDate: now.AddDate(0, 0, -10),
	}
	result4 := payment4.WasPaidLate()
	require.True(t, result4, "Expected payment paid 30 days late to be late")

	// Test 5: not yet paid
	payment5 := payment{
		DueDate:  now.AddDate(0, 0, -10),
		PaidDate: time.Time{}, // Zero time
	}
	result5 := payment5.WasPaidLate()
	require.False(t, result5, "Expected unpaid payment to not be considered late")
}

func TestPaymentDaysLate(t *testing.T) {
	now := time.Now().UTC()

	// Test 1: paid 5 days late
	payment1 := payment{
		DueDate:  now.AddDate(0, 0, -15),
		PaidDate: now.AddDate(0, 0, -10),
	}
	result1 := payment1.DaysLate()
	require.Equal(t, 5, result1, "Expected 5 days late")

	// Test 2: paid 30 days late
	payment2 := payment{
		DueDate:  now.AddDate(0, 0, -40),
		PaidDate: now.AddDate(0, 0, -10),
	}
	result2 := payment2.DaysLate()
	require.Equal(t, 30, result2, "Expected 30 days late")

	// Test 3: paid on time returns 0
	payment3 := payment{
		DueDate:  now.AddDate(0, 0, -10),
		PaidDate: now.AddDate(0, 0, -11),
	}
	result3 := payment3.DaysLate()
	require.Equal(t, 0, result3, "Expected 0 days late for on-time payment")

	// Test 4: not yet paid returns 0
	payment4 := payment{
		DueDate:  now.AddDate(0, 0, -10),
		PaidDate: time.Time{},
	}
	result4 := payment4.DaysLate()
	require.Equal(t, 0, result4, "Expected 0 days late for unpaid payment")
}
