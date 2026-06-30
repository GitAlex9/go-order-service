# Go Order Service

Projeto desenvolvido para estudo da linguagem Go utilizando princípios de Clean Architecture e Rich Domain Model.

Este repositório está sendo construído de forma incremental (MVPs) como parte de um desafio acadêmico, priorizando a evolução da arquitetura e das boas práticas ao longo do desenvolvimento.

---

## Objetivos

- Aprender a linguagem Go
- Aplicar os princípios SOLID
- Estudar Clean Architecture
- Modelar um domínio rico (Rich Domain Model)
- Utilizar interfaces para desacoplamento
- Implementar repositories em memória
- Trabalhar com tratamento de erros utilizando `errors.Is`, `errors.Join` e `fmt.Errorf`
- Evoluir o projeto através de MVPs

---

## Tecnologias

- Go
- Git
- GitHub

Sem utilização de:

- HTTP
- Frameworks
- Banco de dados
- Docker

---

## Estrutura do Projeto

```text
cmd/
internal/
docs/
```

A organização segue uma adaptação da Clean Architecture para projetos em Go.

---

## Roadmap

- [x] MVP 0 - Estrutura inicial do projeto
- [ ] MVP 1 - Modelagem do domínio
- [ ] MVP 2 - Erros do domínio
- [ ] MVP 3 - Interfaces dos repositories
- [ ] MVP 4 - Repositories em memória
- [ ] MVP 5 - Application Services
- [ ] MVP 6 - Fluxo principal
- [ ] MVP 7 - Tratamento de erros
- [ ] MVP 8 - Refatoração
- [ ] MVP 9 - Desafios bônus

---

## Arquitetura

O projeto segue uma adaptação da Clean Architecture:

```

main
↓
Application Services
↓
Domain
↓
Repository Interfaces
↓
Infrastructure

```

As entidades do domínio concentram as regras de negócio (Rich Domain Model), enquanto os serviços apenas coordenam os casos de uso.

---

## Status

🚧 Projeto em desenvolvimento.