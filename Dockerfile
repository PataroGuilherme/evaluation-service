# ============================
# 1 - Build da aplicação
# ============================
FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o evaluation-service .

# ============================
# 2 - Imagem final
# ============================
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/evaluation-service /app/evaluation-service

EXPOSE 8005

CMD ["/app/evaluation-service"]
