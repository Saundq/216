FROM golang:alpine as builder 

WORKDIR /build

ADD backend/go.mod .

COPY backend .

RUN go build -o agent cmd/agent/main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/agent /build/agent

COPY backend/.env /build/.env

CMD ["sleep 10"]

CMD ["./agent"]