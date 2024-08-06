package http

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *Handler) *gin.Engine {
	router := gin.Default()
	router.POST("/update/:type/:name/:value", handler.Update)
	router.GET("/", handler.DerivationMetrics)
	router.GET("/value/:type/:name", handler.DerivationMetric)
	return router
}
