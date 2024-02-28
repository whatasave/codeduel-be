package utils

import (
	"log"
	"os"
	"time"
)

func GetEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		// fmt.Printf("Environment variable %s not found, using default value %s\n", key, defaultValue)
		log.Printf("[WARN] Environment variable %s not found, using default value %s\n", key, defaultValue)
		return defaultValue
	}
	
	return value
}

func UnixTimeToTime(unixTime int64) time.Time {
	return time.Unix(unixTime, 0)
}
