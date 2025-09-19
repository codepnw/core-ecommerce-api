package products

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/codepnw/core-ecommerce-system/internal/utils/errs"
)

type IProductRepository interface {
	Create(ctx context.Context, input *Product) (*Product, error)
	GetByID(ctx context.Context, id int64) (*Product, error)
	List(ctx context.Context, limit, offset uint) ([]*Product, error)
	Update(ctx context.Context, id int64, input *ProductUpdate) error
	Delete(ctx context.Context, id int64) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) IProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, input *Product) (*Product, error) {
	query := `
		INSERT INTO products (category_id, name, description, price, stock, image_url)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		input.CategoryID,
		input.Name,
		input.Description,
		input.Price,
		input.Stock,
		input.ImageURL,
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

const selectProductQuery = `
	SELECT id, category_id, name, description, price, stock, image_url, created_at, updated_at
	FROM products
`

func (r *productRepository) GetByID(ctx context.Context, id int64) (*Product, error) {
	p := new(Product)
	query := fmt.Sprintf("%s WHERE id = $1", selectProductQuery)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.CategoryID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Stock,
		&p.ImageURL,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrProductNotFound
		}
		return nil, err
	}

	return p, nil
}

func (r *productRepository) List(ctx context.Context, limit, offset uint) ([]*Product, error) {
	var products []*Product

	query := fmt.Sprintf("%s LIMIT $1 OFFSET $2", selectProductQuery)
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		p := new(Product)
		err = rows.Scan(
			&p.ID,
			&p.CategoryID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.Stock,
			&p.ImageURL,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrProductNotFound
	}

	return nil
}

func (r *productRepository) Update(ctx context.Context, id int64, input *ProductUpdate) error {
	query, args, err := r.buildUpdateQuery(id, input)
	if err != nil {
		return err
	}
	log.Println(query)

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrProductNotFound
	}

	return nil
}

func (r *productRepository) buildUpdateQuery(id int64, p *ProductUpdate) (string, []any, error) {
	var columns []string
	var args []any
	idx := 1

	if p.CategoryID != nil {
		columns = append(columns, fmt.Sprintf("category_id = $%d", idx))
		args = append(args, p.CategoryID)
		idx++
	}

	if p.Name != nil {
		columns = append(columns, fmt.Sprintf("name = $%d", idx))
		args = append(args, p.Name)
		idx++
	}

	if p.Description != nil {
		columns = append(columns, fmt.Sprintf("description = $%d", idx))
		args = append(args, p.Description)
		idx++
	}

	if p.Price != nil {
		columns = append(columns, fmt.Sprintf("price = $%d", idx))
		args = append(args, p.Price)
		idx++
	}

	if p.Stock != nil {
		columns = append(columns, fmt.Sprintf("stock = $%d", idx))
		args = append(args, p.Stock)
		idx++
	}

	if p.ImageURL != nil {
		columns = append(columns, fmt.Sprintf("image_url = $%d", idx))
		args = append(args, p.ImageURL)
		idx++
	}

	if len(columns) == 0 {
		return "", nil, errs.ErrNoFieldUpdate
	}

	setColumns := strings.Join(columns, ", ")
	query := fmt.Sprintf("UPDATE products SET %s, updated_at = NOW() WHERE id = $%d", setColumns, idx)
	args = append(args, id)

	return query, args, nil
}
