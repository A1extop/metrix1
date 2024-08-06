package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/A1extop/metrix1/internal/storage"
	"github.com/A1extop/metrix1/internal/usecase"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	storage storage.MetricStorage
}

func NewHandler(storage storage.MetricStorage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) Update(c *gin.Context) {
	metricType := c.Param("type")
	metricName := c.Param("name")
	metricValue := c.Param("value")

	err := usecase.UpdateMetric(h.storage, metricType, metricValue, metricName)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	currentTime := time.Now().Format(time.RFC1123)
	c.Header("Date", currentTime)
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Status(http.StatusOK)
}

func (h *Handler) DerivationMetrics(c *gin.Context) {
	h.storage.ServerSendAllMetrics(c)
	c.Status(http.StatusOK)
}
func (h *Handler) DerivationMetric(c *gin.Context) {
	metricType := c.Param("type")
	metricName := c.Param("name")
	value, err := h.storage.ServerSendMetric(metricName, metricType)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	fmt.Printf("type - %s, name - %s, value - %s", metricType, metricName, value) // написано он возвращал текущее значение метрики в текстовом виде, я так и не понял, каким образом
	c.Status(http.StatusOK)
}
