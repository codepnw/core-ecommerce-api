package routes

import (
	"fmt"

	"github.com/codepnw/core-ecommerce-system/internal/features/auth"
	"github.com/codepnw/core-ecommerce-system/internal/features/users"
	"github.com/codepnw/core-ecommerce-system/internal/middleware"
)

func (cfg *RoutesConfig) registerUserRoutes() error {
	// ------- User Setup ----------
	uRepo := users.NewUserRepository(cfg.DB)
	uService := users.NewUserService(uRepo)
	uHandler := users.NewUserHandler(uService)

	const userID = "/:user_id"
	userPath := fmt.Sprintf("%s/users", cfg.Prefix)
	authPath := fmt.Sprintf("%s/auth", cfg.Prefix)

	u := cfg.Router.Group(userPath, cfg.Mid.Authorized())
	staff := u.Group("", cfg.Mid.RoleRequired(middleware.RoleAdmin, middleware.RoleStaff))
	admin := u.Group("", cfg.Mid.RoleRequired(middleware.RoleAdmin))

	// Admin & Staff
	staff.Post("/", uHandler.CreateUser)
	staff.Get("/", uHandler.GetUsers)
	staff.Get(userID, uHandler.GetUser)
	staff.Patch(userID, uHandler.UpdateUser)

	// Admin Only
	admin.Delete(userID, uHandler.DeleteUser)

	// ------- Auth Setup ----------
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

	// Public Auth
	public := cfg.Router.Group(authPath)
	public.Post("/register", aHandler.Register)
	public.Post("/login", aHandler.Login)

	// Private Auth
	private := public.Group("", cfg.Mid.Authorized())
	private.Post("/refresh-token", aHandler.RefreshToken)
	private.Get("/logout", aHandler.Logout)

	return nil
}
