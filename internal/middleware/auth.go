package middleware

type RoleType string

const (
	RoleCustomer RoleType = "customer"
	RoleStaff    RoleType = "staff"
	RoleAdmin    RoleType = "admin"
)

