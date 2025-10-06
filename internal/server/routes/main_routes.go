package routes

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/codepnw/core-ecommerce-system/internal/database"
	"github.com/codepnw/core-ecommerce-system/internal/middleware"
	"github.com/codepnw/core-ecommerce-system/internal/utils/security"
	"github.com/codepnw/core-ecommerce-system/internal/utils/validate"
	"github.com/gofiber/fiber/v2"
)

type RoutesConfig struct {
	DB     *sql.DB                      `validate:"required"`
	Tx     *database.TxManager          `validate:"required"`
	Router *fiber.App                   `validate:"required"`
	Prefix string                       `validate:"required"`
	Mid    *middleware.MiddlewareConfig `validate:"required"`
	Token  *security.JWTToken           `validate:"required"`
}

func InitRoutes(cfg *RoutesConfig) error {
	if err := validate.Struct(cfg); err != nil {
		return errors.New("routes config required: DB, Router, Prefix")
	}

	cfg.categoryRoutes()
	cfg.addressRoutes()
	cfg.cartRoutes()
	cfg.orderRoutes()

	cfg.registerProductRoutes()

	if err := cfg.registerUserRoutes(); err != nil {
		return fmt.Errorf("UserRoutes: %w", err)
	}

	return nil
}
