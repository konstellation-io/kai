package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

const (
	ComponentsKey = "components"

	CfgFilePathKey          = "CONFIG_FILE_PATH"
	RegistryHostKey         = "registry.host"
	VersionStatusTimeoutKey = "application.versionStatusTimeout"
	S3EndpointKey           = "s3.endpoint"
	S3BucketKey             = "s3.bucket"
	S3TierKey               = "s3.tier"
)

func InitConfig() error {
	setDefaultConfig()
	return loadConfig()
}

func setDefaultConfig() {
	viper.SetDefault(CfgFilePathKey, "config.yml")
	viper.SetDefault(VersionStatusTimeoutKey, 20*time.Minute)
	viper.SetDefault(S3BucketKey, "kai")
}

func loadConfig() error {
	viper.SetEnvPrefix("KAI")

	viper.RegisterAlias(RegistryHostKey, "REGISTRY_HOST")
	viper.RegisterAlias(S3EndpointKey, "S3_ENDPOINT_URL")
	viper.RegisterAlias(S3BucketKey, "S3_BUCKET")
	viper.RegisterAlias(S3BucketKey, "S3_TIER")

	viper.AutomaticEnv()

	viper.SetConfigFile(viper.GetString(CfgFilePathKey))

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}

	return nil
}
