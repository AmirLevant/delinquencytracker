package database

type Person struct {
	name   string
	email  string
	phone  string
	loanId string
	loan   Loan
}
