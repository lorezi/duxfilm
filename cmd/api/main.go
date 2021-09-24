package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Declare a string containing the application version number.
const version = "1.0.0"

// Define a config struct to hold all the configuration settings for out application

type config struct {
	port int
	env  string
}

// Define an application struct to build the dependencies for our HTTP handlers, helpers, and middleware.
type application struct {
	config config
	logger *log.Logger
}

func main() {
	// Declare an instance of the config struct
	var cfg config

	// Read the value of the port and env command-line flags into the config struct
	// default port 3000
	// default environment "development"

	flag.IntVar(&cfg.port, "port", 3000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	// Initialize a new logger which writes messages to the standard out stream prefixed with current date and time.
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Declare an instance of the application struct, containing the config struct and the logger
	app := &application{
		config: cfg,
		logger: logger,
	}

	// Declare a new servemux and add a /v1/healthcheck route which dispatches requests to the healthcheckHandler method
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)

	// Declare a HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start the HTTP server
	logger.Printf("starting %s server on %s 😀😀😀", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}
