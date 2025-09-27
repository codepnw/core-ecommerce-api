package server

import (
	"fmt"
	"log"

	"github.com/codepnw/core-ecommerce-system/config"
	"github.com/codepnw/core-ecommerce-system/internal/database"
	"github.com/codepnw/core-ecommerce-system/internal/server/routes"
	"github.com/gofiber/fiber/v2"
)

func Run(cfg *config.EnvConfig) error {
	db, err := database.ConnectPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := fiber.New()

	// Setup Routes
	routeCfg := &routes.RoutesConfig{
		DB:     db,
		Tx:     database.NewTxManager(db),
		Router: app,
		Prefix: fmt.Sprintf("/api/v%d", cfg.APP.Version),
	}
	if err = routes.InitRoutes(routeCfg); err != nil {
		return err
	}

	url := fmt.Sprintf("%s:%d%s", cfg.APP.Host, cfg.APP.Port, routeCfg.Prefix)
	log.Printf("database %s connected...", cfg.DB.Name)
	log.Printf("server running at %s", url)

	return app.Listen(fmt.Sprintf(":%d", cfg.APP.Port))
}
