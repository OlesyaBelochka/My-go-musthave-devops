package files

import (
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func saveMetricsIntoFile() {
	log.Println("новый врайтер")
	new_writer, err := variables.NewWriter(variables.Conf.StoreFile)

	if err != nil {
		log.Fatal(err)
	}

	defer new_writer.Close()

	for k, v := range variables.MG {
		log.Println("в цикле заполняем структуру")
		v_fl := float64(v)
		str := variables.Metrics{
			ID:    k,
			MType: "gauge",
			Value: &v_fl,
		}
		log.Println("вот эти данные будем писать", str)
		if err := new_writer.WriteData(&str); err != nil {
			log.Fatal(err)
		}

	}

	for k, v := range variables.MC {

		v_int := int64(v)
		str := variables.Metrics{
			ID:    k,
			MType: "counter",
			Delta: &v_int,
		}
		log.Println("вот эти данные будем писать", str)
		if err := new_writer.WriteData(&str); err != nil {
			log.Fatal(err)
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
			log.Println(" прочто 300 сек, открываем файл чтобы в него записать")
			saveMetricsIntoFile()

		case <-osSigChan:
			log.Println("открываем файл чтобы в него записать")
			saveMetricsIntoFile()

			os.Exit(1)

			return

		}

	}

}
