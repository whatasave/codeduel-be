BINARY_NAME=codeduel-be.exe

build:
	go build -o ./bin/$(BINARY_NAME) -v

run: build
	./bin/$(BINARY_NAME)

dev:
	go run .

test:
	go test -v ./...

docker-build:
	docker build -t codeduel-be .

# docker run -d -p 5000:5000 -v $(PWD)\.env.docker:/.env --name codeduel-be codeduel-be
docker-up:
	docker run -d -p 5000:5000 --name codeduel-be --env-file .env.docker codeduel-be

docker-down:
	docker stop codeduel-be
	docker rm codeduel-be

docker-restart: docker-down docker-up

clean:
	go clean
	rm -f bin/$(BINARY_NAME)
