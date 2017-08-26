package helpers

import (
	"crypto/rand"
	"encoding/hex"
)

// RandomString return random string
func RandomString(len int) (string, error) {
	bytes := make([]byte, len/2)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
