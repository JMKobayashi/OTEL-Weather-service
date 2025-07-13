FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copiar go.mod e go.sum
COPY service-b/go.mod service-b/go.sum ./

# Baixar dependências
RUN go mod download

# Copiar código fonte
COPY service-b/ .

# Compilar a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Imagem final
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar o binário compilado
COPY --from=builder /app/main .

# Copiar arquivo .env da raiz do projeto
COPY .env .env

# Expor porta
EXPOSE 8081

# Comando para executar a aplicação
CMD ["./main"] 