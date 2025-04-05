# Dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY ./cmd ./cmd
COPY ./config ./config
COPY ./database ./database
COPY ./errors ./errors
COPY ./handlers ./handlers
COPY ./middlewares ./middlewares
COPY ./models ./models
COPY ./utils ./utils
COPY ./web ./web
COPY go.mod go.sum ./
COPY ./.env ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o ./app ./cmd/server

EXPOSE 7540

CMD ["./app"]