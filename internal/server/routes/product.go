package routes

import "github.com/codepnw/core-ecommerce-system/internal/features/products"

func (cfg *RoutesConfig) productRoutes() {
	repo := products.NewProductRepository(cfg.DB)
	service := products.NewProductService(repo)
	handler := products.NewProductHandler(service)

	r := cfg.Router.Group(cfg.Prefix + "/products")

	r.Post("/", handler.CreateProduct)
	r.Get("/", handler.GetProducts)
	r.Get("/:id", handler.GetProduct)
	r.Patch("/:id", handler.UpdateProduct)
	r.Patch("/:id/stock", handler.UpdateStock)
	r.Delete("/:id", handler.DeleteProduct)
}
