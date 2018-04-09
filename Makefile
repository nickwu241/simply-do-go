HEROKU_APP_NAME=simply-do

build-linux-binary:
	GO_ENABLED=0 GOOS=linux go build

build-docker: build-linux-binary
	docker build --rm -t simply-do:latest .
	make clean

deploy: build-docker
	docker tag simply-do:latest registry.heroku.com/$(HEROKU_APP_NAME)/web
	docker push registry.heroku.com/$(HEROKU_APP_NAME)/web

clean:
	rm simply-do

.PHONY: run build-linux-binary build-docker deploy clean
