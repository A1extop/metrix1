package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMetric(t *testing.T) {

	t.Run("Valid Metric", func(t *testing.T) {
		metric, err := NewMetric("metric1", Gauge, 3.14)
		assert.NoError(t, err)
		assert.NotNil(t, metric)
		assert.Equal(t, "metric1", metric.Name)
		assert.Equal(t, Gauge, metric.Type)
		assert.Equal(t, 3.14, metric.Value)

		metric, err = NewMetric("metric2", Counter, int64(100))
		assert.NoError(t, err)
		assert.NotNil(t, metric)
		assert.Equal(t, "metric2", metric.Name)
		assert.Equal(t, Counter, metric.Type)
		assert.Equal(t, int64(100), metric.Value)
	})

	t.Run("Invalid Metric Type", func(t *testing.T) {
		metric, err := NewMetric("metric3", "invalid_type", 123)
		assert.Error(t, err)
		assert.Nil(t, metric)
		assert.Equal(t, ErrInvalidMetricType, err)
	})
}

func TestValidateValue(t *testing.T) {
	t.Run("Valid Gauge", func(t *testing.T) {
		metric := &Metric{Name: "gauge1", Type: Gauge, Value: 3.14}
		err := metric.ValidateValue()
		assert.NoError(t, err)
	})

	t.Run("Invalid Gauge Value", func(t *testing.T) {
		metric := &Metric{Name: "gauge2", Type: Gauge, Value: "invalid_value"}
		err := metric.ValidateValue()
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidMetricValue, err)
	})

	t.Run("Valid Counter", func(t *testing.T) {
		metric := &Metric{Name: "counter1", Type: Counter, Value: int64(10)}
		err := metric.ValidateValue()
		assert.NoError(t, err)
	})

	t.Run("Invalid Counter Value", func(t *testing.T) {
		metric := &Metric{Name: "counter2", Type: Counter, Value: "invalid_value"}
		err := metric.ValidateValue()
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidMetricValue, err)
	})

	t.Run("Invalid Metric Type", func(t *testing.T) {
		metric := &Metric{Name: "invalid", Type: "invalid_type", Value: 123}
		err := metric.ValidateValue()
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidMetricType, err)
	})
}
