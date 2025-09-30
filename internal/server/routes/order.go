package routes

import (
	"github.com/codepnw/core-ecommerce-system/internal/features/addresses"
	"github.com/codepnw/core-ecommerce-system/internal/features/carts"
	"github.com/codepnw/core-ecommerce-system/internal/features/orders"
)

func (cfg *RoutesConfig) orderRoutes() {
	cRepo := carts.NewCartRepository(cfg.DB)
	cService := carts.NewCartService(cRepo)

	aRepo := addresses.NewAddressRepository(cfg.DB)
	aService := addresses.NewAddressSerivce(aRepo)

	oRepo := orders.NewOrderRepository(cfg.DB)
	oService := orders.NewOrderService(cfg.Tx, oRepo, cService, aService)
	handler := orders.NewOrderHandler(oService)

	r := cfg.Router.Group(cfg.Prefix + "/orders")

	r.Post("/", handler.CreateOrder)
	r.Get("/", handler.ListOrders)
	r.Get("/:order_id", handler.UpdateOrderStatus)
}
