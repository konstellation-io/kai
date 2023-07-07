package manager

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/konstellation-io/kai/engine/nats-manager/internal"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/entity"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/interfaces"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/logging"
)

type NatsManager struct {
	logger logging.Logger
	client interfaces.NatsClient
}

func NewNatsManager(logger logging.Logger, client interfaces.NatsClient) *NatsManager {
	return &NatsManager{
		logger: logger,
		client: client,
	}
}

func (m *NatsManager) CreateStreams(
	productID,
	versionName string,
	workflows []entity.Workflow,
) (entity.WorkflowsStreamsConfig, error) {
	if len(workflows) <= 0 {
		return nil, errors.New("no workflows defined")
	}

	workflowsStreamsConfig := entity.WorkflowsStreamsConfig{}

	for _, workflow := range workflows {
		stream := m.getStreamName(productID, versionName, workflow.Name)
		processesStreamConfig := m.getProcessesStreamConfig(stream, workflow.Processes)

		streamConfig := &entity.StreamConfig{
			Stream:    stream,
			Processes: processesStreamConfig,
		}

		err := m.client.CreateStream(streamConfig)
		if err != nil {
			return nil, fmt.Errorf("error creating stream %q: %w", stream, err)
		}

		workflowsStreamsConfig[workflow.Name] = streamConfig
	}

	return workflowsStreamsConfig, nil
}

func (m *NatsManager) CreateObjectStores(
	productID,
	versionName string,
	workflows []entity.Workflow,
) (entity.WorkflowsObjectStoresConfig, error) {
	if len(workflows) <= 0 {
		return nil, fmt.Errorf("no workflows defined")
	}

	if err := m.validateWorkflows(workflows); err != nil {
		return nil, fmt.Errorf("error validating worklfows: %w", err)
	}

	workflowsObjectStoresConfig := entity.WorkflowsObjectStoresConfig{}

	for _, workflow := range workflows {
		processesObjectStoresConfig := entity.ProcessesObjectStoresConfig{}

		for _, process := range workflow.Processes {
			if process.ObjectStore != nil {
				objectStore, err := m.getObjectStoreName(productID, versionName, workflow.Name, process.ObjectStore)
				if err != nil {
					return nil, err
				}

				err = m.client.CreateObjectStore(objectStore)
				if err != nil {
					return nil, fmt.Errorf("error creating object store %q: %w", objectStore, err)
				}

				processesObjectStoresConfig[process.Name] = objectStore
			}
		}
		workflowsObjectStoresConfig[workflow.Name] = &entity.WorkflowObjectStoresConfig{
			Processes: processesObjectStoresConfig,
		}
	}

	return workflowsObjectStoresConfig, nil
}

func (m *NatsManager) DeleteStreams(productID, versionName string) error {
	versionStreamsRegExp := m.getVersionStreamFilter(productID, versionName)

	allStreams, err := m.client.GetStreamNames(versionStreamsRegExp)
	if err != nil {
		return fmt.Errorf("error getting streams: %w", err)
	}
	for _, stream := range allStreams {
		err := m.client.DeleteStream(stream)
		if err != nil {
			return fmt.Errorf("error deleting stream %q: %w", stream, err)
		}
	}
	return nil
}

func (m *NatsManager) getObjectStoreName(productID, versionName, workflowName string, objectStore *entity.ObjectStore) (string, error) {
	switch objectStore.Scope {
	case entity.ObjStoreScopeProject:
		return m.joinWithUnderscores(productID, versionName, objectStore.Name), nil
	case entity.ObjStoreScopeWorkflow:
		return m.joinWithUnderscores(productID, versionName, workflowName, objectStore.Name), nil
	default:
		return "", internal.ErrInvalidObjectStoreScope
	}
}

func (m *NatsManager) DeleteObjectStores(productID, versionName string) error {
	versionObjStoreRegExp := m.getVersionStreamFilter(productID, versionName)

	allObjectStores, err := m.client.GetObjectStoreNames(versionObjStoreRegExp)
	if err != nil {
		return fmt.Errorf("error getting object store names: %w", err)
	}

	for _, objectStore := range allObjectStores {
		m.logger.Debugf("Deleting object store %q", objectStore)

		err := m.client.DeleteObjectStore(objectStore)
		if err != nil {
			return fmt.Errorf("error deleting object store %q: %w", objectStore, err)
		}
	}

	return nil
}

