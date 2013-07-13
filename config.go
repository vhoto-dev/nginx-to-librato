package main

import (
	"fmt"
	"github.com/pearkes/goconfig/config"
	"log"
	"os"
	"time"
)

type conf struct {
	libUser          string        // The librato user
	libToken         string        // The librato token
	libSource        string        // The source of the metric reported
	url              string        // The url to the nginx_status page
	flushInterval    time.Duration // The interval to flush to librato
	rawFlushInterval string        // The raw interval specified by the user
}

func NewConf(path string) *conf {
	con := &conf{}

	// We may have multiple errors here
	errs := make([]error, 0)

	c, err := config.ReadDefault(path)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read config file: \n%s\n", err)
		os.Exit(1)
	}

	token, err := c.String("settings", "token")
	if err != nil {
		errs = append(errs, fmt.Errorf("Failed parsing token: %s", err))
	}
	con.libToken = token

	con.libUser, err = c.String("settings", "user")
	if err != nil {
		errs = append(errs, fmt.Errorf("Failed parsing user: %s", err))
	}

	con.libSource, err = c.String("settings", "source")
	if err != nil {
		errs = append(errs, fmt.Errorf("Failed parsing source: %s", err))
	}

	con.url, err = c.String("settings", "url")
	if err != nil {
		errs = append(errs, fmt.Errorf("Failed parsing url: %s", err))
	}

	con.rawFlushInterval, err = c.String("settings", "flush_interval")
	if err != nil {
		errs = append(errs, fmt.Errorf("Failed parsing flush_interval: %s", err))
	}

	con.flushInterval, err = time.ParseDuration(con.rawFlushInterval)
	if err != nil {
		errs = append(errs, fmt.Errorf("Failed parsing flush_interval: %s", err))
	}

	// Print the errors and exit if there are any
	if len(errs) > 0 {
		fmt.Fprintf(os.Stderr, "Configuration errors:\n")
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "* %s\n", e.Error())
		}
		os.Exit(1)
	}

	log.Printf("Loaded configuration: %v", con)

	return con
}
