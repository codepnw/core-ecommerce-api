package routes

import (
	"github.com/codepnw/core-ecommerce-system/internal/features/categories"
	"github.com/codepnw/core-ecommerce-system/internal/middleware"
)

func (cfg *RoutesConfig) registerCategoryRoutes() {
	repo := categories.NewCategoryRepository(cfg.DB)
	service := categories.NewCategoryService(repo)
	handler := categories.NewCategoryHandler(service)

	const categoryID = "/:category_id"

	public := cfg.Router.Group(cfg.Prefix + "/categories")
	staff := public.Group("", cfg.Mid.Authorized(), cfg.Mid.RoleRequired(middleware.RoleAdmin, middleware.RoleStaff))
	
	public.Get("/", handler.List)

	// Admin & Staff
	staff.Post("/", handler.Create)
	staff.Patch(categoryID, handler.Update)
	staff.Delete(categoryID, handler.Delete)
}
