//go:build integration

package version

import (
	"context"
	"fmt"
	"testing"

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
var versionName = "versionName"
var creatorID = "creatorID"

type VersionRepositoryTestSuite struct {
	suite.Suite
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

	s.mongoDBContainer = mongoDBContainer
	s.mongoClient = client
	s.versionRepo = NewVersionRepoMongoDB(cfg, logger, client)
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

func (s *VersionRepositoryTestSuite) TestCreateIndexes() {
	err := s.versionRepo.CreateIndexes(context.Background(), productID)

	s.Require().NoError(err)
}

func (s *VersionRepositoryTestSuite) TestCreate() {
	testVersion := &entity.Version{
		Name: versionName,
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

func (s *VersionRepositoryTestSuite) TestGetByID() {
	testVersion := &entity.Version{
		Name: versionName,
	}

	createdVer, err := s.versionRepo.Create(creatorID, productID, testVersion)
	s.Require().NoError(err)

	ver, err := s.versionRepo.GetByID(productID, createdVer.ID)
	s.Require().NoError(err)

	s.Equal(testVersion.Name, ver.Name)
}

func (s *VersionRepositoryTestSuite) TestGetByName() {
	testVersion := &entity.Version{
		Name: versionName,
	}

	_, err := s.versionRepo.Create(creatorID, productID, testVersion)
	s.Require().NoError(err)

	ver, err := s.versionRepo.GetByName(context.Background(), productID, testVersion.Name)
	s.Require().NoError(err)

	s.Equal(testVersion.Name, ver.Name)
}

func (s *VersionRepositoryTestSuite) TestUpdate() {
	testVersion := &entity.Version{
		Name: versionName,
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

func (s *VersionRepositoryTestSuite) TestListVersionsByProduct() {
	testVersion := &entity.Version{
		Name: versionName,
	}

	testVersion2 := &entity.Version{
		Name: versionName + "2",
	}

	_, err := s.versionRepo.Create(creatorID, productID, testVersion)
	s.Require().NoError(err)
	_, err = s.versionRepo.Create(creatorID, productID, testVersion2)
	s.Require().NoError(err)

	versions, err := s.versionRepo.ListVersionsByProduct(context.Background(), productID)
	s.Require().NoError(err)

	s.Require().Len(versions, 2)
	s.Equal(testVersion.Name, versions[0].Name)
	s.Equal(testVersion2.Name, versions[1].Name)
}

func (s *VersionRepositoryTestSuite) TestSetStatus() {
	testVersion := &entity.Version{
		Name: versionName,
	}

	createdVer, err := s.versionRepo.Create(creatorID, productID, testVersion)
	s.Require().NoError(err)

	err = s.versionRepo.SetStatus(context.Background(), productID, createdVer.ID, entity.VersionStatusCreated)
	s.Require().NoError(err)

	updatedVer, err := s.versionRepo.GetByID(productID, createdVer.ID)
	s.Require().NoError(err)

	s.Equal(entity.VersionStatusCreated, updatedVer.Status)
}

func (s *VersionRepositoryTestSuite) TestSetErrors() {
	testVersion := &entity.Version{
		Name: versionName,
	}

	createdVer, err := s.versionRepo.Create(creatorID, productID, testVersion)
	s.Require().NoError(err)

	_, err = s.versionRepo.SetErrors(context.Background(), productID, createdVer, []string{"error1", "error2"})
	s.Require().NoError(err)

	updatedVer, err := s.versionRepo.GetByID(productID, createdVer.ID)
	s.Require().NoError(err)

	s.Equal([]string{"error1", "error2"}, updatedVer.Errors)
}

func (s *VersionRepositoryTestSuite) TestClearPublishedVersion() {
	testVersion := &entity.Version{
		Name: versionName,
	}

	createdVersion, err := s.versionRepo.Create(creatorID, productID, testVersion)
	s.Require().NoError(err)

	err = s.versionRepo.SetStatus(context.Background(), productID, createdVersion.ID, entity.VersionStatusPublished)
	s.Require().NoError(err)

	oldPublishedVErsion, err := s.versionRepo.ClearPublishedVersion(context.Background(), productID)
	s.Require().NoError(err)

	s.Equal(testVersion.Name, oldPublishedVErsion.Name)
	s.Equal(entity.VersionStatusStarted, oldPublishedVErsion.Status)
}
