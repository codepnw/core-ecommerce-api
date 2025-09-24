package users

import (
	"errors"

	"github.com/codepnw/core-ecommerce-system/internal/utils/commons"
	"github.com/codepnw/core-ecommerce-system/internal/utils/errs"
	"github.com/codepnw/core-ecommerce-system/internal/utils/response"
	"github.com/codepnw/core-ecommerce-system/internal/utils/validate"
	"github.com/gofiber/fiber/v2"
)

const userIDKey = "user_id"

type userHandler struct {
	srv IUserService
}

func NewUserHandler(srv IUserService) *userHandler {
	return &userHandler{srv: srv}
}

func (h *userHandler) CreateUser(ctx *fiber.Ctx) error {
	req := new(UserCreate)
	if err := ctx.BodyParser(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err := validate.Struct(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	created, err := h.srv.CreateUser(ctx.Context(), req)
	if err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.Created(ctx, "added new user", created)
}

func (h *userHandler) GetUser(ctx *fiber.Ctx) error {
	id, err := commons.GetParamIDStr(ctx, userIDKey)
	if err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	user, err := h.srv.GetUser(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return response.NotFound(ctx, err.Error())
		}
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "", user)
}

func (h *userHandler) GetUsers(ctx *fiber.Ctx) error {
	limit := ctx.QueryInt("limit")
	offset := ctx.QueryInt("offset")

	users, err := h.srv.GetUsers(ctx.Context(), uint(limit), uint(offset))
	if err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "", users)
}

func (h *userHandler) UpdateUser(ctx *fiber.Ctx) error {
	id, err := commons.GetParamIDStr(ctx, userIDKey)
	if err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	req := new(UserUpdate)
	if err := ctx.BodyParser(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err := validate.Struct(req); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err = h.srv.UpdateUser(ctx.Context(), id, req); err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return response.NotFound(ctx, err.Error())
		}
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "user updated", nil)
}

func (h *userHandler) DeleteUser(ctx *fiber.Ctx) error {
	id, err := commons.GetParamIDStr(ctx, userIDKey)
	if err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	if err = h.srv.DeleteUser(ctx.Context(), id); err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return response.NotFound(ctx, err.Error())
		}
		return response.InternalServerError(ctx, err)
	}

	return response.Success(ctx, "user deleted", nil)
}
