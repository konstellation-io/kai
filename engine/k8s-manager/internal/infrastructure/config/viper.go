package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

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
	MinioSSLEnabledKey      = "minio.ssl"

	AuthEndpointKey = "auth.endpoint"
	AuthRealmKey    = "auth.realm"
	AuthClientIDKey = "auth.clientID"
	//nolint:gosec // False positive
	AuthClientSecretKey = "auth.clientSecret"

	TriggersRequestTimeoutKey         = "networking.trigger.requestTimeout"
	TriggersB64IngressesAnnotaionsKey = "networking.trigger.b64Annotations"
	TriggersIngressClassNameKey       = "networking.trigger.ingressClassName"
	TriggersTLSEnabledKey             = "networking.trigger.tls.isEnabled"
	TLSSecretNameKey                  = "networking.trigger.tls.secretName"

	ProcessTimeoutKey         = "processes.timeout"
	AutoscaleCPUPercentageKey = "autoescale.cpu.percentage"

	FluentBitImageKey      = "fluentbit.image"
	FluentBitTagKey        = "fluentbit.tag"
	FluentBitPullPolicyKey = "fluentbit.pullPolicy"

	TelegrafImageKey      = "telegraf.image"
	TelegrafTagKey        = "telegraf.tag"
	TelegrafPullPolicyKey = "telegraf.pullPolicy"
	TelegrafMetricsPort   = "telegraf.port"

	PrometheusURLKey = "prometheus.url"

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

	viper.RegisterAlias(AuthEndpointKey, "KEYCLOAK_BASE_URL")
	viper.RegisterAlias(AuthRealmKey, "KEYCLOAK_REALM")
	viper.RegisterAlias(AuthClientIDKey, "KEYCLOAK_MINIO_CLIENT_ID")
	viper.RegisterAlias(AuthClientSecretKey, "KEYCLOAK_MINIO_CLIENT_SECRET")

	viper.RegisterAlias(TriggersTLSEnabledKey, "TRIGGERS_TLS_ENABLED")
	viper.RegisterAlias(TriggersTLSEnabledKey, "TRIGGERS_TLS_CERT_SECRET_NAME")
	viper.RegisterAlias(TriggersIngressClassNameKey, "TRIGGERS_INGRESS_CLASS_NAME")
	viper.RegisterAlias(TriggersRequestTimeoutKey, "TRIGGERS_REQUEST_TIMEOUT")
	viper.RegisterAlias(TriggersB64IngressesAnnotaionsKey, "TRIGGERS_BASE64_INGRESSES_ANNOTATIONS")

	viper.RegisterAlias(AutoscaleCPUPercentageKey, "AUTOSCALE_CPU_PERCENTAGE")

	viper.RegisterAlias(FluentBitImageKey, "FLUENTBIT_IMAGE_REPOSITORY")
	viper.RegisterAlias(FluentBitTagKey, "FLUENTBIT_IMAGE_TAG")
	viper.RegisterAlias(FluentBitPullPolicyKey, "FLUENTBIT_IMAGE_PULLPOLICY")

	viper.RegisterAlias(TelegrafImageKey, "TELEGRAF_IMAGE_REPOSITORY")
	viper.RegisterAlias(TelegrafTagKey, "TELEGRAF_IMAGE_TAG")
	viper.RegisterAlias(TelegrafPullPolicyKey, "TELEGRAF_IMAGE_PULLPOLICY")
	viper.RegisterAlias(TelegrafMetricsPort, "TELEGRAF_METRICS_PORT")

	viper.RegisterAlias(PrometheusURLKey, "PROMETHEUS_URL")

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
	viper.SetDefault(MinioSSLEnabledKey, false)

	viper.SetDefault("kubernetes.isInsideCluster", true)
	viper.SetDefault(KubeNamespaceKey, "kai")

	viper.SetDefault(AutoscaleCPUPercentageKey, 80)
	viper.SetDefault(ProcessTimeoutKey, 5*time.Minute)

	viper.SetDefault(FluentBitImageKey, "fluent/fluent-bit")
	viper.SetDefault(FluentBitTagKey, "1.3")
	viper.SetDefault(FluentBitPullPolicyKey, "IfNotPresent")

	viper.SetDefault(TelegrafImageKey, "telegraf")
	viper.SetDefault(TelegrafTagKey, "1.28.5")
	viper.SetDefault(TelegrafPullPolicyKey, "IfNotPresent")
	viper.SetDefault(TelegrafMetricsPort, 9191)

	userHome, ok := os.LookupEnv("HOME")
	if ok {
		viper.SetDefault(KubeConfigPathKey, filepath.Join(userHome, ".kube", "config"))
	}
}
