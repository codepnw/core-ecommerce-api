package orders

import (
	"log"

	"github.com/codepnw/core-ecommerce-system/internal/utils/commons"
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

func (h *orderHandler) ListOrders(ctx *fiber.Ctx) error {
	filter := new(OrderFilter)

	if status := ctx.Query("status"); status != "" {
		filter.Status = &status
	}

	if userID := ctx.Query("user_id"); userID != "" {
		filter.UserID = &userID
	}

	if limit := ctx.QueryInt("limit", 0); limit != 0 {
		filter.Limit = &limit
	}

	if offset := ctx.QueryInt("offset", 0); offset != 0 {
		filter.Offset = &offset
	}

	log.Printf("%+v\n", filter)

	res, err := h.srv.ListOrders(ctx.Context(), filter)
	if err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "", res)
}

func (h *orderHandler) UpdateOrderStatus(ctx *fiber.Ctx) error {
	status := ctx.Query("status")
	id, err := commons.GetParamIDInt(ctx, "order_id")
	if err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	s := OrderStatus(status)
	if err = h.srv.UpdateOrderStatus(ctx.Context(), id, s); err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "order status updated", nil)
}
