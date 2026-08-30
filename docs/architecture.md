# Arquitetura

Este projeto segue Clean Architecture / Ports & Adapters com um modelo de domínio rico (Rich Domain Model). As camadas nunca dependem de camadas mais externas — a seta de dependência sempre aponta para dentro, em direção ao domínio.

interfaces/http → application → domain
↑
infrastructure


## Status atual

### ✅ Completo — Domain

- Entidades ricas (`Customer`, `Product`, `Order`, `OrderItem`, `User`), com validação e comportamento próprios, campos privados, construtores (`NewX`) e reconstrutores (`RebuildX`/`RestoreX`)
- Value Objects: `Email`, `CPF`, `Money`, `Password`, `Role`, `OrderStatus`
- Domain Events (`OrderPaidEvent`, `ProductStockDecreasedEvent`, `CustomerRenamedEvent`, etc.), emitidos pelas entidades e acumulados até serem despachados
- Catálogo único de erros de domínio (`domain/errors`)
- Interfaces de repositório (`domain/repositories`), sem dependência de infraestrutura

### ✅ Completo — Application

- Commands e Queries (padrão CQRS simplificado: cada operação é um `XHandler` com um `Handle`, sem uma camada de dispatcher genérico separada)
- DTOs de request/response por entidade
- Validação agregada por caso de uso (`application/validation`), devolvendo todos os erros de campo de uma vez
- Mappers (entidade → DTO)
- Contratos de serviço (`application/contracts`) — interfaces que escondem a orquestração interna
- **Unit of Work**: garante atomicidade entre agregados diferentes (ex: `Order` e `Product` durante a criação/cancelamento de um pedido)
- **Domain Event Dispatcher**: despacha eventos acumulados pelas entidades somente após o `UnitOfWork` confirmar a transação; falha em um handler não interrompe os demais nem propaga para quem disparou o evento
- Autenticação via JWT (geração e validação de token, hash de senha com bcrypt)
- Factory (`application/factory`) — único ponto de injeção de dependências entre `application` e `infrastructure`

### ✅ Completo — Infrastructure

- Repositórios PostgreSQL (via `pgx`), com tradução de erros de banco (violação de unicidade, foreign key) para erros de domínio
- Migrator e Resetter (uso restrito a ambiente de desenvolvimento/teste)
- Configuração via variáveis de ambiente

### ✅ Completo — Cross-cutting

- Logger estruturado (`log/slog`) por trás de uma interface própria (`pkg/logger`), propagado via contexto
- Testes unitários com padrão *table-driven* (`want`/`got`), usando dublês de teste (fakes) em vez de banco real — cobertura de domínio e aplicação acima de 40%

### 🚧 Em andamento

- Reconstrução da camada `interfaces/http` (API REST com Chi) e do binário `cmd/api`

### 📋 Planejado (fora do escopo desta etapa)

- Integração com ViaCEP + Value Object `CEP` (com fallback entre provedores)
- Dockerização completa da aplicação (hoje só o banco roda em container)
- Front-end simples para testes manuais
- Exploração futura, em projeto/etapa separada: decomposição em microsserviços e padrão Saga