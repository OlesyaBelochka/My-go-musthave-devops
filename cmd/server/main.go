package main

import (
	"fmt"
	config "github.com/OlesyaBelochka/My-go-musthave-devops/internal"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/files"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/handlers"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/storage/inmemory"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

var MGServer = inmemory.NewGaugeMS()

var MCServer = inmemory.NewCounterMS()

func init() {
	//loads values from .env into the system
	//if err := godotenv.Load(); err != nil {
	//	log.Print("No .env file found")
	//}

}

func main() {

	config.ConfS = config.NewS()

	if config.ConfS.Restore {
		fmt.Println("start RestoreMetricsFromFile")
		if config.ConfS.DatabaseURL == "" {
			//Использование этого параметра имеет приоритет над параметром
			//file-storage-path и автоматически задействует функциональность
			//сервера БД
			go files.RestoreMetricsFromFile()
		}
	}

	log.Println("Server has started, listening IP: " + config.ConfS.Address)

	r := chi.NewRouter()

	go files.Start()
	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	r.Get("/", handlers.HandleGetAllMetrics)

	r.Route("/value", func(r chi.Router) {
		r.Get("/{mType}/{mName}", handlers.HandleGetMetric)

	})
	r.Get("/ping", handlers.HandlePingDb)

	r.Post("/update/{mType}/{mName}/{mValue}", handlers.HandleUpdateMetrics)
	r.Post("/update", handlers.HandleUpdateMetricsJSON)
	r.Post("/value", handlers.HandleGetMetricJSON)
	r.Post("/value", handlers.HandleGetMetricJSON)

	if config.ConfS.Address != "" {
		err := http.ListenAndServe(config.ConfS.Address, r)
		variables.PrinterErr(err, "Error server's listetning")
	}

}
