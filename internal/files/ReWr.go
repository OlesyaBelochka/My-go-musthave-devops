package files

import (
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/storage"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func saveMetricsIntoFile() {

	newWriter, err := variables.NewWriter(internal.ConfS.StoreFile)

	if err != nil {
		log.Println("can't open file, error: ", err)
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

		log.Printf("id : %s, type:%s, value:%v", k, "gauge", vFl)
		if err := newWriter.WriteData(&str); err != nil {
			log.Println("mistake while writening file gauge metrics ", err)
		}

	}

	for k, v := range variables.MC {

		vInt := int64(v)
		str := variables.Metrics{
			ID:    k,
			MType: "counter",
			Delta: &vInt,
		}
		log.Printf("id : %s, type:%s, value:%v", k, "counter", vInt)
		if err := newWriter.WriteData(&str); err != nil {

			log.Println("mistake while writening file counter metrics ", err)
		}

	}

}

func Start() {

	timerStore := time.NewTimer(time.Duration(internal.ConfS.StoreInterval) * time.Second)

	osSigChan := make(chan os.Signal, 4)
	signal.Notify(osSigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {

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
	fmt.Println("StoreFile = ", internal.ConfS.StoreFile)
	readerM, err := variables.NewReader(internal.ConfS.StoreFile)
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

			storage.MGServer.Set(readedData.ID, []byte(strconv.FormatFloat(float64(*readedData.Value), 'f', -1, 64)))

		case "counter":

			storage.MGServer.Set(readedData.ID, []byte(strconv.FormatInt(int64(*readedData.Delta), 10)))

		}
	}

}
