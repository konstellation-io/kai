package entity

import "errors"

var (
	ErrMissingWorkflowKeyValueStore = errors.New("missing workflow's key-value store")
	ErrMissingProcessKeyValueStore  = errors.New("missing process' key-value store")
)

type KeyValueStores struct {
	GlobalKeyValueStore  string
	VersionKeyValueStore string
	Workflows            map[string]*WorkflowKeyValueStores
}

type WorkflowKeyValueStores struct {
	KeyValueStore string
	Processes     map[string]string
}

func (w *KeyValueStores) GetWorkflowKeyValueStore(workflow string) (string, error) {
	store, ok := w.Workflows[workflow]
	if !ok {
		return "", ErrMissingWorkflowKeyValueStore
	}

	return store.KeyValueStore, nil
}

func (w *WorkflowKeyValueStores) GetProcessKeyValueStore(process string) (string, error) {
	store, ok := w.Processes[process]
	if !ok {
		return "", ErrMissingProcessKeyValueStore
	}

	return store, nil
}
