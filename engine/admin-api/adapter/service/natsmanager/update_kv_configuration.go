package natsmanager

import (
	"context"
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/natspb"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

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
