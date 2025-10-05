package middleware

import (
	"errors"
	"fmt"
	"strings"

	"github.com/codepnw/core-ecommerce-system/internal/utils/response"
	"github.com/codepnw/core-ecommerce-system/internal/utils/security"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type RoleType string

const (
	RoleCustomer RoleType = "customer"
	RoleStaff    RoleType = "staff"
	RoleAdmin    RoleType = "admin"
)

const UserContextKey = "user-context"

type UserContext struct {
	UserID string
	Email  string
	Role   RoleType
	Exp    *jwt.NumericDate
}

type MiddlewareConfig struct {
	token *security.JWTToken
}

func InitMiddleware(token *security.JWTToken) *MiddlewareConfig {
	return &MiddlewareConfig{token: token}
}

func (m *MiddlewareConfig) Authorized() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return response.Unauthorized(ctx, "auth header is missing")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.Unauthorized(ctx, "invalid authorizarion format")
		}

		claims, err := m.token.VerifyAccessToken(parts[1])
		if err != nil {
			msg := fmt.Sprintf("invalid token or expired: %v", err)
			return response.Unauthorized(ctx, msg)
		}

		user := &UserContext{
			UserID: claims.UserID,
			Email:  claims.Email,
			Role:   RoleType(claims.Role),
			Exp:    claims.ExpiresAt,
		}

		ctx.Locals(UserContextKey, user)
		return ctx.Next()
	}
}

func (m *MiddlewareConfig) RoleRequired(roles ...RoleType) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userCtx := ctx.Locals(UserContextKey)
		if userCtx == nil {
			return response.Unauthorized(ctx, "user context is missing")
		}

		user, ok := userCtx.(*UserContext)
		if !ok {
			return response.Unauthorized(ctx, "invalid user context")
		}

		for _, role := range roles {
			if user.Role == role {
				return ctx.Next()
			}
		}

		return response.Forbidden(ctx, "no permissions")
	}
}

func GetUserFromContext(ctx *fiber.Ctx) (*UserContext, error) {
	val := ctx.Locals(UserContextKey)
	user, ok := val.(*UserContext)
	if !ok {
		return nil, errors.New("invalid user from context")
	}
	return user, nil
}
