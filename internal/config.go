package internal

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

var (
	FRstor              bool
	FStrFile            string
	FStrInterv          time.Duration
	FАddr               string
	FRpInterv, FPInterv int64
)

const (
	DefaultAddress        = "127.0.0.1:8080"
	DefaultStoreInterval  = time.Duration(300 * time.Second)
	DefaultStoreFile      = "/tmp/devops-metrics-db.json"
	DefaultRestore        = true
	DefaultPollInterval   = 2
	DefaultReportInterval = 10
)

type ConfigServer struct {
	Address       string
	StoreInterval int64
	StoreFile     string
	Restore       bool
}

type ConfigAgent struct {
	Address        string
	PollInterval   int64
	ReportInterval int64
}

func NewS() *ConfigServer {

	return &ConfigServer{
		Address:       getEnv("ADDRESS", FАddr),
		StoreInterval: getEnvAsInt("STORE_INTERVAL", int64(FStrInterv)),
		StoreFile:     getEnv("STORE_FILE", FStrFile),
		Restore:       getEnvAsBool("RESTORE", FRstor),
	}

}

func NewA() *ConfigAgent {

	return &ConfigAgent{
		Address:        getEnv("ADDRESS", DefaultAddress),
		PollInterval:   getEnvAsInt("POLL_INTERVAL", DefaultPollInterval),
		ReportInterval: getEnvAsInt("REPORT_INTERVAL", DefaultReportInterval),
	}
}

func getEnv(key string, defaultVal string) string {

	if value, exists := os.LookupEnv(key); exists {
		fmt.Println("Получили переменную окружения ", key, " в значение = ", value)
		return value
	}
	fmt.Println("Взяли дефолтное значение", key, " = ", defaultVal)
	return defaultVal

}

func getEnvAsInt(name string, defaultVal int64) int64 {
	valueStr := getEnv(name, "")
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return int64(value)
	}

	fmt.Println("Взяли дефолтное значение", name, " = ", defaultVal)
	return defaultVal

}

func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	fmt.Println("Взяли дефолтное значение", name, " = ", defaultVal)
	return defaultVal

}
