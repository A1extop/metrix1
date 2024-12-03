package hash

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestWorkingWithHash(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	key := "secretkey"

	r.POST("/test", WorkingWithHash(key), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	t.Run("Valid Hash", func(t *testing.T) {
		body := []byte("test payload")
		hash := hmac.New(sha256.New, []byte(key))
		hash.Write(body)
		expectedHash := hex.EncodeToString(hash.Sum(nil))

		req, _ := http.NewRequest("POST", "/test", bytes.NewReader(body))
		req.Header.Set("HashSHA256", expectedHash)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status OK; got %v", w.Code)
		}
	})

	t.Run("Invalid Hash Format", func(t *testing.T) {
		body := []byte("test payload")
		req, _ := http.NewRequest("POST", "/test", bytes.NewReader(body))
		req.Header.Set("HashSHA256", "invalidhashformat")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status BadRequest; got %v", w.Code)
		}
		if w.Body.String() != `{"error":"Invalid hash format"}` {
			t.Errorf("Expected error message 'Invalid hash format'; got %v", w.Body.String())
		}
	})

	t.Run("Invalid Hash", func(t *testing.T) {
		body := []byte("test payload")
		hash := hmac.New(sha256.New, []byte("wrongkey"))
		hash.Write(body)
		expectedHash := hex.EncodeToString(hash.Sum(nil))

		req, _ := http.NewRequest("POST", "/test", bytes.NewReader(body))
		req.Header.Set("HashSHA256", expectedHash)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status BadRequest; got %v", w.Code)
		}
		if w.Body.String() != `{"error":"Invalid hash"}` {
			t.Errorf("Expected error message 'Invalid hash'; got %v", w.Body.String())
		}
	})
}
