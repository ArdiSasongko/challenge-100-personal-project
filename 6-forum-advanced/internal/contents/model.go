package contents

import (
	"mime/multipart"
	"time"

	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/internal/comments"
)

type (
	ContentModel struct {
		ID             uint64    `db:"id" json:"id"`
		UserID         int64     `db:"user_id" json:"user_id"`
		ContentTitle   string    `db:"content_title" json:"content_title"`
		ContentBody    string    `db:"content_body" json:"content_body"`
		ContentHastags string    `db:"content_hastags" json:"content_hastags"`
		CreatedAt      time.Time `db:"created_at" json:"created_at"`
		UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
		CreatedBy      string    `db:"created_by" json:"created_by"`
		UpdatedBy      string    `db:"updated_by" json:"updated_by"`
	}

	ImageModel struct {
		ID        uint64    `db:"id" json:"id"`
		ContentID int64     `db:"content_id" json:"content_id"`
		FileName  string    `db:"file_name" json:"file_name"`
		FileUrl   string    `db:"file_url" json:"file_url"`
		CreatedAt time.Time `db:"created_at" json:"created_at"`
		UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
		CreatedBy string    `db:"created_by" json:"created_by"`
		UpdatedBy string    `db:"updated_by" json:"updated_by"`
	}
)

type (
	ContentRequest struct {
		ContentTitle   string                  `json:"content_title" form:"content_title" validate:"required,min=5,max=255"`
		ContentBody    string                  `json:"content_body" form:"content_body" validate:"required"`
		ContentHastags []string                `json:"content_hastags" form:"content_hastags" validate:"omitempty"`
		File           []*multipart.FileHeader `json:"file" form:"file" validate:"omitempty"`
	}

	ContentUpdateRequest struct {
		ContentTitle   string   `json:"content_title" form:"content_title" validate:"omitempty,min=5,max=255"`
		ContentBody    string   `json:"content_body" form:"content_body" validate:"omitempty"`
		ContentHastags []string `json:"content_hastags" form:"content_hastags" validate:"omitempty"`
	}
)

type (
	GetContent struct {
		Content    ContentResponse             `json:"content"`
		LikesCount int                         `json:"likes_count"`
		Comment    []comments.CommentsResponse `json:"comments"`
	}
	ContentResponse struct {
		Data      Data      `json:"data"`
		IsLiked   bool      `json:"is_liked"`
		CreatedAt time.Time `json:"created_at"`
	}
	ContentsResponse struct {
		Data       []Data     `json:"data"`
		Pagination Pagination `json:"pagination"`
	}

	Data struct {
		ID             uint64   `json:"id"`
		ContentTitle   string   `json:"content_title"`
		ContentBody    string   `json:"content_body"`
		FileUrl        string   `json:"file_url"`
		ContentHastags []string `json:"content_hastags"`
		CreatedBy      string   `json:"created_by"`
	}

	Pagination struct {
		Limit  int64 `json:"limit"`
		Offset int64 `json:"offset"`
	}
)
