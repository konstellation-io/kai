package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

const (
	ComponentsKey = "components"

	CfgFilePathKey          = "CONFIG_FILE_PATH"
	LogLevelKey             = "application.logLevel"
	VersionStatusTimeoutKey = "application.versionStatusTimeout"
	ApplicationPortKey      = "application.port"
	CORSEnabledKey          = "application.corsEnabled"

	RegistryHostKey       = "registry.host"
	GlobalRegistryKey     = "registry.global"
	RegistryAuthSecretKey = "registry.basicAuthSecret"

	MinioEndpointKey           = "minio.endpoint"
	MinioTierEnabledKey        = "minio.tier.enabled"
	MinioTierNameKey           = "minio.tier.name"
	MinioTierTransitionDaysKey = "minio.tier.transition.days"
	MinioRootUserKey           = "minio.credentials.user"
	//nolint:gosec // false positive
	MinioRootPasswordKey = "minio.credentials.password"

	KeycloakAdminUserKey = "keycloak.admin.user"
	//nolint:gosec // false positive
	KeycloakAdminPasswordKey   = "keycloak.admin.password"
	KeycloakAdminClientIDKey   = "keycloak.admin.clientID"
	KeycloakMasterRealmKey     = "keycloak.masterRealm"
	KeycloakRealmKey           = "keycloak.realm"
	KeycloakURLKey             = "keycloak.url"
	KeycloakPolicyAttributeKey = "keycloak.attributes.policy"

	MongoDBEndpointKey    = "mongodb.endpoint"
	MongoDBKaiDatabaseKey = "mongodb.dbName"

	RedisEndpointKey         = "redis.endpoint"
	RedisUsernameKey         = "redis.username"
	RedisPasswordKey         = "redis.password"
	RedisPredictionsIndexKey = "redis.predictionsIndexName"

	K8sManagerEndpointKey  = "services.k8sManager.endpoint"
	NatsManagerEndpointKey = "services.natsManager.endpoint"

	LokiEndpointKey = "loki.endpoint"
)

func InitConfig() error {
	viper.SetEnvPrefix("KAI")
	viper.SetDefault(CfgFilePathKey, "config.yml")

	viper.RegisterAlias(LogLevelKey, "API_LOG_LEVEL")
	viper.RegisterAlias(ApplicationPortKey, "ADMIN_API_PORT")
	viper.RegisterAlias(CORSEnabledKey, "ADMIN_CORS_ENABLED")

	viper.RegisterAlias(RegistryHostKey, "REGISTRY_HOST")
	viper.RegisterAlias(RegistryAuthSecretKey, "REGISTRY_BASIC_AUTH")

	viper.RegisterAlias(MongoDBEndpointKey, "MONGODB_URI")
	viper.RegisterAlias(MongoDBKaiDatabaseKey, "MONGODB_DATABASE")

	viper.RegisterAlias(MinioEndpointKey, "MINIO_ENDPOINT_URL")
	viper.RegisterAlias(MinioTierNameKey, "MINIO_TIER_NAME")
	viper.RegisterAlias(MinioTierEnabledKey, "MINIO_TIER_ENABLED")
	viper.RegisterAlias(MinioRootUserKey, "MINIO_ROOT_USER")
	viper.RegisterAlias(MinioRootPasswordKey, "MINIO_ROOT_PASSWORD")
	viper.RegisterAlias(MinioTierTransitionDaysKey, "MINIO_TIER_TRANSITION_DAYS")

	viper.RegisterAlias(KeycloakAdminUserKey, "KEYCLOAK_ADMIN_USERNAME")
	viper.RegisterAlias(KeycloakAdminPasswordKey, "KEYCLOAK_ADMIN_PASSWORD")
	viper.RegisterAlias(KeycloakMasterRealmKey, "KEYCLOAK_MASTER_REALM")
	viper.RegisterAlias(KeycloakAdminClientIDKey, "KEYCLOAK_ADMIN_CLIENT_ID")
	viper.RegisterAlias(KeycloakRealmKey, "KEYCLOAK_REALM")
	viper.RegisterAlias(KeycloakURLKey, "KEYCLOAK_BASE_URL")
	viper.RegisterAlias(KeycloakPolicyAttributeKey, "KEYCLOAK_POLICY_ATTRIBUTE")

	viper.RegisterAlias(RedisEndpointKey, "REDIS_MASTER_ADDRESS")
	viper.RegisterAlias(RedisUsernameKey, "REDIS_USERNAME")
	viper.RegisterAlias(RedisPasswordKey, "REDIS_PASSWORD")
	viper.RegisterAlias(RedisPredictionsIndexKey, "REDIS_PREDICTIONS_INDEX")

	viper.RegisterAlias(LokiEndpointKey, "LOKI_ADDRESS")

	viper.RegisterAlias(K8sManagerEndpointKey, "SERVICES_K8S_MANAGER")
	viper.RegisterAlias(NatsManagerEndpointKey, "SERVICES_NATS_MANAGER")

	viper.SetConfigFile(viper.GetString(CfgFilePathKey))

	err := viper.ReadInConfig()

	switch {
	case os.IsNotExist(err):
		fmt.Printf("Config file %q not found\n", viper.GetString(CfgFilePathKey))
	case err != nil:
		return fmt.Errorf("read config file: %w", err)
	}

	viper.AutomaticEnv()
	setDefaultConfig()

	return nil
}

func setDefaultConfig() {
	viper.SetDefault(LogLevelKey, "INFO")
	viper.SetDefault(VersionStatusTimeoutKey, 20*time.Minute)
	viper.SetDefault(MinioTierEnabledKey, false)
	viper.SetDefault(MinioTierTransitionDaysKey, 0)
	viper.SetDefault(KeycloakPolicyAttributeKey, "policy")
	viper.SetDefault(MongoDBKaiDatabaseKey, "kai")
	viper.SetDefault(GlobalRegistryKey, "kai")
	viper.SetDefault(RedisUsernameKey, "default")
	viper.SetDefault(RedisPredictionsIndexKey, "predictionsIdx")
	viper.SetDefault(CORSEnabledKey, false)
}
