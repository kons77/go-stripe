package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const version = "1.0"

type config struct {
	port int
	env  string
	smtp struct {
		host     string
		port     int
		username string
		password string
	}
	frontend string // url to frontend
}

// Application struct (holds app configuration)
type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	version  string
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.infoLog.Printf("Starting invoice microservice on port %d", app.config.port)

	return srv.ListenAndServe()
}

func main() {

	var cfg config

	// set flags
	flag.IntVar(&cfg.port, "port", 5000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Applicetion environment {development|production|maintenance}")

	flag.StringVar(&cfg.smtp.host, "smtphost", "", "smtp host")
	flag.IntVar(&cfg.smtp.port, "smtpport", 587, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtpuser", "", "smtp username")
	flag.StringVar(&cfg.smtp.password, "smtppassword", "", "smtp password")
	flag.StringVar(&cfg.frontend, "frontend", "http://127.0.0.1:4000", "url to frontend")

	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Println("Cannot load .env file", err)
	}

	cfg.smtp.host = os.Getenv("SMTP_HOST")
	cfg.smtp.username = os.Getenv("SMTP_USER")
	cfg.smtp.password = os.Getenv("SMTP_PASWORD")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// set app config values
	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		version:  version,
	}

	err = app.serve()
	if err != nil {
		app.errorLog.Println(err)
		log.Fatal()
	}
}
