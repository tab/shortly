FROM golang:1.22-alpine as base-backend

ENV GOFLAGS="-mod=vendor"
ENV CGO_ENABLED=0

RUN apk add --no-cache --update git tzdata ca-certificates

WORKDIR /backend

COPY go.mod go.sum ./
COPY vendor ./vendor

COPY ./cmd ./cmd
COPY ./internal ./internal

RUN go build -o /backend/shortly /backend/cmd/shortener/main.go

CMD ["/backend/shortly"]
