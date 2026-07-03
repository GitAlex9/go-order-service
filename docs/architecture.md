# Architecture

## Overview

O projeto foi desenvolvido seguindo os princípios da Clean Architecture e do Rich Domain Model.

O objetivo é manter o domínio desacoplado de detalhes de infraestrutura, permitindo que a aplicação evolua de um repositório em memória para banco de dados ou APIs externas sem alterar as regras de negócio.

Atualmente a arquitetura está organizada da seguinte forma:

```
cmd/
    app/
        main.go

internal/
    application/
        services/
        id/

    domain/
        entities/
        repositories/
        errors/

    infrastructure/
        memory/

    dto/
```

---

## Camadas

### cmd

Responsável apenas por inicializar a aplicação.

O `main.go` cria as dependências, instancia os serviços e executa os casos de uso.

Nenhuma regra de negócio deve existir nesta camada.

---

### domain

Representa o coração da aplicação.

Contém:

- entidades
- regras de negócio
- contratos (repositories)
- erros de domínio

O domínio não conhece banco de dados, HTTP, JSON ou frameworks.

---

### application

Responsável por coordenar os casos de uso.

Os Services utilizam entidades e repositories através de interfaces.

Nesta camada também ficam componentes reutilizáveis da aplicação, como o gerador de IDs.

---

### infrastructure

Implementações concretas.

No momento:

- repositories em memória utilizando `map`.

No futuro poderão existir implementações para PostgreSQL, MongoDB ou outros bancos sem alterar o domínio.

---

### dto

Objetos utilizados para entrada e saída dos casos de uso.

Atualmente não implementado.

---

## Dependências

A direção das dependências segue o princípio:

```
Infrastructure
        ↓
Application
        ↓
Domain
```

O domínio nunca depende das camadas externas.

---

## Modelo de Domínio

Entidades atuais:

- Product
- Customer
- Order
- OrderItem

Value Objects:

- OrderStatus

Componentes de aplicação:

- CounterGenerator