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
	MinioEndpointKey        = "minio.endpoint"
	MinioTieringEnabledKey  = "minio.tier"
	MinioTierKey            = "minio.tier"
)

func InitConfig() error {
	setDefaultConfig()
	return loadConfig()
}

func setDefaultConfig() {
	viper.SetDefault(CfgFilePathKey, "config.yml")
	viper.SetDefault(VersionStatusTimeoutKey, 20*time.Minute)
	viper.SetDefault(MinioTieringEnabledKey, false)
}

func loadConfig() error {
	viper.SetEnvPrefix("KAI")

	viper.RegisterAlias(RegistryHostKey, "REGISTRY_HOST")
	viper.RegisterAlias(MinioEndpointKey, "MINIO_ENDPOINT_URL")
	viper.RegisterAlias(MinioTierKey, "MINIO_TIER")
	viper.RegisterAlias(MinioTieringEnabledKey, "MINIO_TIERING_ENABLED")

	viper.AutomaticEnv()

	viper.SetConfigFile(viper.GetString(CfgFilePathKey))

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}

	return nil
}
