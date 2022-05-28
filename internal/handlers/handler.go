package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	config "github.com/OlesyaBelochka/My-go-musthave-devops/internal"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/compression"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/prhash"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/storage"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/storage/db"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/storage/inmemory"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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

func sendResponceJSON(w http.ResponseWriter, status int, needCompression bool, e error, h bool) {
	resp := variables.ServResponses{}

	if status == http.StatusOK {

		resp = variables.ServResponses{
			Result: "Data update succesfully",
			Error:  "",
		}
	} else {
		resp = variables.ServResponses{
			Result: "Unsuccesfully",
			Error:  e.Error(),
		}
	}

	strJSON, err := json.Marshal(resp)

	fmt.Println("#(sendResponceJSON) " + string(strJSON))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("(sendResponceJSON) Marshal error: " + err.Error()))
		fmt.Println("# (sendResponceJSON) Marshal error: " + err.Error())
		return
	}

	if needCompression {
		w.Header().Set("Content-Encoding", "gzip")
		strJSON, err = compression.Compress(strJSON)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("# (sendResponceJSON) Compress error : " + err.Error()))
			fmt.Println("# (sendResponceJSON) Compress error : " + err.Error())
			return
		}
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	if _, err = w.Write(strJSON); err != nil {
		variables.PrinterErr(err, ""+"- Send error")
		return
	}

	fmt.Println("#(sendResponceJSON) responce succesfully")

}

func HandleGetAllMetrics(w http.ResponseWriter, r *http.Request) {

	if variables.ShowLog {
		log.Print("HandleGetAllMetrics")
	}
	html := ""
	for s, c := range variables.MG {
		html += fmt.Sprintf("%s : %.3f\n", s, c)
	}

	for s, c := range variables.MC {

		html += fmt.Sprintf("%s : %d\n", s, c)

	}

	w.Header().Set("Content-Type", "text/html")

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		fmt.Println("(HandleGetAllMetrics)  сжимает файл чтобы отправить ответ")

		data, err := compression.Compress([]byte(html))
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		if _, err := w.Write(data); err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		if _, err := w.Write([]byte(html)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)

}

