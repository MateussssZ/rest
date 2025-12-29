# Этап сборки
FROM golang:1.25-alpine3.23 AS builder

WORKDIR /app

# Устанавливаем зависимости для компиляции
RUN apk add --no-cache git ca-certificates

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/main ./cmd/main.go

# Этап запуска
FROM alpine:latest AS runner

WORKDIR /app

# Устанавливаем зависимости для runtime
RUN apk --no-cache add ca-certificates tzdata curl

# Создаем непривилегированного пользователя
RUN addgroup -g 1001 -S appuser && \
    adduser -u 1001 -S appuser -G appuser

# Копируем бинарник из этапа сборки
COPY --from=builder --chown=appuser:appuser /app/main /app/main

COPY --chown=appuser:appuser .env /app/.en

# Переключаемся на непривилегированного пользователя
USER appuser

# Открываем порт
EXPOSE 8080

# Команда запуска
CMD ["/app/main"]