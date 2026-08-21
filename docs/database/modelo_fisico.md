# Modelo físico do banco de dados

O modelo físico é implementado em PostgreSQL por meio da [migração SQL](migration.sql). O diagrama correspondente está em [UML de banco de dados](../diagramas/diagrama_banco_dados.puml).

| Tabela | Chave primária | Chaves estrangeiras | Finalidade |
| --- | --- | --- | --- |
| `users` | `id` | — | Clientes, entregadores e administradores. |
| `products` | `id` | — | Catálogo, preço e caminho da imagem do cupcake. |
| `orders` | `id` | `customer_id → users.id`; `delivery_id → users.id` | Pedido, endereço, total e status. |
| `order_items` | `id` | `order_id → orders.id`; `product_id → products.id` | Itens, quantidade e preço registrado no pedido. |
| `notifications` | `id` | `user_id → users.id`; `order_id → orders.id` | Avisos de status do pedido. |

## Regras físicas relevantes

- Tipos enumerados: `user_type`, `order_status` e `notification_type`.
- `orders.customer_id`, `orders.delivery_id`, `order_items.order_id`, `order_items.product_id`, `notifications.user_id` e `notifications.order_id` possuem integridade referencial.
- Os campos de data e exclusão lógica são mantidos pelo modelo do ORM e previstos no esquema.
- Índices, gatilhos, restrições e dados iniciais estão no arquivo de migração.
- As imagens do catálogo são arquivos locais servidos pelo front-end em `/images/cupcakes/`.
