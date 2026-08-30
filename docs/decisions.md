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
Positivas: baixo acoplamento, facilidade para testes, facilidade para trocar infraestrutura.
Negativas: maior quantidade de arquivos, curva de aprendizado maior.

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

Exemplos: `Product` controla seu estoque; `Order` controla sua mudança de status via `OrderStatus.CanTransitionTo`; `OrderItem` calcula seu subtotal; `Role`/`OrderStatus` validam a si mesmos (`IsValid()`).

## Consequências
Positivas: domínio coeso, regras centralizadas, menor duplicação.

---

# ADR-003

## Título
IDs gerados na camada de Application via `CounterGenerator`.

## Status
**Superseded por ADR-013**

## Contexto
Decisão original, tomada antes da implementação das entidades.

## Decisão original
Um componente reutilizável de geração de identificadores, na camada de aplicação.

## Motivo da revisão
Na prática, adotou-se `uuid.New()` gerado dentro do próprio construtor de cada entidade (`NewCustomer`, `NewProduct`, `NewOrder`, `NewUser`), não na aplicação. Ver ADR-013.

---

# ADR-004

## Título
Repositories definidos por interfaces.

## Status
Accepted

## Contexto
Os Services não devem depender de implementações concretas.

## Decisão
Todos os repositories são definidos como interfaces no domínio (`domain/repositories`). Implementações concretas ficam na infraestrutura (`infrastructure/repositories/postgres`).

## Consequências
Facilidade para testes (dublês/fakes) e para troca de infraestrutura.

---

# ADR-005

## Título
Dinheiro representado como Value Object `Money`, em centavos (`int64`).

## Status
Accepted

## Contexto
Representar valores monetários como `float64` introduz erro de arredondamento em somas e multiplicações sucessivas.

## Decisão
`Money` armazena o valor internamente em centavos (`int64`). Conversão para/de `float64` só acontece na fronteira (DTOs de entrada/saída), nunca dentro do domínio.

## Consequências
Positivas: aritmética exata, sem erro de ponto flutuante.
Negativas: precisa de conversão explícita (`Amount()`/`NewMoneyFromFloat()`) sempre que dinheiro atravessa a fronteira domínio ↔ mundo externo.

---

# ADR-006

## Título
Senha representada como Value Object `Password`, com hash via bcrypt.

## Status
Accepted

## Contexto
Senha em texto puro nunca deve ser persistida nem exposta.

## Decisão
`Password` só existe internamente como hash bcrypt. `NewPassword(plain)` valida força e já devolve o VO já hasheado; `NewPasswordFromHash` reconstrói a partir do hash persistido, sem revalidar nem re-hashear.

## Consequências
A entidade `User` nunca manipula string de senha em texto puro além do momento de criação/troca — reduz superfície de vazamento acidental.

---

# ADR-007

## Título
CQRS simplificado: commands/queries fundidos com handlers, sem dispatcher genérico de commands.

## Status
Accepted

## Contexto
CQRS "de livro" separa Command (dado), Handler (execução) e Dispatcher/Mediator (roteamento). Para o tamanho deste projeto, essa separação plena adicionaria estrutura sem resolver um problema real.

## Decisão
Cada operação é uma struct `XHandler` (em `commands` ou `queries`) com um único método `Handle`. `services` orquestra os handlers atrás de uma interface (`contracts`) — sem um dispatcher central roteando por tipo.

## Consequências
Menos indireção para o tamanho atual do projeto. Se um dispatcher genérico de commands (com middlewares transversais, ex: logging automático de todo command) se tornar necessário, essa decisão pode ser revisitada.

---

# ADR-008

## Título
Unit of Work para consistência entre agregados.

## Status
Accepted

## Contexto
Algumas operações (ex: criar um `Order`, que decrementa estoque de múltiplos `Product`) envolvem mais de um Aggregate Root. Sem uma transação compartilhada, uma falha parcial deixaria dado inconsistente entre `orders` e `products`.

## Decisão
`UnitOfWork.Execute` abre uma transação `pgx.Tx`, injeta repositórios vinculados a essa transação (via a interface `DBTX`, satisfeita tanto por `*pgxpool.Pool` quanto por `pgx.Tx`) e garante commit/rollback atômico.

## Consequências
Positivas: atomicidade real entre agregados, sem transação aninhada acidental.
Negativas: `UnitOfWork` não sobrevive a uma futura decomposição em microsserviços — nesse cenário, o padrão equivalente seria Saga (compensação), não transação distribuída.

---

