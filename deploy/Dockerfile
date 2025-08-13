# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder

WORKDIR /go/src/service

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd

FROM alpine:latest

WORKDIR /go/src/service

COPY --from=builder /go/src/service/server .

ENV CACHE_SPACE=30\
    SERVER_ADDR=0.0.0.0\
    SERVER_PORT=9090\
    LOGGER_LOG_LEVEL=info

EXPOSE ${SERVER_PORT}

CMD ["./server"]