package config

import (
	"github.com/spf13/viper"
)

const (
	DevelopmentMode          = "development_mode"
	NatsManagerPort          = "server.port"
	NatsURL                  = "nats.url"
	ObjectStoreDefaultTTLMin = "object_store.default_ttl_min"
)

func Initialize() {
	viper.AutomaticEnv()

	viper.RegisterAlias("KRE_DEVELOPMENT_MODE", DevelopmentMode)
	viper.RegisterAlias("KRE_NATS_MANAGER_PORT", NatsManagerPort)
	viper.RegisterAlias("KRE_NATS_URL", NatsURL)
	viper.RegisterAlias("KRE_OBJECT_STORE_DEFAULT_TTL", ObjectStoreDefaultTTLMin)

	viper.SetDefault(DevelopmentMode, false)
	viper.SetDefault(NatsManagerPort, 50051)
	viper.SetDefault(NatsURL, "localhost:4222")
	viper.SetDefault(ObjectStoreDefaultTTLMin, 5)
}
