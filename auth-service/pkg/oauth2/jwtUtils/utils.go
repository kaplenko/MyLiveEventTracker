package jwtUtils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateState() (string, error) {
	state := make([]byte, 32)
	if _, err := rand.Read(state); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(state), nil
}
