package http

import (
	"github.com/A1extop/metrix1/internal/server/compress"
	"github.com/A1extop/metrix1/internal/server/hash"
	"github.com/A1extop/metrix1/internal/server/logging"
	psql "github.com/A1extop/metrix1/internal/server/store/postgrestore"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *Handler, repos *psql.Repository, key string) *gin.Engine {
	router := gin.New()
	log := logging.New()

	router.POST("/update/:type/:name/:value", hash.WorkingWithHash(key), logging.LoggingPost(log), handler.Update)
	router.POST("/update/", hash.WorkingWithHash(key), compress.DeCompressData(), logging.LoggingPost(log), handler.UpdateJSON)

	router.POST("/value/", hash.WorkingWithHash(key), compress.CompressData(), logging.LoggingPost(log), handler.GetJSON)

	router.POST("/updates/", hash.WorkingWithHash(key), compress.DeCompressData(), logging.LoggingPost(log), handler.UpdatePacketMetricsJSON)

	router.GET("/", hash.WorkingWithHash(key), compress.CompressData(), logging.LoggingGet(log), handler.DerivationMetrics)
	router.GET("/value/:type/:name", hash.WorkingWithHash(key), logging.LoggingGet(log), handler.DerivationMetric)

	router.GET("/ping", hash.WorkingWithHash(key), logging.LoggingGet(log), repos.Ping)
	return router
}
