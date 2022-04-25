package main

import (
	"encoding/json"
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"net/http"
	"strconv"
	"strings"
)

func sendStatus(w http.ResponseWriter, status int) {

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status) // 404
	//fmt.Println(status)

	http.Error(w, strconv.Itoa(status), status)
	//http.Error(w, strconv.Itoa(status),status)

}
func HandleGetAllMetrics(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandleUpdateMetrics")

	json.NewEncoder(w).Encode(variables.MG)
	json.NewEncoder(w).Encode(variables.MC)
}

//func HandleGetMetric(w http.ResponseWriter, r *http.Request) {
//	fmt.Println("HandleGetMetric")
//
//	var a = strings.Split(r.URL.String(), "/")
//	var answer string
//
//	fmt.Println(r.URL.String())
//
//	//fmt.Println(a)
//	for i, s := range a {
//		fmt.Println(i, s)
//	}
//	fmt.Println(len(a))
//	fmt.Println(strings.ToLower(a[2]))
//	if len(a) >= 4 && (strings.ToLower(a[2]) == "gauge" || strings.ToLower(a[2]) == "counter") {
//
//		if strings.ToLower(a[2]) == "gauge" {
//			if value, inMap := variables.MG[a[3]]; inMap {
//				fmt.Println("нашли имя"+a[3]+" в мапе и его значение = ", value)
//
//				answer = strconv.FormatFloat(float64(value), 'f', 10, 64)
//
//			} else {
//				fmt.Println("не найдено имя " + a[3] + " в мапе")
//				sendStatus(w, http.StatusNotFound) //404
//				return
//			}
//
//		} else {
//			if value, inMap := variables.MC[a[3]]; inMap {
//
//				fmt.Println("нашли имя"+a[3]+" в мапе и его значение = ", value)
//				answer = strconv.FormatInt(int64(value), 10)
//			} else {
//				fmt.Println("не найдено имя " + a[3] + " в мапе")
//				sendStatus(w, http.StatusNotFound) //404
//				return
//			}
//		}
//
//		w.Header().Set("Content-Type", "text/plain")
//		w.WriteHeader(http.StatusOK)
//		w.Write([]byte(answer))
//
//	} else {
//
//		sendStatus(w, http.StatusBadRequest) // 400
//
//	}
//
//}

func getMetric(a []string) (string, int, error) {

	var answer string
	var st int
	var err error

	switch strings.ToLower(a[2]) {
	case "gauge":

		if value, inMap := variables.MG[a[3]]; inMap {
			//	fmt.Println("нашли имя"+a[3]+" в мапе и его значение = ", value)
			answer = fmt.Sprintf("%.3f", value)
			st = http.StatusOK
		} else {
			//	fmt.Println("не найдено имя " + a[3] + " в мапе")
			err = fmt.Errorf("не найдено имя %v", a[3])
			st = http.StatusNotFound //404
			answer = ""
		}

	case "counter":
		if value, inMap := variables.MC[a[3]]; inMap {
			//fmt.Println("нашли имя"+a[3]+" в мапе и его значение = ", value)
			answer = fmt.Sprintf("%d", value)
			st = http.StatusOK
		} else {
			//fmt.Println("не найдено имя " + a[3] + " в мапе")
			err = fmt.Errorf("не найдено имя %v", a[3])
			st = http.StatusNotFound //404
			answer = ""
		}
	default:

		st = http.StatusBadRequest

	}

	return answer, st, err
}

func HandleGetMetric(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("HandleUpdateMetrics")
	var a = strings.Split(r.URL.String(), "/")

	//fmt.Println(r.URL.String())

	val, code, err := getMetric(a)

	if err != nil {
		http.Error(w, err.Error(), code)
	} else {
		w.Write([]byte(val))
	}

}

func HandleUpdateMetrics(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("HandleUpdateMetrics")
	var a = strings.Split(r.URL.String(), "/")

	if len(a) == 5 && (strings.ToLower(a[2]) == "gauge" || strings.ToLower(a[2]) == "counter") {
		if strings.ToLower(a[2]) == "gauge" {
			val2, err := strconv.ParseFloat(a[4], 64)

			if err != nil {
				sendStatus(w, http.StatusBadRequest) // 400
				return
			}

			variables.MG[a[3]] = variables.Gauge(val2)
		} else {
			val, err := strconv.Atoi(a[4])

			if err != nil {
				sendStatus(w, http.StatusBadRequest) // 400
				return
			}

			variables.MC["PollCount"] += variables.Counter(val)

		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		//	w.Write([]byte("метки обновились."))
	} else {

		sendStatus(w, http.StatusNotImplemented) // 401
	}

}
