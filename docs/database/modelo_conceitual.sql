-- Modelo Conceitual do Banco de Dados
-- Sistema de Delivery de Cupcakes (Implementação Final)

/*
DECISÕES DE DESIGN DA IMPLEMENTAÇÃO FINAL:

1. SIMPLICIDADE: Foco no essencial do delivery de cupcakes
2. GORM: Uso do ORM Go com convenções específicas
3. POSTGRESQL: SGBD com recursos avançados
4. SOFT DELETE: Auditoria sem perda de dados

ENTIDADES IMPLEMENTADAS:

1. USERS (Usuários do sistema)
   - Tabela unificada com campo "type" para diferentes perfis
   - Tipos: customer (cliente), delivery (entregador), admin (administrador)

2. PRODUCTS (Cupcakes disponíveis)
   - Catálogo simples com nome, descrição, preço e imagem

3. ORDERS (Pedidos realizados)
   - Cabeçalho do pedido com cliente, entregador (opcional) e endereço como texto

4. ORDER_ITEMS (Itens do pedido)
   - Relacionamento N:N entre pedidos e produtos

5. NOTIFICATIONS (Notificações do sistema)
   - Sistema de comunicação em tempo real sobre mudanças de status

ENTIDADES NÃO IMPLEMENTADAS (por simplicidade):
- ADDRESSES: Endereço armazenado como texto no pedido
- CART: Implementado no frontend (localStorage)
- PAYMENTS: Fora do escopo (simulado)
- DELIVERY_PROOF: Fora do escopo

RELACIONAMENTOS PRINCIPAIS:

1. Um USER (customer) pode ter vários ORDERS
2. Um USER (delivery) pode aceitar vários ORDERS
3. Um ORDER tem vários ORDER_ITEMS
4. Um PRODUCT pode estar em vários ORDER_ITEMS
5. Um USER pode receber várias NOTIFICATIONS
6. Um ORDER pode gerar várias NOTIFICATIONS
*/

-- Modelo conceitual implementado (GORM + PostgreSQL)

CREATE TABLE users (
    id SERIAL PRIMARY KEY,                    -- Auto increment (GORM)
    name VARCHAR(255) NOT NULL,               -- Nome completo
    email VARCHAR(255) NOT NULL UNIQUE,       -- Email único para login
    password VARCHAR(255) NOT NULL,           -- Hash bcrypt
    type user_type NOT NULL,                  -- customer, delivery, admin
    vehicle VARCHAR(255),                     -- Veículo (só entregadores)
    created_at TIMESTAMP NOT NULL,            -- GORM timestamp
    updated_at TIMESTAMP NOT NULL,            -- GORM timestamp
    deleted_at TIMESTAMP                      -- GORM soft delete
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,                    -- Auto increment (GORM)
    name VARCHAR(255) NOT NULL,               -- Nome do cupcake
    description TEXT,                         -- Descrição opcional
    price DECIMAL(10,2) NOT NULL,             -- Preço unitário
    image_url VARCHAR(500),                   -- URL da imagem
    created_at TIMESTAMP NOT NULL,            -- GORM timestamp
    updated_at TIMESTAMP NOT NULL,            -- GORM timestamp
    deleted_at TIMESTAMP                      -- GORM soft delete
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,                    -- Auto increment (GORM)
    customer_id INTEGER NOT NULL,             -- FK para users (cliente)
    delivery_id INTEGER,                      -- FK para users (entregador, opcional)
    status order_status NOT NULL,             -- pending, preparing, ready, delivering, delivered
    total DECIMAL(10,2) NOT NULL,             -- Valor total calculado
    address TEXT NOT NULL,                    -- Endereço como texto simples
    created_at TIMESTAMP NOT NULL,            -- GORM timestamp
    updated_at TIMESTAMP NOT NULL,            -- GORM timestamp
    deleted_at TIMESTAMP                      -- GORM soft delete
);

CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,                    -- Auto increment (GORM)
    order_id INTEGER NOT NULL,                -- FK para orders
    product_id INTEGER NOT NULL,              -- FK para products
    quantity INTEGER NOT NULL,                -- Quantidade
    price DECIMAL(10,2) NOT NULL,             -- Preço no momento da compra
    created_at TIMESTAMP NOT NULL,            -- GORM timestamp
    updated_at TIMESTAMP NOT NULL,            -- GORM timestamp
    deleted_at TIMESTAMP                      -- GORM soft delete
);

CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,                    -- Auto increment (GORM)
    user_id INTEGER NOT NULL,                 -- FK para users (destinatário)
    order_id INTEGER,                         -- FK para orders (opcional)
    type notification_type NOT NULL,          -- Tipo da notificação
    title VARCHAR(255) NOT NULL,              -- Título
    message TEXT NOT NULL,                    -- Conteúdo
    is_read BOOLEAN DEFAULT FALSE,            -- Status de leitura
    created_at TIMESTAMP NOT NULL,            -- Timestamp de criação
    updated_at TIMESTAMP NOT NULL             -- Timestamp de atualização
);

