FROM golang:alpine

RUN apk add build-base

WORKDIR /backend

COPY backend .

RUN go mod download

RUN go build -o /agent cmd/agent/main.go