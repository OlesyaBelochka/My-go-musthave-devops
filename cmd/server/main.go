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
	"log"
	"net/http"
	"time"
)

var (
	fRstor              bool
	fAddr, fStrFile     string
	fStrInterv          time.Duration
	fRpInterv, fPInterv int64
)

func init() {

	flag.BoolVar(&config.FRstor, "r", config.DefaultRestore, "RESTORE=<ЗНАЧЕНИЕ>")
	flag.StringVar(&config.FАddr, "a", config.DefaultAddress, "ADDRESS=<ЗНАЧЕНИЕ>")
	flag.StringVar(&config.FStrFile, "f", config.DefaultStoreFile, "STORE_FILE=<ЗНАЧЕНИЕ>")
	flag.DurationVar(&config.FStrInterv, "i", config.DefaultStoreInterval, "STORE_INTERVAL=<ЗНАЧЕНИЕ>")

}

func main() {
	flag.Parse()

	fmt.Println("парсим флаги начало")
	fmt.Println("config.UseFlagRstor = ", config.FRstor)
	fmt.Println("config.FАddr = ", config.FАddr)
	fmt.Println("config.FStrFile = ", config.FStrFile)
	fmt.Println("config.FStrInterv = ", config.FStrInterv)
	fmt.Println("парсим флаги конец")
	//setFlags()

	variables.ConfS = config.NewS()

	if variables.ConfS.Restore {
		fmt.Println("start RestoreMetricsFromFile")
		go files.RestoreMetricsFromFile()
	}

	log.Println("Server has started, listening IP: " + variables.ConfS.Address)

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

	if variables.ConfS.Address != "" {
		http.ListenAndServe(variables.ConfS.Address, r)
		//http.ListenAndServe("127.0.0.1:8080", r)
	}

}
