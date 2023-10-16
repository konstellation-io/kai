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

type MinioObjectStorage struct {
	logger logr.Logger
	client *minio.Client
}

func NewMinioObjectStorage(logger logr.Logger, client *minio.Client) *MinioObjectStorage {
	return &MinioObjectStorage{
		logger: logger,
		client: client,
	}
}

func (os *MinioObjectStorage) CreateBucket(ctx context.Context, bucket string) error {
	os.logger.Info("Creating bucket in MinIO", "bucket", bucket)

	err := os.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
	if err != nil {
		return fmt.Errorf("creating bucket: %w", err)
	}

	if viper.GetBool(config.MinioTieringEnabledKey) {
		err = os.client.SetBucketLifecycle(ctx, bucket, &lifecycle.Configuration{
			Rules: []lifecycle.Rule{
				{
					ID: fmt.Sprintf("%s-transition-rule", bucket),
					Transition: lifecycle.Transition{
						StorageClass: viper.GetString(config.MinioTierKey),
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
