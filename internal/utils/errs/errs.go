package errs

import "errors"

// Categories
var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrNoFieldUpdate    = errors.New("no fields to update")
)

// Products
var (
	ErrProductNotFound   = errors.New("product not found")
	ErrProductOutOfStock = errors.New("product out of stock")
)

// Users
var (
	ErrUserNotFound = errors.New("user not found")
)
