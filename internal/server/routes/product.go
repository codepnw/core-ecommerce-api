package routes

import "github.com/codepnw/core-ecommerce-system/internal/features/products"

func (cfg *RoutesConfig) productRoutes() {
	repo := products.NewProductRepository(cfg.DB)
	service := products.NewProductService(repo)
	handler := products.NewProductHandler(service)

	const productID = "/:product_id"
	const categoryID = "/:category_id"

	r := cfg.Router.Group(cfg.Prefix + "/products")
	// Prouducts
	r.Post("/", handler.CreateProduct)
	r.Get("/", handler.GetProducts)
	r.Get(productID, handler.GetProduct)
	r.Patch(productID, handler.UpdateProduct)
	r.Delete(productID, handler.DeleteProduct)
	r.Patch(productID+"/stock", handler.UpdateStock)

	// Product Categories path /products/{product_id}/categories
	r.Get(productID+"/categories", handler.GetCategoriesByProduct)
	r.Delete(productID+"/categories"+categoryID, handler.DelCategoryByProduct)
	r.Post(productID+"/categories", handler.AssignCategories)
}
