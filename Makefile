HEROKU_APP_NAME=simply-do

help: ## List targets & descriptions
	@cat Makefile* | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build-linux-binary: ## Builds the linux binary
	GO_ENABLED=0 GOOS=linux go build

build-docker: build-linux-binary ## Builds the docker image
	docker build --rm -t simply-do:latest .
	make clean

clean: ## Cleans up the built binaries
	rm simply-do

deploy: build-docker ## Deploys to Heroku. Requires to be logged in on Heroku Registry
	docker tag simply-do:latest registry.heroku.com/$(HEROKU_APP_NAME)/web
	docker push registry.heroku.com/$(HEROKU_APP_NAME)/web
	heroku container:release web --app $(HEROKU_APP_NAME)

run: ## Starts the server locally on port 8080 or $PORT if set
	go run main.go server

.PHONY: help build-linux-binary build-docker clean deploy run
