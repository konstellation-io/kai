package entity

type WorkflowsStreamsConfig map[string]*StreamConfig

type StreamConfig struct {
	Stream    string
	Processes ProcessesStreamConfig
}

type ProcessesStreamConfig map[string]ProcessStreamConfig

type ProcessStreamConfig struct {
	Subject       string
	Subscriptions []string
}
