package objectstorage

import (
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
)

func NewS3Client() (*minio.Client, error) {
	client, err := minio.New(viper.GetString(config.S3EndpointKey), &minio.Options{
		//Creds:              nil,
		BucketLookup: minio.BucketLookupPath,
	})
	//&aws.Config{
	//	Endpoint:         aws.String(viper.GetString(config.S3EndpointKey)),
	//	S3ForcePathStyle: aws.Bool(true),
	//})
	if err != nil {
		return nil, fmt.Errorf("initializing S3 session: %w", err)
	}

	return client, nil
}
