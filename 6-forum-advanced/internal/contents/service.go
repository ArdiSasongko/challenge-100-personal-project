package contents

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/internal/comments"
	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/internal/usersactivities"
	cld "github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/pkg/cloudinary"
	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/utils"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/sirupsen/logrus"
)

type service struct {
	repo       Repository
	uar        usersactivities.Repository
	cr         comments.Repository
	cloudinary *cld.CloudService
	db         *sql.DB
}

func NewService(repo Repository, uar usersactivities.Repository, cr comments.Repository, cloudinary *cld.CloudService, db *sql.DB) *service {
	return &service{
		repo:       repo,
		uar:        uar,
		cr:         cr,
		cloudinary: cloudinary,
		db:         db,
	}
}

type Service interface {
	CreateContent(ctx context.Context, userID int64, username string, req ContentRequest) error
	GetContents(ctx context.Context, pageSize, pageIndex int64) (*ContentsResponse, error)
	GetContent(ctx context.Context, userID, contentID int64) (*GetContent, error)
	UpdateContent(ctx context.Context, userID, contentID int64, username string, req ContentUpdateRequest) error
	DeleteContent(ctx context.Context, userID, contentID int64, username string) error
}

func (s *service) uploadToCloudInary(ctx context.Context, file *multipart.FileHeader) (string, string, error) {
	fileName := file.Filename
	logrus.Info(fileName)

	src, err := file.Open()
	if err != nil {
		logrus.WithField("open file", err.Error()).Error("failed to open file")
		return "", "", err
	}

	defer src.Close()

	fileUpload := "./assets/upload/" + fileName
	dst, err := os.Create(fileUpload)
	if err != nil {
		logrus.WithField("create file", err.Error()).Error("failed create temp file")
		return "", "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		logrus.WithField("copy file", err.Error()).Error("failed copy temp file")
		return "", "", err
	}

	defer func() {
		os.Remove(fileUpload)
	}()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cloud, err := s.cloudinary.Client.Upload.Upload(ctx, fileUpload, uploader.UploadParams{Folder: "go-folder"})
	if err != nil {
		logrus.WithField("upload file", err.Error()).Error(err.Error())
		return "", "", err
	}

	return cloud.SecureURL, cloud.PublicID, nil
}

func (s *service) getPublicID(fileUrl string) (string, error) {
	filePath := strings.Split(fileUrl, "/go-folder/")[1]
	publicID := strings.TrimSuffix(filePath, path.Ext(filePath))
	return publicID, nil
}

func (s *service) CreateContent(ctx context.Context, userID int64, username string, req ContentRequest) error {
	tx, err := s.db.Begin()
	if err != nil {
		logrus.WithField("tx err", err.Error()).Error(err.Error())
		return err
	}
	defer utils.Tx(tx, err)

	now := time.Now()
	contentModel := ContentModel{
		UserID:         userID,
		ContentTitle:   req.ContentTitle,
		ContentBody:    req.ContentBody,
		ContentHastags: strings.Join(req.ContentHastags, ","),
		CreatedAt:      now,
		UpdatedAt:      now,
		CreatedBy:      username,
		UpdatedBy:      username,
	}
	id, err := s.repo.InsertContent(ctx, tx, contentModel)
	if err != nil {
		logrus.WithField("insert content", err.Error()).Error(err.Error())
		return err
	}

	var publicIDs []string
	for _, file := range req.File {
		filename := file.Filename

		url, publicID, err := s.uploadToCloudInary(ctx, file)
		if err != nil {
			logrus.WithField("upload image", err.Error()).Error("failed upload image to cloudinary")
			return err
		}

		publicIDs = append(publicIDs, publicID)

		imageModel := ImageModel{
			ContentID: int64(id),
			FileName:  filename,
			FileUrl:   url,
			CreatedAt: now,
			UpdatedAt: now,
			CreatedBy: username,
			UpdatedBy: username,
		}

		err = s.repo.InsertImage(ctx, tx, imageModel)
		if err != nil {
			logrus.WithField("insert image", err.Error()).Error("failed insert image")

			for _, id := range publicIDs {
				s.cloudinary.Client.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: id})
			}

			if err := s.repo.DeleteContent(ctx, tx, int64(imageModel.ContentID), contentModel.UserID); err != nil {
				logrus.WithField("delete content", err.Error()).Error("failed delete image")
				return err
			}
			return err
		}
	}

	return nil
}

func (s *service) GetContents(ctx context.Context, pageSize, pageIndex int64) (*ContentsResponse, error) {
	limit := pageSize
	offset := pageSize * (pageIndex - 1)

	tx, err := s.db.Begin()
	if err != nil {
		logrus.WithField("tx err", err.Error()).Error(err.Error())
		return nil, err
	}
	defer utils.Tx(tx, err)

	response, err := s.repo.GetContents(ctx, tx, limit, offset)
	if err != nil {
		logrus.WithField("get contents", err.Error()).Error("failed get contents")
		return nil, err
	}

	return response, nil
}

