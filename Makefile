include dev.env

ENV_FILE = dev.env
MIGRATE_PATH = internal/database/migrations
MIGRATE_DB = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

docker-up:
	@docker compose --env-file $(ENV_FILE) up -d

migrate-create:
	@migrate create -ext sql -dir $(MIGRATE_PATH) -seq $(name)

migrate-up:
	@migrate -database $(MIGRATE_DB) -path $(MIGRATE_PATH) up

migrate-down:
	@migrate -database $(MIGRATE_DB) -path $(MIGRATE_PATH) down

migrate-force:
	@migrate -database $(MIGRATE_DB) -path $(MIGRATE_PATH) force 1