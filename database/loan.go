package database

type Loan struct {
	id        string
	totalLoan float64
	datetaken string

	payments []payment
}
