package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/updater"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"github.com/go-chi/chi/v5"
)

func sendStatus(w http.ResponseWriter, status int) {

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status) // 404

	if status != http.StatusOK {
		http.Error(w, strconv.Itoa(status), status)
	}

}

func sendStatusJSON(w http.ResponseWriter, status int) {

	if status != http.StatusOK {

		strJSON, err := json.Marshal(variables.Metrics{})
		variables.PrinterErr(err, "HandleUpdateMetricsJSON "+"- Marshal error")

		//	fmt.Println("ответ в файле JSON: " + string(strJSON))
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Status-URI", string(status))

		_, err = w.Write(strJSON)
		variables.PrinterErr(err, "HandleUpdateMetricsJSON"+"- Send error")
	}

}

func HandleGetAllMetrics(w http.ResponseWriter, r *http.Request) {

	if variables.ShowLog {
		log.Print("HandleGetAllMetrics")
	}

	w.Header().Set("Content-Type", "text/plain")

	html := ""
	for s, c := range variables.MG {

		html += fmt.Sprintf("%s : %.3f\n", s, c)

	}
	for s, c := range variables.MC {

		html += fmt.Sprintf("%s : %d\n", s, c)
	}
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte(html))

	if err != nil {
		fmt.Println(err)
	}
}

func getMetric(mType, mName string, format bool) (string, int, error) {

	var answer string
	var st int
	var err error

	//if variables.ShowLog {
	//	fmt.Println(mName)
	//	fmt.Println(mType)
	//}

	if mName == "" || mType == "" {
		st = http.StatusBadRequest
		return answer, st, err
	}

	switch strings.ToLower(mType) {
	case "gauge":

		if value, inMap := variables.MG[mName]; inMap {

			if format {
				//fmt.Println("getMetric with format")
				answer = fmt.Sprintf("%0.3f", value)
			} else {

				answer = fmt.Sprintf("%0.10f", value)
			}

			st = http.StatusOK
		} else {
			err = fmt.Errorf("не найдено имя %s", mName)
			st = http.StatusNotFound //404
			answer = ""
		}

	case "counter":
		if value, inMap := variables.MC[mName]; inMap {
			if format {
				//fmt.Println("getMetric with format")
				answer = fmt.Sprintf("%d", value)
			} else {
				//fmt.Println("getMetric except format")
				answer = fmt.Sprintf("%d", value)
			}
			st = http.StatusOK
		} else {
			err = fmt.Errorf("не найдено имя %s", mName)
			st = http.StatusNotFound //404
			answer = ""
		}
	default:

		st = http.StatusBadRequest
	}

	//if variables.ShowFullLog {
	fmt.Println("вот такой ответ дала процедура getMetric", answer, st, err)
	//}
	return answer, st, err
}

func HandleGetMetric(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("star HandleGetMetric old")

	mType := chi.URLParam(r, "mType")
	mName := chi.URLParam(r, "mName")

	//fmt.Println("type metric: ", mType, " name metric: ", mName)

	val, code, err := getMetric(mType, mName, true)

	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}
	fmt.Println("подобрали по типу: ", mType, " и  имени : ", mName, " значение метрики ", val, " в HandleGetMetric old")
	_, err = w.Write([]byte(val))

	if err != nil {
		fmt.Println(err)

	}

}

func HandleGetMetricJSON(w http.ResponseWriter, r *http.Request) {

	//fmt.Println("star HandleGetMetric Json")

	var resp variables.Metrics

	body, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println(w, "can't unmarshal: ", err.Error())
		return
	}

	mType := resp.MType
	mName := resp.ID

	val, code, err := getMetric(mType, mName, false)

	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}
	switch mType {

	case "gauge":
		valFl, err := strconv.ParseFloat(val, 64)

		if err != nil {
			fmt.Println(err)
			return
		}
		resp.Value = &valFl

		if variables.ShowLog {
			fmt.Println("подобрали по типу: ", mType, " и  имени : ", mName, " значение метрики ", valFl, " в HandleGetMetric Json")
		}

	case "counter":
		valInt, err := strconv.ParseInt(val, 10, 64)

		if err != nil {
			fmt.Println(err)
			return
		}
		resp.Delta = &valInt

		if variables.ShowLog {
			fmt.Println("подобрали по типу: ", mType, " и  имени : ", mName, " значение метрики ", valInt, " в HandleGetMetric Json")
		}
		//default:
		//	st := http.StatusBadRequest

	}

	strJSON, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err)
		return
	}

	//	fmt.Println("ответ в файле JSON: " + string(strJSON))

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(strJSON)

	if err != nil {
		fmt.Println(err)
	}

}

func HandleUpdateMetrics(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("HandleUpdateMetrics old")
	//var a = strings.Split(r.URL.String(), "/")

	mType := chi.URLParam(r, "mType")
	mName := chi.URLParam(r, "mName")
	mVal := chi.URLParam(r, "mValue")

	if variables.ShowLog {
		fmt.Println("получен URL ", r.URL)
		fmt.Printf("распознали: mType=%s, mName=%s, mVal = %v", mType, mName, mVal)
	}

	if mName == "" || mVal == "" || (mType != "gauge" && mType != "counter") {
		sendStatus(w, http.StatusNotImplemented) // 501
		return
	}

	switch strings.ToLower(mType) {

	case "gauge":
		val, err := strconv.ParseFloat(mVal, 64)
		fmt.Println(val)
		if err != nil {
			sendStatus(w, http.StatusBadRequest) // 400
			return
		}

		updater.UpdateGaugeMetric(mName, variables.Gauge(val))

		sendStatus(w, http.StatusOK)

	case "counter":

		val, err := strconv.Atoi(mVal)

		if err != nil {
			sendStatus(w, http.StatusBadRequest) // 400
			return
		}

		updater.UpdateCountMetric(mName, variables.Counter(val))

		sendStatus(w, http.StatusOK)

	default:
		sendStatus(w, http.StatusNotImplemented) // 501
	}

}

func HandleUpdateMetricsJSON(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandleUpdateMetricsJSON")

	//var a = strings.Split(r.URL.String(), "/")
	var resp variables.Metrics

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
	err = json.Unmarshal(body, &resp)

	if err != nil {
		fmt.Println(w, "can't unmarshal: ", err.Error())
	}

	mType := resp.MType
	mName := resp.ID

	//fmt.Println("type metric: ", mType, " name metric: ", mName)
	fmt.Println("сюда 1")

	if mName == "" || (mType != "gauge" && mType != "counter") {
		fmt.Println("сюда 2")
		sendStatusJSON(w, http.StatusNotImplemented) // 501

		return
	}

	switch strings.ToLower(mType) {

	case "gauge":
		fmt.Println("попали в gauge")
		val := *resp.Value

		if err != nil {
			fmt.Println(err)
			sendStatusJSON(w, http.StatusBadRequest) // 400
			return
		}

		updater.UpdateGaugeMetric(mName, variables.Gauge(val))
		sendStatusJSON(w, http.StatusOK)

	case "counter":
		fmt.Println("попали в counter")
		//fmt.Println(*resp.Delta)
		val := *resp.Delta

		updater.UpdateCountMetric(mName, variables.Counter(val))
		sendStatusJSON(w, http.StatusOK)

	default:
		sendStatusJSON(w, http.StatusNotImplemented) // 501
	}

}
