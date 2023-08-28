//go:build integration

package processregistry

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

var (
	productID      = "productID"
	ownerID        = "ownerID"
	processVersion = "v1.0.0"
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
		Version:    processVersion,
		Type:       "trigger",
		Image:      "process_image",
		UploadDate: testUploadDate,
		Owner:      ownerID,
	}

	createdProcessRegistry, err := s.processRegistryRepo.Create(productID, testProcessRegistry)
	s.Require().NoError(err)

	s.Equal(testProcessRegistry, createdProcessRegistry)

	// Check if the version is created in the DB
	collection := s.mongoClient.Database(productID).Collection(processRegistryCollectionName)
	filter := bson.M{"_id": createdProcessRegistry.ID}

	var processRegistryDTO processRegistryDTO
	err = collection.FindOne(context.Background(), filter).Decode(&processRegistryDTO)
	s.Require().NoError(err)

	fmt.Println(testProcessRegistry.UploadDate)
	fmt.Println(processRegistryDTO.UploadDate)
}

func (s *ProcessRegistryRepositoryTestSuite) TestListByProductWithTypeFilter() {
	ctx := context.Background()

	testTriggerProcess := &entity.ProcessRegistry{
		ID:         "test_trigger_id",
		Name:       "test_trigger",
		Version:    processVersion,
		Type:       "trigger",
		Image:      "test_trigger_image",
		UploadDate: testUploadDate,
		Owner:      ownerID,
	}

	testTriggerProcess2 := &entity.ProcessRegistry{
		ID:         "test_trigger_id_2",
		Name:       "test_trigger_2",
		Version:    processVersion,
		Type:       "trigger",
		Image:      "test_trigger_image_2",
		UploadDate: testUploadDate,
		Owner:      ownerID,
	}

	testTaskProcess := &entity.ProcessRegistry{
		ID:         "test_task_id",
		Name:       "test_task",
		Version:    processVersion,
		Type:       "task",
		Image:      "test_task_image",
		UploadDate: testUploadDate,
		Owner:      ownerID,
	}

	processRegistries := []*entity.ProcessRegistry{
		testTriggerProcess,
		testTriggerProcess2,
		testTaskProcess,
	}

	for _, p := range processRegistries {
		_, err := s.processRegistryRepo.Create(productID, p)
		s.Require().NoError(err)
	}

	processRegistries, err := s.processRegistryRepo.ListByProductWithTypeFilter(ctx, productID, "task")
	s.Require().NoError(err)

	s.Require().Len(processRegistries, 1)
	s.Equal(testTaskProcess, processRegistries[0])

	processRegistries, err = s.processRegistryRepo.ListByProductWithTypeFilter(ctx, productID, "")
	s.Require().NoError(err)

	s.Require().Len(processRegistries, 3)
}

func (s *ProcessRegistryRepositoryTestSuite) TestListByProductWithUnexistingProduct() {
	ctx := context.Background()

	processRegistries, err := s.processRegistryRepo.ListByProductWithTypeFilter(ctx, "unexisting", "task")
	s.Require().NoError(err)

	s.Empty(processRegistries)
}
