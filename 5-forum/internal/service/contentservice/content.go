package contentservice

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/model/contentmodel"
	"github.com/sirupsen/logrus"
)

func (s *service) CreateContent(ctx context.Context, req contentmodel.ContentRequest, userID int64, username string) error {
	now := time.Now()
	contentHastags := strings.Join(req.ContentHastags, ",")

	model := contentmodel.ContentModel{
		UserID:         userID,
		ContentTitle:   req.ContentTitle,
		Content:        req.Content,
		ContentHastags: contentHastags,
		CreatedAt:      now,
		UpdatedAt:      now,
		CreatedBy:      username,
		UpdatedBy:      username,
	}

	err := s.repo.CreateContent(ctx, model)

	if err != nil {
		logrus.WithField("content service layer", "failed create content").Error(err.Error())
		return err
	}

	return nil
}

func (s *service) GetContent(ctx context.Context, contentId, userID int64) (*contentmodel.GetResponse, error) {
	contentData, err := s.repo.GetContent(ctx, contentId, userID)
	if err != nil {
		logrus.WithField("content service layer", err.Error()).Error(err.Error())
		return nil, err
	}

	liked, err := s.ur.CountLikes(ctx, contentId)
	if err != nil {
		logrus.WithField("content service layer", err.Error()).Error(err.Error())
		return nil, err
	}

	comments, err := s.ur.GetCommentByContent(ctx, contentId)
	if err != nil {
		logrus.WithField("content service layer", err.Error()).Error(err.Error())
		return nil, err
	}

	return &contentmodel.GetResponse{
		Data:      *contentData,
		LikeCount: liked,
		Comment:   *comments,
	}, nil
}

func (s *service) GetContents(ctx context.Context, pagiSize, pageIndex int64) (*contentmodel.GetContents, error) {
	limit := pagiSize
	offset := pagiSize * (pageIndex - 1)
	response, err := s.repo.GetContents(ctx, limit, offset)
	if err != nil {
		logrus.WithField("content service layer", err.Error()).Error(err.Error())
		return nil, errors.New("failed get all contents")
	}

	return response, nil
}
