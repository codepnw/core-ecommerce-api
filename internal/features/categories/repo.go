package categories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/codepnw/core-ecommerce-system/internal/utils/errs"
)

type CategoryRepository interface {
	Create(ctx context.Context, input *Category) error
	List(ctx context.Context) ([]*Category, error)
	Update(ctx context.Context, input *Category) error
	Delete(ctx context.Context, id int64) error
}

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, input *Category) error {
	query := `INSERT INTO categories (name, description) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, input.Name, input.Description)
	if err != nil {
		return err
	}
	return nil
}

func (r *categoryRepository) List(ctx context.Context) ([]*Category, error) {
	query := `SELECT id, name, description FROM categories`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var categories []*Category

	for rows.Next() {
		c := new(Category)
		err = rows.Scan(
			&c.ID,
			&c.Name,
			&c.Description,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepository) Update(ctx context.Context, input *Category) error {
	columns := []string{}
	args := []any{}
	idx := 1

	if input.Name != "" {
		columns = append(columns, fmt.Sprintf("name = $%d", idx))
		args = append(args, input.Name)
		idx++
	}

	if input.Description != "" {
		columns = append(columns, fmt.Sprintf("description = $%d", idx))
		args = append(args, input.Description)
		idx++
	}

	if len(columns) == 0 {
		return errs.ErrNoFieldUpdate
	}

	setColumns := strings.Join(columns, ", ")
	query := fmt.Sprintf(`UPDATE categories SET %s, updated_at = now() WHERE id = $%d`, setColumns, idx)
	args = append(args, input.ID)

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrCategoryNotFound
	}

	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrCategoryNotFound
	}

	return nil
}
