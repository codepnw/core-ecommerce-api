package products

type ProductCreate struct {
	CategoryID  int64   `json:"category_id" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description,omitempty" validate:"omitempty"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock,omitempty" validate:"omitempty"`
	ImageURL    string  `json:"image_url,omitempty" validate:"omitempty"`
}

type ProductUpdate struct {
	CategoryID  *int64   `json:"category_id,omitempty" validate:"omitempty"`
	Name        *string  `json:"name,omitempty" validate:"omitempty"`
	Description *string  `json:"description,omitempty" validate:"omitempty"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,gt=0"`
	Stock       *int     `json:"stock,omitempty" validate:"omitempty"`
	ImageURL    *string  `json:"image_url,omitempty" validate:"omitempty"`
}
