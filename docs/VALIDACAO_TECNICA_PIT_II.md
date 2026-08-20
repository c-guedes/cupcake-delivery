# Validação técnica - PIT II

Data: 20/08/2026

| Item exigido | Evidência | Resultado |
| --- | --- | --- |
| UML e documentação | [diagramas finais](diagramas/README.md) | Validado, renderizado e consolidado |
| Banco, migração e dicionário | [migração SQL](database/migration.sql) e [dicionário](database/dicionario_dados.md) | Executado em PostgreSQL 17 limpo |
| Back-end | `go test ./...` | Aprovado |
| Front-end | `npm run build` e `npm test -- --runInBand` | Aprovado; 14 testes |
| Fluxo integrado | API local, PostgreSQL e [capturas técnicas](pit_ii/evidencias/README.md) | Aprovado |

## Banco de dados

A migração foi executada no banco temporário `pit2_validation_20260820`: cinco tabelas, três enums, índices, triggers e seeds foram criados. Foram validados pedido com itens/endereço/total e a restrição de veículo obrigatório para entregador. O conflito entre os enums SQL e os modelos GORM foi corrigido.

## Fluxo integrado

Cliente e entregador foram cadastrados e autenticados por JWT. O cliente criou pedido com dois cupcakes; o administrador o alterou para `preparing` e `ready`; o entregador o alterou para `delivering` e `delivered`. Resultado final: pedido `delivered`, total `17.00` e cinco notificações para o cliente.

## Validação visual dos UML

Os diagramas final de classes, casos de uso, sequência e banco foram renderizados em PNG e inspecionados. O diagrama de banco foi simplificado para apresentar entidades, chaves, relacionamentos e enums; índices e triggers continuam documentados no dicionário de dados e na migração SQL.

## Pendências externas

- cinco avaliações reais;
- capturas de tela;
- link público de produção;
- vídeo narrado.
