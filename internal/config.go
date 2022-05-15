package internal

import (
	"os"
	"strconv"
)

type Config struct {
	Address        string
	PollInterval   int64
	ReportInterval int64
	StoreInterval  int64
	StoreFile      string
	Restore        bool
}

// New returns a new Config struct
func New() *Config {

	return &Config{
		Address:        getEnv("ADDRESS", "127.0.0.1:8080"),
		PollInterval:   getEnvAsInt("POLL_INTERVAL", 2),
		ReportInterval: getEnvAsInt("REPORT_INTERVAL", 10),
		StoreInterval:  getEnvAsInt("STORE_INTERVAL", 300),
		StoreFile:      getEnv("STORE_FILE", "/tmp/devops-metrics-db.json"),
		Restore:        getEnvAsBool("RESTORE", true),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsInt(name string, defaultVal int64) int64 {
	valueStr := getEnv(name, "")
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return int64(value)
	}

	return defaultVal
}

func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}
