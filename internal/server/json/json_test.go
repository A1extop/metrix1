package json

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetParametersJSON_Valid(t *testing.T) {
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

func TestGetParametersJSON_Invalid(t *testing.T) {
	gin.SetMode(gin.TestMode)

	inputJSON := `{"id":"test_metric","type":"gauge","value":"invalid"}`
	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(inputJSON))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(response)
	c.Request = request

	metric, err := GetParametersJSON(c)

	assert.Error(t, err, "expected error for invalid JSON")
	assert.Nil(t, metric, "metric should be nil")
}

func TestGetParametersMassiveJSON_Valid(t *testing.T) {
	gin.SetMode(gin.TestMode)

	inputJSON := `[{"id":"metric1","type":"gauge","value":12.34},{"id":"metric2","type":"counter","delta":10}]`
	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(inputJSON))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(response)
	c.Request = request

	metrics, err := GetParametersMassiveJSON(c)

	assert.NoError(t, err, "expected no error")
	assert.Len(t, metrics, 2, "expected two metrics")

	assert.Equal(t, "metric1", metrics[0].ID)
	assert.Equal(t, "gauge", metrics[0].MType)
	assert.NotNil(t, metrics[0].Value)
	assert.Equal(t, 12.34, *metrics[0].Value)
	assert.Nil(t, metrics[0].Delta)

	assert.Equal(t, "metric2", metrics[1].ID)
	assert.Equal(t, "counter", metrics[1].MType)
	assert.NotNil(t, metrics[1].Delta)
	assert.Equal(t, int64(10), *metrics[1].Delta)
	assert.Nil(t, metrics[1].Value)
}

func TestGetParametersMassiveJSON_Invalid(t *testing.T) {
	gin.SetMode(gin.TestMode)

	inputJSON := `[{"id":"metric1","type":"gauge","value":12.34},{"id":"metric2","type":"counter","delta":"invalid"}]`
	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(inputJSON))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(response)
	c.Request = request

	metrics, err := GetParametersMassiveJSON(c)

	assert.Error(t, err, "expected error for invalid JSON")
	assert.Nil(t, metrics, "metrics should be nil")
}

func TestGetParametersMassiveJSON_Empty(t *testing.T) {
	gin.SetMode(gin.TestMode)

	inputJSON := `[]`
	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(inputJSON))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(response)
	c.Request = request

	metrics, err := GetParametersMassiveJSON(c)

	assert.NoError(t, err, "expected no error for empty array")
	assert.Empty(t, metrics, "expected metrics to be empty")
}
