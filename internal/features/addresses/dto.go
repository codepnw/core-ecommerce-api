package addresses

type AddressCreate struct {
	UserID      string `json:"user_id" validate:"required"`
	AddressLine string `json:"address_line" validate:"required"`
	City        string `json:"city" validate:"required"`
	State       string `json:"state" validate:"required"`
	PostalCode  string `json:"postal_code" validate:"required,min=5"`
	Phone       string `json:"phone" validate:"required,min=10"`
	IsDefault   bool   `json:"is_default"`
}

type AddressUpdate struct {
	AddressLine *string `json:"address_line,omitempty" validate:"omitempty"`
	City        *string `json:"city,omitempty" validate:"omitempty"`
	State       *string `json:"state,omitempty" validate:"omitempty"`
	PostalCode  *string `json:"postal_code,omitempty" validate:"omitempty,min=5"`
	Phone       *string `json:"phone,omitempty" validate:"omitempty,min=10"`
	IsDefault   *bool   `json:"is_default,omitempty"`
}
