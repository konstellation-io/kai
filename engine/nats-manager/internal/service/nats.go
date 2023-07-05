package service

import (
	"context"
	"fmt"

	"github.com/konstellation-io/kai/engine/nats-manager/internal/config"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/interfaces"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/logging"
	"github.com/konstellation-io/kai/engine/nats-manager/proto/natspb"
)

// NatsService basic server.
type NatsService struct {
	config  *config.Config
	logger  logging.Logger
	manager interfaces.NatsManager
	natspb.UnimplementedNatsManagerServiceServer
}

// NewNatsService instantiates the GRPC server implementation.
func NewNatsService(
	cfg *config.Config,
	logger logging.Logger,
	manager interfaces.NatsManager,
) *NatsService {
	return &NatsService{
		cfg,
		logger,
		manager,
		natspb.UnimplementedNatsManagerServiceServer{},
	}
}

// CreateStreams create streams for given workflows.
func (n *NatsService) CreateStreams(
	_ context.Context,
	req *natspb.CreateStreamsRequest,
) (*natspb.CreateStreamsResponse, error) {
	n.logger.Info("CreateStreams request received")

	streamConfig, err := n.manager.CreateStreams(req.ProductId, req.VersionName, n.dtoToWorkflows(req.Workflows))
	if err != nil {
		n.logger.Errorf("Error creating streams: %s", err)
		return nil, err
	}

	return &natspb.CreateStreamsResponse{
		Workflows: n.workflowsStreamConfigToDto(streamConfig),
	}, nil
}

// CreateObjectStores creates object stores for given workflows.
func (n *NatsService) CreateObjectStores(
	_ context.Context,
	req *natspb.CreateObjectStoresRequest,
) (*natspb.CreateObjectStoresResponse, error) {
	n.logger.Info("CreateObjectStores request received")

	objectStores, err := n.manager.CreateObjectStores(req.ProductId, req.VersionName, n.dtoToWorkflows(req.Workflows))
	if err != nil {
		n.logger.Errorf("Error creating object store: %s", err)
		return nil, err
	}
	return &natspb.CreateObjectStoresResponse{
		Workflows: n.mapWorkflowsObjStoreToDTO(objectStores),
	}, nil
}

// DeleteStreams delete streams for given workflows.
func (n *NatsService) DeleteStreams(
	_ context.Context,
	req *natspb.DeleteStreamsRequest,
) (*natspb.DeleteResponse, error) {
	n.logger.Info("Delete streams request received")

	err := n.manager.DeleteStreams(req.ProductId, req.VersionName)
	if err != nil {
		n.logger.Errorf("Error deleting streams: %s", err)
		return nil, err
	}

	return &natspb.DeleteResponse{
		Message: fmt.Sprintf("Streams and subjects for version %q on product %s deleted", req.VersionName, req.ProductId),
	}, nil
}

// DeleteObjectStores delete object stores for given workflows.
func (n *NatsService) DeleteObjectStores(
	_ context.Context,
	req *natspb.DeleteObjectStoresRequest,
) (*natspb.DeleteResponse, error) {
	n.logger.Info("Delete object stores request received")

	err := n.manager.DeleteObjectStores(req.ProductId, req.VersionName)
	if err != nil {
		n.logger.Errorf("Error deleting object stores: %s", err)
		return nil, err
	}

	return &natspb.DeleteResponse{
		Message: fmt.Sprintf("Object stores for version %q on product %s deleted", req.VersionName, req.ProductId),
	}, nil
}

func (n *NatsService) CreateKeyValueStores(
	_ context.Context,
	req *natspb.CreateKeyValueStoresRequest,
) (*natspb.CreateKeyValueStoreResponse, error) {
	n.logger.Info("CreateKeyValueStores request received")

	keyValueStores, err := n.manager.CreateKeyValueStores(req.ProductId, req.VersionName, n.dtoToWorkflows(req.Workflows))
	if err != nil {
		n.logger.Errorf("Error creating object store: %s", err)
		return nil, err
	}
	return n.mapKeyValueStoresToDTO(keyValueStores), nil
}
