package main

import (
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/reporters"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/updater"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
)

var (
	fАddr               string
	fRpInterv, fPInterv int64
)

func init() {
	// loads values from .env into the system
	//if err := godotenv.Load(); err != nil {
	//	log.Print("No .env file found")
	//}
	//
	//flag.StringVar(&fАddr, "a", "", "ADDRESS=<ЗНАЧЕНИЕ>")
	//flag.Int64Var(&fRpInterv, "f", 10, "REPORT_INTERVAL=<ЗНАЧЕНИЕ>")
	//flag.Int64Var(&fPInterv, "f", 2, "POLL_INTERVAL=<ЗНАЧЕНИЕ>")

}

//func setFlags() {
//	flag.Parse()
//
//	if fАddr != "" {
//		variables.Conf.Address = f_addr
//	}
//
//	if fRpInterv != 0 {
//		variables.Conf.ReportInterval = f_rp_interv
//	}
//
//	if fPInterv != 0 {
//		variables.Conf.PollInterval = f_p_interv
//	}
//}

func main() {
	log.Println("Клиент запустился, обновляет и отправляет")
	//conf := config.New()

	//setFlags()

	if variables.ShowLog {
		fmt.Printf("Address %v, ReportInterval = %v, PollInterval =  %v", variables.Address, variables.ReportInterval, variables.PollInterval)
	}

	endpoint := "/update/"

	client := http.Client{}

	timer10 := time.NewTimer(time.Duration(variables.ReportInterval) * time.Second)

	for {
		osSigChan := make(chan os.Signal)
		signal.Notify(osSigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		timer := time.NewTimer(time.Duration(variables.PollInterval) * time.Second)

		select {

		case <-timer10.C:

			timer10 = time.NewTimer(time.Duration(variables.ReportInterval) * time.Second)
			updater.Pall()
			reporters.Report("http://127.0.0.1:8080"+endpoint, client)

		case <-timer.C:
			updater.Pall()
		case <-osSigChan:

			os.Exit(1)
			return
		}
	}

}
