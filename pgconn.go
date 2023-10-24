package main

import (
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

type config struct {
	n, maxconn int
	sleep      time.Duration
}

func run(c config) (err error) {
	db.SetMaxOpenConns(c.maxconn)

	done := make(chan struct{})
	g := new(errgroup.Group)
	for i := 0; i < c.n; i++ {
		g.Go(func() error { return checkAndSleep(c.sleep) })
	}
	go func() {
		err = g.Wait()
		close(done)
	}()

	// wait a bit for connections and report the stats each 1s
	ticker := time.NewTimer(100 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			ticker.Reset(1 * time.Second)
			fmt.Printf("%+v\n", db.Stats())
		case <-done:
			ticker.Stop()
			return
		}
	}
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
