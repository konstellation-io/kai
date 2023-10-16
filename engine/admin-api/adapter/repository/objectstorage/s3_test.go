//go:build integration

package objectstorage_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/repository/objectstorage"
	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	_testBucket = "test-bucket"
)

type ObjectStorageSuite struct {
	suite.Suite

	minioContainer  testcontainers.Container
	client          *minio.Client
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

	minioEndpoint := fmt.Sprintf("%s:%d", host, port.Int())

	viper.Set(config.S3EndpointKey, minioEndpoint)
	//viper.Set(config.S3BucketKey, "kai")

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

	s.minioContainer = minioContainer
	s.client = client
	s.s3ObjectStorage = s3ObjectStorage
}

func (s *ObjectStorageSuite) TearDownSuite() {
	s.Require().NoError(s.minioContainer.Terminate(context.Background()))
}

func (s *ObjectStorageSuite) TearDownTest() {
	ctx := context.Background()

	//objects := s.client.ListObjects(ctx, _testBucket, minio.ListObjectsOptions{})
	//
	//resCh := s.client.RemoveObjectsWithResult(ctx, _testBucket, objects, minio.RemoveObjectsOptions{})
	//for res := range resCh {
	//	s.Require().NoError(res.Err)
	//}

	err := s.client.RemoveBucketWithOptions(ctx, _testBucket, minio.RemoveBucketOptions{ForceDelete: true})
	s.Require().NoError(err)
}

func (s *ObjectStorageSuite) TestCreateBucket() {
	ctx := context.Background()

	err := s.s3ObjectStorage.CreateBucket(ctx, _testBucket)
	s.Assert().NoError(err)

	bucketExists, err := s.client.BucketExists(ctx, _testBucket)
	s.Require().NoError(err)
	s.Assert().True(bucketExists)
}
