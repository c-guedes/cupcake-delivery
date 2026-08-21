# 📊 Especificação do Banco de Dados

Sistema de Delivery de Cupcakes - Documentação Técnica do Banco de Dados

## 📋 Visão Geral

O banco de dados foi projetado para suportar um sistema completo de delivery de cupcakes, incluindo gestão de usuários, produtos, pedidos e notificações em tempo real.

### 🔧 Tecnologia
- **SGBD**: PostgreSQL 17 (Supabase)
- **ORM**: GORM (Go)
- **Extensões**: uuid-ossp
- **Hospedagem**: Supabase Cloud Database
- **Acesso**: privado pela API Go; a URL de conexão não é publicada.

## 📊 Estrutura das Tabelas

### 👥 users
Tabela central para todos os tipos de usuários do sistema.

| Campo | Tipo | Constraints | Descrição |
|-------|------|-------------|-----------|
| id | SERIAL | PRIMARY KEY | Identificador único |
| name | VARCHAR(255) | NOT NULL | Nome completo |
| email | VARCHAR(255) | UNIQUE, NOT NULL | Email único |
| password | VARCHAR(255) | NOT NULL | Senha hash (bcrypt) |
| type | user_type | NOT NULL | Tipo: customer/delivery/admin |
| vehicle | VARCHAR(255) | NULL | Veículo (obrigatório para entregadores) |
| created_at | TIMESTAMP | DEFAULT NOW() | Data de criação |
| updated_at | TIMESTAMP | DEFAULT NOW() | Data de atualização |
| deleted_at | TIMESTAMP | NULL | Soft delete (GORM) |

**Índices:**
- `idx_users_email` - Email único (excludes soft deleted)
- `idx_users_type` - Filtro por tipo de usuário
- `idx_users_deleted_at` - Otimização para soft delete

**Constraints:**
- `check_delivery_vehicle` - Entregadores devem ter veículo

### 🧁 products
Catálogo de produtos (cupcakes) disponíveis.

| Campo | Tipo | Constraints | Descrição |
|-------|------|-------------|-----------|
| id | SERIAL | PRIMARY KEY | Identificador único |
| name | VARCHAR(255) | NOT NULL | Nome do produto |
| description | TEXT | NULL | Descrição detalhada |
| price | DECIMAL(10,2) | NOT NULL, >= 0 | Preço unitário |
| image_url | VARCHAR(500) | NULL | URL da imagem |
| created_at | TIMESTAMP | DEFAULT NOW() | Data de criação |
| updated_at | TIMESTAMP | DEFAULT NOW() | Data de atualização |
| deleted_at | TIMESTAMP | NULL | Soft delete |

**Índices:**
- `idx_products_name` - Busca por nome
- `idx_products_deleted_at` - Soft delete

### 📦 orders
Pedidos realizados pelos clientes.

| Campo | Tipo | Constraints | Descrição |
|-------|------|-------------|-----------|
| id | SERIAL | PRIMARY KEY | Identificador único |
| customer_id | INTEGER | FK users(id), NOT NULL | Cliente que fez o pedido |
| delivery_id | INTEGER | FK users(id), NULL | Entregador (quando aceito) |
| status | order_status | NOT NULL, DEFAULT 'pending' | Status atual |
| total | DECIMAL(10,2) | NOT NULL, >= 0 | Valor total |
| address | TEXT | NOT NULL | Endereço de entrega |
| created_at | TIMESTAMP | DEFAULT NOW() | Data de criação |
| updated_at | TIMESTAMP | DEFAULT NOW() | Data de atualização |
| deleted_at | TIMESTAMP | NULL | Soft delete |

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

### 🛒 order_items
Itens individuais de cada pedido.

| Campo | Tipo | Constraints | Descrição |
|-------|------|-------------|-----------|
| id | SERIAL | PRIMARY KEY | Identificador único |
| order_id | INTEGER | FK orders(id), NOT NULL | Pedido pai |
| product_id | INTEGER | FK products(id), NOT NULL | Produto referenciado |
| quantity | INTEGER | NOT NULL, > 0 | Quantidade |
| price | DECIMAL(10,2) | NOT NULL, >= 0 | Preço no momento da compra |
| created_at | TIMESTAMP | DEFAULT NOW() | Data de criação |
| updated_at | TIMESTAMP | DEFAULT NOW() | Data de atualização |
| deleted_at | TIMESTAMP | NULL | Soft delete |

**Índices:**
- `idx_order_items_order_id` - Itens por pedido
- `idx_order_items_product_id` - Itens por produto
- `idx_order_items_deleted_at` - Soft delete

**Constraints:**
- `ON DELETE CASCADE` - Remove itens quando pedido é deletado
- **Nota**: Constraint implementada diretamente no PostgreSQL/Supabase

### 🔔 notifications
Sistema de notificações para usuários.

| Campo | Tipo | Constraints | Descrição |
|-------|------|-------------|-----------|
| id | SERIAL | PRIMARY KEY | Identificador único |
| user_id | INTEGER | FK users(id), NOT NULL | Usuário destinatário |
| order_id | INTEGER | FK orders(id), NULL | Pedido relacionado |
| type | notification_type | NOT NULL | Tipo da notificação |
| title | VARCHAR(255) | NOT NULL | Título |
| message | TEXT | NOT NULL | Conteúdo da mensagem |
| is_read | BOOLEAN | DEFAULT FALSE | Status de leitura |
| created_at | TIMESTAMP | DEFAULT NOW() | Data de criação |
| updated_at | TIMESTAMP | DEFAULT NOW() | Data de atualização |

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

