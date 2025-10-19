package main

import (
	"fmt"
	"log"
	"time"

	"github.com/amirlevant/delinquencytracker/dbconnection"
)

func main() {
	currentTime := time.Now()
	var printy string = currentTime.Format(time.DateOnly)
	fmt.Println("the date is ", printy)

	config := dbconnection.DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "amir",
		DBName:   "loan_tracker",
	}

	db, err := dbconnection.ConnectDB(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbconnection.CloseDB(db)

	fmt.Println("Success! connected to the database")
	fmt.Println("Database: loan_tracker")
	fmt.Println("Host: localhost:5432")
}
