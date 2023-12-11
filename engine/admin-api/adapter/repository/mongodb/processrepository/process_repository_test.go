//go:build integration

package processrepository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/repository"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
)

var (
	_kaiProduct        = "kai"
	productID          = "productID"
	ownerID            = "ownerID"
	processVersion     = "v1.0.0"
	testRepoUploadDate = time.Now().Add(-time.Hour).Truncate(time.Millisecond).UTC()
)

type ProcessRepositoryTestSuite struct {
	suite.Suite
	mongoDBContainer testcontainers.Container
	mongoClient      *mongo.Client
	processRepo      *MongoDBProcessRepository
}

func TestProcessRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ProcessRepositoryTestSuite))
}

func (s *ProcessRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()
	logger := testr.NewWithOptions(s.T(), testr.Options{Verbosity: -1})

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

	s.mongoDBContainer = mongoDBContainer
	s.mongoClient = client
	s.processRepo = New(logger, client)

	err = s.processRepo.CreateIndexes(context.Background(), productID)
	s.Require().NoError(err)

	viper.Set(config.MongoDBKaiDatabaseKey, _kaiProduct)

	monkey.Patch(time.Now, func() time.Time {
		return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	})
}

func (s *ProcessRepositoryTestSuite) TearDownSuite() {
	monkey.UnpatchAll()
	s.Require().NoError(s.mongoDBContainer.Terminate(context.Background()))
}

func (s *ProcessRepositoryTestSuite) TearDownTest() {
	filter := bson.D{}

	_, err := s.mongoClient.Database(productID).
		Collection(registeredProcessesCollectionName).
		DeleteMany(context.Background(), filter)
	s.Require().NoError(err)
}

func (s *ProcessRepositoryTestSuite) TestCreate() {
	testRegisteredProcess := &entity.RegisteredProcess{
		ID:         "process_id",
		Name:       "test_trigger",
		Version:    processVersion,
		Type:       "trigger",
		Image:      "process_image",
		UploadDate: testRepoUploadDate,
		Owner:      ownerID,
	}

	createdRegisteredProcess, err := s.processRepo.Create(productID, testRegisteredProcess)
	s.Require().NoError(err)

	s.Equal(testRegisteredProcess, createdRegisteredProcess)

	// Check if the version is created in the DB
	collection := s.mongoClient.Database(productID).Collection(registeredProcessesCollectionName)
	filter := bson.M{"_id": createdRegisteredProcess.ID}

	var registeredProcessDTO registeredProcessDTO
	err = collection.FindOne(context.Background(), filter).Decode(&registeredProcessDTO)
	s.Require().NoError(err)
}

func (s *ProcessRepositoryTestSuite) TestSearchByProduct() {
	ctx := context.Background()

	testTriggerProcess := &entity.RegisteredProcess{
		ID:         "test_trigger_id",
		Name:       "test_trigger",
		Version:    processVersion,
		Type:       "trigger",
		Image:      "test_trigger_image",
		UploadDate: testRepoUploadDate,
		Owner:      ownerID,
	}

	testTriggerProcess2 := &entity.RegisteredProcess{
		ID:         "test_trigger_id_2",
		Name:       "test_trigger_2",
		Version:    processVersion,
		Type:       "trigger",
		Image:      "test_trigger_image_2",
		UploadDate: testRepoUploadDate,
		Owner:      ownerID,
	}

	testTaskProcess := &entity.RegisteredProcess{
		ID:         "test_task_id",
		Name:       "test_task",
		Version:    processVersion,
		Type:       "task",
		Image:      "test_task_image",
		UploadDate: testRepoUploadDate,
		Owner:      ownerID,
	}

	registeredProcesses := []*entity.RegisteredProcess{
		testTriggerProcess,
		testTriggerProcess2,
		testTaskProcess,
	}

	for _, p := range registeredProcesses {
		_, err := s.processRepo.Create(productID, p)
		s.Require().NoError(err)
	}

	registeredProcesses, err := s.processRepo.SearchByProduct(ctx, productID, repository.SearchFilter{
		ProcessType: entity.ProcessTypeTask,
	})
	s.Require().NoError(err)

	s.Require().Len(registeredProcesses, 1)
	s.Equal(testTaskProcess, registeredProcesses[0])

	registeredProcesses, err = s.processRepo.SearchByProduct(ctx, productID, repository.SearchFilter{})
	s.Require().NoError(err)

	s.Require().Len(registeredProcesses, 3)
}

