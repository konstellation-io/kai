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
	MinioTierEnabledKey     = "minio.tier.enabled"
	MinioTierNameKey        = "minio.tier.name"
	MinioRootUserKey        = "minio.credentials.user"
	MinioRootPasswordKey    = "minio.credentials.password"
)

func InitConfig() error {
	setDefaultConfig()
	return loadConfig()
}

func setDefaultConfig() {
	viper.SetDefault(CfgFilePathKey, "config.yml")
	viper.SetDefault(VersionStatusTimeoutKey, 20*time.Minute)
	viper.SetDefault(MinioTierEnabledKey, false)
}

func loadConfig() error {
	viper.SetEnvPrefix("KAI")

	viper.RegisterAlias(RegistryHostKey, "REGISTRY_HOST")

	viper.RegisterAlias(MinioEndpointKey, "MINIO_ENDPOINT_URL")
	viper.RegisterAlias(MinioTierNameKey, "MINIO_TIER_NAME")
	viper.RegisterAlias(MinioTierEnabledKey, "MINIO_TIER_ENABLED")
	viper.RegisterAlias(MinioRootUserKey, "MINIO_ROOT_USER")
	viper.RegisterAlias(MinioRootPasswordKey, "MINIO_ROOT_PASSWORD")

	viper.AutomaticEnv()

	viper.SetConfigFile(viper.GetString(CfgFilePathKey))

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}

	return nil
}
