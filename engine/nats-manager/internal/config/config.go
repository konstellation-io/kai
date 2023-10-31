package config

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	DevelopmentMode           = "DEVELOPMENT_MODE"
	NatsManagerPort           = "NATS_MANAGER_PORT"
	NatsURL                   = "NATS_URL"
	ObjectStoreDefaultTTLDays = "OBJECT_STORE_DEFAULT_TTL"
)

func Initialize() {
	viper.SetDefault(DevelopmentMode, true)
	viper.SetDefault(NatsManagerPort, 50051)
	viper.SetDefault(NatsURL, "localhost:4222")
	viper.SetDefault(ObjectStoreDefaultTTLDays, 5)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("KAI")
	viper.AutomaticEnv()
}
