# Checklist de entrega - PIT II

Este documento separa o que foi comprovado no repositório do que ainda precisa de evidência para a entrega.

## Situação 1 - Revisão do planejamento

| Entrega exigida | Evidência encontrada | Situação |
| --- | --- | --- |
| Escopo e ideia revisados | Documentação consolidada em `docs/` e evidências do MVP | Concluído |
| UML revisada | Versões finais identificadas em `docs/diagramas/README.md` | Concluído e renderizado visualmente |
| IHC e mockups | `docs/mockups/ihc_validada.html` e interface publicada | Concluído; somente fluxos implementados são apresentados |
| Projeto conceitual, lógico e físico | `docs/database/` com SQL, modelo físico e dicionário | Concluído; migração validada em PostgreSQL |
| Dicionário de dados | `docs/database/dicionario_dados.md` | Concluído documentalmente |

## Situação 2 - Solução funcional

| Entrega exigida | Evidência encontrada | Situação |
| --- | --- | --- |
| Front-end | React, TypeScript, Vite e Tailwind em `frontend/` | Build, testes e produção aprovados |
| Back-end | Go, Gin, GORM e PostgreSQL em `backend/` | Testes, integração e produção aprovados |
| Arquitetura em camadas/MVC | handlers, services, models e middleware | Parcial |
| Testes | Testes Go e Jest presentes | `go test ./...`, `npm test` e `npm run build` aprovados |
| Git público | [GitHub](https://github.com/c-guedes/cupcake-delivery) e [produção](https://cupcake-delivery-pit-ii.vercel.app) | Repositório público e aplicação ativa; falta mesclar a branch final na `main` |

## Situação 3 - Validação e qualidade

| Entrega exigida | Evidência encontrada | Situação |
| --- | --- | --- |
| Cinco testes com colegas | Apenas o modelo de coleta | Pendente: não fabricar respostas ou evidências |
| Laudo de qualidade com snapshots | Não encontrado | Pendente após testes reais |
| Correções decorrentes dos testes | Não encontrado | Pendente após feedback real |
| Vídeo narrado de até cinco minutos | Não encontrado | Pendente após validação final |

## Ajustes iniciados

- Notificações passam a ser incluídas na migração automática do GORM.
- Pedido exige endereço, grava itens e total, e gera notificação inicial.
- Atualização de status aceita dados do front-end e restringe o administrador ao fluxo `pending -> preparing -> ready`.
- Removido arquivo de teste vazio que impedia a execução da suíte Go.

## Próxima validação

1. Executar testes e build do front-end/back-end.
2. Subir PostgreSQL e testar o fluxo cliente, administrador e entregador.
3. Consolidar os UML para exibirem apenas funcionalidades comprovadas.
4. Aplicar o formulário a cinco pessoas e registrar os resultados reais.
