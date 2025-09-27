package orders

import "time"

type Order struct {
	ID         int64     `json:"id"`
	UserID     string    `json:"user_id"`
	AddressID  string    `json:"address_id"`
	TotalPrice int64     `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type OrderItem struct {
	ID        int64     `json:"id"`
	OrderID   int64     `json:"order_id"`
	ProductID int64     `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     int64     `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrderAddress struct {
	ID          int64  `json:"id"`
	OrderID     int64  `json:"order_id"`
	AddressID   string `json:"address_id"`
	AddressLine string `json:"address_line"`
	City        string `json:"city"`
	State       string `json:"state"`
	PostalCode  string `json:"postal_code"`
	Phone       string `json:"phone"`
}
