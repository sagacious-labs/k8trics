package store

import (
	"errors"
	"fmt"

	v1 "k8s.io/api/core/v1"
)

// K8tricsPod is a wrapper around v1.Pod struct and adds
// a few helper methods to it
type K8tricsPod struct {
	v1.Pod
}

// Endpoint returns endpoint which can be hit to communicate with the
// given pod without assuming any protocol
//
// The method returns the first endpoint it can find and will return an
// error if no endpoints are found
func (kp K8tricsPod) Endpoint() (string, error) {
	podip := kp.Status.PodIP

	for _, cont := range kp.Spec.Containers {
		for _, port := range cont.Ports {
			return fmt.Sprintf("%s:%d", podip, port.ContainerPort), nil
		}
	}

	return "", errors.New("failed to retrieve endpoint for the pod")
}

// GetContainerIDs return the SHA256 IDs of all of the containers within
// the pod
func (kp K8tricsPod) ContainerIDs() (ids []string) {
	for _, cs := range kp.Status.ContainerStatuses {
		ids = append(ids, cs.ContainerID)
	}

	return
}

// GetPod returns a pointer to the internal Pod struct
func (kp K8tricsPod) GetPod() *v1.Pod {
	return &kp.Pod
}
