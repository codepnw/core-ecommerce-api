package errs

import "errors"

// Categories
var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrNoFieldUpdate    = errors.New("no fields to update")
)

var (
	ErrProductNotFound = errors.New("product not found")
)
