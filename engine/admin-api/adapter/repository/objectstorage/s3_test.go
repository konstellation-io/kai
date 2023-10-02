//go:build integration

package objectstorage_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/repository/objectstorage"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type ObjectStorageSuite struct {
	suite.Suite

	minioContainer  testcontainers.Container
	client          *s3.S3
	s3ObjectStorage *objectstorage.S3ObjectStorage
}

func TestObjectStorageSuite(t *testing.T) {
	suite.Run(t, new(ObjectStorageSuite))
}

func (s *ObjectStorageSuite) SetupSuite() {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "minio/minio:latest",
		ExposedPorts: []string{"9000/tcp"},
		Cmd: []string{
			"server",
			"/data",
		},
		Env:        map[string]string{},
		WaitingFor: wait.ForLog("Status:         1 Online, 0 Offline."),
	}

	minioContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	s.Require().NoError(err)

	host, err := minioContainer.Host(ctx)
	s.Require().NoError(err)

	port, err := minioContainer.MappedPort(ctx, "9000/tcp")
	s.Require().NoError(err)

	minioEndpoint := fmt.Sprintf("http://%s:%d", host, port.Int())

	viper.Set(config.S3EndpointKey, minioEndpoint)
	viper.Set(config.S3BucketKey, "kai")

	err = os.Setenv("AWS_REGION", "us-east-1")
	s.Require().NoError(err)

	err = os.Setenv("AWS_ACCESS_KEY_ID", "minioadmin")
	s.Require().NoError(err)

	err = os.Setenv("AWS_SECRET_ACCESS_KEY", "minioadmin")
	s.Require().NoError(err)

	client, err := objectstorage.NewS3Client()
	s.Require().NoError(err)

	logger := testr.NewWithOptions(s.T(), testr.Options{Verbosity: -1})

	s3ObjectStorage := objectstorage.NewS3ObjectStorage(logger, client)

	_, err = client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(viper.GetString(config.S3BucketKey)),
	})
	s.Require().NoError(err)

	s.minioContainer = minioContainer
	s.client = client
	s.s3ObjectStorage = s3ObjectStorage
}

func (s *ObjectStorageSuite) TearDownSuite() {
	s.Require().NoError(s.minioContainer.Terminate(context.Background()))
}

func (s *ObjectStorageSuite) TearDownTest() {
	fmt.Println("tearing down test")

	objects, err := s.client.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(viper.GetString(config.S3BucketKey))})
	s.Require().NoError(err)

	objectsIdentifiers := make([]*s3.ObjectIdentifier, 0, len(objects.Contents))

	for _, object := range objects.Contents {
		objectsIdentifiers = append(objectsIdentifiers, &s3.ObjectIdentifier{
			Key: object.Key,
		})
	}

	_, err = s.client.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(viper.GetString(config.S3BucketKey)),
		Delete: &s3.Delete{
			Objects: objectsIdentifiers,
		},
	})
	s.Require().NoError(err)
}

func (s *ObjectStorageSuite) TestCreateFolder() {
	folderName := "test-folder"

	err := s.s3ObjectStorage.CreateFolder(folderName)
	s.Assert().NoError(err)

	foundObject, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(viper.GetString(config.S3BucketKey)),
		Key:    aws.String(folderName + "/"),
	})
	s.Require().NoError(err)

	fmt.Println(foundObject.String())
}
