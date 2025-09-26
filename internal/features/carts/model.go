package carts

import "time"

type Cart struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ProductID int64     `json:"product_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
