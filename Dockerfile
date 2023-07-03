FROM golang:1.19 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

COPY go.* .

RUN go mod download

# Копируем исходный код приложения
COPY . .

# Собираем бинарный файл приложения
RUN go build -o short-url /app/cmd/app/main.go;

EXPOSE 8000

# Запускаем приложение
ENTRYPOINT ["./short-url"]