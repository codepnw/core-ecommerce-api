package routes

import (
	"database/sql"
	"errors"

	"github.com/codepnw/core-ecommerce-system/internal/features/categories"
	"github.com/codepnw/core-ecommerce-system/internal/utils/validate"
	"github.com/gofiber/fiber/v2"
)

type RoutesConfig struct {
	DB     *sql.DB    `validate:"required"`
	Router *fiber.App `validate:"required"`
	Prefix string     `validate:"required"`
}

func InitRoutes(cfg *RoutesConfig) (*RoutesConfig, error) {
	if err := validate.Struct(cfg); err != nil {
		return nil, errors.New("routes config required: DB, Router, Prefix")
	}

	cfg.CategoryRoutes()

	return &RoutesConfig{
		DB:     cfg.DB,
		Router: cfg.Router,
		Prefix: cfg.Prefix,
	}, nil
}

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
