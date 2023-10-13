package objectstorage

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
	"github.com/spf13/viper"
)

type S3ObjectStorage struct {
	logger logr.Logger
	client *minio.Client
}

func NewS3ObjectStorage(logger logr.Logger, client *minio.Client) *S3ObjectStorage {
	return &S3ObjectStorage{
		logger: logger,
		client: client,
	}
}

func (os *S3ObjectStorage) CreateBucket(ctx context.Context, bucket string) error {
	os.logger.Info("Creating bucket in S3", "bucket", bucket)

	err := os.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
	if err != nil {
		return fmt.Errorf("creating bucket: %w", err)
	}

	if viper.GetString(config.S3TierKey) != "" {
		err = os.client.SetBucketLifecycle(ctx, bucket, &lifecycle.Configuration{
			Rules: []lifecycle.Rule{
				{
					ID: fmt.Sprintf("%s-transition-rule", bucket),
					Transition: lifecycle.Transition{
						StorageClass: viper.GetString(config.S3TierKey),
						Days:         0,
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("setting bucket's lifecyle policy: %w", err)
		}
	}

	return nil
}
