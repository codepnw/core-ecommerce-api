package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/codepnw/core-ecommerce-system/internal/utils/errs"
)

const (
	selectUserQuery = `
		SELECT id, email, full_name, role, created_at, updated_at
		FROM users
	`
)

type IUserRepository interface {
	Create(ctx context.Context, input *User) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context, limit, offset uint) ([]*User, error)
	Update(ctx context.Context, input *User) error
	Delete(ctx context.Context, id string) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) IUserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, input *User) (*User, error) {
	query := `
		INSERT INTO users (email, password_hash, full_name, role)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		input.Email,
		input.PasswordHash,
		input.FullName,
		input.Role,
	).Scan(
		&input.ID,
		&input.CreatedAt,
		&input.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return input, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	u := new(User)

	query := fmt.Sprintf("%s WHERE email = $1 LIMIT 1", selectUserQuery)
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.FullName,
		&u.Role,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*User, error) {
	u := new(User)

	query := fmt.Sprintf("%s WHERE id = $1 LIMIT 1", selectUserQuery)
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID,
		&u.Email,
		&u.FullName,
		&u.Role,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *userRepository) List(ctx context.Context, limit uint, offset uint) ([]*User, error) {
	query := fmt.Sprintf("%s LIMIT $1 OFFSET $2", selectUserQuery)
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		u := new(User)
		err = rows.Scan(
			&u.ID,
			&u.Email,
			&u.FullName,
			&u.Role,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) Update(ctx context.Context, input *User) error {
	query := `
		UPDATE users SET full_name = $1, updated_at = now() 
		WHERE id = $2
	`
	res, err := r.db.ExecContext(ctx, query, input.FullName, input.ID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrUserNotFound
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrUserNotFound
	}

	return nil
}
