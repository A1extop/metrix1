package domain

import (
	"errors"

	js "github.com/A1extop/metrix1/internal/server/json"
	"github.com/gin-gonic/gin"
)

type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

var (
	ErrInvalidMetricType  = errors.New("invalid metric type")
	ErrInvalidMetricValue = errors.New("invalid metric value")
)

type Metric struct {
	Name  string
	Type  MetricType
	Value interface{}
}

func NewMetric(name string, metricType MetricType, value interface{}) (*Metric, error) {
	if metricType != Gauge && metricType != Counter {
		return nil, ErrInvalidMetricType
	}
	return &Metric{
		Name:  name,
		Type:  metricType,
		Value: value,
	}, nil
}

func (m *Metric) ValidateValue() error {
	switch m.Type {
	case Gauge:
		if _, ok := m.Value.(float64); !ok {
			return ErrInvalidMetricValue
		}
	case Counter:
		if _, ok := m.Value.(int64); !ok {
			return ErrInvalidMetricValue
		}
	default:
		return ErrInvalidMetricType
	}
	return nil
}
func Validate(metricsJs *js.Metrics, c *gin.Context) error {
	switch MetricType(metricsJs.MType) {

	case Gauge:
		if metricsJs.Value == nil {
			return ErrInvalidMetricValue
		}
	case Counter:
		if metricsJs.Delta == nil {
			return ErrInvalidMetricValue
		}
	default:
		return ErrInvalidMetricType
	}
	return nil
}
