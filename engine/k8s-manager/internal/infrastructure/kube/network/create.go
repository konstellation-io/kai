package network

import (
	"context"
	"fmt"
	"strings"

	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/service"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/config"
	"github.com/spf13/viper"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	_servicePortName           = "trigger"
	_metricsPortName           = "metrics"
	_serviceProtocolAnnotation = "konghq.com/protocol"
)

func (kn KubeNetwork) CreateNetwork(ctx context.Context, params service.CreateNetworkParams) error {
	kn.logger.Info("Creating network service",
		"product", params.Product,
		"version", params.Version,
		"workflow", params.Workflow,
		"process", params.Process.Name,
		"protocol", params.Process.Networking.Protocol,
	)

	networking := params.Process.Networking

	_, err := kn.client.CoreV1().Services(kn.namespace).Create(ctx, &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: kn.getServiceName(params.Product, params.Version, params.Workflow, params.Process.Name),
			Labels: kn.getServiceLabels(
				params.Product, params.Version, params.Workflow, params.Process.Name, string(params.Process.Networking.Protocol)),
			Annotations: kn.getServiceAnnotations(networking.Protocol),
		},
		Spec: corev1.ServiceSpec{
			Selector: kn.getSelector(params.Product, params.Version, params.Workflow, params.Process.Name),
			Type:     corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       _servicePortName,
					TargetPort: intstr.FromInt(networking.TargetPort),
					Port:       int32(networking.SourcePort),
				},
				{
					Name:       _metricsPortName,
					TargetPort: intstr.FromInt(viper.GetInt(config.TelegrafMetricsOutputPortKey)),
					Port:       viper.GetInt32(config.TelegrafMetricsOutputPortKey),
				},
			},
		},
	}, metav1.CreateOptions{})

	return err
}

func (kn KubeNetwork) getSelector(product, version, workflow, process string) map[string]string {
	return map[string]string{
		"product":  product,
		"version":  version,
		"workflow": workflow,
		"process":  process,
	}
}

func (kn KubeNetwork) getServiceLabels(product, version, workflow, process, protocol string) map[string]string {
	return map[string]string{
		"product":  product,
		"version":  version,
		"workflow": workflow,
		"process":  process,
		"protocol": protocol,
		"type":     "network",
	}
}

func (kn KubeNetwork) getServiceAnnotations(protocol domain.NetworkingProtocol) map[string]string {
	annotations := make(map[string]string)

	if protocol == domain.NetworkingProtocolGRPC {
		annotations[_serviceProtocolAnnotation] = strings.ToLower(string(domain.NetworkingProtocolGRPC))
	}

	return annotations
}

func (kn KubeNetwork) getServiceName(product, version, workflow, process string) string {
	fullName := fmt.Sprintf("%s-%s-%s-%s", product, version, workflow, process)

	return strings.ReplaceAll(fullName, ".", "-")
}
