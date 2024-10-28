package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewSQLStorage(config string) (*sql.DB, error) {
	db, err := sql.Open("postgres", config)
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}
