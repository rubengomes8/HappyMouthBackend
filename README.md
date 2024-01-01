# HappyMouthBackend

This is the backend that supports this [react native application](https://github.com/rubengomes8/HappyMouth). 

The goal of the app is to allow a user to create a recipe by providing some informations, such as the included ingredients, the ingredients to exclude from the recipe and the recipe type. 

This backend comprises a PostgreSQL database, a Redis cache and a REST API that does the following:
1. handles the user authentication (signup / login / change password)
2. handles the ingredients requests
3. handles the recipes requests


### Demo
Check the app demo video [here](quick_demo.mov)!

## Setup
1. `make dynamo-down`
2. `make docker-down`
3. `make docker-up`
4. `make dynamo-up`

## TODOs
- [x] Add time.Sleep in mocked response
- [x] Add ingredients database
- [x] Add authentication
- [x] Add users database
- [x] Use ~MemCached or~ Redis to store the GPT Recipes instead of Kafka
