package usersactivities

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/utils"
	"github.com/sirupsen/logrus"
)

type service struct {
	r  Repository
	db *sql.DB
}

func NewService(r Repository, db *sql.DB) *service {
	return &service{
		r:  r,
		db: db,
	}
}

type Service interface {
	UpsertUserActivities(ctx context.Context, userID int64, username string, contentID int64, req UserActivitiesRequest) error
}

func (s *service) UpsertUserActivities(ctx context.Context, userID int64, username string, contentID int64, req UserActivitiesRequest) error {
	tx, err := s.db.Begin()
	if err != nil {
		logrus.WithField("tx err", err.Error()).Error(err.Error())
		return err
	}
	defer utils.Tx(tx, err)

	now := time.Now()
	model := UserActivitiesModel{
		UserID:    userID,
		ContentID: contentID,
		IsLiked:   req.IsLiked,
		IsSaved:   req.IsSaved,
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: username,
		UpdatedBy: username,
	}

	userActivity, err := s.r.GetUserActivites(ctx, tx, model)
	if err != nil {
		logrus.WithField("get user activities", err.Error()).Error(err.Error())
		return err
	}

	if userActivity == nil {
		if !req.IsLiked {
			return errors.New("you've never liked this post")
		}
		err = s.r.InsertUserActivities(ctx, tx, model)
	} else {
		err = s.r.UpdatedUserActivities(ctx, tx, model)
	}

	if err != nil {
		logrus.WithField("upsert user activities", err.Error()).Error(err.Error())
		return err
	}

	return nil
}
