package process

import (
	"fmt"

	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/config"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

const _configFilesVolume = "version-conf-files"

func getAppContainer(configMapName string, process *domain.Process) corev1.Container {
	container := corev1.Container{
		Name:            process.Name,
		Image:           process.Image,
		ImagePullPolicy: corev1.PullIfNotPresent,
		Env: []corev1.EnvVar{
			{
				Name:  "KAI_APP_CONFIG_PATH",
				Value: viper.GetString("configPath"),
			},
		},
		EnvFrom: []corev1.EnvFromSource{
			{
				ConfigMapRef: &corev1.ConfigMapEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: configMapName,
					},
				},
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      _configFilesVolume,
				ReadOnly:  true,
				MountPath: viper.GetString("configPath"),
			},
			{
				Name:      "app-log-volume",
				MountPath: "/var/log/app",
			},
		},
		Resources: getContainerResources(process.EnableGpu, process.ResourceLimits),
	}

	if process.Networking != nil {
		container.Ports = []corev1.ContainerPort{
			{
				ContainerPort: int32(process.Networking.SourcePort),
			},
		}
	}

	return container
}

func getFluentBitContainer(spec *processSpec) corev1.Container {
	fluentBitImage := fmt.Sprintf("%s:%s", viper.GetString(config.FluentBitImageKey), viper.GetString(config.FluentBitTagKey))
	envVars := []corev1.EnvVar{
		{Name: "KAI_LOKI_HOST", Value: viper.GetString(config.LokiHostKey)},
		{Name: "KAI_LOKI_PORT", Value: viper.GetString(config.LokiPortKey)},
		{Name: "KAI_PRODUCT_ID", Value: spec.Product},
		{Name: "KAI_VERSION_TAG", Value: spec.Version},
		{Name: "KAI_WORKFLOW_NAME", Value: spec.Workflow},
		{Name: "KAI_PROCESS_NAME", Value: spec.Process.Name},
	}

	return corev1.Container{
		Name:            "fluent-bit",
		Image:           fluentBitImage,
		ImagePullPolicy: corev1.PullPolicy(viper.GetString(config.FluentBitPullPolicyKey)),
		Command: []string{
			"/fluent-bit/bin/fluent-bit",
			"-c",
			"/fluent-bit/etc/fluent-bit.conf",
			"-v",
		},
		Env: envVars,
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      _configFilesVolume,
				ReadOnly:  true,
				MountPath: "/fluent-bit/etc/fluent-bit.conf",
				SubPath:   "fluent-bit.conf",
			},
			{
				Name:      _configFilesVolume,
				ReadOnly:  true,
				MountPath: "/fluent-bit/etc/parsers.conf",
				SubPath:   "parsers.conf",
			},
			{
				Name:      "app-log-volume",
				ReadOnly:  true,
				MountPath: "/var/log/app",
			},
		},
	}
}

func getContainerResources(isGPUEnabled bool, resourceLimits *domain.ProcessResourceLimits) corev1.ResourceRequirements {
	requests := corev1.ResourceList{
		corev1.ResourceCPU:    resource.MustParse(resourceLimits.CPU.Request),
		corev1.ResourceMemory: resource.MustParse(resourceLimits.Memory.Request),
	}
	limits := corev1.ResourceList{
		corev1.ResourceCPU:    resource.MustParse(resourceLimits.CPU.Limit),
		corev1.ResourceMemory: resource.MustParse(resourceLimits.Memory.Limit),
	}

	if isGPUEnabled {
		limits[ResourceNameNvidia] = resource.MustParse("1")
		requests[ResourceNameNvidia] = resource.MustParse("1")
	}

	return corev1.ResourceRequirements{
		Limits:   limits,
		Requests: requests,
	}
}
