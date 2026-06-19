FROM golang:1.26-alpine AS builder

WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Генерируем Swagger
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/server/main.go -o ./docs

# Собираем приложение
RUN go build -o /subscription-server ./cmd/server/main.go

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /subscription-server .
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./subscription-server"]