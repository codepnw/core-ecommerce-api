package products

import (
	"errors"

	"github.com/codepnw/core-ecommerce-system/internal/utils/commons"
	"github.com/codepnw/core-ecommerce-system/internal/utils/errs"
	"github.com/codepnw/core-ecommerce-system/internal/utils/response"
	"github.com/codepnw/core-ecommerce-system/internal/utils/validate"
	"github.com/gofiber/fiber/v2"
)

const (
	productIDKey  = "product_id"
	categoryIDKey = "category_id"
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
	id, err := commons.GetParamIDInt(ctx, productIDKey)
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
	categoryID := int64(ctx.QueryInt(categoryIDKey))
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
	id, err := commons.GetParamIDInt(ctx, productIDKey)
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
	id, err := commons.GetParamIDInt(ctx, productIDKey)
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
	id, err := commons.GetParamIDInt(ctx, productIDKey)
	if err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err = h.srv.Delete(ctx.Context(), id); err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "product deleted", nil)
}

func (h *productHandler) GetCategoriesByProduct(ctx *fiber.Ctx) error {
	id, err := commons.GetParamIDInt(ctx, productIDKey)
	if err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	cats, err := h.srv.GetCategoriesByProduct(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, errs.ErrProductOrCategoryNotFound) {
			return response.NotFound(ctx, err.Error())
		}
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "", cats)
}

func (h *productHandler) AssignCategories(ctx *fiber.Ctx) error {
	req := new(ProductCategoryRequest)
	if err := ctx.BodyParser(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err := validate.Struct(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err := h.srv.AssignCategories(ctx.Context(), req); err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "added product to category success", nil)
}

func (h *productHandler) DelCategoryByProduct(ctx *fiber.Ctx) error {
	pID, err := commons.GetParamIDInt(ctx, productIDKey)
	if err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	cID, err := commons.GetParamIDInt(ctx, "category_id")
	if err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err = h.srv.DelCategoryByProduct(ctx.Context(), pID, cID); err != nil {
		if errors.Is(err, errs.ErrProductOrCategoryNotFound) {
			return response.NotFound(ctx, err.Error())
		}
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "delete category product success", nil)
}
