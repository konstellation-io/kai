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

	viper.RegisterAlias("KAI_DEVELOPMENT_MODE", DevelopmentMode)
	viper.RegisterAlias("KAI_NATS_MANAGER_PORT", NatsManagerPort)
	viper.RegisterAlias("KAI_NATS_URL", NatsURL)
	viper.RegisterAlias("KAI_OBJECT_STORE_DEFAULT_TTL", ObjectStoreDefaultTTLDays)

	viper.SetDefault(DevelopmentMode, false)
	viper.SetDefault(NatsManagerPort, 50051)
	viper.SetDefault(NatsURL, "localhost:4222")
	viper.SetDefault(ObjectStoreDefaultTTLDays, 5)
}
