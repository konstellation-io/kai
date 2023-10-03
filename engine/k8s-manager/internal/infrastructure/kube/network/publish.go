package network

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/service"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/config"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

const (
	_apiVersion              = "networking.k8s.io/v1"
	_kindIngress             = "Ingress"
	_kongStripPathAnnotation = "konghq.com/strip-path"
)

func (kn KubeNetwork) PublishNetwork(ctx context.Context, params service.PublishNetworkParams) (map[string]string, error) {
	res, err := kn.client.CoreV1().Services(kn.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("product=%s,version=%s", params.Product, params.Version),
	})
	if err != nil {
		return nil, fmt.Errorf("listing services: %w", err)
	}

	ingressName := strings.ReplaceAll(fmt.Sprintf("%s-%s", params.Product, params.Version), ".", "-")

	triggerHost := fmt.Sprintf("%s.%s", ingressName, viper.GetString(config.BaseDomainNameKey))

	ingressPaths := make([]networkingv1.HTTPIngressPath, 0, len(res.Items))
	publishedEndpoints := make(map[string]string, len(res.Items))

	pathType := networkingv1.PathTypePrefix

	annotations, err := kn.getIngressAnnotations()
	if err != nil {
		return nil, fmt.Errorf("parsing ingress annotations: %w", err)
	}

	for _, svc := range res.Items {
		triggerPath := fmt.Sprintf("/%s-%s", svc.Labels["workflow"], svc.Labels["process"])
		publishedEndpoints[svc.Name] = fmt.Sprintf("%s%s", triggerHost, triggerPath)

		ingressPaths = append(ingressPaths, networkingv1.HTTPIngressPath{
			Path:     triggerPath,
			PathType: &pathType,
			Backend: networkingv1.IngressBackend{
				Service: &networkingv1.IngressServiceBackend{
					Name: svc.Name,
					Port: networkingv1.ServiceBackendPort{
						Name: _servicePortName,
					},
				},
			},
		})
	}

	ingress := &networkingv1.Ingress{
		TypeMeta: metav1.TypeMeta{
			APIVersion: _apiVersion,
			Kind:       _kindIngress,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: ingressName,
			Labels: map[string]string{
				"product": params.Product,
				"version": params.Version,
				"type":    "network",
			},
			Annotations: annotations,
		},
		Spec: kn.getIngressSpec(triggerHost, ingressPaths),
	}

	_, err = kn.client.NetworkingV1().Ingresses(kn.namespace).Create(ctx, ingress, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("creating ingress: %w", err)
	}

	return publishedEndpoints, nil
}

func (kn KubeNetwork) getIngressSpec(triggerHost string, ingressPaths []networkingv1.HTTPIngressPath) networkingv1.IngressSpec {
	ingressSpec := networkingv1.IngressSpec{
		IngressClassName: pointer.String(viper.GetString(config.TriggersIngressClassNameKey)),
		Rules: []networkingv1.IngressRule{
			{
				Host: triggerHost,
				IngressRuleValue: networkingv1.IngressRuleValue{
					HTTP: &networkingv1.HTTPIngressRuleValue{
						Paths: ingressPaths,
					},
				},
			},
		},
	}

	if viper.GetBool(config.TriggersTLSEnabledKey) {
		var tlsSecretName string
		if viper.GetString(config.TLSSecretNameKey) != "" {
			tlsSecretName = viper.GetString(config.TLSSecretNameKey)
		} else {
			tlsSecretName = fmt.Sprintf("%s-tls", triggerHost)
		}

		ingressSpec.TLS = []networkingv1.IngressTLS{
			{
				Hosts:      []string{triggerHost},
				SecretName: tlsSecretName,
			},
		}
	}

	return ingressSpec
}

func (kn KubeNetwork) getIngressAnnotations() (map[string]string, error) {
	annotations, err := base64.StdEncoding.DecodeString(viper.GetString(config.TriggersB64IngressesAnnotaionsKey))
	if err != nil {
		return nil, err
	}

	annotationsMap := make(map[string]string)

	err = yaml.Unmarshal(annotations, &annotationsMap)
	if err != nil {
		return nil, err
	}

	defaultAnnotations := map[string]string{
		_kongStripPathAnnotation: "true",
	}

	if annotationsMap == nil {
		return defaultAnnotations, nil
	}

	return mergeAnnotations(annotationsMap, defaultAnnotations), nil
}

func mergeAnnotations(annotations1, annotations2 map[string]string) map[string]string {
	annotations := make(map[string]string, len(annotations1)+len(annotations2))

	for key, val := range annotations1 {
		annotations[key] = val
	}

	for key, val := range annotations2 {
		annotations[key] = val
	}

	return annotations
}
