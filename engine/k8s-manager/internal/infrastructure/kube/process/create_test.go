//go:build unit

package process_test

import (
	"context"
	"errors"
	"testing"

	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/service"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/config"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/kube"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/testhelpers"

	"github.com/go-logr/logr/testr"
	"github.com/sebdah/goldie/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

const _namespace = "test"

func TestStartProcess(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	clientset := fake.NewSimpleClientset()

	viper.Set(config.KubeNamespaceKey, _namespace)

	svc := kube.NewK8sContainerService(logger, clientset)

	ctx := context.Background()

	process := testhelpers.NewProcessBuilder().Build()

	params := service.CreateProcessParams{
		ConfigName: "configmap-name",
		Product:    "test-product",
		Version:    "v1.0.0",
		Workflow:   "test-workflow",
		Process:    process,
	}

	err := svc.CreateProcess(ctx, params)
	require.NoError(t, err)

	deployment, err := clientset.AppsV1().Deployments(_namespace).List(ctx, v1.ListOptions{})
	require.NoError(t, err)

	deploymentYaml, err := yaml.Marshal(deployment)
	require.NoError(t, err)

	g := goldie.New(t)
	g.Assert(t, "StartProcess", deploymentYaml)
}

func TestStartProcess_WithHPA(t *testing.T) {
	var (
		logger    = testr.NewWithOptions(t, testr.Options{Verbosity: -1})
		clientset = fake.NewSimpleClientset()
		svc       = kube.NewK8sContainerService(logger, clientset)
		ctx       = context.Background()
	)

	viper.Set(config.KubeNamespaceKey, _namespace)

	process := testhelpers.NewProcessBuilder().
		WithReplicas(5).
		Build()

	params := service.CreateProcessParams{
		ConfigName: "configmap-name",
		Product:    "test-product",
		Version:    "v1.0.0",
		Workflow:   "test-workflow",
		Process:    process,
	}

	err := svc.CreateProcess(ctx, params)
	require.NoError(t, err)

	deployment, err := clientset.AppsV1().Deployments(_namespace).List(ctx, v1.ListOptions{})
	require.NoError(t, err)

	deploymentYaml, err := yaml.Marshal(deployment)
	require.NoError(t, err)

	autoscaler, err := clientset.AutoscalingV2().HorizontalPodAutoscalers(_namespace).List(ctx, v1.ListOptions{})
	require.NoError(t, err)

	autoscalerYaml, err := yaml.Marshal(autoscaler)
	require.NoError(t, err)

	g := goldie.New(t)
	g.Assert(t, "StartProcess_WithHPA_Deployment", deploymentYaml)
	g.Assert(t, "StartProcess_WithHPA_HPA", autoscalerYaml)
}

func TestStartProcess_WithNetwork(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	clientset := fake.NewSimpleClientset()

	viper.Set(config.KubeNamespaceKey, _namespace)

	svc := kube.NewK8sContainerService(logger, clientset)

	ctx := context.Background()

	process := testhelpers.NewProcessBuilder().
		WithNetworking(domain.Networking{
			SourcePort: 80,
			Protocol:   "tcp",
			TargetPort: 80,
		}).
		Build()

	params := service.CreateProcessParams{
		ConfigName: "configmap-name",
		Product:    "test-product",
		Version:    "v1.0.0",
		Workflow:   "test-workflow",
		Process:    process,
	}

	err := svc.CreateProcess(ctx, params)
	require.NoError(t, err)

	deployment, err := clientset.AppsV1().Deployments(_namespace).List(ctx, v1.ListOptions{})
	require.NoError(t, err)

	deploymentYaml, err := yaml.Marshal(deployment)
	require.NoError(t, err)

	g := goldie.New(t)
	g.Assert(t, "StartProcess_WithNetworking", deploymentYaml)
}

func TestStartProcess_WithResourceLimits(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	clientset := fake.NewSimpleClientset()

	viper.Set(config.KubeNamespaceKey, _namespace)

	svc := kube.NewK8sContainerService(logger, clientset)

	ctx := context.Background()

	process := testhelpers.NewProcessBuilder().
		WithResourceLimits(&domain.ProcessResourceLimits{
			CPU: &domain.ResourceLimit{
				Request: "300m",
				Limit:   "600m",
			},
			Memory: &domain.ResourceLimit{
				Request: "300Mi",
				Limit:   "600Mi",
			},
		}).
		Build()

	params := service.CreateProcessParams{
		ConfigName: "configmap-name",
		Product:    "test-product",
		Version:    "v1.0.0",
		Workflow:   "test-workflow",
		Process:    process,
	}

	err := svc.CreateProcess(ctx, params)
	require.NoError(t, err)

	deployment, err := clientset.AppsV1().Deployments(_namespace).List(ctx, v1.ListOptions{})
	require.NoError(t, err)

	deploymentYaml, err := yaml.Marshal(deployment)
	require.NoError(t, err)

	g := goldie.New(t)
	g.Assert(t, "StartProcess_WithCPU", deploymentYaml)
}

func TestStartProcess_WithMoreReplicas(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	clientset := fake.NewSimpleClientset()
	viper.Set(config.KubeNamespaceKey, _namespace)
	svc := kube.NewK8sContainerService(logger, clientset)

	ctx := context.Background()

	process := testhelpers.NewProcessBuilder().
		WithReplicas(4).
		Build()

	params := service.CreateProcessParams{
		ConfigName: "configmap-name",
		Product:    "test-product",
		Version:    "v1.0.0",
		Workflow:   "test-workflow",
		Process:    process,
	}

	err := svc.CreateProcess(ctx, params)
	require.NoError(t, err)

	deployment, err := clientset.AppsV1().Deployments(_namespace).List(ctx, v1.ListOptions{})
	require.NoError(t, err)

	deploymentYaml, err := yaml.Marshal(deployment)
	require.NoError(t, err)

	g := goldie.New(t)
	g.Assert(t, "StartProcess_WithMoreReplicas", deploymentYaml)
}

func TestStartProcess_WithGpuEnabled(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	clientset := fake.NewSimpleClientset()
	viper.Set(config.KubeNamespaceKey, _namespace)
	svc := kube.NewK8sContainerService(logger, clientset)

	ctx := context.Background()

	process := testhelpers.NewProcessBuilder().
		WithEnableGpu(true).
		Build()

	params := service.CreateProcessParams{
		ConfigName: "configmap-name",
		Product:    "test-product",
		Version:    "v1.0.0",
		Workflow:   "test-workflow",
		Process:    process,
	}

	err := svc.CreateProcess(ctx, params)
	require.NoError(t, err)

	deployment, err := clientset.AppsV1().Deployments(_namespace).List(ctx, v1.ListOptions{})
	require.NoError(t, err)

	deploymentYaml, err := yaml.Marshal(deployment)
	require.NoError(t, err)

	g := goldie.New(t)
	g.Assert(t, "StartProcess_WithEnableGpu", deploymentYaml)
}

func TestStartProcess_ClientError(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	namespace := "test"
	viper.Set(config.KubeNamespaceKey, namespace)
	clientset := fake.NewSimpleClientset()

	expectedError := errors.New("kubernetes client error")
	testhelpers.SetMockCall(clientset, testhelpers.MockCallParams{
		Action:   "create",
		Resource: "deployments",
		Obj:      nil,
		Err:      expectedError,
	})

	process := testhelpers.NewProcessBuilder().Build()

	params := service.CreateProcessParams{
		ConfigName: "configmap-name",
		Product:    "test-product",
		Version:    "v1.0.0",
		Workflow:   "test-workflow",
		Process:    process,
	}

	svc := kube.NewK8sContainerService(logger, clientset)

	ctx := context.Background()

	err := svc.CreateProcess(ctx, params)
	require.Error(t, err)
}
