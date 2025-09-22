package users

type Role string

const (
	RoleCustomer Role = "customer"
	RoleStaff    Role = "staff"
	RoleAdmin    Role = "admin"
)

type UserCreate struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"`
	Role     Role   `json:"role,omitempty" validate:"omitempty"`
}

type UserUpdate struct {
	FullName *string `json:"full_name,omitempty" validate:"omitempty"`
}

type UserUpdateForAdmin struct {
	Role *Role `json:"role,omitempty" validate:"omitempty,oneof=customer staff admin"`
}
