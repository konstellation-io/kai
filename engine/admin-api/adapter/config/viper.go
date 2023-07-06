package config

import "github.com/spf13/viper"

const (
	ReleaseVersionKey = "releaseVersion"
)

func InitConfig() {
	viper.SetConfigFile("config.yml")
	viper.SetDefault(ReleaseVersionKey, "latest")
}
