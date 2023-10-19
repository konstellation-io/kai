package objectstorage

import (
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/minio/madmin-go/v3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

func NewMinioClient() (*minio.Client, error) {
	var (
		minioKeyID     = viper.GetString(config.MinioRootUserKey)
		minioKeySecret = viper.GetString(config.MinioRootPasswordKey)
	)

	client, err := minio.New(viper.GetString(config.MinioEndpointKey), &minio.Options{
		Creds:        credentials.NewStaticV4(minioKeyID, minioKeySecret, ""),
		BucketLookup: minio.BucketLookupPath,
	})
	if err != nil {
		return nil, fmt.Errorf("initializing MinIO client: %w", err)
	}

	return client, nil
}

func NewAdminMinioClient() (*madmin.AdminClient, error) {
	var (
		minioKeyID     = viper.GetString(config.MinioRootUserKey)
		minioKeySecret = viper.GetString(config.MinioRootPasswordKey)
	)

	client, err := madmin.New(viper.GetString(config.MinioEndpointKey), minioKeyID, minioKeySecret, false)
	if err != nil {
		return nil, fmt.Errorf("initializing MinIO admin client: %w", err)
	}

	return client, nil
}
