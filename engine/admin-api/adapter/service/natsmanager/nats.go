package natsmanager

import (
	"context"
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/natspb"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NatsManagerClient struct {
	cfg    *config.Config
	client natspb.NatsManagerServiceClient
	logger logging.Logger
}

func NewNatsManagerClient(cfg *config.Config, logger logging.Logger) (*NatsManagerClient, error) {
	cc, err := grpc.Dial(cfg.Services.NatsManager, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	client := natspb.NewNatsManagerServiceClient(cc)

	return &NatsManagerClient{
		cfg,
		client,
		logger,
	}, nil
}

// CreateStreams calls nats-manager to create NATS streams for given version.
//
//nolint:dupl // this is not being duplicated
func (n *NatsManagerClient) CreateStreams(
	ctx context.Context,
	productID string,
	version *entity.Version,
) (*entity.VersionStreamsConfig, error) {
	workflows, err := n.mapWorkflowsToDTO(version.Workflows)
	if err != nil {
		return nil, err
	}

	req := natspb.CreateStreamsRequest{
		ProductId:   productID,
		VersionName: version.Name,
		Workflows:   workflows,
	}

	res, err := n.client.CreateStreams(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error creating streams: %w", err)
	}

	return n.mapDTOToVersionStreamConfig(res.Workflows), err
}

// CreateObjectStores calls nats-manager to create NATS Object Stores for given version.
//
//nolint:dupl // this is not being duplicated
func (n *NatsManagerClient) CreateObjectStores(
	ctx context.Context,
	productID string,
	version *entity.Version,
) (*entity.VersionObjectStoresConfig, error) {
	workflows, err := n.mapWorkflowsToDTO(version.Workflows)
	if err != nil {
		return nil, err
	}

	req := natspb.CreateObjectStoresRequest{
		ProductId:   productID,
		VersionName: version.Name,
		Workflows:   workflows,
	}

	res, err := n.client.CreateObjectStores(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error creating object stores: %w", err)
	}

	return n.mapDTOToVersionObjectStoreConfig(res.Workflows), err
}

// CreateKeyValueStores calls nats-manager to create NATS Key Value Stores for given version.
func (n *NatsManagerClient) CreateKeyValueStores(
	ctx context.Context,
	productID string,
	version *entity.Version,
) (*entity.KeyValueStoresConfig, error) {
	workflows, err := n.mapWorkflowsToDTO(version.Workflows)
	if err != nil {
		return nil, err
	}

	req := natspb.CreateKeyValueStoresRequest{
		ProductId:   productID,
		VersionName: version.Name,
		Workflows:   workflows,
	}

	res, err := n.client.CreateKeyValueStores(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error creating key value stores: %w", err)
	}

	return n.mapDTOToVersionKeyValueStoreConfig(res.KeyValueStore, res.Workflows), err
}

// DeleteStreams calls nats-manager to delete NATS streams for given version.
func (n *NatsManagerClient) DeleteStreams(ctx context.Context, productID, versionName string) error {
	req := natspb.DeleteStreamsRequest{
		ProductId:   productID,
		VersionName: versionName,
	}

	_, err := n.client.DeleteStreams(ctx, &req)
	if err != nil {
		return fmt.Errorf("error deleting version %q NATS streams: %w", versionName, err)
	}

	return nil
}

// DeleteObjectStores calls nats-manager to delete NATS Object Stores for given version.
func (n *NatsManagerClient) DeleteObjectStores(ctx context.Context, productID, versionName string) error {
	req := natspb.DeleteObjectStoresRequest{
		ProductId:   productID,
		VersionName: versionName,
	}

	_, err := n.client.DeleteObjectStores(ctx, &req)
	if err != nil {
		return fmt.Errorf("error deleting version %q NATS object stores: %w", versionName, err)
	}

	return nil
}
