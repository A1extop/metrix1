package usecase

import (
	"fmt"
	"strconv"

	"github.com/A1extop/metrix1/internal/server/domain"
	"github.com/A1extop/metrix1/internal/server/storage"
)

func UpdateMetric(storage storage.MetricStorage, metricType, metricValue, metricName string) error {
	var value interface{}
	var err error
	switch domain.MetricType(metricType) {
	case domain.Gauge:
		value, err = strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return fmt.Errorf("invalid metric value: %w", domain.ErrInvalidMetricValue)
		}
	case domain.Counter:
		value, err = strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid metric value: %w", domain.ErrInvalidMetricValue)
		}
	default:
		return domain.ErrInvalidMetricType
	}

	metric, err := domain.NewMetric(metricName, domain.MetricType(metricType), value)
	if err != nil {
		return err
	}

	err = metric.ValidateValue()
	if err != nil {
		return err
	}

	switch metric.Type {
	case domain.Gauge:
		storage.UpdateGauge(metric.Name, metric.Value.(float64))
	case domain.Counter:
		storage.UpdateCounter(metric.Name, metric.Value.(int64))
	}

	return nil
}
