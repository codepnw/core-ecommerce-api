package orders

import (
	"context"
	"database/sql"
)

type IOrderRepository interface {
	// Table orders
	InsertOrder(ctx context.Context, tx *sql.Tx, input *Order) (int64, error)

	// Table order_items
	InsertOrderItem(ctx context.Context, tx *sql.Tx, input *OrderItem) error

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

// ------------ Table order_items ------------

func (r *orderRepository) InsertOrderItem(ctx context.Context, tx *sql.Tx, input *OrderItem) error {
	query := `
		INSERT INTO order_items (order_id, product_id, quantity, price)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		input.OrderID,
		input.ProductID,
		input.Quantity,
		input.Price,
	)
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
