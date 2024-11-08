package jwt

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateToken() string {
	random := make([]byte, 18)
	_, err := rand.Read(random)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(random)
}
