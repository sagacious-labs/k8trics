package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sagacious-labs/k8trics/pkg/apis/rest/handlers"
	"github.com/sagacious-labs/k8trics/pkg/apis/rest/routes"
	"github.com/sagacious-labs/k8trics/pkg/store"
)

func Run(store *store.PodStore) {
	router := gin.Default()
	handlers := handlers.New(store)

	routes.NewRoutes(router, handlers)

	router.Run()
}
