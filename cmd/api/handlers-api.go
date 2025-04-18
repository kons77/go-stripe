package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kons77/go-stripe/internal/cards"
	"github.com/kons77/go-stripe/internal/encryption"
	"github.com/kons77/go-stripe/internal/models"
	"github.com/kons77/go-stripe/internal/urlsigner"
	"github.com/kons77/go-stripe/internal/validator"
	"github.com/stripe/stripe-go/v81"
	"golang.org/x/crypto/bcrypt"
)

type stripePayload struct {
	Currency      string `json:"currency"`
	Amount        string `json:"amount"`
	PaymentMethod string `json:"payment_method"`
	Email         string `json:"email"`
	CardBrand     string `json:"card_brand"`
	ExpiryMonth   int    `json:"exp_month"`
	ExpiryYear    int    `json:"exp_year"`
	LastFour      string `json:"last_four"`
	Plan          string `json:"plan"`
	ProductID     string `json:"product_id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
}

type jsonResponce struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"` //omitempty means that if the field is empty (""), it will not be included in the JSON.
	Content string `json:"content,omitempty"`
	ID      int    `json:"id,omitempty"`
}

// GetPaymentIntent gets a payment intent, and returns json (or error json)
func (app *application) GetPaymentIntent(w http.ResponseWriter, r *http.Request) {
	var payload stripePayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.errorLog.Println(err)
		// We should be sending back a properly formatted Json error message, but we'll take care of that later.
		return
	}

	amount, err := strconv.Atoi(payload.Amount)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	card := cards.Card{
		Secret:   app.config.stripe.secret,
		Key:      app.config.stripe.key,
		Currency: payload.Currency,
	}

	okay := true
	pi, msg, err := card.Charge(payload.Currency, amount)
	if err != nil {
		okay = false
	}

	if okay {
		out, err := json.MarshalIndent(pi, "", "   ")
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	} else {
		j := jsonResponce{
			OK:      false,
			Message: msg,
			Content: "",
		}

		out, err := json.MarshalIndent(j, "", "   ")
		if err != nil {
			app.errorLog.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	}
}

