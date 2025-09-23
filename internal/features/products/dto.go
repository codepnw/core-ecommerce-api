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

type ProductUpdateStock struct {
	Quantity int `json:"quantity" validate:"required,gt=0"`
}

type ProductFilter struct {
	CategoryID *int64  `json:"category_id,omitempty"`
	OrderBy    *string `json:"order_by,omitempty"`
	Sort       *string `json:"sort,omitempty"`
	Limit      *int    `json:"limit,omitempty"`
	Offset     *int    `json:"offset,omitempty"`
}

type ProductListParams struct {
	CategoryID int64
	OrderBy    *string
	Sort       *string
	Limit      int
	Offset     int
}
