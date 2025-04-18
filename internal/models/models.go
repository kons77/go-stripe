package models

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// DBModel is the type for database connection values
type DBModel struct {
	DB *sql.DB
}

// Models is the wrapper for all models
type Models struct {
	DB DBModel
}

// NewModels returns a model type with database connection pool
func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

// Widget is the type for all widgets
type Widget struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	InventoryLevel int       `json:"inventory_level"`
	Price          int       `json:"price"`
	Image          string    `json:"image"`
	IsReccuring    bool      `json:"is_reccuring"`
	PlanID         string    `json:"plan_id"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

// Order is the type for all orders
type Order struct {
	ID            int         `json:"id"`
	WidgetID      int         `json:"widget_id"`
	TransactionID int         `json:"transaction_id"`
	CustomerID    int         `json:"customer_id"`
	StatusID      int         `json:"status_id"`
	Quantity      int         `json:"quantity"`
	Amount        int         `json:"amount"`
	CreatedAt     time.Time   `json:"-"`
	UpdatedAt     time.Time   `json:"-"`
	Widget        Widget      `json:"widget"`
	Transaction   Transaction `json:"transaction"`
	Customer      Customer    `json:"customer"`
}

// Status is the type for statuses
type Status struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// TransactionStatus is the type for transaction statuses
type TransactionStatus struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// Transaction is the type for transactions
type Transaction struct {
	ID                  int       `json:"id"`
	Amount              int       `json:"amount"`
	Currency            string    `json:"currency"`
	LastFour            string    `json:"last_four"`
	ExpiryMonth         int       `json:"expiry_month"`
	ExpiryYear          int       `json:"expiry_year"`
	PaymentIntent       string    `json:"payment_intent"`
	PaymentMethod       string    `json:"payment_method"`
	BankReturnCode      string    `json:"bank_return_code"`
	TransactionStatusID int       `json:"transaction_status_id"`
	CreatedAt           time.Time `json:"-"`
	UpdatedAt           time.Time `json:"-"`
}

// Users is the type for all users
type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// Customer is the type for all customers
type Customer struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (m *DBModel) GetWidget(id int) (Widget, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var widget Widget

	// Both MySql and MariaDB required ?
	q := `
		select id, name, description, inventory_level, price, is_recurring, 
			plan_id, created_at, updated_at, coalesce(image, '')
		from widgets where id = ?
	`
	row := m.DB.QueryRowContext(ctx, q, id)
	err := row.Scan(
		&widget.ID,
		&widget.Name,
		&widget.Description,
		&widget.InventoryLevel,
		&widget.Price,
		&widget.IsReccuring,
		&widget.PlanID,
		&widget.CreatedAt,
		&widget.UpdatedAt,
		&widget.Image,
	)
	if err != nil {
		return widget, err
	}

	return widget, nil
}

// InsertTransaction inserts a new transaction, and returns its id
func (m *DBModel) InsertTransaction(txn Transaction) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		INSERT INTO transactions
			(amount, currency, last_four, expiry_month, expiry_year, payment_intent, payment_method,
			bank_return_code, transaction_status_id, created_at, updated_at)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := m.DB.ExecContext(ctx, stmt,
		txn.Amount,
		txn.Currency,
		txn.LastFour,
		txn.ExpiryMonth,
		txn.ExpiryYear,
		txn.PaymentIntent,
		txn.PaymentMethod,
		txn.BankReturnCode,
		txn.TransactionStatusID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// InsertOrder inserts a new order, and returns its id
func (m *DBModel) InsertOrder(order Order) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		INSERT INTO orders
			(widget_id, transaction_id, status_id, customer_id, quantity, amount, created_at, updated_at)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := m.DB.ExecContext(ctx, stmt,
		order.WidgetID,
		order.TransactionID,
		order.StatusID,
		order.CustomerID,
		order.Quantity,
		order.Amount,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// InsertCustomer inserts a new customer, and returns its id
func (m *DBModel) InsertCustomer(customer Customer) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		INSERT INTO customers
			(first_name, last_name, email, created_at, updated_at)
		VALUES(?, ?, ?, ?, ?)
	`

	result, err := m.DB.ExecContext(ctx, stmt,
		customer.FirstName,
		customer.LastName,
		customer.Email,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// GetUserByEmail gets a user by email address
func (m *DBModel) GetUserByEmail(email string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	email = strings.ToLower(email)
	var u User

	q := `
		select id, first_name, last_name, email, password, created_at, updated_at 
		from users 
		where email = ?
	`

	row := m.DB.QueryRowContext(ctx, q, email)

	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return u, err
	}

	return u, nil
}

