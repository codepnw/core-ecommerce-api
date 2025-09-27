package addresses

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/codepnw/core-ecommerce-system/internal/utils/errs"
)

const (
	selectAddress = `
		SELECT id, user_id, address_line, city, state, postal_code, phone, created_at, updated_at
		FROM addresses
	`
)

type IAddressRepository interface {
	Create(ctx context.Context, input *Address) error
	GetByID(ctx context.Context, id string) (*Address, error)
	List(ctx context.Context, userID string) ([]*Address, error)
	Update(ctx context.Context, id string, input *AddressUpdate) error
	Delete(ctx context.Context, id string) error
	SetDefault(ctx context.Context, addressID, userID string) error
}

type addressRepository struct {
	db *sql.DB
}

func NewAddressRepository(db *sql.DB) IAddressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) Create(ctx context.Context, input *Address) error {
	query := `
		INSERT INTO addresses (user_id, address_line, city, state, postal_code, phone)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		input.UserID,
		input.AddressLine,
		input.City,
		input.State,
		input.PostalCode,
		input.Phone,
	)
	return err
}

func (r *addressRepository) GetByID(ctx context.Context, id string) (*Address, error) {
	address := new(Address)
	query := fmt.Sprintf("%s WHERE id = $1 LIMIT 1", selectAddress)
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&address.ID,
		&address.UserID,
		&address.AddressLine,
		&address.City,
		&address.State,
		&address.PostalCode,
		&address.Phone,
		&address.CreatedAt,
		&address.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrAddressNotFound
		}
		return nil, err
	}
	return address, nil
}

func (r *addressRepository) List(ctx context.Context, userID string) ([]*Address, error) {
	query := fmt.Sprintf("%s WHERE user_id = $1 ORDER BY created_at DESC", selectAddress)
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var adds []*Address
	for rows.Next() {
		a := new(Address)
		err = rows.Scan(
			&a.ID,
			&a.UserID,
			&a.AddressLine,
			&a.City,
			&a.State,
			&a.PostalCode,
			&a.Phone,
			&a.CreatedAt,
			&a.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		adds = append(adds, a)
	}

	return adds, rows.Err()
}

func (r *addressRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM addresses WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrAddressNotFound
	}

	return nil
}

func (r *addressRepository) Update(ctx context.Context, id string, input *AddressUpdate) error {
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
		return errs.ErrAddressNotFound
	}

	return nil
}

func (r *addressRepository) buildUpdateQuery(id string, input *AddressUpdate) (string, []any, error) {
	columns := []string{}
	args := []any{}
	idx := 1

	if input.AddressLine != nil {
		columns = append(columns, fmt.Sprintf("address_line = $%d", idx))
		args = append(args, *input.AddressLine)
		idx++
	}

	if input.City != nil {
		columns = append(columns, fmt.Sprintf("city = $%d", idx))
		args = append(args, *input.City)
		idx++
	}

	if input.State != nil {
		columns = append(columns, fmt.Sprintf("state = $%d", idx))
		args = append(args, *input.State)
		idx++
	}

	if input.PostalCode != nil {
		columns = append(columns, fmt.Sprintf("postal_code = $%d", idx))
		args = append(args, *input.PostalCode)
		idx++
	}

	if input.Phone != nil {
		columns = append(columns, fmt.Sprintf("phone = $%d", idx))
		args = append(args, *input.Phone)
		idx++
	}

	if len(columns) == 0 {
		return "", nil, errs.ErrNoFieldUpdate
	}

	setColumns := strings.Join(columns, ", ")
	args = append(args, id)

	query := fmt.Sprintf("UPDATE addresses SET %s, updated_at = NOW() WHERE id = $%d", setColumns, idx)
	return query, args, nil
}

func (r *addressRepository) SetDefault(ctx context.Context, addressID, userID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Reset Default
	_, err = tx.ExecContext(
		ctx,
		"UPDATE addresses SET is_default = false WHERE user_id = $1",
		userID,
	)
	if err != nil {
		return err
	}

	// Set New Default
	res, err := tx.ExecContext(
		ctx,
		"UPDATE addresses SET is_default = true WHERE id = $1 AND user_id = $2",
		addressID,
		userID,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrAddressNotFound
	}

	return tx.Commit()
}
