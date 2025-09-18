package response

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// ------- SUCCESS --------

func Created(ctx *fiber.Ctx, msg string, data any) error {
	return ctx.Status(http.StatusCreated).JSON(&fiber.Map{
		"message": msg,
		"data":    data,
	})
}

func Success(ctx *fiber.Ctx, msg string, data any) error {
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": msg,
		"data":    data,
	})
}

// ------- ERROR ----------

func BadRequest(ctx *fiber.Ctx, msg string) error {
	return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": msg})
}

func NotFound(ctx *fiber.Ctx, msg string) error {
	return ctx.Status(http.StatusNotFound).JSON(&fiber.Map{"message": msg})
}

func InternalServerError(ctx *fiber.Ctx, err error) error {
	return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{"error": err.Error()})
}
