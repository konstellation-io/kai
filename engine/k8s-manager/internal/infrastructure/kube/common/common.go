package common

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetLabelSelector(product, version string) string {
	return fmt.Sprintf("product=%s,version=%s", product, version)
}

func GetDeleteOptions() metav1.DeleteOptions {
	deletePolicy := metav1.DeletePropagationForeground
	gracePeriod := int64(0)

	return metav1.DeleteOptions{
		PropagationPolicy:  &deletePolicy,
		GracePeriodSeconds: &gracePeriod,
	}
}
