package main

import (
	"context"
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/config"
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

	config.VarConfServer = config.NewServerConfig()

	storage.MGServer = inmemory.NewGaugeMS()
	storage.MCServer = inmemory.NewCounterMS()

	if config.VarConfServer.Restore {
		if config.VarConfServer.DatabaseURL == "" {
			//Использование этого параметра имеет приоритет над параметром
			//file-storage-path и автоматически задействует функциональность
			//сервера БД
			variables.FShowLog("Данные на сервере читаются из памяти")

			go files.RestoreMetricsFromFile()
		} else {
			variables.FShowLog("Данные на сервере читаются из БД")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			dataBase, err := db.OpenDB(config.VarConfServer)

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

	log.Println("Server has started, listening IP: " + config.VarConfServer.Address)

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

	if config.VarConfServer.Address != "" {
		err := http.ListenAndServe(config.VarConfServer.Address, r)
		variables.PrinterErr(err, "Error server's listetning")
	}

}
