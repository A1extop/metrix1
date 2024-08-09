package storage

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MetricStorage interface {
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int64)
	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)
	ServerSendMetric(metricName string, metricType string) (interface{}, error)
	ServerSendAllMetrics(c *gin.Context)
}

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

func (m *MemStorage) ServerSendMetric(metricName string, metricType string) (interface{}, error) {
	switch metricType {
	case "gauge":
		if value, ok := m.gauges[metricName]; ok {
			return value, nil
		}
	case "counter":
		if value, ok := m.counters[metricName]; ok {
			return value, nil
		}
	}
	return "", fmt.Errorf("metric not found")
}

func (m *MemStorage) ServerSendAllMetrics(c *gin.Context) {

	metricsTemplate, err := template.ParseFiles("templates/metrics.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template: %v", err)
		return
	}

	err = metricsTemplate.Execute(c.Writer, m)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error executing template: %v", err)
		return
	}
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (m *MemStorage) UpdateGauge(name string, value float64) {
	m.gauges[name] = value
}

func (m *MemStorage) UpdateCounter(name string, value int64) {
	m.counters[name] += value
}

func (m *MemStorage) GetGauge(name string) (float64, bool) {
	value, exists := m.gauges[name]
	return value, exists
}

func (m *MemStorage) GetCounter(name string) (int64, bool) {
	value, exists := m.counters[name]
	return value, exists
}
