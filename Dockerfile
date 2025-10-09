# ---------- STAGE 1: BUILD ----------
FROM golang:1.24 AS builder

WORKDIR /app

# Копируем весь проект
COPY . .

# Статическая сборка бинарника
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/subservice-bin ./subservice/cmd

# ---------- STAGE 2: RUNTIME ----------
FROM alpine:latest

# Устанавливаем сертификаты для HTTPS и curl
RUN apk add --no-cache ca-certificates curl bash

WORKDIR /app

# Копируем бинарник и нужные директории
COPY --from=builder /app/subservice-bin .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/subservice/internal/config ./internal/config
COPY --from=builder /app/subservice/internal/migrations /internal/migrations

# Скачиваем wait-for-it.sh для ожидания готовности базы
ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

EXPOSE 8080

# CMD: ждём Postgres, затем запускаем сервис
CMD ["/wait-for-it.sh", "database:5432", "--timeout=30", "--strict", "--", "./subservice-bin"]
