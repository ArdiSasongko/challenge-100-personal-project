package utils

import (
	"database/sql"
	"log"
)

// Tx mengelola commit atau rollback untuk transaksi yang diberikan.
func Tx(err error, tx *sql.Tx) {
	if err != nil {
		// Jika ada error, lakukan rollback
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Printf("tx rollback error: %v", rbErr)
		} else {
			log.Println("tx rollback successful")
		}
	} else {
		// Jika tidak ada error, lakukan commit
		if cmErr := tx.Commit(); cmErr != nil {
			log.Printf("tx commit error: %v", cmErr)
		} else {
			log.Println("tx commit successful")
		}
	}
}
