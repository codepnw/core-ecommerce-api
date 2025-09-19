package routes

import "github.com/codepnw/core-ecommerce-system/internal/features/categories"

func (cfg *RoutesConfig) CategoryRoutes() {
	repo := categories.NewCategoryRepository(cfg.DB)
	service := categories.NewCategoryService(repo)
	handler := categories.NewCategoryHandler(service)

	r := cfg.Router.Group(cfg.Prefix + "/categories")

	r.Post("/", handler.Create)
	r.Get("/", handler.List)
	r.Patch("/:id", handler.Update)
	r.Delete("/:id", handler.Delete)
}
