package errs

import "errors"

// Categories
var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrNoFieldUpdate    = errors.New("no fields to update")
)

// Products
var (
	ErrProductNotFound = errors.New("product not found")
)

// Users
var (
	ErrUserNotFound = errors.New("user not found")
)
