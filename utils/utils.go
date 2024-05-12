package utils

import (
	"crypto/rand"
	"encoding/base64"
	"strconv"
	"strings"
)

func GetLogTag(tag string) string {
	if tag == "error" {
		return "\033[31m[ERROR]\033[0m"
	}
	if tag == "warn" {
		return "\033[33m[WARN]\033[0m"
	}
	if tag == "info" {
		return "\033[32m[INFO]\033[0m"
	}
	if tag == "debug" {
		return "\033[34m[DEBUG]\033[0m"
	}

	return "\033[35m[" + strings.ToUpper(tag) + "]\033[0m"
}

func ToInt(s string, def int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return i
}

func GenerateRandomString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
