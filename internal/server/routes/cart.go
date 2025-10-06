package routes

import "github.com/codepnw/core-ecommerce-system/internal/features/carts"

func (cfg *RoutesConfig) registerCartRoutes() {
	repo := carts.NewCartRepository(cfg.DB)
	service := carts.NewCartService(repo)
	handler := carts.NewCartHandler(service)

	r := cfg.Router.Group(cfg.Prefix + "/cart", cfg.Mid.Authorized())

	r.Post("/", handler.AddItem)
	r.Get("/", handler.GetCart)
	r.Delete("/clear", handler.ClearCart)
	r.Delete("/remove/:product_id", handler.RemoveItem)
}
