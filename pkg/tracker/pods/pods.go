package pods

import (
	"github.com/sagacious-labs/k8trics/pkg/store"
	corev1 "k8s.io/api/core/v1"
	coreinformer "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/tools/cache"
)

// Tracker is a struct representing a Pod tracker
type Tracker struct {
	store    *store.PodStore
	informer coreinformer.PodInformer
}

// New returns pointer to a Pod Tracker
func New(informer coreinformer.PodInformer, store *store.PodStore) *Tracker {
	return &Tracker{
		store:    store,
		informer: informer,
	}
}

// Start attahes the tracker handlers
func (t *Tracker) Start() {
	t.informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    t.handleAdd,
		DeleteFunc: t.handleDelete,
	})
}

func (t *Tracker) handleAdd(obj interface{}) {
	casted, ok := obj.(*corev1.Pod)
	if ok {
		t.store.Upsert(*casted)
	}
}

func (t *Tracker) handleDelete(obj interface{}) {
	casted, ok := obj.(*corev1.Pod)
	if ok {
		t.store.Delete(casted.GetName(), casted.GetNamespace())
	}
}
