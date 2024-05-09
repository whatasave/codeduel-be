FROM golang:1.22 as build-stage

ENV GO_ENV=production
ENV ENV=production

RUN useradd -u 1001 -m codeduel-user

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/codeduel-be -v

FROM build-stage AS run-test-stage
RUN go test -v ./...


FROM gcr.io/distroless/base-debian11 AS release-stage

COPY --from=build-stage /usr/src/app/bin /usr/local/bin
COPY --from=build-stage /etc/passwd /etc/passwd

ENV ENV="production"
ENV COOKIE_DOMAIN="*"
ENV CORS_ORIGIN="*"
ENV HOST=0.0.0.0
ENV PORT=443
ENV PORT_HTTP=80
ENV SSL_KEY="/ssl/server.pem"
ENV SSL_CERT="/ssl/server.pem"

USER 1001

EXPOSE 443
EXPOSE 80

VOLUME [ "/ssl/" ]

ENTRYPOINT ["codeduel-be"]
