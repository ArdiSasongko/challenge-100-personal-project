package usersactivities

import "time"

type (
	UserActivitiesModel struct {
		ID        uint64    `db:"id" json:"id"`
		UserID    int64     `db:"user_id" json:"user_id"`
		ContentID int64     `db:"content_id" json:"content_id"`
		IsLiked   bool      `db:"is_liked" json:"is_liked"`
		IsSaved   bool      `db:"is_saved" json:"is_saved"`
		CreatedAt time.Time `db:"created_at" json:"created_at"`
		UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
		CreatedBy string    `db:"created_by" json:"created_by"`
		UpdatedBy string    `db:"updated_by" json:"updated_by"`
	}
)

type (
	UserActivitiesRequest struct {
		IsLiked bool `json:"is_liked" validate:"omitempty,boolean"`
		IsSaved bool `json:"is_saved" validate:"omitempty,boolean"`
	}
)
