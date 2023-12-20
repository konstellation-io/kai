package entity

type VersionObjectStores struct {
	Workflows map[string]WorkflowObjectStoresConfig
}

type WorkflowObjectStoresConfig struct {
	Processes ProcessObjectStoresConfig
}

type ProcessObjectStoresConfig map[string]string

func (n ProcessObjectStoresConfig) GetProcessObjectStoreConfig(process string) *string {
	processObjStore, ok := n[process]
	if !ok {
		return nil
	}

	return &processObjStore
}