// GetWidgetByID gets one widget by id and returns as JSON
func (app *application) GetWidgetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	widgetID, _ := strconv.Atoi(id)

	widget, err := app.DB.GetWidget(widgetID)
	if err != nil {
		app.errorLog.Println(err) //we'll fix that later on and we'll write this out as Json.
		return
	}

	out, err := json.MarshalIndent(widget, "", "  ") // you shouldn't use Indent in production but for development it's fine.
	if err != nil {
		app.errorLog.Println(err) //we'll fix that later on and we'll write this out as Json.
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Invoice describes the JSON payload sent to the microservice
type Invoice struct {
	ID        int       `json:"id"`
	Quantity  int       `json:"quantity"`
	Amount    int       `json:"amount"`
	Product   string    `json:"product"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// callInvoiceMicro calls the invoicing microservice
func (app *application) callInvoiceMicro(inv Invoice) {
	go func() {
		url := "http://127.0.0.1:5000/invoice/create-and-send"
		out, err := json.MarshalIndent(inv, "", "\t")
		if err != nil {
			app.errorLog.Println("Error marshaling invoice:", err)
			return
		}

		// Create a client with a timeout
		client := &http.Client{
			Timeout: 5 * time.Second,
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(out))
		if err != nil {
			app.errorLog.Println("Error creating request to invoice microservice:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			app.errorLog.Println("Error calling invoice microservice:", err)
			return
		}
		defer resp.Body.Close()

		app.infoLog.Println("Invoice microservice called successfully")

	}()
}

// CreateCustomerAndSubscribeToPlan is the handler for subscribing to the bronze plan
func (app *application) CreateCustomerAndSubscribeToPlan(w http.ResponseWriter, r *http.Request) {
	var data stripePayload
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	// validate data
	v := validator.New()
	v.Check(len(data.FirstName) > 1, "first_name", "must be at least 2 characters")

	if !v.Valid() {
		app.failedValidation(w, r, v.Errors)
		return
	}

	card := cards.Card{
		Secret:   app.config.stripe.secret,
		Key:      app.config.stripe.key,
		Currency: data.Currency,
	}

	okay := true
	var subscription *stripe.Subscription
	txnMsg := "Transcation successful"

	stripeCustomer, msg, err := card.CreateCustomer(data.PaymentMethod, data.Email)
	if err != nil {
		app.errorLog.Println(err)
		okay = false
		txnMsg = msg
	}

	if okay {
		subscription, err = card.SubscribeToPlan(stripeCustomer, data.Plan, data.Email, data.LastFour, "")
		if err != nil {
			app.errorLog.Println(err)
			okay = false
			txnMsg = "Error subscribing customer"
		}

		app.infoLog.Println("subscription id is", subscription.ID)
	}

	var ordID int

	if okay {
		// create a new customer - later need to check if customer exists then login
		productID, _ := strconv.Atoi(data.ProductID)
		customerID, err := app.SaveCustomer(data.FirstName, data.LastName, data.Email)
		if err != nil {
			app.errorLog.Println(err)
			return
		}

		// create a new transaction
		amount, _ := strconv.Atoi(data.Amount)
		txn := models.Transaction{
			Amount:              amount,
			Currency:            "usd",
			LastFour:            data.LastFour,
			ExpiryMonth:         data.ExpiryMonth,
			ExpiryYear:          data.ExpiryYear,
			TransactionStatusID: 2,
			PaymentIntent:       subscription.ID,
			PaymentMethod:       data.PaymentMethod,
		}

		txnID, err := app.SaveTransaction(txn)
		if err != nil {
			app.errorLog.Println(err)
			return
		}

		// create order
		order := models.Order{
			WidgetID:      productID,
			TransactionID: txnID,
			CustomerID:    customerID,
			StatusID:      1,
			Quantity:      1,
			Amount:        amount,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		ordID, err = app.SaveOrder(order)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
	}

	// call microservice
	inv := Invoice{
		ID:        ordID,
		Amount:    2000,
		Product:   "Bronze Plan monthly subscription",
		Quantity:  1,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		CreatedAt: time.Now(),
	}

	app.callInvoiceMicro(inv)

	resp := jsonResponce{
		OK:      okay,
		Message: txnMsg,
	}

	out, err := json.MarshalIndent(resp, "", "   ")
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// SaveCustomer saves a customer and returns id - may be  different for frontend and for backend
func (app *application) SaveCustomer(firstName, lastName, email string) (int, error) {
	customer := models.Customer{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	id, err := app.DB.InsertCustomer(customer)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// SaveTransaction saves a transaction and returns id
func (app *application) SaveTransaction(txn models.Transaction) (int, error) {
	id, err := app.DB.InsertTransaction(txn)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// SaveOrder saves a order and returns id
func (app *application) SaveOrder(order models.Order) (int, error) {
	id, err := app.DB.InsertOrder(order)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// CreateAuthToken creates and sends an auth token, if user supplies valid information
func (app *application) CreateAuthToken(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &userInput)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	// get the user from the database by email; send error if invalid mail
	user, err := app.DB.GetUserByEmail(userInput.Email)
	if err != nil {
		app.invalidCredentials(w)
		return
	}

	// validate the password; send error if invalid password
	validPassword, err := app.passwordMatches(user.Password, userInput.Password)
	if err != nil {
		app.invalidCredentials(w)
		return
	}

	if !validPassword {
		app.invalidCredentials(w)
		return
	}

	// generate the token
	token, err := models.GenerateToken(user.ID, 24*time.Hour, models.ScopeAuthentication)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	//save to db
	err = app.DB.InsertToken(token, user)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	// send response
	var payload struct {
		Error   bool          `json:"error"`
		Message string        `json:"message"`
		Token   *models.Token `json:"authentication_token"`
	}

	payload.Error = false
	payload.Message = fmt.Sprintf("token for %s created", userInput.Email)
	payload.Token = token

	_ = app.writeJSON(w, http.StatusOK, payload)
}

// authenticateToken checks an auth token for validity
func (app *application) authenticateToken(r *http.Request) (*models.User, error) {

	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return nil, errors.New("no authorization header received")
	}

	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("no authorization header received") // don't want to end user a lot about auth
	}

	token := headerParts[1]
	if len(token) != 26 {
		return nil, errors.New("authentication token wrong size")
	}

	// get the user from the tokens table
	user, err := app.DB.GetUserForToken(token)
	if err != nil {
		return nil, errors.New("no matchin user found")
	}

	return user, nil
}

// CheckAuthentication checks auth status
func (app *application) CheckAuthentication(w http.ResponseWriter, r *http.Request) {
	// validate the token, and get associated user
	user, err := app.authenticateToken(r)
	if err != nil {
		app.invalidCredentials(w)
		return
	}

	// valid user
	payload := struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}{
		Error:   false,
		Message: fmt.Sprintf("authenticated user %s", user.Email),
	}

	app.writeJSON(w, http.StatusOK, payload)
}

// VirtualTerminalPaymentSucceeded displays a page with receipt information
func (app *application) VirtualTerminalPaymentSucceeded(w http.ResponseWriter, r *http.Request) {
	//  send JSON as the body of the post request from the front end
	var txnData struct {
		PaymentAmount   int    `json:"amount"`
		PaymentCurrency string `json:"currency"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		Email           string `json:"email"`
		PaymentIntent   string `json:"payment_intent"`
		PaymentMethod   string `json:"payment_method"`
		BankReturnCode  string `json:"bank_return_code"`
		ExpiryMonth     int    `json:"expiry_month"`
		ExpiryYear      int    `json:"expiry_year"`
		LastFour        string `json:"last_four"`
	}

	// if you want to extract these two shorter functions or methods, you feel free to do so.

	err := app.readJSON(w, r, &txnData)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	card := cards.Card{
		Secret: app.config.stripe.secret,
		Key:    app.config.stripe.key,
	}

	pi, err := card.RetrievePaymentIntent(txnData.PaymentIntent)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	pm, err := card.GetPaymentMethod(txnData.PaymentMethod)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	txnData.LastFour = pm.Card.Last4
	txnData.ExpiryMonth = int(pm.Card.ExpMonth)
	txnData.ExpiryYear = int(pm.Card.ExpYear)

	txn := models.Transaction{
		Amount:              txnData.PaymentAmount,
		Currency:            txnData.PaymentCurrency,
		LastFour:            txnData.LastFour,
		ExpiryMonth:         txnData.ExpiryMonth,
		ExpiryYear:          txnData.ExpiryYear,
		PaymentIntent:       txnData.PaymentIntent,
		PaymentMethod:       txnData.PaymentMethod,
		BankReturnCode:      pi.LatestCharge.ID,
		TransactionStatusID: 2,
	}

	_, err = app.SaveTransaction(txn)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, txn)

}

