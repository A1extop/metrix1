package http

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"os"

	js "github.com/A1extop/metrix1/internal/server/json"
	"testing"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) UpdateGauge(name string, value float64) {
	m.Called(name, value)
}

func (m *MockStorage) UpdateCounter(name string, value int64) {
	m.Called(name, value)
}

func (m *MockStorage) ServerFindMetric(name, mType string) (interface{}, error) {
	args := m.Called(name, mType)
	return args.Get(0), args.Error(1)
}

func (m *MockStorage) GetGauge(name string) (float64, bool) {
	args := m.Called(name)
	return args.Get(0).(float64), args.Bool(1)
}

func (m *MockStorage) GetCounter(name string) (int64, bool) {
	args := m.Called(name)
	return args.Get(0).(int64), args.Bool(1)
}

func (m *MockStorage) ServerSendAllMetricsHTML(c *gin.Context) {
	m.Called(c)
}

func (m *MockStorage) ServerSendAllMetricsToFile(file *os.File) error {
	args := m.Called(file)
	return args.Error(0)
}

func (m *MockStorage) ReadingMetricsFile(file *os.File) error {
	args := m.Called(file)
	return args.Error(0)
}

func (m *MockStorage) RecordingMetricsDB(db *sql.DB) error {
	args := m.Called(db)
	return args.Error(0)
}
func BenchmarkUpdateJSON(b *testing.B) {

	mockStorage := new(MockStorage)
	mockStorage.On("UpdateGauge", mock.Anything, mock.Anything).Return(nil)
	mockStorage.On("UpdateCounter", mock.Anything, mock.Anything).Return(nil)
	mockStorage.On("ServerFindMetric", mock.Anything, mock.Anything).Return(int64(100), nil)

	h := Handler{storage: mockStorage}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.POST("/update", h.UpdateJSON)

	metric := js.Metrics{
		ID:    "test_metric",
		MType: "counter",
		Delta: func(i int64) *int64 { return &i }(5),
	}
	body, _ := json.Marshal(metric)

	req := httptest.NewRequest(http.MethodPost, "/update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		router.ServeHTTP(rec, req)

		rec.Body.Reset()
	}
}

func BenchmarkUpdatePacketMetricsJSON(b *testing.B) {
	// Mock storage setup
	mockStorage := new(MockStorage)
	mockStorage.On("UpdateGauge", mock.Anything, mock.Anything).Return(nil)
	mockStorage.On("UpdateCounter", mock.Anything, mock.Anything).Return(nil)
	mockStorage.On("ServerFindMetric", mock.Anything, mock.Anything).Return(int64(100), nil)

	h := Handler{storage: mockStorage}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.POST("/updates/", h.UpdatePacketMetricsJSON)

	metrics := []js.Metrics{
		{
			ID:    "test_metric_1",
			MType: "counter",
			Delta: func(i int64) *int64 { return &i }(10),
		},
		{
			ID:    "test_metric_2",
			MType: "gauge",
			Value: func(f float64) *float64 { return &f }(3.14),
		},
	}
	body, _ := json.Marshal(metrics)

	req := httptest.NewRequest(http.MethodPost, "/updates/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		router.ServeHTTP(rec, req)

		rec.Body.Reset()
	}
}
func BenchmarkUpdate(b *testing.B) {

	mockStorage := new(MockStorage)
	mockStorage.On("UpdateGauge", mock.Anything, mock.Anything).Return(nil)
	mockStorage.On("UpdateCounter", mock.Anything, mock.Anything).Return(nil)

	h := Handler{storage: mockStorage}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.POST("/update/:type/:name/:value", h.Update)

	testCases := []struct {
		metricType  string
		metricName  string
		metricValue string
	}{
		{"gauge", "temperature", "25.5"},
		{"counter", "requests", "10"},
	}

	b.ResetTimer()
	for _, tc := range testCases {
		b.Run(tc.metricType, func(b *testing.B) {
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/update/%s/%s/%s", tc.metricType, tc.metricName, tc.metricValue), nil)
			rec := httptest.NewRecorder()

			for i := 0; i < b.N; i++ {

				router.ServeHTTP(rec, req)

				rec.Body.Reset()
			}
		})
	}
}
func BenchmarkDerivationMetric(b *testing.B) {

	mockStorage := new(MockStorage)
	mockStorage.On("ServerFindMetric", "temperature", "gauge").Return(25.5, nil)
	mockStorage.On("ServerFindMetric", "requests", "counter").Return(int64(100), nil)

	h := Handler{storage: mockStorage}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.GET("/value/:type/:name", h.DerivationMetric)

	testCases := []struct {
		metricType string
		metricName string
	}{
		{"gauge", "temperature"},
		{"counter", "requests"},
	}

	for _, tc := range testCases {
		b.Run(fmt.Sprintf("%s_%s", tc.metricType, tc.metricName), func(b *testing.B) {
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/value/%s/%s", tc.metricType, tc.metricName), nil)
			rec := httptest.NewRecorder()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				router.ServeHTTP(rec, req)
				rec.Body.Reset()
			}
		})
	}
}
