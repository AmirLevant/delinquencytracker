package delinquencytracker

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMakeMonthlyPayment(t *testing.T) {

	//arrange
	ln := loan{}
	// pyment, err := CreatePayment(db, ln.ID, 1, 1000, 900, dueDate, paidDate)
	ln.TotalAmount = 10000
	ln.InterestRate = 0.05
	ln.TermMonths = 36

	//act
	mnthlypaymentAmount := calculateMonthlyPayment(ln.TotalAmount, ln.InterestRate, ln.TermMonths)

	fmt.Printf("The monthly payment is: %v", mnthlypaymentAmount)
	//assert
	require.NotEqual(t, nil, mnthlypaymentAmount)

}
