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
	MinioEndpointKey        = "minio.endpoint"
	MinioTierEnabledKey     = "minio.tier.enabled"
	MinioTierNameKey        = "minio.tier.name"
	MinioRootUserKey        = "minio.credentials.user"
	MinioRootPasswordKey    = "minio.credentials.password"

	KeycloakAdminUserKey     = "keycloak.admin.user"
	KeycloakAdminPasswordKey = "keycloak.admin.password"
	KeycloakAdminClientIDKey = "keycloak.admin.clientID"
	KeycloakMasterRealmKey   = "keycloak.masterRealm"
	KeycloakRealmKey         = "keycloak.realm"
	KeycloakURLKey           = "keycloak.url"
)

func InitConfig() error {
	setDefaultConfig()
	return loadConfig()
}

func setDefaultConfig() {
	viper.SetDefault(CfgFilePathKey, "config.yml")
	viper.SetDefault(VersionStatusTimeoutKey, 20*time.Minute)
	viper.SetDefault(MinioTierEnabledKey, false)
}

func loadConfig() error {
	viper.SetEnvPrefix("KAI")

	viper.RegisterAlias(RegistryHostKey, "REGISTRY_HOST")

	viper.RegisterAlias(MinioEndpointKey, "MINIO_ENDPOINT_URL")
	viper.RegisterAlias(MinioTierNameKey, "MINIO_TIER_NAME")
	viper.RegisterAlias(MinioTierEnabledKey, "MINIO_TIER_ENABLED")
	viper.RegisterAlias(MinioRootUserKey, "MINIO_ROOT_USER")
	viper.RegisterAlias(MinioRootPasswordKey, "MINIO_ROOT_PASSWORD")

	viper.RegisterAlias(KeycloakAdminUserKey, "MINIO_ROOT_PASSWORD")
	viper.RegisterAlias(KeycloakAdminPasswordKey, "MINIO_ROOT_PASSWORD")
	viper.RegisterAlias(KeycloakMasterRealmKey, "KEYCLOAK_MASTER_REALM")
	viper.RegisterAlias(KeycloakAdminClientIDKey, "KEYCLOAK_ADMIN_CLIENT_ID")
	viper.RegisterAlias(KeycloakRealmKey, "KEYCLOAK_REALM")
	viper.RegisterAlias(KeycloakURLKey, "KEYCLOAK_BASE_URL")

	viper.AutomaticEnv()

	viper.SetConfigFile(viper.GetString(CfgFilePathKey))

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}

	return nil
}
