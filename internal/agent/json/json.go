package json

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func NewMetrics() *Metrics {
	return &Metrics{}
}
func GetParametersJSON(c *gin.Context) (*Metrics, error) {
	metric := NewMetrics()
	if err := json.NewDecoder(c.Request.Body).Decode(metric); err != nil {
		return nil, err
	}
	return metric, nil
}
