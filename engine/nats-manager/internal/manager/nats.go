package manager

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/konstellation-io/kai/engine/nats-manager/internal"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/entity"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/interfaces"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/logging"
)

const (
	_keyValueStorePrefix = "key-store_"
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
	versionTag string,
	workflows []entity.Workflow,
) (entity.WorkflowsStreamsConfig, error) {
	if len(workflows) == 0 {
		return nil, internal.ErrNoWorkflowsDefined
	}

	workflowsStreamsConfig := entity.WorkflowsStreamsConfig{}

	for _, workflow := range workflows {
		stream := m.getStreamName(productID, versionTag, workflow.Name)
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
	versionTag string,
	workflows []entity.Workflow,
) (entity.WorkflowsObjectStoresConfig, error) {
	if len(workflows) == 0 {
		return nil, internal.ErrNoWorkflowsDefined
	}

	if err := m.validateWorkflows(workflows); err != nil {
		return nil, fmt.Errorf("error validating worklfows: %w", err)
	}

	workflowsObjectStoresConfig := entity.WorkflowsObjectStoresConfig{}

	for _, workflow := range workflows {
		processesObjectStoresConfig := entity.ProcessesObjectStoresConfig{}

		for _, process := range workflow.Processes {
			if process.ObjectStore == nil {
				continue
			}

			objectStore, err := m.getObjectStoreName(productID, versionTag, workflow.Name, process.ObjectStore)
			if err != nil {
				return nil, err
			}

			err = m.client.CreateObjectStore(objectStore)
			if err != nil {
				return nil, fmt.Errorf("error creating object store %q: %w", objectStore, err)
			}

			processesObjectStoresConfig[process.Name] = objectStore
		}

		workflowsObjectStoresConfig[workflow.Name] = &entity.WorkflowObjectStoresConfig{
			Processes: processesObjectStoresConfig,
		}
	}

	return workflowsObjectStoresConfig, nil
}

func (m *NatsManager) DeleteStreams(productID, versionTag string) error {
	versionStreamsRegExp := m.getVersionStreamFilter(productID, versionTag)

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

func (m *NatsManager) getObjectStoreName(productID, versionTag, workflowName string, objectStore *entity.ObjectStore) (string, error) {
	versionTag = strings.ReplaceAll(versionTag, ".", "_")

	switch objectStore.Scope {
	case entity.ObjStoreScopeProject:
		return m.joinWithUnderscores(productID, versionTag, objectStore.Name), nil
	case entity.ObjStoreScopeWorkflow:
		return m.joinWithUnderscores(productID, versionTag, workflowName, objectStore.Name), nil
	case entity.ObjStoreScopeUndefined:
		return "", internal.ErrInvalidObjectStoreScope
	default:
		return "", internal.ErrInvalidObjectStoreScope
	}
}

func (m *NatsManager) DeleteObjectStores(productID, versionTag string) error {
	versionObjStoreRegExp := m.getVersionStreamFilter(productID, versionTag)

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

func (m *NatsManager) getStreamName(productID, versionTag, workflowID string) string {
	versionTag = strings.ReplaceAll(versionTag, ".", "_")
	return m.joinWithUnderscores(productID, versionTag, workflowID)
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

func (m *NatsManager) getVersionStreamFilter(productID, versionTag string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("^%s", m.joinWithUnderscores(productID, versionTag, ".*")))
}

func (m *NatsManager) joinWithUnderscores(elements ...string) string {
	return strings.Join(elements, "_")
}
