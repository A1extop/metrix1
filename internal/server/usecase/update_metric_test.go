package usecase

import (
	"errors"
	"testing"

	"github.com/A1extop/metrix1/internal/server/domain"
	store "github.com/A1extop/metrix1/internal/server/storage"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetric(t *testing.T) {

	storage := store.NewMemStorage()

	tests := []struct {
		name        string
		metricType  string
		metricValue string
		metricName  string
		expectedErr error
		verifyFunc  func(t *testing.T, storage *store.MemStorage)
	}{
		{
			name:        "Valid Gauge Metric",
			metricType:  "gauge",
			metricValue: "42.5",
			metricName:  "testGauge",
			expectedErr: nil,
			verifyFunc: func(t *testing.T, storage *store.MemStorage) {
				value, exists := storage.Gauges["testGauge"]
				assert.True(t, exists)
				assert.Equal(t, 42.5, value)
			},
		},
		{
			name:        "Valid Counter Metric",
			metricType:  "counter",
			metricValue: "10",
			metricName:  "testCounter",
			expectedErr: nil,
			verifyFunc: func(t *testing.T, storage *store.MemStorage) {
				value, exists := storage.Counters["testCounter"]
				assert.True(t, exists)
				assert.Equal(t, int64(10), value)
			},
		},
		{
			name:        "Invalid Metric Type",
			metricType:  "unknown",
			metricValue: "10",
			metricName:  "testInvalidType",
			expectedErr: domain.ErrInvalidMetricType,
			verifyFunc:  nil,
		},
		{
			name:        "Invalid Gauge Value",
			metricType:  "gauge",
			metricValue: "not-a-number",
			metricName:  "testInvalidGauge",
			expectedErr: domain.ErrInvalidMetricValue,
			verifyFunc:  nil,
		},
		{
			name:        "Invalid Counter Value",
			metricType:  "counter",
			metricValue: "not-an-int",
			metricName:  "testInvalidCounter",
			expectedErr: domain.ErrInvalidMetricValue,
			verifyFunc:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UpdateMetric(storage, tt.metricType, tt.metricValue, tt.metricName)

			if tt.expectedErr != nil {
				assert.True(t, errors.Is(err, tt.expectedErr))
			} else {
				assert.NoError(t, err)
			}

			if tt.verifyFunc != nil {
				tt.verifyFunc(t, storage)
			}
		})
	}
}
