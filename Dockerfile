FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o evaluation-service main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/evaluation-service .

# Config Redis via variáveis de ambiente
ENV REDIS_HOST=evaluation-redis
ENV REDIS_PORT=6379

EXPOSE 8000 8080 80 443

CMD ["./evaluation-service"]
