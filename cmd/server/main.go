package main

import (
	"fmt"
	config "github.com/OlesyaBelochka/My-go-musthave-devops/internal"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/files"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/handlers"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	//if err := godotenv.Load(); err != nil {
	//	log.Print("No .env file found")
	//}

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
	//fmt.Println("Restore = ", variables.Conf.Restore)
	//RESTORE=true
}

func setFlags() {
	//
	//flag.Parse()
	//if !fRstor {
	//	fmt.Println("variables.Conf.Restore = ", variables.Conf.Restore, " fRstor = ", fRstor)
	//
	//	fmt.Println("Server set flag Restore", fRstor)
	//
	//	if !variables.Conf.Restore {
	//		variables.Conf.Restore = fRstor
	//	}
	//
	//}
	//
	//if fAddr != "" {
	//
	//	fmt.Println("variables.Conf.Address = ", variables.Conf.Address, " fAddr = ", fAddr)
	//
	//	if variables.Conf.Address == "" {
	//		fmt.Println("Server set flag Addres", fAddr)
	//		variables.Conf.Address = fAddr
	//	}
	//}
	//
	//if fStrFile != "" {
	//	fmt.Println("variables.Conf.StoreFileexport = ", variables.Conf.StoreFile, " fStrFile = ", fStrFile)
	//
	//	if variables.Conf.StoreFile == "" {
	//		fmt.Println("Server set flag StoreFile", fStrFile)
	//		variables.Conf.StoreFile = fStrFile
	//	}
	//}
	//
	//if fStrInterv != 0 {
	//	fmt.Println("variables.Conf.StoreInterval = ", variables.Conf.StoreInterval, " fStrInterv = ", fStrInterv)
	//
	//	if variables.Conf.StoreInterval == 0 {
	//		fmt.Println("Server set flag StoreInterval", fStrInterv)
	//		variables.Conf.StoreInterval = int64(fStrInterv)
	//	}
	//}

	if !fRstor && !variables.Conf.Restore {
		variables.Conf.Restore = fRstor
		fmt.Println("Server set flag Restore", fRstor)
	} else {
		variables.Conf.Restore = config.DefaultRestore
	}

	fmt.Println("variables.Conf.Restore = ", variables.Conf.Restore)

	if fAddr != "" && variables.Conf.Address == "" {
		variables.Conf.Address = fAddr
		fmt.Println("Server set flag Addres", fAddr)
	} else {
		variables.Conf.Address = config.DefaultAddress
	}

	fmt.Println("variables.Conf.Address = ", variables.Conf.Address)

	if fStrFile != "" && variables.Conf.StoreFile == "" {
		fmt.Println("Server set flag StoreFile", fStrFile)
		variables.Conf.StoreFile = fStrFile
	} else {
		variables.Conf.StoreFile = config.DefaultStoreFile
	}

	fmt.Println("variables.Conf.StoreFileExport = ", variables.Conf.StoreFile)

	if fStrInterv != 0 && variables.Conf.StoreInterval == 0 {
		variables.Conf.StoreInterval = int64(fStrInterv)
		fmt.Println("Server set flag StoreInterval", fStrInterv)
	} else {
		variables.Conf.StoreInterval = config.DefaultStoreInterval

	}

	fmt.Println("variables.Conf.StoreInterval = ", variables.Conf.StoreInterval)

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
