# Build stage
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

# Копируем go mod файлы
COPY go.mod go.sum ./
RUN go mod tidy

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Копируем бинарник из builder stage
COPY --from=builder /app/main .

# Копируем миграции
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./main"]

