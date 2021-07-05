package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
)

type App struct {
	DB *DB
}

func (a *App) hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func (a *App) add(w http.ResponseWriter, req *http.Request) {
	if err := a.DB.createTable(); err != nil {
		fmt.Fprintf(w, "err: %s\n", err.Error())
	}

	if err := a.DB.addNumber(); err != nil {
		fmt.Fprintf(w, "err: %s\n", err.Error())
	}

	fmt.Fprintf(w, "ok\n")
}

func (a *App) get(w http.ResponseWriter, req *http.Request) {
	if err := a.DB.createTable(); err != nil {
		fmt.Fprintf(w, "err: %s\n", err.Error())
	}

	numbers, err := a.DB.getNumbers()
	if err != nil {
		fmt.Fprintf(w, "err: %s\n", err.Error())
	}

	fmt.Fprintf(w, "numbers: %v\n", numbers)
}

func main() {
	//db, err := New(os.Getenv("DSN"))
	//if err != nil {
	//	panic(err)
	//}

	app := App{DB: nil}

	http.HandleFunc("/add", app.add)
	http.HandleFunc("/get", app.get)
	http.HandleFunc("/", app.hello)

	http.ListenAndServe(":10000", nil)
}

type DB struct {
	*pgxpool.Pool
}

func New(dataSourceName string) (*DB, error) {
	var (
		poolConfig *pgxpool.Config
		db         *pgxpool.Pool
		err        error
	)

	if poolConfig, err = pgxpool.ParseConfig(dataSourceName); err != nil {
		return nil, err
	}

	if db, err = pgxpool.ConnectConfig(context.Background(), poolConfig); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) createTable() error {
	sql := `CREATE TABLE IF NOT EXISTS numbers (number integer);`

	if _, err := db.Exec(context.Background(), sql); err != nil {
		return err
	}

	return nil
}

func (db *DB) addNumber() error {
	sql := `INSERT INTO numbers (number) VALUES ($1);`

	if _, err := db.Exec(context.Background(), sql, 10); err != nil {
		return err
	}

	return nil
}

func (db *DB) getNumbers() ([]int, error) {
	sql := `SELECT number FROM numbers`

	rows, err := db.Query(context.Background(), sql)
	if err != nil {
		return nil, err
	}

	numbers := make([]int, 0, 100)

	for rows.Next() {
		var number int

		if err := rows.Scan(
			&number,
		); err != nil {
			return nil, err
		}

		numbers = append(numbers, number)
	}

	return numbers, nil
}
