package delinquencytracker

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// sets up the test database connection
func setupTestDB(t *testing.T) *sql.DB {
	config := "host=localhost port=5432 user=postgres password=amir dbname=loan_tracker sslmode=disable"
	db, err := sql.Open("postgres", config)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	return db
}

// cleanup
func teardownTestDB(db *sql.DB) {
	db.Exec("DELETE FROM payments")
	db.Exec("DELETE FROM loans")
	db.Exec("DELETE FROM users")
	db.Close()
}

// 20/10/25, test will fail since GetUserByID does not exist yet
// 21/10/25 test will pass since GetUserByID exists now
func TestGetUserByID(t *testing.T) {

	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange, creating a test user
	usr, err := CreateUser(db, "Test User", "test@test.com", "555-4444")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	//Act: Get the user by ID
	usr, err = GetUserByID(db, usr.ID)

	//Assert: Check results
	if err != nil {
		t.Errorf("GetUserByID failed: %v", err)
	}
	if usr.Name != "Test User" {
		t.Fatalf("Expected name 'Test User', got '%s'", usr.Name)
	}
	if usr.Email != "test@test.com" {
		t.Fatalf("Expected email 'test@test.com', got '%s'", usr.Email)
	}
}

// we provide the GetUserByID an ID of a user that does not exist
// we expect the GetUserByID to fail and return
func TestGetUserByID_UserNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Act, Trying to get a user that does not exist

	usr, err := GetUserByID(db, 99999)

	//Assert, Should return error
	assert.Error(t, err, "Expected error for non-existent user")
	require.Equal(t, user{}, usr, "Expected empty user struct")
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange, Creating the first user
	_, err := CreateUser(db, "User One", "duplicate@test.com", "555-0001")
	if err != nil {
		t.Fatalf("Failed to create first user: %v", err)
	}

	//Act, Creating another user with the same email
	_, err = CreateUser(db, "User Two", "duplicate@test.com", "555-0002")

	// Assert, Should return Error
	if err == nil {
		t.Error("Expected error for duplicate email, got nil")
	}
}

// 21/10/25, Creating the Test for updating a user
//
//	The function of updating user still does not exist
//	As expected the test fails

// 21/10/25 I created UpdateUser in logic file
// 21/10/25 test passed as expected
// 22/10/25 updated my CRUD operations to be in db.go
// 23/10/25 need to update the return type of certain functions, tests need to change accordingly
// 23/10/25 updated the test to take the id from the returned user of CreateUser()

func TestUpdateUser(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange
	// creating user
	usr, _ := CreateUser(db, "Old Name", "old@test.com", "555-0000")

	// Act
	// updating user
	err := UpdateUser(db, usr.ID, "New Name", "new@test.com", "555-9999")

	// Assert
	// update should succeed
	if err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}

	// ensuring the update worked by calling the ID
	updatedUsr, err := GetUserByID(db, usr.ID)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}

	// esnuring the updated user is as we expected
	if updatedUsr.Name != "New Name" || updatedUsr.Email != "new@test.com" || updatedUsr.Phone != "555-9999" {
		t.Fatalf("Updated User is not as Expected")
	}
}

// 22/10/25 created the test in order to create GetAllUsers
// 22/10/25 realized that due to the composite struct definition i cannot compare users
// 23/10/25 solution, use "testify" a testing oriented module that allows deep testing
func TestGetAllUsers(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Creating multiple test users
	// Arrange
	user1, err := CreateUser(db, "Amir M", "amir@example.com", "111")

	if err != nil {
		t.Fatalf("Failed to create user1: %v", err)
	}

	user2, err := CreateUser(db, "Ori J", "ori@example.com", "333")

	if err != nil {
		t.Fatalf("Failed to create user2: %v", err)
	}

	user3, err := CreateUser(db, "Seb I", "seb@example.com", "222")

	if err != nil {
		t.Fatalf("Failed to create user3: %v", err)
	}

	// Act
	actualusers, err := GetAllUsers(db)
	expectedusers := []user{user1, user2, user3}

	// Assert
	if err != nil {
		t.Errorf("Error is not nil %v", err)
	}
	require.Equal(t, expectedusers, actualusers, "Expected both user slices to be equal")

}

