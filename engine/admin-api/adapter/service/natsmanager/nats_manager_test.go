//go:build unit

package natsmanager_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/natsmanager"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/natspb"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
	"github.com/stretchr/testify/suite"
)

type createStreamsRequestMatcher struct {
	expectedCreateStreamsRequest *natspb.CreateStreamsRequest
}

func newCreateStreamsRequestMatcher(expectedStreamConfig *natspb.CreateStreamsRequest) *createStreamsRequestMatcher {
	return &createStreamsRequestMatcher{
		expectedCreateStreamsRequest: expectedStreamConfig,
	}
}
func (m createStreamsRequestMatcher) String() string {
	return fmt.Sprintf("is equal to %v", m.expectedCreateStreamsRequest)
}

func (m createStreamsRequestMatcher) Matches(actual interface{}) bool {
	actualCfg, ok := actual.(*natspb.CreateStreamsRequest)
	if !ok {
		return false
	}

	return reflect.DeepEqual(actualCfg, m.expectedCreateStreamsRequest)
}

type NatsManagerTestSuite struct {
	suite.Suite
	cfg               *config.Config
	logger            logging.Logger
	mockService       *mocks.MockNatsManagerServiceClient
	natsManagerClient *natsmanager.NatsManagerClient
}

func TestNatsManagerTestSuite(t *testing.T) {
	suite.Run(t, new(NatsManagerTestSuite))
}

func (s *NatsManagerTestSuite) SetupSuite() {
	mockController := gomock.NewController(s.T())
	cfg := &config.Config{}
	logger := mocks.NewMockLogger(mockController)
	mocks.AddLoggerExpects(logger)
	service := mocks.NewMockNatsManagerServiceClient(mockController)

	k8sVersionClient, err := natsmanager.NewNatsManagerClient(cfg, logger, service)
	s.Require().NoError(err)

	s.cfg = cfg
	s.logger = logger
	s.mockService = service
	s.natsManagerClient = k8sVersionClient
}

func (s *NatsManagerTestSuite) TestCreateStreams() {
	ctx := context.Background()

	var (
		productID = "test-product"

		process = testhelpers.NewProcessBuilder().
			WithObjectStore(&entity.ProcessObjectStore{
				Name:  "test-object-store",
				Scope: entity.ObjectStoreScopeWorkflow,
			}).
			WithNetworking(&entity.ProcessNetworking{
				TargetPort:      8080,
				DestinationPort: 8080,
				Protocol:        "TCP",
			}).
			Build()

		workflow = testhelpers.NewWorkflowBuilder().
				WithProcesses([]entity.Process{process}).
				Build()

		version = testhelpers.NewVersionBuilder().
			WithWorkflows([]entity.Workflow{workflow}).
			Build()
	)

	req := &natspb.CreateStreamsRequest{
		ProductId:   productID,
		VersionName: version.Name,
		Workflows: []*natspb.Workflow{
			{
				Name: workflow.Name,
				Processes: []*natspb.Process{
					{
						Name: process.Name,
						ObjectStore: &natspb.ObjectStore{
							Name:  process.ObjectStore.Name,
							Scope: natspb.ObjectStoreScope_SCOPE_WORKFLOW,
						},
						Subscriptions: process.Subscriptions,
					},
				},
			},
		},
	}

	customMatcher := newCreateStreamsRequestMatcher(req)
	s.mockService.EXPECT().CreateStreams(ctx, customMatcher).Return(&natspb.CreateStreamsResponse{}, nil)

	_, err := s.natsManagerClient.CreateStreams(ctx, productID, &version)
	s.Require().NoError(err)
}
