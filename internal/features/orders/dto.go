package orders

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusPaid      OrderStatus = "paid"
	StatusShipped   OrderStatus = "shipped"
	StatusComplated OrderStatus = "completed"
	StatusCancelled OrderStatus = "cancelled"
)

type OrdersResponse struct {
	OrderID    int64   `json:"order_id"`
	Email      string  `json:"email"`
	FullName   string  `json:"full_name"`
	TotalPrice float64 `json:"total_price"`
	Phone      string  `json:"phone"`
	City       string  `json:"city"`
	State      string  `json:"state"`
	Status     string  `json:"status"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type OrderFilter struct {
	Status *string
	UserID *string
	Limit  *int
	Offset *int
}

type OrderRequest struct {
	UserID    string `json:"user_id"`
	AddressID string `json:"address_id"`
}

type OrderItemRequest struct {
	OrderID   int64   `json:"order_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type OrderAddressRequest struct {
	OrderID     int64  `json:"order_id"`
	AddressID   string `json:"address_id"`
	AddressLine string `json:"address_line"`
	City        string `json:"city"`
	State       string `json:"state"`
	PostalCode  string `json:"postal_code"`
	Phone       string `json:"phone"`
}
