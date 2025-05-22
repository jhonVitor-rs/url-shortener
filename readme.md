# 🔗 Shorter URL - API de Encurtamento de URLs

Uma API moderna e performática para encurtamento de URLs, construída em **Go**, com **PostgreSQL**, **Redis** para cache, e **Swagger** para documentação. Infraestrutura gerenciada com **Docker Compose** e roteamento via **Traefik**.

---

## 🚀 Tecnologias e Ferramentas

- **Go** → linguagem principal.
- **SQLC** → geração de código SQL tipado para acesso ao PostgreSQL.
- **PostgreSQL** → banco de dados relacional.
- **Redis** → cache para acelerar resolução de URLs.
- **Swagger** → documentação interativa da API.
- **Traefik** → proxy reverso, roteamento e segurança.
- **Docker Compose** → orquestração de serviços.
- **Testes de Integração** → cobertura completa da aplicação.

---

## 🏗️ Arquitetura

- **Traefik** → Proxy reverso, roteia requisições com middlewares de segurança e rate limiting.
- **Go API** → Expõe endpoints RESTful para encurtamento e resolução de URLs.
- **PostgreSQL** → Armazena URLs persistentes.
- **Redis** → Cache intermediário, reduz consultas ao banco.
- **Swagger** → Interface de documentação automática disponível via HTTP.

---

## 📂 Estrutura dos Serviços

| Serviço    | Porta          | Descrição                                 |
| ---------- | -------------- | ----------------------------------------- |
| Traefik    | 80, 443, 8080  | Proxy reverso e dashboard                 |
| API (Go)   | 8080 (interno) | Serviço principal de encurtamento de URLs |
| PostgreSQL | 5432           | Banco de dados relacional                 |
| Redis      | 6379           | Armazenamento em cache                    |

---

## ✅ Pré-requisitos

- **Docker** e **Docker Compose** instalados.
- **Go** instalado (opcional, apenas se for rodar/testar sem docker).
- Arquivo `/etc/hosts` atualizado.

---

## 🛠️ Configuração Local

1. **Clone o repositório**:

```bash
git clone https://github.com/seu-usuario/shorter-url.git
cd shorter-url
```

2. **Suba os containers**:

```bash
docker-compose up -d
```

3. **Verifique os serviços:**

- API: [http://app.localhost](http://app.localhost)
- Swagger: [http://app.localhost/swagger/index.html](http://app.localhost/swagger/index.html)
- Traefik Dashboard: [http://traefik.localhost:8080](http://traefik.localhost:8080)

---

## 🚨 Segurança via Traefik

- Middlewares ativos:

  - **Security Headers**: proteção contra vulnerabilidades comuns (XSS, clickjacking, etc.).
  - **Rate Limiting**: limite de requisições por cliente.

Traefik gerencia o roteamento para o serviço Go, aplicando segurança antes de encaminhar o tráfego.

---

## 🧪 Executando Testes

1. Suba os serviços (caso não tenha subido):

```bash
docker-compose up -d
```

2. Execute os testes de integração:

```bash
docker exec -it shorter-url-app go test ./... -v
```

> ⚠️ Os testes incluem integração com PostgreSQL e Redis.

---

## 📖 Documentação (Swagger)

A documentação interativa da API está disponível em:

👉 [http://app.localhost/swagger/index.html](http://app.localhost/swagger/index.html)

**Endpoints principais:**

- `POST /short_url` → Cria uma nova URL encurtada.
- `GET /{slug}` → Redireciona para a URL original.

---

## 🔌 Variáveis de Ambiente

Configuradas no `docker-compose.yml` ou `.env`:

| Variável                                                | Descrição                                                                        |
| ------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `DATABASE_USER`                                         | Nome de usuário para o postgres                                                  |
| `DATABASE_PASSWORD`                                     | Senha do banco de dados                                                          |
| `DATABASE_NAME`                                         | Banco de dados no container postgres                                             |
| `POSTGRES_MULTIPLE_DATABASESPOSTGRES_MULTIPLE_DATABASE` | Banco de dados test, utilizar o mesmo nome do banco de dados com o sufixo \_test |
| `REDIS_PASSWOR`                                         | Senha do redis utilizado para cache                                              |
| `MY_SECRET_KEY`                                         | Secret key utilizada como hash pelo token                                        |

---

## 🌍 Acessos Importantes

| URL                                                                                | Descrição         |
| ---------------------------------------------------------------------------------- | ----------------- |
| [http://app.localhost](http://app.localhost)                                       | API principal     |
| [http://app.localhost/swagger/index.html](http://app.localhost/swagger/index.html) | Swagger           |
| [http://traefik.localhost:8080](http://traefik.localhost:8080)                     | Traefik Dashboard |

---

## ⚙️ Como funciona o roteamento com Traefik?

1. O **Traefik** detecta automaticamente os serviços via **labels** no `docker-compose.yml`.
2. Quando você acessa `http://app.localhost`, o Traefik:

   - Aplica middlewares de segurança.
   - Encaminha para o container **shorter-url-app** na porta **8080**.

3. O mesmo ocorre para o dashboard e outros serviços.

---

## 👥 Contribuindo

1. Fork este repositório.
2. Crie uma branch: `feature/sua-feature`.
3. Faça commit das suas alterações.
4. Envie um Pull Request.

---

## 📫 Contato

- Desenvolvedor: **\[João Vitor]**
- Email: **\[[joaovitor.jvrs6@gmail.com](mailto:joaovitor.jvrs6@gmail.com)]**
- LinkedIn: **\[linkedin.com/in/joão-vitor-rankel-siben-7141a310b]**

---
