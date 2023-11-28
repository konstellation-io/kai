package configuration

type ProcessConfig struct {
	Metadata          Metadata           `yaml:"metadata"`
	Nats              NatsConfig         `yaml:"nats"`
	CentralizedConfig CentralizedConfig  `yaml:"centralized_configuration"`
	Minio             MinioConfig        `yaml:"minio"`
	Auth              AuthConfig         `yaml:"auth"`
	Measurements      MeasurementsConfig `yaml:"measurements"`
	Predictions       PredictionsConfig  `yaml:"predictions"`
}

type Metadata struct {
	ProductID    string `yaml:"product_id"`
	VersionTag   string `yaml:"version_tag"`
	WorkflowName string `yaml:"workflow_name"`
	ProcessName  string `yaml:"process_name"`
	BasePath     string `yaml:"base_path"`
	ProcessType  string `yaml:"process_type"`
	WorkflowType string `yaml:"workflow_type"`
}

type CentralizedConfig struct {
	Global   ConfigDefinition `yaml:"global"`
	Version  ConfigDefinition `yaml:"product"` // rename this to version
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

type MinioConfig struct {
	Endpoint       string `yaml:"endpoint"`
	ClientUser     string `yaml:"client_user"`
	ClientPassword string `yaml:"client_password"`
	SSL            bool   `yaml:"ssl"`
	Bucket         string `yaml:"bucket"`
}

type AuthConfig struct {
	Endpoint     string `yaml:"endpoint"`
	Client       string `yaml:"client"`
	ClientSecret string `yaml:"client_secret"`
	Realm        string `yaml:"realm"`
}

type MeasurementsConfig struct {
	Endpoint        string `yaml:"endpoint"`
	Insecure        bool   `yaml:"insecure"`
	Timeout         int    `yaml:"timeout"`
	MetricsInterval int    `yaml:"metrics_interval"`
}

type PredictionsConfig struct {
	Endpoint string `yaml:"endpoint"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
