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

	// Arrange, creating a test User
	usr, err := CreateUser(db, "Test User", "test@test.com", "555-4444")
	if err != nil {
		t.Fatalf("Failed to create test User: %v", err)
	}

	//Act: Get the User by ID
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

// we provide the GetUserByID an ID of a User that does not exist
// we expect the GetUserByID to fail and return
func TestGetUserByID_UserNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Act, Trying to get a User that does not exist

	usr, err := GetUserByID(db, 99999)

	//Assert, Should return error
	assert.Error(t, err, "Expected error for non-existent User")
	require.Equal(t, User{}, usr, "Expected empty User struct")
}

func TestGetUserByEmail(t *testing.T) {

	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange, creating a test User
	usr, err := CreateUser(db, "Test User", "test@test.com", "555-4444")
	if err != nil {
		t.Fatalf("Failed to create test User: %v", err)
	}

	//Act: Get the User by ID
	usr, err = GetUserByEmail(db, usr.Email)

	//Assert: Check results
	if err != nil {
		t.Errorf("GetUserByEmail failed: %v", err)
	}
	if usr.Name != "Test User" {
		t.Fatalf("Expected name 'Test User', got '%s'", usr.Name)
	}
	if usr.Email != "test@test.com" {
		t.Fatalf("Expected email 'test@test.com', got '%s'", usr.Email)
	}
}

func TestGetUserByPhone(t *testing.T) {

	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange, creating a test User
	usr, err := CreateUser(db, "Test User", "test@test.com", "555-4444")
	if err != nil {
		t.Fatalf("Failed to create test User: %v", err)
	}

	//Act: Get the User by ID
	usr, err = GetUserByPhone(db, usr.Phone)

	//Assert: Check results
	if err != nil {
		t.Errorf("GetUserByEmail failed: %v", err)
	}
	if usr.Name != "Test User" {
		t.Fatalf("Expected name 'Test User', got '%s'", usr.Name)
	}
	if usr.Email != "test@test.com" {
		t.Fatalf("Expected email 'test@test.com', got '%s'", usr.Email)
	}
	if usr.Phone != "555-4444" {
		t.Fatalf("Expected phone '555-4444', got '%s'", usr.Phone)
	}
}

func TestCountUsers(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange, creating a test User
	usr1, err := CreateUser(db, "Test User", "test@test.com", "555-4444")
	if err != nil {
		t.Fatalf("Failed to create test User: %v", err)
	}

	// Arrange, creating a test User
	usr2, err := CreateUser(db, "Test User2", "test2@test.com", "222-4444")
	if err != nil {
		t.Fatalf("Failed to create test User: %v", err)
	}

	_ = usr1 // setting them to empty to remove the err
	_ = usr2 // setting them to empty to remove the err

	expectedCount := int64(2)
	actualCount, err := CountUsers(db)

	if err != nil {
		t.Fatalf("Failed to Count Users: %v", err)
	}

	require.Equal(t, expectedCount, actualCount)

}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange, Creating the first User
	_, err := CreateUser(db, "User One", "duplicate@test.com", "555-0001")
	if err != nil {
		t.Fatalf("Failed to create first User: %v", err)
	}

	//Act, Creating another User with the same email
	_, err = CreateUser(db, "User Two", "duplicate@test.com", "555-0002")

	// Assert, Should return Error
	if err == nil {
		t.Error("Expected error for duplicate email, got nil")
	}
}

// 21/10/25, Creating the Test for updating a User
//
//	The function of updating User still does not exist
//	As expected the test fails

// 21/10/25 I created UpdateUser in logic file
// 21/10/25 test passed as expected
// 22/10/25 updated my CRUD operations to be in db.go
// 23/10/25 need to update the return type of certain functions, tests need to change accordingly
// 23/10/25 updated the test to take the id from the returned User of CreateUser()

