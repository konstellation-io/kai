package kube

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/service"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/kube/configuration"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/kube/network"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/kube/process"
	"github.com/spf13/viper"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type K8sContainerService struct {
	logger    logr.Logger
	namespace string
	client    kubernetes.Interface

	processService       process.KubeProcess
	configurationService configuration.KubeConfiguration
	networkService       network.KubeNetwork
}

var _ service.ContainerService = (*K8sContainerService)(nil)

func NewK8sContainerService(logger logr.Logger, client kubernetes.Interface) *K8sContainerService {
	namespace := viper.GetString("kubernetes.namespace")

	return &K8sContainerService{
		logger:    logger,
		namespace: namespace,
		client:    client,

		processService:       process.NewKubeProcess(logger, client, namespace),
		configurationService: configuration.NewKubeConfiguration(logger, client, namespace),
		networkService:       network.NewKubeNetwork(logger, client, namespace),
	}
}

func (k *K8sContainerService) CreateProcess(ctx context.Context, params service.CreateProcessParams) error {
	return k.processService.Create(ctx, params)
}

func (k *K8sContainerService) DeleteProcesses(ctx context.Context, product, version string) error {
	return k.processService.DeleteProcesses(ctx, product, version)
}

func (k *K8sContainerService) CreateVersionConfiguration(ctx context.Context, version domain.Version) (string, error) {
	return k.configurationService.CreateVersionConfiguration(ctx, version)
}
func (k *K8sContainerService) DeleteConfiguration(ctx context.Context, product, version string) error {
	return k.configurationService.DeleteConfiguration(ctx, product, version)
}

func (k *K8sContainerService) CreateNetwork(ctx context.Context, params service.CreateNetworkParams) error {
	return k.networkService.CreateNetwork(ctx, params)
}

func (k *K8sContainerService) DeleteNetwork(ctx context.Context, product, version string) error {
	return k.networkService.DeleteNetwork(ctx, product, version)
}

func (k *K8sContainerService) ProcessRegister(ctx context.Context, name string, sources []byte) error {
	jobName := strings.ReplaceAll(name, "http://", "")
	jobName = strings.ReplaceAll(jobName, "https://", "")
	jobName = strings.ReplaceAll(jobName, "/", "-")
	jobName = strings.ReplaceAll(jobName, ".", "-")
	jobName = strings.ReplaceAll(jobName, ":", "-")

	fmt.Println(name)

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("image-builder-%s-config", jobName),
		},
		BinaryData: map[string][]byte{
			"file.tar.gz": sources,
		},
	}

	_, err := k.client.CoreV1().ConfigMaps(k.namespace).Create(ctx, configMap, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("creating image builder config: %w", err)
	}

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf("image-builder-%s", jobName),
			Labels: nil,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "kaniko",
							Image:   "gcr.io/kaniko-project/executor:latest",
							Command: nil,
							Args: []string{
								"--context=tar:///sources/file.tar.gz",
								"--insecure",
								fmt.Sprintf("--destination=%s", name),
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/sources",
								},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
					Volumes: []corev1.Volume{
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								HostPath: nil,
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: configMap.Name,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	_, err = k.client.BatchV1().Jobs(k.namespace).Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("creating image building job: %w", err)
	}

	//job := &corev1.Pod{
	//	ObjectMeta: metav1.ObjectMeta{
	//		Name:   fmt.Sprintf("image-builder-%s", jobName),
	//		Labels: nil,
	//	},
	//	Spec: corev1.PodSpec{
	//		Containers: []corev1.Container{
	//			{
	//				Name:    "kaniko",
	//				Image:   "debian",
	//				Command: []string{"/bin/bash", "-c", "--"},
	//				Args:    []string{"while true; do sleep 30; done;"},
	//				//Command: nil,
	//				//Args: []string{
	//				//	"--context=tar:///sources/file.tar.gz",
	//				//	fmt.Sprintf("--destination=%s", name),
	//				//},
	//				VolumeMounts: []corev1.VolumeMount{
	//					{
	//						Name:      "config",
	//						MountPath: "/sources",
	//					},
	//				},
	//			},
	//		},
	//		RestartPolicy: corev1.RestartPolicyNever,
	//		Volumes: []corev1.Volume{
	//			{
	//				Name: "config",
	//				VolumeSource: corev1.VolumeSource{
	//					HostPath: nil,
	//					ConfigMap: &corev1.ConfigMapVolumeSource{
	//						LocalObjectReference: corev1.LocalObjectReference{
	//							Name: configMap.Name,
	//						},
	//					},
	//				},
	//			},
	//		},
	//	},
	//}
	//
	//_, err = k.client.CoreV1().Pods(k.namespace).Create(ctx, job, metav1.CreateOptions{})
	//if err != nil {
	//	return fmt.Errorf("creating image building job: %w", err)
	//}

	return nil
}
