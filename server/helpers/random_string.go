package helpers

import (
	"crypto/rand"
	"encoding/hex"
)

// RandomString return random string
func RandomString(len int) string {
	bytes := make([]byte, len/2)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(bytes)
}
