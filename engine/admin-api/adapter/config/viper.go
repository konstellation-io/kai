package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	ComponentsKey = "components"

	CfgFilePathKey = "CONFIG_FILE_PATH"

	RegistryURLKey         = "registry.url"
	RegistriesConfPathKey  = "registry.registriesConfPath"
	SignaturePolicyPathKey = "registry.signaturePolicyPath"
)

func InitConfig() error {
	setDefaultConfig()
	registerAliases()

	return loadConfig()

}

func setDefaultConfig() {
	viper.SetDefault(CfgFilePathKey, "config.yml")
	viper.SetDefault(RegistryURLKey, "localhost:5000")
	viper.SetDefault(RegistriesConfPathKey, "/etc/containers/registries.conf")
	viper.SetDefault(SignaturePolicyPathKey, "/etc/containers/policy.json")
}

func registerAliases() {
	viper.RegisterAlias("REGISTRY_URL", RegistriesConfPathKey)
}

func loadConfig() error {
	viper.SetEnvPrefix("KAI")
	viper.AutomaticEnv()

	viper.SetConfigFile(viper.GetString(CfgFilePathKey))

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}

	return nil
}