# ADR-009

## Título
Domain Event Dispatcher em processo, síncrono, "melhor esforço".

## Status
Accepted

## Contexto
Entidades já emitem eventos de domínio (`AddEvent`/`Events()`/`ClearEvents()`), mas nada os consumia.

## Decisão
Um `Dispatcher` percorre os eventos acumulados por uma entidade e invoca os handlers inscritos no `Registry` (por nome do evento, via `EventName()`). O dispatch só acontece **depois** que o `UnitOfWork` confirma a transação principal — nunca antes, para não reagir a algo que pode ser revertido. Falha em um handler é logada, mas não interrompe os demais handlers nem propaga para quem disparou o evento.

## Consequências
Positivas: entidades já emitem eventos no formato certo para uma futura migração a um message broker (o nome do evento e o payload já são equivalentes a um *integration event*).
Negativas: hoje é só em memória — não há garantia de entrega/retry como um message broker real ofereceria.

---

# ADR-010

## Título
Logging estruturado com `log/slog`, abstraído por uma interface própria.

## Status
Accepted

## Contexto
Logs precisam ser correlacionáveis por requisição e ter nível (`Debug`/`Info`/`Warn`/`Error`), sem acoplar o projeto a uma biblioteca de terceiros.

## Decisão
`pkg/logger.Logger` é uma interface implementada por um wrapper sobre `log/slog` (biblioteca padrão desde Go 1.21). Um logger "carimbado" com `request_id` é propagado via `context.Context` a partir do middleware HTTP.

## Consequências
Sem dependência externa; troca de implementação (ex: logger silencioso em teste) é possível sem tocar em nenhum outro pacote. Deliberadamente **sem** um método `Fatal` na interface — encerrar o processo é responsabilidade só de `main.go`/scripts de bootstrap, nunca de um handler ou command.

---

# ADR-011

## Título
Autenticação stateless via JWT.

## Status
Accepted

## Contexto
A API precisa de autenticação sem manter sessão em memória/banco.

## Decisão
Login gera um JWT assinado (HS256) contendo `user_id`, `email`, `role` e expiração. O middleware `Authenticate` valida a assinatura e injeta os claims no contexto da requisição; `RequireRole` restringe rotas por papel.

## Consequências
Positivas: sem estado de sessão no servidor.
Negativas conhecidas, aceitas por ora: sem revogação de token (logout real) nem refresh token — um token válido não pode ser invalidado antes de expirar. Mensagens de erro de login são deliberadamente genéricas (mesma resposta para "email não existe" e "senha errada"), para não permitir enumeração de usuários.

---

# ADR-012

## Título
Estratégia de testes: dublês de teste (fakes/mocks) em vez de banco real, priorizando domínio e aplicação.

## Status
Accepted

## Contexto
Testes de integração com testcontainers funcionam, mas são lentos e sujeitos a flakiness de ambiente (ex: negociação TLS entre `pgx` e um Postgres em container sem certificado).

## Decisão
`Migrator`/`Resetter` dependem de uma interface mínima (`Execer`), satisfeita tanto por `*pgxpool.Pool` (produção) quanto por um dublê em memória (testes). O mesmo padrão se aplica aos `commands`: dublês em memória das interfaces de `domain/repositories`, sem necessidade de banco. Prioridade de cobertura: Value Objects e Entidades (lógica pura, maior densidade de regra por linha) primeiro; `application/validation`, `pkg/jwt`, `application/events` em seguida; `commands`/`queries` por último, dado o maior custo de setup.

## Consequências
Testes rodam em milissegundos, sem Docker. Testes de integração real (contra Postgres de verdade) continuam existindo informalmente via `cmd/app`, mas não fazem parte da suíte automatizada.

---

# ADR-013

## Título
IDs de entidade gerados via UUID v4, dentro do construtor de cada entidade.

## Status
Accepted — substitui ADR-003

## Contexto
Ver ADR-003.

## Decisão
Cada `NewX` (`NewCustomer`, `NewProduct`, `NewOrder`, `NewUser`) chama `uuid.New()` internamente. `RebuildX`/`RestoreX` recebem o ID já existente, vindo do banco.

## Consequências
Positivas: nenhuma dependência de um serviço central de geração de ID; IDs podem ser gerados offline, sem round-trip ao banco.
Negativas: UUID v4 é maior (16 bytes) que um inteiro sequencial e não é ordenável por tempo de criação — não é um problema neste projeto, mas seria uma consideração relevante para um sistema de altíssimo volume de escrita.