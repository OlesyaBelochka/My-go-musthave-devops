package main

import (
	"flag"
	"fmt"
	config "github.com/OlesyaBelochka/My-go-musthave-devops/internal"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/files"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/handlers"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"time"
)

var (
	fRstor          bool
	fAddr, fStrFile string
	fStrInterv      time.Duration
)

func init() {
	//os.Setenv("RESTORE", "true")
	//os.Setenv("ADDRESS", "127.0.0.1:8080")
	//os.Setenv("STORE_FILE", "/tmp/devops-metrics-db.json")

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	//path, exists := os.LookupEnv("RESTORE")
	//
	//if exists {
	//	// Print the value of the environment variable
	//	fmt.Println("Print the value of the environment variable", path)
	//}

	variables.Conf = config.New()
	//
	//flag.BoolVar(&fRstor, "r", false, "RESTORE=<ЗНАЧЕНИЕ>")
	//flag.StringVar(&fAddr, "a", "", "ADDRESS=<ЗНАЧЕНИЕ>")
	//flag.StringVar(&fStrFile, "f", "/tmp/devops-metrics-db.json", "STORE_FILE=<ЗНАЧЕНИЕ>")
	//flag.DurationVar(&fStrInterv, "i", 300, "STORE_INTERVAL=<ЗНАЧЕНИЕ>")
	//
	//flag.BoolVar(&fRstor, "r", variables.Conf.Restore, "RESTORE=<ЗНАЧЕНИЕ>")
	//flag.StringVar(&fAddr, "a", variables.Conf.Address, "ADDRESS=<ЗНАЧЕНИЕ>")
	//flag.StringVar(&fStrFile, "f", variables.Conf.StoreFile, "STORE_FILE=<ЗНАЧЕНИЕ>")
	//flag.DurationVar(&fStrInterv, "i", 300, "STORE_INTERVAL=<ЗНАЧЕНИЕ>")

	flag.BoolVar(&fRstor, "r", config.DefaultRestore, "RESTORE=<ЗНАЧЕНИЕ>")
	flag.StringVar(&fAddr, "a", config.DefaultAddress, "ADDRESS=<ЗНАЧЕНИЕ>")
	flag.StringVar(&fStrFile, "f", config.DefaultStoreFile, "STORE_FILE=<ЗНАЧЕНИЕ>")
	flag.DurationVar(&fStrInterv, "i", 300, "STORE_INTERVAL=<ЗНАЧЕНИЕ>")

	fmt.Println("Restore = ", variables.Conf.Restore)
	//RESTORE=true
}

func setFlags() {
	flag.Parse()

	if fRstor && variables.Conf.Restore == config.DefaultRestore {
		fmt.Println("variables.Conf.Restore set = ", fRstor)
		variables.Conf.Restore = fRstor
	}

	fmt.Println("variables.Conf.Address = ", variables.Conf.Address, "DefaultAddress = ", config.DefaultAddress, variables.Conf.Address == config.DefaultAddress)

	if fAddr != "" && variables.Conf.Address == config.DefaultAddress {
		fmt.Println("variables.Conf.Address set = ", fAddr)
		variables.Conf.Address = fAddr
	}

	if fStrFile != "" && variables.Conf.StoreFile == config.DefaultStoreFile {
		fmt.Println("Svariables.Conf StoreFile set = ", fStrFile)
		variables.Conf.StoreFile = fStrFile
	}
	fmt.Println(variables.Conf.StoreInterval, config.DefaultStoreInterval, variables.Conf.StoreInterval == config.DefaultStoreInterval)
	if fStrInterv != 0 && variables.Conf.StoreInterval == config.DefaultStoreInterval {
		fmt.Println("variables.Conf.StoreInterval set= ", fStrInterv)
		variables.Conf.StoreInterval = int64(fStrInterv)

	}

}

func main() {

	setFlags()

	if variables.Conf.Restore {
		fmt.Println("start RestoreMetricsFromFile")
		go files.RestoreMetricsFromFile()
	}

	log.Println("Server has started, listening IP: " + variables.Conf.Address)

	r := chi.NewRouter()

	go files.Start()
	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	//if variables.ShowLog {
	//	r.Use(middleware.Logger)
	//}
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	r.Get("/", handlers.HandleGetAllMetrics)

	r.Route("/value", func(r chi.Router) {
		r.Get("/{mType}/{mName}", handlers.HandleGetMetric)
	})

	r.Post("/update/{mType}/{mName}/{mValue}", handlers.HandleUpdateMetrics)
	r.Post("/update", handlers.HandleUpdateMetricsJSON)
	r.Post("/value", handlers.HandleGetMetricJSON)

	if variables.Conf.Address != "" {
		http.ListenAndServe(variables.Conf.Address, r)
		//http.ListenAndServe("127.0.0.1:8080", r)
	}

}
