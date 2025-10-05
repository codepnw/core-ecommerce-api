package auth

import (
	"context"
	"database/sql"

	"github.com/codepnw/core-ecommerce-system/internal/database"
	"github.com/codepnw/core-ecommerce-system/internal/utils/errs"
)

type IAuthRepository interface {
	Save(ctx context.Context, exec database.DBExec, input *AuthToken) error
	GetToken(ctx context.Context, userID, token string) (*AuthToken, error)
	Delete(ctx context.Context, userID string) error
}

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) IAuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) Save(ctx context.Context, exec database.DBExec, input *AuthToken) error {
	query := `
		INSERT INTO auth_tokens (user_id, refresh_token, expired_at)
		VALUES ($1, $2, $3) 
		ON CONFLICT (user_id) 
		DO UPDATE SET
			refresh_token = EXCLUDED.refresh_token,
			expired_at = EXCLUDED.expired_at,
			updated_at = NOW()
	`
	_, err := exec.ExecContext(
		ctx,
		query,
		input.UserID,
		input.Token,
		input.ExpiredAt,
	)
	return err
}

func (r *authRepository) GetToken(ctx context.Context, userID, token string) (*AuthToken, error) {
	query := `
		SELECT user_id , refresh_token , expired_at
		FROM auth_tokens
		WHERE user_id = $1 AND refresh_token = $2
		LIMIT 1
	`
	auth := new(AuthToken)
	err := r.db.QueryRowContext(ctx, query, userID, token).Scan(
		&auth.UserID,
		&auth.Token,
		&auth.ExpiredAt,
	)
	if err != nil {
		return nil, err
	}
	return auth, nil
}

func (r *authRepository) Delete(ctx context.Context, userID string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM auth_tokens WHERE user_id = $1", userID)
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
