package contentmodel

import "time"

type (
	ContentRequest struct {
		ContentTitle   string   `json:"content_title" validate:"required"`
		Content        string   `json:"content" validate:"required"`
		ContentHastags []string `json:"content_hastags" validate:"omitempty"`
	}
)

type (
	ContentModel struct {
		ID             int64     `db:"id" json:"id"`
		UserID         int64     `db:"user_id" json:"user_id"`
		ContentTitle   string    `db:"content_title" json:"content_title"`
		Content        string    `db:"content" json:"content"`
		ContentHastags string    `db:"content_hastags" json:"content_hastags"`
		CreatedAt      time.Time `db:"created_at" json:"created_at"`
		UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
		CreatedBy      string    `db:"created_by" json:"created_by"`
		UpdatedBy      string    `db:"updated_by" json:"updated_by"`
	}
)
