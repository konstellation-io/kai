package network

import (
	"context"
	"encoding/base64"
	"fmt"
	"path"
	"strings"

	applymetav1 "k8s.io/client-go/applyconfigurations/meta/v1"
	applynetworkingv1 "k8s.io/client-go/applyconfigurations/networking/v1"

	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/service"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
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
	_fieldManager            = "k8s-manager"
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

	ingressName := kn.getIngressName(params.Product)

	ingressRules, publishedEndpoints := kn.getIngressRules(params.Product, servicesToPublish)

	ingress := &applynetworkingv1.IngressApplyConfiguration{
		TypeMetaApplyConfiguration: applymetav1.TypeMetaApplyConfiguration{
			APIVersion: pointer.String(_apiVersion),
			Kind:       pointer.String(_kindIngress),
		},
		ObjectMetaApplyConfiguration: &applymetav1.ObjectMetaApplyConfiguration{
			Name: &ingressName,
			Labels: map[string]string{
				"product": params.Product,
				"version": params.Version,
				"type":    "network",
			},
			Annotations: annotations,
		},
		Spec: &applynetworkingv1.IngressSpecApplyConfiguration{
			IngressClassName: pointer.String(viper.GetString(config.TriggersIngressClassNameKey)),
			Rules:            ingressRules,
			TLS:              kn.getIngressTLSConfiguration(ingressRules, ingressName),
		},
	}

	_, err = kn.client.NetworkingV1().Ingresses(kn.namespace).Apply(ctx, ingress, metav1.ApplyOptions{
		FieldManager: _fieldManager,
	})
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

func (kn KubeNetwork) getIngressName(product string) string {
	return strings.ReplaceAll(product, ".", "-")
}

func (kn KubeNetwork) getIngressRules(
	product string, servicesToPublish *corev1.ServiceList,
) (rules []applynetworkingv1.IngressRuleApplyConfiguration, endpoints map[string]string) {
	var (
		httpHost = kn.getHTTPHost(product)

		httpPaths          = make([]applynetworkingv1.HTTPIngressPathApplyConfiguration, 0, len(servicesToPublish.Items))
		publishedEndpoints = make(map[string]string, len(servicesToPublish.Items))
		ingressRules       []applynetworkingv1.IngressRuleApplyConfiguration
	)

	for _, svc := range servicesToPublish.Items {
		workflow := svc.Labels["workflow"]
		process := svc.Labels["process"]

		if kn.isGrpc(svc) {
			grpcHost := kn.getGRPCHost(product, workflow, process)
			publishedEndpoints[svc.Name] = grpcHost
			ingressRules = append(ingressRules, kn.getGRPCIngressRule(grpcHost, svc.Name))
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
	return svc.Labels["protocol"] == string(domain.NetworkingProtocolGRPC)
}

func (kn KubeNetwork) getGRPCHost(product, workflow, process string) string {
	return fmt.Sprintf("%s-%s-%s.%s",
		replaceDotsWithHyphen(product),
		replaceDotsWithHyphen(workflow),
		replaceDotsWithHyphen(process),
		viper.GetString(config.BaseDomainNameKey),
	)
}

func (kn KubeNetwork) getGRPCIngressRule(grpcHost, serviceName string) applynetworkingv1.IngressRuleApplyConfiguration {
	pathType := networkingv1.PathTypePrefix

	return applynetworkingv1.IngressRuleApplyConfiguration{
		Host: pointer.String(grpcHost),
		IngressRuleValueApplyConfiguration: applynetworkingv1.IngressRuleValueApplyConfiguration{
			HTTP: &applynetworkingv1.HTTPIngressRuleValueApplyConfiguration{
				Paths: []applynetworkingv1.HTTPIngressPathApplyConfiguration{
					{
						Path:     pointer.String("/"),
						PathType: &pathType,
						Backend: &applynetworkingv1.IngressBackendApplyConfiguration{
							Service: &applynetworkingv1.IngressServiceBackendApplyConfiguration{
								Name: pointer.String(serviceName),
								Port: &applynetworkingv1.ServiceBackendPortApplyConfiguration{
									Name: pointer.String(_servicePortName),
								},
							},
						},
					},
				},
			},
		},
	}
}

func (kn KubeNetwork) getHTTPHost(product string) string {
	return fmt.Sprintf("%s.%s", replaceDotsWithHyphen(product), viper.GetString(config.BaseDomainNameKey))
}

func (kn KubeNetwork) getHTTPIngressRule(httpHost string, httpPaths []applynetworkingv1.HTTPIngressPathApplyConfiguration) applynetworkingv1.IngressRuleApplyConfiguration {
	return applynetworkingv1.IngressRuleApplyConfiguration{
		Host: pointer.String(httpHost),
		IngressRuleValueApplyConfiguration: applynetworkingv1.IngressRuleValueApplyConfiguration{
			HTTP: &applynetworkingv1.HTTPIngressRuleValueApplyConfiguration{
				Paths: httpPaths,
			},
		},
	}
}

func (kn KubeNetwork) getTriggerPath(workflow, process string) string {
	return fmt.Sprintf("/%s-%s", replaceDotsWithHyphen(workflow), replaceDotsWithHyphen(process))
}

func (kn KubeNetwork) getTriggerIngressPath(triggerPath, serviceName string) applynetworkingv1.HTTPIngressPathApplyConfiguration {
	pathType := networkingv1.PathTypePrefix

	return applynetworkingv1.HTTPIngressPathApplyConfiguration{
		Path:     pointer.String(triggerPath),
		PathType: &pathType,
		Backend: &applynetworkingv1.IngressBackendApplyConfiguration{
			Service: &applynetworkingv1.IngressServiceBackendApplyConfiguration{
				Name: pointer.String(serviceName),
				Port: &applynetworkingv1.ServiceBackendPortApplyConfiguration{
					Name: pointer.String(_servicePortName),
				},
			},
		},
	}
}

func (kn KubeNetwork) getIngressTLSConfiguration(
	rules []applynetworkingv1.IngressRuleApplyConfiguration,
	ingressName string,
) []applynetworkingv1.IngressTLSApplyConfiguration {
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
		publishedHosts = append(publishedHosts, *rule.Host)
	}

	return []applynetworkingv1.IngressTLSApplyConfiguration{
		{
			Hosts:      publishedHosts,
			SecretName: pointer.String(tlsSecretName),
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
