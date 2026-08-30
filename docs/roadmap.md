# Roadmap

## Sprint 1 ✅
- Estrutura inicial do projeto
- Documentação
- Repositório GitHub

## Sprint 2 ✅
- Rich Domain Model
- Domain Errors
- Value Objects (`Email`, `CPF`, `Money`, `Password`)
- Repository Interfaces

## Sprint 3 ✅
- Docker (PostgreSQL)
- Conexão com banco (`pgx`)
- Migrator / Resetter
- Implementação dos repositórios com `pgx`

## Sprint 4 ✅
- Application layer completa: Commands, Queries, DTOs, Validation, Mappers, Contracts, Services, Factory
- Autenticação JWT (login, hash de senha, middleware de autenticação e de papel/role)
- Unit of Work (consistência entre `Order` e `Product`)
- Domain Event Dispatcher (Registry + Dispatcher + handlers de exemplo)
- Logger estruturado (`log/slog`)

## Sprint 5 ✅
- Suíte de testes unitários (domínio, aplicação, infraestrutura), com dublês de teste em vez de banco real
- Cobertura de testes acima de 40% do projeto

## Sprint 6 🚧 (etapa atual)
- Reconstrução da camada HTTP (`interfaces/http`) com Chi Router
- Novo binário `cmd/api`
- Fechamento da documentação de arquitetura

## Sprint 7 (planejado, fora desta etapa)
- Dockerização completa da aplicação (API + banco)
- Front-end simples para testes manuais da API

## Sprint 8 (exploração futura, projeto separado)
- Avaliação de decomposição em microsserviços por agregado
- Padrão Saga para consistência entre serviços, substituindo o Unit of Work local