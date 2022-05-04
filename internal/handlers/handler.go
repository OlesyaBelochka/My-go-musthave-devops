package handlers

import (
	"encoding/json"
	"fmt"
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
		panic(err)
	}
}

func getMetric(mType, mName string) (string, int, error) {

	var answer string
	var st int
	var err error

	if variables.ShowLog {
		fmt.Println(mName)
		fmt.Println(mType)
	}

	if mName == "" || mType == "" {
		st = http.StatusBadRequest
		return answer, st, err
	}

	switch strings.ToLower(mType) {
	case "gauge":

		if value, inMap := variables.MG[mName]; inMap {
			answer = fmt.Sprintf("%.3f", value)
			st = http.StatusOK
		} else {
			err = fmt.Errorf("не найдено имя %s", mName)
			st = http.StatusNotFound //404
			answer = ""
		}

	case "counter":

		if value, inMap := variables.MC[mName]; inMap {

			answer = fmt.Sprintf("%d", value)
			st = http.StatusOK
		} else {
			err = fmt.Errorf("не найдено имя %s", mName)
			st = http.StatusNotFound //404
			answer = ""
		}

	default:

		st = http.StatusBadRequest
	}

	if variables.ShowLog {
		fmt.Println(answer, st, err)
	}
	return answer, st, err
}

func HandleGetMetric(w http.ResponseWriter, r *http.Request) {

	mType := chi.URLParam(r, "mType")
	mName := chi.URLParam(r, "mName")

	val, code, err := getMetric(mType, mName)

	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	_, err = w.Write([]byte(val))

	if err != nil {
		panic(err)
	}

}

func HandleGetMetricJson(w http.ResponseWriter, r *http.Request) {
	var resp variables.Metrics

	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println(w, "can't unmarshal: ", err.Error())
	}

	mType := resp.MType
	mName := resp.ID

	fmt.Println("type metric: ", mType, " name metric: ", mName)

	val, code, err := getMetric(mType, mName)

	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}
	switch mType {
	case "gauge":
		val_fl, err := strconv.ParseFloat(val, 64)

		if err != nil {
			panic(err)
		}
		resp.Value = &val_fl
	case "counter":
		val_int, err := strconv.ParseInt(val, 10, 1)

		if err != nil {
			panic(err)
		}
		resp.Delta = &val_int
	}

	strJSON, err := json.MarshalIndent(resp, "", "	")

	fmt.Println(string(strJSON))

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(strJSON)

	if err != nil {
		panic(err)
	}

}

func HandleUpdateMetrics(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandleUpdateMetrics")
	//var a = strings.Split(r.URL.String(), "/")

	mType := chi.URLParam(r, "mType")
	mName := chi.URLParam(r, "mName")
	mVal := chi.URLParam(r, "mValue")

	if variables.ShowLog {
		fmt.Println("HandleUpdateMetrics")
		fmt.Println("mType", mType)
		fmt.Println("mName", mName)
		fmt.Println("mVal", mVal)
		fmt.Println(mName == "", mVal == "", (mType != "gauge" && mType != "counter"))
	}

	if mName == "" || mVal == "" || (mType != "gauge" && mType != "counter") {
		sendStatus(w, http.StatusNotImplemented) // 501
		return
	}

	switch strings.ToLower(mType) {

	case "gauge":
		val, err := strconv.ParseFloat(mVal, 64)

		if err != nil {
			sendStatus(w, http.StatusBadRequest) // 400
			return
		}
		variables.MG[mName] = variables.Gauge(val)
		sendStatus(w, http.StatusOK)

	case "counter":

		val, err := strconv.Atoi(mVal)

		if err != nil {
			sendStatus(w, http.StatusBadRequest) // 400
			return
		}

		variables.MC[mName] += variables.Counter(val)
		sendStatus(w, http.StatusOK)

	default:
		sendStatus(w, http.StatusNotImplemented) // 501
	}

}

func HandleUpdateMetricsJson(w http.ResponseWriter, r *http.Request) {

	fmt.Println("HandleUpdateMetricsJson")
	//var a = strings.Split(r.URL.String(), "/")
	var resp variables.Metrics

	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &resp)

	if err != nil {
		fmt.Println(w, "can't unmarshal: ", err.Error())
	}

	mType := resp.MType
	mName := resp.ID

	//if variables.ShowLog {
	//	fmt.Println("HandleUpdateMetrics")
	//	fmt.Println("mType", mType)
	//	fmt.Println("mName", mName)
	//	fmt.Println("mVal", mVal)
	//	fmt.Println(mName == "", mVal == "", (mType != "gauge" && mType != "counter"))
	//}

	if mName == "" || (mType != "gauge" && mType != "counter") {
		sendStatus(w, http.StatusNotImplemented) // 501
		return
	}

	switch strings.ToLower(mType) {

	case "gauge":
		val := *resp.Value

		//	val, err := strconv.ParseFloat(mVal, 64)

		if err != nil {
			sendStatus(w, http.StatusBadRequest) // 400
			return
		}
		variables.MG[mName] = variables.Gauge(val)
		sendStatus(w, http.StatusOK)

	case "counter":
		val := *resp.Delta

		//val, err := strconv.Atoi(mVal)

		//if err != nil {
		//	sendStatus(w, http.StatusBadRequest) // 400
		//	return
		//}

		variables.MC[mName] += variables.Counter(val)
		sendStatus(w, http.StatusOK)

	default:
		sendStatus(w, http.StatusNotImplemented) // 501
	}

}
