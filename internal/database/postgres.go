package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/codepnw/core-ecommerce-system/config"
	_ "github.com/lib/pq"
)

func ConnectPostgres(cfg *config.EnvConfig) (*sql.DB, error) {
	const dbStr = "postgres://%s:%s@%s:%d/%s?sslmode=disable"

	dsn := fmt.Sprintf(
		dbStr,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres connect failed: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("postgres ping failed: %w", err)
	}

	return db, nil
}

type TxManager struct {
	db *sql.DB
}

func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{db: db}
}

func (m *TxManager) Transaction(ctx context.Context, fn func(tx *sql.Tx) error) (err error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return
}
