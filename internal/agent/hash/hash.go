package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func SignRequestWithSHA256(metrics []byte, key string) (string, error) {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(metrics)
	dst := h.Sum(nil)
	hashHex := hex.EncodeToString(dst)
	return hashHex, nil
}
