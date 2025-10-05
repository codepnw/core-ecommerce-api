package auth

import (
	"github.com/codepnw/core-ecommerce-system/internal/features/users"
	"github.com/codepnw/core-ecommerce-system/internal/middleware"
	"github.com/codepnw/core-ecommerce-system/internal/utils/response"
	"github.com/codepnw/core-ecommerce-system/internal/utils/validate"
	"github.com/gofiber/fiber/v2"
)

type authHandler struct {
	srv IAuthService
}

func NewAuthHandler(srv IAuthService) *authHandler {
	return &authHandler{srv: srv}
}

func (h *authHandler) Register(ctx *fiber.Ctx) error {
	req := new(users.UserCreate)
	if err := ctx.BodyParser(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err := validate.Struct(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	res, err := h.srv.Register(ctx.Context(), req)
	if err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.Created(ctx, "", res)
}

func (h *authHandler) Login(ctx *fiber.Ctx) error {
	req := new(LoginRequest)
	if err := ctx.BodyParser(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err := validate.Struct(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	res, err := h.srv.Login(ctx.Context(), req)
	if err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "", res)
}

func (h *authHandler) RefreshToken(ctx *fiber.Ctx) error {
	user, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		return response.Unauthorized(ctx, "")
	}

	req := new(RefreshTokenRequest)
	if err := ctx.BodyParser(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err := validate.Struct(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	res, err := h.srv.RefreshToken(ctx.Context(), user.UserID, req)
	if err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "", res)
}

func (h *authHandler) Logout(ctx *fiber.Ctx) error {
	user, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		return response.Unauthorized(ctx, "")
	}

	if err := h.srv.Logout(ctx.Context(), user.UserID); err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.NoContent(ctx)
}
