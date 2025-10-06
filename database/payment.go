package database

import (
	"time"
)

type Payment struct {
	latePayment bool // is this payment late or incomplete based on paymentDueDate and paymentActualDate

	loanId string // which loan are we associated to

	paymentNumber int64 // Sequential counter, which payment is it (1st, 2nd etc...)

	paymentDue float64 // the amount of money owed in the payment

	paymentPaid float64 // the amount of money that was actually paid in this payment

	paymentDueDate time.Time // when was the payment due

	paymentPaidDate time.Time // when was the payment actually paid

	paymentLate bool // is this specific payment late or not

	// lateSettled bool // if the payment is late, was it settled?

}

func (p *Payment) LateCheck() {
	if(p.paymentDueDate)
}
