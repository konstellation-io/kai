package process

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (kp *KubeProcess) DeleteProcesses(ctx context.Context, product, version string) error {
	labelSelector := kp.getLabelSelector(product, version)

	err := kp.client.AppsV1().Deployments(kp.namespace).DeleteCollection(
		ctx,
		kp.getDeleteOptions(),
		metav1.ListOptions{LabelSelector: labelSelector},
	)
	if err != nil {
		return err
	}

	err = kp.client.CoreV1().Pods(kp.namespace).DeleteCollection(
		ctx,
		kp.getDeleteOptions(),
		metav1.ListOptions{LabelSelector: labelSelector},
	)
	if err != nil {
		return err
	}

	return nil
}

func (kp *KubeProcess) getDeleteOptions() metav1.DeleteOptions {
	deletePolicy := metav1.DeletePropagationForeground
	gracePeriod := int64(0)

	return metav1.DeleteOptions{
		PropagationPolicy:  &deletePolicy,
		GracePeriodSeconds: &gracePeriod,
	}
}

func (kp *KubeProcess) getLabelSelector(product, version string) string {
	return fmt.Sprintf("product=%s,version=%s", product, version)
}
