FROM golang:alpine3.15 AS packages
ARG SVC

WORKDIR /go/src/github.com/pseudoincorrect/bariot
COPY ../go.mod .
COPY ../go.sum .

RUN go mod download

FROM packages as builder

COPY ../. .

WORKDIR /go/src/github.com/pseudoincorrect/bariot/build
WORKDIR /go/src/github.com/pseudoincorrect/bariot/cmd/$SVC

RUN go build -o /exec main.go

WORKDIR /go/src/github.com/pseudoincorrect/bariot/build

ENTRYPOINT ["/exec"]

# Build example for "users" service (use "--no-cache" for full rebuild)
# docker build --progress=plain -t users_service --build-arg SVC=users -f Dockerfile ..  