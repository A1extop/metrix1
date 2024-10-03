package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json"
	"log"

	"github.com/A1extop/metrix1/internal/server/compress"
	"github.com/A1extop/metrix1/internal/server/hash"
	http2 "github.com/A1extop/metrix1/internal/server/http"
	js "github.com/A1extop/metrix1/internal/server/json"
	"github.com/A1extop/metrix1/internal/server/logging"
	"github.com/A1extop/metrix1/internal/server/storage"
	psql "github.com/A1extop/metrix1/internal/server/store/postgrestore"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetric(t *testing.T) {
	gin.SetMode(gin.TestMode)

	memStorage := storage.NewMemStorage()
	handler := http2.NewHandler(memStorage)
	storageDB, err := psql.ConnectDB("")
	if err != nil {
		log.Printf("error connecting to database - %v", err)
	}
	store := psql.NewStore(storageDB)
	repos := psql.NewRepository(store)
	router := http2.NewRouter(handler, repos, "key")

	tests := []struct {
		metricType  string
		metricName  string
		metricValue string
		statusCode  int
	}{
		{"gauge", "testGauge", "100.45", http.StatusOK},
		{"counter", "testCounter", "100", http.StatusOK},
		{"invalidType", "testInvalid", "100", http.StatusBadRequest},
		{"gauge", "testGaugeInvalid", "yaAlex", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.metricName, func(t *testing.T) {
			url := "/update/" + tt.metricType + "/" + tt.metricName + "/" + tt.metricValue
			req, err := http.NewRequest("POST", url, nil)
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}
			req.Header.Set("Content-Type", "text/plain")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

func TestGetMetric(t *testing.T) {
	gin.SetMode(gin.TestMode)

	memStorage := storage.NewMemStorage()
	handler := http2.NewHandler(memStorage)

	storageDB, err := psql.ConnectDB("")
	if err != nil {
		log.Printf("error connecting to database - %v", err)
	}
	store := psql.NewStore(storageDB)
	repos := psql.NewRepository(store)
	router := http2.NewRouter(handler, repos, "key")

	memStorage.UpdateGauge("testGauge", 100.3434)

	tests := []struct {
		metricType string
		metricName string
		expected   string
		statusCode int
	}{
		{"gauge", "testGauge", "100.3434", http.StatusOK},
		{"counter", "testCounter", "0", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.metricName, func(t *testing.T) {
			url := "/value/" + tt.metricType + "/" + tt.metricName
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

func TestUpdateJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := logging.New()
	memStorage := storage.NewMemStorage()
	handler := http2.NewHandler(memStorage)

	router := gin.New()
	router.POST("/update/", hash.WorkingWithHash("key"), compress.DeCompressData(), logging.LoggingPost(log), handler.UpdateJSON)

	tests := []struct {
		name       string
		input      js.Metrics
		statusCode int
		expected   string
	}{
		{
			name: "Valid gauge metric",
			input: js.Metrics{
				ID:    "testGauge",
				MType: "gauge",
				Value: float64Pointer(100.3434),
			},
			statusCode: http.StatusOK,
			expected:   "application/json",
		},
		{
			name: "Invalid metric type",
			input: js.Metrics{
				ID:    "testMetric",
				MType: "invalidType",
			},
			statusCode: http.StatusBadRequest,
			expected:   "Invalid metric type",
		},
		{
			name: "Counter metric",
			input: js.Metrics{
				ID:    "testCounter",
				MType: "counter",
				Delta: int64Pointer(5),
			},
			statusCode: http.StatusOK,
			expected:   "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("Could not marshal input: %v", err)
			}

			req, err := http.NewRequest("POST", "/update/", bytes.NewBuffer(body))
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.statusCode, w.Code)

			if w.Code == http.StatusOK {
				assert.Equal(t, tt.expected, w.Header().Get("Content-Type"))
				assert.NotEmpty(t, w.Header().Get("Date"))
			}

			if tt.statusCode == http.StatusOK {
				var response js.Metrics
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("Could not unmarshal response: %v", err)
				}
				assert.Equal(t, tt.input.ID, response.ID)
				assert.Equal(t, tt.input.MType, response.MType)
			}
		})
	}
}

func TestGetMetricJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := logging.New()
	memStorage := storage.NewMemStorage()
	handler := http2.NewHandler(memStorage)

	router := gin.New()
	router.POST("/value/", hash.WorkingWithHash("key"), compress.CompressData(), logging.LoggingPost(log), handler.GetJSON)

	memStorage.UpdateGauge("testGauge", 100.3434)
	memStorage.UpdateCounter("testCounter", 10)

	tests := []struct {
		metricType string
		metricName string
		expected   string
		statusCode int
		body       *js.Metrics
	}{
		{
			"gauge",
			"testGauge",
			`{"id":"testGauge","type":"gauge","value":100.3434}`,
			http.StatusOK,
			&js.Metrics{
				ID:    "testGauge",
				MType: "gauge",
			},
		},
		{
			"counter",
			"testCounter",
			`{"id":"testCounter","type":"counter","delta":10}`,
			http.StatusOK,
			&js.Metrics{
				ID:    "testCounter",
				MType: "counter",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.metricName, func(t *testing.T) {
			body, err := json.Marshal(tt.body)
			if err != nil {
				t.Fatalf("Could not marshal input: %v", err)
			}
			req, err := http.NewRequest("POST", "/value/", bytes.NewBuffer(body))
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)

			if w.Code == http.StatusOK {
				assert.JSONEq(t, tt.expected, w.Body.String())
			}
		})
	}
}

func TestUpdatePacketMetricsJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := logging.New()
	memStorage := storage.NewMemStorage()
	handler := http2.NewHandler(memStorage)

	router := gin.New()

	router.POST("/updates/", hash.WorkingWithHash("key"), compress.DeCompressData(), logging.LoggingPost(log), handler.UpdatePacketMetricsJSON)

	tests := []struct {
		name       string
		metrics    []js.Metrics
		expected   string
		statusCode int
	}{
		{
			"UpdateGauge",
			[]js.Metrics{
				{
					ID:    "testGauge",
					MType: "gauge",
					Value: float64Pointer(104.4),
				},
				{
					ID:    "testGauge1",
					MType: "gauge",
					Value: float64Pointer(105.5),
				},
			},
			`[{"id":"testGauge","type":"gauge","value":104.4},{"id":"testGauge1","type":"gauge","value":105.5}]`,
			http.StatusOK,
		},
		{
			"UpdateCounter",
			[]js.Metrics{
				{
					ID:    "testCounter",
					MType: "counter",
					Delta: int64Pointer(10),
				},
				{
					ID:    "testCounter1",
					MType: "counter",
					Delta: int64Pointer(15),
				},
			},
			`[{"id":"testCounter","type":"counter","delta":10},{"id":"testCounter1","type":"counter","delta":15}]`,
			http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.metrics)
			if err != nil {
				t.Fatalf("Could not marshal input: %v", err)
			}

			req, err := http.NewRequest("POST", "/updates/", bytes.NewBuffer(body))
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)

		})
	}
}

func float64Pointer(v float64) *float64 {
	return &v
}
func int64Pointer(v int64) *int64 {
	return &v
}
