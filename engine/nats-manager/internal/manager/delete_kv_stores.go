package manager

import (
	"fmt"

	"github.com/konstellation-io/kai/engine/nats-manager/internal"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/entity"
)

func (m *NatsManager) DeleteVersionKeyValueStores(
	productID,
	versionTag string,
	workflows []entity.Workflow,
) error {
	if len(workflows) == 0 {
		return internal.ErrNoWorkflowsDefined
	}

	m.logger.Info("Deleting key-value stores")

	productKeyValueStore := m.getVersionKeyValueStoreName(productID, versionTag)

	err := m.client.DeleteKeyValueStore(productKeyValueStore)
	if err != nil {
		return fmt.Errorf("delete product key-value store %q: %w", productKeyValueStore, err)
	}

	for _, workflow := range workflows {
		workflowKeyValueStore := m.getWorkflowKeyValueStoreName(productID, versionTag, workflow.Name)

		err = m.client.DeleteKeyValueStore(workflowKeyValueStore)
		if err != nil {
			return fmt.Errorf("delete workflow key-value store %q: %w", workflowKeyValueStore, err)
		}

		for _, process := range workflow.Processes {
			processKeyValueStore := m.getProcessKeyValueStoreName(productID, versionTag, workflow.Name, process.Name)

			err = m.client.DeleteKeyValueStore(processKeyValueStore)
			if err != nil {
				return fmt.Errorf("delete process key-value store %q: %w", processKeyValueStore, err)
			}
		}
	}

	return nil
}

func (m *NatsManager) DeleteGlobalKeyValueStore(productID string) error {
	m.logger.Info("Deleting global key-value store")

	keyValueStore := m.getProductKeyValueStoreName(productID)

	err := m.client.DeleteKeyValueStore(keyValueStore)
	if err != nil {
		return fmt.Errorf("deleting global key-value store %q: %w", keyValueStore, err)
	}

	return nil
}
