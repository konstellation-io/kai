package configuration

import (
	"context"
	"fmt"
	"strings"

	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/config"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (kc KubeConfiguration) CreateVersionConfiguration(ctx context.Context, version domain.Version) (string, error) {
	kc.logger.Info("Creating version config files",
		"product", version.Product,
		"version", version.Tag,
	)

	processYamlConfigs := make(map[string]string, getProcessesAmount(version))

	for _, workflow := range version.Workflows {
		for _, process := range workflow.Processes {
			// This should be a MustMarshal kind of function
			processYaml, err := yaml.Marshal(kc.getProcessConfig(version, workflow, process))
			if err != nil {
				return "", err
			}

			processYamlConfigs[kc.getFullProcessIdentifier(version.Product, version.Tag, workflow.Name, process.Name)] = string(processYaml)
		}
	}

	configMap := GetAppConfig(version, processYamlConfigs)

	_, err := kc.client.CoreV1().ConfigMaps(kc.namespace).Create(ctx, &configMap, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}

	return configMap.Name, nil
}

func (kc KubeConfiguration) getFullProcessIdentifier(product, version, workflow, process string) string {
	fullName := fmt.Sprintf("%s-%s-%s-%s", product, version, workflow, process)

	return strings.ReplaceAll(fullName, ".", "-")
}

func (kc KubeConfiguration) getProcessConfig(
	version domain.Version,
	workflow *domain.Workflow,
	process *domain.Process,
) ProcessConfig {
	return ProcessConfig{
		Metadata: Metadata{
			ProductID:    version.Product,
			VersionTag:   version.Tag,
			WorkflowName: workflow.Name,
			ProcessName:  process.Name,
			ProcessType:  process.Type.ToString(),
		},
		Nats: NatsConfig{
			URL:           viper.GetString(config.NatsEndpointKey),
			Stream:        workflow.Stream,
			Subject:       process.Subject,
			Subscriptions: process.Subscriptions,
			ObjectStore:   process.ObjectStore,
		},
		CentralizedConfig: CentralizedConfig{
			Global: ConfigDefinition{
				Bucket: version.GlobalKeyValueStore,
			},
			Version: ConfigDefinition{
				Bucket: version.VersionKeyValueStore,
			},
			Workflow: ConfigDefinition{
				Bucket: workflow.KeyValueStore,
			},
			Process: ConfigDefinition{
				Bucket: process.KeyValueStore,
				Config: process.Config,
			},
		},
		Minio: MinioConfig{
			Endpoint:       viper.GetString(config.MinioEndpointKey),
			ClientUser:     version.ServiceAccount.Username,
			ClientPassword: version.ServiceAccount.Password,
			SSL:            viper.GetBool(config.MinioSSLEnabledKey),
			Bucket:         version.MinioConfiguration.Bucket,
		},
		Auth: AuthConfig{
			Endpoint:     viper.GetString(config.AuthEndpointKey),
			Client:       viper.GetString(config.AuthClientIDKey),
			ClientSecret: viper.GetString(config.AuthClientSecretKey),
			Realm:        viper.GetString(config.AuthRealmKey),
		},
		Predictions: PredictionsConfig{
			Endpoint: viper.GetString(config.PredictionsEndpointKey),
			Username: version.ServiceAccount.Username,
			Password: version.ServiceAccount.Password,
		},
	}
}

func getProcessesAmount(v domain.Version) int {
	amount := 0
	for _, w := range v.Workflows {
		amount += len(w.Processes)
	}

	return amount
}
