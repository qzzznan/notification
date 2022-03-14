package util

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/satori/go.uuid"
)

func GenUID() string {
	return uuid.NewV4().String()
}

func GenToken(token string) string {
	checkSum := sha256.Sum256([]byte(token))
	arr := make([]byte, 0, len(checkSum))
	for _, v := range checkSum {
		arr = append(arr, v)
	}
	return hex.EncodeToString(arr) + "-SHA256"
}
