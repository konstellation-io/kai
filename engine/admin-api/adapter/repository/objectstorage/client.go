package objectstorage

import (
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

func NewS3Client() (*minio.Client, error) {
	client, err := minio.New(viper.GetString(config.MinioEndpointKey), &minio.Options{
		Creds:        credentials.NewEnvAWS(),
		BucketLookup: minio.BucketLookupPath,
	})
	if err != nil {
		return nil, fmt.Errorf("initializing S3 session: %w", err)
	}

	return client, nil
}
