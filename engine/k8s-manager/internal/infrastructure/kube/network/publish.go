package network

import (
	"context"
	"encoding/base64"
	"fmt"
	"path"
	"strings"

	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/service"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/config"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
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
	servicesToPublish, err := kn.client.CoreV1().Services(kn.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("product=%s,version=%s", params.Product, params.Version),
	})
	if err != nil {
		return nil, fmt.Errorf("listing services: %w", err)
	}

	annotations, err := kn.getIngressAnnotations()
	if err != nil {
		return nil, fmt.Errorf("parsing ingress annotations: %w", err)
	}

	ingressName := kn.getIngressName(params.Product, params.Version)

	ingressRules, publishedEndpoints := kn.getIngressRules(params.Product, params.Version, servicesToPublish)

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
		Spec: networkingv1.IngressSpec{
			IngressClassName: pointer.String(viper.GetString(config.TriggersIngressClassNameKey)),
			Rules:            ingressRules,
			TLS:              kn.getIngressTLSConfiguration(ingressRules, ingressName),
		},
	}

	_, err = kn.client.NetworkingV1().Ingresses(kn.namespace).Create(ctx, ingress, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("creating ingress: %w", err)
	}

	return publishedEndpoints, nil
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

func (kn KubeNetwork) getIngressName(product, version string) string {
	return strings.ReplaceAll(fmt.Sprintf("%s-%s", product, version), ".", "-")
}

func (kn KubeNetwork) getIngressRules(
	product, version string, servicesToPublish *corev1.ServiceList,
) (rules []networkingv1.IngressRule, endpoints map[string]string) {
	var (
		httpHost = kn.getHTTPHost(product, version)

		httpPaths          = make([]networkingv1.HTTPIngressPath, 0, len(servicesToPublish.Items))
		publishedEndpoints = make(map[string]string, len(servicesToPublish.Items))
		ingressRules       []networkingv1.IngressRule
	)

	for _, svc := range servicesToPublish.Items {
		workflow := svc.Labels["workflow"]
		process := svc.Labels["process"]

		if kn.isGrpc(svc) {
			grpcHost := kn.getGRPCHost(workflow, process, httpHost)
			publishedEndpoints[svc.Name] = grpcHost
			ingressRules = append(ingressRules, kn.getGRPCIngressRule(svc.Name, grpcHost))
		} else {
			triggerPath := kn.getTriggerPath(workflow, process)
			publishedEndpoints[svc.Name] = path.Join(httpHost, triggerPath)
			httpPaths = append(httpPaths, kn.getTriggerIngressPath(triggerPath, svc.Name))
		}
	}

	if len(httpPaths) > 0 {
		ingressRules = append(ingressRules, kn.getHTTPIngressRule(httpHost, httpPaths))
	}

	return ingressRules, publishedEndpoints
}

func (kn KubeNetwork) isGrpc(svc corev1.Service) bool {
	return svc.Labels["protocol"] == "grpc"
}

func (kn KubeNetwork) getGRPCHost(workflow, process, httpHost string) string {
	return fmt.Sprintf("%s.%s.%s", replaceDotsWithHyphen(process), replaceDotsWithHyphen(workflow), httpHost)
}

func (kn KubeNetwork) getGRPCIngressRule(grpcHost, serviceName string) networkingv1.IngressRule {
	pathType := networkingv1.PathTypePrefix

	return networkingv1.IngressRule{
		Host: grpcHost,
		IngressRuleValue: networkingv1.IngressRuleValue{
			HTTP: &networkingv1.HTTPIngressRuleValue{
				Paths: []networkingv1.HTTPIngressPath{
					{
						Path:     "/",
						PathType: &pathType,
						Backend: networkingv1.IngressBackend{
							Service: &networkingv1.IngressServiceBackend{
								Name: serviceName,
								Port: networkingv1.ServiceBackendPort{
									Name: _servicePortName,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (kn KubeNetwork) getHTTPHost(product, version string) string {
	return fmt.Sprintf("%s.%s.%s",
		replaceDotsWithHyphen(version), replaceDotsWithHyphen(product), viper.GetString(config.BaseDomainNameKey),
	)
}

func (kn KubeNetwork) getHTTPIngressRule(httpHost string, httpPaths []networkingv1.HTTPIngressPath) networkingv1.IngressRule {
	return networkingv1.IngressRule{
		Host: httpHost,
		IngressRuleValue: networkingv1.IngressRuleValue{
			HTTP: &networkingv1.HTTPIngressRuleValue{
				Paths: httpPaths,
			},
		},
	}
}

func (kn KubeNetwork) getTriggerPath(workflow, process string) string {
	return fmt.Sprintf("%s-%s", replaceDotsWithHyphen(workflow), replaceDotsWithHyphen(process))
}

func (kn KubeNetwork) getTriggerIngressPath(triggerPath, serviceName string) networkingv1.HTTPIngressPath {
	pathType := networkingv1.PathTypePrefix

	return networkingv1.HTTPIngressPath{
		Path:     triggerPath,
		PathType: &pathType,
		Backend: networkingv1.IngressBackend{
			Service: &networkingv1.IngressServiceBackend{
				Name: serviceName,
				Port: networkingv1.ServiceBackendPort{
					Name: _servicePortName,
				},
			},
		},
	}
}

func (kn KubeNetwork) getIngressTLSConfiguration(rules []networkingv1.IngressRule, ingressName string) []networkingv1.IngressTLS {
	if !viper.GetBool(config.TriggersTLSEnabledKey) {
		return nil
	}

	var tlsSecretName string
	if viper.GetString(config.TLSSecretNameKey) != "" {
		tlsSecretName = viper.GetString(config.TLSSecretNameKey)
	} else {
		tlsSecretName = fmt.Sprintf("%s-tls", ingressName)
	}

	publishedHosts := make([]string, 0, len(rules))

	for _, rule := range rules {
		publishedHosts = append(publishedHosts, rule.Host)
	}

	return []networkingv1.IngressTLS{
		{
			Hosts:      publishedHosts,
			SecretName: tlsSecretName,
		},
	}
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

func replaceDotsWithHyphen(str string) string {
	return strings.ReplaceAll(str, ".", "-")
}
