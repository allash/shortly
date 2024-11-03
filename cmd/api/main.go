package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Port int
	Environment string
}

type Application struct {
	Config Config
	Logger *log.Logger	
}

func main() {
	var cfg Config

	flag.IntVar(&cfg.Port, "port", 3000, "Port")
	flag.StringVar(&cfg.Environment, "env", "dev", "Environment")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &Application{
		Config: cfg,
		Logger: logger,
	}

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Port),
		Handler: app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("Starting %s server on %s", cfg.Environment, srv.Addr)

	err := srv.ListenAndServe()
	logger.Fatal(err)
}
