package utils

import (
	"crypto/sha512"
	"encoding/hex"
)

func GenerateSHAString(pwd string) string {
	algorithm := sha512.New()
	algorithm.Write([]byte(pwd))
	hashString := hex.EncodeToString(algorithm.Sum(nil))
	return hashString
}

func CompareSHAString(pwd string, pwd2 string) bool {
	return false
}
