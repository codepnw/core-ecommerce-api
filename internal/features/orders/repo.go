package orders

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/codepnw/core-ecommerce-system/internal/utils/errs"
)

type IOrderRepository interface {
	// Table orders
	InsertOrder(ctx context.Context, tx *sql.Tx, input *Order) (int64, error)
	ListOrders(ctx context.Context, filter *OrderFilter) ([]*OrdersResponse, error)
	UpdateStatus(ctx context.Context, orderID int64, status string) error

	// Table order_items
	InsertOrderItems(ctx context.Context, tx *sql.Tx, items []*OrderItem) error

	// Table order_addresses
	InsertOrderAddress(ctx context.Context, tx *sql.Tx, input *OrderAddress) error
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) IOrderRepository {
	return &orderRepository{db: db}
}

// ------------ Table orders ------------

func (r *orderRepository) InsertOrder(ctx context.Context, tx *sql.Tx, input *Order) (int64, error) {
	query := `
		INSERT INTO orders (user_id, address_id, total_price, status)
		VALUES ($1, $2, $3, $4) RETURNING id
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		input.UserID,
		input.AddressID,
		input.TotalPrice,
		input.Status,
	).Scan(&input.ID)
	return input.ID, err
}

func (r *orderRepository) ListOrders(ctx context.Context, filter *OrderFilter) ([]*OrdersResponse, error) {
	var sb strings.Builder
	var args []any

	sb.WriteString(`
		SELECT o.id, u.email, u.full_name, o.total_price, a.phone, a.city, a.state, o.status, o.created_at, o.updated_at
		FROM orders o
		JOIN users u ON u.id = o.user_id
		JOIN addresses a ON a.id = o.address_id
		WHERE 1=1
	`)

	if filter.Status != nil {
		sb.WriteString(fmt.Sprintf(" AND o.status = $%d", len(args)+1))
		args = append(args, *filter.Status)
	}

	if filter.UserID != nil {
		sb.WriteString(fmt.Sprintf(" AND o.user_id = $%d", len(args)+1))
		args = append(args, *filter.UserID)
	}

	sb.WriteString(" ORDER BY o.created_at DESC")

	if filter.Limit != nil {
		sb.WriteString(fmt.Sprintf(" LIMIT $%d", len(args)+1))
		args = append(args, *filter.Limit)
	}

	if filter.Offset != nil {
		sb.WriteString(fmt.Sprintf(" OFFSET $%d", len(args)+1))
		args = append(args, *filter.Offset)
	}

	rows, err := r.db.QueryContext(ctx, sb.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var os []*OrdersResponse
	for rows.Next() {
		o := new(OrdersResponse)
		err = rows.Scan(
			&o.OrderID,
			&o.Email,
			&o.FullName,
			&o.TotalPrice,
			&o.Phone,
			&o.City,
			&o.State,
			&o.Status,
			&o.CreatedAt,
			&o.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		os = append(os, o)
	}

	return os, rows.Err()
}

func (r *orderRepository) UpdateStatus(ctx context.Context, orderID int64, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = now() WHERE id = $2`
	res, err := r.db.ExecContext(ctx, query, status, orderID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrOrderNotFound
	}

	return nil
}

// ------------ Table order_items ------------

func (r *orderRepository) InsertOrderItems(ctx context.Context, tx *sql.Tx, items []*OrderItem) error {
	var cols []string
	var vals []any
	query := `INSERT INTO order_items (order_id, product_id, quantity, price, sub_total) VALUES `

	for i, item := range items {
		cols = append(cols, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5))
		vals = append(vals, item.OrderID, item.ProductID, item.Quantity, item.Price, item.SubTotal)
	}

	query += strings.Join(cols, ", ")
	_, err := tx.ExecContext(ctx, query, vals...)
	return err
}

// ------------ Table order_addresses ------------

func (r *orderRepository) InsertOrderAddress(ctx context.Context, tx *sql.Tx, input *OrderAddress) error {
	query := `
		INSERT INTO order_addresses (order_id, address_id, address_line, city, state, postal_code, phone)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		input.OrderID,
		input.AddressID,
		input.AddressLine,
		input.City,
		input.State,
		input.PostalCode,
		input.Phone,
	)
	return err
}
