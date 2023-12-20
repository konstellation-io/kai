//go:build unit

package configuration_test

import (
	"context"
	"errors"
	"testing"

	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/config"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/kube"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/testhelpers"

	"github.com/go-logr/logr/testr"
	"github.com/sebdah/goldie/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

const _namespace = "test"

func TestConfigCreation(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	clientset := fake.NewSimpleClientset()

	viper.Set(config.KubeNamespaceKey, _namespace)
	viper.Set(config.MinioEndpointKey, "test-minio-endpoint")
	viper.Set(config.AuthEndpointKey, "test-auth-endpoint")
	viper.Set(config.AuthRealmKey, "test-auth-realm")
	viper.Set(config.AuthClientIDKey, "test-auth-client-id")
	viper.Set(config.AuthClientSecretKey, "test-auth-client-secret")
	viper.Set(config.PredictionsIndexKey, "predictionsIdx")

	version := testhelpers.NewVersionBuilder().Build()

	svc := kube.NewK8sContainerService(logger, clientset)

	ctx := context.Background()
	configMapName, err := svc.CreateVersionConfiguration(ctx, version)
	assert.NoError(t, err)

	configMap, err := clientset.CoreV1().ConfigMaps(_namespace).Get(ctx, configMapName, v1.GetOptions{})
	assert.NoError(t, err)

	configMapYaml, err := yaml.Marshal(configMap)
	require.NoError(t, err)

	g := goldie.New(t)
	g.Assert(t, "ConfigMapCreation", configMapYaml)
}

func TestConfigCreation_ClientError(t *testing.T) {
	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	viper.Set(config.KubeNamespaceKey, _namespace)

	clientset := fake.NewSimpleClientset()
	clientError := errors.New("error creating confimap")
	testhelpers.SetMockCall(
		clientset,
		testhelpers.MockCallParams{Action: "create", Resource: "configmaps", Obj: nil, Err: clientError},
	)

	ctx := context.Background()

	svc := kube.NewK8sContainerService(logger, clientset)
	version := testhelpers.NewVersionBuilder().Build()

	_, err := svc.CreateVersionConfiguration(ctx, version)
	assert.ErrorIs(t, err, clientError)
}
