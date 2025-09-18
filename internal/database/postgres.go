package database

import (
	"database/sql"
	"fmt"
	"log"

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

	log.Printf("database %s connected...", cfg.DB.Name)
	return db, nil
}
