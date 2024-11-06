package userservice

import (
	"context"
	"errors"
	"time"

	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/model/usermodel"
	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/pkg/jwt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) CreateUser(ctx context.Context, req usermodel.UserRequest) error {
	now := time.Now()

	exitingUser, _ := s.repo.GetUser(ctx, 0, req.Username, req.Email)

	if exitingUser != nil {
		return errors.New("username or email already used")
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed generated hash password")
	}

	model := usermodel.UserModel{
		Name:      req.Name,
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(pass),
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: req.Username,
		UpdatedBy: req.Username,
	}

	err = s.repo.CreateUser(ctx, &model)
	if err != nil {
		logrus.WithField(
			"error", "failed create user",
		).Error(err.Error())

		return err
	}

	return nil
}

func (s *service) LoginUser(ctx context.Context, req usermodel.LoginRequest) (*usermodel.LoginResponse, error) {
	user, err := s.repo.GetUser(ctx, 0, req.Username, req.Email)
	if err != nil {
		logrus.Error(err.Error())
		return nil, errors.New("username or email are invalid")
	}

	if user == nil {
		return nil, errors.New("users didnt exists")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("password are invalid")
	}

	accessToken, err := jwt.CreateJWT(user.ID, user.Username, user.Email, s.cfg.Service.SecretJWT)
	if err != nil {
		logrus.WithField("error", "failed generated jwt token").Error(err.Error())
		return nil, err
	}

	// check refreshToken di database
	exitingRefreshToken, err := s.repo.GetToken(ctx, user.ID, time.Now())
	if err != nil {
		logrus.WithField("error", "failed get refresh token").Error(err.Error())
		return nil, err
	}

	if exitingRefreshToken != nil {
		return &usermodel.LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: exitingRefreshToken.RefreshToken,
		}, nil
	}

	refreshToken := jwt.GenerateToken()
	if refreshToken == "" {
		return nil, errors.New("failed generated refresh token")
	}

	now := time.Now()
	model := usermodel.RefreshTokenModel{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiredAt:    time.Now().Add(10 * 24 * time.Hour),
		CreatedAt:    now,
		UpdatedAt:    now,
		CreatedBy:    user.Username,
		UpdatedBy:    user.Username,
	}

	err = s.repo.InsertToken(ctx, model)
	if err != nil {
		logrus.WithField("error", "failed insert refresh token").Error(err.Error())
		return nil, err
	}

	return &usermodel.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
