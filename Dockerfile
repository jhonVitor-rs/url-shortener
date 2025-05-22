FROM golang:1.24-alpine

RUN apk add --no-cache git bash netcat-openbsd

WORKDIR /app

# Instalar git, tern e sqlc
RUN apk add --no-cache git bash \
  && go install github.com/jackc/tern@latest \
  && go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Copiar arquivos de dependências primeiro (para cache do Docker)
COPY go.mod go.sum ./
RUN go mod download

# Copiar o resto dos arquivos
COPY . .

RUN go mod tidy 
RUN chmod +x /app/wait-for-it.sh \
  && chmod +x /app/wait-for-db.sh \
  && chmod +x init-multiple-dbs.sh 


# Construir o binário Go
RUN go build -o main ./cmd/server/main.go

EXPOSE 8080

# Usar o script de entrada
CMD ["./wait-for-db.sh"]