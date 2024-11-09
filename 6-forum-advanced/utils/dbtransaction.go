package utils

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

func Tx(tx *sql.Tx, err error) {
	if err != nil {
		if errRB := tx.Rollback(); errRB != nil {
			logrus.WithField("tx rollback", errRB.Error()).Error(errRB.Error())
		} else {
			logrus.WithField("tx rollback", "rollback success").Info("rollback success")
		}
	} else {
		if errCM := tx.Commit(); errCM != nil {
			logrus.WithField("tx commit", errCM.Error()).Error(errCM.Error())
		} else {
			logrus.WithField("tx commit", "commit success").Info("commit success")
		}
	}
}
