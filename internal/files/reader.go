package files

import (
	"encoding/json"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"os"
)

type readerM struct {
	file    *os.File
	decoder *json.Decoder
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

func (r *readerM) ReadData() (*variables.Metrics, error) {
	met := &variables.Metrics{}
	if err := r.decoder.Decode(&met); err != nil {
		return nil, err
	}
	return met, nil
}

func (r *readerM) Close() error {
	return r.file.Close()
}
