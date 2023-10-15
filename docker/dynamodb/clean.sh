#!/bin/bash
aws dynamodb delete-table --endpoint-url http://localhost:8000 --table-name generated_recipes
aws dynamodb delete-table --endpoint-url http://localhost:8000 --table-name ingredients
aws dynamodb delete-table --endpoint-url http://localhost:8000 --table-name users