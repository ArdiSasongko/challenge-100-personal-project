package users

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	db   *sql.DB
	repo Repository
}

func NewService(db *sql.DB, repo Repository) *service {
	return &service{
		db:   db,
		repo: repo,
	}
}

type Service interface {
	Register(ctx context.Context, req UserRequest) error
}

func (s *service) Register(ctx context.Context, req UserRequest) error {
	tx, err := s.db.Begin()
	if err != nil {
		logrus.WithField("tx err", err.Error()).Error(err.Error())
		return err
	}

	defer utils.Tx(tx, err)
	// find username and email
	exitingUser, err := s.repo.GetUser(ctx, tx, 0, req.Username, req.Email)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}

	if exitingUser != nil {
		logrus.WithField("get user", "user exiting").Warn("user exiting")
		return errors.New("username or email already exist please try another one")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.WithField("hash password", err.Error()).Error(err.Error())
		return err
	}

	now := time.Now()
	model := UserModel{
		Name:      req.Name,
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hash),
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: req.Username,
		UpdatedBy: req.Username,
	}

	err = s.repo.InsertUser(ctx, tx, model)
	if err != nil {
		logrus.WithField("insert user", err.Error()).Error(err.Error())
		return err
	}

	return nil
}
