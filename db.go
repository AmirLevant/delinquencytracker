package delinquencytracker

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// dbConfig holds the database connections params
type dbConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// ConnectDB establishes a connection to the Postgres database
// It retusn a *sql.DB connection pool and any error encountered

func newDB(config dbConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * 60)

	return db, nil
}

// we pass db connection and the user information
// we return the new user's ID and any error
func CreateUser(db *sql.DB, name, email, phone string) (user, error) {
	query := `
	INSERT INTO users (name, email, phone)
	VALUES ($1, $2, $3)
	RETURNING id, created_at
	`

	var userID int64
	var createdAt time.Time

	err := db.QueryRow(query, name, email, phone).Scan(&userID, &createdAt)
	if err != nil {
		return user{}, fmt.Errorf("failed to create user: %w", err)
	}

	usr := user{userID, name, email, phone, createdAt, nil}

	return usr, nil
}

func UpdateUser(db *sql.DB, userID int64, name, email, phone string) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, phone = $3
		WHERE id = $4
		`

	_, err := db.Exec(query, name, email, phone, userID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func GetUserByID(db *sql.DB, userID int64) (user, error) {
	query := `
	SELECT id, name, email, phone, created_at
	FROM users
	WHERE id = $1
	`

	usr := user{}

	err := db.QueryRow(query, userID).Scan(
		&usr.ID,
		&usr.Name,
		&usr.Email,
		&usr.Phone,
		&usr.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return user{}, fmt.Errorf("user with ID %d not found", userID)
	}
	if err != nil {
		return user{}, fmt.Errorf("failed to get user: %w", err)
	}

	return usr, nil
}

func GetAllUsers(db *sql.DB) ([]user, error) {
	query :=
		`
	SELECT id, name, email, phone, created_at
	FROM users
	ORDER BY name
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user

	for rows.Next() {
		var usr user
		err := rows.Scan(
			&usr.ID,
			&usr.Name,
			&usr.Email,
			&usr.Phone,
			&usr.CreatedAt)

		// if nil then scan was correct
		if err != nil {
			return nil, err
		}

		users = append(users, usr)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func DeleteUser(db *sql.DB, userID int64) error {
	query :=
		`
	DELETE FROM users
	WHERE id = $1
	`
	_, err := db.Exec(query, userID)

	if err != nil {
		return fmt.Errorf("failed to delete user %w", err)
	}

	return nil

}
