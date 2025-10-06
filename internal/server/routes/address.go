package routes

import (
	"github.com/codepnw/core-ecommerce-system/internal/features/addresses"
	"github.com/codepnw/core-ecommerce-system/internal/middleware"
)

func (cfg *RoutesConfig) registerAddressRoutes() {
	repo := addresses.NewAddressRepository(cfg.DB)
	service := addresses.NewAddressSerivce(repo)
	handler := addresses.NewAddressHandler(service)

	const (
		userID           = "/:user_id"
		userIDAddress    = "/:user_id/address"
		addressID        = "/:address_id"
		addressIDDefault = "/:address_id/default"
	)

	r := cfg.Router.Group(cfg.Prefix+"/addresses", cfg.Mid.Authorized())
	staff := r.Group("", cfg.Mid.RoleRequired(middleware.RoleAdmin, middleware.RoleStaff))

	r.Post("/", handler.CreateAddress)
	r.Get(addressID, handler.GetAddressByID)
	r.Patch(addressID, handler.UpdateAddress)
	r.Delete(addressID, handler.DeleteAddress)
	r.Patch(addressIDDefault, handler.SetAddressDefault)

	// Admin & Staff
	staff.Get(userIDAddress, handler.GetAddressByUserID)
}
