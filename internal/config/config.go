package config

import (
	"os"
	"strconv"
)

type Config struct {
	Workers   int
	BatchSize int
	USDPivot  string
}

func Load() *Config {
	return &Config{
		Workers:   getInt("FOREX_WORKERS", 2),
		BatchSize: getInt("FOREX_BATCH_SIZE", 2),
		USDPivot:  "USD",
	}
}

func getInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return def
	}
	return n
}
