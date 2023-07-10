package testhelpers

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	fakecorev1 "k8s.io/client-go/kubernetes/typed/core/v1/fake"
	kubetesting "k8s.io/client-go/testing"
)

type MockCallParams struct {
	Action   string
	Resource string
	Obj      runtime.Object
	Err      error
}

func SetMockCall(clientset *fake.Clientset, params MockCallParams) {
	clientset.CoreV1().(*fakecorev1.FakeCoreV1).
		PrependReactor(params.Action, params.Resource, func(action kubetesting.Action) (handled bool, ret runtime.Object, err error) {
			return true, params.Obj, params.Err
		})
}
