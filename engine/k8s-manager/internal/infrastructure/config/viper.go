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

	MinioEndpointKey        = "minio.endpoint"
	MinioAccessKeyIDKey     = "minio.accessKeyID"
	MinioAccessKeySecretKey = "minio.accessKeySecret"
	MinioRegionKey          = "minio.regionKey"

	TriggersRequestTimeoutKey         = "networking.trigger.requestTimeout"
	TriggersB64IngressesAnnotaionsKey = "networking.trigger.b64Annotations"
	TriggersIngressClassNameKey       = "networking.trigger.ingressClassName"
	TriggersTLSEnabledKey             = "networking.trigger.tls.isEnabled"
	TLSSecretNameKey                  = "networking.trigger.tls.secretName"

	AutoscaleCPUPercentageKey = "autoescale.cpu.percentage"

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

	viper.RegisterAlias(ImageRegistryAuthSecretKey, "REGISTRY_AUTH_SECRET_NAME")
	viper.RegisterAlias(ImageRegistryInsecureKey, "REGISTRY_INSECURE")
	viper.RegisterAlias(ImageBuilderImageKey, "IMAGE_BUILDER_IMAGE")

	viper.RegisterAlias(MinioEndpointKey, "MINIO_ENDPOINT_URL")
	viper.RegisterAlias(MinioAccessKeyIDKey, "MINIO_ROOT_USER")
	viper.RegisterAlias(MinioAccessKeySecretKey, "MINIO_ROOT_PASSWORD")
	viper.RegisterAlias(MinioRegionKey, "MINIO_REGION")

	viper.RegisterAlias(TriggersTLSEnabledKey, "TRIGGERS_TLS_ENABLED")
	viper.RegisterAlias(TriggersTLSEnabledKey, "TRIGGERS_TLS_CERT_SECRET_NAME")
	viper.RegisterAlias(TriggersIngressClassNameKey, "TRIGGERS_INGRESS_CLASS_NAME")
	viper.RegisterAlias(TriggersRequestTimeoutKey, "TRIGGERS_REQUEST_TIMEOUT")
	viper.RegisterAlias(TriggersB64IngressesAnnotaionsKey, "TRIGGERS_BASE64_INGRESSES_ANNOTATIONS")

	viper.RegisterAlias(AutoscaleCPUPercentageKey, "AUTOSCALE_CPU_PERCENTAGE")

	viper.AutomaticEnv()

	setDefaultValues()

	return nil
}

func setDefaultValues() {
	viper.SetDefault("releaseName", "kai")
	viper.SetDefault("server.port", _defaultServerPort)
	viper.SetDefault(TriggersTLSEnabledKey, false)
	viper.SetDefault(TLSSecretNameKey, "")
	viper.SetDefault(TriggersRequestTimeoutKey, _defaultRequestTimeout)
	viper.SetDefault(TriggersIngressClassNameKey, "kong")

	viper.SetDefault(ImageBuilderImageKey, "gcr.io/kaniko-project/executor:latest")
	viper.SetDefault(ImageBuilderLogLevel, "error")
	viper.SetDefault(MinioRegionKey, "us-east-1")

	viper.SetDefault("kubernetes.isInsideCluster", true)
	viper.SetDefault(KubeNamespaceKey, "kai")

	viper.SetDefault(AutoscaleCPUPercentageKey, 80)

	userHome, ok := os.LookupEnv("HOME")
	if ok {
		viper.SetDefault(KubeConfigPathKey, filepath.Join(userHome, ".kube", "config"))
	}
}
