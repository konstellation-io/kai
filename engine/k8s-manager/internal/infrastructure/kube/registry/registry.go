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
	_jobStatusUnknown = iota
	_jobStatusFailed
	_jobStatusComplete

	_ttlSecondsAfterFinishedJob = 100
)

var (
	ErrFailedImageBuild = errors.New("error building image")
	ErrParsingJob       = errors.New("unable to parse Kubernetes Job from Annotation watcher")
	ErrErrorEvent       = errors.New("error event received")
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

func (ib *KanikoImageBuilder) BuildImage(ctx context.Context, processID, processImage string, sources []byte) (string, error) {
	jobName := ib.getJobNameForImage(processID)
	jobConfigName := fmt.Sprintf("%s-config", jobName)

	job := ib.getImageBuilderJob(jobName, processImage, jobConfigName)

	createdJob, err := ib.client.BatchV1().Jobs(ib.namespace).Create(ctx, job, metav1.CreateOptions{})
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

	configMap := ib.getJobConfigMap(jobConfigName, createdJob, sources)

	_, err = ib.client.CoreV1().ConfigMaps(ib.namespace).Create(ctx, configMap, metav1.CreateOptions{})
	if err != nil {
		return "", fmt.Errorf("creating image builder config: %w", err)
	}

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

func (ib *KanikoImageBuilder) getImageBuilderJob(jobName, imageWithDestination, jobConfigName string) *batchv1.Job {
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
							Name:    "kaniko",
							Image:   viper.GetString(config.ImageBuilderImageKey),
							Command: nil,
							Args: []string{
								"--context=tar:///sources/file.tar.gz",
								"--insecure",
								fmt.Sprintf("--verbosity=%s", viper.GetString(config.ImageBuilderLogLevel)),
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
			return errors.New("unexpected context close")
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
