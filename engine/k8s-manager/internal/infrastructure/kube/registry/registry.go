package registry

import (
	"context"
	"errors"
	"fmt"
	"net/url"
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
	_jobStatusUnknown = iota
	_jobStatusFailed
	_jobStatusComplete
)

var ErrFailedImageBuild = errors.New("error building image")

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

func (ib *KanikoImageBuilder) BuildImage(ctx context.Context, imageName string, sources []byte) (string, error) {
	jobName := ib.getJobNameForImage(imageName)
	jobConfigName := fmt.Sprintf("%s-config", jobName)

	imageWithDestination, err := ib.getImageWithDestination(imageName)
	if err != nil {
		return "", err
	}

	job := ib.getImageBuilderJob(jobName, imageWithDestination, jobConfigName)

	createdJob, err := ib.client.BatchV1().Jobs(ib.namespace).Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		return "", fmt.Errorf("creating image building job: %w", err)
	}

	defer func() {
		deletePol := metav1.DeletePropagationBackground

		if err := ib.client.BatchV1().Jobs(ib.namespace).Delete(ctx, job.Name, metav1.DeleteOptions{PropagationPolicy: &deletePol}); err != nil {
			ib.logger.Error(err, "Error deleting job, try to delete it manually", "job", job.Name)
			return
		}

		ib.logger.Info("Job successfully deleted", "job", job.Name)
	}()

	configMap := ib.getJobConfigMap(jobConfigName, createdJob, sources)

	_, err = ib.client.CoreV1().ConfigMaps(ib.namespace).Create(ctx, configMap, metav1.CreateOptions{})
	if err != nil {
		return "", fmt.Errorf("creating image builder config: %w", err)
	}

	err = ib.watchForJob(ctx, job)
	if err != nil {
		return "", err
	}

	return imageWithDestination, nil
}

func (ib *KanikoImageBuilder) getJobNameForImage(imageName string) string {
	normalizedImageName := regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(imageName, "-")
	return fmt.Sprintf("image-builder-%s", normalizedImageName)
}

func (ib *KanikoImageBuilder) getImageWithDestination(imageName string) (string, error) {
	registryURL, err := url.Parse(viper.GetString(config.ImageRegistryURLKey))
	if err != nil {
		return "", fmt.Errorf("parsing registry url: %w", err)
	}

	return fmt.Sprintf("%s/%s", registryURL.Host, imageName), nil
}

func (ib *KanikoImageBuilder) getJobConfigMap(jobConfigName string, createdJob *batchv1.Job, sources []byte) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobConfigName,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "batch/v1",
					Kind:       "Job",
					Name:       createdJob.Name,
					UID:        createdJob.UID,
				},
			},
		},
		BinaryData: map[string][]byte{
			"file.tar.gz": sources,
		},
	}
}

func (ib *KanikoImageBuilder) getImageBuilderJob(jobName string, imageWithDestination string, jobConfigName string) *batchv1.Job {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
			Labels: map[string]string{
				"job-id": jobName,
			},
		},
		Spec: batchv1.JobSpec{
			BackoffLimit:            pointer.Int32(0),
			TTLSecondsAfterFinished: pointer.Int32(100),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "kaniko",
							Image:   viper.GetString(config.ImageBuilderImageKey),
							Command: nil,
							Args: []string{
								"--context=tar:///sources/file.tar.gz",
								"--insecure",
								"--verbosity=error",
								fmt.Sprintf("--destination=%s", imageWithDestination),
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
										Name: jobConfigName,
									},
								},
							},
						},
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
		event := <-rw.ResultChan()

		job, ok := event.Object.(*batchv1.Job)
		if !ok {
			return fmt.Errorf("unable to parse Kubernetes Job from Annotation watcher")
		}

		switch event.Type {
		case watch.Added, watch.Modified:
			status := ib.getJobStatus(job)

			switch status {
			case _jobStatusFailed:
				return ErrFailedImageBuild
			case _jobStatusComplete:
				return nil
			default:
			}

		case watch.Deleted:
			return errors.New("job unexpectedly deleted")

		case watch.Error:
			ib.logger.Info("Error attempting to watch Kubernetes Jobs")

			// This round trip allows us to handle unstructured status
			errObject := apierrors.FromObject(event.Object)
			statusErr, ok := errObject.(*apierrors.StatusError)
			if !ok {
				ib.logger.Info(fmt.Sprintf("received an error which is not *metav1.Status but %+v", event.Object))
			}

			status := statusErr.ErrStatus
			ib.logger.Info("Received an error event", "status", status)
		default:
			ib.logger.Info("Unknown event received")
		}
	}
}

func (ib *KanikoImageBuilder) getJobStatus(job *batchv1.Job) jobStatus {
	for _, condition := range job.Status.Conditions {
		if condition.Status != corev1.ConditionTrue {
			return _jobStatusUnknown
		}

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
