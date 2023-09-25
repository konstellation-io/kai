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

	_, err = s.service.PublishNetwork(ctx, service.PublishNetworkParams{
		Product: product,
		Version: version,
	})

	err = s.service.UnpublishNetwork(ctx, product, version)
	require.NoError(s.T(), err)
}

func (s *networkSuite) TestUnpublish_FailedToDeleteIngressError() {
	ctx := context.Background()
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
	expectedErr := fmt.Errorf("error deleting ingress")

	testhelpers.SetMockCall(s.clientset, testhelpers.MockCallParams{
		Action:   "delete-collection",
		Resource: "ingresses",
		Obj:      nil,
		Err:      expectedErr,
	})

	_, err = s.service.PublishNetwork(ctx, service.PublishNetworkParams{
		Product: product,
		Version: version,
	})

	err = s.service.UnpublishNetwork(ctx, product, version)
	s.Require().ErrorIs(err, expectedErr)
}
