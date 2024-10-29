package userservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"restapi/internal/model/domain"
	"restapi/internal/model/payload"
	"restapi/internal/model/web"
	userrepository "restapi/internal/repository/user"
	"restapi/internal/utils"

	"github.com/go-playground/validator/v10"
)

var ErrNotFound = errors.New("user not found")

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
		if errValid := utils.ValidError(err); errValid != nil {
			return nil, errValid
		}
		return nil, err
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
		log.Println(err)
		return nil, fmt.Errorf("error :%s", err.Error())
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

func (s *UserService) UpdateUser(ctx context.Context, user payload.UpdateUser, userID int) (*web.ToUser, error) {
	if err := s.validator.Struct(user); err != nil {
		if errValid := utils.ValidError(err); errValid != nil {
			return nil, errValid
		}
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}

	defer utils.Tx(err, tx)

	// check user
	result, err := s.repository.FindByID(ctx, tx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	_, err = s.repository.FindByEmail(ctx, tx, user.Email)

	if err == nil {
		return nil, fmt.Errorf("user with %v already exists", user.Email)
	}

	if err != sql.ErrNoRows {
		log.Println(err)
		return nil, fmt.Errorf("error :%s", err.Error())
	}

	name := utils.DefaultValue[string](result.Name, user.Name)
	age := utils.DefaultValue[int](result.Age, user.Age)
	username := utils.DefaultValue[string](result.Username, user.Username)
	email := utils.DefaultValue[string](result.Email, user.Email)

	if user.Password != "" {
		newPass, err := utils.HashPassword(user.Password)
		if err != nil {
			return nil, err
		}
		user.Password = newPass
	}

	password := utils.DefaultValue[string](result.Password, user.Password)

	// update user
	updateUser := domain.User{
		Name:     name,
		Age:      age,
		Username: username,
		Email:    email,
		Password: password,
	}

	resultUpdate, err := s.repository.Update(ctx, tx, updateUser, userID)
	if err != nil {
		return nil, err
	}

	return &web.ToUser{
		ID:        resultUpdate.ID,
		Name:      resultUpdate.Name,
		Username:  resultUpdate.Username,
		CreatedAt: resultUpdate.CreatedAt,
	}, nil
}

func (s *UserService) FindByID(ctx context.Context, userID int) (*web.ToUserDetail, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}

	defer utils.Tx(err, tx)

	user, err := s.repository.FindByID(ctx, tx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &web.ToUserDetail{
		ID:        user.ID,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, userID int) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	defer utils.Tx(err, tx)

	err = s.repository.Delete(ctx, tx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}

	return nil
}
