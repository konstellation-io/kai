//go:build integration

package processregistry

import (
	"context"
	"fmt"
	"testing"
	"time"

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

var (
	productID = "productID"
	ownerID   = "ownerID"
)

type ProcessRegistryRepositoryTestSuite struct {
	suite.Suite
	cfg                 *config.Config
	mongoDBContainer    testcontainers.Container
	mongoClient         *mongo.Client
	processRegistryRepo *ProcessRegistryRepoMongoDB
}

func TestGocloakTestSuite(t *testing.T) {
	suite.Run(t, new(ProcessRegistryRepositoryTestSuite))
}

func (s *ProcessRegistryRepositoryTestSuite) SetupSuite() {
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
	s.processRegistryRepo = NewProcessRegistryRepoMongoDB(cfg, logger, client)

	err = s.processRegistryRepo.CreateIndexes(context.Background(), productID)
	s.Require().NoError(err)
}

func (s *ProcessRegistryRepositoryTestSuite) TearDownSuite() {
	s.Require().NoError(s.mongoDBContainer.Terminate(context.Background()))
}

func (s *ProcessRegistryRepositoryTestSuite) TearDownTest() {
	filter := bson.D{}

	_, err := s.mongoClient.Database(productID).
		Collection(processRegistryCollectionName).
		DeleteMany(context.Background(), filter)
	s.Require().NoError(err)
}

func (s *ProcessRegistryRepositoryTestSuite) TestCreate() {
	testProcessRegistry := &entity.ProcessRegistry{
		ID:         "process_id",
		Name:       "test_trigger",
		Version:    "v1.0.0",
		Type:       "trigger",
		Image:      "process_image",
		UploadDate: time.Now().Add(-time.Hour),
		Owner:      ownerID,
	}

	createdProcessRegistry, err := s.processRegistryRepo.Create(productID, testProcessRegistry)
	s.Require().NoError(err)

	s.Equal(testProcessRegistry.ID, createdProcessRegistry.ID)
	s.Equal(testProcessRegistry.Name, createdProcessRegistry.Name)
	s.Equal(testProcessRegistry.Version, createdProcessRegistry.Version)
	s.Equal(testProcessRegistry.Type, createdProcessRegistry.Type)
	s.Equal(testProcessRegistry.Image, createdProcessRegistry.Image)
	s.Equal(testProcessRegistry.UploadDate, createdProcessRegistry.UploadDate)
	s.Equal(testProcessRegistry.Owner, createdProcessRegistry.Owner)

	// Check if the version is created in the DB
	collection := s.mongoClient.Database(productID).Collection(processRegistryCollectionName)
	filter := bson.M{"_id": createdProcessRegistry.ID}

	var processRegistryDTO processRegistryDTO
	err = collection.FindOne(context.Background(), filter).Decode(&processRegistryDTO)
	s.Require().NoError(err)
}
