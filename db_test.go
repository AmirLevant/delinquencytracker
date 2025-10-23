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
