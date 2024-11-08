package userservice

import (
	"context"
	"errors"
	"time"

	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/model/usermodel"
	"github.com/sirupsen/logrus"
)

func (s *service) UpsertUserActivities(ctx context.Context, userID, content_id int64, username string, req usermodel.UserActivitiesRequest) error {
	now := time.Now()
	model := usermodel.UserActivities{
		UserID:    userID,
		ContentID: content_id,
		IsLiked:   req.IsLiked,
		IsSaved:   req.IsSaved,
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: username,
		UpdatedBy: username,
	}

	userActivity, err := s.repo.GetUserActivity(ctx, model)

	if err != nil {
		logrus.WithField("user service layer", "failed generated jwt token").Error(err.Error())
		return err
	}

	if userActivity == nil {
		if !req.IsLiked {
			req.IsLiked = false
			return errors.New("you've never liked this post")
		}
		err = s.repo.InsertLikeAndSave(ctx, model)
	} else {
		err = s.repo.UpdateUserActivities(ctx, model)
	}

	if err != nil {
		logrus.WithField("user service layer", "failed generated jwt token").Error(err.Error())
		return err
	}

	return nil
}
