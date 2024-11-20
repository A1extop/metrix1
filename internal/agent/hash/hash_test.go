package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignRequestWithSHA256(t *testing.T) {
	metrics := []byte("test metrics data")
	key := "secret-key"

	h := hmac.New(sha256.New, []byte(key))
	h.Write(metrics)
	expectedHash := hex.EncodeToString(h.Sum(nil))

	resultHash, err := SignRequestWithSHA256(metrics, key)

	assert.NoError(t, err, "expected no error from SignRequestWithSHA256")

	assert.Equal(t, expectedHash, resultHash, "hash should match the expected value")
}
