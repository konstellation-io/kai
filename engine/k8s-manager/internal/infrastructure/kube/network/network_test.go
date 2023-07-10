//go:build unit

package network_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/go-logr/logr/testr"
	"github.com/sebdah/goldie/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/service"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/kube"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/testhelpers"
)

const _namespace = "test"

func TestCreateNetwork(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	clientset := fake.NewSimpleClientset()
	viper.Set("kubernetes.namespace", _namespace)

	svc := kube.NewK8sContainerService(logger, clientset)
	process := testhelpers.NewProcessBuilder().
		WithNetworking(domain.Networking{
			SourcePort: 80,
			TargetPort: 80,
			Protocol:   "TCP",
		}).
		Build()

	params := service.CreateNetworkParams{
		Product:  "test-product",
		Version:  "test-version",
		Workflow: "test-workflow",
		Process:  process,
	}

	ctx := context.Background()
	err := svc.CreateNetwork(ctx, params)
	require.NoError(t, err)

	s, err := clientset.CoreV1().Services(_namespace).List(ctx, v1.ListOptions{})
	require.NoError(t, err)

	serviceYaml, err := yaml.Marshal(s)
	require.NoError(t, err)

	g := goldie.New(t)
	g.Assert(t, "CreateNetwork", serviceYaml)
}

func TestCreateNetwork_ClientError(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	clientset := fake.NewSimpleClientset()
	ctx := context.Background()

	viper.Set("kubernetes.namespace", _namespace)

	expectedErr := errors.New("error creating service")

	testhelpers.SetMockCall(clientset, testhelpers.MockCallParams{
		Action:   "create",
		Resource: "services",
		Err:      expectedErr,
	})

	svc := kube.NewK8sContainerService(logger, clientset)
	process := testhelpers.NewProcessBuilder().
		WithNetworking(domain.Networking{
			SourcePort: 80,
			TargetPort: 80,
			Protocol:   "TCP",
		}).
		Build()

	params := service.CreateNetworkParams{
		Product:  "test-product",
		Version:  "test-version",
		Workflow: "test-workflow",
		Process:  process,
	}

	err := svc.CreateNetwork(ctx, params)
	require.ErrorIs(t, err, expectedErr)
}

func TestDeleteNetwork(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	clientset := fake.NewSimpleClientset()
	ctx := context.Background()

	viper.Set("kubernetes.namespace", _namespace)

	svc := kube.NewK8sContainerService(logger, clientset)

	process := testhelpers.NewProcessBuilder().
		WithNetworking(domain.Networking{
			SourcePort: 80,
			TargetPort: 80,
			Protocol:   "TCP",
		}).
		Build()

	params := service.CreateNetworkParams{
		Product:  "test-product",
		Version:  "test-version",
		Workflow: "test-workflow",
		Process:  process,
	}

	err := svc.CreateNetwork(ctx, params)
	require.NoError(t, err)

	err = svc.DeleteNetwork(ctx, params.Product, params.Version)
	assert.NoError(t, err)

	svcs, err := clientset.CoreV1().Services(_namespace).List(ctx, v1.ListOptions{})
	require.NoError(t, err)
	assert.Empty(t, svcs.Items)
}

func TestDelete_ClientErrorListingServices(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	clientset := fake.NewSimpleClientset()

	viper.Set("kubernetes.namespace", _namespace)

	expectedErr := errors.New("error listing services")
	testhelpers.SetMockCall(clientset, testhelpers.MockCallParams{
		Action:   "list",
		Resource: "services",
		Err:      expectedErr,
	})

	product := faker.UUIDHyphenated()
	version := faker.UUIDHyphenated()

	svc := kube.NewK8sContainerService(logger, clientset)

	ctx := context.Background()
	err := svc.DeleteNetwork(ctx, product, version)
	require.ErrorIs(t, err, expectedErr)
}

func TestDelete_ClientErrorDeletingSomeService(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	clientset := fake.NewSimpleClientset()

	viper.Set("kubernetes.namespace", _namespace)

	expectedErr := errors.New("error deleting service")
	testhelpers.SetMockCall(clientset, testhelpers.MockCallParams{
		Action:   "delete",
		Resource: "services",
		Err:      expectedErr,
	})

	svc := kube.NewK8sContainerService(logger, clientset)

	process := testhelpers.NewProcessBuilder().
		WithNetworking(domain.Networking{
			SourcePort: 80,
			TargetPort: 80,
			Protocol:   "TCP",
		}).
		Build()

	params := service.CreateNetworkParams{
		Product:  "test-product",
		Version:  "test-version",
		Workflow: "test-workflow",
		Process:  process,
	}

	ctx := context.Background()

	err := svc.CreateNetwork(ctx, params)
	require.NoError(t, err)

	err = svc.DeleteNetwork(ctx, params.Product, params.Version)
	assert.ErrorIs(t, err, expectedErr)
}
