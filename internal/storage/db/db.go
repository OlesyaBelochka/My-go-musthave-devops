package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	_ "github.com/jackc/pgx/stdlib"
	"strconv"
)

type CounterBDStorage struct {
	bd *sql.DB
}

type GaugeBDStorage struct {
	bd *sql.DB
}

func NewCounterMS(bd1 *sql.DB) *CounterBDStorage {
	return &CounterBDStorage{bd: bd1}
}

func NewGaugeMS(bd1 *sql.DB) *GaugeBDStorage {
	return &GaugeBDStorage{bd: bd1}
}

const (
	// TODO(spencer): update the CREATE DATABASE statement in the schema
	//   to pull out the database specified in the DB URL and use it instead
	//   of "photos" below.
	photosSchema = `
CREATE TABLE IF NOT EXISTS metrics (
  id           TEXT,
  mtype 	  TEXT,
  delta		   BIGINT,
  val        DOUBLE PRECISION,
  PRIMARY KEY (id, mtype)
);`
)

func OpenDB(config *internal.ConfigServer) (*sql.DB, error) {
	return sql.Open("pgx", config.DatabaseURL)
}

func InitSchema(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, photosSchema)
	return err
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

func (r CounterBDStorage) Set(name string, val []byte) {

	byteToInt, _ := strconv.ParseInt(string(val), 10, 64)
	insertSQL := `INSERT INTO metrics VALUES ($1, $2, $3 ,0)
	ON CONFLICT (id,mtype) DO UPDATE SET delta =(metrics.delta + ($3));`

	_, err := r.bd.Exec(insertSQL, name, "counter", byteToInt)

	if err != nil {
		variables.PrinterErr(err, fmt.Sprintf("ошибка при вставке counter %s, %f", name, byteToInt))
		return
	}

	variables.FShowLog(fmt.Sprintf("Add  in BDCounter %s, in val = %v \n", name, byteToInt))

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

func (r GaugeBDStorage) Get(name string) ([]byte, bool) {
	var value float64
	selectSQL := `
	SELECT val
	FROM metrics 
	WHERE id=$1 AND mtype=$2;
	`
	r.bd.QueryRow(selectSQL, name, "gauge").Scan(&value)

	variables.FShowLog(fmt.Sprintf("(Get: GaugeBDStorage) получили значение метрики из БД с типом Gauge  и менем ", name, " значение = ", value))

	if value != 0 {
		return []byte(strconv.FormatFloat(value, 'f', -1, 64)), true
	}

	return []byte(""), false // пустой список байт

}

func (r CounterBDStorage) Get(name string) ([]byte, bool) {
	fmt.Print("Зашли в функцию Get Counter")

	var value int64
	selectSQL := `
	SELECT delta
	FROM metrics 
	WHERE id=$1 AND mtype=$2;
	`
	//fmt.Print("выполняем запрос Get Counter :", selectSQL)
	r.bd.QueryRow(selectSQL, name, "counter").Scan(&value)

	variables.FShowLog(fmt.Sprintf("(Get: CounterBDStorage) получили значение метрики из БД с типом Counter  и менем ", name, " значение = ", value))

	if value != 0 {
		return []byte(strconv.FormatInt(value, 10)), true
	}

	return []byte(""), false // пустой список байт

}
