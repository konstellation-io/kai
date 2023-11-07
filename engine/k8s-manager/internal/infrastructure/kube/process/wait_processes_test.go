//go:build unit

package process_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/service"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/config"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/kube"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/kube/process"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/testhelpers"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestWaitProcesses_Timeout(t *testing.T) {
	var (
		logger    = testr.NewWithOptions(t, testr.Options{Verbosity: 1})
		clientset = fake.NewSimpleClientset()
		ctx       = context.Background()

		version = testhelpers.NewVersionBuilder().Build()
	)

	viper.Set(config.KubeNamespaceKey, _namespace)
	viper.Set(config.ProcessTimeoutKey, 10*time.Millisecond)

	svc := kube.NewK8sContainerService(logger, clientset)

	err := svc.CreateProcess(ctx, service.CreateProcessParams{
		ConfigName: "test",
		Product:    version.Product,
		Version:    version.Tag,
		Workflow:   version.Workflows[0].Name,
		Process:    version.Workflows[0].Processes[0],
	})
	require.NoError(t, err)

	err = svc.WaitProcesses(ctx, version)
	require.ErrorIs(t, err, process.ErrTimeoutWaitingProcesses)
}

func TestWaitProcesses_ContextCancelled(t *testing.T) {
	var (
		logger    = testr.NewWithOptions(t, testr.Options{Verbosity: 1})
		clientset = fake.NewSimpleClientset()

		version = testhelpers.NewVersionBuilder().Build()
	)

	ctx, cancel := context.WithCancel(context.Background())

	viper.Set(config.KubeNamespaceKey, _namespace)
	viper.Set(config.ProcessTimeoutKey, 10*time.Millisecond)

	svc := kube.NewK8sContainerService(logger, clientset)

	err := svc.CreateProcess(ctx, service.CreateProcessParams{
		ConfigName: "test",
		Product:    version.Product,
		Version:    version.Tag,
		Workflow:   version.Workflows[0].Name,
		Process:    version.Workflows[0].Processes[0],
	})
	require.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = svc.WaitProcesses(ctx, version)
		require.ErrorIs(t, err, process.ErrUnexpectedContextClosed)
	}()

	cancel()
}

func cleanDeployments(clientset *fake.Clientset) error {
	ctx := context.Background()
	deployments, err := clientset.AppsV1().Deployments(_namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, deployment := range deployments.Items {
		err := clientset.AppsV1().Deployments(_namespace).Delete(ctx, deployment.Name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}
