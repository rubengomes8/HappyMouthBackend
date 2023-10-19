aws dynamodb create-table --endpoint-url http://localhost:8000 --table-name ingredients \
	--attribute-definitions AttributeName=name,AttributeType=S \
  	--key-schema AttributeName=name,KeyType=HASH \
  	--provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5

aws dynamodb create-table --endpoint-url http://localhost:8000  --table-name users \
	--attribute-definitions AttributeName=email,AttributeType=S \
  	--key-schema AttributeName=email,KeyType=HASH \
  	--provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5

aws dynamodb put-item --endpoint-url http://localhost:8000 --table-name ingredients \
	--item file://docker/dynamodb/ingredients/garlic.json

aws dynamodb put-item --endpoint-url http://localhost:8000 --table-name ingredients \
	--item file://docker/dynamodb/ingredients/mushroom.json

aws dynamodb put-item --endpoint-url http://localhost:8000 --table-name ingredients \
	--item file://docker/dynamodb/ingredients/onion.json
	
aws dynamodb put-item --endpoint-url http://localhost:8000 --table-name ingredients \
	--item file://docker/dynamodb/ingredients/potato.json

aws dynamodb put-item --endpoint-url http://localhost:8000 --table-name ingredients \
	--item file://docker/dynamodb/ingredients/tomato.json

aws dynamodb put-item --endpoint-url http://localhost:8000 --table-name ingredients \
	--item file://docker/dynamodb/ingredients/bean.json

aws dynamodb put-item --endpoint-url http://localhost:8000 --table-name ingredients \
	--item file://docker/dynamodb/ingredients/broccoli.json

aws dynamodb put-item --endpoint-url http://localhost:8000 --table-name ingredients \
	--item file://docker/dynamodb/ingredients/cabbage.json

aws dynamodb put-item --endpoint-url http://localhost:8000 --table-name ingredients \
	--item file://docker/dynamodb/ingredients/codfish.json

aws dynamodb put-item --endpoint-url http://localhost:8000 --table-name ingredients \
	--item file://docker/dynamodb/ingredients/ground_beef.json

aws dynamodb put-item --endpoint-url http://localhost:8000 --table-name ingredients \
	--item file://docker/dynamodb/ingredients/lettuce.json

aws dynamodb put-item --endpoint-url http://localhost:8000 --table-name ingredients \
	--item file://docker/dynamodb/ingredients/pasta.json

aws dynamodb put-item --endpoint-url http://localhost:8000 --table-name ingredients \
	--item file://docker/dynamodb/ingredients/red_onion.json

aws dynamodb put-item --endpoint-url http://localhost:8000 --table-name ingredients \
	--item file://docker/dynamodb/ingredients/rice.json

aws dynamodb put-item --endpoint-url http://localhost:8000 --table-name ingredients \
	--item file://docker/dynamodb/ingredients/spinach.json

aws dynamodb put-item --endpoint-url http://localhost:8000 --table-name ingredients \
	--item file://docker/dynamodb/ingredients/tuna.json

aws dynamodb put-item --endpoint-url http://localhost:8000 --table-name ingredients \
	--item file://docker/dynamodb/ingredients/almond.json

aws dynamodb put-item --endpoint-url http://localhost:8000 --table-name ingredients \
	--item file://docker/dynamodb/ingredients/chicken.json

aws dynamodb put-item --endpoint-url http://localhost:8000 --table-name ingredients \
	--item file://docker/dynamodb/ingredients/chocolate.json

aws dynamodb put-item --endpoint-url http://localhost:8000 --table-name ingredients \
	--item file://docker/dynamodb/ingredients/egg.json

aws dynamodb put-item --endpoint-url http://localhost:8000 --table-name ingredients \
	--item file://docker/dynamodb/ingredients/lamb.json

aws dynamodb put-item --endpoint-url http://localhost:8000 --table-name ingredients \
	--item file://docker/dynamodb/ingredients/meat_ball.json

aws dynamodb put-item --endpoint-url http://localhost:8000 --table-name ingredients \
	--item file://docker/dynamodb/ingredients/meat.json

aws dynamodb put-item --endpoint-url http://localhost:8000 --table-name ingredients \
	--item file://docker/dynamodb/ingredients/oregon.json

aws dynamodb put-item --endpoint-url http://localhost:8000 --table-name users \
	--item file://docker/dynamodb/users/ruben.d.s.gomes11@gmail.com.json