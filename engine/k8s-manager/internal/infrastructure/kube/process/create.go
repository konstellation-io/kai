package process

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/service"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/config"
	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type processSpec struct {
	Product  string
	Version  string
	Workflow string
	Process  *domain.Process
}

func (kp *KubeProcess) Create(
	ctx context.Context,
	params service.CreateProcessParams,
) error {
	kp.logger.Info("Starting process",
		"product", params.Product,
		"version", params.Version,
		"process", params.Process.Name,
	)

	process := &processSpec{
		Product:  params.Product,
		Version:  params.Version,
		Workflow: params.Workflow,
		Process:  params.Process,
	}

	createdDeployment, err := kp.createProcessDeployment(ctx, params.ConfigName, process)
	if err != nil {
		return fmt.Errorf("creating deployment: %w", err)
	}

	if params.Process.Replicas > 1 {
		if err := kp.createAutoscaler(ctx, createdDeployment, params.Process); err != nil {
			return fmt.Errorf("creating autoscaler: %w", err)
		}
	}

	return nil
}

func (kp *KubeProcess) getDeploymentSpec(configMapName string, spec *processSpec) *appsv1.Deployment {
	labels := kp.getProcessLabels(spec)

	processIdentifier := getDeploymentName(spec.Product, spec.Version, spec.Workflow, spec.Process.Name)

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      processIdentifier,
			Namespace: kp.namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: kp.getPrometheusAnnotations(),
				},
				Spec: corev1.PodSpec{
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: viper.GetString(config.ImageRegistryAuthSecretKey),
						},
					},
					Containers:   kp.getContainers(configMapName, spec),
					NodeSelector: kp.getNodeSelector(spec.Process),
					Tolerations:  kp.getTolerations(spec.Process.EnableGpu),
					Volumes:      kp.getVolumes(configMapName, processIdentifier),
				},
			},
		},
	}
}

func (kp *KubeProcess) getProcessLabels(process *processSpec) map[string]string {
	return map[string]string{
		"product":  process.Product,
		"version":  process.Version,
		"workflow": process.Workflow,
		"process":  process.Process.Name,
		"type":     process.Process.Type.ToString(),
	}
}

func (kp *KubeProcess) getContainers(configmapName string, spec *processSpec) []corev1.Container {
	return []corev1.Container{
		kp.getFluentBitContainer(spec),
		kp.getAppContainer(configmapName, spec.Process),
		kp.getTelegrafContainer(),
	}
}

func (kp *KubeProcess) createProcessDeployment(ctx context.Context, configMapName string, spec *processSpec) (*appsv1.Deployment, error) {
	return kp.client.AppsV1().Deployments(kp.namespace).
		Create(ctx, kp.getDeploymentSpec(configMapName, spec), metav1.CreateOptions{})
}

func (kp *KubeProcess) getPrometheusAnnotations() map[string]string {
	return map[string]string{
		"kai.prometheus/scrape": "true",
		"kai.prometheus/scheme": "http",
		"kai.prometheus/path":   "/metrics",
		"kai.prometheus/port":   strconv.Itoa(viper.GetInt(config.TelegrafMetricsOutputPortKey)),
	}
}

func getDeploymentName(product, version, workflow, process string) string {
	fullName := fmt.Sprintf("%s-%s-%s-%s", product, version, workflow, process)

	return strings.ReplaceAll(fullName, ".", "-")
}
