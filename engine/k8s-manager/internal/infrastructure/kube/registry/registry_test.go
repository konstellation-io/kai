package registry_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/config"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/kube/registry"
	"github.com/sebdah/goldie/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

const (
	_namespace    = "test"
	_registryHost = "test.local"
	_imageName    = "test-image:v1.0.0"
)

func TestBuildImage_SucceedJob(t *testing.T) {
	var (
		logger    = testr.NewWithOptions(t, testr.Options{Verbosity: -1})
		clientset = fake.NewSimpleClientset()
		ctx       = context.Background()
	)

	expectedImageRef := fmt.Sprintf("%s/%s", _registryHost, _imageName)

	viper.Set(config.KubeNamespaceKey, _namespace)
	viper.Set(config.ImageRegistryURLKey, fmt.Sprintf("http://%s", _registryHost))
	viper.Set(config.ImageBuilderLogLevel, "error")

	imageBuilder := registry.NewKanikoImageBuilder(logger, clientset)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		imageRef, err := imageBuilder.BuildImage(ctx, _imageName, []byte{})
		require.NoError(t, err)
		require.Equal(t, expectedImageRef, imageRef)
		wg.Done()
	}()

	time.Sleep(500 * time.Millisecond)

	// Check if the job is created
	job, err := clientset.BatchV1().Jobs(_namespace).Get(ctx, "image-builder-test-image-v1-0-0", metav1.GetOptions{})
	require.NoError(t, err)

	jobYaml, err := yaml.Marshal(job)
	require.NoError(t, err)

	g := goldie.New(t)
	g.Assert(t, "BuildImage_Job", jobYaml)

	// Check if the configmap is created
	configMap, err := clientset.CoreV1().ConfigMaps(_namespace).Get(ctx, "image-builder-test-image-v1-0-0-config", metav1.GetOptions{})
	require.NoError(t, err)

	configMapYAML, err := yaml.Marshal(configMap)
	require.NoError(t, err)

	g.Assert(t, "BuildImage_ConfigMap", configMapYAML)

	// Update job status to complete
	err = updateJobStatus(ctx, clientset, job, &batchv1.JobCondition{
		Type:   batchv1.JobComplete,
		Status: corev1.ConditionTrue,
	})
	require.NoError(t, err)

	wg.Wait()

	// Check if Job is deleted
	_, err = clientset.BatchV1().Jobs(_namespace).Get(ctx, job.Name, metav1.GetOptions{})
	require.True(t, errors.IsNotFound(err))
}

func TestBuildImage_FailedJob(t *testing.T) {
	var (
		logger    = testr.NewWithOptions(t, testr.Options{Verbosity: -1})
		clientset = fake.NewSimpleClientset()
		ctx       = context.Background()
	)

	viper.Set(config.KubeNamespaceKey, _namespace)
	viper.Set(config.ImageRegistryURLKey, fmt.Sprintf("http://%s", _registryHost))
	viper.Set(config.ImageBuilderLogLevel, "error")

	imageBuilder := registry.NewKanikoImageBuilder(logger, clientset)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		_, err := imageBuilder.BuildImage(ctx, _imageName, []byte{})
		require.ErrorIs(t, err, registry.ErrFailedImageBuild)
		wg.Done()
	}()

	time.Sleep(500 * time.Millisecond)

	// Check if the job is created
	job, err := clientset.BatchV1().Jobs(_namespace).Get(ctx, "image-builder-test-image-v1-0-0", metav1.GetOptions{})
	require.NoError(t, err)

	jobYaml, err := yaml.Marshal(job)
	require.NoError(t, err)

	g := goldie.New(t)
	g.Assert(t, "BuildImage_Job", jobYaml)

	// Check if the configmap is created
	configMap, err := clientset.CoreV1().ConfigMaps(_namespace).Get(ctx, "image-builder-test-image-v1-0-0-config", metav1.GetOptions{})
	require.NoError(t, err)

	configMapYAML, err := yaml.Marshal(configMap)
	require.NoError(t, err)

	g.Assert(t, "BuildImage_ConfigMap", configMapYAML)

	// Update job status to failed
	err = updateJobStatus(ctx, clientset, job, &batchv1.JobCondition{
		Type:   batchv1.JobFailed,
		Status: corev1.ConditionTrue,
	})
	require.NoError(t, err)

	wg.Wait()

	// Check if Job is deleted
	_, err = clientset.BatchV1().Jobs(_namespace).Get(ctx, job.Name, metav1.GetOptions{})
	require.True(t, errors.IsNotFound(err))
}

