package entity

import "errors"

type VersionStreamsConfig struct {
	KeyValueStore string
	Workflows     map[string]WorkflowStreamConfig
}

type WorkflowStreamConfig struct {
	Stream            string
	Processes         map[string]ProcessStreamConfig
	EntrypointSubject string
	KeyValueStore     string
}

func (w *WorkflowStreamConfig) GetProcessConfig(process string) (*ProcessStreamConfig, error) {
	processConfig, ok := w.Processes[process]
	if !ok {
		//nolint:goerr113 // The error needs to be dynamic
		return nil, errors.New("process configuration not found")
	}

	return &processConfig, nil
}

type ProcessStreamConfig struct {
	Subject       string
	ObjectStore   *string
	KeyValueStore string
	Subscriptions []string
}

type StreamInfo struct {
	Stream           string
	ProcesssSubjects map[string]string
}
