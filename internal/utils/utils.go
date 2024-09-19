package utils

import (
	"crypto/sha512"
	"encoding/hex"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

const (
	AccountIDContextKey = contextKey("AccountID")
)

func GenerateSHAString(pwd string) string {
	algorithm := sha512.New()
	algorithm.Write([]byte(pwd))
	hashString := hex.EncodeToString(algorithm.Sum(nil))
	return hashString
}
