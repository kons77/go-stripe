package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kons77/go-stripe/internal/driver"
	"github.com/kons77/go-stripe/internal/models"
)

const version = "1.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string // data source name
	}
	stripe struct {
		secret string
		key    string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
	}
	secretkey string // secret key to sign reset password link
	frontend  string // url to frontend
}

// Application struct (holds app configuration)
type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	version  string
	DB       models.DBModel
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

	app.infoLog.Printf("Starting Back end server in %s mode on port %d", app.config.env, app.config.port)

	return srv.ListenAndServe()
}

func main() {
	var cfg config

	// set flags
	flag.IntVar(&cfg.port, "port", 4001, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Applicetion environment {development|production|maintenance}")
	flag.StringVar(&cfg.db.dsn, "dsn", "", "Database DSN connection string")
	flag.StringVar(&cfg.smtp.host, "smtphost", "", "smtp host")
	flag.IntVar(&cfg.smtp.port, "smtpport", 587, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtpuser", "", "smtp username")
	flag.StringVar(&cfg.smtp.password, "smtppassword", "", "smtp password")
	//maybe AUTH and TLS
	//flag.StringVar(&cfg.secretkey, "secret", "", "secret key")
	flag.StringVar(&cfg.frontend, "frontend", "http://127.0.0.1:4000", "url to frontend")

	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Println("Cannot load .env file", err)
	}

	// read from .env
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")
	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	cfg.db.dsn = os.Getenv("DSN")
	cfg.smtp.host = os.Getenv("SMTP_HOST")
	cfg.smtp.username = os.Getenv("SMTP_USER")
	cfg.smtp.password = os.Getenv("SMTP_PASWORD")
	cfg.secretkey = os.Getenv("SIGN_KEY")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// create db connection
	conn, err := driver.OpenDB(cfg.db.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer conn.Close()

	// set app config values
	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		version:  version,
		DB:       models.DBModel{DB: conn},
	}

	err = app.serve()
	if err != nil {
		app.errorLog.Println(err)
		log.Fatal()
	}

}
