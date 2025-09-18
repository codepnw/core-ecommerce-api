package categories

type CategoryCreate struct {
	Name        string `json:"name" validate:"required,gte=3"`
	Description string `json:"description"`
}

type CategoryUpdate struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
