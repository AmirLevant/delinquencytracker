package delinquencytracker

import "time"

// IsOverdue checks if a payment is past its due date and not fully paid
func (p *payment) IsOverdue() bool {
	now := time.Now().UTC()
	return now.After(p.DueDate) && !p.IsFullyPaid()
}

// DaysOverdue calculates how many days past the due date this payment is
// Returns 0 if not overdue
func (p *payment) DaysOverdue() int {
	if !p.IsOverdue() {
		return 0
	}
	now := time.Now().UTC()
	duration := now.Sub(p.DueDate)
	return int(duration.Hours() / 24)
}

// IsFullyPaid checks if the payment has been paid in full
func (p *payment) IsFullyPaid() bool {
	return p.AmountPaid >= p.AmountDue
}

// RemainingBalance returns how much is still owed on this payment
func (p *payment) RemainingBalance() float64 {
	remaining := p.AmountDue - p.AmountPaid
	if remaining < 0 {
		return 0
	}
	return remaining
}

// IsPartiallyPaid checks if some payment has been made but not the full amount
func (p *payment) IsPartiallyPaid() bool {
	return p.AmountPaid > 0 && p.AmountPaid < p.AmountDue
}

// IsPaid checks if any payment has been recorded (even partial)
func (p *payment) IsPaid() bool {
	return p.AmountPaid > 0
}
