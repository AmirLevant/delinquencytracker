package dbconnection

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// DBConfig holds the database connections params
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// ConnectDB establishes a connection to the Postgres database
// It retusn a *sql.DB connection pool and any error encountered

func ConnectDB(config DBConfig) (*sql.DB, error) {
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

func CloseDB(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}
