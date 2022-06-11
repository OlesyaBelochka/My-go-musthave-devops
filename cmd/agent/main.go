package main

import (
	"context"
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/config"
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

func main() {
	config.Client = http.Client{}
	config.VarConfAgent = config.NewAgentConfig()
	log.Println("Client started, update and report to IP ", config.VarConfAgent.Address)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if variables.ShowLog {
		fmt.Printf("Address %s, ReportInterval = %d, PollInterval =  %d \n", config.VarConfAgent.Address, config.VarConfAgent.ReportInterval, config.VarConfAgent.PollInterval)
	}
	config.EndpointAgent = "/update/"

	osSigChan := make(chan os.Signal, 1)
	signal.Notify(osSigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go poller.PallStart(ctx)

	go func() {
		<-osSigChan
		fmt.Println("Get signal")
		cancel()
	}()

	for {
		timerReport := time.NewTimer(config.VarConfAgent.ReportInterval)
		select {
		case <-timerReport.C:
			gR := reporters.NewGaugeReporter()
			cR := reporters.NewCounterReporter()

			var reporters []reporters.Reporter
			rep := append(reporters, gR, cR)

			for _, r := range rep {
				r.Report(config.VarConfAgent.Key)
			}

			variables.FShowLog("#reporting..")

		case <-ctx.Done():
			variables.FShowLog("ctx.Done(): Report")
			return

		}

	}
}
