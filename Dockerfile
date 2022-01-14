## Taken from https://docs.docker.com/language/golang/build-images/
FROM golang:1.17-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go install github.com/githubnemo/CompileDaemon@latest

COPY . ./

RUN go build -o /main

ENTRYPOINT CompileDaemon -log-prefix=false -build="go build -o /main" -command="/main"
