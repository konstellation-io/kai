package manager

import (
	"errors"
	"fmt"

	"github.com/konstellation-io/kai/engine/nats-manager/internal/entity"
)

func (m *NatsManager) UpdateKeyValueStoresConfiguration(configurations []entity.KeyValueConfiguration) error {
	if len(configurations) == 0 {
		m.logger.Info("No configurations to update")
		return nil
	}

	m.logger.Info("Updating key-value stores configurations")

	var errs error

	fmt.Println(configurations)

	for _, cfg := range configurations {
		err := m.client.UpdateConfiguration(cfg.KeyValueStore, cfg.Configuration)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}

	return errs
}
