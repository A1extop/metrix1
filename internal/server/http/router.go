package http

import (
	"github.com/A1extop/metrix1/internal/server/logging"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *Handler) *gin.Engine {
	router := gin.New()
	log := logging.New()

	router.POST("/update/:type/:name/:value", logging.LoggingPost(log), handler.Update)
	router.POST("/update/", logging.LoggingPost(log), handler.UpdateJSON)

	router.POST("/value/", logging.LoggingPost(log), handler.GetJSON)

	router.GET("/", logging.LoggingGet(log), handler.DerivationMetrics)
	router.GET("/value/:type/:name", logging.LoggingGet(log), handler.DerivationMetric)
	return router
}
