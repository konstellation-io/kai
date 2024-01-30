//go:build unit

package network_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/config"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/kube"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	networkingv1 "k8s.io/api/networking/v1"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/fake"
	kubetesting "k8s.io/client-go/testing"
)

const (
	_labelFormat = "product=%s,version=%s"
	_product     = "test-product"
	_version     = "v1.0.0"
	_workflow    = "test-workflow"
	_process     = "test-process"
)

var _fullProcessIdentifier = getFullProcessIdentifier(_product, _version, _workflow, _process)

type networkSuite struct {
	suite.Suite

	namespace string
	logger    logr.Logger
	clientset *fake.Clientset
	service   *kube.K8sContainerService
}

func TestNetworkSuite(t *testing.T) {
	suite.Run(t, new(networkSuite))
}

func (s *networkSuite) SetupSuite() {
	s.namespace = "test"
	viper.Set(config.KubeNamespaceKey, s.namespace)
	viper.Set(config.BaseDomainNameKey, "test")

	s.logger = testr.NewWithOptions(s.T(), testr.Options{Verbosity: -1})
	s.clientset = fake.NewSimpleClientset()

	s.clientset.PrependReactor(
		"patch",
		"ingresses",
		func(action kubetesting.Action) (handled bool, ret runtime.Object, err error) {
			pa := action.(kubetesting.PatchAction)
			if pa.GetPatchType() == types.ApplyPatchType {
				// Apply patches are supposed to upsert, but fake client fails if the object doesn't exist,
				// if an apply patch occurs for a deployment that doesn't yet exist, create it.
				// However, we already hold the fakeclient lock, so we can't use the front door.
				rfunc := kubetesting.ObjectReaction(s.clientset.Tracker())
				_, obj, err := rfunc(
					kubetesting.NewGetAction(pa.GetResource(), pa.GetNamespace(), pa.GetName()),
				)
				if kubeerrors.IsNotFound(err) || obj == nil {
					_, _, _ = rfunc(
						kubetesting.NewCreateAction(
							pa.GetResource(),
							pa.GetNamespace(),
							&networkingv1.Ingress{
								ObjectMeta: metav1.ObjectMeta{
									Name:      pa.GetName(),
									Namespace: pa.GetNamespace(),
								},
							},
						),
					)
				}
				return rfunc(kubetesting.NewPatchAction(
					pa.GetResource(),
					pa.GetNamespace(),
					pa.GetName(),
					types.StrategicMergePatchType,
					pa.GetPatch()))
			}
			return false, nil, nil
		})

	s.service = kube.NewK8sContainerService(s.logger, s.clientset)

	viper.Set(config.KubeNamespaceKey, s.namespace)
	viper.Set(config.BaseDomainNameKey, "test")
}

func (s *networkSuite) TearDownTest() {
	ctx := context.Background()
	services, err := s.clientset.CoreV1().Services(s.namespace).List(ctx, metav1.ListOptions{})
	s.Require().NoError(err)

	for _, svc := range services.Items {
		err = s.clientset.CoreV1().Services(s.namespace).Delete(ctx, svc.Name, metav1.DeleteOptions{})
		s.Require().NoError(err)
	}

	ingresses, err := s.clientset.NetworkingV1().Ingresses(s.namespace).List(ctx, metav1.ListOptions{})
	s.Require().NoError(err)

	for _, ingress := range ingresses.Items {
		err = s.clientset.NetworkingV1().Ingresses(s.namespace).Delete(ctx, ingress.Name, metav1.DeleteOptions{})
		s.Require().NoError(err)
	}
}

func getFullProcessIdentifier(product, version, workflow, process string) string {
	return strings.ReplaceAll(fmt.Sprintf("%s-%s-%s-%s", product, version, workflow, process), ".", "-")
}
