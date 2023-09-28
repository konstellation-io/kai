package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

const (
	KubeConfigPathKey  = "kubernetes.kubeConfigPath"
	KubeNamespaceKey   = "kubernetes.namespace"
	ServerPortKey      = "server.port"
	BaseDomainNameKey  = "baseDomainName"
	IsInsideClusterKey = "kubernetes.isInsideCluster"

	ImageRegistryURLKey = "registry.url"
	//nolint:gosec // False positive
	ImageRegistryAuthSecretKey = "registry.authSecret"
	ImageBuilderImageKey       = "registry.imageBuilder.image"
	ImageBuilderLogLevel       = "registry.imageBuilder.logLevel"
	ImageRegistryInsecureKey   = "registry.insecure"

	TriggerRequestTimeoutKey = "networking.trigger.requestTimeout"
	IngressClassNameKey      = "networking.trigger.ingressClassName"
	TLSIsEnabledKey          = "networking.trigger.tls.isEnabled"
	TLSSecretNameKey         = "networking.trigger.tls.secretName"

	configType = "yaml"

	_defaultServerPort     = 50051
	_defaultRequestTimeout = 5000
)

func Init(configFilePath string) error {
	configDir := filepath.Dir(configFilePath)
	configFileName := filepath.Base(configFilePath)
	fileNameWithoutExt := regexp.MustCompile(".[a-zA-Z]*$").ReplaceAllString(configFileName, "")

	viper.AddConfigPath(configDir)
	viper.SetConfigName(fileNameWithoutExt)
	viper.SetConfigType(configType)

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	viper.SetEnvPrefix("KAI")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.RegisterAlias("nats.url", "NATS_URL")
	viper.RegisterAlias(ImageRegistryURLKey, "REGISTRY_URL")
	viper.RegisterAlias(KubeNamespaceKey, "KUBERNETES_NAMESPACE")
	viper.RegisterAlias(BaseDomainNameKey, "BASE_DOMAIN_NAME")
	viper.RegisterAlias(IngressClassNameKey, "INGRESS_CLASS_NAME")
	viper.RegisterAlias(ImageRegistryAuthSecretKey, "REGISTRY_AUTH_SECRET_NAME")
	viper.RegisterAlias(ImageRegistryInsecureKey, "REGISTRY_INSECURE")

	viper.AutomaticEnv()

	setDefaultValues()

	return nil
}

func setDefaultValues() {
	viper.SetDefault("releaseName", "kai")
	viper.SetDefault("server.port", _defaultServerPort)
	viper.SetDefault(TLSIsEnabledKey, false)
	viper.SetDefault(TLSSecretNameKey, "")
	viper.SetDefault(TriggerRequestTimeoutKey, _defaultRequestTimeout)
	viper.SetDefault(IngressClassNameKey, "kong")

	viper.SetDefault(ImageBuilderImageKey, "gcr.io/kaniko-project/executor:latest")
	viper.SetDefault(ImageBuilderLogLevel, "error")

	viper.SetDefault("kubernetes.isInsideCluster", true)
	viper.SetDefault(KubeNamespaceKey, "kai")

	userHome, ok := os.LookupEnv("HOME")
	if ok {
		viper.SetDefault(KubeConfigPathKey, filepath.Join(userHome, ".kube", "config"))
	}
}
