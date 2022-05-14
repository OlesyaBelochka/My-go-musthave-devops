package main

import (
	"flag"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/handlers"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

var (
	f_rstor            bool
	f_addr, f_str_file string
	f_str_interv       int64
)

func init() {
	////os.Setenv("RESTORE", "true")
	////os.Setenv("ADDRESS", "127.0.0.1:8080")
	////os.Setenv("STORE_FILE", "/tmp/devops-metrics-db.json")
	//
	////if err := godotenv.Load(); err != nil {
	////	log.Print("No .env file found")
	////}
	//
	//path, exists := os.LookupEnv("RESTORE")
	//
	//if exists {
	//	// Print the value of the environment variable
	//	fmt.Println("Print the value of the environment variable", path)
	//}
	//
	//variables.Conf = config.New()
	//
	//flag.BoolVar(&f_rstor, "r", false, "RESTORE=<ЗНАЧЕНИЕ>")
	//flag.StringVar(&f_addr, "a", "", "ADDRESS=<ЗНАЧЕНИЕ>")
	//flag.StringVar(&f_str_file, "i", "/tmp/devops-metrics-db.json", "STORE_FILE=<ЗНАЧЕНИЕ>")
	//flag.Int64Var(&f_str_interv, "f", 300, "STORE_INTERVAL=<ЗНАЧЕНИЕ>")
	////fmt.Println("Restore = ", variables.Conf.Restore)
	////RESTORE=true
}
func setFlags() {

	flag.Parse()
	if !f_rstor {
		variables.Conf.Restore = f_rstor
	}

	if f_addr != "" {
		variables.Conf.Address = f_addr
	}

	if f_str_file != "" {
		variables.Conf.StoreFile = f_str_file
	}

	if f_str_interv != 0 {
		variables.Conf.StoreInterval = f_str_interv
	}
}

func main() {
	//setFlags()

	//if variables.Conf.Restore {
	//	fmt.Println("start RestoreMetricsFromFile")
	//	go files.RestoreMetricsFromFile()
	//}

	log.Println("Server has started, listening... ")
	r := chi.NewRouter()

	//go files.Start()
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

	//if variables.Conf.Address != "" {
	//http.ListenAndServe(variables.Conf.Address, r)
	http.ListenAndServe("127.0.0.1:8080", r)
	//}

}
