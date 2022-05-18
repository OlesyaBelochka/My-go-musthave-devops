package main

import (
	"context"
	"fmt"
	config "github.com/OlesyaBelochka/My-go-musthave-devops/internal"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/reporters"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/updater"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	variables.ConfA = config.NewA()
	log.Println("Client started, update and report to IP ", variables.ConfA.Address)
	ctx, cancel := context.WithCancel(context.Background())
	if variables.ShowLog {
		fmt.Printf("Address %v, ReportInterval = %v, PollInterval =  %v", variables.ConfA.Address, variables.ConfA.ReportInterval, variables.ConfA.PollInterval)
	}
	endpoint := "/update/"
	osSigChan := make(chan os.Signal, 1)
	signal.Notify(osSigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go updater.Pall(ctx)
	go reporters.Report(ctx, "http://"+variables.ConfA.Address+endpoint)
	sigEnd := <-osSigChan
	fmt.Println("Get signal", sigEnd)
	cancel()

}
