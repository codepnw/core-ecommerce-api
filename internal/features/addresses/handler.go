package addresses

import (
	"errors"

	"github.com/codepnw/core-ecommerce-system/internal/utils/commons"
	"github.com/codepnw/core-ecommerce-system/internal/utils/consts"
	"github.com/codepnw/core-ecommerce-system/internal/utils/errs"
	"github.com/codepnw/core-ecommerce-system/internal/utils/response"
	"github.com/codepnw/core-ecommerce-system/internal/utils/validate"
	"github.com/gofiber/fiber/v2"
)

type addressHandler struct {
	srv IAddressServide
}

func NewAddressHandler(srv IAddressServide) *addressHandler {
	return &addressHandler{srv: srv}
}

func (h *addressHandler) CreateAddress(ctx *fiber.Ctx) error {
	req := new(AddressCreate)
	if err := ctx.BodyParser(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err := validate.Struct(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err := h.srv.CreateAddress(ctx.Context(), req); err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.Created(ctx, "added new address", nil)
}

func (h *addressHandler) GetAddressByID(ctx *fiber.Ctx) error {
	id, err := commons.GetParamIDInt(ctx, consts.KeyAddressParam)
	if err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	res, err := h.srv.GetAddressByID(ctx.Context(), id)
	if err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "", res)
}

func (h *addressHandler) GetAddressByUserID(ctx *fiber.Ctx) error {
	id, err := commons.GetParamIDStr(ctx, consts.KeyUserParam)
	if err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	res, err := h.srv.GetAddressByUserID(ctx.Context(), id)
	if err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "", res)
}

func (h *addressHandler) UpdateAddress(ctx *fiber.Ctx) error {
	id, err := commons.GetParamIDInt(ctx, consts.KeyAddressParam)
	if err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	req := new(AddressUpdate)
	if err = ctx.BodyParser(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err = validate.Struct(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err = h.srv.UpdateAddress(ctx.Context(), id, req); err != nil {
		if errors.Is(err, errs.ErrAddressNotFound) {
			return response.NotFound(ctx, err.Error())
		}
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "address updated", nil)
}

func (h *addressHandler) DeleteAddress(ctx *fiber.Ctx) error {
	id, err := commons.GetParamIDInt(ctx, consts.KeyAddressParam)
	if err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err = h.srv.DeleteAddress(ctx.Context(), id); err != nil {
		if errors.Is(err, errs.ErrAddressNotFound) {
			return response.NotFound(ctx, err.Error())
		}
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "address deleted", nil)
}
