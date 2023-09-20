//go:build integration

package versionrepository

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/version"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/libs/simplelogger"
)

var productID = "productID"
var versionTag = "v1.0.0"
var creatorID = "creatorID"

type VersionRepositoryTestSuite struct {
	suite.Suite
	cfg              *config.Config
	mongoDBContainer testcontainers.Container
	mongoClient      *mongo.Client
	versionRepo      *VersionRepoMongoDB
}

func TestGocloakTestSuite(t *testing.T) {
	suite.Run(t, new(VersionRepositoryTestSuite))
}

func (s *VersionRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()
	cfg := &config.Config{}
	logger := simplelogger.New(simplelogger.LevelInfo)

	cfg.MongoDB.KRTBucket = "krt"

	req := testcontainers.ContainerRequest{
		Image:        "mongo:latest",
		ExposedPorts: []string{"27017/tcp", "27018/tcp"},
		Env: map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": "root",
			"MONGO_INITDB_ROOT_PASSWORD": "root",
		},
		WaitingFor: wait.ForLog("MongoDB starting"),
	}

	mongoDBContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	s.Require().NoError(err)

	host, err := mongoDBContainer.Host(context.Background())
	s.Require().NoError(err)
	p, err := mongoDBContainer.MappedPort(context.Background(), "27017/tcp")
	s.Require().NoError(err)

	port := p.Int()
	uri := fmt.Sprintf("mongodb://%v:%v@%v:%v/", "root", "root", host, port) //NOSONAR not used in secure contexts
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	s.Require().NoError(err)

	s.cfg = cfg
	s.mongoDBContainer = mongoDBContainer
	s.mongoClient = client
	s.versionRepo = New(cfg, logger, client)

	err = s.versionRepo.CreateIndexes(context.Background(), productID)
	s.Require().NoError(err)
}

func (s *VersionRepositoryTestSuite) TearDownSuite() {
	s.Require().NoError(s.mongoDBContainer.Terminate(context.Background()))
}

func (s *VersionRepositoryTestSuite) TearDownTest() {
	filter := bson.D{}

	_, err := s.mongoClient.Database(productID).
		Collection(versionsCollectionName).
		DeleteMany(context.Background(), filter)
	s.Require().NoError(err)
}

func (s *VersionRepositoryTestSuite) TestCreate() {
	testVersion := &entity.Version{
		Tag: versionTag,
	}

	createdVer, err := s.versionRepo.Create(creatorID, productID, testVersion)
	s.Require().NoError(err)

	s.NotEmpty(createdVer.Tag)
	s.NotEmpty(createdVer.CreationDate)
	s.Equal(creatorID, createdVer.CreationAuthor)
	s.Equal(entity.VersionStatusCreated, createdVer.Status)

	// Check if the version is created in the DB
	collection := s.mongoClient.Database(productID).Collection(versionsCollectionName)
	filter := bson.M{"tag": createdVer.Tag}

	var versionDTO versionDTO
	err = collection.FindOne(context.Background(), filter).Decode(&versionDTO)
	s.Require().NoError(err)
}

func (s *VersionRepositoryTestSuite) TestCreateDuplicateTagError() {
	testVersion := &entity.Version{
		Tag: versionTag,
	}

	duplicatedVersion := &entity.Version{
		Tag: versionTag,
	}

	_, err := s.versionRepo.Create(creatorID, productID, testVersion)
	s.Require().NoError(err)

	_, err = s.versionRepo.Create(creatorID, productID, duplicatedVersion)
	s.Require().Error(err)
}

func (s *VersionRepositoryTestSuite) TestGetByTag() {
	testVersion := &entity.Version{
		Tag: versionTag,
	}

	_, err := s.versionRepo.Create(creatorID, productID, testVersion)
	s.Require().NoError(err)

	ver, err := s.versionRepo.GetByTag(context.Background(), productID, testVersion.Tag)
	s.Require().NoError(err)

	s.Equal(testVersion.Tag, ver.Tag)
}

func (s *VersionRepositoryTestSuite) TestGetByTagNotFound() {
	_, err := s.versionRepo.GetByTag(context.Background(), productID, "notfound")
	s.Require().Error(err)
	s.True(errors.Is(err, version.ErrVersionNotFound))
}

