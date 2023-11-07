package network

import (
	"context"
	"fmt"
	"strings"

	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/service"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	_servicePortName = "trigger"
)

func (kn KubeNetwork) CreateNetwork(ctx context.Context, params service.CreateNetworkParams) error {
	kn.logger.Info("Creating network service",
		"product", params.Product,
		"version", params.Version,
		"workflow", params.Workflow,
		"process", params.Process.Name,
	)

	networking := params.Process.Networking

	_, err := kn.client.CoreV1().Services(kn.namespace).Create(ctx, &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: kn.getServiceName(params.Product, params.Version, params.Workflow, params.Process.Name),
			Labels: kn.getServiceLabels(
				params.Product, params.Version, params.Workflow, params.Process.Name, params.Process.Networking.Protocol),
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

func (kn KubeNetwork) getServiceName(product, version, workflow, process string) string {
	fullName := fmt.Sprintf("%s-%s-%s-%s", product, version, workflow, process)

	return strings.ReplaceAll(fullName, ".", "-")
}
