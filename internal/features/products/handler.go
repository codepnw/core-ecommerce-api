package products

import (
	"errors"

	"github.com/codepnw/core-ecommerce-system/internal/utils/commons"
	"github.com/codepnw/core-ecommerce-system/internal/utils/errs"
	"github.com/codepnw/core-ecommerce-system/internal/utils/response"
	"github.com/codepnw/core-ecommerce-system/internal/utils/validate"
	"github.com/gofiber/fiber/v2"
)

type productHandler struct {
	srv IProductService
}

func NewProductHandler(srv IProductService) *productHandler {
	return &productHandler{srv: srv}
}

func (h *productHandler) CreateProduct(ctx *fiber.Ctx) error {
	req := new(ProductCreate)
	if err := ctx.BodyParser(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err := validate.Struct(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	created, err := h.srv.Create(ctx.Context(), req)
	if err != nil {
		if errors.Is(err, errs.ErrProductNotFound) {
			return response.NotFound(ctx, err.Error())
		}
		return response.InternalServerError(ctx, err)
	}

	return response.Created(ctx, "product added", created)
}

func (h *productHandler) GetProduct(ctx *fiber.Ctx) error {
	id, err := commons.GetParamIDInt(ctx)
	if err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	product, err := h.srv.GetByID(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, errs.ErrProductNotFound) {
			return response.NotFound(ctx, err.Error())
		}
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "", product)
}

func (h *productHandler) GetProducts(ctx *fiber.Ctx) error {
	categoryID := int64(ctx.QueryInt("category_id"))
	orderBy := ctx.Query("order_by")
	sort := ctx.Query("sort")
	limit := ctx.QueryInt("limit")
	offset := ctx.QueryInt("offset")

	filter := &ProductFilter{
		CategoryID: &categoryID,
		OrderBy:    &orderBy,
		Sort:       &sort,
		Limit:      &limit,
		Offset:     &offset,
	}

	products, err := h.srv.List(ctx.Context(), filter)
	if err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "", products)
}

func (h *productHandler) UpdateStock(ctx *fiber.Ctx) error {
	id, err := commons.GetParamIDInt(ctx)
	if err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	req := new(ProductUpdateStock)
	if err := ctx.BodyParser(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err := validate.Struct(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	err = h.srv.UpdateStock(ctx.Context(), id, req.Quantity)
	if err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "product stock updated", nil)
}

func (h *productHandler) UpdateProduct(ctx *fiber.Ctx) error {
	id, err := commons.GetParamIDInt(ctx)
	if err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	req := new(ProductUpdate)
	if err := ctx.BodyParser(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err := validate.Struct(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err = h.srv.Update(ctx.Context(), id, req); err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "product updated", nil)
}

func (h *productHandler) DeleteProduct(ctx *fiber.Ctx) error {
	id, err := commons.GetParamIDInt(ctx)
	if err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err = h.srv.Delete(ctx.Context(), id); err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "product deleted", nil)
}
