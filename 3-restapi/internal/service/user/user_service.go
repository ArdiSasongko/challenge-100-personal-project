package userservice

import (
	"context"
	"database/sql"
	"fmt"
	"restapi/internal/model/domain"
	"restapi/internal/model/payload"
	"restapi/internal/model/web"
	userrepository "restapi/internal/repository/user"
	"restapi/internal/utils"

	"github.com/go-playground/validator/v10"
)

type UserService struct {
	repository userrepository.UserRepositoryInterface
	validator  *validator.Validate
	DB         *sql.DB
}

func NewUserService(repository userrepository.UserRepositoryInterface, validator *validator.Validate, DB *sql.DB) *UserService {
	return &UserService{
		repository: repository,
		validator:  validator,
		DB:         DB,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user payload.CreateUser) (*web.ToUser, error) {
	// validate user
	if err := s.validator.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		return nil, errors
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}

	defer utils.Tx(err, tx)

	// check user
	_, err = s.repository.FindByEmail(ctx, tx, user.Email)
	if err == nil {
		return nil, fmt.Errorf("user with %v already exists", user.Email)
	}

	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("database error")
	}

	// hash password
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	// save user
	newUser := domain.User{
		Name:     user.Name,
		Age:      user.Age,
		Username: user.Username,
		Email:    user.Email,
		Password: hash,
	}

	result, err := s.repository.Create(ctx, tx, newUser)

	if err != nil {
		return nil, err
	}

	return &web.ToUser{
		ID:        result.ID,
		Name:      result.Name,
		Username:  result.Username,
		CreatedAt: result.CreatedAt,
	}, nil
}
