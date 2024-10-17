package app

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func NewDatabase() *sql.DB {
	db, err := sql.Open("postgres", "user=postgres password=mypassword dbname=basic_crud host=localhost port=5432 sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(10 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)
	fmt.Println("Successfully connected to the database")
	return db
}
