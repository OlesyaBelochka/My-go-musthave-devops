package files

import (
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/config"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/storage"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"io"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func saveMetricsIntoFile() {

	newWriter, err := NewWriter(config.VarConfServer.StoreFile)

	if err != nil {
		variables.PrinterErr(err, "(RestoreMetricsFromFile)can't open file, error: ")
		return

	}

	defer newWriter.Close()

	for k, v := range variables.MG {

		vFl := float64(v)
		str := variables.Metrics{
			ID:    k,
			MType: "gauge",
			Value: &vFl,
		}

		if err := newWriter.WriteData(&str); err != nil {
			variables.PrinterErr(err, "(RestoreMetricsFromFile)mistake while writening file gauge metrics: ")
		}
	}

	for k, v := range variables.MC {

		vInt := int64(v)
		str := variables.Metrics{
			ID:    k,
			MType: "counter",
			Delta: &vInt,
		}

		if err := newWriter.WriteData(&str); err != nil {
			variables.PrinterErr(err, "(RestoreMetricsFromFile) mistake while writening file counter metrics")
		}

	}

}

func Start() {

	timerStore := time.NewTimer(config.VarConfServer.StoreInterval * time.Second)

	osSigChan := make(chan os.Signal, 4)
	signal.Notify(osSigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {

		select {
		case <-timerStore.C:
			saveMetricsIntoFile()

		case <-osSigChan:
			saveMetricsIntoFile()
			os.Exit(1)
			return
		}
	}
}

func RestoreMetricsFromFile() {

	readerM, err := NewReader(config.VarConfServer.StoreFile)
	variables.PrinterErr(err, "(RestoreMetricsFromFile) не смогли создать NewReader, ошибка: ")

	defer readerM.Close()

	for {
		readedData, err := readerM.ReadData()
		if err == io.EOF { // если конец файла
			break // выходим из цикла
		} else if err != nil {
			variables.PrinterErr(err, "(RestoreMetricsFromFile) произошла ошибка в процессе чтения файла: ")
		}

		switch readedData.MType {

		case "gauge":
			storage.MGServer.Set(readedData.ID, []byte(strconv.FormatFloat(float64(*readedData.Value), 'f', -1, 64)))

		case "counter":
			storage.MGServer.Set(readedData.ID, []byte(strconv.FormatInt(*readedData.Delta, 10)))

		}
	}

}
