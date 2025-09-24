package commons

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

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
