package process

import (
	"fmt"

	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
)

func GeKrtFilesDownloaderContainer(commonEnvVars []corev1.EnvVar) corev1.Container {
	image := fmt.Sprintf("%s:%s", viper.Get("krtFilesDownloader.image"), viper.Get("krtFilesDownloader.tag"))

	return corev1.Container{
		Name:  "krt-files-downloader",
		Image: image,
		Env:   commonEnvVars,
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "krt-base-path",
				ReadOnly:  false,
				MountPath: "krt-files",
			},
		},
	}
}
