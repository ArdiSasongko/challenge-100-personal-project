package comments

import "time"

type (
	CommentModel struct {
		ID          uint64    `db:"id" json:"id"`
		UserID      int64     `db:"user_id" json:"user_id"`
		ContentID   int64     `db:"content_id" json:"content_id"`
		CommentBody string    `db:"comment_body" json:"comment_body"`
		CreatedAt   time.Time `db:"created_at" json:"created_at"`
		UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
		CreatedBy   string    `db:"created_by" json:"created_by"`
		UpdatedBy   string    `db:"updated_by" json:"updated_by"`
	}
)

type (
	CommentRequest struct {
		CommentBody string `json:"comment_body" form:"comment_body" validate:"required"`
	}
)
type (
	CommentsResponse struct {
		ID          uint64    `json:"id"`
		UserID      int64     `json:"user_id"`
		ContentID   int64     `json:"content_id"`
		CommentBody string    `json:"comment_body"`
		UpdatedAt   time.Time `json:"updated_at"`
		CreatedBy   string    `json:"created_by"`
	}
)
