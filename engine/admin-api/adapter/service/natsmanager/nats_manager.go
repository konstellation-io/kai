package natsmanager

import (
	"context"
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/natspb"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/logging"
)

//go:generate mockgen -source=../proto/natspb/nats_grpc.pb.go -destination=../../../mocks/${GOFILE} -package=mocks

type Client struct {
	cfg    *config.Config
	client natspb.NatsManagerServiceClient
	logger logging.Logger
}

func NewClient(cfg *config.Config, logger logging.Logger, client natspb.NatsManagerServiceClient) (*Client, error) {
	return &Client{
		cfg,
		client,
		logger,
	}, nil
}

// CreateStreams calls nats-manager to create NATS streams for given version.
func (n *Client) CreateStreams(
	ctx context.Context,
	productID string,
	version *entity.Version,
) (*entity.VersionStreamsConfig, error) {
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
) (*entity.VersionObjectStoresConfig, error) {
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

func (n *Client) CreateGlobalKeyValueStore(ctx context.Context, product string) (string, error) {
	req := natspb.CreateGlobalKeyValueStoreRequest{
		ProductId: product,
	}

	res, err := n.client.CreateGlobalKeyValueStore(ctx, &req)
	if err != nil {
		return "", fmt.Errorf("creating global key-value store: %w", err)
	}

	return res.GlobalKeyValueStore, err
}

func (n *Client) UpdateKeyValueConfiguration(ctx context.Context, configurations []entity.KeyValueConfiguration) error {
	req := natspb.UpdateKeyValueConfigurationRequest{
		KeyValueStoresConfig: n.mapKeyValueConfigToDTO(configurations),
	}

	_, err := n.client.UpdateKeyValueConfiguration(ctx, &req)
	if err != nil {
		return fmt.Errorf("creating global key-value store: %w", err)
	}

	return err
}

func (n *Client) mapKeyValueConfigToDTO(configurations []entity.KeyValueConfiguration) []*natspb.KeyValueConfiguration {
	dtoConfigurations := make([]*natspb.KeyValueConfiguration, 0, len(configurations))

	for _, configuration := range configurations {
		dtoConfigurations = append(dtoConfigurations, &natspb.KeyValueConfiguration{
			KeyValueStore: configuration.Store,
			Configuration: n.mapConfigurationToDTO(configuration.Configuration),
		})
	}

	return dtoConfigurations
}

func (n *Client) mapConfigurationToDTO(configuration []entity.ConfigurationVariable) map[string]string {
	configDTO := make(map[string]string, len(configuration))

	for _, cfg := range configuration {
		configDTO[cfg.Key] = cfg.Value
	}

	return configDTO
}
