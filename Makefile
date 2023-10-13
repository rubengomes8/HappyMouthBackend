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

dynamo-create-table-recipes:
	aws dynamodb create-table --table-name GeneratedRecipes --endpoint-url http://localhost:8000 \
	--attribute-definitions AttributeName=PrimaryKey,AttributeType=S \
  	--key-schema AttributeName=PrimaryKey,KeyType=HASH \
  	--provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5
