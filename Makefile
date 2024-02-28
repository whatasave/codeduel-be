BINARY_NAME=codeduel-be.exe
DOCKERHUB_USERNAME=xedom
DOCKER_IMAGE_NAME=codeduel-be
DOCKER_CONTAINER_NAME=codeduel-be

build:
	go build -o ./bin/$(BINARY_NAME) -v

run: build
	./bin/$(BINARY_NAME)

dev:
	go run .

test:
	go test -v ./...

docker-build:
	docker build -t $(DOCKERHUB_USERNAME)/$(DOCKER_IMAGE_NAME) .

docker-push:
	docker push $(DOCKERHUB_USERNAME)/$(DOCKER_IMAGE_NAME)

# docker run -d -p 5000:5000 -v $(PWD)\.env.docker:/.env --name $(DOCKER_CONTAINER_NAME) $(DOCKERHUB_USERNAME)/$(DOCKER_IMAGE_NAME)
docker-up:
	docker run -d -p 5000:5000 --name $(DOCKER_CONTAINER_NAME) --env-file .env.docker $(DOCKERHUB_USERNAME)/$(DOCKER_IMAGE_NAME)

docker-down:
	docker stop $(DOCKER_CONTAINER_NAME)
	docker rm $(DOCKER_CONTAINER_NAME)

docker-restart: docker-down docker-up

release:
	git checkout release
	git merge main
	git push origin release
	git checkout main

clean:
	go clean
	rm -f bin/$(BINARY_NAME)
