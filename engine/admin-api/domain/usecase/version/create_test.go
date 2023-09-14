//go:build unit

package version_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/version"
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
	"github.com/stretchr/testify/suite"
)

type createVersionSuite struct {
	suite.Suite

	logger      logr.Logger
	versionRepo *mocks.MockVersionRepo
	productRepo *mocks.MockProductRepo
	//versionService   *mocks.MockVersionService
	userActivity  *mocks.MockUserActivityInteracter
	accessControl *mocks.MockAccessControl
	handler       *version.Handler
}

func TestCreateVersionSuite(t *testing.T) {
	suite.Run(t, new(createVersionSuite))
}

func (s *createVersionSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())

	s.logger = testr.NewWithOptions(s.T(), testr.Options{Verbosity: -1})
	s.versionRepo = mocks.NewMockVersionRepo(ctrl)
	s.productRepo = mocks.NewMockProductRepo(ctrl)
	s.userActivity = mocks.NewMockUserActivityInteracter(ctrl)
	s.accessControl = mocks.NewMockAccessControl(ctrl)

	natsManagerService := mocks.NewMockNatsManagerService(ctrl)
	versionService := mocks.NewMockVersionService(ctrl)
	processLogRepo := mocks.NewMockProcessLogRepository(ctrl)

	s.handler = version.NewHandler(
		s.logger,
		s.versionRepo,
		s.productRepo,
		versionService,
		natsManagerService,
		s.userActivity,
		s.accessControl,
		processLogRepo)
}

func (s *createVersionSuite) TestCreateVersion() {
	var (
		ctx             = context.Background()
		user            = testhelpers.NewUserBuilder().Build()
		expectedVersion = getClassificatorVersion()
		product         = &entity.Product{
			ID: "test-product",
		}
	)

	file, err := os.Open("../../../testdata/classificator_krt.yaml")
	s.Require().NoError(err)

	defer file.Close()

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActCreateVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(product, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, product.ID, expectedVersion.Tag).Return(nil, version.ErrVersionNotFound)
	s.versionRepo.EXPECT().Create(user.Email, product.ID, expectedVersion).Return(expectedVersion, nil)
	s.userActivity.EXPECT().RegisterCreateAction(user.Email, product.ID, expectedVersion).Return(nil)

	createdVersion, err := s.handler.Create(ctx, user, product.ID, file)
	s.Require().NoError(err)

	s.Assert().Equal(expectedVersion, createdVersion)
}

func (s *createVersionSuite) TestCreateVersion_FailsIfUserIsNotAuthorized() {
	var (
		ctx     = context.Background()
		user    = testhelpers.NewUserBuilder().Build()
		product = &entity.Product{
			ID: "test-product",
		}
	)

	file, err := os.Open("../../../testdata/classificator_krt.yaml")
	s.Require().NoError(err)

	defer file.Close()

	expectedError := errors.New("unauthorized")

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActCreateVersion).Return(expectedError)

	_, err = s.handler.Create(ctx, user, product.ID, file)
	s.Require().ErrorIs(err, expectedError)
}

func (s *createVersionSuite) TestCreateVersion_FailsIfProductNotFound() {
	var (
		ctx     = context.Background()
		user    = testhelpers.NewUserBuilder().Build()
		product = &entity.Product{
			ID: "test-product",
		}
	)

	file, err := os.Open("../../../testdata/classificator_krt.yaml")
	s.Require().NoError(err)

	defer file.Close()

	expectedError := errors.New("product not found")

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActCreateVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(nil, expectedError)

	_, err = s.handler.Create(ctx, user, product.ID, file)
	s.Require().ErrorIs(err, expectedError)
}

func (s *createVersionSuite) TestCreateVersion_FailsIfThereIsAnErrorCreatingInRepo() {
	var (
		ctx        = context.Background()
		user       = testhelpers.NewUserBuilder().Build()
		newVersion = getClassificatorVersion()
		product    = &entity.Product{
			ID: "test-product",
		}
	)

	file, err := os.Open("../../../testdata/classificator_krt.yaml")
	s.Require().NoError(err)

	defer file.Close()

	expectedError := errors.New("error creating version")

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActCreateVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(product, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, product.ID, newVersion.Tag).Return(nil, version.ErrVersionNotFound)
	s.versionRepo.EXPECT().Create(user.Email, product.ID, newVersion).Return(nil, expectedError)

	_, err = s.handler.Create(ctx, user, product.ID, file)
	s.Require().ErrorIs(err, expectedError)
}

