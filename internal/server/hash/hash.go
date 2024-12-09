package hash

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// WorkingWithHash performs work with hash.
func WorkingWithHash(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		headerValue := c.GetHeader("HashSHA256")
		if headerValue == "" {
			c.Next()
			return
		}

		receivedHash, err := hex.DecodeString(headerValue)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hash format"})
			c.Abort()
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to read request body"})
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		h := hmac.New(sha256.New, []byte(key))
		_, err = h.Write(body)
		if err != nil {
			return
		}
		expectedHash := h.Sum(nil)

		if !hmac.Equal(receivedHash, expectedHash) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hash"})
			c.Abort()
			return
		}

		c.Next()
	}
}
