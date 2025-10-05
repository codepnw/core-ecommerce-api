package users

import (
	"context"
	"database/sql"
	"strings"

	"github.com/codepnw/core-ecommerce-system/internal/utils/consts"
	"github.com/codepnw/core-ecommerce-system/internal/utils/errs"
	"github.com/codepnw/core-ecommerce-system/internal/utils/security"
)

type IUserService interface {
	CreateUser(ctx context.Context, req *UserCreate) (*User, error)
	CreateUserTx(ctx context.Context, tx *sql.Tx, req *UserCreate) (*User, error)
	GetUser(ctx context.Context, id string) (*User, error)
	GetUsers(ctx context.Context, limit, offset uint) ([]*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, id string, req *UserUpdate) error
	DeleteUser(ctx context.Context, id string) error
}

type userService struct {
	repo IUserRepository
}

func NewUserService(repo IUserRepository) IUserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, req *UserCreate) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		Role:         string(RoleCustomer),
	}
	resp, err := s.repo.Create(ctx, user)
	if err != nil {
		if strings.Contains(err.Error(), "users_email_key") {
			return nil, errs.ErrEmailAlreadyExists
		}
		return nil, err
	}

	return resp, nil
}

func (s *userService) CreateUserTx(ctx context.Context, tx *sql.Tx, req *UserCreate) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		Role:         string(RoleCustomer),
	}
	resp, err := s.repo.CreateTx(ctx, tx, user)
	if err != nil {
		if strings.Contains(err.Error(), "users_email_key") {
			return nil, errs.ErrEmailAlreadyExists
		}
		return nil, err
	}

	return resp, nil
}

func (s *userService) GetUser(ctx context.Context, id string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	return s.repo.GetByID(ctx, id)
}

func (s *userService) GetUsers(ctx context.Context, limit uint, offset uint) ([]*User, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	if limit < 10 {
		limit = 10
	}

	return s.repo.List(ctx, limit, offset)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	return s.repo.GetByEmail(ctx, email)
}

func (s *userService) UpdateUser(ctx context.Context, id string, req *UserUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	update := &User{
		ID: id,
	}
	if req.FullName != nil {
		update.FullName = *req.FullName
	}

	return s.repo.Update(ctx, update)
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	return s.repo.Delete(ctx, id)
}