func (s *VersionRepositoryTestSuite) TestUpdate() {
	testVersion := &entity.Version{
		Tag: versionTag,
	}
	ctx := context.Background()

	createdVer, err := s.versionRepo.Create(creatorID, productID, testVersion)
	s.Require().NoError(err)

	createdVer.Description = "updated description"

	err = s.versionRepo.Update(productID, createdVer)
	s.Require().NoError(err)

	updatedVer, err := s.versionRepo.GetByTag(ctx, productID, createdVer.Tag)
	s.Require().NoError(err)

	s.Equal(createdVer.Description, updatedVer.Description)
}

func (s *VersionRepositoryTestSuite) TestUpdateNotFound() {
	testVersion := &entity.Version{
		Tag: versionTag,
	}

	err := s.versionRepo.Update(productID, testVersion)
	s.Require().Error(err)
	s.True(errors.Is(err, version.ErrVersionNotFound))
}

func (s *VersionRepositoryTestSuite) TestListVersionsByProduct() {
	testVersion := &entity.Version{
		Tag: versionTag,
	}

	testVersion2 := &entity.Version{
		Tag: "v2.0.0",
	}

	_, err := s.versionRepo.Create(creatorID, productID, testVersion)
	s.Require().NoError(err)
	_, err = s.versionRepo.Create(creatorID, productID, testVersion2)
	s.Require().NoError(err)

	versions, err := s.versionRepo.ListVersionsByProduct(context.Background(), productID)
	s.Require().NoError(err)

	s.Require().Len(versions, 2)
	s.Equal(testVersion.Tag, versions[0].Tag)
	s.Equal(testVersion2.Tag, versions[1].Tag)
}

func (s *VersionRepositoryTestSuite) TestSetStatus() {
	testVersion := &entity.Version{
		Tag: versionTag,
	}

	createdVer, err := s.versionRepo.Create(creatorID, productID, testVersion)
	s.Require().NoError(err)

	err = s.versionRepo.SetStatus(context.Background(), productID, createdVer.Tag, entity.VersionStatusCreated)
	s.Require().NoError(err)

	updatedVer, err := s.versionRepo.GetByTag(context.Background(), productID, createdVer.Tag)
	s.Require().NoError(err)

	s.Equal(entity.VersionStatusCreated, updatedVer.Status)
}

func (s *VersionRepositoryTestSuite) TestSetStatusWithPreviousError() {
	testVersion := &entity.Version{
		Tag:   versionTag,
		Error: "dummy error",
	}
	ctx := context.Background()

	createdVer, err := s.versionRepo.Create(creatorID, productID, testVersion)
	s.Require().NoError(err)

	err = s.versionRepo.SetStatus(context.Background(), productID, createdVer.Tag, entity.VersionStatusCreated)
	s.Require().NoError(err)

	updatedVer, err := s.versionRepo.GetByTag(ctx, productID, createdVer.Tag)
	s.Require().NoError(err)

	s.Equal(entity.VersionStatusCreated, updatedVer.Status)
	s.Empty(updatedVer.Error)
}

func (s *VersionRepositoryTestSuite) TestSetStatusNotFound() {
	err := s.versionRepo.SetStatus(context.Background(), productID, "notfound", entity.VersionStatusCreated)
	s.Assert().ErrorIs(err, version.ErrVersionNotFound)
}

func (s *VersionRepositoryTestSuite) TestSetError() {
	testVersion := &entity.Version{
		Tag: versionTag,
	}
	ctx := context.Background()

	createdVer, err := s.versionRepo.Create(creatorID, productID, testVersion)
	s.Require().NoError(err)

	_, err = s.versionRepo.SetError(context.Background(), productID, createdVer, "error1")
	s.Require().NoError(err)

	updatedVer, err := s.versionRepo.GetByTag(ctx, productID, createdVer.Tag)
	s.Require().NoError(err)

	s.Equal("error1", updatedVer.Error)
}

func (s *VersionRepositoryTestSuite) TestSetErrorNotFound() {
	_, err := s.versionRepo.SetError(context.Background(), productID, &entity.Version{Tag: "notfound"}, "error1")
	s.Require().Error(err)
	s.True(errors.Is(err, version.ErrVersionNotFound))
}
