# Validação técnica - PIT II

Data: 20/08/2026

| Item exigido | Evidência | Resultado |
| --- | --- | --- |
| UML e documentação | [diagramas finais](diagramas/README.md) | Validado, renderizado e consolidado |
| Banco, migração e dicionário | [modelo físico](database/modelo_fisico.md), [migração SQL](database/migration.sql) e [dicionário](database/dicionario_dados.md) | Executado novamente em PostgreSQL 17 limpo |
| Back-end | `go test ./...` | Aprovado |
| Front-end | `npm run build` e `npm test -- --runInBand` | Aprovado; 14 testes |
| Fluxo integrado | API, PostgreSQL e [capturas técnicas](pit_ii/evidencias/README.md) | Aprovado localmente e em produção |
| Produção pública | [Vercel](https://cupcake-delivery-pit-ii.vercel.app), [health](https://cupcake-delivery-pit-ii.vercel.app/api/health) e [products](https://cupcake-delivery-pit-ii.vercel.app/api/products) | Retestado: HTTP 200 nos quatro pontos (raiz, login, health e catálogo) |

## Banco de dados

A migração foi executada novamente no banco temporário `pit2_audit_20260820`: cinco tabelas, três enums, índices, triggers e seeds foram criados. Os seis produtos receberam caminhos de imagem locais em `/images/cupcakes/`. Foram validados pedido com itens/endereço/total e a restrição de veículo obrigatório para entregador. O conflito entre os enums SQL e os modelos GORM foi corrigido.

## Fluxo integrado

Cliente e entregador foram cadastrados e autenticados por JWT. O cliente criou pedido com dois cupcakes; o administrador o alterou para `preparing` e `ready`; o entregador o alterou para `delivering` e `delivered`. Resultado final: pedido `delivered`, total `17.00` e cinco notificações para o cliente.

## Validação visual dos UML

Os diagramas final de classes, casos de uso, sequência e banco foram renderizados em PNG e inspecionados. O diagrama de banco foi simplificado para apresentar entidades, chaves, relacionamentos e enums; índices e triggers continuam documentados no dicionário de dados e na migração SQL.

## Validação visual da aplicação

O catálogo foi executado novamente após a troca dos placeholders. As seis imagens de cupcake locais foram exibidas corretamente, assim como o fluxo de pedido criado e o pedido pronto no painel administrativo. As capturas estão em [evidências técnicas](pit_ii/evidencias/README.md).

## Produção

A aplicação é entregue como serviço público no Vercel. O front-end React é servido na raiz e a API Go responde sob o prefixo `/api`; o PostgreSQL é hospedado no Supabase, com credencial limitada à aplicação. A verificação de produção realizada em 20/08/2026 confirmou `GET /`, `GET /login`, `GET /api/health` e `GET /api/products` com HTTP 200.

## Pendências externas

- cinco avaliações reais;
- vídeo narrado.
