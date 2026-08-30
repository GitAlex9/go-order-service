# go-order-service

Projeto de estudo em Go, construído para aprender e aplicar na prática **Domain-Driven Design (DDD)**, **Rich Domain Model**, **Clean Architecture** e boas práticas de Go idiomático — do domínio até uma API HTTP completa, passando por testes automatizados.

> Este é um projeto de aprendizado. Decisões de design são documentadas propositalmente em `docs/decisions.md`, incluindo trade-offs e limitações conhecidas — a ideia é que o histórico do projeto (e do repositório) sirva de material de estudo, não só o código final.

## Stack

- **Go** 1.22+
- **PostgreSQL** (via [`pgx`](https://github.com/jackc/pgx))
- **Chi** — roteamento HTTP
- **JWT** (`golang-jwt/jwt`) — autenticação
- **bcrypt** — hash de senha
- **log/slog** — logging estruturado
- **testify-free**: testes unitários com a biblioteca padrão (`testing`), padrão *table-driven* (`want`/`got`)

## Arquitetura

O projeto segue Clean Architecture / Ports & Adapters, com a regra de dependência sempre apontando para dentro:

```
cmd/api, cmd/seed, cmd/app          → binários executáveis
        │
internal/interfaces/http            → HTTP (Chi): handlers, middleware, routes
        │
internal/application                → commands, queries, dto, validation,
        │                              mapper, contracts, services, factory,
        │                              unit of work, event dispatcher
        │
internal/domain                     → entities, value objects, domain events,
        │                              domain errors, repository interfaces
        │                              (não depende de nenhuma camada externa)
        │
internal/infrastructure             → implementação Postgres dos repositórios,
                                       migrator, config de conexão

internal/pkg                        → infraestrutura transversal (jwt, logger)
```

Documentação detalhada, incluindo o histórico de decisões arquiteturais (ADRs), backlog de produto e roadmap de sprints, está em [`docs/`](./docs):

- [`docs/architecture.md`](./docs/architecture.md) — visão geral e status atual por camada
- [`docs/decisions.md`](./docs/decisions.md) — ADRs (Architectural Decision Records)
- [`docs/roadmap.md`](./docs/roadmap.md) — histórico de sprints
- [`docs/backlog.md`](./docs/backlog.md) — evoluções de produto planejadas, fora do escopo atual

## Principais decisões de design (resumo)

- **Modelo rico**: entidades protegem seu próprio estado (`Product` controla estoque, `Order` controla transição de status via `OrderStatus.CanTransitionTo`) — nenhuma regra de negócio vive em um "service anêmico".
- **Value Objects** para conceitos com regra própria: `Email`, `CPF`, `Money` (dinheiro em centavos, nunca `float64`, para evitar erro de arredondamento), `Password` (hash bcrypt, nunca texto puro fora do momento de criação).
- **Unit of Work** garante atomicidade entre agregados diferentes (ex: criar um pedido decrementa estoque de múltiplos produtos — tudo confirma ou tudo desfaz junto).
- **Domain Event Dispatcher**: eventos emitidos pelas entidades (`OrderPaidEvent`, `ProductStockDecreasedEvent`, etc.) são despachados **depois** da transação confirmar, nunca antes.
- **CQRS simplificado**: cada operação é um `XHandler` (`commands`/`queries`), sem uma camada de dispatcher genérico de commands — decisão consciente para o tamanho atual do projeto (ver ADR-007).
- **Autenticação stateless via JWT**, com mensagens de erro de login deliberadamente genéricas (proteção contra enumeração de usuários).

Ver `docs/decisions.md` para o racional completo de cada decisão, incluindo o que foi avaliado e não implementado (ex: refresh token, revogação de token).

## Rodando o projeto

### Pré-requisitos

- Go 1.22+
- Docker (para o PostgreSQL)

### 1. Subir o banco de dados

```bash
docker compose up -d
```

### 2. Configurar variáveis de ambiente

Copie `.env.example` para `.env` e ajuste se necessário:

```bash
cp .env.example .env
```

### 3. Popular dados iniciais (usuário admin + catálogo de produtos)

```bash
go run ./cmd/seed/
```

Credenciais padrão criadas pelo seed (customizáveis via `SEED_ADMIN_EMAIL`/`SEED_ADMIN_PASSWORD`):

| Campo | Valor |
|---|---|
| E-mail | `admin@example.com` |
| Senha | `SenhaForte123!` |
| Role | `admin` |

### 4. Subir a API

```bash
go run ./cmd/api/
```

A API sobe em `http://localhost:8080/api/v1`.

### Testando a API

Um guia completo de testes manuais via `curl`, cobrindo todos os endpoints, está disponível em [`docs/guia-testes-curl.pdf`](./docs/guia-testes-curl.pdf) (ou no arquivo `api-tests.http`, compatível com a extensão REST Client do VS Code).

Fluxo rápido:

```bash
# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"SenhaForte123!"}'

# Usa o token retornado nas próximas requisições
curl http://localhost:8080/api/v1/customers \
  -H "Authorization: Bearer <token>"
```

### Playground manual (`cmd/app`)

Além da API, existe um binário de testes manuais que reseta o banco, roda um fluxo completo (criar/atualizar/listar em todas as entidades, incluindo casos que devem falhar propositalmente) e imprime o resultado no terminal — útil para validar mudanças rapidamente sem precisar de `curl`:

```bash
APP_ENV=development go run ./cmd/app/
```

> ⚠️ `cmd/app` **reseta o banco de dados** a cada execução. Nunca aponte `APP_ENV` para `production` — o `Resetter` se recusa a rodar fora de `development`/`test` por design.

## Testes automatizados

```bash
go test ./... -cover
```

Para ver o relatório de cobertura em HTML:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

Estratégia de testes (ver ADR-012 em `docs/decisions.md`):
- Testes unitários com padrão *table-driven* (`want`/`got`), sem framework externo de asserção
- Nenhuma dependência de banco de dados real nos testes — repositórios são substituídos por dublês (fakes) em memória
- Prioridade de cobertura: domínio (Value Objects, Entities) e regras de aplicação (validation, JWT, dispatcher de eventos) primeiro, por serem lógica pura de maior densidade de regra por linha

## Estrutura de pastas

```
cmd/
  api/      → binário da API HTTP
  seed/     → popula dados iniciais (idempotente)
  app/      → playground de testes manuais (reseta o banco a cada execução)
internal/
  domain/           → entidades, value objects, eventos, erros, interfaces de repositório
  application/       → commands, queries, dto, validation, mapper, contracts,
                        services, factory, unit of work, event dispatcher
  infrastructure/    → implementação Postgres, migrator, config
  interfaces/http/   → handlers, middleware, routes, server (Chi)
  pkg/               → jwt, logger
docs/                → arquitetura, decisões (ADRs), roadmap, backlog, guia de testes
```

## Status

Etapa atual concluída: **monolito modular completo** — domínio rico, camada de aplicação com Unit of Work e Domain Event Dispatcher, API HTTP com autenticação JWT, e suíte de testes automatizados.

Próximos passos planejados estão documentados em [`docs/roadmap.md`](./docs/roadmap.md) e [`docs/backlog.md`](./docs/backlog.md) — incluindo, para uma etapa futura e separada, uma exploração de decomposição em microsserviços com padrão Saga.
