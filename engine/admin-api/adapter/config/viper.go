package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	ComponentsKey = "components"

	CfgFilePathKey = "CONFIG_FILE_PATH"
)

func InitConfig() error {
	viper.SetEnvPrefix("KAI")

	viper.SetDefault("CONFIG_FILE_PATH", "config.yml")

	viper.AutomaticEnv()

	viper.SetConfigFile(viper.GetString(CfgFilePathKey))

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}

	return nil
}
