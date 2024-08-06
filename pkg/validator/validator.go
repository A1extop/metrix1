package validator

import (
	"net/http"
)

func ValidateRequest(w http.ResponseWriter, req *http.Request, expectedContentType, metricName string) bool {
	if req.Header.Get("Content-Type") != expectedContentType {
		http.Error(w, "Invalid Content-Type", http.StatusUnsupportedMediaType)
		return false
	}

	if metricName == "" {
		http.Error(w, "Empty metricName", http.StatusNotFound)
		return false
	}
	return true
}
