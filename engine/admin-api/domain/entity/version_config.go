package entity

import (
	"errors"
)

var (
	ErrNilVersionStreamConfig      = errors.New("stream config cannot be nil")
	ErrNilKeyValueStoreConfig      = errors.New("key-value store config cannot be nil")
	ErrWorkflowStreamNotFound      = errors.New("workflow stream config not found")
	ErrWorkflowKVStoreNotFound     = errors.New("workflow key-value store config not found")
	ErrWorkflowObjectStoreNotFound = errors.New("workflow key-value store config not found")
)

type VersionStreamingResources struct {
	Streams        *VersionStreams
	ObjectStores   *VersionObjectStores
	KeyValueStores *KeyValueStores
}

func NewVersionConfig(streamsConfig *VersionStreams, objectStoresConfig *VersionObjectStores,
	keyValueStoresConfig *KeyValueStores) (*VersionStreamingResources, error) {
	if streamsConfig == nil {
		return nil, ErrNilVersionStreamConfig
	}

	if keyValueStoresConfig == nil {
		return nil, ErrNilKeyValueStoreConfig
	}

	return &VersionStreamingResources{
		Streams:        streamsConfig,
		ObjectStores:   objectStoresConfig,
		KeyValueStores: keyValueStoresConfig,
	}, nil
}

func (v *VersionStreamingResources) GetWorkflowStreamConfig(workflow string) (*WorkflowStreamConfig, error) {
	w, ok := v.Streams.Workflows[workflow]
	if !ok {
		return nil, ErrWorkflowStreamNotFound
	}

	return &w, nil
}

func (v *VersionStreamingResources) GetWorkflowKeyValueStoresConfig(workflow string) (*WorkflowKeyValueStores, error) {
	w, ok := v.KeyValueStores.Workflows[workflow]
	if !ok {
		return nil, ErrWorkflowKVStoreNotFound
	}

	return w, nil
}

func (v *VersionStreamingResources) GetWorkflowObjectStoresConfig(workflow string) (*WorkflowObjectStoresConfig, error) {
	w, ok := v.ObjectStores.Workflows[workflow]
	if !ok {
		return nil, ErrWorkflowObjectStoreNotFound
	}

	return &w, nil
}
