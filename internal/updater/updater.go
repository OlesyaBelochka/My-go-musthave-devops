package updater

import (
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
)

func UpdateGaugeMetric(name string, val variables.Gauge) {
	//if variables.ShowFullLog {
	//	log.Printf("обновляем метку %v  в значение %v", name, val)
	//}

	variables.MG[name] = val

}

func UpdateCountMetric(name string, val variables.Counter) {

	//if variables.ShowFullLog {
	//	log.Printf("обновляем сounter метку %v  если уже существует добавляем %v", name, val)
	//}

	variables.MC[name] += val
}
