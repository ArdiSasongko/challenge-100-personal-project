package users

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/config"
	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/pkg/jwt"
	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	db   *sql.DB
	repo Repository
	cfg  *config.Config
}

func NewService(db *sql.DB, repo Repository, cfg *config.Config) *service {
	return &service{
		db:   db,
		repo: repo,
		cfg:  cfg,
	}
}

type Service interface {
	Register(ctx context.Context, req UserRequest) error
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	GetRefreshToken(ctx context.Context, userID int64, req TokenRequest) (string, error)
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

func (s *service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	tx, err := s.db.Begin()
	if err != nil {
		logrus.WithField("tx err", err.Error()).Error(err.Error())
		return nil, err
	}

	defer utils.Tx(tx, err)

	user, err := s.repo.GetUser(ctx, tx, 0, "", req.Email)
	if err != nil {
		logrus.WithField("get user", err.Error()).Error(err.Error())
		return nil, errors.New("invalid credentials")
	}

	if user == nil {
		logrus.WithField("get user", "user didnt exists").Error("user didnt exists")
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		logrus.WithField("compare password", err.Error()).Error(err.Error())
		return nil, errors.New("invalid credentials")
	}

	// generated access token and refresh token
	claims := jwt.ClaimsToken{
		ID:       int64(user.ID),
		Username: user.Username,
		Email:    user.Email,
	}

	token, err := jwt.GeneratedToken(claims, s.cfg.Service.SecretJWT)
	if err != nil {
		logrus.WithField("generate jwt", err.Error()).Error(err.Error())
		return nil, err
	}

	// check refresh token
	exitingRefreshToken, err := s.repo.GetToken(ctx, tx, int64(user.ID), time.Now())
	if err != nil {
		logrus.WithField("get refresh token", err.Error()).Error(err.Error())
		return nil, err
	}

	// generated refreshToken
	refreshToken := jwt.GeneratedRefreshToken()
	if refreshToken == "" {
		logrus.WithField("created refresh token", "failed create refresh token").Error("failed create refresh token")
		return nil, errors.New("failed generate refresh token")
	}

	model := RefreshTokenModel{
		UserID:    int64(user.ID),
		Token:     refreshToken,
		ExpiredAt: time.Now().Add(10 * 24 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: user.Username,
		UpdatedBy: user.Username,
	}

	if exitingRefreshToken != nil {
		if exitingRefreshToken.ExpiredAt.Before(time.Now()) {
			err := s.repo.UpdateToken(ctx, tx, model)
			if err != nil {
				logrus.WithField("update refresh token", err.Error()).Error(err.Error())
				return nil, err
			}
		}

		return &LoginResponse{
			AccessToken:  token,
			RefreshToken: exitingRefreshToken.Token,
		}, nil
	}

	err = s.repo.InsertToken(ctx, tx, model)
	if err != nil {
		logrus.WithField("insert refresh token", err.Error()).Error(err.Error())
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) GetRefreshToken(ctx context.Context, userID int64, req TokenRequest) (string, error) {
	tx, err := s.db.Begin()
	if err != nil {
		logrus.WithField("tx err", err.Error()).Error(err.Error())
		return "", err
	}

	defer utils.Tx(tx, err)
	exitingToken, err := s.repo.GetToken(ctx, tx, userID, time.Now())
	if err != nil {
		logrus.WithField("get refresh token", err.Error()).Error(err.Error())
		return "", err
	}

	if exitingToken == nil {
		logrus.WithField("get refresh token", "refresh token didnt exists").Error("refresh token didnt exists")
		return "", err
	}

	if exitingToken.Token != req.Token {
		logrus.WithField("get refresh token", "refresh token invalid").Error("refresh token invalid")
		return "", err
	}

	user, err := s.repo.GetUser(ctx, tx, userID, "", "")
	if err != nil {
		logrus.WithField("get user", err.Error()).Error(err.Error())
		return "", errors.New("invalid credentials")
	}

	if user == nil {
		logrus.WithField("get user", "user didnt exists").Error("user didnt exists")
		return "", errors.New("invalid credentials")
	}

	claims := jwt.ClaimsToken{
		ID:       int64(user.ID),
		Username: user.Username,
		Email:    user.Email,
	}

	token, err := jwt.GeneratedToken(claims, s.cfg.Service.SecretJWT)
	if err != nil {
		logrus.WithField("generate jwt", err.Error()).Error(err.Error())
		return "", err
	}

	return token, nil
}
