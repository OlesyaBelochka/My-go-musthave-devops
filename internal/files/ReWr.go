package files

import (
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/updater"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func saveMetricsIntoFile() {

	new_writer, err := variables.NewWriter(variables.Conf.StoreFile)

	if err != nil {
		log.Println("can't open file, error: ", err)
		return

	}

	defer new_writer.Close()

	for k, v := range variables.MG {

		v_fl := float64(v)
		str := variables.Metrics{
			ID:    k,
			MType: "gauge",
			Value: &v_fl,
		}

		log.Printf("id : %s, type:%s, value:%v", k, "gauge", v_fl)
		if err := new_writer.WriteData(&str); err != nil {
			log.Println("mistake while writening file gauge metrics ", err)
		}

	}

	for k, v := range variables.MC {

		v_int := int64(v)
		str := variables.Metrics{
			ID:    k,
			MType: "counter",
			Delta: &v_int,
		}
		log.Printf("id : %s, type:%s, value:%v", k, "counter", v_int)
		if err := new_writer.WriteData(&str); err != nil {

			log.Println("mistake while writening file counter metrics ", err)
		}

	}

}
func Start() {

	timerStore := time.NewTimer(time.Duration(variables.Conf.StoreInterval) * time.Second)

	for {
		osSigChan := make(chan os.Signal)
		signal.Notify(osSigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-timerStore.C:
			log.Println("in 300 sec, open file to write")
			saveMetricsIntoFile()

		case <-osSigChan:
			log.Println("Break signal,  open file to write")
			saveMetricsIntoFile()

			os.Exit(1)

			return

		}

	}

}

func RestoreMetricsFromFile() {
	fmt.Println("StoreFile = ", variables.Conf.StoreFile)
	readerM, err := variables.NewReader(variables.Conf.StoreFile)
	if err != nil {

		log.Println("can't create NewReader from func RestoreMetricsFromFile,  error: ", err)

		//log.Fatal(err)
	}

	defer readerM.Close()

	for {
		readedData, err := readerM.ReadData()

		if err == io.EOF { // если конец файла
			log.Println("", err)
			break // выходим из цикла
		} else if err != nil {
			log.Println("error while reading file, error: ", err)
		}

		switch readedData.MType {

		case "gauge":

			updater.UpdateGaugeMetric(readedData.ID, variables.Gauge(*readedData.Value))

		case "counter":

			updater.UpdateCountMetric(readedData.ID, variables.Counter(*readedData.Delta))

		}
	}

}
