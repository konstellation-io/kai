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
	KubeConfigPathKey    = "kubernetes.kubeConfigPath"
	KubeNamespaceKey     = "kubernetes.namespace"
	ServerPortKey        = "server.port"
	IsInsideClusterKey   = "kubernetes.isInsideCluster"
	ImageRegistryURLKey  = "registry.url"
	ImageBuilderImageKey = "registry.imageBuilder.image"
	ImageBuilderLogLevel = "registry.imageBuilder.logLevel"

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

	viper.AutomaticEnv()

	setDefaultValues()

	return nil
}

func setDefaultValues() {
	viper.SetDefault("releaseName", "kai")
	viper.SetDefault("server.port", _defaultServerPort)
	viper.SetDefault("networking.trigger.tls.isEnabled", false)
	viper.SetDefault("networking.trigger.tls.secretName", "")
	viper.SetDefault("networking.trigger.requestTimeout", _defaultRequestTimeout)
	viper.SetDefault("networking.trigger.ingressClassName", "")

	viper.SetDefault(ImageBuilderImageKey, "gcr.io/kaniko-project/executor:latest")
	viper.SetDefault(ImageBuilderLogLevel, "error")

	viper.SetDefault("kubernetes.isInsideCluster", true)
	viper.SetDefault(KubeNamespaceKey, "kai")

	userHome, ok := os.LookupEnv("HOME")
	if ok {
		viper.SetDefault(KubeConfigPathKey, filepath.Join(userHome, ".kube", "config"))
	}
}
