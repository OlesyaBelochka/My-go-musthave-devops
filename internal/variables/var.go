package variables

import (
	"encoding/json"
	"log"
	"os"
)

type Gauge float64
type Counter int64

const (
	ShowLog = true
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
		errorLog.Println(srt, err)
	}

}
func FShowLog(s string) {
	if ShowLog {
		infoLog.Println(s)
	}
}
