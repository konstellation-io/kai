package process

import (
	"context"

	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/config"
	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
	autoscalilngv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

const (
	_kindDeployment   = "Deployment"
	_appsV1APIVersion = "apps/v1"
)

func (kp *KubeProcess) createAutoscaler(ctx context.Context, deployment *appsv1.Deployment, process *domain.Process) error {
	autoscaler := &autoscalilngv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:   deployment.Name,
			Labels: deployment.Labels,
			OwnerReferences: []metav1.OwnerReference{
				{
					Kind:       _kindDeployment,
					Name:       deployment.Name,
					APIVersion: _appsV1APIVersion,
					UID:        deployment.UID,
				},
			},
		},
		Spec: autoscalilngv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalilngv2.CrossVersionObjectReference{
				Kind:       _kindDeployment,
				Name:       deployment.Name,
				APIVersion: _appsV1APIVersion,
			},
			MinReplicas: pointer.Int32(1),
			MaxReplicas: process.Replicas,
			Metrics: []autoscalilngv2.MetricSpec{
				{
					Type: autoscalilngv2.ContainerResourceMetricSourceType,
					ContainerResource: &autoscalilngv2.ContainerResourceMetricSource{
						Name:      corev1.ResourceCPU,
						Container: process.Name,
						Target: autoscalilngv2.MetricTarget{
							Type:               autoscalilngv2.UtilizationMetricType,
							AverageUtilization: pointer.Int32(viper.GetInt32(config.AutoscaleCPUPercentageKey)),
						},
					},
				},
			},
		},
	}

	_, err := kp.client.AutoscalingV2().HorizontalPodAutoscalers(kp.namespace).Create(ctx, autoscaler, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}
