package entity

import "fmt"

const entrypointProcessName = "entrypoint"

type VersionStreamsConfig struct {
	KeyValueStore string
	Workflows     map[string]*WorkflowStreamConfig
}

type WorkflowStreamConfig struct {
	Stream            string
	Processes         map[string]*ProcessStreamConfig
	EntrypointSubject string
	KeyValueStore     string
}

func (w *WorkflowStreamConfig) GetProcessConfig(processName string) (*ProcessStreamConfig, error) {
	processConfig, ok := w.Processes[processName]
	if !ok {
		//nolint:goerr113 // The error needs to be dynamic
		return nil, fmt.Errorf("error obtaining stream config for process %q", processName)
	}

	return processConfig, nil
}

func (w *WorkflowStreamConfig) GetEntrypointSubject() (string, error) {
	entrypointConfig, err := w.GetProcessConfig(entrypointProcessName)
	if err != nil {
		return "", err
	}

	return entrypointConfig.Subject, nil
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
