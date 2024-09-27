FROM golang:1.22-alpine as backend

ENV GOFLAGS="-mod=vendor"
ENV CGO_ENABLED=0

RUN apk add --no-cache --update git tzdata ca-certificates

ADD . /build
WORKDIR /build
