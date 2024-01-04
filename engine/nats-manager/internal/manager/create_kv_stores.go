package manager

import (
	"fmt"
	"strings"

	"github.com/konstellation-io/kai/engine/nats-manager/internal"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/entity"
)

func (m *NatsManager) CreateVersionKeyValueStores(
	productID,
	versionTag string,
	workflows []entity.Workflow,
) (*entity.VersionKeyValueStores, error) {
	if len(workflows) == 0 {
		return nil, internal.ErrNoWorkflowsDefined
	}

	m.logger.Info("Creating key-value stores")

	productKeyValueStore := m.getVersionKeyValueStoreName(productID, versionTag)

	err := m.client.CreateKeyValueStore(productKeyValueStore)
	if err != nil {
		return nil, fmt.Errorf("create product key-value store %q: %w", productKeyValueStore, err)
	}

	workflowsKeyValueStores := make(map[string]*entity.WorkflowKeyValueStores, len(workflows))

	for _, workflow := range workflows {
		workflowKeyValueStore := m.getWorkflowKeyValueStoreName(productID, versionTag, workflow.Name)

		err = m.client.CreateKeyValueStore(workflowKeyValueStore)
		if err != nil {
			return nil, fmt.Errorf("create workflow key-value store %q: %w", workflowKeyValueStore, err)
		}

		processesKeyValueStores := make(map[string]string, len(workflow.Processes))

		for _, process := range workflow.Processes {
			processKeyValueStore := m.getProcessKeyValueStoreName(productID, versionTag, workflow.Name, process.Name)

			err = m.client.CreateKeyValueStore(processKeyValueStore)
			if err != nil {
				return nil, fmt.Errorf("create process key-value store %q: %w", processKeyValueStore, err)
			}

			processesKeyValueStores[process.Name] = processKeyValueStore
		}

		workflowsKeyValueStores[workflow.Name] = &entity.WorkflowKeyValueStores{
			WorkflowStore: workflowKeyValueStore,
			Processes:     processesKeyValueStores,
		}
	}

	return &entity.VersionKeyValueStores{
		ProjectStore:    productKeyValueStore,
		WorkflowsStores: workflowsKeyValueStores,
	}, nil
}

func (m *NatsManager) CreateGlobalKeyValueStore(product string) (string, error) {
	m.logger.Info("Creating global key-value store for product", "product", product)

	keyValueStoreName := m.getProductKeyValueStoreName(product)

	err := m.client.CreateKeyValueStore(keyValueStoreName)
	if err != nil {
		return "", fmt.Errorf("creating global key-value store: %w", err)
	}

	return keyValueStoreName, nil
}

func (m *NatsManager) getProductKeyValueStoreName(product string) string {
	return m.replaceInvalidCharacters(fmt.Sprintf("%s%s", _keyValueStorePrefix, product))
}

func (m *NatsManager) getVersionKeyValueStoreName(product, versionTag string) string {
	return m.replaceInvalidCharacters(fmt.Sprintf("%s%s_%s", _keyValueStorePrefix, product, versionTag))
}

func (m *NatsManager) getWorkflowKeyValueStoreName(product, versionTag, workflow string) string {
	return m.replaceInvalidCharacters(fmt.Sprintf("%s%s_%s_%s", _keyValueStorePrefix, product, versionTag, workflow))
}

func (m *NatsManager) getProcessKeyValueStoreName(product, versionTag, workflow, process string) string {
	return m.replaceInvalidCharacters(fmt.Sprintf("%s%s_%s_%s_%s", _keyValueStorePrefix, product, versionTag, workflow, process))
}

func (m *NatsManager) replaceInvalidCharacters(name string) string {
	return strings.ReplaceAll(name, ".", "_")
}
