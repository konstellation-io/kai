package config

import "github.com/spf13/viper"

const (
	ReleaseVersionKey   = "releaseVersion"
	VersionsFilePathKey = "VERSIONS_FILE_PATH"
)

func InitConfig() {
	viper.SetEnvPrefix("KAI")

	viper.SetConfigFile("config.yml")
	viper.SetDefault(ReleaseVersionKey, "latest")

	viper.AutomaticEnv()

	viper.SetDefault(VersionsFilePathKey, "versions.yaml")
}
