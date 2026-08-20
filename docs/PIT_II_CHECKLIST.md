# Checklist de entrega - PIT II

Este documento separa o que foi comprovado no repositório do que ainda precisa de evidência para a entrega.

## Situação 1 - Revisão do planejamento

| Entrega exigida | Evidência encontrada | Situação |
| --- | --- | --- |
| Escopo e ideia revisados | `DOCUMENTACAO_PLANEJAMENTO_COMPLETA.md` compara planejamento e MVP | Concluído, precisa remover promessas não implementadas |
| UML revisada | Diagramas de classe, caso de uso, sequência e banco em `docs/diagramas/` | Parcial: coexistem versões antigas e finais |
| IHC e mockups | Mockups e interface React responsiva | Parcial: validar em execução antes de registrar como evidência |
| Projeto conceitual, lógico e físico | `docs/database/` com SQL, migração e dicionário | Parcial: validar migração no PostgreSQL |
| Dicionário de dados | `docs/database/dicionario_dados.md` | Concluído documentalmente |

## Situação 2 - Solução funcional

| Entrega exigida | Evidência encontrada | Situação |
| --- | --- | --- |
| Front-end | React, TypeScript, Vite e Tailwind em `frontend/` | Pendente de build e execução |
| Back-end | Go, Gin, GORM e PostgreSQL em `backend/` | Em ajuste e pendente de integração com banco |
| Arquitetura em camadas/MVC | handlers, services, models e middleware | Parcial |
| Testes | Testes Go e Jest presentes | Back-end corrigido para voltar a executar; front pendente de validação |
| Git público | Repositórios remotos configurados | Pendente confirmar acesso público e links de produção |

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
