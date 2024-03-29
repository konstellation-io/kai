package service

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/interfaces"
	"github.com/konstellation-io/kai/engine/nats-manager/proto/natspb"
)

// NatsService basic server.
type NatsService struct {
	logger  logr.Logger
	manager interfaces.NatsManager
	natspb.UnimplementedNatsManagerServiceServer
}

// NewNatsService instantiates the GRPC server implementation.
func NewNatsService(
	logger logr.Logger,
	manager interfaces.NatsManager,
) *NatsService {
	return &NatsService{
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

	streamConfig, err := n.manager.CreateStreams(req.ProductId, req.VersionTag, n.dtoToWorkflows(req.Workflows))
	if err != nil {
		n.logger.Error(err, "Error creating streams")
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

	objectStores, err := n.manager.CreateObjectStores(req.ProductId, req.VersionTag, n.dtoToWorkflows(req.Workflows))
	if err != nil {
		n.logger.Error(err, "Error creating object store")
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

	err := n.manager.DeleteStreams(req.ProductId, req.VersionTag)
	if err != nil {
		n.logger.Error(err, "Error deleting streams")
		return nil, err
	}

	return &natspb.DeleteResponse{
		Message: fmt.Sprintf("Streams and subjects for version %q on product %s deleted", req.VersionTag, req.ProductId),
	}, nil
}

// DeleteObjectStores delete object stores for given workflows.
func (n *NatsService) DeleteObjectStores(
	_ context.Context,
	req *natspb.DeleteObjectStoresRequest,
) (*natspb.DeleteResponse, error) {
	n.logger.Info("Delete object stores request received")

	err := n.manager.DeleteObjectStores(req.ProductId, req.VersionTag)
	if err != nil {
		n.logger.Error(err, "Error deleting object stores")
		return nil, err
	}

	return &natspb.DeleteResponse{
		Message: fmt.Sprintf("Object stores for version %q on product %s deleted", req.VersionTag, req.ProductId),
	}, nil
}

func (n *NatsService) CreateVersionKeyValueStores(
	_ context.Context,
	req *natspb.CreateVersionKeyValueStoresRequest,
) (*natspb.CreateVersionKeyValueStoresResponse, error) {
	n.logger.Info("CreateVersionKeyValueStores request received")

	keyValueStores, err := n.manager.CreateVersionKeyValueStores(req.ProductId, req.VersionTag, n.dtoToWorkflows(req.Workflows))
	if err != nil {
		n.logger.Error(err, "Error creating version's key-value store")
		return nil, err
	}

	return n.mapKeyValueStoresToDTO(keyValueStores), nil
}

func (n *NatsService) CreateGlobalKeyValueStore(
	_ context.Context,
	req *natspb.CreateGlobalKeyValueStoreRequest,
) (*natspb.CreateGlobalKeyValueStoreResponse, error) {
	n.logger.Info("CreateGlobalKeyValueStore request received")

	keyValueStore, err := n.manager.CreateGlobalKeyValueStore(req.ProductId)
	if err != nil {
		n.logger.Error(err, "Error creating global key-value store")
		return nil, err
	}

	return &natspb.CreateGlobalKeyValueStoreResponse{GlobalKeyValueStore: keyValueStore}, nil
}

func (n *NatsService) DeleteVersionKeyValueStores(
	_ context.Context,
	req *natspb.DeleteVersionKeyValueStoresRequest,
) (*natspb.DeleteResponse, error) {
	n.logger.Info("DeleteVersionKeyValueStores request received")

	err := n.manager.DeleteVersionKeyValueStores(req.ProductId, req.VersionTag, n.dtoToWorkflows(req.Workflows))
	if err != nil {
		n.logger.Error(err, "Error deleting version's key-value stores")
		return nil, err
	}

	return &natspb.DeleteResponse{
		Message: fmt.Sprintf("Key-value stores for version %q on product %s deleted", req.VersionTag, req.ProductId),
	}, nil
}

func (n *NatsService) UpdateKeyValueConfiguration(
	_ context.Context,
	req *natspb.UpdateKeyValueConfigurationRequest,
) (*natspb.UpdateKeyValueConfigurationResponse, error) {
	n.logger.Info("CreateGlobalKeyValueStore request received")

	err := n.manager.UpdateKeyValueStoresConfiguration(n.mapDTOToKeyValueStoreConfigurations(req.KeyValueStoresConfig))
	if err != nil {
		n.logger.Error(err, "Error creating object store")
		return nil, err
	}

	return &natspb.UpdateKeyValueConfigurationResponse{
		Message: "Configurations successfully updated!",
	}, nil
}

func (n *NatsService) DeleteGlobalKeyValueStore(
	_ context.Context,
	req *natspb.DeleteGlobalKeyValueStoreRequest,
) (*natspb.DeleteResponse, error) {
	n.logger.Info("DeleteVersionKeyValueStores request received")

	err := n.manager.DeleteGlobalKeyValueStore(req.ProductId)
	if err != nil {
		n.logger.Error(err, "Error deleting global key-value store", "product", req.ProductId)
		return nil, err
	}

	return &natspb.DeleteResponse{
		Message: fmt.Sprintf("Global key-value store for product %q deleted", req.ProductId),
	}, nil
}
