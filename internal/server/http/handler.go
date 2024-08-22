package http

import (
	"encoding/json"
	"net/http"
	"time"

	"fmt"

	js "github.com/A1extop/metrix1/internal/server/json"
	"github.com/A1extop/metrix1/internal/server/storage"
	"github.com/A1extop/metrix1/internal/server/usecase"
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
	c.JSON(http.StatusOK, value)
}

func (h *Handler) UpdateJSON(c *gin.Context) {
	metricsJs, err := js.GetParametersJSON(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	var metricValue string

	switch metricsJs.MType {
	case "gauge":
		if metricsJs.Value == nil {
			c.String(http.StatusBadRequest, "missing value for gauge metric")
			return
		}
		metricValue = fmt.Sprintf("%g", *metricsJs.Value)
	case "counter":
		if metricsJs.Delta == nil {
			c.String(http.StatusBadRequest, "missing value for counter metric")
			return
		}
		metricValue = fmt.Sprintf("%d", *metricsJs.Delta)
	default:
		c.String(http.StatusBadRequest, "invalid metric type")
		return
	}
	err = usecase.UpdateMetric(h.storage, metricsJs.MType, metricValue, metricsJs.ID)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	//
	if metricsJs.MType == "counter" {
		valueInterface, err := h.storage.ServerSendMetric(metricsJs.ID, metricsJs.MType)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		if delta, ok := valueInterface.(int64); ok {
			metricsJs.Delta = &delta
		} else {
			c.String(http.StatusBadRequest, "Invalid type for counter")
			return
		}
	}

	metric, err := json.Marshal(metricsJs)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.Header("Content-Type", "application/json")
	c.Writer.Write(metric)
	currentTime := time.Now().Format(time.RFC1123)
	c.Header("Date", currentTime)
	c.Status(http.StatusOK)
}

func (h *Handler) GetJSON(c *gin.Context) {
	metrics, err := js.GetParametersJSON(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	valueInterface, err := h.storage.ServerSendMetric(metrics.ID, metrics.MType)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	switch metrics.MType {
	case "counter":
		if delta, ok := valueInterface.(int64); ok {
			metrics.Delta = &delta
		} else {
			c.String(http.StatusBadRequest, "Invalid type for counter")
			return
		}
	case "gauge":
		if value, ok := valueInterface.(float64); ok {
			metrics.Value = &value
		} else {
			c.String(http.StatusBadRequest, "Invalid type for gauge")
			return
		}
	}
	metricJs, err := json.Marshal(metrics)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.Header("Content-Type", "application/json")
	c.Writer.Write(metricJs)
	currentTime := time.Now().Format(time.RFC1123)
	c.Header("Date", currentTime)
	c.Status(http.StatusOK)
}
