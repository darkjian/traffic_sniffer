# Используем официальный образ Golang как базовый образ для сборки
FROM golang:1.23.1-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы исходного кода в рабочую директорию контейнера
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем остальные файлы исходного кода
COPY . .

# Собираем Go-приложение
RUN go build -o traffic_sniffer ./cmd/traffic_sniffer.go

# Используем минимальный образ Alpine для финального контейнера
FROM alpine:latest

# Устанавливаем необходимые зависимости
RUN apk add --no-cache iproute2 iputils netcat-openbsd

COPY --from=builder /app/traffic_sniffer /usr/local/bin/traffic_sniffer

COPY setup_and_test.sh /usr/local/bin/setup_and_test.sh

RUN chmod +x /usr/local/bin/setup_and_test.sh

RUN mkdir -p /logs

ENTRYPOINT ["traffic_sniffer"]

