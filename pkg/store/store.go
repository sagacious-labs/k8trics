package store

import (
	"fmt"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

// PodStore is an in memory store for storing pod info
type PodStore struct {
	internal map[string]K8tricsPod

	lock sync.RWMutex
}

// New returns a new instance of pod store - pod store can
// be used to store pod info in memory
func New() *PodStore {
	return &PodStore{
		internal: make(map[string]K8tricsPod),
	}
}

// Upsert takes in pod and adds that pod to the pod store
func (ps *PodStore) Upsert(pod v1.Pod) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	logrus.Println("Found pod with containers: ", (K8tricsPod{pod}).ContainerIDs())
	ps.internal[fmt.Sprintf("%s.%s", pod.GetNamespace(), pod.GetName())] = K8tricsPod{pod}
}

// Get takes in name and namespace of a pod and returns the pod
func (ps *PodStore) Get(name, namespace string) (K8tricsPod, bool) {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	pod, ok := ps.internal[generateKey(namespace, name)]
	return pod, ok
}

// GetByLabels and returns all of the pods which have the label attached
func (ps *PodStore) GetByLabels(labels map[string]string) (pods []K8tricsPod) {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	for _, pod := range ps.internal {
		if compareLabels(labels, pod.GetLabels()) {
			pods = append(pods, pod)
		}
	}

	return
}

// GetByContainerID takes in a container ID and returns the pod that is running the container
func (ps *PodStore) GetByContainerID(containerID string) (*K8tricsPod, bool) {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	for _, pod := range ps.internal {
		for _, id := range pod.ContainerIDs() {
			splitted := strings.Split(id, "://")
			if len(splitted) != 2 {
				continue
			}

			id = splitted[1]

			if id == containerID {
				return &pod, true
			}
		}
	}

	return nil, false
}

// Delete takes in a name and namespace of a pod and deletes the
// entry corresponding to the name and namespace
func (ps *PodStore) Delete(name, namespace string) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	delete(ps.internal, generateKey(namespace, name))
}

// compareLabels takes in a "from" labels and "with" labels
// and check if all the labels in the "from" map are also in the
// "with" labels
func compareLabels(from, with map[string]string) bool {
	for k, v := range from {
		wv, ok := with[k]
		if !ok || v != wv {
			return false
		}
	}

	return true
}

// generateKey takes in name and namespace of a pod
// and returns a key for the store
func generateKey(namespace, name string) string {
	return fmt.Sprintf("%s.%s", namespace, name)
}
