package contentrepository

import (
	"context"
	"database/sql"
	"strings"

	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/model/contentmodel"
	"github.com/sirupsen/logrus"
)

func (r *repository) CreateContent(ctx context.Context, model contentmodel.ContentModel) error {
	query := `INSERT INTO contents(user_id, content_title, content, content_hastags, created_at, updated_at, created_by, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	err := r.db.QueryRowContext(ctx, query, &model.UserID, &model.ContentTitle, &model.Content, &model.ContentHastags, &model.CreatedAt, &model.UpdatedAt, &model.CreatedBy, &model.UpdatedBy).Scan(&model.ID)

	if err != nil {
		logrus.WithField(
			"error", "error failed created content",
		).Error(err.Error())
		return err
	}

	return nil
}

func (r *repository) GetContent(ctx context.Context, contentID, userID int64) (*contentmodel.GetContentResponse, error) {
	query := `
	SELECT c.id, c.content_title, c.content, c.content_hastags, c.created_at, c.created_by, 
	       COALESCE(ua.is_liked, FALSE) AS is_liked
	FROM contents c
	LEFT JOIN users_activities ua 
	       ON c.id = ua.content_id AND ua.user_id = $2
	WHERE c.id = $1
	`

	var (
		model    contentmodel.ContentModel
		response = &contentmodel.GetContentResponse{}
		like     bool
	)

	row := r.db.QueryRowContext(ctx, query, contentID, userID)
	err := row.Scan(&model.ID, &model.ContentTitle, &model.Content, &model.ContentHastags, &model.CreatedAt, &model.CreatedBy, &like)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	response.Creator = contentmodel.CreatorModel{
		CreatedBy: model.CreatedBy,
		CreatedAt: model.CreatedAt,
	}

	response.Content = contentmodel.ContentDetail{
		ContentID:      model.ID,
		ContentTitle:   model.ContentTitle,
		Content:        model.Content,
		ContentHastags: strings.Split(model.ContentHastags, ","),
	}

	response.IsLike = like

	return response, nil
}

func (r *repository) GetContents(ctx context.Context, limit, offset int64) (*contentmodel.GetContents, error) {
	query := `SELECT id, content_title, content, content_hastags, created_at, created_by FROM contents ORDER BY updated_at DESC LIMIT $1 OFFSET $2`

	var (
		content  contentmodel.ContentModel
		response contentmodel.GetContents
		data     = make([]contentmodel.ContentDetail, 0)
	)

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&content.ID, &content.ContentTitle, &content.Content, &content.ContentHastags, &content.CreatedAt, &content.CreatedBy)

		if err != nil {
			return nil, err
		}

		data = append(data, contentmodel.ContentDetail{
			ContentID:      content.ID,
			ContentTitle:   content.ContentTitle,
			Content:        content.Content,
			ContentHastags: strings.Split(content.ContentHastags, ","),
		})
	}

	response.Data = data
	response.Pagination = contentmodel.Pagination{
		Limit:  limit,
		Offset: offset,
	}

	return &response, nil
}
