package validator

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateRequest(c *gin.Context, expectedContentType, metricName string) bool {
	if c.GetHeader("Content-Type") != expectedContentType {
		c.String(http.StatusUnsupportedMediaType, "Invalid Content-Type")
		return false
	}

	if metricName == "" {
		c.String(http.StatusNotFound, "Empty metricName")
		return false
	}
	return true
}
