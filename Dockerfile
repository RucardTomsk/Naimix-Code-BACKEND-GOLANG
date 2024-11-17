# syntax=docker/dockerfile:1

# Stage 1: build
FROM golang:1.22-alpine AS builder

# Установка зависимостей
RUN apk add --no-cache libwebp-dev libwebp-tools

# Создаем рабочую директорию
RUN mkdir /build
WORKDIR /build

# Копируем исходный код
COPY . .

# Сборка приложения
RUN GOOS=linux GOARCH=amd64 go build -v -o /build/promitent-api .

# Stage 2: run binary
FROM alpine:3.11


# Установка зависимостей
RUN apk add --no-cache --update libpng-dev libjpeg-turbo-dev giflib-dev tiff-dev autoconf automake make gcc g++ wget
RUN wget https://storage.googleapis.com/downloads.webmproject.org/releases/webp/libwebp-0.6.0.tar.gz && \
    tar -xvzf libwebp-0.6.0.tar.gz && \
    mv libwebp-0.6.0 libwebp && \
    rm libwebp-0.6.0.tar.gz && \
    cd /libwebp && \
    ./configure && \
    make && \
    make install && \
    rm -rf libwebp \

# Открываем порт 80
EXPOSE 80

# Копируем собранный бинарник в финальный образ
COPY --from=builder /build/promitent-api /usr/local/bin/promitent-api

# Копируем конфигурационные файлы и шаблоны
COPY ./internal/common/defaults.yml /usr/local/bin/internal/common/defaults.yml
COPY ./internal/mail/template /usr/local/bin/internal/mail/template
COPY ./config /usr/local/bin/config

# Устанавливаем рабочую директорию
WORKDIR /usr/local/bin

# Устанавливаем точку входа для запуска бинарного файла
ENTRYPOINT ["/usr/local/bin/promitent-api", "dev.yml", "./config"]
