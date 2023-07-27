//go:build unit

package grpc_test

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/usecase"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/grpc"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/grpc/proto/versionpb"
	"github.com/konstellation-io/kai/engine/k8s-manager/mocks"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type VersionServiceTestSuite struct {
	suite.Suite
	cfg                *config.Config
	logger             logr.Logger
	versionServiceMock *mocks.VersionServiceMock
	versionGRPCService *grpc.VersionService
}

func TestVersionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(VersionServiceTestSuite))
}

func (s *VersionServiceTestSuite) SetupSuite() {
	cfg := &config.Config{}
	logger := testr.NewWithOptions(s.T(), testr.Options{Verbosity: -1})

	s.versionServiceMock = mocks.NewVersionServiceMock(s.T())

	s.versionGRPCService = grpc.NewVersionService(
		logger,
		s.versionServiceMock,
		s.versionServiceMock,
	)

	s.cfg = cfg
	s.logger = logger
}

func (s *VersionServiceTestSuite) TestStart() {
	ctx := context.Background()

	objectStore := "test-object-store"
	req := &versionpb.StartRequest{
		ProductId:     "test-product",
		VersionTag:    "test-version",
		KeyValueStore: "test-kv-store",
		Workflows: []*versionpb.Workflow{
			{
				Name: "test-workflow",
				Processes: []*versionpb.Process{
					{
						Name:          "test-process",
						Image:         "test-image",
						Gpu:           true,
						Subscriptions: []string{"test-subject"},
						Subject:       "test-subject",
						Replicas:      1,
						ObjectStore:   &objectStore,
						KeyValueStore: "test-kv-store",
						Type:          versionpb.ProcessType_ProcessTypeExit,
						Networking: &versionpb.Network{
							TargetPort: 8080,
							Protocol:   "TCP",
							SourcePort: 8080,
						},
						Config: map[string]string{
							"test-key": "test-value",
						},
						ResourceLimits: &versionpb.ProcessResourceLimits{
							Cpu: &versionpb.ResourceLimit{
								Request: "500m",
								Limit:   "1000m",
							},
							Memory: &versionpb.ResourceLimit{
								Request: "500Mi",
								Limit:   "1000Mi",
							},
						},
					},
				},
			},
		},
	}

	expectedVersion := domain.Version{
		Product:       req.ProductId,
		Tag:           req.VersionTag,
		KeyValueStore: req.KeyValueStore,
		Workflows: []*domain.Workflow{
			{
				Name: req.Workflows[0].Name,
				Processes: []*domain.Process{
					{
						Name:          req.Workflows[0].Processes[0].Name,
						Image:         req.Workflows[0].Processes[0].Image,
						EnableGpu:     req.Workflows[0].Processes[0].Gpu,
						Subscriptions: req.Workflows[0].Processes[0].Subscriptions,
						Subject:       req.Workflows[0].Processes[0].Subject,
						Replicas:      req.Workflows[0].Processes[0].Replicas,
						ObjectStore:   req.Workflows[0].Processes[0].ObjectStore,
						KeyValueStore: req.Workflows[0].Processes[0].KeyValueStore,
						Type:          domain.ProcessType(req.Workflows[0].Processes[0].Type),
						Networking: &domain.Networking{
							Protocol:   req.Workflows[0].Processes[0].Networking.Protocol,
							SourcePort: int(req.Workflows[0].Processes[0].Networking.SourcePort),
							TargetPort: int(req.Workflows[0].Processes[0].Networking.TargetPort),
						},
						Config: req.Workflows[0].Processes[0].Config,
						ResourceLimits: &domain.ProcessResourceLimits{

							CPU: &domain.ResourceLimit{
								Request: req.Workflows[0].Processes[0].ResourceLimits.Cpu.Request,
								Limit:   req.Workflows[0].Processes[0].ResourceLimits.Cpu.Limit,
							},
							Memory: &domain.ResourceLimit{
								Request: req.Workflows[0].Processes[0].ResourceLimits.Memory.Request,
								Limit:   req.Workflows[0].Processes[0].ResourceLimits.Memory.Limit,
							},
						},
					},
				},
			},
		},
	}

	s.versionServiceMock.EXPECT().StartVersion(ctx, expectedVersion).Return(nil)

	res, err := s.versionGRPCService.Start(ctx, req)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), res)
}

func (s *VersionServiceTestSuite) TestStop() {
	ctx := context.Background()

	req := &versionpb.StopRequest{
		Product:    "test-product",
		VersionTag: "test-version",
	}

	expectedParams := usecase.StopParams{
		Product: req.Product,
		Version: req.VersionTag,
	}

	s.versionServiceMock.EXPECT().StopVersion(ctx, expectedParams).Return(nil)

	res, err := s.versionGRPCService.Stop(ctx, req)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), res)
}
