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

