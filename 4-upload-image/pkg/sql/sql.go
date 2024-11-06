package sql

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func Connect(dataSource string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		logrus.WithField("error", err.Error()).Fatal("failed connect database")
		return nil, err
	}
	return db, nil
}

func InitStorage(db *sql.DB) {
	if err := db.Ping(); err != nil {
		logrus.WithField("error", err.Error()).Fatal("failed connect database")
	}
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(10 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	logrus.Info("Success Connected Database")
}
