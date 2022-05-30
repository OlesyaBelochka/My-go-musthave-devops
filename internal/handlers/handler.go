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

func sendResponceJSON(w http.ResponseWriter, status int, needCompression bool, e error) {
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

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("(sendResponceJSON) Marshal error: " + err.Error()))

		variables.PrinterErr(err, "Marshal error:")
		return
	}

	if needCompression {
		w.Header().Set("Content-Encoding", "gzip")
		strJSON, err = compression.Compress(strJSON)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("# (sendResponceJSON) Compress error : " + err.Error()))
			variables.PrinterErr(err, "(sendResponceJSON) Compress error :")
			return
		}
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	if _, err = w.Write(strJSON); err != nil {
		variables.PrinterErr(err, ""+"- Send error")
		return
	}
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

		data, err := compression.Compress([]byte(html))
		if err != nil {
			variables.PrinterErr(err, "(HandleGetAllMetrics) Ошибка сжатия : ")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		if _, err := w.Write(data); err != nil {
			variables.PrinterErr(err, "(HandleGetAllMetrics) Ошибка отправки : ")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		if _, err := w.Write([]byte(html)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			variables.PrinterErr(err, "(HandleGetAllMetrics) Ошибка отправки : ")
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

	mType := chi.URLParam(r, "mType")
	mName := chi.URLParam(r, "mName")

	val, code, err := getMetric(mType, mName, true)

	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	_, err = w.Write([]byte(val))

	if err != nil {
		variables.PrinterErr(err, "(HandleGetMetric) Ошибка отправки : ")
	}

}

func HandleUpdateMetrics(w http.ResponseWriter, r *http.Request) {

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

func readBodyJSONRequest(w http.ResponseWriter, r *http.Request, resp *variables.Metrics, needCompression *bool, doVerificationHash bool) (int, error) {

	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		*needCompression = true
	}

	body, err := io.ReadAll(r.Body)
	variables.PrinterErr(err, "(readBodyJSONRequest) Ошибка чтения тела запроса : ")

	if err != nil {
		return http.StatusInternalServerError, err
	}

	if *needCompression {
		body, err = compression.Decompress(body)
		variables.PrinterErr(err, "#(readBodyJSONRequest)ошибка декомпрессии: ")
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	variables.FShowLog(fmt.Sprintf(" Сервер получил json %s \n выполняет дальнейшую обработку \n", string(body)))

	err = json.Unmarshal(body, &resp)
	if err != nil {
		variables.PrinterErr(err, "can't unmarshal: ")
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return http.StatusInternalServerError, err
	}

	// Проверим пришла ли хэш функция и если она пришла но не равна получаемой на сервере
	// то отправляем лесом

	if doVerificationHash {
		err = verificationHash(*resp, config.ConfS.Key)
	}

	if err != nil {
		variables.PrinterErr(err, "")
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil

}

func readBodySliceJSONRequest(w http.ResponseWriter, r *http.Request, resp *[]variables.Metrics, needCompression *bool) (int, error) {

	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		*needCompression = true
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		variables.PrinterErr(err, "#(readBodySliceJSONRequest) ошибка чтения тела запроса: ")

		return http.StatusInternalServerError, err
	}

	if *needCompression {
		body, err = compression.Decompress(body)
		variables.PrinterErr(err, "#HandleGetMetricJSON mistake decompression: ")

	}

	variables.FShowLog(fmt.Sprintf(" Сервер получил json %s \n выполняет дальнейшую обработку \n", string(body)))

	err = json.Unmarshal(body, &resp)
	if err != nil {
		variables.PrinterErr(err, "can't unmarshal: ")
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// HandleGetMetricJSON возвращает метрику в виде JSON
func HandleGetMetricJSON(w http.ResponseWriter, r *http.Request) {

	var (
		resp            variables.Metrics
		needCompression bool
	)

	st, err := readBodyJSONRequest(w, r, &resp, &needCompression, false)

	if err != nil {
		// это значит что мы по какой-то причине не смогли выполнить процедуру выше и отправляем статус
		sendStatus(w, st) // 400
		return
	}

	if config.ConfS.Key != "" {
		variables.FShowLog(fmt.Sprintf("при отправке с сервера вычислять хеш = %t", config.ConfS.Key))
	}

	mType := resp.MType
	mName := resp.ID

	val, code, err := getMetric(mType, mName, false)

	if err != nil {
		variables.PrinterErr(err, "(HandleGetMetricJSON) ошибка выполнения функции getMetric: ")
		http.Error(w, err.Error(), code)
		return
	}

	switch mType {

	case "gauge":
		valFl, err := strconv.ParseFloat(val, 64)

		if err != nil {
			variables.PrinterErr(err, "(HandleGetMetricJSON) не можем сделать ParseFloat для gauge")
			return
		}
		resp.Value = &valFl
		if config.ConfS.Key != "" {
			resp.Hash = prhash.Hash(fmt.Sprintf("%s:gauge:%d", resp.ID, resp.Value), config.ConfS.Key)
		}

	case "counter":
		valInt, err := strconv.ParseInt(val, 10, 64)

		if err != nil {
			variables.PrinterErr(err, "(HandleGetMetricJSON) не можем сделать ParseInt для counter")
			return
		}

		resp.Delta = &valInt
		if config.ConfS.Key != "" {
			resp.Hash = prhash.Hash(fmt.Sprintf("%s:counter:%d", resp.ID, resp.Delta), config.ConfS.Key)
		}
	}

	strJSON, err := json.Marshal(resp)
	if err != nil {
		variables.PrinterErr(err, "(HandleGetMetricJSON) не можем сделать json.Marshal для структуры")
		return
	}

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		variables.FShowLog("(HandleGetMetricJSON) Увидели заголовок Accept-Encoding,  сжимает файл чтобы отправить ответ ")
		w.Header().Set("Content-Encoding", "gzip")
		strJSON, err = compression.Compress(strJSON)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("# (HandleGetMetricJSON) Compress error : " + err.Error()))
			variables.PrinterErr(err, "# (HandleGetMetricJSON) Compress error : ")
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if config.ConfS.Key != "" {
		variables.FShowLog("отправили структуру на сервер с хэшами так как сервер имеет ключ" + string(strJSON))
	}

	_, err = w.Write(strJSON)

	if err != nil {
		variables.PrinterErr(err, "#(HandleGetMetricJSON) ошибка отправки ответного файла от сервера: ")
	}
}

func HandleUpdateMetricsJSON(w http.ResponseWriter, r *http.Request) {

	var (
		metrics         variables.Metrics
		needCompression bool
	)

	st, err := readBodyJSONRequest(w, r, &metrics, &needCompression, true)

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

	sendResponceJSON(w, st, needCompression, err)

}

func HandlePingDB(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)

	defer cancel()

	dataBase, err := db.OpenDB(config.ConfS)

	if err != nil {
		variables.PrinterErr(err, "#(HandlePingDB) ошибка при открытии БД: ")
		sendStatus(w, http.StatusInternalServerError)
	}

	defer func() { _ = dataBase.Close() }()

	if err := db.InitSchema(ctx, dataBase); err != nil {
		variables.PrinterErr(err, "#(HandlePingDB)ошибка при создании схемы инициализации: ")
		sendStatus(w, http.StatusInternalServerError)
	}

	sendStatus(w, http.StatusOK)

}

type gaugeM struct {
	id  string
	val float64
}

type gaugeC struct {
	id  string
	val int64
}

func HandleUpdatesSliceMetricsJSON(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)

	defer cancel()

	var (
		slMetrics       []variables.Metrics
		needCompression bool

		mmNameG []string
		mmNameC []string
		mmValG  [][]byte
		mmValC  [][]byte
	)

	st, err := readBodySliceJSONRequest(w, r, &slMetrics, &needCompression)

	if err == nil {

		for _, metrics := range slMetrics {

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

				err := verificationHash(metrics, config.ConfS.Key)

				if err != nil {
					variables.PrinterErr(err, "")
					return
				}
				mmNameG = append(mmNameG, mName)
				mmValG = append(mmValG, []byte(strconv.FormatFloat(float64(val), 'f', -1, 64)))

			case "counter":

				val := *metrics.Delta

				err := verificationHash(metrics, config.ConfS.Key)
				if err != nil {
					variables.PrinterErr(err, "")
					return
				}
				mmNameC = append(mmNameC, mName)
				mmValC = append(mmValC, []byte(strconv.FormatInt(int64(val), 10)))

			}
		}

		if len(mmNameG) > 0 {
			storage.MGServer.SetSlice(ctx, mmNameG, mmValG)
			st = http.StatusOK
		}

		if len(mmNameC) > 0 {
			storage.MCServer.SetSlice(ctx, mmNameC, mmValC)
			st = http.StatusOK
		}

		sendResponceJSON(w, st, needCompression, err)

	}
}

func verificationHash(resp variables.Metrics, partKey2 string) error {

	if resp.Hash != "" {

		getSHash := ""
		partKey1 := ""
		switch resp.MType {

		case "gauge":
			partKey1 = fmt.Sprintf("%s:%s:%f", resp.ID, resp.MType, *resp.Value)
		case "counter":
			partKey1 = fmt.Sprintf("%s:%s:%d", resp.ID, resp.MType, *resp.Delta)
		default:
			return errors.New("не смогли посчитать хэш на сервере так как получен неверны тип метрики")
		}

		getSHash = prhash.Hash(partKey1, partKey2)

		str := fmt.Sprintf("полученный хеш: %s, посчитанный хеш: %s,считали для метки: %s,использовали ключ: %s", resp.Hash, getSHash, partKey1, partKey2)

		if getSHash != resp.Hash {
			return errors.New("Хеши не равны: " + str)
		}
	}

	return nil

}
