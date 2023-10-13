#!/bin/bash
aws dynamodb create-table --table-name generated_recipes --endpoint-url http://localhost:8000 \
	--attribute-definitions AttributeName=recipe_key,AttributeType=S \
  	--key-schema AttributeName=recipe_key,KeyType=HASH \
  	--provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5

aws dynamodb create-table --table-name ingredients --endpoint-url http://localhost:8000 \
	--attribute-definitions AttributeName=name,AttributeType=S \
  	--key-schema AttributeName=name,KeyType=HASH \
  	--provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5