package config

import (
	"github.com/spf13/viper"
)

const (
	DEVELOPMENT_MODE             = "development_mode"
	NATS_MANAGER_PORT            = "server.port"
	NATS_URL                     = "nats.url"
	OBJECT_STORE_DEFAULT_TTL_MIN = "object_store.default_ttl_min"
)

func Initialize() {
	viper.AutomaticEnv()

	viper.RegisterAlias("KRE_DEVELOPMENT_MODE", DEVELOPMENT_MODE)
	viper.RegisterAlias("KRE_NATS_MANAGER_PORT", NATS_MANAGER_PORT)
	viper.RegisterAlias("KRE_NATS_URL", NATS_URL)
	viper.RegisterAlias("KRE_OBJECT_STORE_DEFAULT_TTL", OBJECT_STORE_DEFAULT_TTL_MIN)

	viper.SetDefault(DEVELOPMENT_MODE, false)
	viper.SetDefault(NATS_MANAGER_PORT, 50051)
	viper.SetDefault(NATS_URL, "localhost:4222")
	viper.SetDefault(OBJECT_STORE_DEFAULT_TTL_MIN, 5)
}
