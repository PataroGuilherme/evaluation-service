# ---------------------------
# Etapa 1: Build da aplicação
# ---------------------------
FROM golang:1.22 AS builder

# Define o diretório de trabalho
WORKDIR /app

# Copia go.mod e go.sum
COPY go.mod go.sum ./

# Baixa as dependências
RUN go mod download

# Copia todo o código
COPY . .

# Compila a aplicação (binário estático)
RUN CGO_ENABLED=0 GOOS=linux go build -o evaluation-service main.go


# --------------------------------
# Etapa 2: Imagem final (mais leve)
# --------------------------------
FROM alpine:3.20

# Define diretório de trabalho
WORKDIR /app

# Copia o binário da etapa anterior
COPY --from=builder /app/evaluation-service .

# Expõe as portas solicitadas
EXPOSE 8000
EXPOSE 8080
EXPOSE 80
EXPOSE 443

# Comando padrão ao iniciar o container
CMD ["./evaluation-service"]
