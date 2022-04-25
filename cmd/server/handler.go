package main

import (
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"net/http"
	"strconv"
	"strings"
)

func sendStatus(w http.ResponseWriter, status int) {

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status) // 404

	if status != http.StatusOK {
		http.Error(w, strconv.Itoa(status), status)
	}

}

func HandleGetAllMetrics(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandleUpdateMetrics")

	w.Header().Set("Content-Type", "text/plain")
	html := ""

	for s, counter := range variables.MG {
		fmt.Println(s)
		fmt.Println(counter)

		html += fmt.Sprintf("\n %s : %.3f", s, counter)

	}
	for s, counter := range variables.MC {
		fmt.Println(s)
		fmt.Println(counter)

		html += fmt.Sprintf("\n %s : %d", s, counter)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))

}

func getMetric(a []string) (string, int, error) {

	var answer string
	var st int
	var err error

	if len(a[2]) < 4 {
		st = http.StatusBadRequest
		return answer, st, err
	}

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

	var a = strings.Split(r.URL.String(), "/")

	val, code, err := getMetric(a)

	if err != nil {
		http.Error(w, err.Error(), code)
	} else {
		w.Write([]byte(val))
	}

}

func HandleUpdateMetrics(w http.ResponseWriter, r *http.Request) {

	var a = strings.Split(r.URL.String(), "/")

	if len(a) != 5 && (strings.ToLower(a[2]) != "gauge" || strings.ToLower(a[2]) != "counter") {
		sendStatus(w, http.StatusNotImplemented)
		return
	}

	switch strings.ToLower(a[2]) {

	case "gauge":
		val2, err := strconv.ParseFloat(a[4], 64)

		if err != nil {
			sendStatus(w, http.StatusBadRequest) // 400
			return
		}
		variables.MG[a[3]] = variables.Gauge(val2)
		sendStatus(w, http.StatusOK)

	case "counter":

		val, err := strconv.Atoi(a[4])

		if err != nil {
			sendStatus(w, http.StatusBadRequest) // 400
			return
		}

		variables.MC[a[3]] += variables.Counter(val)
		sendStatus(w, http.StatusOK)

	default:
		sendStatus(w, http.StatusNotImplemented) // 401
	}

}
