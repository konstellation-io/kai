package objectstorage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/spf13/viper"
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

func (os *S3ObjectStorage) CreateFolder(name string) error {
	os.logger.Info("Creating folder in S3", "bucket", viper.GetString(config.S3BucketKey), "folder", name)

	_, err := os.client.PutObject(&s3.PutObjectInput{
		Key:    aws.String(name + "/"),
		Bucket: aws.String(viper.GetString(config.S3BucketKey)),
	})

	return err
}
