HEROKU_APP_NAME="simply-do"

build-linux-binary:
	GO_ENABLED=0 GOOS=linux go build

build-docker: build-linux-binary
	docker build --rm -t simply-do:latest .

deploy: build-linux-binary
	docker build --rm -t registry.heroku.com/$(HEROKU_APP_NAME)/web .
	docker push registry.heroku.com/$(HEROKU_APP_NAME)/web

.PHONY: build-linux-binary build-docker deploy