// SendPasswordResetEmail sends an email with a signed url to allow user to reset password
func (app *application) SendPasswordResetEmail(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email string `json:"email"`
	}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	// verify that email exists
	_, err = app.DB.GetUserByEmail(payload.Email)
	if err != nil {
		var resp struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}
		resp.Error = true
		resp.Message = "No matching email found on our system"
		app.writeJSON(w, http.StatusAccepted, resp)
		return
	}

	link := fmt.Sprintf("%s/reset-password?email=%s", app.config.frontend, payload.Email)

	sign := urlsigner.Signer{
		Secret: []byte(app.config.secretkey),
	}

	signedLink := sign.GenerateTokenFromString(link)

	var data struct {
		Link string
	}

	data.Link = signedLink

	// send mail
	err = app.SendMail("info@widgets.com", payload.Email, "Password Reset Request", "password-reset", data)
	if err != nil {
		app.errorLog.Println(err)
		app.badRequest(w, r, err)
		return
	}

	resp := struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}{
		Error:   false,
		Message: "",
	}

	app.writeJSON(w, http.StatusCreated, resp)
}

// ResetPassword resets a user's password in the database
func (app *application) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	ecryptor := encryption.Encryption{
		Key: []byte(app.config.secretkey),
	}

	realEmail, err := ecryptor.Decrypt(payload.Email)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	user, err := app.DB.GetUserByEmail(realEmail)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 12)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	err = app.DB.UpdatePasswordForUser(user, string(newHash))
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	resp := struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}{
		Error:   false,
		Message: "password changed",
	}

	app.writeJSON(w, http.StatusCreated, resp)

}

