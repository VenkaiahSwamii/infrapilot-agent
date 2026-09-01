package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateAPIKey() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return "ip_live_" + hex.EncodeToString(bytes)
}