func TestUpdateUser(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange
	// creating User
	usr, _ := CreateUser(db, "Old Name", "old@test.com", "555-0000")

	// Act
	// updating User
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

	// esnuring the updated User is as we expected
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
	expectedusers := []User{user1, user2, user3}

	// Assert
	if err != nil {
		t.Errorf("Error is not nil %v", err)
	}
	require.Equal(t, expectedusers, actualusers, "Expected both User slices to be equal")

}

// 23/10/25 create test before Delete User
// 23/10/25 expected that the User gets created then deleted
// 23/10/25 expected that the GetUserByID() will return an empty User
// 23/10/25 expected
func TestDeleteUser(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange, creating the User

	usr, err := CreateUser(db, "Deleted Usr", "deleted@example.com", "555")

	if err != nil {
		t.Fatal("err is not nil in CreateUser %w", err)
	}

	// Act
	// deleting the User
	err = DeleteUser(db, usr.ID)

	if err != nil {
		t.Fatal("err is not nil in DeleteUser but %w", err)
	}

	// verifying that such a User does not exist
	// function should return an empty User and NOT nil
	deletedUsr, err := GetUserByID(db, usr.ID)

	if err == nil {
		t.Fatalf("Error should not be nil, but a message saying User not found")
	}
	require.Equal(t, User{}, deletedUsr, "User has been deleted and not found")

}

// 23/10/25 create test for CreateLoan
// expecting to create a Loan for a User and verify all fields are set correctly
func TestCreateLoan(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange, creating a test User first
	usr, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test User: %v", err)
	}

	// Act, creating a Loan for this User
	dateTaken := time.Now()
	ln, err := CreateLoan(db, usr.ID, 10000.00, 0.05, 36, 15, "active", dateTaken)

	// Assert, Loan creation should succeed
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}
	if ln.ID == 0 {
		t.Error("Expected Loan ID to be set")
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
	// Creating a test User first
	usr, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test User: %v", err)
	}

	// Creating a Loan for this User
	dateTaken := time.Now().UTC().Truncate(24 * time.Hour)
	ln, err := CreateLoan(db, usr.ID, 10000.00, 0.05, 36, 15, "active", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	// Act
	// Updating the Loan with new values
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

	// Should have exactly one Loan
	if len(loans) != 1 {
		t.Fatalf("Expected 1 Loan, got %d", len(loans))
	}

	updatedLoan := loans[0]

	// Ensuring the updated Loan has the new values
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

func TestGetLoanByLoanID(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	dateTaken := time.Now().UTC().Truncate(24 * time.Hour)

	// Arrange, creating a test User first
	usr, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test User: %v", err)
	}

	// Creating a Loan for the test User
	createdLoan, err := CreateLoan(db, usr.ID, 10000.00, 0.05, 36, 15, "active", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	// Act
	retrievedLoan, err := GetLoanByLoanID(db, createdLoan.ID)

	// Assert
	if err != nil {
		t.Fatalf("GetLoanByID failed: %v", err)
	}

	// Verify all fields match
	require.Equal(t, createdLoan, retrievedLoan, "Retrieved Loan should match created Loan")
}

func TestGetLoanByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Act, Trying to get a Loan that does not exist
	ln, err := GetLoanByLoanID(db, 99999)

	// Assert, Should return error
	assert.Error(t, err, "Expected error for non-existent Loan")
	require.Equal(t, Loan{}, ln, "Expected empty Loan struct")
}

func TestGetLoansByUserID_OneLoan(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange, creating a test User first
	usr, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test User: %v", err)
	}

	// Act, creating a Loan for this User
	dateTaken := time.Now()
	expectedln, err := CreateLoan(db, usr.ID, 10000.00, 0.05, 36, 15, "active", dateTaken)

	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	// Act, we query all the loans that belong to the userID
	loans, err := GetLoansByUserID(db, usr.ID)

	actualLn := loans[0]

	// Assert, Loan creation should succeed
	if err != nil {
		t.Fatalf("GetLoansByUserID failed: %v", err)
	}
	if expectedln.ID != actualLn.ID {
		t.Error("Expected Loan ID to match")
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

	// Arrange, creating a test User first
	usr, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test User: %v", err)
	}

	// Act, creating a Loan for this User
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

	var expectedLoans = []Loan{expectedln1, expectedln2, expectedln3}

	require.Equal(t, expectedLoans, actualLoans)

}

