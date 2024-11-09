package database

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func Database(dataSource string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		logrus.WithField("database connect", err.Error()).Error(err.Error())
		return nil, err
	}

	return db, nil
}

func IntiStorage(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		logrus.WithField("database connect", "failed connected").Error(err.Error())
		return err
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(10 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	logrus.WithField("database connect", "success connected").Info("success connected database")
	return nil
}
