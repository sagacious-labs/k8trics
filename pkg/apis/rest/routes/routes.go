package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sagacious-labs/k8trics/pkg/apis/rest/handlers"
)

func NewRoutes(r *gin.Engine, handlers *handlers.Handlers) {
	// Setup v1 api routes
	v1ApiRoutes(r, handlers)
}

func v1ApiRoutes(r *gin.Engine, handlers *handlers.Handlers) {
	v1 := r.Group("/api/v1")

	v1.GET("/module", handlers.List)
	v1.GET("/module/:name", handlers.Get)
	v1.GET("/module/:name/log", handlers.WatchLog)
	v1.GET("/module/:name/data", handlers.WatchData)
	v1.DELETE("/module/:name", handlers.Delete)
	v1.POST("/module/:name", handlers.Apply)
}
