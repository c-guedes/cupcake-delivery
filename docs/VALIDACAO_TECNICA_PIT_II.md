# Validação técnica - PIT II

Data: 20/08/2026

| Item exigido | Evidência | Resultado |
| --- | --- | --- |
| UML e documentação | [diagramas finais](diagramas/README.md) | Validado, renderizado e consolidado |
| Banco, migração e dicionário | [modelo físico](database/modelo_fisico.md), [migração SQL](database/migration.sql) e [dicionário](database/dicionario_dados.md) | Executado novamente em PostgreSQL 17 limpo |
| Back-end | `go test ./...` | Aprovado |
| Front-end | `npm run build` e `npm test -- --runInBand` | Aprovado; 14 testes |
| Fluxo integrado | API local, PostgreSQL e [capturas técnicas](pit_ii/evidencias/README.md) | Aprovado |

## Banco de dados

A migração foi executada novamente no banco temporário `pit2_audit_20260820`: cinco tabelas, três enums, índices, triggers e seeds foram criados. Os seis produtos receberam caminhos de imagem locais em `/images/cupcakes/`. Foram validados pedido com itens/endereço/total e a restrição de veículo obrigatório para entregador. O conflito entre os enums SQL e os modelos GORM foi corrigido.

## Fluxo integrado

Cliente e entregador foram cadastrados e autenticados por JWT. O cliente criou pedido com dois cupcakes; o administrador o alterou para `preparing` e `ready`; o entregador o alterou para `delivering` e `delivered`. Resultado final: pedido `delivered`, total `17.00` e cinco notificações para o cliente.

## Validação visual dos UML

Os diagramas final de classes, casos de uso, sequência e banco foram renderizados em PNG e inspecionados. O diagrama de banco foi simplificado para apresentar entidades, chaves, relacionamentos e enums; índices e triggers continuam documentados no dicionário de dados e na migração SQL.

## Validação visual da aplicação

O catálogo foi executado novamente após a troca dos placeholders. As seis imagens de cupcake locais foram exibidas corretamente, assim como o fluxo de pedido criado e o pedido pronto no painel administrativo. As capturas estão em [evidências técnicas](pit_ii/evidencias/README.md).

## Pendências externas

- cinco avaliações reais;
- capturas de tela;
- link público de produção;
- vídeo narrado.
