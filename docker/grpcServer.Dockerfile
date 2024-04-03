FROM golang:alpine

RUN apk add build-base

WORKDIR /backend

COPY backend .

RUN go mod download

RUN go build -o /grpcServer cmd/grpcServer/main.go