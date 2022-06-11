package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"strconv"
)

type CounterBDStorage struct {
	bd *sql.DB
}

func NewCounterMS(bd1 *sql.DB) *CounterBDStorage {
	return &CounterBDStorage{bd: bd1}
}

func (r CounterBDStorage) Set(name string, val []byte) {

	byteToInt, _ := strconv.ParseInt(string(val), 10, 64)
	insertSQL := `INSERT INTO metrics VALUES ($1, $2, $3 ,0)
	ON CONFLICT (id,mtype) DO UPDATE SET delta =(metrics.delta + ($3));`

	_, err := r.bd.Exec(insertSQL, name, "counter", byteToInt)

	if err != nil {
		variables.PrinterErr(err, fmt.Sprintf("ошибка при вставке counter %s, %d", name, byteToInt))
		return
	}

	variables.FShowLog(fmt.Sprintf("Add  in BDCounter %s, in val = %v \n", name, byteToInt))

}

func (r CounterBDStorage) SetSlice(ctx context.Context, name []string, val [][]byte) {
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
	insertSQL := `INSERT INTO metrics VALUES ($1, $2, $3 ,0)
	ON CONFLICT (id,mtype) DO UPDATE SET delta =(metrics.delta + ($3));`

	stmt, err := tx.PrepareContext(ctx, insertSQL)
	if err != nil {
		variables.PrinterErr(err, "(SetSlice) Произошла ошибка на шаге 2 :")
		return
	}
	// шаг 2.1 — не забываем закрыть инструкцию, когда она больше не нужна
	defer stmt.Close()

	fmt.Println("// шаг 3 — указываем в цикле, что каждая метрика будет добавлена в транзакцию")
	for i := 0; i < len(name); i++ {

		// шаг 3 — указываем, что каждая метрика будет добавлена в транзакцию
		byteToInt, _ := strconv.ParseInt(string(val[i]), 10, 64)

		if _, err = stmt.ExecContext(ctx, name[i], "counter", byteToInt); err != nil {
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

func (r CounterBDStorage) Get(name string) ([]byte, bool) {

	var value int64
	selectSQL := `
	SELECT delta
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
	row := stmt.QueryRow(name, "counter")
	if err := row.Scan(check); err != sql.ErrNoRows {
		r.bd.QueryRow(selectSQL, name, "counter").Scan(&value)
		variables.FShowLog(fmt.Sprintf("(Get: CounterBDStorage) получили значение метрики из БД с типом Counter  и менем, %s, значение = %v", name, value))
		return []byte(strconv.FormatInt(value, 10)), true
	}

	return []byte(""), false // пустой список байт

}
