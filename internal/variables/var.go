package variables

type Gauge float64
type Counter int64

const IPServer = "127.0.0.1:8080"
const ShowLog = true

var MG = map[string]Gauge{}

var MC = map[string]Counter{}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
