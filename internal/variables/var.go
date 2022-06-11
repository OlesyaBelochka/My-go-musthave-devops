package variables

import (
	"log"
	"os"
)

type Gauge float64
type Counter int64

const (
	ShowLog = false
)

var errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

var MG = map[string]Gauge{}

var MC = map[string]Counter{}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`
}

type ServResponses struct {
	Result string `json:"rez"`
	Error  string `json:"err"`
}

func PrinterErr(err error, srt string) {
	if err != nil {
		errorLog.Println(srt, err)
	}

}
func FShowLog(s string) {
	if ShowLog {
		infoLog.Println(s)
	}
}
