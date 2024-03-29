FROM golang:1.21 as build-stage

ENV BINARY_NAME=codeduel-be
ENV ENV=production

RUN useradd -u 1001 -m codeduel-user

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/$BINARY_NAME -v


FROM build-stage AS run-test-stage
RUN go test -v ./...

# FROM alpine:3.14
# FROM scratch
FROM gcr.io/distroless/base-debian11 AS release-stage

# RUN apk add --no-cache ca-certificates
COPY --from=build-stage /usr/src/app/bin /usr/local/bin
# COPY --from=build-stage /usr/src/app/bin/.env /.env
COPY --from=build-stage /etc/passwd /etc/passwd

USER 1001
# USER codeduel-user:codeduel-user
EXPOSE 5000

ENTRYPOINT ["codeduel-be"]
