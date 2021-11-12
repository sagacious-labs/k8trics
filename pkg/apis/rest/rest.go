package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sagacious-labs/k8trics/pkg/apis/rest/handlers"
	"github.com/sagacious-labs/k8trics/pkg/apis/rest/routes"
)

func Run() {
	router := gin.Default()
	handlers := handlers.New()

	routes.NewRoutes(router, handlers)

	router.Run()
}
