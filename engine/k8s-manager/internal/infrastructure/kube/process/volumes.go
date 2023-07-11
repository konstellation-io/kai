package process

import (
	corev1 "k8s.io/api/core/v1"
)

const (
	versionConfFiles = "version-conf-files"
	appLogsVolume    = "app-log-volume"
	krtFilesVolume   = "krt-base-path"
)

func GetVolumes(configRef, configKey string) []corev1.Volume {
	return []corev1.Volume{
		{
			Name: krtFilesVolume,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
		{
			Name: versionConfFiles,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: configRef,
					},
					Items: []corev1.KeyToPath{
						{
							Key:  configKey,
							Path: "app.yaml",
						},
						{
							Key:  "parsers.conf",
							Path: "parsers.conf",
						},
						{
							Key:  "fluent-bit.conf",
							Path: "fluent-bit.conf",
						},
					},
				},
			},
		},
		{
			Name: appLogsVolume,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
	}
}
