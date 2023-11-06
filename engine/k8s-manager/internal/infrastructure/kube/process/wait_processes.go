package process

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/config"
	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	watchtools "k8s.io/client-go/tools/watch"
)

var (
	ErrParsingDeployment       = errors.New("error parsing deployment from event")
	ErrUnexpectedContextClosed = errors.New("error unexpected context close")
	ErrTimeoutWaitingProcesses = errors.New("error timeout waiting processes")
)

func (kp *KubeProcess) WaitProcesses(ctx context.Context, version domain.Version) error {
	return kp.watchForVersionDeployments(ctx, version)
}

func (kp *KubeProcess) watchForVersionDeployments(ctx context.Context, version domain.Version) error {
	wCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	deployments, err := kp.client.AppsV1().Deployments(kp.namespace).
		List(ctx, metav1.ListOptions{LabelSelector: fmt.Sprintf("product=%s,version=%s", version.Product, version.Tag)})
	if err != nil {
		return fmt.Errorf("listing deployments: %w", err)
	}

	readyDeployments := make(map[string]bool)

	for _, deployment := range deployments.Items {
		if deployment.Status.UnavailableReplicas == 0 {
			readyDeployments[deployment.Name] = true
		}
	}

	rw, err := watchtools.NewRetryWatcher("1", &cache.ListWatch{
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return kp.client.AppsV1().Deployments(kp.namespace).Watch(wCtx, metav1.ListOptions{
				LabelSelector: fmt.Sprintf(
					"product=%s,version=%s", version.Product, version.Tag,
				),
			})
		},
	})
	if err != nil {
		return fmt.Errorf("error creating watcher: %w", err)
	}

	defer rw.Stop()

	for {
		select {
		case event := <-rw.ResultChan():
			deployment, ok := event.Object.(*appsv1.Deployment)
			if !ok {
				return ErrParsingDeployment
			}

			if kp.isDeploymentReady(deployment) {
				readyDeployments[deployment.Name] = true
			}

			if len(readyDeployments) == len(deployments.Items) {
				return nil
			}
		case <-rw.Done():
			return nil
		case <-wCtx.Done():
			return ErrUnexpectedContextClosed
		case <-time.After(viper.GetDuration(config.ProcessTimeoutKey)):
			return ErrTimeoutWaitingProcesses
		}
	}
}

func (kp *KubeProcess) isDeploymentReady(deployment *appsv1.Deployment) bool {
	return deployment.Status.UnavailableReplicas == 0
}
