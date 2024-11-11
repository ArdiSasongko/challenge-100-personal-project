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
	GetContentByID(ctx context.Context, tx *sql.Tx, userID, contentID int64) (*ContentResponse, error)
	UpdateContent(ctx context.Context, tx *sql.Tx, userID, contentID int64, model ContentModel) error
	GetImagebyContent(ctx context.Context, tx *sql.Tx, contentID int64) (*[]ImageModel, error)
	DeleteImages(ctx context.Context, tx *sql.Tx, contentID int64) error
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

func (r *repository) GetContentByID(ctx context.Context, tx *sql.Tx, userID, contentID int64) (*ContentResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
	SELECT c.id, c.content_title, c.content_body, i.file_url, c.content_hastags, c.created_by, c.created_at,
	       COALESCE(ua.is_liked, false) AS is_liked
	FROM contents c 
	LEFT JOIN images i ON i.content_id = c.id
	LEFT JOIN users_activities ua ON c.id = ua.content_id AND ua.user_id = $1
	WHERE c.id = $2
	`

	var (
		fileurl  sql.NullString
		isliked  bool
		content  ContentModel
		response = &ContentResponse{}
	)

	row := tx.QueryRowContext(ctx, query, userID, contentID)
	err := row.Scan(&content.ID, &content.ContentTitle, &content.ContentBody, &fileurl, &content.ContentHastags, &content.CreatedBy, &content.CreatedAt, &isliked)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		logrus.WithContext(ctx).WithField("error", err.Error()).Error(err.Error())
		return nil, err
	}

	response.Data = Data{
		ID:             content.ID,
		ContentTitle:   content.ContentTitle,
		ContentBody:    content.ContentBody,
		FileUrl:        fileurl.String,
		ContentHastags: strings.Split(content.ContentHastags, ","),
		CreatedBy:      content.CreatedBy,
	}

	response.IsLiked = isliked
	response.CreatedAt = content.CreatedAt
	return response, nil
}

func (r *repository) UpdateContent(ctx context.Context, tx *sql.Tx, userID, contentID int64, model ContentModel) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `UPDATE contents SET content_title = $1, content_body = $2, content_hastags = $3, updated_at = $4, updated_by = $5 WHERE user_id = $6 AND id = $7`

	_, err := tx.ExecContext(ctx, query, model.ContentTitle, model.ContentBody, model.ContentHastags, model.UpdatedAt, model.UpdatedBy, userID, contentID)

	if err != nil {
		logrus.WithContext(ctx).WithField("error", err.Error()).Error(err.Error())
		return err
	}

	return nil
}

func (r *repository) GetImagebyContent(ctx context.Context, tx *sql.Tx, contentID int64) (*[]ImageModel, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `SELECT id, content_id, file_name, file_url, created_at, updated_at, created_by, updated_by FROM images WHERE content_id = $1`

	var image ImageModel
	var images = make([]ImageModel, 0)

	row := tx.QueryRowContext(ctx, query, contentID)
	err := row.Scan(&image.ID, &image.ContentID, &image.FileName, &image.FileUrl, &image.CreatedAt, &image.UpdatedAt, &image.CreatedBy, &image.UpdatedBy)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		logrus.WithContext(ctx).WithField("error", err.Error()).Error(err.Error())
		return nil, err
	}

	images = append(images, image)

	return &images, nil
}

func (r *repository) DeleteImages(ctx context.Context, tx *sql.Tx, contentID int64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `DELETE FROM images WHERE content_id = $1`

	_, err := tx.ExecContext(ctx, query, contentID)
	if err != nil {
		logrus.WithContext(ctx).WithField("error", err.Error()).Error(err.Error())
		return err
	}

	return nil
}
