//go:build unit

package network_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/service"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/config"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/kube"
	"github.com/sebdah/goldie/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

type publishNetworkSuite struct {
	suite.Suite

	namespace string
	logger    logr.Logger
	clientset *fake.Clientset
	service   *kube.K8sContainerService
}

func TestPublishNetworkSuite(t *testing.T) {
	suite.Run(t, new(publishNetworkSuite))
}

func (s *publishNetworkSuite) SetupSuite() {
	s.namespace = "test"
	viper.Set(config.KubeNamespaceKey, s.namespace)
	viper.Set(config.BaseDomainNameKey, "test")

	s.logger = testr.NewWithOptions(s.T(), testr.Options{Verbosity: -1})
	s.clientset = fake.NewSimpleClientset()
	s.service = kube.NewK8sContainerService(s.logger, s.clientset)

	viper.Set(config.KubeNamespaceKey, s.namespace)
	viper.Set(config.BaseDomainNameKey, "test")

}

func (s *publishNetworkSuite) TestPublish() {
	var (
		ctx      = context.Background()
		product  = "test-product"
		version  = "v1.0.0"
		workflow = "test-workflow"
		process  = "test-process"

		fullProcessIdentifier = strings.ReplaceAll(fmt.Sprintf("%s-%s-%s-%s", product, version, workflow, process), ".", "-")
	)

	err := s.service.CreateNetwork(ctx, service.CreateNetworkParams{
		Product:  product,
		Version:  version,
		Workflow: workflow,
		Process: &domain.Process{
			Name: process,
			Networking: &domain.Networking{
				SourcePort: 8080,
				Protocol:   "TCP",
				TargetPort: 8080,
			},
		},
	})
	s.Require().NoError(err)

	versionIdentifier := strings.ReplaceAll(fmt.Sprintf("%s-%s", product, version), ".", "-")

	expectedPublishedURLs := map[string]string{
		fullProcessIdentifier: fmt.Sprintf("%s.%s/%s-%s", versionIdentifier, viper.GetString(config.BaseDomainNameKey), workflow, process),
	}

	publishedURLs, err := s.service.PublishNetwork(ctx, service.PublishNetworkParams{
		Product: product,
		Version: version,
	})
	s.Require().NoError(err)
	s.Assert().Equal(expectedPublishedURLs, publishedURLs)

	res, err := s.clientset.NetworkingV1().Ingresses(s.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("product=%s,version=%s", product, version),
	})
	s.Require().NoError(err)

	res, err = s.clientset.NetworkingV1().Ingresses(s.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("product=%s,version=%s", product, version),
	})
	s.Require().NoError(err)

	ingressesYaml, err := yaml.Marshal(res)
	require.NoError(s.T(), err)

	g := goldie.New(s.T())
	g.Assert(s.T(), "PublishNetwork", ingressesYaml)
}

func (s *publishNetworkSuite) TestPublish_WithTLS() {
	var (
		ctx      = context.Background()
		product  = "test-product"
		version  = "v1.0.0"
		workflow = "test-workflow"
		process  = "test-process"

		fullProcessIdentifier = strings.ReplaceAll(fmt.Sprintf("%s-%s-%s-%s", product, version, workflow, process), ".", "-")
	)

	viper.Set(config.TLSIsEnabledKey, true)

	err := s.service.CreateNetwork(ctx, service.CreateNetworkParams{
		Product:  product,
		Version:  version,
		Workflow: workflow,
		Process: &domain.Process{
			Name: process,
			Networking: &domain.Networking{
				SourcePort: 8080,
				Protocol:   "TCP",
				TargetPort: 8080,
			},
		},
	})
	s.Require().NoError(err)

	versionIdentifier := strings.ReplaceAll(fmt.Sprintf("%s-%s", product, version), ".", "-")

	expectedPublishedURLs := map[string]string{
		fullProcessIdentifier: fmt.Sprintf("%s.%s/%s-%s", versionIdentifier, viper.GetString(config.BaseDomainNameKey), workflow, process),
	}

	publishedURLs, err := s.service.PublishNetwork(ctx, service.PublishNetworkParams{
		Product: product,
		Version: version,
	})
	s.Require().NoError(err)
	s.Assert().Equal(expectedPublishedURLs, publishedURLs)

	res, err := s.clientset.NetworkingV1().Ingresses(s.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("product=%s,version=%s", product, version),
	})
	s.Require().NoError(err)

	res, err = s.clientset.NetworkingV1().Ingresses(s.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("product=%s,version=%s", product, version),
	})
	s.Require().NoError(err)

	ingressesYaml, err := yaml.Marshal(res)
	require.NoError(s.T(), err)

	g := goldie.New(s.T())
	g.Assert(s.T(), "PublishNetwork_WithTLS", ingressesYaml)
}

func (s *publishNetworkSuite) TestPublish_WithTLS_WithTLSSecret() {
	var (
		ctx      = context.Background()
		product  = "test-product"
		version  = "v1.0.0"
		workflow = "test-workflow"
		process  = "test-process"

		fullProcessIdentifier = strings.ReplaceAll(fmt.Sprintf("%s-%s-%s-%s", product, version, workflow, process), ".", "-")
	)

	viper.Set(config.TLSIsEnabledKey, true)
	viper.Set(config.TLSSecretNameKey, "test-secret")

	err := s.service.CreateNetwork(ctx, service.CreateNetworkParams{
		Product:  product,
		Version:  version,
		Workflow: workflow,
		Process: &domain.Process{
			Name: process,
			Networking: &domain.Networking{
				SourcePort: 8080,
				Protocol:   "TCP",
				TargetPort: 8080,
			},
		},
	})
	s.Require().NoError(err)

	versionIdentifier := strings.ReplaceAll(fmt.Sprintf("%s-%s", product, version), ".", "-")

	expectedPublishedURLs := map[string]string{
		fullProcessIdentifier: fmt.Sprintf("%s.%s/%s-%s", versionIdentifier, viper.GetString(config.BaseDomainNameKey), workflow, process),
	}

	publishedURLs, err := s.service.PublishNetwork(ctx, service.PublishNetworkParams{
		Product: product,
		Version: version,
	})
	s.Require().NoError(err)
	s.Assert().Equal(expectedPublishedURLs, publishedURLs)

	res, err := s.clientset.NetworkingV1().Ingresses(s.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("product=%s,version=%s", product, version),
	})
	s.Require().NoError(err)

	res, err = s.clientset.NetworkingV1().Ingresses(s.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("product=%s,version=%s", product, version),
	})
	s.Require().NoError(err)

	ingressesYaml, err := yaml.Marshal(res)
	require.NoError(s.T(), err)

	g := goldie.New(s.T())
	g.Assert(s.T(), "PublishNetwork_WithTLS_WithTLSSecret", ingressesYaml)
}

func (s *publishNetworkSuite) TearDownTest() {
	ctx := context.Background()
	services, err := s.clientset.CoreV1().Services(s.namespace).List(ctx, metav1.ListOptions{})
	s.Require().NoError(err)

	for _, svc := range services.Items {
		err = s.clientset.CoreV1().Services(s.namespace).Delete(ctx, svc.Name, metav1.DeleteOptions{})
		s.Require().NoError(err)
	}

	ingresses, err := s.clientset.NetworkingV1().Ingresses(s.namespace).List(ctx, metav1.ListOptions{})
	s.Require().NoError(err)

	for _, ingress := range ingresses.Items {
		err = s.clientset.NetworkingV1().Ingresses(s.namespace).Delete(ctx, ingress.Name, metav1.DeleteOptions{})
		s.Require().NoError(err)
	}
}
