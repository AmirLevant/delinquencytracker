package logic

import (
	"database/sql"
	"fmt"
)

// we pass db connection and the user information
// we return the new user's ID and any error
func CreateUser(db *sql.DB, name, email, phone string) (int64, error) {
	query := `
	INSERT INTO users (name, email, phone)
	VALUES (1$, 2$, 3$)
	RETURNING id
	`

	var userID int64

	err := db.QueryRow(query, name, email, phone).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("Failed to create user: %w", err)
	}

	return userID, nil
}
