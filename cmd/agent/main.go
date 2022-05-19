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
	go poller.PallStart(ctx)
	go reporters.ReportAgent(ctx)

	osSigChan := make(chan os.Signal, 1)
	signal.Notify(osSigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	//go func() {
	//	<-osSigChan
	//	fmt.Println("Get signal")
	//	cancel()
	//}()

	sigEnd := <-osSigChan
	fmt.Println("Get signal", sigEnd)

}
