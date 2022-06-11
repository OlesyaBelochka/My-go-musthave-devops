package config

import (
	"flag"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	FRstor                          bool
	FStrFile                        string
	FStrInterv, FRpInterv, FPInterv time.Duration
	FАddr, FKey, FDb                string
	EndpointAgent                   string
)

const (
	DefaultAddress        = "127.0.0.1:8080"
	DefaultStoreInterval  = 300 * time.Second
	DefaultStoreFile      = "/tmp/devops-metrics-db.json"
	DefaultRestore        = true
	DefaultPollInterval   = 2 * time.Second
	DefaultReportInterval = 10 * time.Second
	DefaultKey            = ""
	DefaultDB             = "" //host=localhost dbname=ya_pr_devops
)

var (
	VarConfAgent  *ConfigAgent
	VarConfServer *ConfigServer
	Client        http.Client
)

type ConfigServer struct {
	Address       string
	StoreInterval time.Duration
	StoreFile     string
	Restore       bool
	Key           string
	DatabaseURL   string
}

type ConfigAgent struct {
	Address        string
	PollInterval   time.Duration
	ReportInterval time.Duration
	Key            string
}

func NewServerConfig() *ConfigServer {

	flag.BoolVar(&FRstor, "r", DefaultRestore, "RESTORE=<ЗНАЧЕНИЕ>")
	flag.StringVar(&FАddr, "a", DefaultAddress, "ADDRESS=<ЗНАЧЕНИЕ>")
	flag.StringVar(&FStrFile, "f", DefaultStoreFile, "STORE_FILE=<ЗНАЧЕНИЕ>")
	flag.DurationVar(&FStrInterv, "i", DefaultStoreInterval, "STORE_INTERVAL=<ЗНАЧЕНИЕ>")
	flag.StringVar(&FKey, "k", DefaultKey, "k=<КЛЮЧ>")
	flag.StringVar(&FDb, "d", DefaultDB, "путь к базе данных")

	flag.Parse()

	cnf := ConfigServer{
		Address:       getEnv("ADDRESS", FАddr),
		StoreInterval: getEnvAsDur("STORE_INTERVAL", FStrInterv),
		StoreFile:     getEnv("STORE_FILE", FStrFile),
		Restore:       getEnvAsBool("RESTORE", FRstor),
		Key:           getEnv("KEY", FKey),
		DatabaseURL:   getEnv("DATABASE_DSN", FDb),
	}

	return &cnf
}

func NewAgentConfig() *ConfigAgent {

	flag.StringVar(&FАddr, "a", DefaultAddress, "ADDRESS=<ЗНАЧЕНИЕ>")
	flag.DurationVar(&FRpInterv, "r", DefaultReportInterval, "REPORT_INTERVAL=<ЗНАЧЕНИЕ>")
	flag.DurationVar(&FPInterv, "p", DefaultPollInterval, "POLL_INTERVAL=<ЗНАЧЕНИЕ>")
	flag.StringVar(&FKey, "k", DefaultKey, "k=<КЛЮЧ>")

	flag.Parse()

	cnf := &ConfigAgent{
		Address:        getEnv("ADDRESS", FАddr),
		PollInterval:   getEnvAsDur("POLL_INTERVAL", FPInterv),
		ReportInterval: getEnvAsDur("REPORT_INTERVAL", FRpInterv),
		Key:            getEnv("KEY", FKey),
	}

	return cnf
}

func getEnv(key string, defaultVal string) string {

	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal

}

func getEnvAsDur(name string, defaultVal time.Duration) time.Duration {
	valueStr := getEnv(name, "")

	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
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
