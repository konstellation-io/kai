//go:build unit

package configuration_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/config"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/kube"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/testhelpers"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes/fake"
)

const (
	_testProduct = "test-product"
	_testVersion = "v1.0.0"
)

func TestDeleteConfiguration(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	clientset := fake.NewSimpleClientset()

	viper.Set(config.KubeNamespaceKey, _namespace)

	svc := kube.NewK8sContainerService(logger, clientset)

	ctx := context.Background()
	// NOTE: Fake client delete resources on delete-collection actions
	err := svc.DeleteConfiguration(ctx, _testProduct, _testVersion)
	require.NoError(t, err)
}

func TestDeleteConfiguration_ClientError(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	clientset := fake.NewSimpleClientset()

	viper.Set(config.KubeNamespaceKey, _namespace)

	expectedErr := errors.New("error deleting configmaps")

	testhelpers.SetMockCall(clientset, testhelpers.MockCallParams{
		Action:   "delete-collection",
		Resource: "configmaps",
		Obj:      nil,
		Err:      expectedErr,
	})

	svc := kube.NewK8sContainerService(logger, clientset)

	ctx := context.Background()
	err := svc.DeleteConfiguration(ctx, _testProduct, _testVersion)
	require.ErrorIs(t, err, expectedErr)
}
