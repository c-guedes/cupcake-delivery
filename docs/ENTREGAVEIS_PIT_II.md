# Auditoria de entregáveis — PIT II

Revisado em 20/08/2026 contra o enunciado da atividade e os arquivos do repositório.

| Exigência | Evidência no repositório | Situação |
| --- | --- | --- |
| Escopo, UML, IHC/mockups, dados e dicionário revisados | [UML](diagramas/README.md), [IHC validada](mockups/README.md), [modelo físico](database/modelo_fisico.md), [migração](database/migration.sql), [dicionário](database/dicionario_dados.md) | Concluído documentalmente |
| Código front-end e back-end testado | [relatório técnico](VALIDACAO_TECNICA_PIT_II.md) | Concluído tecnicamente |
| Evidência visual de execução | [capturas técnicas](pit_ii/evidencias/README.md) | Parcial: login, cadastro, catálogo com imagens, pedido criado e painel administrativo |
| Repositório GitHub aberto | [repositório](https://github.com/c-guedes/cupcake-delivery) | Confirmado público em 20/08/2026. A documentação final está na branch `pit-ii-documentacao-validacao`; antes do envio, mesclar na `main` para que o link padrão do repositório exiba a versão final. |
| Link público da aplicação funcionando | [produção no Vercel](https://cupcake-delivery-pit-ii.vercel.app) e endpoints [health](https://cupcake-delivery-pit-ii.vercel.app/api/health) / [products](https://cupcake-delivery-pit-ii.vercel.app/api/products) | Concluído e retestado em 20/08/2026: respostas HTTP 200; API conectada ao PostgreSQL hospedado no Supabase. |
| Cinco fichas/opiniões reais em PDF e evidências | [formulário](pit_ii/formulario_cinco_testes_usuarios.pdf) e [laudo modelo](pit_ii/laudo_qualidade_template.pdf) | Pendente avaliação real; modelos não são evidência de teste realizado |
| Vídeo narrado de aproximadamente cinco minutos | Não encontrado | Pendente |

## Roteiro objetivo para finalizar

1. Mesclar a branch `pit-ii-documentacao-validacao` na `main`, mantendo a visibilidade pública já confirmada.
2. Aplicar o formulário a cinco participantes reais, gerar os cinco PDFs preenchidos e completar o laudo de qualidade com as evidências.
3. Gravar o vídeo narrado de até cinco minutos mostrando: login, catálogo, criação do pedido, painel administrativo, entrega, testes e documentação.

## Observação de integridade acadêmica

As evidências técnicas comprovam execução do software. Avaliações de usuários só devem ser registradas após participação real; não devem ser simuladas ou atribuídas a pessoas inexistentes.
