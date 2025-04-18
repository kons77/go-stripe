package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/joho/godotenv"
	"github.com/kons77/go-stripe/internal/driver"
	"github.com/kons77/go-stripe/internal/models"
)

const version = "1.0"

// const cssVersion = "1"

var session *scs.SessionManager

type config struct {
	port int
	env  string
	api  string
	db   struct {
		dsn string // data source name
	}
	stripe struct {
		secret string
		key    string
	}
	secretkey string // secret key to sign reset password link
	frontend  string // url to frontend
}

// Application struct (holds app configuration)
type application struct {
	config        config
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCahce map[string]*template.Template // html/template package, not text/template
	version       string
	DB            models.DBModel
	Session       *scs.SessionManager
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

	app.infoLog.Printf("Starting HTTP server in %s mode on port %d", app.config.env, app.config.port)
	//app.infoLog.Println(app.config.db.dsn)

	return srv.ListenAndServe()
}

func main() {
	gob.Register(TransactionData{})

	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Applicetion environment {development|production}")
	flag.StringVar(&cfg.api, "api", "http://127.0.0.1:4001", "URL to api")
	flag.StringVar(&cfg.db.dsn, "dsn", "", "Database DSN connection string")
	//flag.StringVar(&cfg.secretkey, "secret", "", "secret key")
	flag.StringVar(&cfg.frontend, "frontend", "http://127.0.0.1:4000", "url to frontend")

	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Println("Cannot load .env file", err)
	}

	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")
	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	cfg.db.dsn = os.Getenv("DSN")
	cfg.secretkey = os.Getenv("SIGN_KEY") // must be exactly 32 characters long

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// make database connection
	conn, err := driver.OpenDB(cfg.db.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer conn.Close()

	// set up session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Store = mysqlstore.New(conn)

	tc := make(map[string]*template.Template)

	app := &application{
		config:        cfg,
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCahce: tc,
		version:       version,
		DB:            models.DBModel{DB: conn},
		Session:       session,
	}

	go app.ListenToWsChannel()

	err = app.serve()
	if err != nil {
		app.errorLog.Println(err)
		log.Fatal()
	}
}
