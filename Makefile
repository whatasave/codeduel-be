BINARY_NAME=codeduel-be.exe
DOCKERHUB_USERNAME=xedom
DOCKER_IMAGE_NAME=codeduel-be
DOCKER_CONTAINER_NAME=codeduel-be
PWD=$(CURDIR)

build:
	go build -o ./bin/$(BINARY_NAME) -v

run: build
	./bin/$(BINARY_NAME)

dev:
	swag fmt -g api/api.go
	swag init -g api/api.go
	go run .

test:
	go test -v ./...

gen-ssl:
	mkdir -p ssl
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ssl/server.key -out ssl/server.crt -subj "/C=US/ST=State/L=City/O=Organization/OU=Department/CN=codeduel.it"

docker-build:
	docker build -t $(DOCKERHUB_USERNAME)/$(DOCKER_IMAGE_NAME) .

docker-push:
	docker push $(DOCKERHUB_USERNAME)/$(DOCKER_IMAGE_NAME)

docker-up:
	docker run --name $(DOCKER_CONTAINER_NAME) \
		-p 5000:443 \
		-p 5001:80 \
		--env-file .env.docker \
		-v $(PWD)/ssl:/ssl \
		$(DOCKERHUB_USERNAME)/$(DOCKER_IMAGE_NAME)

docker-down:
	-docker stop $(DOCKER_CONTAINER_NAME)
	-docker rm $(DOCKER_CONTAINER_NAME)

docker-restart: docker-down docker-up

release:
	git checkout release
	git merge main
	git push origin release
	git checkout main

clean:
	go clean
	go mod tidy
	-rm -f bin/$(BINARY_NAME)
