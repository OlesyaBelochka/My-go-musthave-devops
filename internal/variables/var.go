package variables

import (
	"encoding/json"
	"fmt"
	config "github.com/OlesyaBelochka/My-go-musthave-devops/internal"
	"os"
	"runtime"
	"time"
)

//const (
//	ReportInterval = 10
//	Address        = "127.0.0.1:8080"
//	PollInterval   = 2
//)

type Gauge float64
type Counter int64

const (
	ShowLog     = false
	ShowFullLog = false
)

var MG = map[string]Gauge{}

var MemSt = new(runtime.MemStats)

var MC = map[string]Counter{}

var ConfS *config.ConfigServer
var ConfA *config.ConfigAgent

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type ServResponses struct {
	Result string `json:"rez"`
	Error  string `json:"err"`
}

type writerM struct {
	file    *os.File
	encoder *json.Encoder
}

type readerM struct {
	file    *os.File
	decoder *json.Decoder
}

func NewWriter(fileName string) (*writerM, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &writerM{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}
func (w *writerM) WriteData(metric *Metrics) error {
	return w.encoder.Encode(&metric)
}

func (w *writerM) Close() error {
	return w.file.Close()
}

func NewReader(fileName string) (*readerM, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &readerM{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (r *readerM) ReadData() (*Metrics, error) {
	met := &Metrics{}
	if err := r.decoder.Decode(&met); err != nil {
		return nil, err
	}
	return met, nil
}

func (r *readerM) Close() error {
	return r.file.Close()
}

func PrinterErr(err error, srt string) {
	if err != nil {
		fmt.Println(srt, err)
		return
	}

}
func FShowLog(s string) {
	if ShowLog {
		fmt.Println(s, time.Now())
	}
}
