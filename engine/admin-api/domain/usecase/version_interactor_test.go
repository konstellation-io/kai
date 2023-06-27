//go:build unit

package usecase_test

import (
	"context"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/errors"
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
	"github.com/konstellation-io/krt/pkg/krt"
)

type versionSuiteMocks struct {
	cfg              *config.Config
	logger           *mocks.MockLogger
	versionRepo      *mocks.MockVersionRepo
	productRepo      *mocks.MockProductRepo
	k8sService       *mocks.MockK8sService
	userActivityRepo *mocks.MockUserActivityRepo
	accessControl    *mocks.MockAccessControl
	dashboardService *mocks.MockDashboardService
}

type VersionInteractorSuite struct {
	suite.Suite
	ctrl              *gomock.Controller
	mocks             versionSuiteMocks
	versionInteractor *usecase.VersionInteractor
	ctx               context.Context
}

func TestVersionInteractorSuite(t *testing.T) {
	suite.Run(t, new(VersionInteractorSuite))
}

// SetupSuite will create a mock controller and will initialize all required mock interfaces.
func (s *VersionInteractorSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())

	cfg := &config.Config{}
	logger := mocks.NewMockLogger(ctrl)
	versionRepo := mocks.NewMockVersionRepo(ctrl)
	productRepo := mocks.NewMockProductRepo(ctrl)
	k8sService := mocks.NewMockK8sService(ctrl)
	natsManagerService := mocks.NewMockNatsManagerService(ctrl)
	userActivityRepo := mocks.NewMockUserActivityRepo(ctrl)
	accessControl := mocks.NewMockAccessControl(ctrl)
	dashboardService := mocks.NewMockDashboardService(ctrl)
	processLogRepo := mocks.NewMockProcessLogRepository(ctrl)

	mocks.AddLoggerExpects(logger)

	userActivityInteractor := usecase.NewUserActivityInteractor(
		logger,
		userActivityRepo,
		accessControl,
	)

	versionInteractor := usecase.NewVersionInteractor(
		cfg, logger, versionRepo, productRepo, k8sService, natsManagerService,
		userActivityInteractor, accessControl, dashboardService, processLogRepo)

	s.ctrl = ctrl
	s.mocks = versionSuiteMocks{
		cfg,
		logger,
		versionRepo,
		productRepo,
		k8sService,
		userActivityRepo,
		accessControl,
		dashboardService,
	}
	s.versionInteractor = versionInteractor
	s.ctx = context.Background()
}

// TearDownSuite finish controller.
func (s *VersionInteractorSuite) TearDownSuite() {
	s.ctrl.Finish()
}

func (s *VersionInteractorSuite) getTestVersion() *entity.Version {
	commonObjectStore := &entity.ProcessObjectStore{
		Name:  "emails",
		Scope: "workflow",
	}

	return &entity.Version{
		ID:          "", // ID to be given after create
		Name:        "email-classificator",
		Description: "Email classificator for branching features.",
		Version:     "v1.0.0",
		Config: []entity.ConfigurationVariable{
			{
				Key:   "keyA",
				Value: "value1",
			},
			{
				Key:   "keyB",
				Value: "value2",
			},
		},
		Workflows: []entity.Workflow{
			{
				Name: "go-classificator",
				Type: "data",
				Config: []entity.ConfigurationVariable{
					{
						Key:   "keyA",
						Value: "value1",
					},
					{
						Key:   "keyB",
						Value: "value2",
					},
				},
				Processes: []entity.Process{
					{
						Name:          "entrypoint",
						Type:          "trigger",
						Image:         "konstellation/kai-grpc-trigger:latest",
						Replicas:      krt.DefaultNumberOfReplicas,
						GPU:           krt.DefaultGPUValue,
						Subscriptions: []string{"exitpoint"},
						Networking: &entity.ProcessNetworking{
							TargetPort:      9000,
							DestinationPort: 9000,
							Protocol:        "TCP",
						},
					},
					{
						Name:          "etl",
						Type:          "task",
						Image:         "konstellation/kai-etl-task:latest",
						Replicas:      krt.DefaultNumberOfReplicas,
						GPU:           krt.DefaultGPUValue,
						ObjectStore:   commonObjectStore,
						Subscriptions: []string{"entrypoint"},
					},
					{
						Name:          "email-classificator",
						Type:          "task",
						Image:         "konstellation/kai-ec-task:latest",
						Replicas:      krt.DefaultNumberOfReplicas,
						GPU:           krt.DefaultGPUValue,
						ObjectStore:   commonObjectStore,
						Subscriptions: []string{"etl"},
					},
					{
						Name:          "exitpoint",
						Type:          "exit",
						Image:         "konstellation/kai-exitpoint:latest",
						Replicas:      krt.DefaultNumberOfReplicas,
						GPU:           krt.DefaultGPUValue,
						ObjectStore:   commonObjectStore,
						Subscriptions: []string{"etl", "stats-storer"},
					},
				},
			},
		},
	}
}

