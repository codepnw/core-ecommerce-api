package main

import (
	"log"

	"github.com/codepnw/core-ecommerce-system/config"
	"github.com/codepnw/core-ecommerce-system/internal/server"
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

	if err = server.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
