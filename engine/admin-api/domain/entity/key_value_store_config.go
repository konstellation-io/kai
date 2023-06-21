package entity

import "fmt"

type KeyValueStoresConfig struct {
	ProjectKeyValueStore    string
	WorkflowsKeyValueStores WorkflowsKeyValueStoresConfig
}

type WorkflowsKeyValueStoresConfig map[string]*WorkflowKeyValueStores

type WorkflowKeyValueStores struct {
	WorkflowKeyValueStore  string
	ProcesssKeyValueStores map[string]string
}

func (w *WorkflowKeyValueStores) GetProcessKeyValueStore(process string) (string, error) {
	store, ok := w.ProcesssKeyValueStores[process]
	if !ok {
		//nolint:goerr113 // error needs to be dynamic
		return "", fmt.Errorf("missing key value store for process %q", process)
	}

	return store, nil
}