func (s *ProcessRepositoryTestSuite) TestSearchByProductWithUnexistingProduct() {
	ctx := context.Background()

	registeredProcesses, err := s.processRepo.SearchByProduct(ctx, "non-existent", repository.SearchFilter{
		ProcessType: entity.ProcessTypeTask,
	})
	s.Require().NoError(err)

	s.Empty(registeredProcesses)
}

func (s *ProcessRepositoryTestSuite) TestGlobalSearch() {
	ctx := context.Background()

	testGlobalProcess := testhelpers.NewRegisteredProcessBuilder(_kaiProduct).Build()

	expectedProcesses := []*entity.RegisteredProcess{
		testGlobalProcess,
	}

	_, err := s.processRepo.Create(_kaiProduct, testGlobalProcess)
	s.Require().NoError(err)

	actualProcesses, err := s.processRepo.GlobalSearch(ctx, repository.SearchFilter{
		ProcessType: entity.ProcessTypeTask,
	})
	s.Require().NoError(err)

	s.Assert().Equal(expectedProcesses, actualProcesses)
}

func (s *ProcessRepositoryTestSuite) TestUpdate() {
	ctx := context.Background()

	testRegisteredProcess := &entity.RegisteredProcess{
		ID:         "process_id",
		Name:       "test_trigger",
		Version:    processVersion,
		Type:       "trigger",
		Image:      "process_image",
		UploadDate: testRepoUploadDate,
		Owner:      ownerID,
	}

	createdRegisteredProcess, err := s.processRepo.Create(productID, testRegisteredProcess)
	s.Require().NoError(err)

	createdRegisteredProcess.Image = "new_process_image"

	err = s.processRepo.Update(ctx, productID, createdRegisteredProcess)
	s.Require().NoError(err)

	// Check if the version is updated in the DB
	collection := s.mongoClient.Database(productID).Collection(registeredProcessesCollectionName)
	filter := bson.M{"_id": createdRegisteredProcess.ID}

	var registeredProcessDTO registeredProcessDTO
	err = collection.FindOne(context.Background(), filter).Decode(&registeredProcessDTO)
	s.Require().NoError(err)
	s.Equal(createdRegisteredProcess, mapDTOToEntity(&registeredProcessDTO))
}

func (s *ProcessRepositoryTestSuite) TestGetByID() {
	ctx := context.Background()

	testRegisteredProcess := &entity.RegisteredProcess{
		ID:         "process_id",
		Name:       "test_trigger",
		Version:    processVersion,
		Type:       "trigger",
		Image:      "process_image",
		UploadDate: testRepoUploadDate,
		Owner:      ownerID,
	}

	createdRegisteredProcess, err := s.processRepo.Create(productID, testRegisteredProcess)
	s.Require().NoError(err)

	registeredProcess, err := s.processRepo.GetByID(ctx, productID, createdRegisteredProcess.ID)
	s.Require().NoError(err)

	s.Equal(createdRegisteredProcess, registeredProcess)
}

func (s *ProcessRepositoryTestSuite) TestGetByID_NoResults() {
	ctx := context.Background()

	_, err := s.processRepo.GetByID(ctx, productID, "nonexistent")
	s.Require().Error(err)

	s.ErrorIs(err, usecase.ErrRegisteredProcessNotFound)
}