// 23/10/25 create test before Delete User
// 23/10/25 expected that the user gets created then deleted
// 23/10/25 expected that the GetUserByID() will return an empty user
// 23/10/25 expected
func TestDeleteUser(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange, creating the user

	usr, err := CreateUser(db, "Deleted Usr", "deleted@example.com", "555")

	if err != nil {
		t.Fatal("err is not nil in CreateUser %w", err)
	}

	// Act
	// deleting the user
	err = DeleteUser(db, usr.ID)

	if err != nil {
		t.Fatal("err is not nil in DeleteUser but %w", err)
	}

	// verifying that such a user does not exist
	// function should return an empty user and NOT nil
	deletedUsr, err := GetUserByID(db, usr.ID)

	if err == nil {
		t.Fatalf("Error should not be nil, but a message saying user not found")
	}
	require.Equal(t, user{}, deletedUsr, "User has been deleted and not found")

}

// 23/10/25 create test for CreateLoan
// expecting to create a loan for a user and verify all fields are set correctly
func TestCreateLoan(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange, creating a test user first
	usr, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Act, creating a loan for this user
	dateTaken := time.Now()
	ln, err := CreateLoan(db, usr.ID, 10000.00, 0.05, 36, 15, "active", dateTaken)

	// Assert, loan creation should succeed
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}
	if ln.ID == 0 {
		t.Error("Expected loan ID to be set")
	}
	if ln.UserID != usr.ID {
		t.Errorf("Expected UserID %d, got %d", usr.ID, ln.UserID)
	}
	if ln.TotalAmount != 10000.00 {
		t.Errorf("Expected TotalAmount 10000.00, got %f", ln.TotalAmount)
	}
	if ln.InterestRate != 0.05 {
		t.Errorf("Expected InterestRate 0.05, got %f", ln.InterestRate)
	}
	if ln.TermMonths != 36 {
		t.Errorf("Expected TermMonths 36, got %d", ln.TermMonths)
	}
	if ln.DayDue != 15 {
		t.Errorf("Expected DayDue 15, got %d", ln.DayDue)
	}
	if ln.Status != "active" {
		t.Errorf("Expected Status 'active', got '%s'", ln.Status)
	}
	if ln.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
}

