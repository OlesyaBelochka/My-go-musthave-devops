package main

import (
	"context"
	"fmt"
	config "github.com/OlesyaBelochka/My-go-musthave-devops/internal"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/files"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/handlers"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/storage"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/storage/db"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/storage/inmemory"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"time"
)

func main() {

	config.ConfS = config.NewS()

	storage.MGServer = inmemory.NewGaugeMS()
	storage.MCServer = inmemory.NewCounterMS()

	if config.ConfS.Restore {
		if config.ConfS.DatabaseURL == "" {
			//Использование этого параметра имеет приоритет над параметром
			//file-storage-path и автоматически задействует функциональность
			//сервера БД
			variables.FShowLog("Данные на сервере читаются из памяти")

			go files.RestoreMetricsFromFile()
		} else {
			variables.FShowLog("Данные на сервере читаются из БД")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			dataBase, err := db.OpenDB(config.ConfS)

			if err != nil {
				fmt.Println("ошибка при открытии БД", err)
			}

			defer func() { _ = dataBase.Close() }()

			if err := db.InitSchema(ctx, dataBase); err != nil {
				fmt.Println("ошибка при создании инициализации схемы", err)
			}
			storage.MGServer = db.NewGaugeMS(dataBase)
			storage.MCServer = db.NewCounterMS(dataBase)
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
	r.Get("/ping", handlers.HandlePingDB)

	r.Post("/update/{mType}/{mName}/{mValue}", handlers.HandleUpdateMetrics)
	r.Post("/update", handlers.HandleUpdateMetricsJSON)
	r.Post("/value", handlers.HandleGetMetricJSON)
	r.Post("/value", handlers.HandleGetMetricJSON)
	r.Post("/updates", handlers.HandleUpdatesSliceMetricsJSON)

	if config.ConfS.Address != "" {
		err := http.ListenAndServe(config.ConfS.Address, r)
		variables.PrinterErr(err, "Error server's listetning")
	}

}
