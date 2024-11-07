package usermodel

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func MustHaveUsernameOrEmail(sl validator.StructLevel) {
	loginRequest := sl.Current().Interface().(LoginRequest)
	if loginRequest.Username == "" && loginRequest.Email == "" {
		sl.ReportError(loginRequest.Username, "Username", "username", "username_or_email", "")
	}
}

type (
	UserRequest struct {
		Name     string `json:"name" validate:"required"`
		Username string `json:"username" validate:"required,alphanum"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,alphanum"`
	}

	LoginRequest struct {
		Username string `json:"username" validate:"omitempty,alphanum"`
		Email    string `json:"email" validate:"omitempty,email"`
		Password string `json:"password" validate:"required,alphanum"`
	}

	CommentRequest struct {
		Comment string `json:"comment" validate:"required"`
	}
)

type (
	UserModel struct {
		ID        int64     `db:"id" json:"id"`
		Name      string    `db:"name" json:"name"`
		Username  string    `db:"username" json:"username"`
		Email     string    `db:"email" json:"email"`
		Password  string    `db:"password" json:"password"`
		CreatedAt time.Time `db:"created_at" json:"created_at"`
		UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
		CreatedBy string    `db:"created_by" json:"created_by"`
		UpdatedBy string    `db:"updated_by" json:"updated_by"`
	}

	RefreshTokenModel struct {
		ID           int64     `db:"id" json:"id"`
		UserID       int64     `db:"user_id" json:"user_id"`
		RefreshToken string    `db:"refresh_token" json:"refresh_token"`
		ExpiredAt    time.Time `db:"expired_at" json:"expired_at"`
		CreatedAt    time.Time `db:"created_at" json:"created_at"`
		UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
		CreatedBy    string    `db:"created_by" json:"created_by"`
		UpdatedBy    string    `db:"updated_by" json:"updated_by"`
	}

	CommentModel struct {
		ID        int64     `db:"id" json:"id"`
		UserID    int64     `db:"user_id" json:"user_id"`
		ContentID int64     `db:"content_id" json:"content_id"`
		Comment   string    `db:"commnet" json:"comment"`
		CreatedAt time.Time `db:"created_at" json:"created_at"`
		UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
		CreatedBy string    `db:"created_by" json:"created_by"`
		UpdatedBy string    `db:"updated_by" json:"updated_by"`
	}
)

type (
	LoginResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
)
