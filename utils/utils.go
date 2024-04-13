package utils

import (
	"log"
	"os"
	"strings"
	"time"
)

func GetEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		// fmt.Printf("Environment variable %s not found, using default value %s\n", key, defaultValue)
		log.Printf("%s Environment variable %s not found, using default value: %s\n", GetLogTag("warn"), key, defaultValue)
		return defaultValue
	}
	
	return value
}

func GetLogTag(tag string) string {
	if tag == "error"	{ return "\033[31m[ERROR]\033[0m" }
	if tag == "warn"	{ return "\033[33m[WARN]\033[0m" }
	if tag == "info"	{ return "\033[32m[INFO]\033[0m" }
	if tag == "debug"	{ return "\033[34m[DEBUG]\033[0m" }

	return "\033[35m[" + strings.ToUpper(tag) + "]\033[0m"
}

func UnixTimeToTime(unixTime int64) time.Time {
	return time.Unix(unixTime, 0)
}