func (s *createVersionSuite) TestCreateVersion_FailsIfVersionTagIsDuplicated() {
	var (
		ctx        = context.Background()
		user       = testhelpers.NewUserBuilder().Build()
		newVersion = getClassificatorVersion()
		product    = &entity.Product{
			ID: "test-product",
		}
	)

	file, err := os.Open("../../../testdata/classificator_krt.yaml")
	s.Require().NoError(err)

	defer file.Close()

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActCreateVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(product, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, product.ID, newVersion.Tag).Return(newVersion, nil)

	_, err = s.handler.Create(ctx, user, product.ID, file)
	s.Require().ErrorIs(err, version.ErrVersionDuplicated)
}

func getClassificatorVersion() *entity.Version {
	return &entity.Version{
		Tag:         "v1.0.0",
		Description: "Email classificator for branching features.",
		Config: []entity.ConfigurationVariable{
			{
				Key:   "keyA",
				Value: "value1",
			},
		},
		Workflows: []entity.Workflow{
			{
				Name: "go-classificator",
				Type: entity.WorkflowTypeData,
				//Stream:    "",
				Config: []entity.ConfigurationVariable{
					{
						Key:   "keyA",
						Value: "value1",
					},
				},
				Processes: []entity.Process{
					{
						Name:          "entrypoint",
						Type:          entity.ProcessTypeTrigger,
						Image:         "konstellation/kai-grpc-trigger:latest",
						Subscriptions: []string{"exitpoint"},
						Replicas:      int32(1),
						Networking: &entity.ProcessNetworking{
							TargetPort:      9000,
							DestinationPort: 9000,
							Protocol:        "TCP",
						},
						ResourceLimits: &entity.ProcessResourceLimits{
							CPU: &entity.ResourceLimit{
								Request: "100m",
								Limit:   "200m",
							},
							Memory: &entity.ResourceLimit{
								Request: "100Mi",
								Limit:   "200Mi",
							},
						},
						Status: entity.RegisterProcessStatusCreated,
					},
					{
						Name:          "etl",
						Type:          entity.ProcessTypeTask,
						Image:         "konstellation/kai-etl-task:latest",
						Subscriptions: []string{"entrypoint"},
						Replicas:      int32(1),
						ObjectStore: &entity.ProcessObjectStore{
							Name:  "emails",
							Scope: "workflow",
						},
						ResourceLimits: &entity.ProcessResourceLimits{
							CPU: &entity.ResourceLimit{
								Request: "100m",
								Limit:   "200m",
							},
							Memory: &entity.ResourceLimit{
								Request: "100Mi",
								Limit:   "200Mi",
							},
						},
						Status: entity.RegisterProcessStatusCreated,
					},
					{
						Name:          "email-classificator",
						Type:          entity.ProcessTypeTask,
						Image:         "konstellation/kai-ec-task:latest",
						Subscriptions: []string{"etl"},
						Replicas:      int32(1),
						ObjectStore: &entity.ProcessObjectStore{
							Name:  "emails",
							Scope: "workflow",
						},
						ResourceLimits: &entity.ProcessResourceLimits{
							CPU: &entity.ResourceLimit{
								Request: "100m",
								Limit:   "200m",
							},
							Memory: &entity.ResourceLimit{
								Request: "100Mi",
								Limit:   "200Mi",
							},
						},
						Status: entity.RegisterProcessStatusCreated,
					},
					{
						Name:          "exitpoint",
						Type:          entity.ProcessTypeExit,
						Image:         "konstellation/kai-exitpoint:latest",
						Subscriptions: []string{"etl", "email-classificator"},
						Replicas:      int32(1),
						ObjectStore: &entity.ProcessObjectStore{
							Name:  "emails",
							Scope: "workflow",
						},
						ResourceLimits: &entity.ProcessResourceLimits{
							CPU: &entity.ResourceLimit{
								Request: "100m",
								Limit:   "200m",
							},
							Memory: &entity.ResourceLimit{
								Request: "100Mi",
								Limit:   "200Mi",
							},
						},
						Status: entity.RegisterProcessStatusCreated,
					},
				},
			},
		},
		Status: entity.VersionStatusCreated,
	}
}