func (m *NatsManager) getStreamName(productID, versionName, workflowID string) string {
	return m.joinWithUnderscores(productID, versionName, workflowID)
}

func (m *NatsManager) CreateKeyValueStores(
	productID,
	versionName string,
	workflows []entity.Workflow,
) (*entity.VersionKeyValueStores, error) {
	if len(workflows) <= 0 {
		return nil, internal.ErrNoWorkflowsDefined
	}

	m.logger.Info("Creating key-value stores")

	// create key-value store for project
	productKeyValueStore, err := m.getKeyValueStoreName(productID, versionName, "", "", entity.KVScopeProject)
	if err != nil {
		return nil, err
	}

	err = m.client.CreateKeyValueStore(productKeyValueStore)
	if err != nil {
		return nil, fmt.Errorf("create product key-value store %q: %w", productKeyValueStore, err)
	}

	workflowsKeyValueStores := map[string]*entity.WorkflowKeyValueStores{}

	for _, workflow := range workflows {
		// create key-value store for workflow
		workflowKeyValueStore, err := m.getKeyValueStoreName(productID, versionName, workflow.Name, "", entity.KVScopeWorkflow)
		if err != nil {
			return nil, err
		}

		err = m.client.CreateKeyValueStore(workflowKeyValueStore)
		if err != nil {
			return nil, fmt.Errorf("create workflow key-value store %q: %w", workflowKeyValueStore, err)
		}

		processesKeyValueStores := map[string]string{}
		for _, process := range workflow.Processes {
			// create key-value store for process
			processKeyValueStore, err := m.getKeyValueStoreName(productID, versionName, workflow.Name, process.Name, entity.KVScopeProcess)
			if err != nil {
				return nil, err
			}

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

func (m *NatsManager) getKeyValueStoreName(
	product, version, workflow, process string,
	keyValueStore entity.KeyValueStoreScope,
) (string, error) {
	switch keyValueStore {
	case entity.KVScopeProject:
		return fmt.Sprintf("key-store_%s_%s", product, version), nil
	case entity.KVScopeWorkflow:
		return fmt.Sprintf("key-store_%s_%s_%s", product, version, workflow), nil
	case entity.KVScopeProcess:
		return fmt.Sprintf("key-store_%s_%s_%s_%s", product, version, workflow, process), nil
	default:
		return "", internal.ErrInvalidKeyValueStoreScope
	}
}

func (m *NatsManager) getProcessesStreamConfig(stream string, processes []entity.Process) entity.ProcessesStreamConfig {
	processesConfig := entity.ProcessesStreamConfig{}
	for _, process := range processes {
		processesConfig[process.Name] = entity.ProcessStreamConfig{
			Subject:       m.getSubjectName(stream, process.Name),
			Subscriptions: m.getSubjectsToSubscribe(stream, process.Subscriptions),
		}
	}
	return processesConfig
}

func (m *NatsManager) getSubjectName(stream, process string) string {
	return fmt.Sprintf("%s.%s", stream, process)
}

func (m *NatsManager) getSubjectsToSubscribe(stream string, subscriptions []string) []string {
	subjectsToSubscribe := make([]string, 0, len(subscriptions))
	for _, processToSubscribe := range subscriptions {
		subjectsToSubscribe = append(subjectsToSubscribe, m.getSubjectName(stream, processToSubscribe))
	}
	return subjectsToSubscribe
}

func (m *NatsManager) validateWorkflows(workflows []entity.Workflow) error {
	for _, workflow := range workflows {
		if err := workflow.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (m *NatsManager) getVersionStreamFilter(productID, versionName string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("^%s", m.joinWithUnderscores(productID, versionName, ".*")))
}

func (m *NatsManager) joinWithUnderscores(elements ...string) string {
	return strings.Join(elements, "_")
}
