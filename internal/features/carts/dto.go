package carts

type CartItemRequest struct {
	ProductID int64 `json:"product_id" validate:"required"`
	Quantity  int   `json:"quantity" validate:"required,gt=0"`
}

type CartItemsResponse struct {
	ProductID       int64  `json:"product_id"`
	ProductName     string `json:"product_name"`
	ProductPrice    float64  `json:"product_price"`
	ProductQuantity int64    `json:"product_quantity"`
}
