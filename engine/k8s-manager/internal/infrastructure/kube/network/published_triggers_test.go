//go:build unit

package network_test

import (
	"context"

	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/service"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/testhelpers"
)

func (s *networkSuite) TestGetPublishedTriggers_OK() {
	ctx := context.Background()
	httpProcess1 := testhelpers.NewProcessBuilder().
		WithName("http-process-1").
		WithNetworking(domain.Networking{
			SourcePort: 8080,
			Protocol:   domain.NetworkingProtocolHTTP,
			TargetPort: 8080,
		}).
		Build()
	httpProcess2 := testhelpers.NewProcessBuilder().
		WithName("http-process-2").
		WithNetworking(domain.Networking{
			SourcePort: 8080,
			Protocol:   domain.NetworkingProtocolHTTP,
			TargetPort: 8080,
		}).
		Build()
	grpcProcess := testhelpers.NewProcessBuilder().
		WithName("grpc-process").
		WithNetworking(domain.Networking{
			SourcePort: 8080,
			Protocol:   domain.NetworkingProtocolGRPC,
			TargetPort: 8080,
		}).
		Build()

	err := s.service.CreateNetwork(ctx, service.CreateNetworkParams{
		Product:  _product,
		Version:  _version,
		Workflow: _workflow,
		Process:  httpProcess1,
	})
	s.Require().NoError(err)

	err = s.service.CreateNetwork(ctx, service.CreateNetworkParams{
		Product:  _product,
		Version:  _version,
		Workflow: _workflow,
		Process:  httpProcess2,
	})
	s.Require().NoError(err)

	err = s.service.CreateNetwork(ctx, service.CreateNetworkParams{
		Product:  _product,
		Version:  _version,
		Workflow: _workflow,
		Process:  grpcProcess,
	})
	s.Require().NoError(err)

	expectedPublishedURLs := map[string]string{
		getFullProcessIdentifier(_product, _version, _workflow, httpProcess1.Name): s.getHTTPEndpoint(_product, _workflow, httpProcess1.Name),
		getFullProcessIdentifier(_product, _version, _workflow, httpProcess2.Name): s.getHTTPEndpoint(_product, _workflow, httpProcess2.Name),
		getFullProcessIdentifier(_product, _version, _workflow, grpcProcess.Name):  s.getGRPCEndpoint(_product, _workflow, grpcProcess.Name),
	}

	_, err = s.service.PublishNetwork(ctx, service.PublishNetworkParams{
		Product: _product,
		Version: _version,
	})
	s.Require().NoError(err)

	publishedURLs, err := s.service.GetPublishedTriggers(ctx, _product)
	s.Require().NoError(err)

	s.Equal(expectedPublishedURLs, publishedURLs)
}

func (s *networkSuite) TestGetPublishedTriggers_NoPublishedIngress() {
	ctx := context.Background()

	publishedURLs, err := s.service.GetPublishedTriggers(ctx, _product)
	s.Require().NoError(err)

	s.Nil(publishedURLs)
}
