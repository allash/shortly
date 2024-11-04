package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"shorturl.allash.com/internal/data"
)

type Config struct {
	port int
	environment string
	db struct {
		dsn string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime	 string
	}
}

type Application struct {
	config Config
	logger *log.Logger
	models data.Models
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file!")
	}

	var cfg Config

	flag.IntVar(&cfg.port, "port", 3000, "Port")
	flag.StringVar(&cfg.environment, "env", "dev", "Environment")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DB_CONNECTION_STRING"), "Postgres DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 30, "Postgres max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 30, "Postgres max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "Postgres max connection idle time")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	app := &Application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.port),
		Handler: app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("Starting %s server on %s", cfg.environment, srv.Addr)

	err = srv.ListenAndServe()
	logger.Fatal(err)
}

func openDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
