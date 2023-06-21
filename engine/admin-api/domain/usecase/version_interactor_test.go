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
	versionService   *mocks.MockVersionService
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
	versionService := mocks.NewMockVersionService(ctrl)
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
		cfg, logger, versionRepo, productRepo, versionService, natsManagerService,
		userActivityInteractor, accessControl, dashboardService, processLogRepo)

	s.ctrl = ctrl
	s.mocks = versionSuiteMocks{
		cfg,
		logger,
		versionRepo,
		productRepo,
		versionService,
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
	return &entity.Version{
		ID:  "version-id",
		Krt: s.getClassificatorKRT(),
	}
}

func (s *VersionInteractorSuite) getClassificatorKRT() *krt.Krt {
	defaultReplicas := krt.DefaultNumberOfReplicas
	defaultGPU := krt.DefaultGPUValue
	commonObjectStore := &krt.ProcessObjectStore{
		Name:  "emails",
		Scope: krt.ObjectStoreScopeWorkflow,
	}

	return &krt.Krt{
		Name:        "email-classificator",
		Description: "Email classificator for branching features.",
		Version:     "v1.0.0",
		Config: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
		Workflows: []krt.Workflow{
			{
				Name: "go-classificator",
				Type: krt.WorkflowTypeData,
				Config: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
				Processes: []krt.Process{
					{
						Name:          "entrypoint",
						Type:          krt.ProcessTypeTrigger,
						Image:         "konstellation/kai-grpc-trigger:latest",
						Replicas:      &defaultReplicas,
						GPU:           &defaultGPU,
						Subscriptions: []string{"exitpoint"},
						Networking: &krt.ProcessNetworking{
							TargetPort:          9000,
							TargetProtocol:      krt.NetworkingProtocolTCP,
							DestinationPort:     9000,
							DestinationProtocol: krt.NetworkingProtocolTCP,
						},
					},
					{
						Name:          "etl",
						Type:          krt.ProcessTypeTask,
						Image:         "konstellation/kai-etl-task:latest",
						Replicas:      &defaultReplicas,
						GPU:           &defaultGPU,
						ObjectStore:   commonObjectStore,
						Subscriptions: []string{"entrypoint"},
					},
					{
						Name:          "email-classificator",
						Type:          krt.ProcessTypeTask,
						Image:         "konstellation/kai-ec-task:latest",
						Replicas:      &defaultReplicas,
						GPU:           &defaultGPU,
						ObjectStore:   commonObjectStore,
						Subscriptions: []string{"etl"},
					},
					{
						Name:          "exitpoint",
						Type:          krt.ProcessTypeExit,
						Image:         "konstellation/kai-exitpoint:latest",
						Replicas:      &defaultReplicas,
						GPU:           &defaultGPU,
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
	krtClassificator := s.getClassificatorKRT()

	file, err := os.Open("../../test_assets/classificator_krt.yaml")
	s.Require().NoError(err)

	s.mocks.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActCreateVersion)
	s.mocks.productRepo.EXPECT().GetByID(s.ctx, productID).Return(product, nil)
	s.mocks.versionRepo.EXPECT().GetByName(s.ctx, productID, testVersion.Name).Return(nil, errors.ErrVersionNotFound)
	s.mocks.versionRepo.EXPECT().SetStatus(s.ctx, productID, testVersion.ID, entity.VersionStatusCreated).Return(nil)
	s.mocks.versionRepo.EXPECT().UploadKRTFile(productID, testVersion, gomock.Any()).Return(nil)
	s.mocks.userActivityRepo.EXPECT().Create(gomock.Any()).Return(nil)
	s.mocks.versionRepo.EXPECT().Create(
		user.ID,
		productID,
		&entity.Version{
			Krt: krtClassificator,
		},
	).Return(testVersion, nil)

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

	file, err := os.Open("../../test_assets/classificator_krt.yaml")
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
