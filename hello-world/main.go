package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/Masterminds/log-go"
)

const out = `prd-demo`

var debug = flag.Bool("debug", false, "display debug output")

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Debug("function %s took %s", name, elapsed)
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer timeTrack(time.Now(), "handler")
	// log.Info("Printing message")
	log.Debug("Request %v", r)
	fmt.Fprint(w, out)
}

func main() {

	flag.Parse()
	logger := log.NewStandard()
	if *debug {
		logger.Level = log.DebugLevel
	}
	log.Current = logger

	log.Info("Starting Hello World")

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
