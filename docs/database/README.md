# Banco de dados - PIT II

Esta pasta é a referência da persistência da aplicação Cupcake Delivery. A implementação usa PostgreSQL hospedado no Supabase e é acessada exclusivamente pela API Go; o cliente React não recebe credenciais do banco.

## Artefatos finais

- [Modelo físico](modelo_fisico.md): tabelas, chaves, cardinalidades, tipos e regras.
- [Modelo conceitual/lógico](modelo_conceitual.sql): estrutura relacional normalizada.
- [Migração SQL](migration.sql): enums, tabelas, chaves, índices, triggers, seeds e RLS.
- [Dicionário de dados](dicionario_dados.md): domínio e significado de cada campo.
- [Especificação técnica](DATABASE_SPEC.md): restrições e decisões de persistência.

## Execução da migração

Use um banco PostgreSQL vazio e execute `migration.sql` uma única vez com uma conta administradora. A aplicação em produção opera com `AUTO_MIGRATE=false`, portanto não cria nem altera a estrutura automaticamente em cada inicialização.

```bash
psql "$DATABASE_URL" -f migration.sql
```

O script cria cinco tabelas (`users`, `products`, `orders`, `order_items` e `notifications`), três enums, índices para as consultas usuais, triggers de integridade e seis produtos iniciais.

## Produção e segurança

- Banco: PostgreSQL gerenciado pelo Supabase.
- API: Go/Gin publicada no Vercel, conectada por uma credencial de runtime com permissões limitadas.
- Acesso público direto via PostgREST: bloqueado por RLS; somente a API aplica as regras de negócio.
- Aplicação: [Cupcake Delivery em produção](https://cupcake-delivery-pit-ii.vercel.app).

As variáveis de conexão devem permanecer no ambiente de implantação; nunca devem ser adicionadas ao repositório.
