package hash

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"log"
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

func decrypt(data []byte, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	if len(data) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)

	return data, nil
}
func WorkingWithDecryption(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if key == "" {
			c.Next()
			return
		}
		encryptedBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to read encrypted body"})
			log.Println(err)
			c.Abort()
			return
		}

		decryptedBody, err := decrypt(encryptedBody, key)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Decryption failed"})
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(decryptedBody))

		c.Next()
	}
}
