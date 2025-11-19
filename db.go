package delinquencytracker

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// we pass db connection and the User information
// we return the new User's ID and any error
func CreateUser(db *sql.DB, name, email, phone string) (User, error) {
	query := `
	INSERT INTO users (name, email, phone)
	VALUES ($1, $2, $3)
	RETURNING id, created_at
	`

	var userID int64
	var createdAt time.Time

	err := db.QueryRow(query, name, email, phone).Scan(&userID, &createdAt)
	if err != nil {
		return User{}, fmt.Errorf("failed to create User: %w", err)
	}

	usr := User{userID, name, email, phone, createdAt, nil}

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
		return fmt.Errorf("failed to update User: %w", err)
	}

	return nil
}

func GetUserByID(db *sql.DB, userID int64) (User, error) {
	query := `
	SELECT id, name, email, phone, created_at
	FROM users
	WHERE id = $1
	`

	usr := User{}

	err := db.QueryRow(query, userID).Scan(
		&usr.ID,
		&usr.Name,
		&usr.Email,
		&usr.Phone,
		&usr.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return User{}, fmt.Errorf("User with ID %d not found", userID)
	}
	if err != nil {
		return User{}, fmt.Errorf("failed to get User: %w", err)
	}

	return usr, nil
}

func GetUserByEmail(db *sql.DB, email string) (User, error) {
	query := `
	SELECT id, name, email, phone, created_at
	FROM users
	WHERE email = $1
	`

	usr := User{}

	err := db.QueryRow(query, email).Scan(
		&usr.ID,
		&usr.Name,
		&usr.Email,
		&usr.Phone,
		&usr.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return User{}, fmt.Errorf("User with Email %s not found", email)
	}
	if err != nil {
		return User{}, fmt.Errorf("failed to get User: %w", err)
	}

	return usr, nil
}

func GetUserByPhone(db *sql.DB, phone string) (User, error) {
	query := `
	SELECT id, name, email, phone, created_at
	FROM users
	WHERE phone = $1
	`

	usr := User{}

	err := db.QueryRow(query, phone).Scan(
		&usr.ID,
		&usr.Name,
		&usr.Email,
		&usr.Phone,
		&usr.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return User{}, fmt.Errorf("User with phone %s not found", phone)
	}
	if err != nil {
		return User{}, fmt.Errorf("failed to get User: %w", err)
	}

	return usr, nil
}

