package entity

import (
	"errors"
	"fmt"
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
		//nolint:goerr113 // The error needs to be dynamic
		return nil, fmt.Errorf("workflow %q stream config not found", workflow)
	}

	return &w, nil
}

func (v *VersionConfig) GetWorkflowKeyValueStoresConfig(workflow string) (*WorkflowKeyValueStores, error) {
	w, ok := v.KeyValueStoresConfig.WorkflowsKeyValueStores[workflow]
	if !ok {
		//nolint:goerr113 // errors need to be dynamic
		return nil, fmt.Errorf("workflow %q stream config not found", workflow)
	}

	return w, nil
}

func (v *VersionConfig) GetWorkflowObjectStoresConfig(workflow string) (*WorkflowObjectStoresConfig, error) {
	w, ok := v.ObjectStoresConfig.Workflows[workflow]
	if !ok {
		//nolint:goerr113 // errors need to be dynamic
		return nil, errors.New("object store config not found")
	}

	return &w, nil
}
