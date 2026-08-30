# Backlog

Itens abaixo representam evoluções de **produto/domínio** ainda não implementadas — não fazem parte da etapa atual (monolito), e não têm data definida.

## Catálogo
- [ ] Categorias
- [ ] Marcas
- [ ] Variantes
- [ ] SKU
- [ ] Imagens

## Clientes
- [ ] Múltiplos telefones
- [ ] Múltiplos endereços
- [ ] Integração com ViaCEP + Value Object `CEP` (avaliado e adiado nesta etapa; ver `decisions.md` ADR-012 sobre integrações externas ainda não implementadas)

## Estoque
- [ ] Reserva de estoque com expiração (hoje o estoque é decrementado na criação do pedido, sem essa granularidade)
- [ ] Movimentações auditáveis (histórico de entradas/saídas)
- [ ] Múltiplos depósitos

## Segurança
- [x] Usuário
- [x] JWT
- [x] Roles
- [ ] OAuth2 / login social
- [ ] Refresh token
- [ ] Revogação de token (logout real)

## Pedidos
- [ ] Cupons
- [ ] Frete
- [ ] Integração real de pagamento (hoje `Pay()` é uma transição de status, sem gateway)
- [ ] Histórico de mudanças de status (hoje só o status atual é persistido)

## Infraestrutura / Plataforma (fora do escopo de produto, mas planejado)
- [ ] Dockerizar a aplicação inteira (hoje só o Postgres roda em container)
- [ ] Front-end simples para testes manuais

## Exploração futura (projeto/etapa separada, não compromissada)
- [ ] Decomposição do monolito em microsserviços por agregado (`Customer`, `Product`, `Order`, `User`)
- [ ] Padrão Saga para consistência entre serviços (substituindo o Unit of Work local)