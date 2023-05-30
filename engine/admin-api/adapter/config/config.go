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
	LogLevel        string `yaml:"logLevel" envconfig:"KRE_API_LOG_LEVEL" default:"DEBUG"`
	DevelopmentMode bool   `yaml:"developmentMode" envconfig:"KRE_DEVELOPMENT_MODE"`
	ReleaseName     string `yaml:"releaseName" envconfig:"KRE_RELEASE_NAME"`
	BaseDomainName  string `yaml:"baseDomainName" envconfig:"KRE_BASE_DOMAIN_NAME"`

	Application struct {
		VersionStatusTimeout time.Duration `yaml:"versionStatusTimeout"`
	} `yaml:"application"`

	Admin struct {
		APIAddress  string `yaml:"apiAddress" envconfig:"KRE_ADMIN_API_ADDRESS"`
		BaseURL     string `yaml:"baseURL" envconfig:"KRE_ADMIN_API_BASE_URL"`
		CORSEnabled bool   `yaml:"corsEnabled" envconfig:"KRE_ADMIN_CORS_ENABLED"`
		StoragePath string `yaml:"storagePath" envconfig:"KRE_ADMIN_STORAGE_PATH"`
	} `yaml:"admin"`
	MongoDB struct {
		Address             string `yaml:"address" envconfig:"KRE_MONGODB_URI"`
		DBName              string `yaml:"dbName"`
		RuntimeDataUser     string `yaml:"runtimeDataUser" envconfig:"KRE_MONGODB_MONGOEXPRESS_USERNAME"`
		KRTBucket           string `yaml:"krtBucket"`
		MongoExpressAddress string `yaml:"mongoExpressAddress" envconfig:"KRE_MONGODB_MONGOEXPRESS_ADDRESS"`
	} `yaml:"mongodb"`

	InfluxDB struct {
		Address string `yaml:"address" envconfig:"KRE_INFLUXDB_ADDRESS"`
	} `yaml:"influxdb"`

	Chronograf struct {
		Address string `yaml:"address" envconfig:"KRE_CHRONOGRAF_ADDRESS"`
	} `yaml:"chronograf"`

	K8s struct {
		Namespace string `yaml:"namespace" envconfig:"POD_NAMESPACE"`
	} `yaml:"k8s"`

	Services struct {
		K8sManager  string `yaml:"k8sManager" envconfig:"KRE_SERVICES_K8S_MANAGER"`
		NatsManager string `yaml:"natsManager" envconfig:"KRE_SERVICES_NATS_MANAGER"`
	} `yaml:"services"`

	Keycloak KeycloakConfig `yaml:"keycloak"`
}

// TODO: Get into an agreement with infra
//
//nolint:godox // this is a task to be done
type KeycloakConfig struct {
	URL           string `yaml:"base_url" envconfig:"KEYCLOAK_BASE_URL"`
	Realm         string `yaml:"realm" envconfig:"KEYCLOAK_REALM"`
	MasterRealm   string `yaml:"master_realm" envconfig:"KEYCLOAK_MASTER_REALM"`
	AdminUsername string `yaml:"admin_username" envconfig:"KEYCLOAK_ADMIN_USERNAME"`
	AdminPassword string `yaml:"admin_password" envconfig:"KEYCLOAK_ADMIN_PASSWORD"`
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