func TestUpdateLoan(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange
	// Creating a test user first
	usr, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Creating a loan for this user
	dateTaken := time.Now().UTC().Truncate(24 * time.Hour)
	ln, err := CreateLoan(db, usr.ID, 10000.00, 0.05, 36, 15, "active", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	// Act
	// Updating the loan with new values
	newDateTaken := time.Now().UTC().Truncate(24*time.Hour).AddDate(0, 0, -30) // 30 days ago
	err = UpdateLoan(db, ln.ID, 15000.00, 0.08, 48, 20, "refinanced", newDateTaken)

	// Assert
	// Update should succeed
	if err != nil {
		t.Fatalf("UpdateLoan failed: %v", err)
	}

	// Ensuring the update worked by querying the loans
	loans, err := GetLoansByUserID(db, usr.ID)
	if err != nil {
		t.Fatalf("GetLoansByUserID failed: %v", err)
	}

	// Should have exactly one loan
	if len(loans) != 1 {
		t.Fatalf("Expected 1 loan, got %d", len(loans))
	}

	updatedLoan := loans[0]

	// Ensuring the updated loan has the new values
	if updatedLoan.TotalAmount != 15000.00 {
		t.Errorf("Expected TotalAmount 15000.00, got %f", updatedLoan.TotalAmount)
	}
	if updatedLoan.InterestRate != 0.08 {
		t.Errorf("Expected InterestRate 0.08, got %f", updatedLoan.InterestRate)
	}
	if updatedLoan.TermMonths != 48 {
		t.Errorf("Expected TermMonths 48, got %d", updatedLoan.TermMonths)
	}
	if updatedLoan.DayDue != 20 {
		t.Errorf("Expected DayDue 20, got %d", updatedLoan.DayDue)
	}
	if updatedLoan.Status != "refinanced" {
		t.Errorf("Expected Status 'refinanced', got '%s'", updatedLoan.Status)
	}
	// Verify DateTaken was updated (comparing truncated dates)
	if !updatedLoan.DateTaken.Equal(newDateTaken) {
		t.Errorf("Expected DateTaken %v, got %v", newDateTaken, updatedLoan.DateTaken)
	}
}

func TestGetLoanByID(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	dateTaken := time.Now().UTC().Truncate(24 * time.Hour)

	// Arrange, creating a test user first
	usr, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Creating a loan for the test user
	createdLoan, err := CreateLoan(db, usr.ID, 10000.00, 0.05, 36, 15, "active", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	// Act
	retrievedLoan, err := GetLoanByID(db, createdLoan.ID)

	// Assert
	if err != nil {
		t.Fatalf("GetLoanByID failed: %v", err)
	}

	// Verify all fields match
	require.Equal(t, createdLoan, retrievedLoan, "Retrieved loan should match created loan")
}

func TestGetLoanByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Act, Trying to get a loan that does not exist
	ln, err := GetLoanByID(db, 99999)

	// Assert, Should return error
	assert.Error(t, err, "Expected error for non-existent loan")
	require.Equal(t, loan{}, ln, "Expected empty loan struct")
}

func TestGetLoansByUserID_OneLoan(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange, creating a test user first
	usr, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Act, creating a loan for this user
	dateTaken := time.Now()
	expectedln, err := CreateLoan(db, usr.ID, 10000.00, 0.05, 36, 15, "active", dateTaken)

	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	// Act, we query all the loans that belong to the userID
	loans, err := GetLoansByUserID(db, usr.ID)

	actualLn := loans[0]

	// Assert, loan creation should succeed
	if err != nil {
		t.Fatalf("GetLoansByUserID failed: %v", err)
	}
	if expectedln.ID != actualLn.ID {
		t.Error("Expected loan ID to match")
	}
	if expectedln.UserID != actualLn.UserID {
		t.Errorf("Expected UserID %d, got %d", expectedln.UserID, actualLn.UserID)
	}
	if expectedln.TotalAmount != actualLn.TotalAmount {
		t.Errorf("Expected TotalAmount 10000.00, got %f", actualLn.TotalAmount)
	}
	if expectedln.InterestRate != actualLn.InterestRate {
		t.Errorf("Expected InterestRate 0.05, got %f", actualLn.InterestRate)
	}
	if expectedln.TermMonths != actualLn.TermMonths {
		t.Errorf("Expected TermMonths 36, got %d", actualLn.TermMonths)
	}
	if expectedln.DayDue != actualLn.DayDue {
		t.Errorf("Expected DayDue 15, got %d", actualLn.DayDue)
	}
	if expectedln.Status != actualLn.Status {
		t.Errorf("Expected Status 'active', got '%s'", actualLn.Status)
	}
	if actualLn.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

}

func TestGetLoansByUserID_MultiLoan(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange, creating a test user first
	usr, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Act, creating a loan for this user
	dateTaken := time.Now().UTC().Truncate(24 * time.Hour)

	expectedln1, err := CreateLoan(db, usr.ID, 10000.00, 0.05, 16, 05, "active", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	expectedln2, err := CreateLoan(db, usr.ID, 20000.00, 0.25, 26, 15, "paid_off", dateTaken)

	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	expectedln3, err := CreateLoan(db, usr.ID, 30000.00, 0.35, 36, 25, "defaulted", dateTaken)

	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	// Act, we query all the loans that belong to the userID
	actualLoans, err := GetLoansByUserID(db, usr.ID)

	if err != nil {
		t.Fatalf("GetLoansByUserID failed: %v", err)
	}

	var expectedLoans = []loan{expectedln1, expectedln2, expectedln3}

	require.Equal(t, expectedLoans, actualLoans)

}

func TestGetLoansByUserID_NoLoan(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange, creating a test user first
	usr, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Act, no loan for this user

	actualLoans, err := GetLoansByUserID(db, usr.ID)

	if err != nil {
		t.Fatalf("GetLoansByUserID failed: %v", err)
	}

	// when comparing actualLoans to an expectedLoans there is an issue since GetLoansByUserID()
	// initializes a slice with nils, which is why it is different
	require.Empty(t, actualLoans)

}

func TestGetAllLoans(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	dateTaken := time.Now().UTC().Truncate(24 * time.Hour)

	// Arrange, creating a multiple test users
	usr1, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test user1: %v", err)
	}

	usr2, err := CreateUser(db, "Test User", "loanuser2@test.com", "555-2222")
	if err != nil {
		t.Fatalf("Failed to create test user2: %v", err)
	}

	usr3, err := CreateUser(db, "User Third", "loanuser3@test.com", "555-3333")
	if err != nil {
		t.Fatalf("Failed to create test user3: %v", err)
	}

	expectedln1, err := CreateLoan(db, usr1.ID, 10000.00, 0.05, 16, 05, "active", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	expectedln2, err := CreateLoan(db, usr2.ID, 20000.00, 0.25, 26, 15, "paid_off", dateTaken)

	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	expectedln3, err := CreateLoan(db, usr3.ID, 30000.00, 0.35, 36, 25, "defaulted", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	var expectedLoans = []loan{expectedln1, expectedln2, expectedln3}

	actualLoans, err := GetAllLoans(db)

	if err != nil {
		t.Fatalf("GetAllLoans failed: %v", err)
	}

	require.Equal(t, expectedLoans, actualLoans)

}

func TestDeleteLoan(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	dateTaken := time.Now().UTC().Truncate(24 * time.Hour)

	// Arrange, creating a test user
	usr1, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test user1: %v", err)
	}

	// Creating a loan for the test user
	expectedln1, err := CreateLoan(db, usr1.ID, 10000.00, 0.05, 16, 05, "active", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	err = DeleteLoan(db, expectedln1.ID)
	if err != nil {
		t.Fatalf("DeleteLoan failed: %v", err)
	}

	checkLn, err := GetLoansByUserID(db, usr1.ID)
	if err != nil {
		t.Fatalf("GetLoansByUserID failed: %v", err)
	}

	require.Empty(t, checkLn)

}

func TestCreatePayment(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)
	dateTaken := time.Now().UTC().Truncate(24 * time.Hour)

	// Arrange

	usr, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Creating a loan for the test user
	ln, err := CreateLoan(db, usr.ID, 10000.00, 0.05, 16, 05, "active", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	// Set up payment dates
	dueDate := dateTaken.Add(30 * 24 * time.Hour) // 30 days after loan was taken
	paidDate := dueDate.Add(-2 * 24 * time.Hour)  // paid 2 days before due date

	pyment, err := CreatePayment(db, ln.ID, 1, 1000, 900, dueDate, paidDate)
	if err != nil {
		t.Fatalf("Create Payment failed %v:", err)
	}

	var expectedPyment = payment{pyment.ID, ln.ID, 1, 1000, 900, dueDate, paidDate, pyment.CreatedAt}

	require.Equal(t, expectedPyment, pyment)

}

func TestUpdatePayment(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)
	dateTaken := time.Now().UTC().Truncate(24 * time.Hour)

	// Arrange

	usr, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Creating a loan for the test user
	ln, err := CreateLoan(db, usr.ID, 10000.00, 0.05, 16, 05, "active", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	// Set up payment dates
	dueDate := dateTaken.Add(30 * 24 * time.Hour) // 30 days after loan was taken
	paidDate := dueDate.Add(-2 * 24 * time.Hour)  // paid 2 days before due date

	pyment, err := CreatePayment(db, ln.ID, 1, 1000, 900, dueDate, paidDate)
	if err != nil {
		t.Fatalf("Create Payment failed %v:", err)
	}

	// Act
	// Updating the payment with new values
	newDueDate := dateTaken.Add(45 * 24 * time.Hour)  // 45 days after loan was taken
	newPaidDate := newDueDate.Add(3 * 24 * time.Hour) // paid 3 days late

	err = UpdatePayment(db, pyment.ID, ln.ID, 2, 1200.00, 1200.00, newDueDate, newPaidDate)

	// Assert
	// Update should succeed
	if err != nil {
		t.Fatalf("UpdatePayment failed: %v", err)
	}

	// Ensuring the update worked by querying the payment
	updatedPayment, err := GetPaymentByID(db, pyment.ID)
	if err != nil {
		t.Fatalf("GetPaymentByID failed: %v", err)
	}

	// Ensuring the updated payment has the new values
	if updatedPayment.LoanID != ln.ID {
		t.Errorf("Expected LoanID %d, got %d", ln.ID, updatedPayment.LoanID)
	}
	if updatedPayment.PaymentNumber != 2 {
		t.Errorf("Expected PaymentNumber 2, got %d", updatedPayment.PaymentNumber)
	}
	if updatedPayment.AmountDue != 1200.00 {
		t.Errorf("Expected AmountDue 1200.00, got %f", updatedPayment.AmountDue)
	}
	if updatedPayment.AmountPaid != 1200.00 {
		t.Errorf("Expected AmountPaid 1200.00, got %f", updatedPayment.AmountPaid)
	}

	// Verify dates were updated (comparing truncated dates)
	if !updatedPayment.DueDate.Equal(newDueDate) {
		t.Errorf("Expected DueDate %v, got %v", newDueDate, updatedPayment.DueDate)
	}
	if !updatedPayment.PaidDate.Equal(newPaidDate) {
		t.Errorf("Expected PaidDate %v, got %v", newPaidDate, updatedPayment.PaidDate)
	}

}

func TestGetPaymentByID(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)
	dateTaken := time.Now().UTC().Truncate(24 * time.Hour)

	// Arrange
	usr, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Creating a loan for the test user
	ln, err := CreateLoan(db, usr.ID, 10000.00, 0.05, 16, 05, "active", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	// Set up payment dates
	dueDate := dateTaken.Add(30 * 24 * time.Hour) // 30 days after loan was taken
	paidDate := dueDate.Add(-2 * 24 * time.Hour)  // paid 2 days before due date

	// Create a payment to retrieve
	createdPayment, err := CreatePayment(db, ln.ID, 1, 1000.00, 900.00, dueDate, paidDate)
	if err != nil {
		t.Fatalf("CreatePayment failed: %v", err)
	}

	// Act
	retrievedPayment, err := GetPaymentByID(db, createdPayment.ID)

	// Assert
	if err != nil {
		t.Fatalf("GetPaymentByID failed: %v", err)
	}

	// Verify all fields match
	require.Equal(t, createdPayment, retrievedPayment, "Retrieved payment should match created payment")
}

func TestGetPaymentByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Act
	_, err := GetPaymentByID(db, 99999) // Non-existent ID

	// Assert
	if err == nil {
		t.Fatal("Expected error for non-existent payment, got nil")
	}
}

func TestGetPaymentsByLoanID_SinglePayment(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)
	dateTaken := time.Now().UTC().Truncate(24 * time.Hour)

	// Arrange
	usr, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Creating a loan for the test user
	ln, err := CreateLoan(db, usr.ID, 10000.00, 0.05, 16, 05, "active", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	// Set up payment dates
	dueDate := dateTaken.Add(30 * 24 * time.Hour) // 30 days after loan was taken
	paidDate := dueDate.Add(-2 * 24 * time.Hour)  // paid 2 days before due date

	expectedPayment, err := CreatePayment(db, ln.ID, 1, 1000.00, 900.00, dueDate, paidDate)
	if err != nil {
		t.Fatalf("CreatePayment failed: %v", err)
	}

	// Act
	payments, err := GetPaymentsByLoanID(db, ln.ID)

	// Assert
	if err != nil {
		t.Fatalf("GetPaymentsByLoanID failed: %v", err)
	}

	require.Len(t, payments, 1, "Should have exactly one payment")
	require.Equal(t, expectedPayment, payments[0], "Payment should match created payment")
}

func TestGetPaymentsByLoanID_MultiplePayments(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)
	dateTaken := time.Now().UTC().Truncate(24 * time.Hour)

	// Arrange
	usr, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Creating a loan for the test user
	ln, err := CreateLoan(db, usr.ID, 10000.00, 0.05, 36, 15, "active", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	// Create multiple payments
	dueDate1 := dateTaken.Add(30 * 24 * time.Hour)
	paidDate1 := dueDate1.Add(-2 * 24 * time.Hour)
	expectedPayment1, err := CreatePayment(db, ln.ID, 1, 300.00, 300.00, dueDate1, paidDate1)
	if err != nil {
		t.Fatalf("CreatePayment 1 failed: %v", err)
	}

	dueDate2 := dateTaken.Add(60 * 24 * time.Hour)
	paidDate2 := dueDate2.Add(-1 * 24 * time.Hour)
	expectedPayment2, err := CreatePayment(db, ln.ID, 2, 300.00, 295.00, dueDate2, paidDate2)
	if err != nil {
		t.Fatalf("CreatePayment 2 failed: %v", err)
	}

	dueDate3 := dateTaken.Add(90 * 24 * time.Hour)
	paidDate3 := dueDate3.Add(2 * 24 * time.Hour) // late payment
	expectedPayment3, err := CreatePayment(db, ln.ID, 3, 300.00, 310.00, dueDate3, paidDate3)
	if err != nil {
		t.Fatalf("CreatePayment 3 failed: %v", err)
	}

	expectedPayments := []payment{expectedPayment1, expectedPayment2, expectedPayment3}

	// Act
	actualPayments, err := GetPaymentsByLoanID(db, ln.ID)

	// Assert
	if err != nil {
		t.Fatalf("GetPaymentsByLoanID failed: %v", err)
	}

	require.Equal(t, expectedPayments, actualPayments, "Payments should match and be ordered by payment_number")
}

func TestGetPaymentsByLoanID_NoPayments(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)
	dateTaken := time.Now().UTC().Truncate(24 * time.Hour)

	// Arrange
	usr, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Creating a loan for the test user with no payments
	ln, err := CreateLoan(db, usr.ID, 10000.00, 0.05, 16, 05, "active", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	// Act - no payments created for this loan
	actualPayments, err := GetPaymentsByLoanID(db, ln.ID)

	// Assert
	if err != nil {
		t.Fatalf("GetPaymentsByLoanID failed: %v", err)
	}

	// when comparing actualPayments to an expectedPayments there is an issue since GetPaymentsByLoanID()
	// initializes a slice with nils, which is why it is different
	require.Empty(t, actualPayments, "Should return empty slice for loan with no payments")
}

func TestGetAllPayments(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	dateTaken := time.Now().UTC().Truncate(24 * time.Hour)

	// Arrange, creating multiple test users
	usr1, err := CreateUser(db, "Loan User 1", "loanuser1@test.com", "555-1111")
	if err != nil {
		t.Fatalf("Failed to create test user1: %v", err)
	}

	usr2, err := CreateUser(db, "Loan User 2", "loanuser2@test.com", "555-2222")
	if err != nil {
		t.Fatalf("Failed to create test user2: %v", err)
	}

	// Creating loans for the test users
	ln1, err := CreateLoan(db, usr1.ID, 10000.00, 0.05, 24, 10, "active", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan 1 failed: %v", err)
	}

	ln2, err := CreateLoan(db, usr2.ID, 20000.00, 0.07, 36, 15, "active", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan 2 failed: %v", err)
	}

	// Creating payments for different loans
	dueDate1 := dateTaken.Add(30 * 24 * time.Hour)
	paidDate1 := dueDate1.Add(-2 * 24 * time.Hour)
	expectedPayment1, err := CreatePayment(db, ln1.ID, 1, 500.00, 500.00, dueDate1, paidDate1)
	if err != nil {
		t.Fatalf("CreatePayment 1 failed: %v", err)
	}

	dueDate2 := dateTaken.Add(30 * 24 * time.Hour)
	paidDate2 := dueDate2.Add(-1 * 24 * time.Hour)
	expectedPayment2, err := CreatePayment(db, ln2.ID, 1, 600.00, 600.00, dueDate2, paidDate2)
	if err != nil {
		t.Fatalf("CreatePayment 2 failed: %v", err)
	}

	dueDate3 := dateTaken.Add(60 * 24 * time.Hour)
	paidDate3 := dueDate3.Add(1 * 24 * time.Hour) // late payment
	expectedPayment3, err := CreatePayment(db, ln1.ID, 2, 500.00, 510.00, dueDate3, paidDate3)
	if err != nil {
		t.Fatalf("CreatePayment 3 failed: %v", err)
	}

	var expectedPayments = []payment{expectedPayment1, expectedPayment2, expectedPayment3}

	// Act
	actualPayments, err := GetAllPayments(db)

	// Assert
	if err != nil {
		t.Fatalf("GetAllPayments failed: %v", err)
	}

	require.Equal(t, expectedPayments, actualPayments)
}

func TestDeletePayment(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	dateTaken := time.Now().UTC().Truncate(24 * time.Hour)

	// Arrange, creating a test user
	usr, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Creating a loan for the test user
	ln, err := CreateLoan(db, usr.ID, 10000.00, 0.05, 16, 05, "active", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	// Set up payment dates
	dueDate := dateTaken.Add(30 * 24 * time.Hour) // 30 days after loan was taken
	paidDate := dueDate.Add(-2 * 24 * time.Hour)  // paid 2 days before due date

	// Creating a payment to delete
	pyment, err := CreatePayment(db, ln.ID, 1, 1000.00, 900.00, dueDate, paidDate)
	if err != nil {
		t.Fatalf("CreatePayment failed: %v", err)
	}

	// Act
	err = DeletePayment(db, pyment.ID)
	if err != nil {
		t.Fatalf("DeletePayment failed: %v", err)
	}

	// Assert - verify payment no longer exists
	checkPayments, err := GetPaymentsByLoanID(db, ln.ID)
	if err != nil {
		t.Fatalf("GetPaymentsByLoanID failed: %v", err)
	}

	require.Empty(t, checkPayments)
}
