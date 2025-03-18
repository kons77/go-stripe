package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/phpdave11/gofpdf"
	"github.com/phpdave11/gofpdf/contrib/gofpdi"
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

	app.CreateDirInNotExist("./invoices")

	err = app.serve()
	if err != nil {
		app.errorLog.Println(err)
		log.Fatal()
	}

}

// createInvoicePDF generates a PDF version of the invoice
func (app *application) createInvoicePDF(order Order) error {
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.SetMargins(10, 13, 10)
	pdf.SetAutoPageBreak(true, 0)

	importer := gofpdi.NewImporter()

	t := importer.ImportPage(pdf, "./pdf-templates/invoice.pdf", 1, "/MediaBox")

	pdf.AddPage()
	importer.UseImportedTemplate(pdf, t, 0, 0, 215.9, 0)

	// write info
	pdf.SetY(50)
	pdf.SetX(10)
	pdf.SetFont("Times", "", 11)
	pdf.CellFormat(97, 8, fmt.Sprintf("Attention: %s %s", order.FirstName, order.LastName), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(97, 8, order.Email, "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(97, 8, order.CreatedAt.Format("2006-01-02"), "", 0, "L", false, 0, "")

	// range throug a slice of products

	pdf.SetX(58)
	pdf.SetY(93)
	pdf.CellFormat(155, 8, order.Product, "", 0, "L", false, 0, "")
	pdf.SetX(166)
	pdf.CellFormat(20, 8, fmt.Sprintf("%d", order.Quantity), "", 0, "C", false, 0, "")
	pdf.SetX(185)
	pdf.CellFormat(20, 8, fmt.Sprintf("$%.2f", float32(order.Amount/100.0)), "", 0, "R", false, 0, "")

	/*
		The next thing we do, my PDF is done for my purposes and
		I could leave the bottom part where you have really not difficult.
	*/

	invoicePath := fmt.Sprintf("./invoices/%d.pdf", order.ID)
	err := pdf.OutputFileAndClose(invoicePath)
	if err != nil {
		return err
	}

	return nil

}
