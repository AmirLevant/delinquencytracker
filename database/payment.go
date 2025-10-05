package database

type payment struct {
	paymentnumber int64 // which payment is it

	totalLoan int64 // total amount of money taken in the loan

	duedate string // when was the payment due

	paiddate string // when was the payment paid

	latepayment bool // is this payment late / not fully fullfilled

}
