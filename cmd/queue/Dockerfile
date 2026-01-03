# Сборка
FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/queue ./cmd/queue/main.go
RUN go build -o /bin/migrator ./cmd/migrator/main.go

# Запуск
FROM alpine
USER root
WORKDIR /home/app
COPY --from=builder /bin/queue ./
COPY --from=builder /bin/migrator ./
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/scripts/docker-entrypoint.sh ./

EXPOSE 8080
ENTRYPOINT ["/bin/sh", "./docker-entrypoint.sh"]
