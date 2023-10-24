package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

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

func checkAndSleep(dt time.Duration) error {
	row := db.QueryRow("select 1")
	var x int
	err := row.Scan(&x)
	if err != nil {
		return err
	}
	if x != 1 {
		return fmt.Errorf("select 1 should return 1")
	}
	var v any
	row = db.QueryRow("select pg_sleep($1)", dt.Seconds())
	return row.Scan(&v)
}

func run(c config) error {
	db.SetMaxOpenConns(c.maxconn)

	g := new(errgroup.Group)
	for i := 0; i < c.n; i++ {
		g.Go(func() error { return checkAndSleep(c.sleep) })
	}

	fmt.Printf("db stats: %+v\n", db.Stats())

	return g.Wait()
}

type config struct {
	n, maxconn int
	sleep      time.Duration
}

func parseArgs(args []string) (c config, err error) {
	usage := `usage: pgconn n-conn max-conn sleep-time`

	if len(args) < 3 {
		return c, fmt.Errorf(usage)
	}
	c.n, err = strconv.Atoi(args[0])
	if err != nil {
		return c, err
	}
	c.maxconn, err = strconv.Atoi(args[1])
	if err != nil {
		return c, err
	}
	c.sleep, err = time.ParseDuration(args[2])
	return c, err
}

func main() {
	c, err := parseArgs(os.Args[1:])
	if err != nil {
		die(2, err)
	}

	err = run(c)
	if err != nil {
		die(1, err)
	}
}

func die(exitcode int, err error) {
	if strings.HasPrefix(err.Error(), "usage:") {
		fmt.Fprintln(os.Stderr, err)
	} else {
		fmt.Fprintln(os.Stderr, "pgconn:", err)
	}
	os.Exit(exitcode)
}
