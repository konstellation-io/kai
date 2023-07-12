package configuration

type ProcessConfig struct {
	Metadata          Metadata          `yaml:"metadata"`
	Nats              NatsConfig        `yaml:"nats"`
	CentralizedConfig CentralizedConfig `yaml:"centralized_configuration"`
}

type Metadata struct {
	ProductID    string `yaml:"product_id"`
	VersionTag   string `yaml:"version_tag"`
	WorkflowName string `yaml:"workflow_name"`
	ProcessName  string `yaml:"process_name"`
	BasePath     string `yaml:"base_path"`
	ProcessType  string `yaml:"process_type"`
}

type CentralizedConfig struct {
	Product  ConfigDefinition `yaml:"product"`
	Workflow ConfigDefinition `yaml:"workflow"`
	Process  ConfigDefinition `yaml:"process"`
}

type ConfigDefinition struct {
	Bucket string            `yaml:"bucket"`
	Config map[string]string `yaml:"config,omitempty"`
}

type NatsConfig struct {
	URL           string   `yaml:"url"`
	Stream        string   `yaml:"stream"`
	Subject       string   `yaml:"output"`
	Subscriptions []string `yaml:"inputs"`
	ObjectStore   *string  `yaml:"object_store,omitempty"`
}
