package storage

import (
	"fmt"
	"html/template"

	"github.com/gin-gonic/gin"
)

type MetricStorage interface {
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int64)
	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)
	ServerSendMetric(metricName string, metricType string) (string, error)
	ServerSendAllMetrics(c *gin.Context)
}

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

func (m *MemStorage) ServerSendMetric(metricName string, metricType string) (string, error) {
	switch metricType {
	case "gauge":
		if value, ok := m.gauges[metricName]; ok {
			return fmt.Sprintf("%f", value), nil
		}
	case "counter":
		if value, ok := m.counters[metricName]; ok {
			return fmt.Sprintf("%d", value), nil
		}
	}
	return "", fmt.Errorf("metric not found")
}

var metricsTemplate = template.Must(template.New("metrics").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Metrics</title>
</head>
<body>
    <h1>Metrics</h1>
    <h2>Gauges</h2>
    <ul>
        {{ range $key, $value := .Gauges }}
            <li>{{ $key }}: {{ $value }}</li>
        {{ end }}
    </ul>
    <h2>Counters</h2>
    <ul>
        {{ range $key, $value := .Counters }}  
            <li>{{ $key }}: {{ $value }}</li>
        {{ end }}
    </ul>
</body>
</html>
`)) // пытался пытался поместить в файл metrics директории template, возникала ошибка

func (m *MemStorage) ServerSendAllMetrics(c *gin.Context) {

	metricsTemplate.Execute(c.Writer, m)
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
