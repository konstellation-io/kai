package process

import (
	corev1 "k8s.io/api/core/v1"
)

const (
	WaitForDeployments = iota
	WaitForConfigMaps
	ResourceNameNvidia corev1.ResourceName = "nvidia.com/gpu"
	ResourceNameKstGpu corev1.ResourceName = "konstellation.io/gpu"
)

func (kp *KubeProcess) getNodeSelector(isGPUEnabled bool) map[string]string {
	if !isGPUEnabled {
		return map[string]string{}
	}

	return map[string]string{
		ResourceNameKstGpu.String(): "true",
	}
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
