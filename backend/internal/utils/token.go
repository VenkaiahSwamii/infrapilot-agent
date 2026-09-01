package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateEnrollmentToken() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
