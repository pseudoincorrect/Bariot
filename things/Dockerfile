FROM golang:alpine3.15

WORKDIR /go/src/app

RUN go install github.com/jackc/tern@latest

COPY go.mod .
COPY go.sum .
# Only download updates if modules files have changed
RUN go mod download
COPY ./ ./
RUN go build main.go
# RUN go get -d -v ./...
RUN go install -v ./...

CMD ["things"]