// AllSales returns all sales or all subscriptions as a slice
func (app *application) AllSalesSubscription(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		PageSize    int `json:"page_size"`
		CurrentPage int `json:"page"`
		IsReccuring int `json:"is_recur"`
	}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	isReccuring := payload.IsReccuring // is it sale (0) or subscription (1)

	allSales, lastPage, totalRecords, err := app.DB.GetAllOrdersPaginated(isReccuring, payload.PageSize, payload.CurrentPage)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	var resp struct {
		CurrentPage  int             `json:"current_page"`
		PageSize     int             `json:"page_size"`
		LastPage     int             `json:"last_page"`
		TotalRecords int             `json:"total_records"`
		Orders       []*models.Order `json:"orders"`
	}

	resp.CurrentPage = payload.CurrentPage
	resp.PageSize = payload.PageSize
	resp.LastPage = lastPage
	resp.TotalRecords = totalRecords
	resp.Orders = allSales

	app.writeJSON(w, http.StatusOK, resp)
}

// GetSale returns one sale as json, by id
func (app *application) GetSale(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	orderID, _ := strconv.Atoi(id)

	order, err := app.DB.GetOrderByID(orderID)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, order)
}

// RefundCharge accepts a json payload and tries to refund a charge
func (app *application) RefundCharge(w http.ResponseWriter, r *http.Request) {
	var chargeToRefund struct {
		ID            int    `json:"id"`
		PaymentIntent string `json:"pi"`
		Amount        int    `json:"amount"`
		Currency      string `json:"currency"`
	}

	err := app.readJSON(w, r, &chargeToRefund)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	// validate
	/*Validate the refund amount against the order by querying the database for the order ID, retrieving the stored amount,
	and comparing it with the amount in the JSON payload. Ensure the refund isn’t greater than the original amount,
	or apply any other logic needed. This is a simple check you can implement on your own.*/

	card := cards.Card{
		Secret:   app.config.stripe.secret,
		Key:      app.config.stripe.key,
		Currency: chargeToRefund.Currency,
	}

	err = card.Refund(chargeToRefund.PaymentIntent, chargeToRefund.Amount)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	//update status in db
	err = app.DB.UpdateOrderStatus(chargeToRefund.ID, 2)
	if err != nil {
		app.badRequest(w, r, errors.New("the charge was refunded, but the database could not be updated"))
		return
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false
	resp.Message = "Charge refunded"

	app.writeJSON(w, http.StatusOK, resp)

}

// CancelSubscription is the handler to cancel a subscription
func (app *application) CancelSubscription(w http.ResponseWriter, r *http.Request) {
	var subToCancel struct {
		ID            int    `json:"id"`
		PaymentIntent string `json:"pi"` //subscritpion id in fact
		Currency      string `json:"currency"`
	}

	err := app.readJSON(w, r, &subToCancel)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	card := cards.Card{
		Secret:   app.config.stripe.secret,
		Key:      app.config.stripe.key,
		Currency: subToCancel.Currency,
	}

	err = card.CancelSubscription(subToCancel.PaymentIntent)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	//update status in db
	err = app.DB.UpdateOrderStatus(subToCancel.ID, 3)
	if err != nil {
		app.badRequest(w, r, errors.New("the subscription is canceled, but the database could not be updated"))
		return
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false
	resp.Message = "Subscription canceled"

	app.writeJSON(w, http.StatusOK, resp)

}

// AllUsers returns a JSON file listing all admin users
func (app *application) AllUsers(w http.ResponseWriter, r *http.Request) {
	allUsers, err := app.DB.GetAllUsers()
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, allUsers)
}

// OneUser gets one user by id (from the url) and returns it as JSON
func (app *application) OneUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	user, err := app.DB.GetOneUser(userID)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, user)
}

// EditUser is the handler for adding or editing an existing user
func (app *application) EditUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	var user models.User

	err = app.readJSON(w, r, &user)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	// add or update user
	if userID > 0 {
		err = app.DB.EditUser(user)
		if err != nil {
			app.badRequest(w, r, err)
			return
		}

		// change password if not empty
		if user.Password != "" {
			newHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
			if err != nil {
				app.badRequest(w, r, err)
				return
			}

			err = app.DB.UpdatePasswordForUser(user, string(newHash))
			if err != nil {
				app.badRequest(w, r, err)
				return
			}
		}
	} else {
		// editing user
		newHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
		if err != nil {
			app.badRequest(w, r, err)
			return
		}
		err = app.DB.AddUser(user, string(newHash))
		if err != nil {
			app.badRequest(w, r, err)
			return
		}
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false
	app.writeJSON(w, http.StatusOK, resp)
}

// DeleteUser deletes a user, and all associated tokens, from the database
func (app *application) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	err = app.DB.DeleteUser(userID)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false
	app.writeJSON(w, http.StatusOK, resp)
}
