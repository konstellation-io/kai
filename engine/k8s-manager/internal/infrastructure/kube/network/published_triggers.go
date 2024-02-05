package network

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (kn KubeNetwork) GetPublishedTriggers(ctx context.Context, product string) (map[string]string, error) {
	ingressName := kn.getIngressName(product)

	ingress, err := kn.client.NetworkingV1().Ingresses(kn.namespace).Get(ctx, ingressName, metav1.GetOptions{})
	if err != nil {
		if kubeerrors.IsNotFound(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("getting published ingress: %w", err)
	}

	publishedTriggers := map[string]string{}

	for _, rule := range ingress.Spec.Rules {
		host := rule.Host

		for _, path := range rule.HTTP.Paths {
			publishedURL, err := url.JoinPath(host, path.Path)
			if err != nil {
				return nil, fmt.Errorf("failed to build path: %w", err)
			}

			publishedTriggers[path.Backend.Service.Name] = strings.Trim(publishedURL, "/")
		}
	}

	return publishedTriggers, nil
}
