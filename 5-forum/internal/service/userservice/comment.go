package userservice

import (
	"context"
	"errors"
	"time"

	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/model/usermodel"
	"github.com/sirupsen/logrus"
)

func (s *service) CreateComment(ctx context.Context, req usermodel.CommentRequest, contentID, userID int64, username string) error {
	now := time.Now()

	// checking post
	content, err := s.cr.GetContent(ctx, contentID)
	if err != nil {
		return err
	}

	if content == nil {
		logrus.WithField("user service layer", "content is nil").Error("content didnt exist")
		return errors.New("content didnt exist")
	}

	model := usermodel.CommentModel{
		UserID:    userID,
		ContentID: contentID,
		Comment:   req.Comment,
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: username,
		UpdatedBy: username,
	}

	err = s.repo.CreateComment(ctx, &model)
	if err != nil {
		logrus.WithField("user service layer", err.Error()).Error(err.Error())
		return errors.New("failed create comments")
	}

	return nil
}
