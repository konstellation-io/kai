package network

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (kn KubeNetwork) UnpublishNetwork(ctx context.Context, product, version string) error {
	err := kn.client.NetworkingV1().Ingresses(kn.namespace).DeleteCollection(ctx, metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("product=%s,version=%s", product, version),
	})
	if err != nil {
		return fmt.Errorf("failed to delete ingresses: %w", err)
	}

	return nil
}
