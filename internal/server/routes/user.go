package routes

import (
	"github.com/codepnw/core-ecommerce-system/internal/features/users"
)

func (cfg *RoutesConfig) userRoutes() {
	repo := users.NewUserRepository(cfg.DB)
	service := users.NewUserService(repo)
	handler := users.NewUserHandler(service)

	r := cfg.Router.Group(cfg.Prefix + "/users")

	r.Post("/", handler.CreateUser)
	r.Get("/", handler.GetUsers)
	r.Get("/:id", handler.GetUser)
	r.Patch("/:id", handler.UpdateUser)
	r.Delete("/:id", handler.DeleteUser)
}
