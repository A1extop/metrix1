package json

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetParametersJSON_ValidGauge(t *testing.T) {
	gin.SetMode(gin.TestMode)

	inputJSON := `{"id":"test_metric","type":"gauge","value":42.42}`
	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(inputJSON))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(response)
	c.Request = request

	metric, err := GetParametersJSON(c)

	assert.NoError(t, err, "expected no error")
	assert.NotNil(t, metric, "metric should not be nil")
	assert.Equal(t, "test_metric", metric.ID)
	assert.Equal(t, "gauge", metric.MType)
	assert.NotNil(t, metric.Value)
	assert.Equal(t, 42.42, *metric.Value)
	assert.Nil(t, metric.Delta)
}

func TestGetParametersJSON_ValidCounter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	inputJSON := `{"id":"test_metric","type":"counter","delta":10}`
	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(inputJSON))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(response)
	c.Request = request

	metric, err := GetParametersJSON(c)

	assert.NoError(t, err, "expected no error")
	assert.NotNil(t, metric, "metric should not be nil")
	assert.Equal(t, "test_metric", metric.ID)
	assert.Equal(t, "counter", metric.MType)
	assert.NotNil(t, metric.Delta)
	assert.Equal(t, int64(10), *metric.Delta)
	assert.Nil(t, metric.Value)
}

func TestGetParametersJSON_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	inputJSON := `{"id":"test_metric","type":"gauge","value":"invalid_value"}`
	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(inputJSON))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(response)
	c.Request = request

	metric, err := GetParametersJSON(c)

	assert.Error(t, err, "expected an error due to invalid JSON")
	assert.Nil(t, metric, "metric should be nil when JSON is invalid")
}

func TestGetParametersJSON_EmptyBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
	response := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(response)
	c.Request = request

	metric, err := GetParametersJSON(c)

	assert.Error(t, err, "expected an error due to empty request body")
	assert.Nil(t, metric, "metric should be nil when body is empty")
}
