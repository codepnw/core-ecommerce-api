package commons

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetParamIDInt(ctx *fiber.Ctx) (int64, error) {
	idStr := ctx.Params("id")
	if idStr == "" {
		return 0, errors.New("id is required")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, errors.New("invalid id format")
	}

	return id, nil
}

func GetParamIDStr(ctx *fiber.Ctx) (string, error) {
	id := ctx.Params("id")
	if id == "" {
		return "", errors.New("id is required")
	}
	return id, nil
}
