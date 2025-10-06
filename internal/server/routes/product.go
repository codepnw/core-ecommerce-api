package routes

import (
	"fmt"

	"github.com/codepnw/core-ecommerce-system/internal/features/products"
	"github.com/codepnw/core-ecommerce-system/internal/middleware"
)

func (cfg *RoutesConfig) registerProductRoutes() {
	repo := products.NewProductRepository(cfg.DB)
	service := products.NewProductService(repo)
	handler := products.NewProductHandler(service)

	const (
		productID         = "/:product_id"
		categoryID        = "/:category_id"
		productCategoryID = "/:product_id/categories"
	)
	path := fmt.Sprintf("%s/products", cfg.Prefix)

	public := cfg.Router.Group(path)
	protected := cfg.Router.Group(path, cfg.Mid.Authorized())
	staff := protected.Group("", cfg.Mid.RoleRequired(middleware.RoleAdmin, middleware.RoleStaff))
	admin := protected.Group("", cfg.Mid.RoleRequired(middleware.RoleAdmin))

	// Public
	public.Get("/", handler.GetProducts)
	public.Get(productID, handler.GetProduct)

	// Admin & Staff
	staff.Post("/", handler.CreateProduct)
	staff.Patch(productID+"/stock", handler.UpdateStock)

	// Admin Only
	admin.Delete(productID, handler.DeleteProduct)
	admin.Patch(productID, handler.UpdateProduct)

	// Product Categories path /products/{product_id}/categories
	admin.Delete(productCategoryID+categoryID, handler.DelCategoryByProduct)
	admin.Get(productCategoryID, handler.GetCategoriesByProduct)
	admin.Post(productCategoryID, handler.AssignCategories)
}