func (s *service) GetContent(ctx context.Context, userID, contentID int64) (*GetContent, error) {
	tx, err := s.db.Begin()
	if err != nil {
		logrus.WithField("tx err", err.Error()).Error(err.Error())
		return nil, err
	}
	defer utils.Tx(tx, err)

	content, err := s.repo.GetContentByID(ctx, tx, userID, contentID)
	if err != nil {
		logrus.WithField("get content", err.Error()).Error("failed get contents")
		return nil, err
	}

	if content == nil {
		logrus.WithField("get contents", "content not available").Error("failed get contents")
		return nil, errors.New("content not available")
	}

	likes, err := s.uar.CountLikes(ctx, tx, contentID)
	if err != nil {
		logrus.WithField("get likes", err.Error()).Error("failed get contents")
		return nil, err
	}

	allComents, err := s.cr.GetCommentsByContent(ctx, tx, contentID)
	if err != nil {
		logrus.WithField("get comments", err.Error()).Error("failed get contents")
		return nil, err
	}

	if allComents == nil {
		allComents = &[]comments.CommentsResponse{}
	}

	return &GetContent{
		Content:    *content,
		LikesCount: likes,
		Comment:    *allComents,
	}, nil
}

func (s *service) UpdateContent(ctx context.Context, userID, contentID int64, username string, req ContentUpdateRequest) error {
	tx, err := s.db.Begin()
	if err != nil {
		logrus.WithField("tx err", err.Error()).Error(err.Error())
		return err
	}
	defer utils.Tx(tx, err)

	content, err := s.repo.GetContentByID(ctx, tx, userID, contentID)
	if err != nil {
		logrus.WithField("get content", err.Error()).Error("failed get contents")
		return err
	}

	if content == nil {
		logrus.WithField("get contents", "content not available").Error("failed get contents")
		return errors.New("content not available")
	}

	if content.Data.CreatedBy != username {
		logrus.WithField("get contents", "user didnt match").Error("user didnt match")
		return errors.New("UNAUTHORIZED for updated this content")
	}

	logrus.Info("req content hastags :", req.ContentHastags)
	logrus.Info("current content hastags :", content.Data.ContentHastags)
	req.ContentTitle = utils.DefaultValue[string](content.Data.ContentTitle, req.ContentTitle)
	req.ContentBody = utils.DefaultValue[string](content.Data.ContentBody, req.ContentBody)
	req.ContentHastags = utils.DefaultValue[[]string](content.Data.ContentHastags, req.ContentHastags)
	logrus.Info("req content hastags :", req.ContentHastags)

	now := time.Now()
	model := ContentModel{
		ContentTitle:   req.ContentTitle,
		ContentBody:    req.ContentBody,
		ContentHastags: strings.Join(req.ContentHastags, ","),
		UpdatedAt:      now,
		UpdatedBy:      username,
	}

	err = s.repo.UpdateContent(ctx, tx, userID, contentID, model)
	if err != nil {
		logrus.WithField("update content", err.Error()).Error("failed update content")
		return err
	}

	return nil
}

func (s *service) DeleteContent(ctx context.Context, userID, contentID int64, username string) error {
	tx, err := s.db.Begin()
	if err != nil {
		logrus.WithField("tx err", err.Error()).Error(err.Error())
		return err
	}
	defer utils.Tx(tx, err)

	content, err := s.repo.GetContentByID(ctx, tx, userID, contentID)
	if err != nil {
		logrus.WithField("get content", err.Error()).Error("failed get contents")
		return err
	}

	if content == nil {
		logrus.WithField("get contents", "content not available").Error("failed get contents")
		return errors.New("content not available")
	}

	if content.Data.CreatedBy != username {
		logrus.WithField("get contents", "user didnt match").Error("user didnt match")
		return errors.New("UNAUTHORIZED for updated this content")
	}

	images, err := s.repo.GetImagebyContent(ctx, tx, contentID)
	if err != nil {
		logrus.WithField("get images", err.Error()).Error("failed get contents")
		return err
	}

	if images == nil {
		logrus.WithField("get images", "images not available").Error("images get contents")
		return errors.New("images not available")
	}

	for _, image := range *images {
		publicId, err := s.getPublicID(image.FileUrl)
		if err != nil {
			logrus.WithField("extract publicID", err.Error()).Error("failed to extract publicID")
			return err
		}

		_, err = s.cloudinary.Client.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicId})
		if err != nil {
			logrus.WithField("delete image", err.Error()).Error("failed to delete image")
			return err
		}

		err = s.repo.DeleteImages(ctx, tx, image.ContentID)
		if err != nil {
			logrus.WithField("delete image", err.Error()).Error("failed to delete image")
			return err
		}
	}

	err = s.repo.DeleteContent(ctx, tx, contentID, userID)
	if err != nil {
		logrus.WithField("delete content", err.Error()).Error("failed to delete content")
		return err
	}

	return nil
}