func GetAllUsers(db *sql.DB) ([]User, error) {
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

	var users []User

	for rows.Next() {
		var usr User
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

func CountUsers(db *sql.DB) (int64, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int64

	err := db.QueryRow(query).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

func DeleteUser(db *sql.DB, userID int64) error {
	query :=
		`
	DELETE FROM users
	WHERE id = $1
	`
	_, err := db.Exec(query, userID)

	if err != nil {
		return fmt.Errorf("failed to delete User %w", err)
	}

	return nil

}

func CreateLoan(db *sql.DB, userID int64, totalAmount, interestRate float64, termMonths, dayDue int, status string, dateTaken time.Time) (Loan, error) {
	query := `
        INSERT INTO loans (user_id, total_amount, interest_rate, term_months, day_due, status, date_taken)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at
    `
	var loanID int64
	var createdAt time.Time

	err := db.QueryRow(query, userID, totalAmount, interestRate, termMonths, dayDue, status, dateTaken).Scan(&loanID, &createdAt)
	if err != nil {
		return Loan{}, fmt.Errorf("failed to create Loan: %w", err)
	}

	ln := Loan{loanID, userID, totalAmount, interestRate, termMonths, dayDue, status, dateTaken.UTC(), createdAt.UTC(), nil}
	return ln, nil
}

func UpdateLoan(db *sql.DB, loanID int64, totalAmount, interestRate float64, termMonths, dayDue int, status string, dateTaken time.Time) error {
	query := `
		UPDATE loans
		SET total_amount = $1, interest_rate = $2, term_months = $3, day_due = $4, status = $5, date_taken = $6
		WHERE id = $7
	`

	result, err := db.Exec(query, totalAmount, interestRate, termMonths, dayDue, status, dateTaken, loanID)
	if err != nil {
		return fmt.Errorf("failed to update Loan: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Loan with ID %d not found", loanID)
	}

	return nil
}

// Get a singular Loan based on it's ID
func GetLoanByLoanID(db *sql.DB, loanID int64) (Loan, error) {
	query := `
	SELECT id, user_id, total_amount, interest_rate, term_months, day_due, status, date_taken, created_at
	FROM loans
	WHERE id = $1
	`

	var l Loan

	err := db.QueryRow(query, loanID).Scan(
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

	if err == sql.ErrNoRows {
		return Loan{}, fmt.Errorf("Loan with ID %d not found", loanID)
	}
	if err != nil {
		return Loan{}, fmt.Errorf("failed to get Loan: %w", err)
	}

	l.DateTaken = l.DateTaken.UTC()
	l.CreatedAt = l.CreatedAt.UTC()

	return l, nil
}

// Get all loans associated to a User
func GetLoansByUserID(db *sql.DB, userID int64) ([]Loan, error) {
	query :=
		`
	SELECT id, user_id, total_amount, interest_rate, term_months, day_due, status, date_taken, created_at
	FROM loans 
	WHERE user_id = $1
	ORDER BY id 
	`

	rows, err := db.Query(query, userID)

	if err != nil {
		return []Loan{}, fmt.Errorf("failed to query loans for User %d: %w", userID, err)
	}
	defer rows.Close()

	var loans []Loan

	for rows.Next() {
		var l Loan

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
			return []Loan{}, fmt.Errorf("failed to scan Loan row: %w", err)
		}

		l.DateTaken = l.DateTaken.UTC()
		l.CreatedAt = l.CreatedAt.UTC()
		loans = append(loans, l) // we add l to loans
	}

	// we must check if the loop exited normally or fell silently
	if err = rows.Err(); err != nil {
		return []Loan{}, fmt.Errorf("error iterating Loan rows: %w", err)
	}

	return loans, nil

}

// Gets all the loans in the database
func GetAllLoans(db *sql.DB) ([]Loan, error) {
	query :=
		`
	SELECT id, user_id, total_amount, interest_rate, term_months, day_due, status, date_taken, created_at
	FROM loans 
	ORDER BY id 
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	loans := []Loan{}

	for rows.Next() {
		var ln Loan

		err := rows.Scan(
			&ln.ID,
			&ln.UserID,
			&ln.TotalAmount,
			&ln.InterestRate,
			&ln.TermMonths,
			&ln.DayDue,
			&ln.Status,
			&ln.DateTaken,
			&ln.CreatedAt,
		)

		if err != nil {
			return nil, err
		}
		ln.DateTaken = ln.DateTaken.UTC()
		ln.CreatedAt = ln.CreatedAt.UTC()

		loans = append(loans, ln)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return loans, nil
}

// GetLoansByStatus retrieves all loans with a specific status
func GetLoansByStatus(db *sql.DB, status string) ([]Loan, error) {
	query := `
	SELECT id, user_id, total_amount, interest_rate, term_months, day_due, status, date_taken, created_at 
	FROM loans
	where status = $1
	ORDER BY id
	`
	rows, err := db.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	loans := []Loan{}

	for rows.Next() {
		var ln Loan

		err := rows.Scan(
			&ln.ID,
			&ln.UserID,
			&ln.TotalAmount,
			&ln.InterestRate,
			&ln.TermMonths,
			&ln.DayDue,
			&ln.Status,
			&ln.DateTaken,
			&ln.CreatedAt,
		)

		if err != nil {
			return nil, err
		}
		ln.DateTaken = ln.DateTaken.UTC()
		ln.CreatedAt = ln.CreatedAt.UTC()

		loans = append(loans, ln)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return loans, nil
}

// CountLoansByStatus returns the count of loans with a specific status
func CountLoansByStatus(db *sql.DB, status string) (int64, error) {
	query := `
	SELECT COUNT(*) 
	FROM loans 
	where status = $1`

	var count int64

	err := db.QueryRow(query, status).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

func DeleteLoan(db *sql.DB, LoanID int64) error {
	query :=
		`
	DELETE FROM loans 
	where id = $1
	`

	_, err := db.Exec(query, LoanID)

	if err != nil {
		return fmt.Errorf("failed to delete Loan %w", err)
	}

	return nil
}

func CreatePayment(db *sql.DB, LoanID, payment_number int64, AmountDue, AmountPaid float64, DueDate, PaidDate time.Time) (Payment, error) {
	query :=
		`
	INSERT INTO payments (loan_id, payment_number, amount_due, amount_paid, due_date, paid_date)
	VALUES ($1, $2, $3, $4, $5, $6)
	returning id, created_at
	`

	var paymentID int64
	var createdAt time.Time

	err := db.QueryRow(query, LoanID, payment_number, AmountDue, AmountPaid, DueDate, PaidDate).Scan(&paymentID, &createdAt)
	if err != nil {
		return Payment{}, fmt.Errorf("failed to create Payment: %w", err)
	}

	pyment := Payment{paymentID, LoanID, payment_number, AmountDue, AmountPaid, DueDate.UTC(), PaidDate.UTC(), createdAt.UTC()}
	return pyment, nil
}

func UpdatePayment(db *sql.DB, UserID, LoanID, payment_number int64, AmountDue, AmountPaid float64, DueDate, PaidDate time.Time) error {
	query :=
		`
	UPDATE payments
	SET loan_id = $1, payment_number = $2, amount_due = $3, amount_paid = $4, due_date = $5, paid_date = $6
	WHERE id = $7
	`

	result, err := db.Exec(query, LoanID, payment_number, AmountDue, AmountPaid, DueDate, PaidDate, UserID)
	if err != nil {
		return fmt.Errorf("failed to update Payment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Payment with ID %d not found", UserID)
	}

	return nil

}

func GetPaymentByID(db *sql.DB, paymentID int64) (Payment, error) {
	query := `
        SELECT id, loan_id, payment_number, amount_due, amount_paid, due_date, paid_date, created_at
        FROM payments
        WHERE id = $1
    `

	var p Payment
	err := db.QueryRow(query, paymentID).Scan(
		&p.ID,
		&p.LoanID,
		&p.PaymentNumber,
		&p.AmountDue,
		&p.AmountPaid,
		&p.DueDate,
		&p.PaidDate,
		&p.CreatedAt,
	)
	p.DueDate = p.DueDate.UTC()
	p.PaidDate = p.PaidDate.UTC()
	p.CreatedAt = p.CreatedAt.UTC()

	if err != nil {
		return Payment{}, fmt.Errorf("failed to get Payment: %w", err)
	}

	return p, nil
}

// Gets all the payments associated with a singular Loan
func GetPaymentsByLoanID(db *sql.DB, loanID int64) ([]Payment, error) {
	query := `
	SELECT id, loan_id, payment_number, amount_due, amount_paid, due_date, paid_date, created_at
	FROM payments
	WHERE loan_id = $1
	ORDER BY payment_number
	`

	rows, err := db.Query(query, loanID)
	if err != nil {
		return []Payment{}, fmt.Errorf("failed to query payments for Loan %d: %w", loanID, err)
	}
	defer rows.Close()

	var payments []Payment

	for rows.Next() {
		var p Payment

		err := rows.Scan(
			&p.ID,
			&p.LoanID,
			&p.PaymentNumber,
			&p.AmountDue,
			&p.AmountPaid,
			&p.DueDate,
			&p.PaidDate,
			&p.CreatedAt,
		)

		if err != nil {
			return []Payment{}, fmt.Errorf("failed to scan Payment row: %w", err)
		}

		p.DueDate = p.DueDate.UTC()
		p.PaidDate = p.PaidDate.UTC()
		p.CreatedAt = p.CreatedAt.UTC()

		payments = append(payments, p)
	}

	// we must check if the loop exited normally or fell silently
	if err = rows.Err(); err != nil {
		return []Payment{}, fmt.Errorf("error iterating Payment rows: %w", err)
	}

	return payments, nil
}

// Gets all the payments in the database, regardless of Loan
func GetAllPayments(db *sql.DB) ([]Payment, error) {
	query :=
		`
	SELECT id, loan_id, payment_number, amount_due, amount_paid, due_date, paid_date, created_at
	FROM payments
	ORDER BY id
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []Payment

	for rows.Next() {
		var p Payment

		err := rows.Scan(
			&p.ID,
			&p.LoanID,
			&p.PaymentNumber,
			&p.AmountDue,
			&p.AmountPaid,
			&p.DueDate,
			&p.PaidDate,
			&p.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		p.DueDate = p.DueDate.UTC()
		p.PaidDate = p.PaidDate.UTC()
		p.CreatedAt = p.CreatedAt.UTC()

		payments = append(payments, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

// GetUnpaidPaymentsByLoanID retrieves all unpaid payments for a Loan
func GetUnpaidPaymentsByLoanID(db *sql.DB, loanID int64) ([]Payment, error) {
	query := `
	SELECT id, loan_id, payment_number, amount_due, amount_paid, due_date, paid_date, created_at
	FROM payments
	WHERE loan_id = $1 
	AND (paid_date IS NULL OR amount_paid < amount_due)
	ORDER BY payment_number
	`

	rows, err := db.Query(query, loanID)
	if err != nil {
		return []Payment{}, fmt.Errorf("failed to query unpaid payments for Loan %d: %w", loanID, err)
	}
	defer rows.Close()

	var payments []Payment

	for rows.Next() {
		var p Payment

		err := rows.Scan(
			&p.ID,
			&p.LoanID,
			&p.PaymentNumber,
			&p.AmountDue,
			&p.AmountPaid,
			&p.DueDate,
			&p.PaidDate,
			&p.CreatedAt,
		)

		if err != nil {
			return []Payment{}, fmt.Errorf("failed to scan Payment row: %w", err)
		}

		p.DueDate = p.DueDate.UTC()
		p.PaidDate = p.PaidDate.UTC()
		p.CreatedAt = p.CreatedAt.UTC()

		payments = append(payments, p)
	}

	// Check if the loop exited normally or fell silently
	if err = rows.Err(); err != nil {
		return []Payment{}, fmt.Errorf("error iterating Payment rows: %w", err)
	}

	return payments, nil
}

// Deletes a singular Payment based on a given ID
func DeletePayment(db *sql.DB, paymentID int64) error {
	query :=
		`
	DELETE FROM payments
	WHERE id = $1
	`
	_, err := db.Exec(query, paymentID)

	if err != nil {
		return fmt.Errorf("failed to delete Payment %w", err)
	}

	return nil
}
