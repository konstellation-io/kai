package configuration

import (
	"context"

	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/kube/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (kc KubeConfiguration) DeleteConfiguration(ctx context.Context, product, version string) error {
	err := kc.client.CoreV1().ConfigMaps(kc.namespace).DeleteCollection(
		ctx,
		common.GetDeleteOptions(),
		metav1.ListOptions{LabelSelector: common.GetLabelSelector(product, version)})
	if err != nil {
		return err
	}

	return nil
}
