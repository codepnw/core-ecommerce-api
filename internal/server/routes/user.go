package routes

import (
	"fmt"

	"github.com/codepnw/core-ecommerce-system/internal/features/auth"
	"github.com/codepnw/core-ecommerce-system/internal/features/users"
	"github.com/codepnw/core-ecommerce-system/internal/middleware"
)

func (cfg *RoutesConfig) userRoutes() error {
	// User Setup
	uRepo := users.NewUserRepository(cfg.DB)
	uService := users.NewUserService(uRepo)
	uHandler := users.NewUserHandler(uService)

	const userID = "/:user_id"

	r := cfg.Router.Group(cfg.Prefix + "/users")
	uPrivate := r.Use(
		cfg.Mid.Authorized(),
		cfg.Mid.RoleRequired(middleware.RoleAdmin, middleware.RoleStaff),
	)

	uPrivate.Post("/", uHandler.CreateUser)
	uPrivate.Get("/", uHandler.GetUsers)
	uPrivate.Get(userID, uHandler.GetUser)
	uPrivate.Patch(userID, uHandler.UpdateUser)
	uPrivate.Delete(userID, uHandler.DeleteUser)

	// Auth Setup
	aRepo := auth.NewAuthRepository(cfg.DB)
	aServiceCfg := &auth.AuthServiceConfig{
		AuthRepo: aRepo,
		UserSrv:  uService,
		Token:    cfg.Token,
		Tx:       cfg.Tx,
		DB:       cfg.DB,
	}
	aService, err := auth.NewAuthService(aServiceCfg)
	if err != nil {
		return fmt.Errorf("auth.NewAuthService Failed: %w", err)
	}
	aHandler := auth.NewAuthHandler(aService)

	aPublic := cfg.Router.Group(cfg.Prefix + "/auth")

	aPublic.Post("/register", aHandler.Register)
	aPublic.Post("/login", aHandler.Login)

	aPrivate := aPublic.Use(cfg.Mid.Authorized())

	aPrivate.Post("/refresh-token", aHandler.RefreshToken)
	aPrivate.Get("/logout", aHandler.Logout)

	return nil
}
