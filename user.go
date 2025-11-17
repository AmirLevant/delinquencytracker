package delinquencytracker

import "time"

type User struct {
	ID        int64     // unique identifier for the user
	Name      string    // full name of the user
	Email     string    // email address
	Phone     string    // phone number
	CreatedAt time.Time // when the user was created

	Loans []Loan // all loans associated with this user
}
