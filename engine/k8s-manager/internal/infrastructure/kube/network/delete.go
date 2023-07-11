package network

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/kube/common"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (kn KubeNetwork) DeleteNetwork(ctx context.Context, product, version string) error {
	services, err := kn.listServices(ctx, common.GetLabelSelector(product, version))
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	errCh := make(chan error, len(services.Items))

	wg.Add(len(services.Items))

	//nolint:gocritic
	for _, svc := range services.Items {
		go func(svcName string) {
			defer wg.Done()

			err := kn.deleteService(ctx, svcName)
			if err != nil {
				errCh <- fmt.Errorf("error deleting service %q: %w", svcName, err)
			}
		}(svc.Name)
	}

	wg.Wait()
	close(errCh)

	var errs error
	for err := range errCh {
		errs = errors.Join(errs, err)
	}

	return errs
}

func (kn KubeNetwork) listServices(ctx context.Context, labelSelector string) (*corev1.ServiceList, error) {
	return kn.client.CoreV1().Services(kn.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
}

func (kn KubeNetwork) deleteService(ctx context.Context, serviceName string) error {
	kn.logger.V(1).Info("Deleting service", "serviceName", serviceName)

	return kn.client.CoreV1().Services(kn.namespace).Delete(
		ctx,
		serviceName,
		common.GetDeleteOptions(),
	)
}
