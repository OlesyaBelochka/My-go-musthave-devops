package main

import (
	"context"
	"fmt"
	config "github.com/OlesyaBelochka/My-go-musthave-devops/internal"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/poller"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/reporters"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	// loads values from .env into the system
	//if err := godotenv.Load(); err != nil {
	//	log.Print("No .env file found")
	//}
}

func main() {
	config.Client = http.Client{}
	config.ConfA = config.NewA()
	log.Println("Client started, update and report to IP ", config.ConfA.Address)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if variables.ShowLog {
		fmt.Printf("Address %s, ReportInterval = %d, PollInterval =  %d \n", config.ConfA.Address, config.ConfA.ReportInterval, config.ConfA.PollInterval)
	}
	config.EndpointAgent = "/update/"

	osSigChan := make(chan os.Signal, 1)
	signal.Notify(osSigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go poller.PallStart(ctx)

	//go reporters.ReportAgentNew(ctx, config.ConfA.Key)
	go func() {
		<-osSigChan
		fmt.Println("Get signal")
		cancel()
	}()

	for {
		timerReport := time.NewTimer(config.ConfA.ReportInterval)
		select {
		case <-timerReport.C:
			gR := reporters.NewGaugeReporter()
			cR := reporters.NewCounterReporter()
			fmt.Println(gR)
			fmt.Println(cR)

			var reporters []reporters.ReporterInterface
			rep := append(reporters, gR, cR)

			for _, reporterInterface := range rep {
				reporterInterface.Report(config.ConfA.Key)
			}

			variables.FShowLog("#reporting..")
			//		reportMetrics()
		case <-ctx.Done():
			variables.FShowLog("ctx.Done(): Report")
			return

		}

	}
}
