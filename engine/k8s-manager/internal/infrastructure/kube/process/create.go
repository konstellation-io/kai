package process

import (
	"context"
	"fmt"
	"strings"

	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/service"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
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

	deploymentSpec := kp.getDeploymentSpec(params.ConfigName, &processSpec{
		Product:  params.Product,
		Version:  params.Version,
		Workflow: params.Workflow,
		Process:  params.Process,
	})

	return kp.createProcessDeployment(ctx, deploymentSpec)
}

func (kp *KubeProcess) getDeploymentSpec(configMapName string, spec *processSpec) *appsv1.Deployment {
	labels := kp.getProcessLabels(spec)

	processIdentifier := getFullProcessIdentifier(spec.Product, spec.Version, spec.Workflow, spec.Process.Name)

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      processIdentifier,
			Namespace: kp.namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32(spec.Process.Replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers:   kp.getContainers(configMapName, spec),
					NodeSelector: kp.getNodeSelector(spec.Process.EnableGpu),
					Tolerations:  kp.getTolerations(spec.Process.EnableGpu),
					Volumes:      GetVolumes(configMapName, processIdentifier),
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
		getFluentBitContainer(spec),
		getAppContainer(configmapName, spec.Process),
	}
}

func (kp *KubeProcess) createProcessDeployment(
	ctx context.Context,
	deploy *appsv1.Deployment,
) error {
	_, err := kp.client.AppsV1().Deployments(kp.namespace).Create(ctx, deploy, metav1.CreateOptions{})

	return err
}

func getFullProcessIdentifier(product, version, workflow, process string) string {
	fullName := fmt.Sprintf("%s-%s-%s-%s", product, version, workflow, process)

	return strings.ReplaceAll(fullName, ".", "-")
}