func TestGetLoansByUserID_NoLoan(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange, creating a test User first
	usr, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test User: %v", err)
	}

	// Act, no Loan for this User

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

	var expectedLoans = []Loan{expectedln1, expectedln2, expectedln3}

	actualLoans, err := GetAllLoans(db)

	if err != nil {
		t.Fatalf("GetAllLoans failed: %v", err)
	}

	require.Equal(t, expectedLoans, actualLoans)

}

func TestGetLoansByStatus(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Arrange, creating a multiple test users
	dateTaken := time.Now().UTC().Truncate(24 * time.Hour)

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

	expectedln2, err := CreateLoan(db, usr2.ID, 20000.00, 0.25, 26, 15, "active", dateTaken)

	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	expectedln3, err := CreateLoan(db, usr3.ID, 30000.00, 0.35, 36, 25, "defaulted", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	var expectedActiveLoans = []Loan{expectedln1, expectedln2}
	var expectedDefaultedLoans = []Loan{expectedln3}
	var expectedPaidOffLoans = []Loan{}

	// Act

	actualActiveLoans, err := GetLoansByStatus(db, "active")
	if err != nil {
		t.Fatalf("Failed to get Loans by Active Status: %v", err)
	}

	require.Equal(t, expectedActiveLoans, actualActiveLoans)

	actualDefaultedLoans, err := GetLoansByStatus(db, "defaulted")
	if err != nil {
		t.Fatalf("Failed to get Loans by Defaulted Status: %v", err)
	}

	require.Equal(t, expectedDefaultedLoans, actualDefaultedLoans)

	actualPaidOffLoans, err := GetLoansByStatus(db, "paid-off")
	if err != nil {
		t.Fatalf("Failed to get Loans by paid-off Status: %v", err)
	}
	require.Equal(t, expectedPaidOffLoans, actualPaidOffLoans)

}

func TestCountLoansByStatus(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	dateTaken := time.Now().UTC().Truncate(24 * time.Hour)

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

	_, err = CreateLoan(db, usr1.ID, 10000.00, 0.05, 16, 05, "active", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	_, err = CreateLoan(db, usr2.ID, 20000.00, 0.25, 26, 15, "active", dateTaken)

	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	_, err = CreateLoan(db, usr3.ID, 30000.00, 0.35, 36, 25, "defaulted", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	var expectedCountActiveLoans = int64(2)
	var expectedCountDefaultedLoans = int64(1)
	var expectedCountPaidOffLoans = int64(0)

	// Act

	actualCountActiveLoans, err := CountLoansByStatus(db, "active")
	if err != nil {
		t.Fatalf("Failed to get Loans by Active Status: %v", err)
	}

	require.Equal(t, expectedCountActiveLoans, actualCountActiveLoans)

	actualDefaultedLoans, err := CountLoansByStatus(db, "defaulted")
	if err != nil {
		t.Fatalf("Failed to get Loans by Defaulted Status: %v", err)
	}

	require.Equal(t, expectedCountDefaultedLoans, actualDefaultedLoans)

	actualPaidOffLoans, err := CountLoansByStatus(db, "paid-off")
	if err != nil {
		t.Fatalf("Failed to get Loans by paid-off Status: %v", err)
	}
	require.Equal(t, expectedCountPaidOffLoans, actualPaidOffLoans)

}

func TestDeleteLoan(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	dateTaken := time.Now().UTC().Truncate(24 * time.Hour)

	// Arrange, creating a test User
	usr1, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test user1: %v", err)
	}

	// Creating a Loan for the test User
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
		t.Fatalf("Failed to create test User: %v", err)
	}

	// Creating a Loan for the test User
	ln, err := CreateLoan(db, usr.ID, 10000.00, 0.05, 16, 05, "active", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	// Set up Payment dates
	dueDate := dateTaken.Add(30 * 24 * time.Hour) // 30 days after Loan was taken
	paidDate := dueDate.Add(-2 * 24 * time.Hour)  // paid 2 days before due date

	pyment, err := CreatePayment(db, ln.ID, 1, 1000, 900, dueDate, paidDate)
	if err != nil {
		t.Fatalf("Create Payment failed %v:", err)
	}

	var expectedPyment = Payment{pyment.ID, ln.ID, 1, 1000, 900, dueDate, paidDate, pyment.CreatedAt}

	require.Equal(t, expectedPyment, pyment)

}

func TestUpdatePayment(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)
	dateTaken := time.Now().UTC().Truncate(24 * time.Hour)

	// Arrange

	usr, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test User: %v", err)
	}

	// Creating a Loan for the test User
	ln, err := CreateLoan(db, usr.ID, 10000.00, 0.05, 16, 05, "active", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	// Set up Payment dates
	dueDate := dateTaken.Add(30 * 24 * time.Hour) // 30 days after Loan was taken
	paidDate := dueDate.Add(-2 * 24 * time.Hour)  // paid 2 days before due date

	pyment, err := CreatePayment(db, ln.ID, 1, 1000, 900, dueDate, paidDate)
	if err != nil {
		t.Fatalf("Create Payment failed %v:", err)
	}

	// Act
	// Updating the Payment with new values
	newDueDate := dateTaken.Add(45 * 24 * time.Hour)  // 45 days after Loan was taken
	newPaidDate := newDueDate.Add(3 * 24 * time.Hour) // paid 3 days late

	err = UpdatePayment(db, pyment.ID, ln.ID, 2, 1200.00, 1200.00, newDueDate, newPaidDate)

	// Assert
	// Update should succeed
	if err != nil {
		t.Fatalf("UpdatePayment failed: %v", err)
	}

	// Ensuring the update worked by querying the Payment
	updatedPayment, err := GetPaymentByID(db, pyment.ID)
	if err != nil {
		t.Fatalf("GetPaymentByID failed: %v", err)
	}

	// Ensuring the updated Payment has the new values
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
		t.Fatalf("Failed to create test User: %v", err)
	}

	// Creating a Loan for the test User
	ln, err := CreateLoan(db, usr.ID, 10000.00, 0.05, 16, 05, "active", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	// Set up Payment dates
	dueDate := dateTaken.Add(30 * 24 * time.Hour) // 30 days after Loan was taken
	paidDate := dueDate.Add(-2 * 24 * time.Hour)  // paid 2 days before due date

	// Create a Payment to retrieve
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
	require.Equal(t, createdPayment, retrievedPayment, "Retrieved Payment should match created Payment")
}

func TestGetPaymentByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Act
	_, err := GetPaymentByID(db, 99999) // Non-existent ID

	// Assert
	if err == nil {
		t.Fatal("Expected error for non-existent Payment, got nil")
	}
}

