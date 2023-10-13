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
mock-recipe-service:
	go-mockgen -f ./internal/recipe/ -i service -d ./internal/recipe/mocks/
