package process

import (
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	corev1 "k8s.io/api/core/v1"
)

const (
	WaitForDeployments = iota
	WaitForConfigMaps
	ResourceNameNvidia corev1.ResourceName = "nvidia.com/gpu"
	ResourceNameKstGpu corev1.ResourceName = "konstellation.io/gpu"
)

func (kp *KubeProcess) getNodeSelector(process *domain.Process) map[string]string {
	nodeSelectors := map[string]string{}

	if process.NodeSelectors != nil {
		nodeSelectors = process.NodeSelectors
	}

	if process.EnableGpu {
		nodeSelectors[ResourceNameKstGpu.String()] = "true"
	}

	return nodeSelectors
}

func (kp *KubeProcess) getTolerations(isGPUEnabled bool) []corev1.Toleration {
	if !isGPUEnabled {
		return []corev1.Toleration{}
	}

	return []corev1.Toleration{
		{
			Key:      ResourceNameKstGpu.String(),
			Operator: corev1.TolerationOpExists,
			Effect:   corev1.TaintEffectNoSchedule,
		},
		{
			Key:      ResourceNameNvidia.String(),
			Operator: corev1.TolerationOpExists,
			Effect:   corev1.TaintEffectNoSchedule,
		},
	}
}
