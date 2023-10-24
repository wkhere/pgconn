package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

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
