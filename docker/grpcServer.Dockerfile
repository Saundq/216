FROM golang:alpine as builder 

WORKDIR /build

ADD backend/go.mod .

COPY backend .

RUN go build -o grpcServer cmd/grpcServer/main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/grpcServer /build/grpcServer

COPY backend/.env /build/.env

CMD ["./grpcServer"]