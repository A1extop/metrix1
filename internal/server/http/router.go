package http

import (
	"github.com/A1extop/metrix1/config/serverconfig"
	"github.com/A1extop/metrix1/internal/server/compress"
	"github.com/A1extop/metrix1/internal/server/hash"
	"github.com/A1extop/metrix1/internal/server/logging"
	psql "github.com/A1extop/metrix1/internal/server/store/postgrestore"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *Handler, repos *psql.Repository, parameters *serverconfig.Parameters) *gin.Engine {
	router := gin.New()
	log := logging.New()

	router.POST("/update/:type/:name/:value", hash.WorkingWithHash(parameters.Key), logging.LoggingPost(log), handler.Update)
	router.POST("/update/", hash.WorkingWithHash(parameters.Key), compress.DeCompressData(), logging.LoggingPost(log), handler.UpdateJSON)

	router.POST("/value/", hash.WorkingWithHash(parameters.Key), compress.CompressData(), logging.LoggingPost(log), handler.GetJSON)

	router.POST("/updates/", hash.WorkingWithHash(parameters.Key), compress.DeCompressData(), logging.LoggingPost(log), handler.UpdatePacketMetricsJSON)

	router.GET("/", hash.WorkingWithHash(parameters.Key), compress.CompressData(), logging.LoggingGet(log), handler.DerivationMetrics)
	router.GET("/value/:type/:name", hash.WorkingWithHash(parameters.Key), logging.LoggingGet(log), handler.DerivationMetric)
	router.GET("/ping", hash.WorkingWithHash(parameters.Key), logging.LoggingGet(log), repos.Ping)
	return router
}
