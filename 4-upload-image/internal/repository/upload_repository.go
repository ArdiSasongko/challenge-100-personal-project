package repository

import (
	"context"
	"database/sql"
	"upload_image/internal/model"

	"github.com/sirupsen/logrus"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{
		db: db,
	}
}

type Repository interface {
	UploadImage(ctx context.Context, model model.UploadModel) error
}

func (r *repository) UploadImage(ctx context.Context, model model.UploadModel) error {
	query := `INSERT INTO files (url, file_name) VALUES ($1, $2) RETURNING id`

	var id int64
	err := r.db.QueryRowContext(ctx, query, model.URL, model.FileName).Scan(&id)
	if err != nil {
		logrus.WithField("error", "failed upload").Error(err.Error())
		return err
	}

	return nil
}
