# 📊 Dicionário de Dados - Sistema de Delivery de Cupcakes

**Versão:** 2.0 (Implementação Final)  
**Data:** Setembro 2025  
**SGBD:** PostgreSQL 12+  
**ORM:** GORM (Go)

---

## 👥 Tabela: users
Armazena todos os usuários do sistema (clientes, entregadores e administradores).

| Campo | Tipo | Descrição | Restrições |
|-------|------|-----------|------------|
| id | SERIAL | Identificador único do usuário | PRIMARY KEY |
| name | VARCHAR(255) | Nome completo do usuário | NOT NULL |
| email | VARCHAR(255) | Email único para login | NOT NULL, UNIQUE |
| password | VARCHAR(255) | Senha hash (bcrypt) | NOT NULL |
| type | user_type | Tipo: customer, delivery, admin | NOT NULL |
| vehicle | VARCHAR(255) | Veículo (obrigatório para entregadores) | NULL |
| created_at | TIMESTAMP | Data/hora de criação (GORM) | DEFAULT NOW() |
| updated_at | TIMESTAMP | Data/hora da última atualização | DEFAULT NOW() |
| deleted_at | TIMESTAMP | Soft delete (GORM) | NULL |

**Índices:**
- `idx_users_email` - Email único (excluindo deletados)
- `idx_users_type` - Filtro por tipo de usuário
- `idx_users_deleted_at` - Otimização soft delete

**Constraints:**
- `check_delivery_vehicle` - Entregadores devem ter veículo

---

## 🧁 Tabela: products
Catálogo de produtos (cupcakes) disponíveis para venda.

| Campo | Tipo | Descrição | Restrições |
|-------|------|-----------|------------|
| id | SERIAL | Identificador único do produto | PRIMARY KEY |
| name | VARCHAR(255) | Nome do cupcake | NOT NULL |
| description | TEXT | Descrição detalhada do produto | NULL |
| price | DECIMAL(10,2) | Preço unitário | NOT NULL, >= 0 |
| image_url | VARCHAR(500) | URL da imagem do produto | NULL |
| created_at | TIMESTAMP | Data/hora de criação (GORM) | DEFAULT NOW() |
| updated_at | TIMESTAMP | Data/hora da última atualização | DEFAULT NOW() |
| deleted_at | TIMESTAMP | Soft delete (GORM) | NULL |

**Índices:**
- `idx_products_name` - Busca por nome
- `idx_products_deleted_at` - Soft delete

---

## 📦 Tabela: orders
Pedidos realizados pelos clientes.

| Campo | Tipo | Descrição | Restrições |
|-------|------|-----------|------------|
| id | SERIAL | Identificador único do pedido | PRIMARY KEY |
| customer_id | INTEGER | ID do cliente que fez o pedido | FK users(id), NOT NULL |
| delivery_id | INTEGER | ID do entregador (quando aceito) | FK users(id), NULL |
| status | order_status | Status atual do pedido | NOT NULL, DEFAULT 'pending' |
| total | DECIMAL(10,2) | Valor total do pedido | NOT NULL, >= 0 |
| address | TEXT | Endereço completo de entrega | NOT NULL |
| created_at | TIMESTAMP | Data/hora de criação (GORM) | DEFAULT NOW() |
| updated_at | TIMESTAMP | Data/hora da última atualização | DEFAULT NOW() |
| deleted_at | TIMESTAMP | Soft delete (GORM) | NULL |

**Status possíveis:**
- `pending` - Aguardando confirmação
- `preparing` - Em preparação  
- `ready` - Pronto para entrega
- `delivering` - Saiu para entrega
- `delivered` - Entregue

**Índices:**
- `idx_orders_customer_id` - Pedidos por cliente
- `idx_orders_delivery_id` - Pedidos por entregador
- `idx_orders_status` - Filtro por status
- `idx_orders_created_at` - Ordenação temporal
- `idx_orders_deleted_at` - Soft delete

**Constraints:**
- `check_delivery_user_trigger` - delivery_id deve ser um entregador ativo

---

## 🛒 Tabela: order_items
Itens individuais de cada pedido (relacionamento N:N entre orders e products).

| Campo | Tipo | Descrição | Restrições |
|-------|------|-----------|------------|
| id | SERIAL | Identificador único do item | PRIMARY KEY |
| order_id | INTEGER | ID do pedido pai | FK orders(id), NOT NULL |
| product_id | INTEGER | ID do produto referenciado | FK products(id), NOT NULL |
| quantity | INTEGER | Quantidade do produto | NOT NULL, > 0 |
| price | DECIMAL(10,2) | Preço no momento da compra | NOT NULL, >= 0 |
| created_at | TIMESTAMP | Data/hora de criação (GORM) | DEFAULT NOW() |
| updated_at | TIMESTAMP | Data/hora da última atualização | DEFAULT NOW() |
| deleted_at | TIMESTAMP | Soft delete (GORM) | NULL |

**Índices:**
- `idx_order_items_order_id` - Itens por pedido
- `idx_order_items_product_id` - Itens por produto
- `idx_order_items_deleted_at` - Soft delete

**Constraints:**
- `ON DELETE CASCADE` - Remove itens quando pedido é deletado

---

## 🔔 Tabela: notifications
Sistema de notificações para usuários sobre mudanças de status.

