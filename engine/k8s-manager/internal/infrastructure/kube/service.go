package kube

import (
	"bytes"
	"context"
	"fmt"
	"io"
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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	watchtools "k8s.io/client-go/tools/watch"
	"k8s.io/utils/pointer"
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

	fmt.Println(jobName)
	jobName = fmt.Sprintf("image-builder-%s", jobName)

	jobConfigName := fmt.Sprintf("%s-config", jobName)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
			Labels: map[string]string{
				"job-id": jobName,
			},
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: pointer.Int32(1),
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
								"--verbosity=error",
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

	createdJob, err := k.client.BatchV1().Jobs(k.namespace).Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("creating image building job: %w", err)
	}

	defer func() {
		deletePol := metav1.DeletePropagationBackground

		if err := k.client.BatchV1().Jobs(k.namespace).Delete(ctx, job.Name, metav1.DeleteOptions{PropagationPolicy: &deletePol}); err != nil {
			k.logger.Error(err, "Error deleting job, try to delete it manually", "job", job.Name)
			return
		}

		k.logger.Info("Job successfully deleted", "job", job.Name)
	}()

	configMap := &corev1.ConfigMap{
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

	_, err = k.client.CoreV1().ConfigMaps(k.namespace).Create(ctx, configMap, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("creating image builder config: %w", err)
	}

	err = k.watchForJob(ctx, jobName, job.Name)
	if err != nil {
		return err
	}

	return nil
}

func (k *K8sContainerService) watchForJob(ctx context.Context, jobID, name string) error {
	wCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	rw, err := watchtools.NewRetryWatcher("1", &cache.ListWatch{
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return k.client.BatchV1().Jobs(k.namespace).Watch(wCtx, metav1.ListOptions{
				LabelSelector: "job-id=" + jobID,
			})
		},
	})
	if err != nil {
		return fmt.Errorf("error creating label watcher: %s", err.Error())
	}

	defer rw.Stop()

	for {
		select {
		case event := <-rw.ResultChan():
			err := k.processJobEvent(wCtx, event, name)
			if err != nil {
				return err
			}
		case <-wCtx.Done():
			return nil
		}
	}
}

func (k *K8sContainerService) processJobEvent(ctx context.Context, event watch.Event, jobName string) error {

	switch event.Type {

	case watch.Added, watch.Modified:
		job, ok := event.Object.(*batchv1.Job)
		if !ok {
			return fmt.Errorf("unable to parse Kubernetes Job from Annotation watcher")
		}

		if k.isJobFailed(job) {
			fmt.Println("Job failed")

			pods, err := k.client.CoreV1().Pods(k.namespace).List(ctx, metav1.ListOptions{
				LabelSelector: fmt.Sprintf("job-name=%s", jobName),
			})
			if err != nil {
				return fmt.Errorf("getting pods: %w", err)
			}

			for _, pod := range pods.Items {
				req := k.client.CoreV1().Pods(k.namespace).GetLogs(pod.Name, &corev1.PodLogOptions{
					TailLines: pointer.Int64(1),
				})

				stream, err := req.Stream(ctx)
				if err != nil {
					return err
				}
				defer stream.Close()

				buf := new(bytes.Buffer)

				_, err = io.Copy(buf, stream)
				if err != nil {
					return err
				}
				return fmt.Errorf("job failed: %s", buf.String())
			}
		}

		ctx.Done()
		return nil

	case watch.Deleted:
		_, ok := event.Object.(*batchv1.Job)
		if !ok {
			return fmt.Errorf("unable to parse Kubernetes Job from Annotation watcher")
		}

	case watch.Error:
		k.logger.Info("Error attempting to watch Kubernetes Jobs")

		// This round trip allows us to handle unstructured status
		errObject := apierrors.FromObject(event.Object)
		statusErr, ok := errObject.(*apierrors.StatusError)
		if !ok {
			k.logger.Info(fmt.Sprintf("received an error which is not *metav1.Status but %+v", event.Object))
		}

		status := statusErr.ErrStatus
		fmt.Printf("%v", status)
	default:
	}
	return nil
}

func (k *K8sContainerService) isJobFailed(job *batchv1.Job) bool {
	for _, condition := range job.Status.Conditions {
		if condition.Status == corev1.ConditionTrue && condition.Type == batchv1.JobFailed {
			return true
		}
	}

	return false
}

