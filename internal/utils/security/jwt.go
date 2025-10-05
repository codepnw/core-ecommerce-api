package security

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/codepnw/core-ecommerce-system/config"
	"github.com/golang-jwt/jwt/v5"
)

type JWTToken struct {
	cfg *config.EnvConfig
}

type jwtTokenClaims struct {
	UserID string
	Email  string
	Role   string
	jwt.RegisteredClaims
}

type UserTokenReq struct {
	UserID string
	Email  string
	Role   string
}

func InitJWT(cfg *config.EnvConfig) *JWTToken {
	if cfg == nil {
		log.Fatal("JWT config is nil")
	}
	return &JWTToken{cfg: cfg}
}

// ------- Generate Token -------

func (j *JWTToken) GenerateAccessToken(req *UserTokenReq) (string, error) {
	duration := time.Hour * 24
	return j.generateToken(j.cfg.JWT.SecretKey, duration, req)
}

func (j *JWTToken) GenerateRefreshToken(req *UserTokenReq) (string, error) {
	duration := time.Hour * 24 * 7
	return j.generateToken(j.cfg.JWT.RefreshKey, duration, req)
}

func (j *JWTToken) generateToken(key string, duration time.Duration, req *UserTokenReq) (string, error) {
	claims := &jwtTokenClaims{
		UserID: req.UserID,
		Email:  req.Email,
		Role:   req.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "ecommerce-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("sign token failed: %w", err)
	}

	return ss, nil
}

// ------- Verify Token -------

func (j *JWTToken) VerifyAccessToken(token string) (*jwtTokenClaims, error) {
	return j.verifyToken(j.cfg.JWT.SecretKey, token)
}

func (j *JWTToken) VerifyRefreshToken(token string) (*jwtTokenClaims, error) {
	return j.verifyToken(j.cfg.JWT.RefreshKey, token)
}

func (j *JWTToken) verifyToken(key, tokenStr string) (*jwtTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwtTokenClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token failed: %w", err)
	}

	claims, ok := token.Claims.(*jwtTokenClaims)
	if !token.Valid || !ok {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

func (j *JWTToken) HashToken(token []byte) string {
	hashArr := sha256.Sum256(token)
	return hex.EncodeToString(hashArr[:])
}
