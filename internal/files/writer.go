package files

import (
	"encoding/json"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"os"
)

type writerM struct {
	file    *os.File
	encoder *json.Encoder
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

func (w *writerM) WriteData(metric *variables.Metrics) error {
	return w.encoder.Encode(&metric)
}

func (w *writerM) Close() error {
	return w.file.Close()
}
