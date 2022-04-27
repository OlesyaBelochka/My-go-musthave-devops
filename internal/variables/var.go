package variables

type Gauge float64
type Counter int64

const IPServer = "127.0.0.1:8080"
const ShowLog = false

var MG = map[string]Gauge{}

var MC = map[string]Counter{}
