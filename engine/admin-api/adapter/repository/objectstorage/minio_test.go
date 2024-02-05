//go:build integration

package objectstorage_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/repository/objectstorage"
	"github.com/minio/madmin-go/v3"
	"github.com/minio/minio-go/v7"
	"github.com/sebdah/goldie/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	_testBucket = "test-bucket"
	_kaiPath    = ".kai"
)

type ObjectStorageSuite struct {
	suite.Suite

	minioContainer testcontainers.Container
	client         *minio.Client
	adminClient    *madmin.AdminClient
	objectStorage  *objectstorage.MinioObjectStorage
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
		Env: map[string]string{},
		WaitingFor: wait.ForAll(wait.ForLog("Status:         1 Online, 0 Offline."), wait.ForExposedPort()).
			WithDeadline(time.Minute * 3),
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

	viper.Set(config.MinioEndpointKey, minioEndpoint)
	viper.Set(config.MinioRootUserKey, "minioadmin")
	viper.Set(config.MinioRootPasswordKey, "minioadmin")

	client, err := objectstorage.NewMinioClient()
	s.Require().NoError(err)

	adminClient, err := objectstorage.NewAdminMinioClient()
	s.Require().NoError(err)

	logger := testr.NewWithOptions(s.T(), testr.Options{Verbosity: -1})

	s3ObjectStorage := objectstorage.NewMinioObjectStorage(logger, client, adminClient)

	s.minioContainer = minioContainer
	s.client = client
	s.adminClient = adminClient
	s.objectStorage = s3ObjectStorage
}

func (s *ObjectStorageSuite) TearDownSuite() {
	s.Require().NoError(s.minioContainer.Terminate(context.Background()))
}

func (s *ObjectStorageSuite) TearDownTest() {
	ctx := context.Background()

	err := s.client.RemoveBucketWithOptions(ctx, _testBucket, minio.RemoveBucketOptions{ForceDelete: true})
	if err != nil {
		var minioErr minio.ErrorResponse
		errors.As(err, &minioErr)
		if minioErr.StatusCode != http.StatusNotFound {
			s.Failf(err.Error(), "Error deleting bucket %q", _testBucket)
		}
	}

	policies, err := s.adminClient.ListCannedPolicies(ctx)
	s.Require().NoError(err)

	for name := range policies {
		err = s.adminClient.RemoveCannedPolicy(ctx, name)
		if err != nil && !strings.Contains(err.Error(), "inbuilt policy") {
			s.Failf(err.Error(), "Error deleting policy non inbuilt policy %q", name)
		}
	}
}

func (s *ObjectStorageSuite) TestCreateBucket() {
	ctx := context.Background()

	err := s.objectStorage.CreateBucket(ctx, _testBucket)
	s.Assert().NoError(err)

	bucketExists, err := s.client.BucketExists(ctx, _testBucket)
	s.Require().NoError(err)
	s.Assert().True(bucketExists)
}

func (s *ObjectStorageSuite) TestCreateBucket_WithLifecycle_ErrorTierDoesntExist() {
	ctx := context.Background()

	viper.Set(config.MinioTierEnabledKey, true)
	viper.Set(config.MinioTierNameKey, "TIER-NAME")

	err := s.objectStorage.CreateBucket(ctx, _testBucket)
	s.Assert().Error(err)

	viper.Set(config.MinioTierEnabledKey, false)
}

func (s *ObjectStorageSuite) TestCreateBucket_InvalidBucketName() {
	ctx := context.Background()

	err := s.objectStorage.CreateBucket(ctx, "invalid bucket")
	s.Assert().Error(err)
}

func (s *ObjectStorageSuite) TestCreateBucketPolicy() {
	var (
		ctx                = context.Background()
		expectedPolicyName = _testBucket
	)

	err := s.objectStorage.CreateBucket(ctx, _testBucket)
	s.Assert().NoError(err)

	policyName, err := s.objectStorage.CreateBucketPolicy(ctx, _testBucket)
	s.Require().NoError(err)
	s.Assert().Equal(expectedPolicyName, policyName)

	policy, err := s.adminClient.InfoCannedPolicy(ctx, policyName)
	s.Require().NoError(err)

	// Unmarshal the policy into a map to make it easier to work with
	var policyMap map[string]interface{}
	err = json.Unmarshal(policy, &policyMap)
	s.Require().NoError(err)

	// Sort the "Resource" array within each "Statement" block
	if statements, ok := policyMap["Statement"].([]interface{}); ok {
		for _, statement := range statements {
			if statementMap, ok := statement.(map[string]interface{}); ok {
				if resources, ok := statementMap["Resource"].([]string); ok {
					sort.Strings(resources)
				}
			}
		}
	}

	// Marshal the sorted policy back to a JSON string
	sortedPolicyBytes, err := json.MarshalIndent(policyMap, "", " ")
	s.Require().NoError(err)

	g := goldie.New(s.T())
	g.Assert(s.T(), "CreateBucketPolicy", sortedPolicyBytes)
}

func (s *ObjectStorageSuite) TestCreateBucketPolicy_InvalidResourceInPolicy() {
	var (
		ctx = context.Background()
	)

	_, err := s.objectStorage.CreateBucketPolicy(ctx, "")
	s.Require().Error(err)
}

func (s *ObjectStorageSuite) TestUploadImageSources() {
	var (
		ctx     = context.Background()
		product = _testBucket
		image   = "test-image"
		sources = []byte("this is a test")
	)

	err := s.client.MakeBucket(ctx, product, minio.MakeBucketOptions{})
	s.Require().NoError(err)

	err = s.objectStorage.UploadImageSources(ctx, product, image, sources)
	s.Assert().NoError(err)

	imagePath := path.Join(_kaiPath, image)

	actualImage, err := s.client.GetObject(ctx, product, imagePath, minio.GetObjectOptions{})
	s.Require().NoError(err)

	content, err := io.ReadAll(actualImage)
	s.Require().NoError(err)

	s.Assert().Equal(sources, content)
}

func (s *ObjectStorageSuite) TestDeleteImageSources() {
	var (
		ctx     = context.Background()
		product = _testBucket
		image   = "test-image"
		sources = []byte("this is a test")
	)

	err := s.client.MakeBucket(ctx, product, minio.MakeBucketOptions{})
	s.Require().NoError(err)

	err = s.objectStorage.UploadImageSources(ctx, product, image, sources)
	s.Assert().NoError(err)

	err = s.objectStorage.DeleteImageSources(ctx, product, image)
	s.Assert().NoError(err)

	objs := s.client.ListObjects(ctx, product, minio.ListObjectsOptions{})
	for obj := range objs {
		s.Fail("Found unexpected object", "Object %q was not deleted", obj.Key)
	}
}
