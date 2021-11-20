package handlers

import "github.com/sagacious-labs/k8trics/pkg/store"

type Handlers struct {
	store *store.PodStore
}

func New(store *store.PodStore) *Handlers {
	return &Handlers{
		store: store,
	}
}
