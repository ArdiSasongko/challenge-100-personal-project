package internalsql

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func Database(dataSource string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSource)

	if err != nil {
		logrus.WithField("error", "failed open database").
			Error(err.Error())
		return nil, err
	}

	return db, nil
}

func InitStorage(db *sql.DB) {
	if err := db.Ping(); err != nil {
		logrus.WithField("error", "failed open database").
			Error(err.Error())
	}
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(10 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	logrus.Info("DB: Successfully Connected")
}
