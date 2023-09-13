//go:build unit

package version_test

import (
	"context"
	"os"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/krt"
	"github.com/konstellation-io/kai/engine/admin-api/internal/errors"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
)

const KRT_PATH = "../../../testdata/classificator_krt.yaml"
const PRODUCT_ID = "product-1"

func (s *VersionUsecaseTestSuite) getTestVersion() *entity.Version {
	commonObjectStore := &entity.ProcessObjectStore{
		Name:  "emails",
		Scope: "workflow",
	}

	return &entity.Version{
		ID:          "", // ID to be given after create
		Tag:         "v1.0.0",
		Description: "Email classificator for branching features.",
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

func (s *VersionUsecaseTestSuite) TestCreateNewVersion() {
	user := testhelpers.NewUserBuilder().Build()
	productID := PRODUCT_ID
	ctx := context.Background()
	product := &entity.Product{
		ID: productID,
	}

	testVersion := s.getTestVersion()

	file, err := os.Open(KRT_PATH)
	s.Require().NoError(err)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActCreateVersion)
	s.productRepo.EXPECT().GetByID(ctx, productID).Return(product, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, testVersion.Tag).Return(nil, errors.ErrVersionNotFound)
	s.versionRepo.EXPECT().Create(user.ID, productID, gomock.Any()).Return(testVersion, nil)
	s.versionRepo.EXPECT().UploadKRTYamlFile(productID, testVersion, gomock.Any()).Return(nil)
	s.versionRepo.EXPECT().SetStatus(ctx, productID, testVersion.ID, entity.VersionStatusCreated).Return(nil)
	s.userActivityInteractor.EXPECT().RegisterCreateAction(user.ID, productID, testVersion).Return(nil)

	_, statusCh, err := s.handler.Create(ctx, user, productID, file)
	s.Require().NoError(err)

	actual := <-statusCh
	expected := testVersion
	expected.Status = entity.VersionStatusCreated
	s.Equal(expected, actual)
}

func (s *VersionUsecaseTestSuite) TestCreateNewVersion_FailsIfVersionTagIsDuplicated() {
	productID := PRODUCT_ID
	user := testhelpers.NewUserBuilder().Build()
	ctx := context.Background()
	product := &entity.Product{
		ID: productID,
	}

	testVersion := s.getTestVersion()

	file, err := os.Open(KRT_PATH)
	s.Require().NoError(err)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActCreateVersion)
	s.productRepo.EXPECT().GetByID(ctx, productID).Return(product, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, testVersion.Tag).Return(testVersion, nil)

	_, _, err = s.handler.Create(context.Background(), user, productID, file)
	s.ErrorIs(err, errors.ErrVersionDuplicated)
}

func (s *VersionUsecaseTestSuite) TestCreateNewVersion_FailsIfProductNotFound() {
	productID := PRODUCT_ID
	user := testhelpers.NewUserBuilder().Build()
	ctx := context.Background()
	file, err := os.Open(KRT_PATH)
	s.Require().NoError(err)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActCreateVersion)
	s.productRepo.EXPECT().GetByID(ctx, productID).Return(nil, usecase.ErrProductNotFound)

	_, _, err = s.handler.Create(context.Background(), user, productID, file)
	s.ErrorIs(err, usecase.ErrProductNotFound)
}

func (s *VersionUsecaseTestSuite) TestCreateNewVersion_FailsIfKrtIsInvalid() {
	productID := PRODUCT_ID
	user := testhelpers.NewUserBuilder().Build()
	ctx := context.Background()
	product := &entity.Product{
		ID: productID,
	}

	file, err := os.Open("./testdata/invalid_krt.yaml")
	s.Require().NoError(err)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActCreateVersion)
	s.productRepo.EXPECT().GetByID(ctx, productID).Return(product, nil)

	_, _, err = s.handler.Create(context.Background(), user, productID, file)

	invalidKrtErr := &errors.KRTValidationError{}
	s.True(s.ErrorAs(err, invalidKrtErr))
}
