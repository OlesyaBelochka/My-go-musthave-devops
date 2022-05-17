package internal

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

var (
	FRstor                          bool
	FStrFile                        string
	FStrInterv, FRpInterv, FPInterv time.Duration
	FАddr                           string
)

const (
	DefaultAddress        = "127.0.0.1:8080"
	DefaultStoreInterval  = 300 * time.Second
	DefaultStoreFile      = "/tmp/devops-metrics-db.json"
	DefaultRestore        = true
	DefaultPollInterval   = 2 * time.Second
	DefaultReportInterval = 10 * time.Second
)

type ConfigServer struct {
	Address       string
	StoreInterval time.Duration
	StoreFile     string
	Restore       bool
}

type ConfigAgent struct {
	Address        string
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func NewS() *ConfigServer {

	//flag.BoolVar(&config.FRstor, "r", config.DefaultRestore, "RESTORE=<ЗНАЧЕНИЕ>")
	//flag.StringVar(&config.FАddr, "a", config.DefaultAddress, "ADDRESS=<ЗНАЧЕНИЕ>")
	//flag.StringVar(&config.FStrFile, "f", config.DefaultStoreFile, "STORE_FILE=<ЗНАЧЕНИЕ>")
	//flag.DurationVar(&config.FStrInterv, "i", config.DefaultStoreInterval, "STORE_INTERVAL=<ЗНАЧЕНИЕ>")

	flag.BoolVar(&FRstor, "r", DefaultRestore, "RESTORE=<ЗНАЧЕНИЕ>")
	flag.StringVar(&FАddr, "a", DefaultAddress, "ADDRESS=<ЗНАЧЕНИЕ>")
	flag.StringVar(&FStrFile, "f", DefaultStoreFile, "STORE_FILE=<ЗНАЧЕНИЕ>")
	flag.DurationVar(&FStrInterv, "i", DefaultStoreInterval, "STORE_INTERVAL=<ЗНАЧЕНИЕ>")

	flag.Parse()

	fmt.Println("парсим флаги начало")
	fmt.Println("config.UseFlagRstor = ", FRstor)
	fmt.Println("config.FАddr = ", FАddr)
	fmt.Println("config.FStrFile = ", FStrFile)
	fmt.Println("config.FStrInterv = ", FStrInterv)
	fmt.Println("парсим флаги конец")

	cnf := ConfigServer{
		Address:       getEnv("ADDRESS", FАddr),
		StoreInterval: getEnvAsDur("STORE_INTERVAL", FStrInterv),
		StoreFile:     getEnv("STORE_FILE", FStrFile),
		Restore:       getEnvAsBool("RESTORE", FRstor),
	}

	fmt.Println(cnf)

	return &cnf
}

func NewA() *ConfigAgent {

	flag.StringVar(&FАddr, "a", DefaultAddress, "ADDRESS=<ЗНАЧЕНИЕ>")
	flag.DurationVar(&FRpInterv, "r", DefaultReportInterval, "REPORT_INTERVAL=<ЗНАЧЕНИЕ>")
	flag.DurationVar(&FPInterv, "p", DefaultPollInterval, "POLL_INTERVAL=<ЗНАЧЕНИЕ>")

	flag.Parse()

	fmt.Println("парсим флаги агент начало")
	fmt.Println("config.FАddr = ", FАddr)
	fmt.Println("config.FRpInterv = ", FRpInterv)
	fmt.Println("config.FPInterv = ", FPInterv)
	fmt.Println("парсим флаги агент конец")

	cnf := ConfigAgent{
		Address:        getEnv("ADDRESS", FАddr),
		PollInterval:   getEnvAsDur("POLL_INTERVAL", FPInterv),
		ReportInterval: getEnvAsDur("REPORT_INTERVAL", FRpInterv),
	}

	fmt.Println(cnf)

	return &cnf

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

func getEnvAsDur(name string, defaultVal time.Duration) time.Duration {
	valueStr := getEnv(name, "")

	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
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
