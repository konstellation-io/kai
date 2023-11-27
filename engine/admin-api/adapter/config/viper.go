package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

const (
	ComponentsKey = "components"

	CfgFilePathKey          = "CONFIG_FILE_PATH"
	RegistryHostKey         = "registry.host"
	VersionStatusTimeoutKey = "application.versionStatusTimeout"

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

	MongoDBKaiDatabaseKey = "mongodb.dbName"

	RedisEndpointKey         = "redis.endpoint"
	RedisPasswordKey         = "redis.password"
	RedisPredictionsIndexKey = "redis.predictionsIndexName"
)

func InitConfig() error {
	setDefaultConfig()
	return loadConfig()
}

func setDefaultConfig() {
	viper.SetDefault(CfgFilePathKey, "config.yml")
	viper.SetDefault(VersionStatusTimeoutKey, 20*time.Minute)
	viper.SetDefault(MinioTierEnabledKey, false)
	viper.SetDefault(MinioTierTransitionDaysKey, 0)
	viper.SetDefault(KeycloakPolicyAttributeKey, "policy")
	viper.SetDefault(MongoDBKaiDatabaseKey, "kai")
	viper.SetDefault(RedisPredictionsIndexKey, "predictionsIdx")
}

func loadConfig() error {
	viper.SetEnvPrefix("KAI")

	viper.RegisterAlias(RegistryHostKey, "REGISTRY_HOST")

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
	viper.RegisterAlias(RedisPasswordKey, "REDIS_PASSWORD")
	viper.RegisterAlias(RedisPredictionsIndexKey, "REDIS_PREDICTIONS_INDEX")

	viper.AutomaticEnv()

	viper.SetConfigFile(viper.GetString(CfgFilePathKey))

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}

	return nil
}
