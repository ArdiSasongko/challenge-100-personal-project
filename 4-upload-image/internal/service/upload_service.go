package service

import (
	"context"
	"upload_image/internal/model"
	"upload_image/internal/repository"
	"upload_image/pkg/cloudinary"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type service struct {
	repository repository.Repository
	cloudinary *cloudinary.CloudService
}

func NewService(repository repository.Repository, cloudinary *cloudinary.CloudService) *service {
	return &service{
		repository: repository,
		cloudinary: cloudinary,
	}
}

type Service interface {
	UploadImage(ctx context.Context, fileName string) (*model.UploadModel, error)
}

func (s *service) UploadImage(ctx context.Context, fileName string) (*model.UploadModel, error) {
	uploadFile, err := s.cloudinary.Client.Upload.Upload(ctx, fileName, uploader.UploadParams{})
	if err != nil {
		return nil, err
	}

	file := model.UploadModel{
		FileName: fileName,
		URL:      uploadFile.SecureURL,
	}

	if err := s.repository.UploadImage(ctx, file); err != nil {
		return nil, err
	}

	return &file, nil
}
