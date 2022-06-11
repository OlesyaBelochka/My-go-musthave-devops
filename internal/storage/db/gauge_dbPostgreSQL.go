package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"strconv"
)

type GaugeBDStorage struct {
	bd *sql.DB
}

func NewGaugeMS(bd1 *sql.DB) *GaugeBDStorage {
	return &GaugeBDStorage{bd: bd1}
}

func (r GaugeBDStorage) Set(name string, val []byte) {
	byteToFloat, _ := strconv.ParseFloat(string(val), 64)

	insertSQL := `INSERT INTO metrics VALUES ($1, $2, 0 ,$3)
	ON CONFLICT (id,mtype) DO UPDATE SET val=EXCLUDED.val;`

	_, err := r.bd.Exec(insertSQL, name, "gauge", byteToFloat)
	if err != nil {
		variables.PrinterErr(err, fmt.Sprintf("ошибка при вставке gauge %s, %f", name, byteToFloat))
		return
	}

	variables.FShowLog(fmt.Sprintf("Set  in BDGauge %s, in val = %f \n", name, byteToFloat))
}

func (r GaugeBDStorage) SetSlice(ctx context.Context, name []string, val [][]byte) {
	// шаг 1 — объявляем транзакцию
	db := r.bd

	tx, err := db.Begin()
	if err != nil {
		variables.PrinterErr(err, "(SetSlice) Произошла ошибка на шаге 1 :")
		return
	}
	// шаг 1.1 — если возникает ошибка, откатываем изменения
	defer tx.Rollback()

	// шаг 2 — готовим инструкцию
	insertSQL := `INSERT INTO metrics VALUES ($1, $2, 0 ,$3)
	ON CONFLICT (id,mtype) DO UPDATE SET val=EXCLUDED.val;`

	stmt, err := tx.PrepareContext(ctx, insertSQL)
	if err != nil {
		variables.PrinterErr(err, "(SetSlice) Произошла ошибка на шаге 2 :")
		return
	}
	// шаг 2.1 — не забываем закрыть инструкцию, когда она больше не нужна
	defer stmt.Close()

	for i := 0; i < len(name); i++ {

		// шаг 3 — указываем, что каждая метрика будет добавлена в транзакцию
		byteToFloat, _ := strconv.ParseFloat(string(val[i]), 64)

		if _, err = stmt.ExecContext(ctx, name[i], "gauge", byteToFloat); err != nil {
			variables.PrinterErr(err, "(SetSlice) Произошла ошибка на шаге 3, не добавили метрику в транзации: ")
			return
		}
	}
	// шаг 4 — сохраняем изменения
	if err = tx.Commit(); err != nil {
		// шаг 4 — сохраняем изменения
		variables.PrinterErr(err, "(SetSlice) Произошла ошибка на шаге 4, не смогли сохранить метрики в транзакции: ")
		return
	} else {
		variables.FShowLog("Запись метрик в базу данных произошла успешно! Ты молодец!")
	}

}

func (r GaugeBDStorage) Get(name string) ([]byte, bool) {
	var value float64
	selectSQL := `
	SELECT val
	FROM metrics 
	WHERE id=$1 AND mtype=$2;
	`
	if r.bd == nil {
		return []byte(""), false // пустой список байт
	}
	check := new(string)
	stmt, err := r.bd.Prepare(selectSQL)
	if err != nil {
		return []byte(""), false
	}

	row := stmt.QueryRow(name, "gauge")

	if err := row.Scan(check); err != sql.ErrNoRows {
		r.bd.QueryRow(selectSQL, name, "gauge").Scan(&value)
		variables.FShowLog(fmt.Sprintf("(Get: GaugeBDStorage) получили значение метрики из БД с типом Gauge  и менем, %s, значение = %f", name, value))
		return []byte(strconv.FormatFloat(value, 'f', -1, 64)), true
	}

	return []byte(""), false // пустой список байт

}
