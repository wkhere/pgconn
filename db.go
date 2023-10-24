package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", dsn(os.Getenv("LOGNAME"), os.Getenv("DB")))
	if err != nil {
		panic(err)
	}
}

func dsn(user, dbname string) string {
	return fmt.Sprintf("postgres://%s@localhost/%s?sslmode=disable", user, dbname)
}
