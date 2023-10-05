package manager

import (
	"errors"

	"github.com/konstellation-io/kai/engine/nats-manager/internal/entity"
)

func (m *NatsManager) UpdateKeyValueStoresConfiguration(configurations []entity.KeyValueConfiguration) error {
	m.logger.Info("Updating key-value stores configurations")

	var errs error

	for _, cfg := range configurations {
		err := m.client.UpdateConfiguration(cfg.KeyValueStore, cfg.Configuration)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}

	return errs
}
