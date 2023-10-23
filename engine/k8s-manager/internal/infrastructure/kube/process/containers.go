package process

import (
	"fmt"
	"strings"

	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

const _configFilesVolume = "version-conf-files"

func getAppContainer(configMapName string, process *domain.Process) corev1.Container {
	var pullPolicy corev1.PullPolicy
	if process.Name == "latest" {
		pullPolicy = corev1.PullAlways
	} else {
		pullPolicy = corev1.PullIfNotPresent
	}

	container := corev1.Container{
		Name:            process.Name,
		Image:           process.Image,
		ImagePullPolicy: pullPolicy,
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
				ReadOnly:  true,
				MountPath: "/app/logs",
			},
		},
		Resources: getContainerResources(process.EnableGpu, process.ResourceLimits),
	}

	if process.Networking != nil {
		container.Ports = []corev1.ContainerPort{
			{
				ContainerPort: int32(process.Networking.SourcePort),
				Protocol:      getProtocol(process.Networking.Protocol),
			},
		}
	}

	return container
}

func getFluentBitContainer(spec *processSpec) corev1.Container {
	fluetBitImage := fmt.Sprintf("%s:%s", viper.GetString("fluentbit.image"), viper.GetString("fluentbit.tag"))
	envVars := []corev1.EnvVar{
		{Name: "KAI_MESSAGING_HOST", Value: viper.GetString("messaging.host")},
		{Name: "KAI_MESSAGING_PORT", Value: viper.GetString("messaging.port")},
		{Name: "KAI_PRODUCT_ID", Value: spec.Product},
		{Name: "KAI_VERSION_TAG", Value: spec.Version},
		{Name: "KAI_WORKFLOW_NAME", Value: spec.Workflow},
		{Name: "KAI_PROCESS_NAME", Value: spec.Process.Name},
	}

	return corev1.Container{
		Name:            "fluent-bit",
		Image:           fluetBitImage,
		ImagePullPolicy: corev1.PullIfNotPresent,
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

func getProtocol(protocol string) corev1.Protocol {
	switch strings.ToUpper(protocol) {
	case "TCP":
		return corev1.ProtocolTCP
	case "UDP":
		return corev1.ProtocolUDP
	case "SCTP":
		return corev1.ProtocolSCTP
	default:
		// Default Kubernetes value
		return corev1.ProtocolTCP
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
