package products

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/codepnw/core-ecommerce-system/internal/features/categories"
	"github.com/codepnw/core-ecommerce-system/internal/utils/errs"
)

const (
	selectProductQuery = `
		SELECT id, category_id, name, description, price, stock, image_url, created_at, updated_at
		FROM products
	`
)

type IProductRepository interface {
	// Products
	Create(ctx context.Context, input *Product) (*Product, error)
	GetByID(ctx context.Context, id int64) (*Product, error)
	List(ctx context.Context, filter *ProductListParams) ([]*Product, error)
	UpdateStock(ctx context.Context, id int64, stock int) error
	Update(ctx context.Context, id int64, input *ProductUpdate) error
	Delete(ctx context.Context, id int64) error

	// Product Categories
	AssignCategory(ctx context.Context, productID, categoryID int64) error
	GetCategoriesByProduct(ctx context.Context, productID int64) ([]*categories.Category, error)
	DelCategoryByProduct(ctx context.Context, productID, categoryID int64) error
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

func (r *productRepository) List(ctx context.Context, filter *ProductListParams) ([]*Product, error) {
	var validateOrderByField = map[string]bool{
		"id":         true,
		"name":       true,
		"price":      true,
		"stock":      true,
		"created_at": true,
	}

	var validateSortField = map[string]bool{
		"asc":  true,
		"desc": true,
	}

	col := "id"
	sort := "DESC"

	if filter.OrderBy != nil && validateOrderByField[*filter.OrderBy] {
		col = *filter.OrderBy
	}

	if filter.Sort != nil && validateSortField[*filter.Sort] {
		sort = *filter.Sort
	}

	var rows *sql.Rows
	var err error

	if filter.CategoryID != 0 {
		query := fmt.Sprintf(`%s WHERE category_id = $1 ORDER BY %s %s LIMIT $2 OFFSET $3`, selectProductQuery, col, sort)
		rows, err = r.db.QueryContext(ctx, query, filter.CategoryID, filter.Limit, filter.Offset)
	} else {
		query := fmt.Sprintf(`%s ORDER BY %s %s LIMIT $1 OFFSET $2`, selectProductQuery, col, sort)
		rows, err = r.db.QueryContext(ctx, query, filter.Limit, filter.Offset)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*Product

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

func (r *productRepository) UpdateStock(ctx context.Context, id int64, stock int) error {
	query := `
		UPDATE products SET stock = stock + $1 
		WHERE id = $2 AND stock + $1 >= 0
	`
	res, err := r.db.ExecContext(ctx, query, stock, id)
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

func (r *productRepository) GetCategoriesByProduct(ctx context.Context, productID int64) ([]*categories.Category, error) {
	query := `
		SELECT c.id, c.name, c.description
		FROM categories c
		JOIN product_categories pc ON pc.category_id = c.id
		WHERE pc.product_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []*categories.Category
	for rows.Next() {
		c := new(categories.Category)
		err = rows.Scan(
			&c.ID,
			&c.Name,
			&c.Description,
		)
		if err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}

	return cats, rows.Err()
}

func (r *productRepository) AssignCategory(ctx context.Context, productID, categoryID int64) error {
	query := `
		INSERT INTO product_categories (product_id, category_id) 
		VALUES ($1, $2)
	`
	res, err := r.db.ExecContext(ctx, query, productID, categoryID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrProductOrCategoryNotFound
	}

	return nil
}

func (r *productRepository) DelCategoryByProduct(ctx context.Context, productID, categoryID int64) error {
	query := `DELETE FROM product_categories WHERE category_id = $1 AND product_id = $2`
	res, err := r.db.ExecContext(ctx, query, categoryID, productID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrProductOrCategoryNotFound
	}

	return nil
}