//
//func (k *K8sContainerService) watchForJob(ctx context.Context, jobID, name string) error {
//	wCtx, cancel := context.WithCancel(ctx)
//	defer cancel()
//
//	rw, err := watchtools.NewRetryWatcher("1", &cache.ListWatch{
//		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
//			return k.client.BatchV1().Jobs(k.namespace).Watch(wCtx, metav1.ListOptions{
//				LabelSelector: "job-id=" + jobID,
//			})
//		},
//	})
//	if err != nil {
//		return fmt.Errorf("error creating label watcher: %s", err.Error())
//	}
//
//	defer rw.Stop()
//
//	//go func() {
//	//	<-ctx.Done()
//	//	// Cancel the context
//	//	rw.Stop()
//	//}()
//
//	//ch := rw.ResultChan()
//	//defer rw.Stop()
//
//	for {
//		select {
//		case event := <-rw.ResultChan():
//			switch event.Type {
//
//			case watch.Added, watch.Modified:
//				job, ok := event.Object.(*batchv1.Job)
//				if !ok {
//					return fmt.Errorf("unable to parse Kubernetes Job from Annotation watcher")
//				}
//
//				done := false
//				message := ""
//				failed := false
//
//				for _, condition := range job.Status.Conditions {
//					switch condition.Type {
//					case batchv1.JobFailed:
//						failed = true
//						message = condition.Message
//
//						done = true
//					case batchv1.JobComplete:
//						failed = false
//						message = condition.Message
//
//						done = true
//					}
//				}
//
//				fmt.Printf(".")
//
//				if done {
//					if failed {
//						fmt.Printf("\nJob %s.%s (%s) failed %s\n", name, k.namespace, jobID, message)
//					} else {
//						fmt.Printf("\nJob %s.%s (%s) succeeded %s\n", name, k.namespace, jobID, message)
//					}
//
//					pods, err := k.client.CoreV1().Pods(k.namespace).List(wCtx, metav1.ListOptions{
//						LabelSelector: fmt.Sprintf("job-name=%s", name),
//					})
//					if err != nil {
//						return fmt.Errorf("getting pods: %w", err)
//					}
//
//					fmt.Printf("Found %d pods\n", len(pods.Items))
//
//					for _, pod := range pods.Items {
//						req := k.client.CoreV1().Pods(k.namespace).GetLogs(pod.Name, &corev1.PodLogOptions{
//							TailLines: pointer.Int64(1),
//						})
//
//						stream, err := req.Stream(ctx)
//						if err != nil {
//							fmt.Println("HERE in req Stream Error")
//							return err
//						}
//						defer stream.Close()
//
//						buf := new(bytes.Buffer)
//
//						_, err = io.Copy(buf, stream)
//						if err != nil {
//							fmt.Println("HERE in Copy Error")
//							return err
//						}
//						fmt.Println("HERE")
//						fmt.Println(buf.String())
//					}
//					//logsOut, err := logs(wCtx, clientset, podNames, namespace)
//					//if err != nil {
//					//	return fmt.Errorf("getting logs: %w", err)
//					//}
//					//logsOut = "Recorded: " + time.Now().UTC().String() + "\n\n" + logsOut
//
//					deletePol := metav1.DeletePropagationBackground
//					if err := k.client.BatchV1().Jobs(k.namespace).Delete(wCtx, name, metav1.DeleteOptions{PropagationPolicy: &deletePol}); err != nil {
//						return fmt.Errorf("deleting job: %w", err)
//					} else {
//						fmt.Printf("Deleted job %s\n", name)
//					}
//
//					cancel()
//
//					return nil
//				}
//
//			case watch.Deleted:
//				_, ok := event.Object.(*batchv1.Job)
//				if !ok {
//					return fmt.Errorf("unable to parse Kubernetes Job from Annotation watcher")
//				}
//
//			case watch.Error:
//				k.logger.Info("Error attempting to watch Kubernetes Jobs")
//
//				// This round trip allows us to handle unstructured status
//				errObject := apierrors.FromObject(event.Object)
//				statusErr, ok := errObject.(*apierrors.StatusError)
//				if !ok {
//					k.logger.Info(spew.Sprintf("received an error which is not *metav1.Status but %#+v", event.Object))
//
//				}
//
//				status := statusErr.ErrStatus
//				fmt.Printf("%v", status)
//			default:
//			}
//		case <-ctx.Done():
//			break
//		}
//	}
//	//
//	//for event := range ch {
//	//	// We need to inspect the event and get ResourceVersion out of it
//	//	switch event.Type {
//	//
//	//	case watch.Added, watch.Modified:
//	//		job, ok := event.Object.(*batchv1.Job)
//	//		if !ok {
//	//			return fmt.Errorf("unable to parse Kubernetes Job from Annotation watcher")
//	//		}
//	//
//	//		done := false
//	//		message := ""
//	//		failed := false
//	//
//	//		for _, condition := range job.Status.Conditions {
//	//			switch condition.Type {
//	//			case batchv1.JobFailed:
//	//				failed = true
//	//				message = condition.Message
//	//
//	//				done = true
//	//			case batchv1.JobComplete:
//	//				failed = false
//	//				message = condition.Message
//	//
//	//				done = true
//	//			}
//	//		}
//	//
//	//		fmt.Printf(".")
//	//
//	//		if done {
//	//			if failed {
//	//				fmt.Printf("\nJob %s.%s (%s) failed %s\n", name, k.namespace, jobID, message)
//	//			} else {
//	//				fmt.Printf("\nJob %s.%s (%s) succeeded %s\n", name, k.namespace, jobID, message)
//	//			}
//	//
//	//			pods, err := k.client.CoreV1().Pods(k.namespace).List(wCtx, metav1.ListOptions{
//	//				LabelSelector: fmt.Sprintf("job-name=%s", name),
//	//			})
//	//			if err != nil {
//	//				return fmt.Errorf("getting pods: %w", err)
//	//			}
//	//
//	//			fmt.Printf("Found %d pods\n", len(pods.Items))
//	//
//	//			for _, pod := range pods.Items {
//	//				req := k.client.CoreV1().Pods(k.namespace).GetLogs(pod.Name, &corev1.PodLogOptions{
//	//					TailLines: pointer.Int64(1),
//	//				})
//	//
//	//				stream, err := req.Stream(ctx)
//	//				if err != nil {
//	//					fmt.Println("HERE in req Stream Error")
//	//					return err
//	//				}
//	//				defer stream.Close()
//	//
//	//				buf := new(bytes.Buffer)
//	//
//	//				_, err = io.Copy(buf, stream)
//	//				if err != nil {
//	//					fmt.Println("HERE in Copy Error")
//	//					return err
//	//				}
//	//				fmt.Println("HERE")
//	//				fmt.Println(buf.String())
//	//			}
//	//			//logsOut, err := logs(wCtx, clientset, podNames, namespace)
//	//			//if err != nil {
//	//			//	return fmt.Errorf("getting logs: %w", err)
//	//			//}
//	//			//logsOut = "Recorded: " + time.Now().UTC().String() + "\n\n" + logsOut
//	//
//	//			deletePol := metav1.DeletePropagationBackground
//	//			if err := k.client.BatchV1().Jobs(k.namespace).Delete(wCtx, name, metav1.DeleteOptions{PropagationPolicy: &deletePol}); err != nil {
//	//				return fmt.Errorf("deleting job: %w", err)
//	//			} else {
//	//				fmt.Printf("Deleted job %s\n", name)
//	//			}
//	//
//	//			cancel()
//	//
//	//			return nil
//	//		}
//	//
//	//	case watch.Deleted:
//	//		_, ok := event.Object.(*batchv1.Job)
//	//		if !ok {
//	//			return fmt.Errorf("unable to parse Kubernetes Job from Annotation watcher")
//	//		}
//	//
//	//	case watch.Error:
//	//		k.logger.Info("Error attempting to watch Kubernetes Jobs")
//	//
//	//		// This round trip allows us to handle unstructured status
//	//		errObject := apierrors.FromObject(event.Object)
//	//		statusErr, ok := errObject.(*apierrors.StatusError)
//	//		if !ok {
//	//			k.logger.Info(spew.Sprintf("received an error which is not *metav1.Status but %#+v", event.Object))
//	//
//	//		}
//	//
//	//		status := statusErr.ErrStatus
//	//		fmt.Printf("%v", status)
//	//	default:
//	//	}
//	//}
//
//	return nil
//
//}
