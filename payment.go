package delinquencytracker

import "time"

type payment struct {
	ID            int64      // unique identifier for the payment
	LoanID        int64      // which loan is this payment for
	PaymentNumber int        // sequential counter (1st, 2nd, 3rd payment, etc.)
	AmountDue     float64    // how much money is owed in this payment
	AmountPaid    float64    // how much money was actually paid
	DueDate       time.Time  // when is this payment due
	PaidDate      *time.Time // when was this payment actually made (nil if unpaid)
	CreatedAt     time.Time  // when was this record created
}
