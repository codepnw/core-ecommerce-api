package auth

import (
	"context"
	"errors"
	"time"

	"github.com/codepnw/core-ecommerce-system/internal/features/users"
	"github.com/codepnw/core-ecommerce-system/internal/utils/consts"
	"github.com/codepnw/core-ecommerce-system/internal/utils/errs"
	"github.com/codepnw/core-ecommerce-system/internal/utils/security"
)

type IAuthService interface {
	Register(ctx context.Context, req *users.UserCreate) (*TokenResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*TokenResponse, error)
	RefreshToken(ctx context.Context, userID string, req *RefreshTokenRequest) (*TokenResponse, error)
	Logout(ctx context.Context, refreshToken string) error
}

type authService struct {
	repo  IAuthRepository
	user  users.IUserService
	token *security.JWTToken
}

func NewAuthService(repo IAuthRepository, user users.IUserService) IAuthService {
	return &authService{
		repo: repo,
		user: user,
	}
}

func (s *authService) Login(ctx context.Context, req *LoginRequest) (*TokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	u, err := s.user.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errs.ErrInvalidEmailOrPassword
	}

	ok := security.ComparePassword(req.Password, u.PasswordHash)
	if !ok {
		return nil, errs.ErrInvalidEmailOrPassword
	}

	tokenReq := s.userTokenReq(u)
	accessToken, err := s.token.GenerateAccessToken(tokenReq)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.token.GenerateRefreshToken(tokenReq)
	if err != nil {
		return nil, err
	}

	hashedToken := s.token.HashToken([]byte(refreshToken))
	err = s.repo.Save(ctx, &AuthToken{
		UserID:    u.ID,
		Token:     hashedToken,
		ExpiredAt: time.Now().Add(time.Hour * 24 * 7),
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

func (s *authService) Register(ctx context.Context, req *users.UserCreate) (*TokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	req.Password = hashedPassword
	u, err := s.user.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	tokenReq := s.userTokenReq(u)
	accessToken, err := s.token.GenerateAccessToken(tokenReq)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.token.GenerateRefreshToken(tokenReq)
	if err != nil {
		return nil, err
	}

	hashedToken := s.token.HashToken([]byte(refreshToken))
	err = s.repo.Save(ctx, &AuthToken{
		UserID:    u.ID,
		Token:     hashedToken,
		ExpiredAt: time.Now().Add(time.Hour * 24 * 7),
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

func (s *authService) RefreshToken(ctx context.Context, userID string, req *RefreshTokenRequest) (*TokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	claims, err := s.token.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	if claims.UserID != userID {
		return nil, errors.New("invalid user token")
	}

	hashedToken := s.token.HashToken([]byte(req.RefreshToken))
	authToken, err := s.repo.GetToken(ctx, claims.UserID, hashedToken)
	if err != nil {
		return nil, err
	}

	if authToken == nil {
		return nil, errs.ErrUserTokenNotFound
	}

	if time.Now().After(authToken.ExpiredAt) {
		return nil, errs.ErrUserTokenExpired
	}

	tokenReq := &security.UserTokenReq{
		UserID: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
	}
	accessToken, err := s.token.GenerateAccessToken(tokenReq)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.token.GenerateRefreshToken(tokenReq)
	if err != nil {
		return nil, err
	}

	newHashedToken := s.token.HashToken([]byte(req.RefreshToken))
	err = s.repo.Save(ctx, &AuthToken{
		UserID:    claims.UserID,
		Token:     newHashedToken,
		ExpiredAt: time.Now().Add(time.Hour * 24 * 7),
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

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	if err := s.repo.Delete(ctx, refreshToken); err != nil {
		return err
	}

	return nil
}

func (s *authService) userTokenReq(u *users.User) *security.UserTokenReq {
	req := &security.UserTokenReq{
		UserID: u.ID,
		Email:  u.Email,
		Role:   u.Role,
	}
	return req
}
