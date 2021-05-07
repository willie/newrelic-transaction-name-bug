package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/willie/newrelic-transaction-name-bug/handlers"
)

const (
	programName = "bob"
	portNumber  = 8082
)

func main() {
	newRelicLicense, defined := os.LookupEnv("NEWRELIC_LICENSE")
	if !defined {
		log.Fatalln("NEWRELIC_LICENSE enviroment variable not found")
	}

	nrApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName(programName),
		newrelic.ConfigLicense(newRelicLicense),
	)
	if err != nil {
		log.Fatalln(err)
	}

	http.DefaultTransport = newrelic.NewRoundTripper(http.DefaultTransport)

	mux := chi.NewMux()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Use(handlers.NewRelic(nrApp))
	mux.Use(handlers.ChiNewRelic())

	mux.Get("/header/{vbk}", func(w http.ResponseWriter, r *http.Request) {
		vbk := chi.URLParam(r, "vbk")
		io.WriteString(w, vbk+" was requested.")
		io.WriteString(w, "\n")
	})

	fmt.Println(programName, "listening on port:", portNumber)
	err = http.ListenAndServe(fmt.Sprintf(":%d", portNumber), mux)
	if err != nil {
		log.Fatalln(err)
	}
}
