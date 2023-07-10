package entity

import "fmt"

type KeyValueStoresConfig struct {
	KeyValueStore string
	Workflows     WorkflowsKeyValueStoresConfig
}

type WorkflowsKeyValueStoresConfig map[string]*WorkflowKeyValueStores

type WorkflowKeyValueStores struct {
	KeyValueStore string
	Processes     map[string]string
}

func (w *WorkflowKeyValueStores) GetProcessKeyValueStore(process string) (string, error) {
	store, ok := w.Processes[process]
	if !ok {
		//nolint:goerr113 // error needs to be dynamic
		return "", fmt.Errorf("missing key value store for process %q", process)
	}

	return store, nil
}
