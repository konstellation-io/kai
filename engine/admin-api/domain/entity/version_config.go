package entity

import (
	"errors"
)

var (
	ErrWorkflowStreamNotFound      = errors.New("workflow stream config not found")
	ErrWorkflowKVStoreNotFound     = errors.New("workflow key-value store config not found")
	ErrWorkflowObjectStoreNotFound = errors.New("workflow key-value store config not found")
)

type VersionConfig struct {
	KeyValueStore string

	StreamsConfig        *VersionStreamsConfig
	ObjectStoresConfig   *VersionObjectStoresConfig
	KeyValueStoresConfig *KeyValueStoresConfig
}

func NewVersionConfig(streamsConfig *VersionStreamsConfig, objectStoresConfig *VersionObjectStoresConfig,
	keyValueStoresConfig *KeyValueStoresConfig) *VersionConfig {
	return &VersionConfig{
		StreamsConfig:        streamsConfig,
		ObjectStoresConfig:   objectStoresConfig,
		KeyValueStoresConfig: keyValueStoresConfig,
	}
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
