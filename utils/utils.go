package utils

import (
	"fmt"
	"os"
	"time"
)

func GetEnv(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists { return "", fmt.Errorf("%s not defined in .env file", key) }

	return value, nil
}

func UnixTimeToTime(unixTime int64) time.Time {
	return time.Unix(unixTime, 0)
}