func TestGetPaymentsByLoanID_SinglePayment(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)
	dateTaken := time.Now().UTC().Truncate(24 * time.Hour)

	// Arrange
	usr, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test User: %v", err)
	}

	// Creating a Loan for the test User
	ln, err := CreateLoan(db, usr.ID, 10000.00, 0.05, 16, 05, "active", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	// Set up Payment dates
	dueDate := dateTaken.Add(30 * 24 * time.Hour) // 30 days after Loan was taken
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

	require.Len(t, payments, 1, "Should have exactly one Payment")
	require.Equal(t, expectedPayment, payments[0], "Payment should match created Payment")
}

func TestGetPaymentsByLoanID_MultiplePayments(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)
	dateTaken := time.Now().UTC().Truncate(24 * time.Hour)

	// Arrange
	usr, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test User: %v", err)
	}

	// Creating a Loan for the test User
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
	paidDate3 := dueDate3.Add(2 * 24 * time.Hour) // late Payment
	expectedPayment3, err := CreatePayment(db, ln.ID, 3, 300.00, 310.00, dueDate3, paidDate3)
	if err != nil {
		t.Fatalf("CreatePayment 3 failed: %v", err)
	}

	expectedPayments := []Payment{expectedPayment1, expectedPayment2, expectedPayment3}

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
		t.Fatalf("Failed to create test User: %v", err)
	}

	// Creating a Loan for the test User with no payments
	ln, err := CreateLoan(db, usr.ID, 10000.00, 0.05, 16, 05, "active", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	// Act - no payments created for this Loan
	actualPayments, err := GetPaymentsByLoanID(db, ln.ID)

	// Assert
	if err != nil {
		t.Fatalf("GetPaymentsByLoanID failed: %v", err)
	}

	// when comparing actualPayments to an expectedPayments there is an issue since GetPaymentsByLoanID()
	// initializes a slice with nils, which is why it is different
	require.Empty(t, actualPayments, "Should return empty slice for Loan with no payments")
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
	paidDate3 := dueDate3.Add(1 * 24 * time.Hour) // late Payment
	expectedPayment3, err := CreatePayment(db, ln1.ID, 2, 500.00, 510.00, dueDate3, paidDate3)
	if err != nil {
		t.Fatalf("CreatePayment 3 failed: %v", err)
	}

	var expectedPayments = []Payment{expectedPayment1, expectedPayment2, expectedPayment3}

	// Act
	actualPayments, err := GetAllPayments(db)

	// Assert
	if err != nil {
		t.Fatalf("GetAllPayments failed: %v", err)
	}

	require.Equal(t, expectedPayments, actualPayments)
}

