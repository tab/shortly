FROM golang:1.24.1-alpine as builder

ENV CGO_ENABLED=0

RUN apk add --no-cache --update git tzdata ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
COPY . ./

ARG BUILD_VERSION="N/A"
ARG BUILD_DATE="N/A"
ARG BUILD_COMMIT="N/A"
ARG VERSION_PACKAGE="shortly/internal/app/version"

RUN go build -ldflags="\
  -s -w \
  -X ${VERSION_PACKAGE}.buildVersion=${BUILD_VERSION} \
  -X ${VERSION_PACKAGE}.buildDate=${BUILD_DATE} \
  -X ${VERSION_PACKAGE}.buildCommit=${BUILD_COMMIT}" \
  -o /app/shortly /app/cmd/shortener/main.go

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/shortly /app/shortly

CMD ["/app/shortly"]