func (s *VersionInteractorSuite) TestCreateNewVersion() {
	user := testhelpers.NewUserBuilder().Build()
	productID := "product-1"

	product := &entity.Product{
		ID: productID,
	}

	testVersion := s.getTestVersion()

	file, err := os.Open("../../testdata/classificator_krt.yaml")
	s.Require().NoError(err)

	s.mocks.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActCreateVersion)
	s.mocks.productRepo.EXPECT().GetByID(s.ctx, productID).Return(product, nil)
	s.mocks.versionRepo.EXPECT().GetByName(s.ctx, productID, testVersion.Name).Return(nil, errors.ErrVersionNotFound)
	s.mocks.versionRepo.EXPECT().Create(user.ID, productID, gomock.Any()).Return(testVersion, nil)
	s.mocks.versionRepo.EXPECT().SetStatus(s.ctx, productID, testVersion.ID, entity.VersionStatusCreated).Return(nil)
	s.mocks.versionRepo.EXPECT().UploadKRTFile(productID, testVersion, gomock.Any()).Return(nil)
	s.mocks.userActivityRepo.EXPECT().Create(gomock.Any()).Return(nil)

	_, statusCh, err := s.versionInteractor.Create(context.Background(), user, productID, file)
	s.Require().NoError(err)

	actual := <-statusCh
	expected := testVersion
	expected.Status = entity.VersionStatusCreated
	s.Equal(expected, actual)
}

func (s *VersionInteractorSuite) TestCreateNewVersion_FailsIfVersionNameIsDuplicated() {
	productID := "product-1"

	user := testhelpers.NewUserBuilder().Build()

	product := &entity.Product{
		ID: productID,
	}

	testVersion := s.getTestVersion()

	file, err := os.Open("../../testdata/classificator_krt.yaml")
	s.Require().NoError(err)

	s.mocks.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActCreateVersion)
	s.mocks.productRepo.EXPECT().GetByID(s.ctx, productID).Return(product, nil)
	s.mocks.versionRepo.EXPECT().GetByName(s.ctx, productID, testVersion.Name).Return(testVersion, nil)

	_, _, err = s.versionInteractor.Create(context.Background(), user, productID, file)
	s.ErrorIs(err, errors.ErrVersionDuplicated)
}

func (s *VersionInteractorSuite) TestGetByName() {
	productID := "product-1"

	user := testhelpers.NewUserBuilder().Build()
	testVersion := s.getTestVersion()

	s.mocks.versionRepo.EXPECT().GetByName(s.ctx, productID, testVersion.Name).Return(testVersion, nil)

	actual, err := s.versionInteractor.GetByName(s.ctx, user, productID, testVersion.Name)
	s.Require().NoError(err)

	s.Equal(testVersion, actual)
}
