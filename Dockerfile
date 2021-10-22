
# syntax=docker/dockerfile:1

FROM golang:alpine AS build
ENV GOPROXY=https://proxy.golang.org

WORKDIR /api

COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY . .

RUN go build  ./cmd/micro-videogame
EXPOSE 8080

CMD ["./micro-videogame"]


