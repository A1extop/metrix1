package http

import (
	"github.com/A1extop/metrix1/internal/server/compress"
	"github.com/A1extop/metrix1/internal/server/logging"
	psql "github.com/A1extop/metrix1/internal/server/store/postgrestore"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *Handler, repos *psql.Repository) *gin.Engine {
	router := gin.New()
	log := logging.New()

	router.POST("/update/:type/:name/:value", logging.LoggingPost(log), handler.Update)
	router.POST("/update/", compress.DeCompressData(), logging.LoggingPost(log), handler.UpdateJSON)

	router.POST("/value/", compress.CompressData(), logging.LoggingPost(log), handler.GetJSON)

	router.POST("/updates/", compress.DeCompressData(), logging.LoggingPost(log), handler.UpdatePacketMetricsJSON)

	router.GET("/", compress.CompressData(), logging.LoggingGet(log), handler.DerivationMetrics)
	router.GET("/value/:type/:name", logging.LoggingGet(log), handler.DerivationMetric)

	router.GET("/ping", logging.LoggingGet(log), repos.Ping)
	return router
}
