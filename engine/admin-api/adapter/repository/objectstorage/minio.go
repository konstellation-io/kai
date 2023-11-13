package objectstorage

import (
	"bytes"
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/minio/madmin-go/v3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
	"github.com/spf13/viper"
)

type MinioObjectStorage struct {
	logger      logr.Logger
	client      *minio.Client
	adminClient *madmin.AdminClient
}

func NewMinioObjectStorage(logger logr.Logger, client *minio.Client, adminClient *madmin.AdminClient) *MinioObjectStorage {
	return &MinioObjectStorage{
		logger:      logger,
		client:      client,
		adminClient: adminClient,
	}
}

func (os *MinioObjectStorage) CreateBucket(ctx context.Context, bucket string) error {
	os.logger.Info("Creating bucket in MinIO", "bucket", bucket)

	err := os.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
	if err != nil {
		return fmt.Errorf("creating bucket: %w", err)
	}

	err = os.client.EnableVersioning(ctx, bucket)
	if err != nil {
		return fmt.Errorf("enabling versioning in bucket: %w", err)
	}

	if viper.GetBool(config.MinioTierEnabledKey) {
		err = os.client.SetBucketLifecycle(ctx, bucket, &lifecycle.Configuration{
			Rules: []lifecycle.Rule{
				{
					ID:     fmt.Sprintf("%s-transition-rule", bucket),
					Status: minio.Enabled,
					Transition: lifecycle.Transition{
						StorageClass: viper.GetString(config.MinioTierNameKey),
						Days:         lifecycle.ExpirationDays(viper.GetInt(config.MinioTierTransitionDaysKey)),
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("setting bucket's lifecyle: %w", err)
		}
	}

	return nil
}

func (os *MinioObjectStorage) CreateBucketPolicy(ctx context.Context, bucket string) (string, error) {
	policyName := bucket

	err := os.adminClient.AddCannedPolicy(
		ctx,
		policyName,
		os.getPolicy(bucket),
	)
	if err != nil {
		return "", fmt.Errorf("creating policy: %w", err)
	}

	return policyName, nil
}

func (os *MinioObjectStorage) UploadImageSources(ctx context.Context, product, image string, sources []byte) error {
	os.logger.Info("Uploading image's sources", "product", product, "image", image)

	object := bytes.NewReader(sources)

	_, err := os.client.PutObject(ctx, product, image, object, object.Size(), minio.PutObjectOptions{})

	return err
}

func (os *MinioObjectStorage) DeleteImageSources(ctx context.Context, product, image string) error {
	os.logger.Info("Deleting image's sources", "product", product, "image", image)

	return os.client.RemoveObject(ctx, product, image, minio.RemoveObjectOptions{})
}

func (os *MinioObjectStorage) getPolicy(bucket string) []byte {
	return []byte(
		fmt.Sprintf(`{
			"Version":"2012-10-17",
			"Statement":[
				{
					"Effect":"Allow",
					"Action":["s3:*"],
					"Resource":["arn:aws:s3:::%s","arn:aws:s3:::%s/*"]
				}
			]
		}`, bucket, bucket))
}
