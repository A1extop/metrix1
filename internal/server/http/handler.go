package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/A1extop/metrix1/internal/server/domain"
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

// UpdatePacketMetricsJSON processes an HTTP request to update metrics
// accepted in JSON format. The function gets an array of metrics from the request body,
// validates each metric and updates their values in the repository.
func (h *Handler) UpdatePacketMetricsJSON(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	metrics, err := js.GetParametersMassiveJSON(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	for _, metricsJs := range metrics {
		err = domain.Validate(&metricsJs, c)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			log.Printf("error UpdateJSON: %v", err)
			return
		}
		metricValue := GetValue(&metricsJs)
		err = usecase.UpdateMetric(h.storage, metricsJs.MType, metricValue, metricsJs.ID) //
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		if domain.MetricType(metricsJs.MType) == domain.Counter {
			valueInterface, err := h.storage.ServerFindMetric(metricsJs.ID, metricsJs.MType)
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
		_, err = c.Writer.Write(metric)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
	}

	currentTime := time.Now().Format(time.RFC1123)
	c.Header("Date", currentTime)
	c.Status(http.StatusOK)
}

// Update processes an HTTP request to update the metric according to the specified parameters.
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

// DerivationMetrics output metrics outputs all metrics in
// HTML format.
func (h *Handler) DerivationMetrics(c *gin.Context) {
	h.storage.ServerSendAllMetricsHTML(c)
	c.Status(http.StatusOK)
}

// DerivationMetric processes the request to get a specific metric.
// It takes the metric type and the metric name from the query parameters,
// and then extracts the metric from storage and returns it in JSON format.
func (h *Handler) DerivationMetric(c *gin.Context) {
	metricType := c.Param("type")
	metricName := c.Param("name")
	value, err := h.storage.ServerFindMetric(metricName, metricType)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, value)
}

func GetValue(metricsJs *js.Metrics) string {
	var metricValue string
	switch domain.MetricType(metricsJs.MType) {

	case domain.Gauge:
		metricValue = fmt.Sprintf("%g", *metricsJs.Value)
	case domain.Counter:
		metricValue = fmt.Sprintf("%d", *metricsJs.Delta)
	default:
		return ""
	}
	return metricValue
}

// UpdateJSON handles updating
// metrics in JSON format.
func (h *Handler) UpdateJSON(c *gin.Context) {
	metricsJs, err := js.GetParametersJSON(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	err = domain.Validate(metricsJs, c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		log.Printf("error UpdateJSON: %v", err)
		return
	}
	metricValue := GetValue(metricsJs)
	err = usecase.UpdateMetric(h.storage, metricsJs.MType, metricValue, metricsJs.ID) //
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if domain.MetricType(metricsJs.MType) == domain.Counter {
		var valueInterface interface{}
		valueInterface, err = h.storage.ServerFindMetric(metricsJs.ID, metricsJs.MType)
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
	_, err = c.Writer.Write(metric)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	currentTime := time.Now().Format(time.RFC1123)
	c.Header("Date", currentTime)
	c.Status(http.StatusOK)
}

// GetJSON processes a request to get
// metrics in JSON format.
func (h *Handler) GetJSON(c *gin.Context) {
	metrics, err := js.GetParametersJSON(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	valueInterface, err := h.storage.ServerFindMetric(metrics.ID, metrics.MType)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	switch domain.MetricType(metrics.MType) {
	case domain.Counter:
		if delta, ok := valueInterface.(int64); ok {
			metrics.Delta = &delta
		} else {
			c.String(http.StatusBadRequest, "Invalid type for counter")
			return
		}
	case domain.Gauge:
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
	_, err = c.Writer.Write(metricJs)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	currentTime := time.Now().Format(time.RFC1123)
	c.Header("Date", currentTime)
	c.Status(http.StatusOK)
}