func TestBuildImage_UnknownEvent(t *testing.T) {
	var (
		logger    = testr.NewWithOptions(t, testr.Options{Verbosity: -1})
		clientset = fake.NewSimpleClientset()
		ctx       = context.Background()
	)

	expectedImageRef := fmt.Sprintf("%s/%s", _registryHost, _imageName)

	viper.Set(config.KubeNamespaceKey, _namespace)
	viper.Set(config.ImageRegistryURLKey, fmt.Sprintf("http://%s", _registryHost))
	viper.Set(config.ImageBuilderLogLevel, "error")

	imageBuilder := registry.NewKanikoImageBuilder(logger, clientset)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		imageRef, err := imageBuilder.BuildImage(ctx, _imageName, []byte{})
		require.NoError(t, err)
		require.Equal(t, expectedImageRef, imageRef)
		wg.Done()
	}()

	time.Sleep(500 * time.Millisecond)

	// Check if the job is created
	job, err := clientset.BatchV1().Jobs(_namespace).Get(ctx, "image-builder-test-image-v1-0-0", metav1.GetOptions{})
	require.NoError(t, err)

	jobYaml, err := yaml.Marshal(job)
	require.NoError(t, err)

	g := goldie.New(t)
	g.Assert(t, "BuildImage_Job", jobYaml)

	// Check if the configmap is created
	configMap, err := clientset.CoreV1().ConfigMaps(_namespace).Get(ctx, "image-builder-test-image-v1-0-0-config", metav1.GetOptions{})
	require.NoError(t, err)

	configMapYAML, err := yaml.Marshal(configMap)
	require.NoError(t, err)

	g.Assert(t, "BuildImage_ConfigMap", configMapYAML)

	// Update job status to suspended
	err = updateJobStatus(ctx, clientset, job, &batchv1.JobCondition{
		Type:   batchv1.JobSuspended,
		Status: corev1.ConditionTrue,
	})
	require.NoError(t, err)

	// waitTimeout returns true if timeout is reached
	require.True(t, waitTimeout(&wg, 1*time.Second))
}

func TestBuildImage_DeletedJob(t *testing.T) {
	var (
		logger    = testr.NewWithOptions(t, testr.Options{Verbosity: -1})
		clientset = fake.NewSimpleClientset()
		ctx       = context.Background()
	)

	viper.Set(config.KubeNamespaceKey, _namespace)
	viper.Set(config.ImageRegistryURLKey, fmt.Sprintf("http://%s", _registryHost))
	viper.Set(config.ImageBuilderLogLevel, "error")

	imageBuilder := registry.NewKanikoImageBuilder(logger, clientset)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		_, err := imageBuilder.BuildImage(ctx, _imageName, []byte{})
		require.ErrorIs(t, err, registry.ErrErrorEvent)
		wg.Done()
	}()

	time.Sleep(500 * time.Millisecond)

	// Check if the job is created
	job, err := clientset.BatchV1().Jobs(_namespace).Get(ctx, "image-builder-test-image-v1-0-0", metav1.GetOptions{})
	require.NoError(t, err)

	jobYaml, err := yaml.Marshal(job)
	require.NoError(t, err)

	g := goldie.New(t)
	g.Assert(t, "BuildImage_Job", jobYaml)

	// Check if the configmap is created
	configMap, err := clientset.CoreV1().ConfigMaps(_namespace).Get(ctx, "image-builder-test-image-v1-0-0-config", metav1.GetOptions{})
	require.NoError(t, err)

	configMapYAML, err := yaml.Marshal(configMap)
	require.NoError(t, err)

	g.Assert(t, "BuildImage_ConfigMap", configMapYAML)

	// Delete job manually
	err = clientset.BatchV1().Jobs(_namespace).Delete(ctx, job.Name, metav1.DeleteOptions{})
	require.NoError(t, err)

	wg.Wait()
}

func updateJobStatus(ctx context.Context, clientset *fake.Clientset, job *batchv1.Job, statusCondition *batchv1.JobCondition) error {
	job.ObjectMeta.ResourceVersion = "batch/v1"
	job.Status.Conditions = []batchv1.JobCondition{
		*statusCondition,
	}

	// This event is ignored
	_, err := clientset.BatchV1().Jobs(_namespace).UpdateStatus(ctx, job, metav1.UpdateOptions{})

	return err
}

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false
	case <-time.After(timeout):
		return true
	}
}
