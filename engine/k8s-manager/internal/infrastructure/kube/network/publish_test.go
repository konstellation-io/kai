//go:build unit

package network_test

import (
	"context"
	"fmt"
	"strings"

	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/service"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/config"
	"github.com/sebdah/goldie/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	_httpEndpointFormat = "%s.%s/%s-%s"
	_grpcEndpointFormat = "%s-%s-%s.%s"
)

func (s *networkSuite) TestPublish() {
	ctx := context.Background()
	err := s.service.CreateNetwork(ctx, service.CreateNetworkParams{
		Product:  _product,
		Version:  _version,
		Workflow: _workflow,
		Process: &domain.Process{
			Name: _process,
			Networking: &domain.Networking{
				SourcePort: 8080,
				Protocol:   domain.NetworkingProtocolHTTP,
				TargetPort: 8080,
			},
		},
	})
	s.Require().NoError(err)

	expectedPublishedURLs := map[string]string{
		_fullProcessIdentifier: s.getHTTPEndpoint(_product, _workflow, _process),
	}

	publishedURLs, err := s.service.PublishNetwork(ctx, service.PublishNetworkParams{
		Product: _product,
		Version: _version,
	})
	s.Require().NoError(err)
	s.Assert().Equal(expectedPublishedURLs, publishedURLs)

	res, err := s.clientset.NetworkingV1().Ingresses(s.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf(_labelFormat, _product, _version),
	})
	s.Require().NoError(err)

	ingressesYaml, err := yaml.Marshal(res)
	require.NoError(s.T(), err)

	g := goldie.New(s.T())
	g.Assert(s.T(), "PublishNetwork", ingressesYaml)
}

func (s *networkSuite) TestPublish_GRPCTrigger() {
	ctx := context.Background()
	err := s.service.CreateNetwork(ctx, service.CreateNetworkParams{
		Product:  _product,
		Version:  _version,
		Workflow: _workflow,
		Process: &domain.Process{
			Name: _process,
			Networking: &domain.Networking{
				SourcePort: 8080,
				Protocol:   domain.NetworkingProtocolGRPC,
				TargetPort: 8080,
			},
		},
	})
	s.Require().NoError(err)

	expectedPublishedURLs := map[string]string{
		_fullProcessIdentifier: s.getGRPCEndpoint(_product, _workflow, _process),
	}

	publishedURLs, err := s.service.PublishNetwork(ctx, service.PublishNetworkParams{
		Product: _product,
		Version: _version,
	})
	s.Require().NoError(err)
	s.Assert().Equal(expectedPublishedURLs, publishedURLs)

	res, err := s.clientset.NetworkingV1().Ingresses(s.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf(_labelFormat, _product, _version),
	})
	s.Require().NoError(err)

	ingressesYaml, err := yaml.Marshal(res)
	require.NoError(s.T(), err)

	g := goldie.New(s.T())
	g.Assert(s.T(), "PublishNetwork_GRPCTrigger", ingressesYaml)
}

func (s *networkSuite) TestPublish_WithTLS() {
	viper.Set(config.TriggersTLSEnabledKey, true)
	ctx := context.Background()

	err := s.service.CreateNetwork(ctx, service.CreateNetworkParams{
		Product:  _product,
		Version:  _version,
		Workflow: _workflow,
		Process: &domain.Process{
			Name: _process,
			Networking: &domain.Networking{
				SourcePort: 8080,
				Protocol:   domain.NetworkingProtocolHTTP,
				TargetPort: 8080,
			},
		},
	})
	s.Require().NoError(err)

	expectedPublishedURLs := map[string]string{
		_fullProcessIdentifier: s.getHTTPEndpoint(_product, _workflow, _process),
	}

	publishedURLs, err := s.service.PublishNetwork(ctx, service.PublishNetworkParams{
		Product: _product,
		Version: _version,
	})
	s.Require().NoError(err)
	s.Assert().Equal(expectedPublishedURLs, publishedURLs)

	res, err := s.clientset.NetworkingV1().Ingresses(s.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf(_labelFormat, _product, _version),
	})
	s.Require().NoError(err)

	ingressesYaml, err := yaml.Marshal(res)
	require.NoError(s.T(), err)

	g := goldie.New(s.T())
	g.Assert(s.T(), "PublishNetwork_WithTLS", ingressesYaml)
}

func (s *networkSuite) TestPublish_WithTLS_WithTLSSecret() {
	viper.Set(config.TriggersTLSEnabledKey, true)
	viper.Set(config.TLSSecretNameKey, "test-secret")
	ctx := context.Background()

	err := s.service.CreateNetwork(ctx, service.CreateNetworkParams{
		Product:  _product,
		Version:  _version,
		Workflow: _workflow,
		Process: &domain.Process{
			Name: _process,
			Networking: &domain.Networking{
				SourcePort: 8080,
				Protocol:   domain.NetworkingProtocolHTTP,
				TargetPort: 8080,
			},
		},
	})
	s.Require().NoError(err)

	expectedPublishedURLs := map[string]string{
		_fullProcessIdentifier: s.getHTTPEndpoint(_product, _workflow, _process),
	}

	publishedURLs, err := s.service.PublishNetwork(ctx, service.PublishNetworkParams{
		Product: _product,
		Version: _version,
	})
	s.Require().NoError(err)
	s.Assert().Equal(expectedPublishedURLs, publishedURLs)

	res, err := s.clientset.NetworkingV1().Ingresses(s.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf(_labelFormat, _product, _version),
	})
	s.Require().NoError(err)

	ingressesYaml, err := yaml.Marshal(res)
	require.NoError(s.T(), err)

	g := goldie.New(s.T())
	g.Assert(s.T(), "PublishNetwork_WithTLS_WithTLSSecret", ingressesYaml)
}

func (s *networkSuite) getHTTPEndpoint(product, workflow, process string) string {
	parsedProduct := strings.ReplaceAll(product, ".", "-")
	parsedWorkflow := strings.ReplaceAll(workflow, ".", "-")
	parsedProcess := strings.ReplaceAll(process, ".", "-")

	return fmt.Sprintf(_httpEndpointFormat, parsedProduct, viper.GetString(config.BaseDomainNameKey), parsedWorkflow, parsedProcess)
}

func (s *networkSuite) getGRPCEndpoint(product, workflow, process string) string {
	parsedProduct := strings.ReplaceAll(product, ".", "-")
	parsedWorkflow := strings.ReplaceAll(workflow, ".", "-")
	parsedProcess := strings.ReplaceAll(process, ".", "-")

	return fmt.Sprintf(_grpcEndpointFormat, parsedProduct, parsedWorkflow, parsedProcess, viper.GetString(config.BaseDomainNameKey))
}
