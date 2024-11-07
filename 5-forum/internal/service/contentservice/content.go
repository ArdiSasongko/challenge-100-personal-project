package contentservice

import (
	"context"
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
		logrus.WithField("error", "failed create content").Error(err.Error())
		return err
	}

	return nil
}
