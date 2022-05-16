package main

import (
	"context"
	"flag"
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
	//flag.Int64Var(&fRpInterv, "r", 10, "REPORT_INTERVAL=<ЗНАЧЕНИЕ>")
	//flag.Int64Var(&fPInterv, "p", 2, "POLL_INTERVAL=<ЗНАЧЕНИЕ>")

}

func setFlags() {
	flag.Parse()

	if fАddr != "" {
		fmt.Println("Agent set flag Addres", fАddr)
		variables.Conf.Address = fАddr
	} else {
		variables.Conf.Address = config.DefaultAddress
	}

	if fRpInterv != 0 {
		fmt.Println("Agent set flag ReportInterval", fRpInterv)
		variables.Conf.ReportInterval = fRpInterv

	} else {
		variables.Conf.ReportInterval = config.DefaultReportInterval
	}

	if fPInterv != 0 {
		fmt.Println("Agent set flag PollInterval", fPInterv)
		variables.Conf.PollInterval = fPInterv
	} else {

		variables.Conf.PollInterval = config.DefaultPollInterval
	}
}

func main() {

	variables.Conf = config.New()
	log.Println("Client started, update and report to IP ", variables.Conf.Address)

	//setFlags()
	ctx, cancel := context.WithCancel(context.Background())

	if variables.ShowLog {
		fmt.Printf("Address %v, ReportInterval = %v, PollInterval =  %v", variables.Conf.Address, variables.Conf.ReportInterval, variables.Conf.PollInterval)
	}

	endpoint := "/update/"

	osSigChan := make(chan os.Signal, 1)
	signal.Notify(osSigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go updater.Pall(ctx)

	go reporters.Report(ctx, "http://"+variables.Conf.Address+endpoint)

	sigEnd := <-osSigChan
	fmt.Println("Get signal", sigEnd)
	cancel()

}
