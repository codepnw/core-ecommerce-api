package routes

import (
	"github.com/codepnw/core-ecommerce-system/internal/features/addresses"
	"github.com/codepnw/core-ecommerce-system/internal/features/carts"
	"github.com/codepnw/core-ecommerce-system/internal/features/orders"
	"github.com/codepnw/core-ecommerce-system/internal/features/products"
)

func (cfg *RoutesConfig) registerOrderRoutes() error {
	pRepo := products.NewProductRepository(cfg.DB)
	pSerivce := products.NewProductService(pRepo)

	cRepo := carts.NewCartRepository(cfg.DB)
	cService := carts.NewCartService(cRepo)

	aRepo := addresses.NewAddressRepository(cfg.DB)
	aService := addresses.NewAddressSerivce(aRepo)

	oRepo := orders.NewOrderRepository(cfg.DB)
	oService, err := orders.NewOrderService(&orders.OrderServiceConfig{
		OrderRepo: oRepo,
		CartSrv:   cService,
		ProdSrv:   pSerivce,
		AddrSrv:   aService,
		Tx:        cfg.Tx,
	})
	if err != nil {
		return err
	}
	handler := orders.NewOrderHandler(oService)

	r := cfg.Router.Group(cfg.Prefix+"/orders", cfg.Mid.Authorized())

	r.Post("/", handler.CreateOrder)
	r.Get("/", handler.ListOrders)
	r.Get("/:order_id", handler.UpdateOrderStatus)

	// TODO: admin get order, update status

	return nil
}