-- Relacionamentos (Foreign Keys)

ALTER TABLE orders
    ADD CONSTRAINT fk_orders_customer
    FOREIGN KEY (customer_id) REFERENCES users(id),
    ADD CONSTRAINT fk_orders_delivery
    FOREIGN KEY (delivery_id) REFERENCES users(id);

ALTER TABLE order_items
    ADD CONSTRAINT fk_order_items_order
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    ADD CONSTRAINT fk_order_items_product
    FOREIGN KEY (product_id) REFERENCES products(id);

ALTER TABLE notifications
    ADD CONSTRAINT fk_notifications_user
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    ADD CONSTRAINT fk_notifications_order
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE;

-- Constraints de negócio

-- Entregadores devem ter veículo
ALTER TABLE users
    ADD CONSTRAINT check_delivery_vehicle
    CHECK (type != 'delivery' OR vehicle IS NOT NULL);

-- Valores monetários devem ser positivos
ALTER TABLE products
    ADD CONSTRAINT check_price_positive
    CHECK (price >= 0);

ALTER TABLE orders
    ADD CONSTRAINT check_total_positive
    CHECK (total >= 0);

ALTER TABLE order_items
    ADD CONSTRAINT check_quantity_positive CHECK (quantity > 0),
    ADD CONSTRAINT check_price_positive CHECK (price >= 0);

-- Índices para otimização de consultas frequentes

-- Busca de usuários
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_type ON users(type) WHERE deleted_at IS NULL;

-- Busca de produtos
CREATE INDEX idx_products_name ON products(name) WHERE deleted_at IS NULL;

-- Busca de pedidos (queries mais frequentes)
CREATE INDEX idx_orders_customer_id ON orders(customer_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_orders_delivery_id ON orders(delivery_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_orders_status ON orders(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_orders_created_at ON orders(created_at) WHERE deleted_at IS NULL;

-- Busca de itens do pedido
CREATE INDEX idx_order_items_order_id ON order_items(order_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_order_items_product_id ON order_items(product_id) WHERE deleted_at IS NULL;

-- Busca de notificações
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);

-- Soft delete
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE INDEX idx_products_deleted_at ON products(deleted_at);
CREATE INDEX idx_orders_deleted_at ON orders(deleted_at);
CREATE INDEX idx_order_items_deleted_at ON order_items(deleted_at);

-- Comentários das tabelas para documentação

COMMENT ON TABLE users IS 'Usuários do sistema (clientes, entregadores e administradores) - Implementação GORM';
COMMENT ON TABLE products IS 'Catálogo de cupcakes disponíveis para venda - Soft delete habilitado';
COMMENT ON TABLE orders IS 'Pedidos realizados pelos clientes - Endereço como texto simples';
COMMENT ON TABLE order_items IS 'Itens incluídos em cada pedido - Relacionamento N:N com products';
COMMENT ON TABLE notifications IS 'Sistema de notificações em tempo real sobre mudanças de status';

-- Comentários sobre decisões de design

COMMENT ON COLUMN users.type IS 'Tipo de usuário: customer (cliente), delivery (entregador), admin (administrador)';
COMMENT ON COLUMN users.vehicle IS 'Informações do veículo - obrigatório apenas para entregadores';
COMMENT ON COLUMN orders.address IS 'Endereço completo como texto - simplificação vs tabela separada';
COMMENT ON COLUMN order_items.price IS 'Preço do produto no momento da compra - snapshot para auditoria';
COMMENT ON COLUMN notifications.is_read IS 'Status de leitura da notificação pelo usuário';

/*
OBSERVAÇÕES SOBRE A IMPLEMENTAÇÃO:

1. GORM CONVENTIONS:
   - IDs como SERIAL (auto increment)
   - Campos created_at, updated_at, deleted_at automáticos
   - Soft delete habilitado

2. POSTGRESQL FEATURES:
   - Enums para tipos específicos
   - Constraints avançadas
   - Índices parciais (WHERE deleted_at IS NULL)

3. PERFORMANCE:
   - Índices em foreign keys
   - Índices em campos de busca frequente
   - Soft delete indexado

4. BUSINESS RULES:
   - Entregadores devem ter veículo
   - Valores monetários >= 0
   - delivery_id deve ser tipo delivery

5. SIMPLICIDADE:
   - Endereço como texto ao invés de tabela normalizada
   - Carrinho no frontend (localStorage)
   - Pagamento simulado (não persistido)
   - Foco no fluxo essencial de delivery
*/
