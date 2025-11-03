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

func CreateLoan(db *sql.DB, userID int64, totalAmount, interestRate float64, termMonths, dayDue int, status string, dateTaken time.Time) (loan, error) {
	query := `
        INSERT INTO loans (user_id, total_amount, interest_rate, term_months, day_due, status, date_taken)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at
    `
	var loanID int64
	var createdAt time.Time

	err := db.QueryRow(query, userID, totalAmount, interestRate, termMonths, dayDue, status, dateTaken).Scan(&loanID, &createdAt)
	if err != nil {
		return loan{}, fmt.Errorf("failed to create loan: %w", err)
	}

	ln := loan{loanID, userID, totalAmount, interestRate, termMonths, dayDue, status, dateTaken.UTC(), createdAt.UTC(), nil}
	return ln, nil
}

func GetLoansByUserID(db *sql.DB, userID int64) ([]loan, error) {
	query :=
		`
	SELECT id, user_id, total_amount, interest_rate, term_months, day_due, status, date_taken, created_at
	FROM loans 
	WHERE user_id = $1
	ORDER BY id 
	`

	rows, err := db.Query(query, userID)

	if err != nil {
		return []loan{}, fmt.Errorf("failed to query loans for user %d: %w", userID, err)
	}
	defer rows.Close()

	var loans []loan

	for rows.Next() {
		var l loan

		err := rows.Scan(
			&l.ID,
			&l.UserID,
			&l.TotalAmount,
			&l.InterestRate,
			&l.TermMonths,
			&l.DayDue,
			&l.Status,
			&l.DateTaken,
			&l.CreatedAt,
		)

		if err != nil {
			return []loan{}, fmt.Errorf("failed to scan loan row: %w", err)
		}

		l.DateTaken = l.DateTaken.UTC()
		l.CreatedAt = l.CreatedAt.UTC()
		loans = append(loans, l) // we add l to loans
	}

	// we must check if the loop exited normally or fell silently
	if err = rows.Err(); err != nil {
		return []loan{}, fmt.Errorf("error iterating loan rows: %w", err)
	}

	return loans, nil

}
