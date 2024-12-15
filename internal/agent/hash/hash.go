package hash

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func SignRequestWithSHA256(metrics []byte, key string) (string, error) {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(metrics)
	dst := h.Sum(nil)
	hashHex := hex.EncodeToString(dst)
	return hashHex, nil
}
func Encrypt(data []byte, key string) ([]byte, error) {
	// Создаем новый AES шифратор с ключом
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании шифратора: %v", err)
	}

	// Создаем массив для шифрованных данных
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]

	// Генерируем случайный вектор инициализации (IV)
	if _, err := rand.Read(iv); err != nil {
		return nil, fmt.Errorf("ошибка при генерации IV: %v", err)
	}

	// Создаем новый шифровальный поток для режима CFB
	stream := cipher.NewCFBEncrypter(block, iv)

	// Шифруем данные
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return ciphertext, nil
}
