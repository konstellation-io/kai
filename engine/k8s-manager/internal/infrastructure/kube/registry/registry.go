package registry

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/config"
	"github.com/spf13/viper"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	watchtools "k8s.io/client-go/tools/watch"
	"k8s.io/utils/pointer"
)

type jobStatus int32

const (
	_kaiFolder = ".kai"

	_jobStatusUnknown = iota
	_jobStatusFailed
	_jobStatusComplete

	_ttlSecondsAfterFinishedJob = 100

	_registryAuthSecretVolume  = "registry-auth-secret"  //nolint:gosec // False positive
	_registryNetrcSecretVolume = "registry-netrc-secret" //nolint:gosec // False positive
)

var (
	ErrFailedImageBuild       = errors.New("error building image")
	ErrParsingJob             = errors.New("unable to parse Kubernetes Job from Annotation watcher")
	ErrErrorEvent             = errors.New("error event received")
	ErrUnexpectedContextClose = errors.New("unexpected context close")
)

type KanikoImageBuilder struct {
	logger    logr.Logger
	client    kubernetes.Interface
	namespace string
}

func NewKanikoImageBuilder(logger logr.Logger, client kubernetes.Interface) *KanikoImageBuilder {
	return &KanikoImageBuilder{
		logger:    logger,
		client:    client,
		namespace: viper.GetString(config.KubeNamespaceKey),
	}
}

func (ib *KanikoImageBuilder) BuildImage(ctx context.Context, productID, processID, processImage string) (string, error) {
	jobName := ib.getJobNameForImage(processID)

	job := ib.getImageBuilderJob(productID, jobName, processImage)

	_, err := ib.client.BatchV1().Jobs(ib.namespace).Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		return "", fmt.Errorf("creating image building job: %w", err)
	}

	defer func() {
		deletePol := metav1.DeletePropagationBackground

		err := ib.client.BatchV1().Jobs(ib.namespace).Delete(ctx, job.Name, metav1.DeleteOptions{PropagationPolicy: &deletePol})
		if err != nil && !apierrors.IsNotFound(err) {
			ib.logger.Error(err, "Error deleting job, try to delete it manually", "job", job.Name)
			return
		}

		ib.logger.Info("Job successfully deleted", "job", job.Name)
	}()

	err = ib.watchForJob(ctx, job)
	if err != nil {
		return "", err
	}

	return processImage, nil
}

func (ib *KanikoImageBuilder) getJobNameForImage(imageName string) string {
	normalizedImageName := regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(imageName, "-")
	return fmt.Sprintf("image-builder-%s", normalizedImageName)
}

func (ib *KanikoImageBuilder) getImageBuilderJob(productID, jobName, imageWithDestination string) *batchv1.Job {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
			Labels: map[string]string{
				"job-id": jobName,
			},
		},
		Spec: batchv1.JobSpec{
			BackoffLimit:            pointer.Int32(0),
			TTLSecondsAfterFinished: pointer.Int32(_ttlSecondsAfterFinishedJob),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "kaniko",
							Image:           fmt.Sprintf("%s:%s", viper.GetString(config.ImageBuilderImageKey), viper.GetString(config.ImageBuilderTagKey)),
							ImagePullPolicy: corev1.PullPolicy(viper.GetString(config.ImageBuilderPullPolicyKey)),
							Command:         nil,
							Args: []string{
								fmt.Sprintf("--context=s3://%s/%s/%s", productID, _kaiFolder, imageWithDestination),
								fmt.Sprintf("--insecure=%s", viper.GetString(config.ImageRegistryInsecureKey)),
								fmt.Sprintf("--verbosity=%s", viper.GetString(config.ImageBuilderLogLevel)),
								fmt.Sprintf("--destination=%s", imageWithDestination),
							},
							Env: []corev1.EnvVar{
								{
									Name:  "S3_ENDPOINT",
									Value: "http://" + viper.GetString(config.MinioEndpointKey),
								},
								{
									Name:  "AWS_ACCESS_KEY_ID",
									Value: viper.GetString(config.MinioAccessKeyIDKey),
								},
								{
									Name:  "AWS_SECRET_ACCESS_KEY",
									Value: viper.GetString(config.MinioAccessKeySecretKey),
								},
								{
									Name:  "AWS_REGION",
									Value: viper.GetString(config.MinioRegionKey),
								},
								{
									Name:  "S3_FORCE_PATH_STYLE",
									Value: "true",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      _registryAuthSecretVolume,
									MountPath: "/kaniko/.docker",
								},
								{
									Name:      _registryNetrcSecretVolume,
									MountPath: "/kaniko/.netrc",
								},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
					Volumes: []corev1.Volume{
						ib.getRegistryAuthVolume(),
						ib.getRegistryNetrcVolume(),
					},
				},
			},
		},
	}
}

func (ib *KanikoImageBuilder) getRegistryAuthVolume() corev1.Volume {
	return corev1.Volume{
		Name: _registryAuthSecretVolume,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: viper.GetString(config.ImageRegistryAuthSecretKey),
				Items: []corev1.KeyToPath{
					{
						Key:  ".dockerconfigjson",
						Path: "config.json",
					},
				},
			},
		},
	}
}

func (ib *KanikoImageBuilder) getRegistryNetrcVolume() corev1.Volume {
	return corev1.Volume{
		Name: _registryNetrcSecretVolume,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: viper.GetString(config.ImageRegistryAuthSecretKey),
				Items: []corev1.KeyToPath{
					{
						Key:  ".netrcconfig",
						Path: ".netrc",
					},
				},
			},
		},
	}
}

func (ib *KanikoImageBuilder) watchForJob(ctx context.Context, job *batchv1.Job) error {
	wCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	rw, err := watchtools.NewRetryWatcher("1", &cache.ListWatch{
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return ib.client.BatchV1().Jobs(ib.namespace).Watch(wCtx, metav1.ListOptions{
				LabelSelector: "job-id=" + job.Name,
			})
		},
	})
	if err != nil {
		return fmt.Errorf("error creating label watcher: %w", err)
	}

	defer rw.Stop()

	for {
		select {
		case event := <-rw.ResultChan():
			done, err := ib.handleJobEvent(event)
			if err != nil {
				return err
			}

			if done {
				return nil
			}
		case <-rw.Done():
			return nil
		case <-wCtx.Done():
			return ErrUnexpectedContextClose
		}
	}
}

func (ib *KanikoImageBuilder) handleJobEvent(event watch.Event) (bool, error) {
	//nolint:exhaustive // Not all event types needs a specific case
	switch event.Type {
	case watch.Added, watch.Modified:
		job, ok := event.Object.(*batchv1.Job)
		if !ok {
			return true, ErrParsingJob
		}

		status := ib.getJobStatus(job)

		switch status {
		case _jobStatusFailed:
			return true, ErrFailedImageBuild
		case _jobStatusComplete:
			return true, nil
		default:
			return false, nil
		}

	case watch.Error:
		return true, ErrErrorEvent

	default:
		ib.logger.Info("Unknown event received")
		return false, nil
	}
}

func (ib *KanikoImageBuilder) getJobStatus(job *batchv1.Job) jobStatus {
	for _, condition := range job.Status.Conditions {
		if condition.Status != corev1.ConditionTrue {
			return _jobStatusUnknown
		}

		//nolint:exhaustive // Not all condition types needs a specific case
		switch condition.Type {
		case batchv1.JobFailed:
			return _jobStatusFailed
		case batchv1.JobComplete:
			return _jobStatusComplete
		default:
			return _jobStatusUnknown
		}
	}

	return _jobStatusUnknown
}
