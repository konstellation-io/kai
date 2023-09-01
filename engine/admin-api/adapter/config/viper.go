package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	ComponentsKey = "components"

	CfgFilePathKey = "CONFIG_FILE_PATH"
	RegistryURLKey = "registry.url"
)

func InitConfig() error {
	setDefaultConfig()
	return loadConfig()
}

func setDefaultConfig() {
	viper.SetDefault(CfgFilePathKey, "config.yml")
}

func loadConfig() error {
	viper.SetEnvPrefix("KAI")

	viper.RegisterAlias(RegistryURLKey, "REGISTRY_URL")

	viper.AutomaticEnv()

	viper.SetConfigFile(viper.GetString(CfgFilePathKey))

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}

	return nil
}
