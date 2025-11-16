package delinquencytracker

import (
	"database/sql"
	"time"
)

// calculateMonthlyPayment calculates the monthly payment using the amortization formula
// Formula: M = P * [r(1+r)^n] / [(1+r)^n - 1]
// Where:
//
//	P = principal (total amount borrowed)
//	r = monthly interest rate (annual rate / 12)
//	n = number of payments (term in months)
func calculateMonthlyPayment(principal, annualRate float64, months int) float64 {

}

// calculateDueDate calculates when a specific payment is due
// It adds 'paymentNum' months to the start date and sets the day to 'dayDue'
func calculateDueDate(startDate time.Time, paymentNum, dayDue int) time.Time {

}

func InitializeUserWithLoan(db *sql.DB, name, email, phone string,
	totalAmount, interestRate float64, termMonths, dayDue int) (user, error) {

}
