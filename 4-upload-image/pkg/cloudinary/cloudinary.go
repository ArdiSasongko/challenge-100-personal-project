package cloudinary

import (
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/sirupsen/logrus"
)

type CloudService struct {
	Client *cloudinary.Cloudinary
}

func Init(cloudUrl string) (*CloudService, error) {
	cloud, err := cloudinary.NewFromURL(cloudUrl)
	if err != nil {
		logrus.WithField("error", "failed connect cloudinary").Error(err.Error())
		return nil, err
	}

	return &CloudService{Client: cloud}, nil
}
