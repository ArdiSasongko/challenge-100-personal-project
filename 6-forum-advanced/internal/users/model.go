package users

import "time"

type (
	UserModel struct {
		ID        uint64    `db:"id" json:"id"`
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
		ID        uint64    `db:"id" json:"id"`
		UserID    int64     `db:"user_id" json:"user_id"`
		Token     string    `db:"token" json:"token"`
		ExpiredAt time.Time `db:"expired_at" json:"expired_at"`
		CreatedAt time.Time `db:"created_at" json:"created_at"`
		UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
		CreatedBy string    `db:"created_by" json:"created_by"`
		UpdatedBy string    `db:"updated_by" json:"updated_by"`
	}
)

type (
	UserRequest struct {
		Name     string `json:"name" validate:"required,min=1,max=255"`
		Username string `json:"username" validate:"required,min=6,max=255"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	TokenRequest struct {
		Token string `json:"token" validate:"required"`
	}

	LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}
)

type (
	LoginResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
)
