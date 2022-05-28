package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal"
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
  delta		   INT,
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
	//insertSQL := `INSERT INTO metrics VALUES ($1, $2, 0 ,$3)

	insertSQL := `INSERT INTO metrics VALUES ($1, $2, 0 ,$3)
	ON CONFLICT (id,mtype) DO UPDATE SET val=EXCLUDED.val;`

	_, err := r.bd.Exec(insertSQL, name, "gauge", byteToFloat)
	if err != nil {
		fmt.Print("ошибка при вставке gauge ", name, " ", byteToFloat, " ошибка :", err)
	}
	fmt.Printf("Set  in BD Gauge %s, in val = %f \n", name, byteToFloat)
}

func (r CounterBDStorage) Set(name string, val []byte) {
	var value int64
	selectSQL := `
	SELECT delta
	FROM metrics 
	WHERE id=$1 AND mtype=$2;
	`

	r.bd.QueryRow(selectSQL, name, "counter").Scan(&value)

	byteToInt, _ := strconv.ParseInt(string(val), 10, 64)
	insertSQL := `INSERT INTO metrics VALUES ($1, $2, $3 ,0)
	ON CONFLICT (id,mtype) DO UPDATE SET delta =($3);`

	_, err := r.bd.Exec(insertSQL, name, "counter", byteToInt+value)
	if err != nil {
		fmt.Print("ошибка при вставке gauge ", name, " ", byteToInt, " ошибка :", err)
		return
	}

	fmt.Printf("Add in BDCounter %s, in val = %d \n", name, byteToInt)

}

func (r GaugeBDStorage) Get(name string) ([]byte, bool) {
	fmt.Print("Зашли в функцию Get Gauge")

	var value float64
	selectSQL := `
	SELECT val
	FROM metrics 
	WHERE id=$1 AND mtype=$2;
	`
	//fmt.Print("выполняем запрос Get Gauge :", selectSQL)
	r.bd.QueryRow(selectSQL, name, "gauge").Scan(&value)

	fmt.Print("получили значение метрики из БД  с типом Gauge  и менем ", name, " значение = ", value)

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

	fmt.Print("получили значение метрики из БД  с типом Counter  и менем ", name, " значение = ", value)

	if value != 0 {
		return []byte(strconv.FormatInt(value, 10)), true
	}

	return []byte(""), false // пустой список байт

}
