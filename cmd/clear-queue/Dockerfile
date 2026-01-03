# Сборка
FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/clear-queue ./cmd/clear-queue/main.go

# Запуск
FROM alpine
USER root
WORKDIR /home/app
COPY --from=builder /bin/clear-queue ./

ENTRYPOINT ["./clear-queue"]
