include app.env

# docker images

run-postgres:
	docker run --name $(DOCKER_DB_IMAGE) --network $(DOCKER_NETWORK) -p 5432:5432 -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASS) -d postgres:14-alpine

# database commands

create-db:
	docker exec -it $(DOCKER_DB_IMAGE) createdb --username=$(DB_USER) --owner=root $(DB_NAME)

drop-db:
	docker exec -it $(DOCKER_DB_IMAGE) dropdb $(DB_NAME)

sqlc:
	./bin/sqlc generate

migrate-new:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_SOURCE)" -verbose up 1

migrate-up-all:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_SOURCE)" -verbose up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_SOURCE)" -verbose down 1

migrate-down-all:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_SOURCE)" -verbose down

# run commands

dev:
	go run cmd/main.go

test:
	go test -cover ./...

.PHONY: dev test run-postgres create-db drop-db sqlc migrate-new migrate-down-all migrate-down migrate-up migrate-up-all