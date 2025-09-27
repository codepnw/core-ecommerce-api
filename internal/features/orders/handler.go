package orders

import (
	"github.com/codepnw/core-ecommerce-system/internal/utils/response"
	"github.com/codepnw/core-ecommerce-system/internal/utils/validate"
	"github.com/gofiber/fiber/v2"
)

type orderHandler struct {
	srv IOrderService
}

func NewOrderHandler(srv IOrderService) *orderHandler {
	return &orderHandler{srv: srv}
}

func (h *orderHandler) CreateOrder(ctx *fiber.Ctx) error {
	req := new(OrderRequest)
	if err := ctx.BodyParser(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err := validate.Struct(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err := h.srv.CreateOrder(ctx.Context(), req); err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "order created", nil)
}