func getMetric(mType, mName string, format bool) (string, int, error) {

	var answer string
	var st int
	var err error

	if mName == "" || mType == "" {
		st = http.StatusBadRequest
		return answer, st, err
	}

	switch strings.ToLower(mType) {
	case "gauge":

		if value, inMap := storage.MGServer.Get(mName); inMap {
			// мы получили байты теперь преобразуем их в float
			byteToFloat, _ := strconv.ParseFloat(string(value), 64)

			if format {
				answer = fmt.Sprintf("%0.3f", byteToFloat)
			} else {

				answer = strconv.FormatFloat(byteToFloat, 'f', -1, 64)
			}

			st = http.StatusOK
		} else {
			err = fmt.Errorf("не найдено имя %s", mName)
			st = http.StatusNotFound //404
			answer = ""
		}

	case "counter":
		if value, inMap := storage.MCServer.Get(mName); inMap {
			// мы получили байты теперь преобразуем их в int
			byteToInt, _ := strconv.ParseInt(string(value), 10, 64)

			if format {
				answer = fmt.Sprintf("%d", byteToInt)
			} else {
				answer = strconv.FormatInt(byteToInt, 10)
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

func HandleUpdateMetrics(w http.ResponseWriter, r *http.Request) {

	if storage.MCServer == nil {
		storage.MGServer = inmemory.NewGaugeMS()
	}

	if storage.MCServer == nil {
		storage.MCServer = inmemory.NewCounterMS()
	}

	mType := chi.URLParam(r, "mType")
	mName := chi.URLParam(r, "mName")
	mVal := chi.URLParam(r, "mValue")

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

		storage.MGServer.Set(mName, []byte(strconv.FormatFloat(float64(val), 'f', -1, 64)))

		sendStatus(w, http.StatusOK)

	case "counter":

		val, err := strconv.Atoi(mVal)

		if err != nil {
			sendStatus(w, http.StatusBadRequest) // 400
			return
		}

		storage.MCServer.Set(mName, []byte(strconv.FormatInt(int64(val), 10)))

		sendStatus(w, http.StatusOK)

	default:
		sendStatus(w, http.StatusNotImplemented) // 501
	}

}

func readBodyJSONRequest(w http.ResponseWriter, r *http.Request, resp *variables.Metrics, needCompression *bool) (int, error, bool) {

	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		*needCompression = true
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Println(err)
		return http.StatusInternalServerError, err, false
	}

	if *needCompression {
		body, err = compression.Decompress(body)
		variables.PrinterErr(err, "#HandleGetMetricJSON mistake decompression: ")

	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("can't unmarshal: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return http.StatusInternalServerError, err, false
	}

	// Проверим пришла ли хэш функция и если она пришла но не равна получаемой на сервере
	// то отправляем лесом
	var isHash bool
	if resp.Hash != "" {

		getSHash := ""
		switch resp.MType {

		case "gauge":
			getSHash = prhash.Hash(fmt.Sprintf("%s:%s:%f", resp.ID, resp.MType, *resp.Value), config.ConfS.Key)
		case "counter":
			getSHash = prhash.Hash(fmt.Sprintf("%s:%s:%d", resp.ID, resp.MType, *resp.Delta), config.ConfS.Key)
		default:
			return http.StatusBadRequest, errors.New("не смогли посчитать хэш на сервере так как получен неверны тип метрики"), false
		}

		fmt.Println("полученный хеш: ", getSHash, " посчитанный хеш: ", resp.Hash)
		if getSHash != resp.Hash {
			return http.StatusBadRequest, errors.New("Хеши не равны"), false
		}
		isHash = true

	}

	return http.StatusOK, nil, isHash

}

// HandleGetMetricJSON возвращает метрику в виде JSON
func HandleGetMetricJSON(w http.ResponseWriter, r *http.Request) {

	var (
		resp            variables.Metrics
		needCompression bool
	)

	st, err, isHash := readBodyJSONRequest(w, r, &resp, &needCompression)

	if err != nil {
		// это значит что мы по какой-то причине не смогли выполнить процедуру выше и отправляем статус
		sendStatus(w, st) // 400
		return
	}

	fmt.Println(isHash)

	mType := resp.MType
	mName := resp.ID

	val, code, err := getMetric(mType, mName, false)

	if err != nil {
		fmt.Println("#mistake getMetric: ", code, err)
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

	case "counter":
		valInt, err := strconv.ParseInt(val, 10, 64)

		if err != nil {
			fmt.Println(err)
			return
		}
		resp.Delta = &valInt

	}

	strJSON, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err)
		return
	}

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		fmt.Println("(HandleGetMetricJSON)  сжимает файл чтобы отправить ответ")
		w.Header().Set("Content-Encoding", "gzip")
		strJSON, err = compression.Compress(strJSON)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("# (sendResponceJSON) Compress error : " + err.Error()))
			fmt.Println("# (sendResponceJSON) Compress error : " + err.Error())
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(strJSON)

	if err != nil {
		fmt.Println(err)
	}
}

func HandleUpdateMetricsJSON(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(HandleUpdateMetricsJSON)

	var (
		metrics         variables.Metrics
		needCompression bool
	)

	st, err, isHash := readBodyJSONRequest(w, r, &metrics, &needCompression)

	if err == nil {

		mType := metrics.MType
		mName := metrics.ID

		if mName == "" || (mType != "gauge" && mType != "counter") {
			mType = ""
			err = errors.New("can't find gauge or counter or empty id")
			st = http.StatusNotImplemented
		}

		switch strings.ToLower(mType) {

		case "gauge":
			val := *metrics.Value

			storage.MGServer.Set(mName, []byte(strconv.FormatFloat(float64(val), 'f', -1, 64)))

			st = http.StatusOK

		case "counter":
			val := *metrics.Delta

			storage.MCServer.Set(mName, []byte(strconv.FormatInt(int64(val), 10)))

			st = http.StatusOK
		}
	}

	sendResponceJSON(w, st, needCompression, err, isHash)

}

func HandlePingDB(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)

	defer cancel()

	dataBase, err := db.OpenDB(config.ConfS)

	if err != nil {
		fmt.Println("ошибка при открытии БД", err)
		sendStatus(w, http.StatusInternalServerError)
	}

	defer func() { _ = dataBase.Close() }()

	if err := db.InitSchema(ctx, dataBase); err != nil {
		fmt.Println("ошибка при создании инициализации схемы", err)
		sendStatus(w, http.StatusInternalServerError)
	}

	sendStatus(w, http.StatusOK)

}
