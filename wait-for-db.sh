#!/bin/bash
set -e  # Fail fast
set -o pipefail

WAIT_TIMEOUT=${WAIT_TIMEOUT:-60}

# Verifica se as variáveis estão definidas
: "${DATABASE_HOST:?DATABASE_HOST não definido}"
: "${DATABASE_PORT:?DATABASE_PORT não definido}"
: "${REDIS_HOST:?REDIS_HOST não definido}"

echo "Aguardando banco de dados em ${DATABASE_HOST}:${DATABASE_PORT}..."
./wait-for-it.sh "${DATABASE_HOST}:${DATABASE_PORT}" --timeout="${WAIT_TIMEOUT}" --strict --

echo "Aguardando Redis em ${REDIS_HOST}..."
./wait-for-it.sh "${REDIS_HOST}" --timeout="${WAIT_TIMEOUT}" --strict --

echo "Executando aplicação..."
echo "Gerando código e criando tabelas no banco..."
go generate ./...

if [ $? -ne 0 ]; then
    echo "Erro ao executar go generate. Verifique as configurações do banco de dados."
    exit 1
fi

echo "Iniciando servidor..."
exec ./main