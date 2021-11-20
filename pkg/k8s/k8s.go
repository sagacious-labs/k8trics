package k8s

import (
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type K8s struct {
	clientset *kubernetes.Clientset
	informers informers.SharedInformerFactory

	stop chan struct{}
}

func New(kubeconfigLoc string) (*K8s, error) {
	cs, err := setupClientset(kubeconfigLoc)
	if err != nil {
		return nil, err
	}

	stop := make(chan struct{})

	return &K8s{
		clientset: cs,
		informers: setupInformerFactory(cs, stop),
		stop:      stop,
	}, nil
}

func setupClientset(kubeconfigLoc string) (*kubernetes.Clientset, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfigLoc)
	if err != nil {
		cfg, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	return kubernetes.NewForConfig(cfg)
}

func setupInformerFactory(cs *kubernetes.Clientset, ch chan struct{}) informers.SharedInformerFactory {
	informer := informers.NewSharedInformerFactory(cs, 1*time.Second)

	return informer
}

func (k8s *K8s) ClientSet() *kubernetes.Clientset {
	return k8s.clientset
}

func (k8s *K8s) Informers() informers.SharedInformerFactory {
	return k8s.informers
}

func (k8s *K8s) Close() {
	k8s.stop <- struct{}{}
}
