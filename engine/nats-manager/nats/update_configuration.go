package nats

import (
	"errors"
	"fmt"
)

func (n *NatsClient) UpdateConfiguration(keyValueStore string, keyValueConfig map[string]string) error {
	n.logger.V(2).Info("Updating config to key-value store", "key-value-store", keyValueStore)

	kvStoreBucket, err := n.js.KeyValue(keyValueStore)
	if err != nil {
		return fmt.Errorf("initializing key-value store bucket: %w", err)
	}

	var errs error

	for key, value := range keyValueConfig {
		_, err := kvStoreBucket.PutString(key, value)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("updating value for key %q: %w", key, err))
		}
	}

	return errs
}
