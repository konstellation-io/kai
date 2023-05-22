package kubernetes

import (
	"fmt"

	configuration "github.com/konstellation-io/kai/engine/k8s-manager/config"
	"github.com/konstellation-io/kai/engine/k8s-manager/entity"
	"github.com/konstellation-io/kai/engine/k8s-manager/kubernetes/node"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"

	"github.com/konstellation-io/kai/libs/simplelogger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Watcher struct {
	config    *configuration.Config
	logger    *simplelogger.SimpleLogger
	clientset *kubernetes.Clientset
}

func NewWatcher(config *configuration.Config, logger *simplelogger.SimpleLogger, clientset *kubernetes.Clientset) *Watcher {
	return &Watcher{
		config,
		logger,
		clientset,
	}
}

func (w *Watcher) WatchNodeStatus(productID, versionName string, statusCh chan<- entity.Node) chan struct{} {
	w.logger.Debugf("[WatchNodeStatus] watching %q", versionName)

	labelSelector := fmt.Sprintf("product-id=%s,version-name=%s,type in (node, entrypoint)", productID, versionName)
	resolver := node.NodeStatusResolver{
		Out:        statusCh,
		Logger:     w.logger,
		LastStatus: map[string]entity.NodeStatus{},
	}

	return w.watchResources(labelSelector, cache.ResourceEventHandlerFuncs{
		AddFunc:    resolver.OnAdd,
		UpdateFunc: resolver.OnUpdate,
		DeleteFunc: resolver.OnDelete,
	})
}

func (w *Watcher) watchResources(labelSelector string, handlers cache.ResourceEventHandler) chan struct{} {
	stopCh := make(chan struct{})

	go func() {
		w.logger.Debugf("Starting informer with labelSelector: %s ", labelSelector)

		factory := informers.NewSharedInformerFactoryWithOptions(w.clientset, 0,
			informers.WithNamespace(w.config.Kubernetes.Namespace),
			informers.WithTweakListOptions(func(options *metav1.ListOptions) {
				options.LabelSelector = labelSelector
			}))

		informer := factory.Core().V1().Pods().Informer()
		informer.AddEventHandler(handlers)
		informer.Run(stopCh)
	}()

	return stopCh
}
