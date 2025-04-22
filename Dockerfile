# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Instalar dependencias necesarias
RUN apk add --no-cache git make

# Copiar archivos de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

# Instalar certificados CA
RUN apk --no-cache add ca-certificates

# Copiar el binario compilado
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Exponer el puerto
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./main"] 