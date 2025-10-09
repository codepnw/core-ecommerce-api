package carts

import (
	"github.com/codepnw/core-ecommerce-system/internal/middleware"
	"github.com/codepnw/core-ecommerce-system/internal/utils/commons"
	"github.com/codepnw/core-ecommerce-system/internal/utils/response"
	"github.com/codepnw/core-ecommerce-system/internal/utils/validate"
	"github.com/gofiber/fiber/v2"
)

type cartHandler struct {
	srv ICartService
}

func NewCartHandler(srv ICartService) *cartHandler {
	return &cartHandler{srv: srv}
}

func (h *cartHandler) AddItem(ctx *fiber.Ctx) error {
	user, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		return response.Unauthorized(ctx, err.Error())
	}

	req := new(CartItemRequest)
	if err := ctx.BodyParser(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err := validate.Struct(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err := h.srv.AddItem(ctx.Context(), user.UserID, req); err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "items added", nil)
}

func (h *cartHandler) GetCart(ctx *fiber.Ctx) error {
	user, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		return response.Unauthorized(ctx, err.Error())
	}

	res, err := h.srv.GetCart(ctx.Context(), user.UserID)
	if err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "", res)
}

func (h *cartHandler) RemoveItem(ctx *fiber.Ctx) error {
	user, err := commons.GetCurrentUser(ctx)
	if err != nil {
		return response.Unauthorized(ctx, err.Error())
	}

	productID, err := commons.GetParamIDInt(ctx, "product_id")
	if err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err := h.srv.RemoveItem(ctx.Context(), user.ID, productID); err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "product remove", nil)
}

func (h *cartHandler) ClearCart(ctx *fiber.Ctx) error {
	user, err := commons.GetCurrentUser(ctx)
	if err != nil {
		return response.Unauthorized(ctx, err.Error())
	}

	if err := h.srv.ClearCart(ctx.Context(), user.ID); err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "cart is empty", nil)
}
