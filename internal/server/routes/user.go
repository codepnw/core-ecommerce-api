package routes

import (
	"github.com/codepnw/core-ecommerce-system/internal/features/users"
)

func (cfg *RoutesConfig) userRoutes() {
	repo := users.NewUserRepository(cfg.DB)
	service := users.NewUserService(repo)
	handler := users.NewUserHandler(service)

	const userID = "/:user_id"

	r := cfg.Router.Group(cfg.Prefix + "/users")

	r.Post("/", handler.CreateUser)
	r.Get("/", handler.GetUsers)
	r.Get(userID, handler.GetUser)
	r.Patch(userID, handler.UpdateUser)
	r.Delete(userID, handler.DeleteUser)
}
