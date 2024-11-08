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

type (
	GetContents struct {
		Data       []ContentDetail `json:"data"`
		Pagination Pagination      `json:"pagination"`
	}

	Pagination struct {
		Limit  int64 `json:"limit"`
		Offset int64 `json:"offset"`
	}

	GetResponse struct {
		Data      GetContentResponse `json:"data"`
		LikeCount int                `json:"liked_count"`
		Comment   []Comment          `json:"comment"`
	}

	GetContentResponse struct {
		Content ContentDetail `json:"content"`
		Creator CreatorModel  `json:"creator"`
		IsLike  bool          `json:"is_liked"`
	}

	ContentDetail struct {
		ContentID      int64    `json:"content_id"`
		ContentTitle   string   `json:"content_title"`
		Content        string   `json:"content"`
		ContentHastags []string `json:"content_hastags"`
	}
	CreatorModel struct {
		CreatedBy string    `json:"created_by"`
		CreatedAt time.Time `json:"created_at"`
	}

	Comment struct {
		ID       int64  `json:"id"`
		UserID   int64  `json:"user_id"`
		Username string `json:"username"`
		Comment  string `json:"comment"`
	}
)
