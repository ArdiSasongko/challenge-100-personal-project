package utils

import (
	"database/sql"
	"log"
)

func Tx(err error, tx *sql.Tx) {
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Printf("tx rollback error: %v", rbErr)
		}
	} else {
		if cmErr := tx.Commit(); cmErr != nil {
			log.Printf("tx commit error: %v", cmErr)
		}
	}
}
