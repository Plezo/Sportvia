package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/Plezo/Sportvia/internal/data"
	"github.com/Plezo/Sportvia/internal/jsonlog"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type config struct {
	bind string
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	templates *template.Template
	scrapers data.Scrapers
}

func main() {

	var cfg config
	// os.Setenv("SPORTVIA_DB_DSN", "postgres://postgres:password@localhost:5432/sportvia?sslmode=disable")

	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.bind, "bind", "localhost:8080", "Server bind address")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Parse()

	if cfg.env == "development" {
		err := godotenv.Load(".env")
		if err != nil {
			fmt.Println("Error loading .env file")
			os.Exit(1)
		}
	}

	if _, ok := os.LookupEnv("SPORTVIA_DB_DSN"); ok {
		cfg.db.dsn = os.Getenv("SPORTVIA_DB_DSN")
	} else {
		fmt.Println("SPORTVIA_DB_DSN not set")
		os.Exit(1)
	}

	// flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("SPORTVIA_DB_DSN"), "PostgreSQL DSN")

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		templates: template.Must(template.ParseGlob("ui/html/*.html")),
		scrapers: data.NewScrapers(),
	}

	srv := &http.Server{
		Addr:         cfg.bind,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.PrintInfo("starting %s server on %s", map[string]string{
		"addr": srv.Addr,
		"env": cfg.env,
	})

	err = srv.ListenAndServe()
	logger.PrintFatal(err, nil)
}

func openDB(cfg config) (*sql.DB, error) {
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