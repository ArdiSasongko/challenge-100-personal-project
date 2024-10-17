package helper

import (
	"database/sql"
	"log"
)

func RollBackCommit(tx *sql.Tx) {
	if p := recover(); p != nil {
		err := tx.Rollback()
		if err != nil {
			log.Printf("Failed to rollback transaction: %v", err)
		}
		panic(p) // re-throw panic after rollback
	} else {
		err := tx.Commit()
		if err != nil {
			log.Printf("Failed to commit transaction: %v", err)
		}
	}
}
