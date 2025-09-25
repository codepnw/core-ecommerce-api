package addresses

import "time"

type Address struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	AddressLine string    `json:"address_line"`
	City        string    `json:"city"`
	State       string    `json:"state"`
	PostalCode  string    `json:"postal_code"`
	Phone       string    `json:"phone"`
	IsDefault   bool      `json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
