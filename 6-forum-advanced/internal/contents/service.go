package contents

import (
	"context"
	"database/sql"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"

	cld "github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/pkg/cloudinary"
	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/utils"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/sirupsen/logrus"
)

type service struct {
	repo       Repository
	cloudinary *cld.CloudService
	db         *sql.DB
}

func NewService(repo Repository, cloudinary *cld.CloudService, db *sql.DB) *service {
	return &service{
		repo:       repo,
		cloudinary: cloudinary,
		db:         db,
	}
}

type Service interface {
	CreateContent(ctx context.Context, userID int64, username string, req ContentRequest) error
	GetContents(ctx context.Context, pageSize, pageIndex int64) (*ContentsResponse, error)
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
