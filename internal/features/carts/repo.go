package carts

import (
	"context"
	"database/sql"

	"github.com/codepnw/core-ecommerce-system/internal/utils/errs"
)

type ICartRepository interface {
	GetByUser(ctx context.Context, userID string) ([]*CartItemsResponse, error)
	AddOrUpdate(ctx context.Context, userID string, productID int64, qty int) error
	RemoveItem(ctx context.Context, userID string, productID int64) error
	ClearCart(ctx context.Context, userID string) error
}

type cartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) ICartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) AddOrUpdate(ctx context.Context, userID string, productID int64, qty int) error {
	query := `
		INSERT INTO carts (user_id, product_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, product_id)
		DO UPDATE SET quantity = carts.quantity + EXCLUDED.quantity, updated_at = now()
	`
	_, err := r.db.ExecContext(ctx, query, userID, productID, qty)
	return err
}

func (r *cartRepository) GetByUser(ctx context.Context, userID string) ([]*CartItemsResponse, error) {
	query := `
		SELECT c.product_id, p.name, p.price, c.quantity
		FROM carts c
		JOIN products p ON p.id = c.product_id
		WHERE c.user_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*CartItemsResponse
	for rows.Next() {
		item := new(CartItemsResponse)
		err = rows.Scan(
			&item.ProductID,
			&item.ProductName,
			&item.ProductPrice,
			&item.ProductQuantity,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *cartRepository) RemoveItem(ctx context.Context, userID string, productID int64) error {
	query := "DELETE FROM carts WHERE user_id = $1 AND product_id = $2"
	res, err := r.db.ExecContext(ctx, query, userID, productID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrCartNotFound
	}

	return nil
}

func (r *cartRepository) ClearCart(ctx context.Context, userID string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM carts WHERE user_id = $1", userID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrCartNotFound
	}

	return nil
}
