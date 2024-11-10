package contents

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type repository struct{}

func NewRepository() *repository {
	return &repository{}
}

type Repository interface {
	InsertContent(ctx context.Context, tx *sql.Tx, model ContentModel) (int, error)
	InsertImage(ctx context.Context, tx *sql.Tx, model ImageModel) error
	DeleteContent(ctx context.Context, tx *sql.Tx, contentID int64, userID int64) error
	GetContents(ctx context.Context, tx *sql.Tx, limit, offset int64) (*ContentsResponse, error)
}

func (r *repository) InsertContent(ctx context.Context, tx *sql.Tx, model ContentModel) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `INSERT INTO contents(user_id, content_title, content_body, content_hastags, created_at, updated_at, created_by, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	err := tx.QueryRowContext(ctx, query, &model.UserID, &model.ContentTitle, &model.ContentBody, &model.ContentHastags, &model.CreatedAt, &model.UpdatedAt, &model.CreatedBy, &model.UpdatedBy).Scan(&model.ID)

	if err != nil {
		logrus.WithContext(ctx).WithField("error", err.Error()).Error(err.Error())
		return 0, err
	}

	return int(model.ID), nil
}

func (r *repository) InsertImage(ctx context.Context, tx *sql.Tx, model ImageModel) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `INSERT INTO images(content_id, file_name, file_url, created_at, updated_at, created_by, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err := tx.QueryRowContext(ctx, query, &model.ContentID, &model.FileName, &model.FileUrl, &model.CreatedAt, &model.UpdatedAt, &model.CreatedBy, &model.UpdatedBy).Scan(&model.ID)

	if err != nil {
		logrus.WithContext(ctx).WithField("error", err.Error()).Error(err.Error())
		return err
	}

	return nil
}

func (r *repository) DeleteContent(ctx context.Context, tx *sql.Tx, contentID int64, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `DELETE FROM contents WHERE id = $1 AND user_id = $2`
	_, err := tx.ExecContext(ctx, query, contentID, userID)

	if err != nil {
		logrus.WithContext(ctx).WithField("error", err.Error()).Error(err.Error())
		return err
	}

	return nil
}

func (r *repository) GetContents(ctx context.Context, tx *sql.Tx, limit, offset int64) (*ContentsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `SELECT c.id, c.content_title, c.content_body, i.file_url, c.content_hastags, c.created_by 
	FROM contents c 
	LEFT JOIN images i ON i.content_id = c.id 
	ORDER BY c.created_at 
	DESC LIMIT $1 OFFSET $2`

	var (
		fileurl  sql.NullString
		content  ContentModel
		contents = make([]Data, 0)
		response ContentsResponse
	)

	rows, err := tx.QueryContext(ctx, query, limit, offset)
	if err != nil {
		logrus.WithContext(ctx).WithField("error", "contents didnt available").Error("content didnt available")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&content.ID, &content.ContentTitle, &content.ContentBody, &fileurl, &content.ContentHastags, &content.CreatedBy)
		if err != nil {
			logrus.WithContext(ctx).WithField("error", err.Error()).Error(err.Error())
			return nil, err
		}

		contents = append(contents, Data{
			ID:             content.ID,
			ContentTitle:   content.ContentTitle,
			ContentBody:    content.ContentBody,
			FileUrl:        fileurl.String,
			ContentHastags: strings.Split(content.ContentHastags, ","),
			CreatedBy:      content.CreatedBy,
		})
	}

	if len(contents) == 0 {
		logrus.WithContext(ctx).Info("no contents available")
		return &ContentsResponse{Data: []Data{}, Pagination: Pagination{Limit: limit, Offset: offset}}, nil
	}

	response.Data = contents
	response.Pagination = Pagination{
		Limit:  limit,
		Offset: offset,
	}

	return &response, nil
}
