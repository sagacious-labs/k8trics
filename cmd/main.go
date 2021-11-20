package main

import (
	"github.com/sagacious-labs/k8trics/pkg/apis/rest"
	"github.com/sagacious-labs/k8trics/pkg/k8s"
	"github.com/sagacious-labs/k8trics/pkg/store"
	"github.com/sagacious-labs/k8trics/pkg/tracker"
	"github.com/sagacious-labs/k8trics/pkg/utils"
)

func main() {
	store := store.New()

	khandler, err := k8s.New(utils.GetEnv("KUBECONFIG", ""))
	if err != nil {
		panic(err)
	}

	tracker.
		New(khandler, store).
		Start()

	rest.Run(store)
}
