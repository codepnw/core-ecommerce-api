package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/codepnw/core-ecommerce-system/internal/database"
	"github.com/codepnw/core-ecommerce-system/internal/features/users"
	"github.com/codepnw/core-ecommerce-system/internal/utils/consts"
	"github.com/codepnw/core-ecommerce-system/internal/utils/errs"
	"github.com/codepnw/core-ecommerce-system/internal/utils/security"
	"github.com/codepnw/core-ecommerce-system/internal/utils/validate"
	"github.com/gofiber/fiber/v2/log"
)

type IAuthService interface {
	Register(ctx context.Context, req *users.UserCreate) (*TokenResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*TokenResponse, error)
	RefreshToken(ctx context.Context, userID string, req *RefreshTokenRequest) (*TokenResponse, error)
	Logout(ctx context.Context, userID string) error
}

type AuthServiceConfig struct {
	AuthRepo IAuthRepository     `validate:"required"`
	UserSrv  users.IUserService  `validate:"required"`
	Token    *security.JWTToken  `validate:"required"`
	Tx       *database.TxManager `validate:"required"`
	DB       *sql.DB             `validate:"required"`
}

func NewAuthService(cfg *AuthServiceConfig) (IAuthService, error) {
	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("AuthServiceConfig required all fields: %w", err)
	}
	return cfg, nil
}

func (s *AuthServiceConfig) Login(ctx context.Context, req *LoginRequest) (*TokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	u, err := s.UserSrv.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errs.ErrInvalidEmailOrPassword
	}

	ok := security.ComparePassword(req.Password, u.PasswordHash)
	if !ok {
		return nil, errs.ErrInvalidEmailOrPassword
	}

	accessToken, refreshToken, err := s.generateToken(u)
	if err != nil {
		return nil, err
	}

	hashedToken := s.Token.HashToken([]byte(refreshToken))
	err = s.AuthRepo.Save(ctx, s.DB, &AuthToken{
		UserID:    u.ID,
		Token:     hashedToken,
		ExpiredAt: time.Now().Add(consts.ExpRefreshToken),
	})
	if err != nil {
		return nil, err
	}

	response := &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return response, nil
}

func (s *AuthServiceConfig) Register(ctx context.Context, req *users.UserCreate) (*TokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	response := new(TokenResponse)

	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	req.Password = hashedPassword

	err = s.Tx.Transaction(ctx, func(tx *sql.Tx) error {
		u, err := s.UserSrv.CreateUserTx(ctx, tx, req)
		if err != nil {
			return err
		}

		accessToken, refreshToken, err := s.generateToken(u)
		if err != nil {
			return err
		}

		hashedToken := s.Token.HashToken([]byte(refreshToken))
		err = s.AuthRepo.Save(ctx, tx, &AuthToken{
			UserID:    u.ID,
			Token:     hashedToken,
			ExpiredAt: time.Now().Add(consts.ExpRefreshToken),
		})
		if err != nil {
			return err
		}

		response.AccessToken = accessToken
		response.RefreshToken = refreshToken
		return nil
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *AuthServiceConfig) RefreshToken(ctx context.Context, userID string, req *RefreshTokenRequest) (*TokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	claims, err := s.Token.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	if claims.UserID != userID {
		log.Errorf("claims.UserID: %v , userID: %v", claims.UserID, userID)
		return nil, errors.New("invalid user token")
	}

	hashedToken := s.Token.HashToken([]byte(req.RefreshToken))
	authToken, err := s.AuthRepo.GetToken(ctx, claims.UserID, hashedToken)
	if err != nil {
		return nil, err
	}

	if authToken == nil {
		return nil, errs.ErrUserTokenNotFound
	}

	if time.Now().After(authToken.ExpiredAt) {
		return nil, errs.ErrUserTokenExpired
	}

	accessToken, refreshToken, err := s.generateToken(&users.User{
		ID:    claims.UserID,
		Email: claims.Email,
		Role:  claims.Role,
	})
	if err != nil {
		return nil, err
	}

	newHashedToken := s.Token.HashToken([]byte(req.RefreshToken))
	err = s.AuthRepo.Save(ctx, s.DB, &AuthToken{
		UserID:    claims.UserID,
		Token:     newHashedToken,
		ExpiredAt: time.Now().Add(consts.ExpRefreshToken),
	})
	if err != nil {
		return nil, err
	}

	response := &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return response, nil
}

func (s *AuthServiceConfig) Logout(ctx context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	if err := s.AuthRepo.Delete(ctx, userID); err != nil {
		return err
	}
	return nil
}

func (s *AuthServiceConfig) generateToken(u *users.User) (string, string, error) {
	req := &security.UserTokenReq{
		UserID: u.ID,
		Email:  u.Email,
		Role:   u.Role,
	}

	accessToken, err := s.Token.GenerateAccessToken(req)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.Token.GenerateRefreshToken(req)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
