# This template is used to build every Bariot docker images.
# to build and run a service (for development)
# $ docker-compose up --build {service_name} 

FROM golang:alpine3.15 AS packages
ARG SVC

WORKDIR /go/src/github.com/pseudoincorrect/bariot

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

FROM packages as builder
ARG SVC

COPY ./. .

WORKDIR /go/src/github.com/pseudoincorrect/bariot/build
WORKDIR /go/src/github.com/pseudoincorrect/bariot/${SVC}/main

RUN go build -o /exec ${SVC}_main.go

WORKDIR /go/src/github.com/pseudoincorrect/bariot/build

ENTRYPOINT ["/exec"]

# Build example for "users" service (use "--no-cache" for full rebuild)
# docker build --progress=plain -t users_service --build-arg SVC=users -f Dockerfile ..  