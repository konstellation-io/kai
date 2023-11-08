//go:build unit

package network_test

import (
	"context"
	"fmt"

	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/service"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/testhelpers"
	"github.com/stretchr/testify/require"
)

func (s *networkSuite) TestUnpublish() {
	ctx := context.Background()
	err := s.service.CreateNetwork(ctx, service.CreateNetworkParams{
		Product:  _product,
		Version:  _version,
		Workflow: _workflow,
		Process: &domain.Process{
			Name: _process,
			Networking: &domain.Networking{
				SourcePort: 8080,
				TargetPort: 8080,
				Protocol:   "GRPC",
			},
		},
	})
	s.Require().NoError(err)

	_, err = s.service.PublishNetwork(ctx, service.PublishNetworkParams{
		Product: _product,
		Version: _version,
	})

	err = s.service.UnpublishNetwork(ctx, _product, _version)
	require.NoError(s.T(), err)
}

func (s *networkSuite) TestUnpublish_FailedToDeleteIngressError() {
	ctx := context.Background()
	err := s.service.CreateNetwork(ctx, service.CreateNetworkParams{
		Product:  _product,
		Version:  _version,
		Workflow: _workflow,
		Process: &domain.Process{
			Name: _process,
			Networking: &domain.Networking{
				SourcePort: 8080,
				TargetPort: 8080,
				Protocol:   "GRPC",
			},
		},
	})
	s.Require().NoError(err)
	expectedErr := fmt.Errorf("error deleting ingress")

	testhelpers.SetMockCall(s.clientset, testhelpers.MockCallParams{
		Action:   "delete-collection",
		Resource: "ingresses",
		Obj:      nil,
		Err:      expectedErr,
	})

	_, err = s.service.PublishNetwork(ctx, service.PublishNetworkParams{
		Product: _product,
		Version: _version,
	})

	err = s.service.UnpublishNetwork(ctx, _product, _version)
	s.Require().ErrorIs(err, expectedErr)
}
