# Этап 1: Сборка (Builder)
FROM golang:1.26.2-alpine AS builder

# Настройка рабочей директории
WORKDIR /app

# Копирование файлов модулей и скачивание зависимостей (кеширование слоев)
COPY go.mod ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка бинарного файла. -ldflags="-s -w" уменьшает размер файла
RUN go build -o link-service cmd/main.go

# Этап 2: Финальный образ (Final)
FROM alpine:latest

# # Установка сертификатов безопасности (если нужно обращение по HTTPS)
# RUN apk --no-cache add ca-certificates

WORKDIR /root/

# # Копирование бинарника из этапа сборки
COPY --from=builder /app/link-service .
COPY --from=builder /app/migrations ./migrations

# Запуск приложения
CMD ["./link-service"]