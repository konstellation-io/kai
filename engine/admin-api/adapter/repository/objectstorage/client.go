package objectstorage

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/spf13/viper"
)

func NewS3Client() (*s3.S3, error) {
	newSession, err := session.NewSession(&aws.Config{
		Endpoint:         aws.String(viper.GetString(config.S3EndpointKey)),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("initializing S3 session: %w", err)
	}

	return s3.New(newSession), nil
}