## 🔗 Relacionamentos

### Principais Relacionamentos
1. **users → orders** (1:N)
   - Um usuário (cliente) pode ter vários pedidos
   - Um usuário (entregador) pode aceitar vários pedidos

2. **orders → order_items** (1:N)
   - Um pedido pode ter vários itens
   - CASCADE DELETE: Remove itens quando pedido é deletado

3. **products → order_items** (1:N)
   - Um produto pode estar em vários itens de pedido

4. **users → notifications** (1:N)
   - Um usuário pode receber várias notificações
   - CASCADE DELETE: Remove notificações quando usuário é deletado

5. **orders → notifications** (1:N)
   - Um pedido pode gerar várias notificações

### Restrições de Integridade
- **Email único**: Cada usuário deve ter email único
- **Tipos válidos**: Enums garantem valores consistentes
- **Valores positivos**: Preços e quantidades devem ser > 0
- **Entregador válido**: delivery_id deve referenciar usuário tipo 'delivery'
- **Veículo obrigatório**: Entregadores devem ter veículo cadastrado

## ⚡ Otimizações de Performance

### Índices Estratégicos
1. **Busca por status**: Queries frequentes filtram por status
2. **Ordenação temporal**: created_at indexado para listagens
3. **Relacionamentos**: Foreign keys indexadas
4. **Soft Delete**: deleted_at indexado para exclusão de registros deletados

### Triggers Automáticos
1. **update_updated_at**: Atualiza timestamp automaticamente
2. **check_delivery_user**: Valida se entregador é válido

### Constraints de Negócio
1. **Valores monetários**: Sempre >= 0
2. **Quantidades**: Sempre > 0
3. **Tipos de usuário**: Enum restringe valores
4. **Status de pedidos**: Enum garante fluxo correto

## 📊 Dados Iniciais (Seeds)

### Usuário Administrador
```sql
INSERT INTO users (name, email, password, type) VALUES 
('Administrador', 'admin@cupcake.com', '$2a$10$...', 'admin');
```

### Produtos Iniciais
- Cupcake de Chocolate - R$ 8,50
- Cupcake de Morango - R$ 9,00
- Cupcake Red Velvet - R$ 10,50
- Cupcake de Limão - R$ 8,00
- Cupcake de Cenoura - R$ 7,50
- Cupcake de Coco - R$ 9,50

## 🔧 Configuração e Manutenção

### Configuração Inicial
1. Executar script `migration.sql`
2. Verificar criação de todos os índices
3. Testar constraints e triggers
4. Inserir dados iniciais

### Backup e Recovery
- **Backup**: responsabilidade do serviço gerenciado Supabase, conforme plano contratado.
- **Recuperação e monitoramento**: configurados no provedor; não há promessa de rotina própria no repositório.

### Monitoramento
- **Slow queries**: Identificar queries lentas
- **Índice usage**: Verificar uso dos índices
- **Tamanho das tabelas**: Monitorar crescimento
- **Locks**: Detectar bloqueios longos

## 📈 Métricas e Análises

### Queries Frequentes
1. **Listagem de produtos**: `SELECT * FROM products WHERE deleted_at IS NULL`
2. **Pedidos por status**: `SELECT * FROM orders WHERE status = ? AND deleted_at IS NULL`
3. **Notificações não lidas**: `SELECT * FROM notifications WHERE user_id = ? AND is_read = false`
4. **Histórico de pedidos**: `SELECT * FROM orders WHERE customer_id = ? ORDER BY created_at DESC`

### Medição de desempenho
O projeto valida correção funcional e integridade das operações. Não declara metas de latência sem uma medição controlada de carga.

## 🚀 Escalabilidade

### Estratégias de Crescimento
1. **Particionamento**: Por data (orders, notifications)
2. **Réplicas de leitura**: Para dashboards e relatórios
3. **Cache**: Redis para sessões e dados frequentes
4. **Arquivamento**: Dados antigos em tabelas históricas

### Limites atuais
Os limites dependem do plano e da capacidade do provedor. Não há limpeza automática de notificações implementada.

---

## 📝 Notas de Implementação

Esta especificação reflete a implementação atual usando GORM com PostgreSQL. O design prioriza:

- **Consistência**: Constraints garantem integridade
- **Performance**: Índices otimizados para queries comuns
- **Flexibilidade**: Estrutura suporta extensões futuras
- **Manutenibilidade**: Soft delete e timestamps automáticos

## 🚀 Integração com Supabase

O sistema foi configurado para usar Supabase como banco de dados PostgreSQL gerenciado:

### 🔗 **Conexão**
- **Banco**: PostgreSQL 17 na nuvem.
- **SSL**: obrigatório na conexão de produção.
- **Segredo**: `DATABASE_URL` permanece configurada somente no ambiente do Vercel.

### 🛡️ **Segurança**
- **Row Level Security (RLS)**: Habilitado em todas as tabelas
- **Policies**: o acesso direto público permanece bloqueado; a API usa uma role limitada.
- **JWT**: emitido e validado pela própria API Go.

### 📊 **Vantagens**
- ✅ PostgreSQL gerenciado
- ✅ Console administrativo do provedor
- ✅ TLS na conexão de produção

### 🔧 **Scripts de Migração**
- `migration.sql` - Migração única usada tanto no PostgreSQL local quanto no Supabase

### 📚 **Documentação Relacionada**
- `modelo_fisico.md` - Modelo físico e relacionamentos
- `dicionario_dados.md` - Dicionário de dados
- `migration.sql` - Estrutura executável e sementes
