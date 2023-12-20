package entity

type WorkflowsObjectStoresConfig map[string]*WorkflowObjectStoresConfig

type WorkflowObjectStoresConfig struct {
	Processes ProcessesObjectStoresConfig
}

type ProcessesObjectStoresConfig map[string]string
