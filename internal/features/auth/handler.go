package auth

import (
	"github.com/codepnw/core-ecommerce-system/internal/features/users"
	"github.com/codepnw/core-ecommerce-system/internal/utils/commons"
	"github.com/codepnw/core-ecommerce-system/internal/utils/response"
	"github.com/codepnw/core-ecommerce-system/internal/utils/validate"
	"github.com/gofiber/fiber/v2"
)

type authHandler struct {
	srv IAuthService
}

func NewAuthHandler() *authHandler {
	return &authHandler{}
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
	user, err := commons.GetCurrentUser(ctx)
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

	res, err := h.srv.RefreshToken(ctx.Context(), user.ID, req)
	if err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "", res)
}

func (h *authHandler) Logout(ctx *fiber.Ctx) error {
	req := new(RefreshTokenRequest)
	if err := ctx.BodyParser(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err := validate.Struct(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err := h.srv.Logout(ctx.Context(), req.RefreshToken); err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "", nil)
}