// Authenticate
func (m *DBModel) Authenticate(email, password string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	q := `select id, password from users where email = ?`

	row := m.DB.QueryRowContext(ctx, q, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, errors.New("incorrect password")
	} else if err != nil {
		return 0, err
	}

	return id, nil
}

// UpdatePasswordForUser updates hash for user
func (m *DBModel) UpdatePasswordForUser(u User, hash string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `update users set password = ? where id = ?`

	_, err := m.DB.ExecContext(ctx, stmt, hash, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// GetAllOrders returns all orders if isReccuring = 0 or all subscriptions if isReccuring = 1
func (m *DBModel) GetAllOrders(isReccuring int) ([]*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var orders []*Order

	query := `
	select 
		o.id, o.widget_id, o.transaction_id, o.customer_id, o.status_id, o.quantity, o.amount, o.created_at, o.updated_at,
		w.id, w.name, 
		t.id, t.amount, t.currency, t.last_four, t.expiry_month, t.expiry_year, t.payment_intent, t.payment_method, t.bank_return_code,
		c.id, c.first_name, c.last_name, c.email
	from 
		orders o
		left join widgets w on (o.widget_id = w.id )
		left join transactions t on (o.transaction_id = t.id)
		left join customers c on (o.customer_id = c.id)
	where 
		w.is_recurring = ?	
	order by 
		o.created_at desc
	`

	rows, err := m.DB.QueryContext(ctx, query, isReccuring)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var o Order
		err = rows.Scan(
			&o.ID, &o.WidgetID, &o.TransactionID, &o.CustomerID,
			&o.StatusID, &o.Quantity, &o.Amount, &o.CreatedAt, &o.UpdatedAt,
			&o.Widget.ID, &o.Widget.Name,
			&o.Transaction.ID, &o.Transaction.Amount, &o.Transaction.Currency,
			&o.Transaction.LastFour, &o.Transaction.ExpiryMonth, &o.Transaction.ExpiryYear,
			&o.Transaction.PaymentIntent, &o.Transaction.PaymentMethod, &o.Transaction.BankReturnCode,
			&o.Customer.ID, &o.Customer.FirstName, &o.Customer.LastName, &o.Customer.Email,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}

	return orders, nil
}

// GetAllOrders returns a slice of subset of  orders if isReccuring = 0 or subscriptions if isReccuring = 1

func (m *DBModel) GetAllOrdersPaginated(isReccuring, pageSize, page int) ([]*Order, int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	offset := (page - 1) * pageSize
	// offset how far from the beginning of the results we should offset what we're looking for
	// pageSize tells us how much to limit our query for

	var orders []*Order

	query := `
	select 
		o.id, o.widget_id, o.transaction_id, o.customer_id, o.status_id, o.quantity, o.amount, o.created_at, o.updated_at,
		w.id, w.name, 
		t.id, t.amount, t.currency, t.last_four, t.expiry_month, t.expiry_year, t.payment_intent, t.payment_method, t.bank_return_code,
		c.id, c.first_name, c.last_name, c.email
	from 
		orders o
		left join widgets w on (o.widget_id = w.id )
		left join transactions t on (o.transaction_id = t.id)
		left join customers c on (o.customer_id = c.id)
	where 
		w.is_recurring = ?	
	order by 
		o.created_at desc
	limit ? offset ?
	`

	rows, err := m.DB.QueryContext(ctx, query, isReccuring, pageSize, offset)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var o Order
		err = rows.Scan(
			&o.ID, &o.WidgetID, &o.TransactionID, &o.CustomerID,
			&o.StatusID, &o.Quantity, &o.Amount, &o.CreatedAt, &o.UpdatedAt,
			&o.Widget.ID, &o.Widget.Name,
			&o.Transaction.ID, &o.Transaction.Amount, &o.Transaction.Currency,
			&o.Transaction.LastFour, &o.Transaction.ExpiryMonth, &o.Transaction.ExpiryYear,
			&o.Transaction.PaymentIntent, &o.Transaction.PaymentMethod, &o.Transaction.BankReturnCode,
			&o.Customer.ID, &o.Customer.FirstName, &o.Customer.LastName, &o.Customer.Email,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		orders = append(orders, &o)
	}

	query = `
		select count(o.id) 
		from orders o
		left join widgets w on (o.widget_id = w.id)
		where w.is_recurring = ?
	`

	var totalRecords int
	countRow := m.DB.QueryRowContext(ctx, query, isReccuring)
	err = countRow.Scan(&totalRecords)
	if err != nil {
		return nil, 0, 0, err
	}

	lastPage := totalRecords / pageSize

	return orders, lastPage, totalRecords, nil
}

// GetAllOrders returns all orders if isReccuring = 0 or all subscriptions if isReccuring = 1
func (m *DBModel) GetOrderByID(orderID int) (Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var o Order

	query := `
	select 
		o.id, o.widget_id, o.transaction_id, o.customer_id, o.status_id, o.quantity, o.amount, o.created_at, o.updated_at,
		w.id, w.name, 
		t.id, t.amount, t.currency, t.last_four, t.expiry_month, t.expiry_year, t.payment_intent, t.payment_method, t.bank_return_code,
		c.id, c.first_name, c.last_name, c.email
	from 
		orders o
		left join widgets w on (o.widget_id = w.id )
		left join transactions t on (o.transaction_id = t.id)
		left join customers c on (o.customer_id = c.id)
	where 
		o.id = ?	
	`

	row := m.DB.QueryRowContext(ctx, query, orderID)

	err := row.Scan(
		&o.ID, &o.WidgetID, &o.TransactionID, &o.CustomerID,
		&o.StatusID, &o.Quantity, &o.Amount, &o.CreatedAt, &o.UpdatedAt,
		&o.Widget.ID, &o.Widget.Name,
		&o.Transaction.ID, &o.Transaction.Amount, &o.Transaction.Currency,
		&o.Transaction.LastFour, &o.Transaction.ExpiryMonth, &o.Transaction.ExpiryYear,
		&o.Transaction.PaymentIntent, &o.Transaction.PaymentMethod, &o.Transaction.BankReturnCode,
		&o.Customer.ID, &o.Customer.FirstName, &o.Customer.LastName, &o.Customer.Email,
	)
	if err != nil {
		return o, err
	}

	return o, nil
}

func (m *DBModel) UpdateOrderStatus(id, statusID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `update orders set status_id = ? where id = ? `

	_, err := m.DB.ExecContext(ctx, stmt, statusID, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *DBModel) GetAllUsers() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var users []*User

	q := `select id, last_name, first_name, email, created_at, updated_at
		from users 
		order by last_name, first_name
	`

	rows, err := m.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		err = rows.Scan(
			&u.ID,
			&u.LastName,
			&u.FirstName,
			&u.Email,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, nil
}

func (m *DBModel) GetOneUser(id int) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u User

	q := `select id, last_name, first_name, email, created_at, updated_at
		from users 
		where id = ?
	`

	row := m.DB.QueryRowContext(ctx, q, id)

	err := row.Scan(
		&u.ID,
		&u.LastName,
		&u.FirstName,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return u, err
	}

	return u, nil

}

func (m *DBModel) EditUser(u User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := ` update users set 
		last_name = ?, first_name = ?, email = ?, updated_at = ? 
		where id = ? `

	_, err := m.DB.ExecContext(ctx, stmt,
		u.LastName,
		u.FirstName,
		u.Email,
		time.Now(),
		u.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) AddUser(u User, hash string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := ` insert into users (last_name, first_name, email, password, created_at, updated_at )
		values (?, ?, ?, ?, ?, ?)`

	_, err := m.DB.ExecContext(ctx, stmt,
		u.LastName,
		u.FirstName,
		u.Email,
		hash,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) DeleteUser(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `delete from users where id = ?`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	// no foreing key between users and tokens so do it manually
	stmt = `delete from tokens where user_id = ?`

	_, err = m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}
