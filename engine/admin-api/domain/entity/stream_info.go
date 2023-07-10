package entity

import "errors"

var ErrProcessStreamNotFound = errors.New("process stream configuration not found")

type VersionStreamsConfig struct {
	Workflows map[string]WorkflowStreamConfig
}

type WorkflowStreamConfig struct {
	Stream    string
	Processes map[string]ProcessStreamConfig
}

func (w *WorkflowStreamConfig) GetProcessConfig(process string) (*ProcessStreamConfig, error) {
	processConfig, ok := w.Processes[process]
	if !ok {
		return nil, ErrProcessStreamNotFound
	}

	return &processConfig, nil
}

type ProcessStreamConfig struct {
	Subject       string
	Subscriptions []string
}

type StreamInfo struct {
	Stream           string
	ProcesssSubjects map[string]string
}
