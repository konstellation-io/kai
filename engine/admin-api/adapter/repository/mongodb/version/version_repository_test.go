//go:build integration

package version

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	apperrors "github.com/konstellation-io/kai/engine/admin-api/domain/usecase/errors"
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
	uri := fmt.Sprintf("mongodb://%v:%v@%v:%v/", "root", "root", host, port)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	s.Require().NoError(err)

	s.cfg = cfg
	s.mongoDBContainer = mongoDBContainer
	s.mongoClient = client
	s.versionRepo = NewVersionRepoMongoDB(cfg, logger, client)

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

	s.NotEmpty(createdVer.ID)
	s.NotEmpty(createdVer.CreationDate)
	s.Equal(creatorID, createdVer.CreationAuthor)
	s.Equal(entity.VersionStatusCreating, createdVer.Status)

	// Check if the version is created in the DB
	collection := s.mongoClient.Database(productID).Collection(versionsCollectionName)
	filter := bson.M{"_id": createdVer.ID}

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

func (s *VersionRepositoryTestSuite) TestGetByID() {
	testVersion := &entity.Version{
		Tag: versionTag,
	}

	createdVer, err := s.versionRepo.Create(creatorID, productID, testVersion)
	s.Require().NoError(err)

	ver, err := s.versionRepo.GetByID(productID, createdVer.ID)
	s.Require().NoError(err)

	s.Equal(testVersion.Tag, ver.Tag)
}

func (s *VersionRepositoryTestSuite) TestGetByIDNotFound() {
	_, err := s.versionRepo.GetByID(productID, "notfound")
	s.Require().Error(err)
	s.True(errors.Is(err, apperrors.ErrVersionNotFound))
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
	s.True(errors.Is(err, apperrors.ErrVersionNotFound))
}

func (s *VersionRepositoryTestSuite) TestUpdate() {
	testVersion := &entity.Version{
		Tag: versionTag,
	}

	createdVer, err := s.versionRepo.Create(creatorID, productID, testVersion)
	s.Require().NoError(err)

	createdVer.Description = "updated description"

	err = s.versionRepo.Update(productID, createdVer)
	s.Require().NoError(err)

	updatedVer, err := s.versionRepo.GetByID(productID, createdVer.ID)
	s.Require().NoError(err)

	s.Equal(createdVer.Description, updatedVer.Description)
}

func (s *VersionRepositoryTestSuite) TestUpdateNotFound() {
	testVersion := &entity.Version{
		Tag: versionTag,
	}

	err := s.versionRepo.Update(productID, testVersion)
	s.Require().Error(err)
	s.True(errors.Is(err, apperrors.ErrVersionNotFound))
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

	err = s.versionRepo.SetStatus(context.Background(), productID, createdVer.ID, entity.VersionStatusCreated)
	s.Require().NoError(err)

	updatedVer, err := s.versionRepo.GetByID(productID, createdVer.ID)
	s.Require().NoError(err)

	s.Equal(entity.VersionStatusCreated, updatedVer.Status)
}

func (s *VersionRepositoryTestSuite) TestSetStatusNotFound() {
	err := s.versionRepo.SetStatus(context.Background(), productID, "notfound", entity.VersionStatusCreated)
	s.Require().Error(err)
	s.True(errors.Is(err, apperrors.ErrVersionNotFound))
}

func (s *VersionRepositoryTestSuite) TestSetErrors() {
	testVersion := &entity.Version{
		Tag: versionTag,
	}

	createdVer, err := s.versionRepo.Create(creatorID, productID, testVersion)
	s.Require().NoError(err)

	_, err = s.versionRepo.SetErrors(context.Background(), productID, createdVer, []string{"error1", "error2"})
	s.Require().NoError(err)

	updatedVer, err := s.versionRepo.GetByID(productID, createdVer.ID)
	s.Require().NoError(err)

	s.Equal([]string{"error1", "error2"}, updatedVer.Errors)
}

func (s *VersionRepositoryTestSuite) TestSetErrorsNotFound() {
	_, err := s.versionRepo.SetErrors(context.Background(), productID, &entity.Version{ID: "notfound"}, []string{"error1", "error2"})
	s.Require().Error(err)
	s.True(errors.Is(err, apperrors.ErrVersionNotFound))
}

func (s *VersionRepositoryTestSuite) TestUploadKRTYamlFile() {
	testVersion := &entity.Version{
		Tag: versionTag,
	}

	createdVer, err := s.versionRepo.Create(creatorID, productID, testVersion)
	s.Require().NoError(err)

	err = s.versionRepo.UploadKRTYamlFile(productID, createdVer, "../../../../testdata/classificator_krt.yaml")
	s.Require().NoError(err)

	bucket, err := gridfs.NewBucket(
		s.mongoClient.Database(productID),
		options.GridFSBucket().SetName(s.cfg.MongoDB.KRTBucket),
	)

	filter := bson.M{"_id": createdVer.ID}
	cursor, err := bucket.Find(filter)
	s.Require().NoError(err)

	defer func() {
		err := cursor.Close(context.TODO())
		s.Require().NoError(err)
	}()

	type gridfsFile struct {
		Name   string `bson:"filename"`
		Length int64  `bson:"length"`
	}
	var foundFiles []gridfsFile
	err = cursor.All(context.TODO(), &foundFiles)
	s.Require().NoError(err)

	s.Len(foundFiles, 1)
}

func (s *VersionRepositoryTestSuite) TestClearPublishedVersion() {
	testVersion := &entity.Version{
		Tag: versionTag,
	}

	createdVersion, err := s.versionRepo.Create(creatorID, productID, testVersion)
	s.Require().NoError(err)

	err = s.versionRepo.SetStatus(context.Background(), productID, createdVersion.ID, entity.VersionStatusPublished)
	s.Require().NoError(err)

	oldPublishedVErsion, err := s.versionRepo.ClearPublishedVersion(context.Background(), productID)
	s.Require().NoError(err)

	s.Equal(testVersion.Tag, oldPublishedVErsion.Tag)
	s.Equal(entity.VersionStatusStarted, oldPublishedVErsion.Status)
}
