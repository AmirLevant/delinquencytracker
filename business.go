package delinquencytracker

import (
	"math"
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
	var mnthlyPayment float64
	var mnthlyInterestRate float64 = annualRate / 12

	mnthlyPayment = (principal * (mnthlyInterestRate * math.Pow(1+mnthlyInterestRate, float64(months)))) /
		(math.Pow(1+mnthlyInterestRate, float64(months)) - 1)

	return mnthlyPayment
}

// calculateDueDate calculates when a specific payment is due
// It adds 'paymentNum' months to the start date and sets the day to 'dayDue'
func calculateDueDate(startDate time.Time, termMonths, dayDue int) time.Time {

	// Add the number of months for this payment
	dueDate := startDate.AddDate(0, termMonths, 0)

	// Get the year and month
	year, month, _ := dueDate.Date()

	// Handle edge case: if dayDue is 31 but month only has 30 days
	// Go will automatically normalize (e.g., Feb 31 becomes Mar 3)
	// To be safe, we clamp to the last day of the month if needed
	lastDayOfMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
	if dayDue > lastDayOfMonth {
		dayDue = lastDayOfMonth
	}

	// Return the due date in UTC
	return time.Date(year, month, dayDue, 0, 0, 0, 0, time.UTC)
}

//func InitializeUserWithLoan(db *sql.DB, name, email, phone string,totalAmount, interestRate float64, termMonths, dayDue int) (user, error) {

//}
