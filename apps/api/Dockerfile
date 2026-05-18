FROM golang:1.22-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/runtimewall-api ./cmd/server

FROM alpine:3.20

RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /out/runtimewall-api /app/runtimewall-api

ENV HTTP_ADDR=:8080
EXPOSE 8080

ENTRYPOINT ["/app/runtimewall-api"]
