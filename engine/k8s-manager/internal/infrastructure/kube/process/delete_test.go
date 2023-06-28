//go:build unit

package process_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/kube"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/testhelpers"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes/fake"
)

func TestDeleteProcess(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	clientset := fake.NewSimpleClientset()

	viper.Set("kubernetes._namespace", _namespace)

	product := faker.UUIDHyphenated()
	version := faker.UUIDHyphenated()

	svc := kube.NewK8sContainerService(logger, clientset)

	ctx := context.Background()
	// NOTE: Fake client delete resources on delete-collection actions
	err := svc.DeleteProcesses(ctx, product, version)
	require.NoError(t, err)
}

func TestDeleteProcess_DeleteDeploymentsError(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	clientset := fake.NewSimpleClientset()

	deleteDeploymentsErr := errors.New("error deleting deployments")

	testhelpers.SetMockCall(clientset, testhelpers.MockCallParams{
		Action:   "delete-collection",
		Resource: "deployments",
		Obj:      nil,
		Err:      deleteDeploymentsErr,
	})

	viper.Set("kubernetes._namespace", _namespace)

	product := faker.UUIDHyphenated()
	version := faker.UUIDHyphenated()

	svc := kube.NewK8sContainerService(logger, clientset)

	ctx := context.Background()
	// NOTE: Fake client delete resources on delete-collection actions
	err := svc.DeleteProcesses(ctx, product, version)
	require.ErrorIs(t, err, deleteDeploymentsErr)
}

func TestDeleteProcess_DeletePodsError(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	clientset := fake.NewSimpleClientset()

	viper.Set("kubernetes._namespace", _namespace)

	deletePodsErr := errors.New("error deleting pods")
	testhelpers.SetMockCall(clientset, testhelpers.MockCallParams{
		Action:   "delete-collection",
		Resource: "pods",
		Obj:      nil,
		Err:      deletePodsErr,
	})

	product := faker.UUIDHyphenated()
	version := faker.UUIDHyphenated()

	svc := kube.NewK8sContainerService(logger, clientset)

	ctx := context.Background()
	// NOTE: Fake client delete resources on delete-collection actions
	err := svc.DeleteProcesses(ctx, product, version)
	require.ErrorIs(t, err, deletePodsErr)
}
