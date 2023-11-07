SWAG_PARAMS = init --parseInternal --parseDependency --parseDepth 3
SWAG_EXCLUDE = --exclude ./db,./docker
SWAG_EXCLUDE_API = $(SWAG_EXCLUDE)

# SWAGGER #
swagger:
	swag $(SWAG_PARAMS) $(SWAG_EXCLUDE_API) -o ./docs -g ./cmd/api/main.go

# GO DEPS #
deps:
	go mod tidy
	go mod vendor

# DOCKER #
docker-up:
	docker-compose -f docker-compose.yaml up --build --detach

docker-down: ## Stop docker containers and clear artefacts.
	docker-compose -f docker-compose.yaml down
	docker system prune

# GENERATE MOCKS #
mock-recipegenerator-service:
	go-mockgen -f ./internal/recipegenerator/ -i service -d ./internal/recipegenerator/mocks/

# DYNAMO DB #
dynamo-up:
	./docker/dynamodb/setup.sh

dynamo-down:
	./docker/dynamodb/clean.sh

# REDIS #
redis-up:
	go run cmd/scripts/populate_cache/main.go

# POSTGRES #
database-migrate:
	DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PWD=passw0rd123 DB_NAME=database DB_MIGRATIONS_PATH=./db/migrations go run ./cmd/dbcli/main.go migrate

database-rollback:
	go run ./cmd/dbcli/main.go rollback

migrate-create:
	migrate create -ext sql -dir db/migrations/ $(MIGRATION)
	ls db/migrations/*.up.sql -r1 | head -n 1 > db/last_migration

