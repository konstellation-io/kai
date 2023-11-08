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

func (s *networkSuite) TestPublish() {
	ctx := context.Background()
	err := s.service.CreateNetwork(ctx, service.CreateNetworkParams{
		Product:  product,
		Version:  version,
		Workflow: workflow,
		Process: &domain.Process{
			Name: process,
			Networking: &domain.Networking{
				SourcePort: 8080,
				TargetPort: 8080,
				Protocol:   "GRPC",
			},
		},
	})
	s.Require().NoError(err)

	versionIdentifier := strings.ReplaceAll(fmt.Sprintf("%s-%s", product, version), ".", "-")

	expectedPublishedURLs := map[string]string{
		fullProcessIdentifier: fmt.Sprintf(_processFormat, versionIdentifier, viper.GetString(config.BaseDomainNameKey), workflow, process),
	}

	publishedURLs, err := s.service.PublishNetwork(ctx, service.PublishNetworkParams{
		Product: product,
		Version: version,
	})
	s.Require().NoError(err)
	s.Assert().Equal(expectedPublishedURLs, publishedURLs)

	res, err := s.clientset.NetworkingV1().Ingresses(s.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf(_labelFormat, product, version),
	})
	s.Require().NoError(err)

	ingressesYaml, err := yaml.Marshal(res)
	require.NoError(s.T(), err)

	g := goldie.New(s.T())
	g.Assert(s.T(), "PublishNetwork", ingressesYaml)
}

func (s *networkSuite) TestPublish_WithTLS() {
	viper.Set(config.TriggersTLSEnabledKey, true)
	ctx := context.Background()

	err := s.service.CreateNetwork(ctx, service.CreateNetworkParams{
		Product:  product,
		Version:  version,
		Workflow: workflow,
		Process: &domain.Process{
			Name: process,
			Networking: &domain.Networking{
				SourcePort: 8080,
				TargetPort: 8080,
				Protocol:   "GRPC",
			},
		},
	})
	s.Require().NoError(err)

	versionIdentifier := strings.ReplaceAll(fmt.Sprintf("%s-%s", product, version), ".", "-")

	expectedPublishedURLs := map[string]string{
		fullProcessIdentifier: fmt.Sprintf(_processFormat, versionIdentifier, viper.GetString(config.BaseDomainNameKey), workflow, process),
	}

	publishedURLs, err := s.service.PublishNetwork(ctx, service.PublishNetworkParams{
		Product: product,
		Version: version,
	})
	s.Require().NoError(err)
	s.Assert().Equal(expectedPublishedURLs, publishedURLs)

	res, err := s.clientset.NetworkingV1().Ingresses(s.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf(_labelFormat, product, version),
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
		Product:  product,
		Version:  version,
		Workflow: workflow,
		Process: &domain.Process{
			Name: process,
			Networking: &domain.Networking{
				SourcePort: 8080,
				TargetPort: 8080,
				Protocol:   "GRPC",
			},
		},
	})
	s.Require().NoError(err)

	versionIdentifier := strings.ReplaceAll(fmt.Sprintf("%s-%s", product, version), ".", "-")

	expectedPublishedURLs := map[string]string{
		fullProcessIdentifier: fmt.Sprintf(_processFormat, versionIdentifier, viper.GetString(config.BaseDomainNameKey), workflow, process),
	}

	publishedURLs, err := s.service.PublishNetwork(ctx, service.PublishNetworkParams{
		Product: product,
		Version: version,
	})
	s.Require().NoError(err)
	s.Assert().Equal(expectedPublishedURLs, publishedURLs)

	res, err := s.clientset.NetworkingV1().Ingresses(s.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf(_labelFormat, product, version),
	})
	s.Require().NoError(err)

	ingressesYaml, err := yaml.Marshal(res)
	require.NoError(s.T(), err)

	g := goldie.New(s.T())
	g.Assert(s.T(), "PublishNetwork_WithTLS_WithTLSSecret", ingressesYaml)
}
