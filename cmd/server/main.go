package main

import (
	config "github.com/OlesyaBelochka/My-go-musthave-devops/internal"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/files"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/handlers"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/updater"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"log"
	"net/http"
)

func init() {
	variables.Conf = config.New()

	if variables.Conf.Restore {

		readerM, err := variables.NewReader(variables.Conf.StoreFile)
		if err != nil {
			log.Fatal(err)
		}

		defer readerM.Close()

		for {
			readedData, err := readerM.ReadData()

			if err == io.EOF { // если конец файла
				break // выходим из цикла
			}

			//fmt.Println(readedData)

			switch readedData.MType {

			case "gauge":

				updater.UpdateGaugeMetric(readedData.ID, variables.Gauge(*readedData.Value))

			case "counter":

				updater.UpdateCountMetric(readedData.ID, variables.Counter(*readedData.Delta))

			}
		}
	}
}

func main() {

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
	r.Post("/update", handlers.HandleUpdateMetricsJson)
	r.Post("/value", handlers.HandleGetMetricJson)

	http.ListenAndServe(variables.Conf.Address, r)

}
