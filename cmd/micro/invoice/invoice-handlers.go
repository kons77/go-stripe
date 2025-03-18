package main

import (
	"fmt"
	"net/http"
	"time"
)

type Products struct {
	Name     string
	Amount   int
	Quantity int
}

type Order struct {
	ID        int       `json:"id"`
	Quantity  int       `json:"quantity"`
	Amount    int       `json:"amount"`
	Product   string    `json:"product"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	// Items     []Products `json:"items"`
}

// CreateAndSendInvoice creates an invoice as a PDF, and emails it to recipient
func (app *application) CreateAndSendInvoice(w http.ResponseWriter, r *http.Request) {
	// receive json
	var order Order

	err := app.readJSON(w, r, &order)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	// generate a pdf invoice
	err = app.createInvoicePDF(order)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	// create mail attachment
	attachments := []string{
		fmt.Sprintf("./invoices/%d.pdf", order.ID),
	}

	// send mail with attachment
	err = app.SendMail("info@widget.com", order.Email, "Your invoice", "invoice", attachments, nil)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	// send response
	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false
	resp.Message = fmt.Sprintf("Invoice %d.pdf created and sent to %s", order.ID, order.Email)
	app.writeJSON(w, http.StatusCreated, resp)

}
