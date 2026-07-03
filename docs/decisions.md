# Architectural Decisions

Este documento registra as principais decisões arquiteturais tomadas durante o desenvolvimento do projeto.

---

# ADR-001

## Título

Utilizar Clean Architecture.

## Status

Accepted

## Contexto

O projeto deverá evoluir para diferentes formas de persistência e poderá receber novas interfaces (HTTP, CLI ou gRPC).

## Decisão

Separar domínio, aplicação e infraestrutura.

## Consequências

Positivas:

- baixo acoplamento
- facilidade para testes
- facilidade para trocar infraestrutura

Negativas:

- maior quantidade de arquivos
- curva de aprendizado maior

---

# ADR-002

## Título

Utilizar Rich Domain Model.

## Status

Accepted

## Contexto

As regras de negócio não devem ficar espalhadas pelos Services.

## Decisão

As entidades são responsáveis por proteger seu próprio estado.

Exemplos:

- Product controla seu estoque.
- Order controla sua mudança de status.
- OrderItem calcula seu subtotal.

## Consequências

Positivas:

- domínio coeso
- regras centralizadas
- menor duplicação

---

# ADR-003

## Título

IDs são gerados na camada de Application.

## Status

Accepted

## Contexto

Entidades não devem conhecer mecanismos de geração de identificadores.

## Decisão

Criar um componente reutilizável (`CounterGenerator`) responsável pela geração dos IDs.

## Consequências

Permite substituir facilmente a implementação por UUID, Snowflake ou IDs gerados pelo banco.

---

# ADR-004

## Título

Repositories definidos por interfaces.

## Status

Accepted

## Contexto

Os Services não devem depender de implementações concretas.

## Decisão

Todos os repositories serão definidos como interfaces no domínio.

Implementações concretas ficarão na infraestrutura.

## Consequências

Facilidade para testes.

Facilidade para troca da infraestrutura.