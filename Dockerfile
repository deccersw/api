# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS builder
WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/api ./cmd/api

FROM alpine:3.21
WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/api /app/api


CMD ["/app/api"]