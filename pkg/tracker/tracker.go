package tracker

import (
	"github.com/sagacious-labs/k8trics/pkg/k8s"
	"github.com/sagacious-labs/k8trics/pkg/store"
	"github.com/sagacious-labs/k8trics/pkg/tracker/pods"
)

// Tracker
type Tracker struct {
	pod      *pods.Tracker
	khandler *k8s.K8s

	stop chan struct{}
}

func New(khandler *k8s.K8s, store *store.PodStore) *Tracker {
	pods.New(khandler.Informers().Core().V1().Pods(), store)

	return &Tracker{
		khandler: khandler,
		pod:      pods.New(khandler.Informers().Core().V1().Pods(), store),
		stop:     make(chan struct{}),
	}
}

func (t *Tracker) Start() {
	println("Attaching helpers")
	t.pod.Start()

	println("Staring the informer")
	t.khandler.Informers().Start(t.stop)
}

func (t *Tracker) Stop() {
	t.stop <- struct{}{}
}
