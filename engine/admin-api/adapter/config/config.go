package config

import (
	"os"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

// Config holds the configuration values for the application.
type Config struct {
	LogLevel        string `yaml:"logLevel" envconfig:"KAI_API_LOG_LEVEL" default:"DEBUG"`
	DevelopmentMode bool   `yaml:"developmentMode" envconfig:"KAI_DEVELOPMENT_MODE"`
	ReleaseName     string `yaml:"releaseName" envconfig:"KAI_RELEASE_NAME"`
	BaseDomainName  string `yaml:"baseDomainName" envconfig:"KAI_BASE_DOMAIN_NAME"`

	Application struct {
		VersionStatusTimeout time.Duration `yaml:"versionStatusTimeout"`
	} `yaml:"application"`

	Admin struct {
		APIAddress  string `yaml:"apiAddress" envconfig:"KAI_ADMIN_API_ADDRESS"`
		BaseURL     string `yaml:"baseURL" envconfig:"KAI_ADMIN_API_BASE_URL"`
		CORSEnabled bool   `yaml:"corsEnabled" envconfig:"KAI_ADMIN_CORS_ENABLED"`
		StoragePath string `yaml:"storagePath" envconfig:"KAI_ADMIN_STORAGE_PATH"`
	} `yaml:"admin"`
	MongoDB struct {
		Address             string `yaml:"address" envconfig:"KAI_MONGODB_URI"`
		DBName              string `yaml:"dbName"`
		RuntimeDataUser     string `yaml:"runtimeDataUser" envconfig:"KAI_MONGODB_MONGOEXPRESS_USERNAME"`
		KRTBucket           string `yaml:"krtBucket"`
		MongoExpressAddress string `yaml:"mongoExpressAddress" envconfig:"KAI_MONGODB_MONGOEXPRESS_ADDRESS"`
	} `yaml:"mongodb"`

	InfluxDB struct {
		Address string `yaml:"address" envconfig:"KAI_INFLUXDB_ADDRESS"`
	} `yaml:"influxdb"`

	Chronograf struct {
		Address string `yaml:"address" envconfig:"KAI_CHRONOGRAF_ADDRESS"`
	} `yaml:"chronograf"`

	K8s struct {
		Namespace string `yaml:"namespace" envconfig:"POD_NAMESPACE"`
	} `yaml:"k8s"`

	Services struct {
		K8sManager  string `yaml:"k8sManager" envconfig:"KAI_SERVICES_K8S_MANAGER"`
		NatsManager string `yaml:"natsManager" envconfig:"KAI_SERVICES_NATS_MANAGER"`
	} `yaml:"services"`
}

// NewConfig will read the config.yml file and override values with env vars.
func NewConfig() *Config {
	var once sync.Once

	var cfg *Config

	once.Do(func() {
		f, err := os.Open("config.yml")
		if err != nil {
			panic(err)
		}

		cfg = &Config{}
		decoder := yaml.NewDecoder(f)

		err = decoder.Decode(cfg)
		if err != nil {
			panic(err)
		}

		err = envconfig.Process("", cfg)
		if err != nil {
			panic(err)
		}
	})

	return cfg
}
