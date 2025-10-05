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

func NoContent(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusNoContent).JSON(nil)
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

func Unauthorized(ctx *fiber.Ctx, msg string) error {
	return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{"message": msg})
}

func Forbidden(ctx *fiber.Ctx, msg string) error {
	return ctx.Status(http.StatusForbidden).JSON(&fiber.Map{"message": msg})
}