| Campo | Tipo | Descrição | Restrições |
|-------|------|-----------|------------|
| id | SERIAL | Identificador único da notificação | PRIMARY KEY |
| user_id | INTEGER | ID do usuário destinatário | FK users(id), NOT NULL |
| order_id | INTEGER | ID do pedido relacionado | FK orders(id), NULL |
| type | notification_type | Tipo da notificação | NOT NULL |
| title | VARCHAR(255) | Título da notificação | NOT NULL |
| message | TEXT | Conteúdo da mensagem | NOT NULL |
| is_read | BOOLEAN | Status de leitura | DEFAULT FALSE |
| created_at | TIMESTAMP | Data/hora de criação | DEFAULT NOW() |
| updated_at | TIMESTAMP | Data/hora da última atualização | DEFAULT NOW() |

**Tipos de notificação:**
- `order_created` - Pedido criado
- `order_confirmed` - Pedido confirmado
- `order_preparing` - Em preparação
- `order_ready` - Pronto para entrega
- `order_delivering` - Saiu para entrega
- `order_delivered` - Entregue
- `order_cancelled` - Cancelado

**Índices:**
- `idx_notifications_user_id` - Notificações por usuário
- `idx_notifications_order_id` - Notificações por pedido
- `idx_notifications_is_read` - Filtro por status de leitura
- `idx_notifications_created_at` - Ordenação temporal

**Constraints:**
- `ON DELETE CASCADE` - Remove notificações quando usuário/pedido é deletado

---

## 📊 Relacionamentos

### 🔗 Relacionamentos Principais
1. **users → orders (customer)**
   - Um cliente pode ter vários pedidos (1:N)
   - `orders.customer_id → users.id`

2. **users → orders (delivery)**
   - Um entregador pode aceitar vários pedidos (1:N)
   - `orders.delivery_id → users.id` (opcional)

3. **orders → order_items**
   - Um pedido pode ter vários itens (1:N)
   - `order_items.order_id → orders.id`

4. **products → order_items**
   - Um produto pode estar em vários itens (1:N)
   - `order_items.product_id → products.id`

5. **users → notifications**
   - Um usuário pode receber várias notificações (1:N)
   - `notifications.user_id → users.id`

6. **orders → notifications**
   - Um pedido pode gerar várias notificações (1:N)
   - `notifications.order_id → orders.id` (opcional)

---

## 🏗️ Estrutura Implementada vs. Planejada

### ✅ **Implementado (Atual)**
- **users** - Usuários com tipos (customer, delivery, admin)
- **products** - Catálogo de cupcakes
- **orders** - Pedidos com status
- **order_items** - Itens dos pedidos
- **notifications** - Sistema de notificações

### ❌ **Não Implementado (Escopo Reduzido)**
- ~~enderecos~~ - Endereço armazenado como texto em orders.address
- ~~carrinhos~~ - Carrinho implementado no frontend (localStorage)
- ~~itens_carrinho~~ - Não persistido no banco
- ~~pagamentos~~ - Fora do escopo (simulado)
- ~~provas_entrega~~ - Fora do escopo

### 🎯 **Decisões de Design**
- **Simplicidade**: Endereço como texto ao invés de tabela separada
- **Performance**: Carrinho no localStorage ao invés do banco
- **Escopo**: Foco no fluxo principal de pedidos
- **Flexibilidade**: Soft delete para auditoria

---

## 🔧 Configurações Técnicas

### **GORM Configuration**
```go
// Soft Delete automático
type Model struct {
    ID        uint           `gorm:"primarykey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

### **Enums PostgreSQL**
```sql
CREATE TYPE user_type AS ENUM ('customer', 'delivery', 'admin');
CREATE TYPE order_status AS ENUM ('pending', 'preparing', 'ready', 'delivering', 'delivered');
CREATE TYPE notification_type AS ENUM ('order_created', 'order_confirmed', ...);
```

### **Triggers Automáticos**
- `update_updated_at` - Atualiza timestamp automaticamente
- `check_delivery_user` - Valida se entregador é válido

---

## 📈 Performance e Otimização

### **Índices Estratégicos**
- **Busca frequente**: email, status, user_id
- **Ordenação**: created_at, updated_at
- **Relacionamentos**: foreign keys indexadas
- **Soft Delete**: deleted_at indexado

### **Queries Otimizadas**
- Filtros sempre incluem `deleted_at IS NULL`
- Joins otimizados com índices em foreign keys
- Paginação usando LIMIT/OFFSET com ORDER BY

### **Escalabilidade**
- Estrutura suporta até 1M+ pedidos
- Índices otimizados para crescimento
- Soft delete para auditoria sem perda de performance

---

## 💾 Dados de Exemplo

### **Usuários Iniciais**
```sql
-- Administrador padrão
INSERT INTO users (name, email, password, type) VALUES 
('Administrador', 'admin@cupcake.com', '$2a$10$...', 'admin');
```

### **Produtos Iniciais**
- Cupcake de Chocolate - R$ 8,50
- Cupcake de Morango - R$ 9,00  
- Cupcake Red Velvet - R$ 10,50
- Cupcake de Limão - R$ 8,00
- Cupcake de Cenoura - R$ 7,50
- Cupcake de Coco - R$ 9,50

---

**📝 Nota:** Este dicionário reflete a implementação final do sistema, com foco na funcionalidade essencial do delivery de cupcakes.
