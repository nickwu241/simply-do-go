HEROKU_APP_NAME=simply-do

help: ## List targets & descriptions
	@cat Makefile* | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build-client: ## Builds the frontend files to public
	@sh -c 'cd client && npm run build && mv dist public && mv public ../'

build-linux-binary: ## Builds the linux binary
	GO_ENABLED=0 GOOS=linux go build

build-docker: build-client build-linux-binary ## Builds the docker image
	docker build --rm -t simply-do:latest .
	make clean

clean: ## Cleans up the built binaries
	rm -rf public
	rm -f simply-do

deploy: build-docker ## Deploys to Heroku. Requires to be logged in on Heroku Registry
	docker tag simply-do:latest registry.heroku.com/$(HEROKU_APP_NAME)/web
	docker push registry.heroku.com/$(HEROKU_APP_NAME)/web
	heroku container:release web --app $(HEROKU_APP_NAME)

run: ## Starts the server locally on port 8080 or $PORT if set
	go run main.go server

run-docker: ## Starts the server with the latest built docker image. Requires .env.secret
	docker run --rm -it -p 8080:8080 --env-file .env.secret simply-do:latest

.PHONY: help build-client build-linux-binary build-docker clean deploy run
