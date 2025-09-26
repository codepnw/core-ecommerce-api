package routes

import (
	"database/sql"
	"errors"

	"github.com/codepnw/core-ecommerce-system/internal/utils/validate"
	"github.com/gofiber/fiber/v2"
)

type RoutesConfig struct {
	DB     *sql.DB    `validate:"required"`
	Router *fiber.App `validate:"required"`
	Prefix string     `validate:"required"`
}

func InitRoutes(cfg *RoutesConfig) error {
	if err := validate.Struct(cfg); err != nil {
		return errors.New("routes config required: DB, Router, Prefix")
	}

	cfg.categoryRoutes()
	cfg.productRoutes()
	cfg.userRoutes()
	cfg.addressRoutes()
	cfg.cartRoutes()

	return nil
}
