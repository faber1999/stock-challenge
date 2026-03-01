# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS builder
WORKDIR /src

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/api ./cmd/api

FROM alpine:3.20
WORKDIR /app

RUN apk add --no-cache ca-certificates && addgroup -S app && adduser -S app -G app

COPY --from=builder /out/api /app/api

EXPOSE 8080
USER app
ENTRYPOINT ["/app/api"]
