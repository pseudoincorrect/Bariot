# syntax=docker/dockerfile:1

# base step
FROM golang:alpine3.15 as base

WORKDIR /go/src/service_1
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Development step
FROM base as dev

COPY . /go/src/service_1
RUN ls
RUN go build -o /service_1

ENTRYPOINT ["/service_1"]