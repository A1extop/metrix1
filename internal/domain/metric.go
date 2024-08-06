package domain

import "errors"

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
