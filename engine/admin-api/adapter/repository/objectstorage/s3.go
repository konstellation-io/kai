package objectstorage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-logr/logr"
)

type S3ObjectStorage struct {
	logger logr.Logger
	client *s3.S3
}

func NewS3ObjectStorage(logger logr.Logger, client *s3.S3) *S3ObjectStorage {
	return &S3ObjectStorage{
		logger: logger,
		client: client,
	}
}

func (os *S3ObjectStorage) CreateBucket(name string) error {
	os.logger.Info("Creating S3 bucket", "name", name)

	_, err := os.client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(name),
	})

	return err
}
