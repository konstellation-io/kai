package natsmanager

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/natspb"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

//go:generate mockgen -source=../proto/natspb/nats_grpc.pb.go -destination=../../../mocks/${GOFILE} -package=mocks

type Client struct {
	client natspb.NatsManagerServiceClient
	logger logr.Logger
}

func NewClient(logger logr.Logger, client natspb.NatsManagerServiceClient) (*Client, error) {
	return &Client{
		client,
		logger,
	}, nil
}

// CreateStreams calls nats-manager to create NATS streams for given version.
func (n *Client) CreateStreams(
	ctx context.Context,
	productID string,
	version *entity.Version,
) (*entity.VersionStreams, error) {
	req := natspb.CreateStreamsRequest{
		ProductId:  productID,
		VersionTag: version.Tag,
		Workflows:  n.mapWorkflowsToDTO(version.Workflows),
	}

	res, err := n.client.CreateStreams(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error creating streams: %w", err)
	}

	return n.mapDTOToVersionStreamConfig(res.Workflows), err
}

// CreateObjectStores calls nats-manager to create NATS Object Stores for given version.
func (n *Client) CreateObjectStores(
	ctx context.Context,
	productID string,
	version *entity.Version,
) (*entity.VersionObjectStores, error) {
	req := natspb.CreateObjectStoresRequest{
		ProductId:  productID,
		VersionTag: version.Tag,
		Workflows:  n.mapWorkflowsToDTO(version.Workflows),
	}

	res, err := n.client.CreateObjectStores(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error creating object stores: %w", err)
	}

	return n.mapDTOToVersionObjectStoreConfig(res.Workflows), err
}

// CreateVersionKeyValueStores calls nats-manager to create NATS Key Value Stores for given version.
func (n *Client) CreateVersionKeyValueStores(
	ctx context.Context,
	productID string,
	version *entity.Version,
) (*entity.KeyValueStores, error) {
	req := natspb.CreateVersionKeyValueStoresRequest{
		ProductId:  productID,
		VersionTag: version.Tag,
		Workflows:  n.mapWorkflowsToDTO(version.Workflows),
	}

	res, err := n.client.CreateVersionKeyValueStores(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error creating key value stores: %w", err)
	}

	return n.mapDTOToVersionKeyValueStoreConfig(res.KeyValueStore, res.Workflows), err
}

// DeleteStreams calls nats-manager to delete NATS streams for given version.
func (n *Client) DeleteStreams(ctx context.Context, productID, versionTag string) error {
	req := natspb.DeleteStreamsRequest{
		ProductId:  productID,
		VersionTag: versionTag,
	}

	_, err := n.client.DeleteStreams(ctx, &req)
	if err != nil {
		return fmt.Errorf("error deleting version %q NATS streams: %w", versionTag, err)
	}

	return nil
}

// DeleteObjectStores calls nats-manager to delete NATS Object Stores for given version.
func (n *Client) DeleteObjectStores(ctx context.Context, productID, versionTag string) error {
	req := natspb.DeleteObjectStoresRequest{
		ProductId:  productID,
		VersionTag: versionTag,
	}

	_, err := n.client.DeleteObjectStores(ctx, &req)
	if err != nil {
		return fmt.Errorf("error deleting version %q NATS object stores: %w", versionTag, err)
	}

	return nil
}

// DeleteVersionKeyValueStores calls nats-manager to delete NATS Key Value Stores for given version.
func (n *Client) DeleteVersionKeyValueStores(
	ctx context.Context,
	productID string,
	version *entity.Version,
) error {
	req := natspb.DeleteVersionKeyValueStoresRequest{
		ProductId:  productID,
		VersionTag: version.Tag,
		Workflows:  n.mapWorkflowsToDTO(version.Workflows),
	}

	_, err := n.client.DeleteVersionKeyValueStores(ctx, &req)
	if err != nil {
		return fmt.Errorf(
			"error deleting product %q version %q NATS key value stores: %w", productID, version.Tag, err,
		)
	}

	return nil
}
