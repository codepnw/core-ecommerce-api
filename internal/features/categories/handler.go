package categories

import (
	"github.com/codepnw/core-ecommerce-system/internal/utils/commons"
	"github.com/codepnw/core-ecommerce-system/internal/utils/response"
	"github.com/codepnw/core-ecommerce-system/internal/utils/validate"
	"github.com/gofiber/fiber/v2"
)

const categoryIDKey = "category_id"

type categoryHandler struct {
	srv CategoryService
}

func NewCategoryHandler(srv CategoryService) *categoryHandler {
	return &categoryHandler{srv: srv}
}

func (h *categoryHandler) Create(ctx *fiber.Ctx) error {
	req := new(CategoryCreate)
	if err := ctx.BodyParser(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err := validate.Struct(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err := h.srv.Create(ctx.Context(), req); err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.Created(ctx, "category created", req.Name)
}

func (h *categoryHandler) List(ctx *fiber.Ctx) error {
	cats, err := h.srv.List(ctx.Context())
	if err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "", cats)
}

func (h *categoryHandler) Update(ctx *fiber.Ctx) error {
	id, err := commons.GetParamIDInt(ctx, categoryIDKey)
	if err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	req := new(CategoryUpdate)
	if err := ctx.BodyParser(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err := validate.Struct(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err := h.srv.Update(ctx.Context(), id, req); err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "category updated", nil)
}

func (h *categoryHandler) Delete(ctx *fiber.Ctx) error {
	id, err := commons.GetParamIDInt(ctx, categoryIDKey)
	if err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err := h.srv.Delete(ctx.Context(), id); err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "category deleted", nil)
}
