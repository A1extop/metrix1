package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"log"

	http2 "github.com/A1extop/metrix1/internal/server/http"
	"github.com/A1extop/metrix1/internal/server/storage"
	psql "github.com/A1extop/metrix1/internal/server/store/postgrestore"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetric(t *testing.T) {
	gin.SetMode(gin.TestMode)

	memStorage := storage.NewMemStorage()
	handler := http2.NewHandler(memStorage)
	storageDB, err := psql.ConnectDB("") //заглушка
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

	storageDB, err := psql.ConnectDB("") //заглушка
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
		expected   string //пока оставил,но так как нет указаний, что надо было делать с отдельными данными, я их вывожу в командной строке
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