func TestGetUnpaidPaymentsByLoanID(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)
	dateTaken := time.Now().UTC().Truncate(24 * time.Hour)

	// Arrange
	usr, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test User: %v", err)
	}

	// Creating a Loan for the test User
	ln, err := CreateLoan(db, usr.ID, 10000.00, 0.05, 36, 15, "active", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	// Create multiple payments with different Payment statuses

	// Payment 1: Fully paid on time
	dueDate1 := dateTaken.Add(30 * 24 * time.Hour)
	paidDate1 := dueDate1.Add(-2 * 24 * time.Hour)
	_, err = CreatePayment(db, ln.ID, 1, 300.00, 300.00, dueDate1, paidDate1)
	if err != nil {
		t.Fatalf("CreatePayment 1 failed: %v", err)
	}

	// Payment 2: Partially paid (unpaid)
	dueDate2 := dateTaken.Add(60 * 24 * time.Hour)
	paidDate2 := dueDate2.Add(-1 * 24 * time.Hour)
	expectedPayment2, err := CreatePayment(db, ln.ID, 2, 300.00, 150.00, dueDate2, paidDate2)
	if err != nil {
		t.Fatalf("CreatePayment 2 failed: %v", err)
	}

	// Payment 3: Not paid at all (PaidDate would be zero/null)
	dueDate3 := dateTaken.Add(90 * 24 * time.Hour)
	expectedPayment3, err := CreatePayment(db, ln.ID, 3, 300.00, 0.00, dueDate3, time.Time{})
	if err != nil {
		t.Fatalf("CreatePayment 3 failed: %v", err)
	}

	// Payment 4: Fully paid late (should not be in unpaid list)
	dueDate4 := dateTaken.Add(120 * 24 * time.Hour)
	paidDate4 := dueDate4.Add(5 * 24 * time.Hour) // 5 days late but fully paid
	_, err = CreatePayment(db, ln.ID, 4, 300.00, 300.00, dueDate4, paidDate4)
	if err != nil {
		t.Fatalf("CreatePayment 4 failed: %v", err)
	}

	// Payment 5: Another unpaid Payment
	dueDate5 := dateTaken.Add(150 * 24 * time.Hour)
	expectedPayment5, err := CreatePayment(db, ln.ID, 5, 300.00, 0.00, dueDate5, time.Time{})
	if err != nil {
		t.Fatalf("CreatePayment 5 failed: %v", err)
	}

	expectedUnpaidPayments := []Payment{expectedPayment2, expectedPayment3, expectedPayment5}

	// Act
	actualUnpaidPayments, err := GetUnpaidPaymentsByLoanID(db, ln.ID)

	// Assert
	if err != nil {
		t.Fatalf("GetUnpaidPaymentsByLoanID failed: %v", err)
	}

	require.Equal(t, len(expectedUnpaidPayments), len(actualUnpaidPayments), "Should have exactly 3 unpaid payments")

	require.Equal(t, expectedUnpaidPayments, actualUnpaidPayments, "Unpaid payments should match expected and be ordered by payment_number")
}

