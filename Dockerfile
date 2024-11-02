FROM golang:1.22-alpine as base-backend

ENV GOFLAGS="-mod=vendor"
ENV CGO_ENABLED=0

RUN apk add --no-cache --update git tzdata ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor ./vendor

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./store.json ./store.json

RUN go build -o /app/shortly /app/cmd/shortener/main.go

CMD ["/app/shortly"]
