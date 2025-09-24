package routes

import "github.com/codepnw/core-ecommerce-system/internal/features/categories"

func (cfg *RoutesConfig) categoryRoutes() {
	repo := categories.NewCategoryRepository(cfg.DB)
	service := categories.NewCategoryService(repo)
	handler := categories.NewCategoryHandler(service)

	const categoryID = "/:category_id"

	r := cfg.Router.Group(cfg.Prefix + "/categories")

	r.Post("/", handler.Create)
	r.Get("/", handler.List)
	r.Patch(categoryID, handler.Update)
	r.Delete(categoryID, handler.Delete)
}
