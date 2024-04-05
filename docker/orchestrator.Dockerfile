FROM golang:alpine as builder 

WORKDIR /build

ADD backend/go.mod .

COPY backend .

RUN go build -o orchestrator cmd/orchestrator/main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/orchestrator /build/orchestrator

COPY backend/.env /build/.env

CMD ["./orchestrator"]