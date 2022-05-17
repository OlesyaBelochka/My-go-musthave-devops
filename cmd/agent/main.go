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

//func setFlags() {
//	//flag.Parse()
//
//	//if fАddr != "" && variables.Conf.Address =="" {
//	//	fmt.Println("Agent set flag Addres", fАddr)
//	//	variables.Conf.Address = fАddr
//	//} else {
//	//	variables.Conf.Address = config.DefaultAddress
//	//}
//
//	if variables.ConfA.Address == "" {
//
//		variables.ConfA.Address = config.DefaultAddress
//		if fАddr != "" {
//			variables.ConfA.Address = fАddr
//		}
//	}
//
//	if fRpInterv != 0 {
//		fmt.Println("Agent set flag ReportInterval", fRpInterv)
//		variables.ConfA.ReportInterval = fRpInterv
//
//	} else {
//		variables.ConfA.ReportInterval = config.DefaultReportInterval
//	}
//
//	if fPInterv != 0 {
//		fmt.Println("Agent set flag PollInterval", fPInterv)
//		variables.ConfA.PollInterval = fPInterv
//	} else {
//
//		variables.ConfA.PollInterval = config.DefaultPollInterval
//	}
//}

func main() {

	variables.ConfA = config.NewA()

	log.Println("Client started, update and report to IP ", variables.ConfA.Address)

	//setFlags()

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
