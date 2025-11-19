package delinquencytracker

import "time"

type Loan struct {
	ID           int64     // unique identifier for the loan
	UserID       int64     // which user this loan belong to
	TotalAmount  float64   // total amount of money borrowed
	InterestRate float64   // annual interest rate (0.05 for 5% etc...)
	TermMonths   int       // how many months is the loan term
	DayDue       int       // what day of the month is payment due (1-31)
	Status       string    // current status: "active", "paid_off", "defaulted"
	DateTaken    time.Time // when was the loan taken
	CreatedAt    time.Time // when was this record created

	Payments []Payment // all payments associated with this loan

}
