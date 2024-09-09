package env

import (
	"os"
	"strconv"
)

const Prefix = "STOCK_PLUS_"

func Int(key string) int {
	s := os.Getenv(Prefix + key)
	if s == "" {
		return 0
	}
	i, _ := strconv.Atoi(s)
	return i
}

func String(key string) string {
	return os.Getenv(Prefix + key)
}
