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

type VersionConfig struct {
	StreamsConfig        *VersionStreamsConfig
	ObjectStoresConfig   *VersionObjectStoresConfig
	KeyValueStoresConfig *KeyValueStores
}

func NewVersionConfig(streamsConfig *VersionStreamsConfig, objectStoresConfig *VersionObjectStoresConfig,
	keyValueStoresConfig *KeyValueStores) (*VersionConfig, error) {
	if streamsConfig == nil {
		return nil, ErrNilVersionStreamConfig
	}

	if keyValueStoresConfig == nil {
		return nil, ErrNilKeyValueStoreConfig
	}

	return &VersionConfig{
		StreamsConfig:        streamsConfig,
		ObjectStoresConfig:   objectStoresConfig,
		KeyValueStoresConfig: keyValueStoresConfig,
	}, nil
}

func (v *VersionConfig) GetWorkflowStreamConfig(workflow string) (*WorkflowStreamConfig, error) {
	w, ok := v.StreamsConfig.Workflows[workflow]
	if !ok {
		return nil, ErrWorkflowStreamNotFound
	}

	return &w, nil
}

func (v *VersionConfig) GetWorkflowKeyValueStoresConfig(workflow string) (*WorkflowKeyValueStores, error) {
	w, ok := v.KeyValueStoresConfig.Workflows[workflow]
	if !ok {
		return nil, ErrWorkflowKVStoreNotFound
	}

	return w, nil
}

func (v *VersionConfig) GetWorkflowObjectStoresConfig(workflow string) (*WorkflowObjectStoresConfig, error) {
	w, ok := v.ObjectStoresConfig.Workflows[workflow]
	if !ok {
		return nil, ErrWorkflowObjectStoreNotFound
	}

	return &w, nil
}
