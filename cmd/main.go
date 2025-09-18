package main

import (
	"log"

	"github.com/codepnw/core-ecommerce-system/config"
	"github.com/codepnw/core-ecommerce-system/internal/database"
	"github.com/joho/godotenv"
)

const envFile = "dev.env"

func init() {
	if err := godotenv.Load(envFile); err != nil {
		log.Fatalf("load env failed: %v", err)
	}
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.ConnectPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
