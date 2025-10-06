package database

type Loan struct {
	id        string  // what is the id corresponding with the loan
	totalLoan float64 // what is the total amount of money borrowed
	dateTaken string  // what date was the loan taken
	dayDue    string  // what day of the month is the payment due

	payments []Payment
}
