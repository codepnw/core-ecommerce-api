package commons

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type MockUser struct {
	ID string
}

func GetCurrentUser(ctx *fiber.Ctx) (*MockUser, error) {
	// TODO: Get User Context later
	u := &MockUser{
		ID: "3128d165-65c6-4ad9-85d0-231c617cb01e",
	}
	return u, nil
}

func GetParamIDInt(ctx *fiber.Ctx, key string) (int64, error) {
	idStr := ctx.Params(key)
	if idStr == "" {
		return 0, errors.New("id is required")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, errors.New("invalid id format")
	}

	return id, nil
}

func GetParamIDStr(ctx *fiber.Ctx, key string) (string, error) {
	id := ctx.Params(key)
	if id == "" {
		return "", errors.New("id is required")
	}
	return id, nil
}