func TestGetUnpaidPaymentsByLoanID_NoUnpaidPayments(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)
	dateTaken := time.Now().UTC().Truncate(24 * time.Hour)

	// Arrange
	usr, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test User: %v", err)
	}

	// Creating a Loan for the test User
	ln, err := CreateLoan(db, usr.ID, 5000.00, 0.04, 12, 10, "active", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	// Create only fully paid payments
	dueDate1 := dateTaken.Add(30 * 24 * time.Hour)
	paidDate1 := dueDate1.Add(-5 * 24 * time.Hour)
	_, err = CreatePayment(db, ln.ID, 1, 450.00, 450.00, dueDate1, paidDate1)
	if err != nil {
		t.Fatalf("CreatePayment 1 failed: %v", err)
	}

	dueDate2 := dateTaken.Add(60 * 24 * time.Hour)
	paidDate2 := dueDate2.Add(-3 * 24 * time.Hour)
	_, err = CreatePayment(db, ln.ID, 2, 450.00, 450.00, dueDate2, paidDate2)
	if err != nil {
		t.Fatalf("CreatePayment 2 failed: %v", err)
	}

	// Act
	actualUnpaidPayments, err := GetUnpaidPaymentsByLoanID(db, ln.ID)

	// Assert
	if err != nil {
		t.Fatalf("GetUnpaidPaymentsByLoanID failed: %v", err)
	}

	require.Empty(t, actualUnpaidPayments, "Should return empty slice when all payments are fully paid")
}

func TestGetUnpaidPaymentsByLoanID_NonExistentLoan(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Act - Query for non-existent Loan ID
	actualUnpaidPayments, err := GetUnpaidPaymentsByLoanID(db, 99999)

	// Assert
	if err != nil {
		t.Fatalf("GetUnpaidPaymentsByLoanID failed: %v", err)
	}

	require.Empty(t, actualUnpaidPayments, "Should return empty slice for non-existent Loan")
}

func TestDeletePayment(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	dateTaken := time.Now().UTC().Truncate(24 * time.Hour)

	// Arrange, creating a test User
	usr, err := CreateUser(db, "Loan User", "loanuser@test.com", "555-1234")
	if err != nil {
		t.Fatalf("Failed to create test User: %v", err)
	}

	// Creating a Loan for the test User
	ln, err := CreateLoan(db, usr.ID, 10000.00, 0.05, 16, 05, "active", dateTaken)
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	// Set up Payment dates
	dueDate := dateTaken.Add(30 * 24 * time.Hour) // 30 days after Loan was taken
	paidDate := dueDate.Add(-2 * 24 * time.Hour)  // paid 2 days before due date

	// Creating a Payment to delete
	pyment, err := CreatePayment(db, ln.ID, 1, 1000.00, 900.00, dueDate, paidDate)
	if err != nil {
		t.Fatalf("CreatePayment failed: %v", err)
	}

	// Act
	err = DeletePayment(db, pyment.ID)
	if err != nil {
		t.Fatalf("DeletePayment failed: %v", err)
	}

	// Assert - verify Payment no longer exists
	checkPayments, err := GetPaymentsByLoanID(db, ln.ID)
	if err != nil {
		t.Fatalf("GetPaymentsByLoanID failed: %v", err)
	}

	require.Empty(t, checkPayments)
}
