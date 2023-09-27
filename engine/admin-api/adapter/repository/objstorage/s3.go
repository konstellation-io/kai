package objstorage

import (
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/go-logr/logr"
)

type S3ObjectStorage struct {
	logger logr.Logger
	client *client.Client
}

func NewS3ObjectStorage(logger logr.Logger, client *client.Client) *S3ObjectStorage {
	return &S3ObjectStorage{
		logger: logger,
		client: client,
	}
}

func (os *S3ObjectStorage) CreateBucket(name string) error {
	return nil
}
