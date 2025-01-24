FROM golang:1.22-alpine as builder

ENV CGO_ENABLED=0

RUN apk add --no-cache --update git tzdata ca-certificates

WORKDIR /app

COPY go.mod go.sum ./

COPY . ./
RUN go build -o /app/shortly /app/cmd/shortener/main.go

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/shortly /app/shortly

CMD ["/app/shortly"]
