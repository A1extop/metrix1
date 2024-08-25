package storage

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/A1extop/metrix1/internal/server/domain"
	"github.com/gin-gonic/gin"
)

type MetricStorage interface {
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int64)
	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)
	ServerSendMetric(metricName string, metricType string) (interface{}, error)
	ServerSendAllMetricsHTML(c *gin.Context)
	ServerSendAllMetrics(*os.File) error
	RecordingMetricsFile(*os.File) error
}

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

func (m *MemStorage) RecordingMetricsFile(file *os.File) error {
	var loadedGauges map[string]float64
	var loadedCounters map[string]int64

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file info: %v", err)
	}

	fileSize := fileInfo.Size()
	if fileSize == 0 {
		return fmt.Errorf("file is empty, no metrics to restore:%v", err)
	}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&loadedGauges); err != nil {
		return fmt.Errorf("error deserializing gauges: %v", err)
	}
	if err := decoder.Decode(&loadedCounters); err != nil {
		return fmt.Errorf("error deserializing counters: %v", err)
	}
	m.gauges = loadedGauges
	m.counters = loadedCounters
	log.Println("Metrics successfully restored from file")
	return nil
}
func (m *MemStorage) ServerSendMetric(metricName string, metricType string) (interface{}, error) {
	switch domain.MetricType(metricType) {
	case domain.Gauge:
		if value, ok := m.gauges[metricName]; ok {
			return value, nil
		}
	case domain.Counter:
		if value, ok := m.counters[metricName]; ok {
			return value, nil
		}
	}
	return "", fmt.Errorf("metric not found")
}

func (m *MemStorage) ServerSendAllMetricsHTML(c *gin.Context) {
	c.Header("Content-Type", "text/html")
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
func (m *MemStorage) ServerSendAllMetrics(file *os.File) error {
	dataGauges, err := json.MarshalIndent(m.gauges, "", " ")
	if err != nil {
		log.Println("Error serializing data:", err)
		return fmt.Errorf("error serializing data: %v", err)
	}
	dataCounters, err := json.MarshalIndent(m.counters, "", " ")
	if err != nil {
		return fmt.Errorf("error serializing data: %v", err)
	}
	if _, err := file.Write(dataGauges); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	if _, err := file.Write(dataCounters); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	log.Println("Metrics successfully written to file ")
	return nil
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
