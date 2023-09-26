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
)

func InitConfig() error {
	setDefaultConfig()
	return loadConfig()
}

func setDefaultConfig() {
	viper.SetDefault(CfgFilePathKey, "config.yml")
	viper.SetDefault(VersionStatusTimeoutKey, 20*time.Minute)
}

func loadConfig() error {
	viper.SetEnvPrefix("KAI")

	viper.RegisterAlias(RegistryHostKey, "REGISTRY_HOST")

	viper.AutomaticEnv()

	viper.SetConfigFile(viper.GetString(CfgFilePathKey))

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}

	return nil
}
