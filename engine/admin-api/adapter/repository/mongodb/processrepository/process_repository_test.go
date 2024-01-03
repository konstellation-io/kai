//go:build integration

package processrepository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/repository/mongodb/processrepository"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/repository"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/process"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	_kaiProduct                     = "kai"
	productID                       = "productID"
	ownerID                         = "ownerID"
	processVersion                  = "v1.0.0"
	registeredProcessCollectionName = "registered_processes"
)

var (
	testRepoUploadDate = time.Now().Add(-time.Hour).Truncate(time.Millisecond).UTC()
)

type ProcessRepositoryTestSuite struct {
	suite.Suite
	mongoDBContainer testcontainers.Container
	mongoClient      *mongo.Client
	processRepo      *processrepository.MongoDBProcessRepository
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
	s.processRepo = processrepository.New(logger, client)

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
		Collection(registeredProcessCollectionName).
		DeleteMany(context.Background(), filter)
	s.Require().NoError(err)
}

func (s *ProcessRepositoryTestSuite) TestCreate() {
	ctx := context.Background()
	expectedProcess := &entity.RegisteredProcess{
		ID:         "process_id",
		Name:       "test_trigger",
		Version:    processVersion,
		Type:       "trigger",
		Image:      "process_image",
		UploadDate: testRepoUploadDate,
		Owner:      ownerID,
	}

	err := s.processRepo.Create(ctx, productID, expectedProcess)
	s.Require().NoError(err)

	// Check if the version is created in the DB
	actualProcess, err := s.processRepo.GetByID(ctx, productID, expectedProcess.ID)
	s.Require().NoError(err)

	s.Equal(expectedProcess, actualProcess)

	//collection := s.mongoClient.Database(productID).Collection(registeredProcessesCollectionName)
	//filter := bson.M{"_id": createdRegisteredProcess.ID}
	//
	//var registeredProcessDTO registeredProcessDTO
	//err = collection.FindOne(context.Background(), filter).Decode(&registeredProcessDTO)
	//s.Require().NoError(err)
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
		err := s.processRepo.Create(ctx, productID, p)
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

	err := s.processRepo.Create(ctx, _kaiProduct, testGlobalProcess)
	s.Require().NoError(err)

	actualProcesses, err := s.processRepo.GlobalSearch(ctx, repository.SearchFilter{
		ProcessType: entity.ProcessTypeTask,
	})
	s.Require().NoError(err)

	s.Assert().Equal(expectedProcesses, actualProcesses)
}

func (s *ProcessRepositoryTestSuite) TestUpdate() {
	ctx := context.Background()

	expectedProcess := &entity.RegisteredProcess{
		ID:         "process_id",
		Name:       "test_trigger",
		Version:    processVersion,
		Type:       "trigger",
		Image:      "process_image",
		UploadDate: testRepoUploadDate,
		Owner:      ownerID,
	}

	err := s.processRepo.Create(ctx, productID, expectedProcess)
	s.Require().NoError(err)

	expectedProcess.Image = "new_process_image"

	err = s.processRepo.Update(ctx, productID, expectedProcess)
	s.Require().NoError(err)

	// Check if the version is updated in the DB
	actualProcess, err := s.processRepo.GetByID(ctx, productID, expectedProcess.ID)
	s.Require().NoError(err)

	s.Equal(expectedProcess, actualProcess)
}

func (s *ProcessRepositoryTestSuite) TestGetByID() {
	ctx := context.Background()

	expectedProcess := &entity.RegisteredProcess{
		ID:         "process_id",
		Name:       "test_trigger",
		Version:    processVersion,
		Type:       "trigger",
		Image:      "process_image",
		UploadDate: testRepoUploadDate,
		Owner:      ownerID,
	}

	err := s.processRepo.Create(ctx, productID, expectedProcess)
	s.Require().NoError(err)

	actualProcess, err := s.processRepo.GetByID(ctx, productID, expectedProcess.ID)
	s.Require().NoError(err)

	s.Equal(expectedProcess, actualProcess)
}

func (s *ProcessRepositoryTestSuite) TestGetByID_NoResults() {
	ctx := context.Background()

	_, err := s.processRepo.GetByID(ctx, productID, "nonexistent")
	s.Require().Error(err)

	s.ErrorIs(err, process.ErrRegisteredProcessNotFound)
}
