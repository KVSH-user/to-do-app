FROM golang:1.21-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext musl-dev

# dependencies
COPY ["go.mod", "go.sum", "./"]
RUN go mod download

# build
COPY ./ ./
RUN go build -o ./bin/app cmd/todo/main.go

FROM alpine AS runner

# Установка базовых пакетов (если ваше приложение их требует)
RUN apk --no-cache add ca-certificates

# Если приложению требуется рабочий каталог, создаем его и задаем как рабочий
WORKDIR /app

# Копируем собранный бинарный файл и конфигурационные файлы в рабочий каталог
COPY --from=builder /usr/local/src/bin/app /app/
COPY config/config.yaml /app/config/config.yaml
COPY db/migrations /app/db/migrations
COPY .env /app/.env

# Указываем CMD с правильным путем к исполняемому файлу
CMD ["/app/app"]

