package comments

import (
	"context"
	"database/sql"
	"time"

	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/utils"
	"github.com/sirupsen/logrus"
)

type service struct {
	repo Repository
	db   *sql.DB
}

func NewService(repo Repository, db *sql.DB) *service {
	return &service{
		repo: repo,
		db:   db,
	}
}

type Service interface {
	InsertComment(ctx context.Context, userID int64, contentID int64, username string, req CommentRequest) error
}

func (s *service) InsertComment(ctx context.Context, userID int64, contentID int64, username string, req CommentRequest) error {
	tx, err := s.db.Begin()
	if err != nil {
		logrus.WithField("tx err", err.Error()).Error(err.Error())
		return err
	}
	defer utils.Tx(tx, err)

	now := time.Now()
	model := CommentModel{
		UserID:      userID,
		ContentID:   contentID,
		CommentBody: req.CommentBody,
		CreatedAt:   now,
		UpdatedAt:   now,
		CreatedBy:   username,
		UpdatedBy:   username,
	}

	if err := s.repo.InsertComment(ctx, tx, model); err != nil {
		if err == sql.ErrNoRows {
			logrus.WithField("insert comment", err.Error()).Error(err.Error())
			return err
		}
		logrus.WithField("insert comment", err.Error()).Error(err.Error())
		return err
	}

	return nil
}
