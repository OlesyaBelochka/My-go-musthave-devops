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
	DefaultStoreInterval  = 300
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

// New returns a new Config struct

func NewS() *ConfigServer {

	//flag.Parse()

	fmt.Println("(ConfigServer FАddr) =", FАddr)
	if FАddr == "" {
		FАddr = DefaultAddress
	}

	if FStrFile == "" {
		FStrFile = DefaultStoreFile
	}

	//if FRstor != DefaultRestore {
	//	FRstor = DefaultRestore
	//}

	if FStrInterv == 0 {
		FStrInterv = DefaultStoreInterval
	}

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

//func New() *Config {
//	return &Config{
//		Address:        getEnv("ADDRESS", ""),
//		PollInterval:   getEnvAsInt("POLL_INTERVAL", 2),
//		ReportInterval: getEnvAsInt("REPORT_INTERVAL", 10),
//		StoreInterval:  getEnvAsInt("STORE_INTERVAL", 300),
//		StoreFile:      getEnv("STORE_FILE", ""),
//		Restore:        getEnvAsBool("RESTORE", false),
//	}
//
//}

//func NewDefault() *Config {
//
//	return &Config{
//		Address:        DefaultAddress,
//		PollInterval:   DefaultPollInterval,
//		ReportInterval: DefaultReportInterval,
//		StoreInterval:  DefaultStoreInterval,
//		StoreFile:      DefaultStoreFile,
//		Restore:        DefaultRestore,
//	}
//}

//func New() *Config {
//
//	return &Config{
//		Address:        getEnv("ADDRESS", "127.0.0.1:8080"),
//		PollInterval:   getEnvAsInt("POLL_INTERVAL", 2),
//		ReportInterval: getEnvAsInt("REPORT_INTERVAL", 10),
//		StoreInterval:  getEnvAsInt("STORE_INTERVAL", 300),
//		StoreFile:      getEnv("STORE_FILE", "/tmp/devops-metrics-db.json"),
//		Restore:        getEnvAsBool("RESTORE", true),
//	}
//}

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
