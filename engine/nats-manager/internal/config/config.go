package config

import (
	"github.com/spf13/viper"
)

const (
	DevelopmentMode           = "development_mode"
	NatsManagerPort           = "server.port"
	NatsURL                   = "nats.url"
	ObjectStoreDefaultTTLDays = "object_store.default_ttl_days"
)

func Initialize() {
	viper.AutomaticEnv()

	viper.RegisterAlias(DevelopmentMode, "KAI_DEVELOPMENT_MODE")
	viper.RegisterAlias(NatsManagerPort, "KAI_NATS_MANAGER_PORT")
	viper.RegisterAlias(NatsURL, "KAI_NATS_URL")
	viper.RegisterAlias(ObjectStoreDefaultTTLDays, "KAI_OBJECT_STORE_DEFAULT_TTL")

	viper.SetDefault(DevelopmentMode, false)
	viper.SetDefault(NatsManagerPort, 50051)
	viper.SetDefault(NatsURL, "localhost:4222")
	viper.SetDefault(ObjectStoreDefaultTTLDays, 5)

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
