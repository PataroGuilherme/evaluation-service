FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o evaluation-service .

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/evaluation-service .

EXPOSE 8000 8080 80 443

CMD ["./evaluation-service"]
