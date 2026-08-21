-- Script de Criação do Banco de Dados
-- Sistema de Delivery de Cupcakes (Implementação Final)
-- PostgreSQL 12+
-- GORM Compatible

-- Criar banco de dados
CREATE DATABASE cupcake_delivery
    WITH 
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

\connect cupcake_delivery

-- Extensões necessárias
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Criar enums para tipos específicos (conforme implementação)
CREATE TYPE user_type AS ENUM ('customer', 'delivery', 'admin');
CREATE TYPE order_status AS ENUM ('pending', 'preparing', 'ready', 'delivering', 'delivered');
CREATE TYPE notification_type AS ENUM (
    'order_created',
    'order_confirmed',
    'order_preparing',
    'order_ready',
    'order_delivering',
    'order_delivered',
    'order_cancelled'
);

-- Tabela de usuários (conforme models.go)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    type user_type NOT NULL,
    vehicle VARCHAR(255), -- Apenas para entregadores
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft delete (GORM)
);

-- Tabela de produtos (conforme models.go)
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    image_url VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft delete (GORM)
);

-- Tabela de pedidos (conforme models.go)
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER NOT NULL REFERENCES users(id),
    delivery_id INTEGER REFERENCES users(id), -- Opcional
    status order_status NOT NULL DEFAULT 'pending',
    total DECIMAL(10,2) NOT NULL CHECK (total >= 0),
    address TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft delete (GORM)
);

-- Tabela de itens do pedido (conforme models.go)
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft delete (GORM)
);

-- Tabela de notificações (conforme notification.go)
CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
    type notification_type NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Índices para otimização de performance
-- Usuários
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_type ON users(type) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- Produtos  
CREATE INDEX idx_products_name ON products(name) WHERE deleted_at IS NULL;
CREATE INDEX idx_products_deleted_at ON products(deleted_at);

-- Pedidos
CREATE INDEX idx_orders_customer_id ON orders(customer_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_orders_delivery_id ON orders(delivery_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_orders_status ON orders(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_orders_created_at ON orders(created_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_orders_deleted_at ON orders(deleted_at);

-- Itens do pedido
CREATE INDEX idx_order_items_order_id ON order_items(order_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_order_items_product_id ON order_items(product_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_order_items_deleted_at ON order_items(deleted_at);

-- Notificações
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_order_id ON notifications(order_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);

-- Função para atualizar updated_at automaticamente
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers para atualização automática de updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_orders_updated_at BEFORE UPDATE ON orders
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_order_items_updated_at BEFORE UPDATE ON order_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_notifications_updated_at BEFORE UPDATE ON notifications
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Constraint para verificar se entregador tem veículo
ALTER TABLE users ADD CONSTRAINT check_delivery_vehicle 
    CHECK (type != 'delivery' OR vehicle IS NOT NULL);

-- Função para verificar se delivery_id é realmente um entregador
CREATE OR REPLACE FUNCTION check_delivery_user()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.delivery_id IS NOT NULL THEN
        IF NOT EXISTS (
            SELECT 1 FROM users 
            WHERE id = NEW.delivery_id 
            AND type = 'delivery' 
            AND deleted_at IS NULL
        ) THEN
            RAISE EXCEPTION 'delivery_id deve referenciar um usuário do tipo delivery ativo';
        END IF;
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER check_delivery_user_trigger 
    BEFORE INSERT OR UPDATE ON orders
    FOR EACH ROW EXECUTE FUNCTION check_delivery_user();

-- Dados iniciais (seeds)
-- Usuário administrador padrão
INSERT INTO users (name, email, password, type) VALUES 
('Administrador', 'admin@cupcake.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin')
ON CONFLICT (email) DO NOTHING;

-- Produtos iniciais
INSERT INTO products (name, description, price, image_url) VALUES 
('Cupcake de Chocolate', 'Delicioso cupcake de chocolate com cobertura de brigadeiro', 8.50, '/images/cupcakes/chocolate.png'),
('Cupcake de Morango', 'Cupcake de baunilha com cobertura de morango e pedaços da fruta', 9.00, '/images/cupcakes/morango.png'),
('Cupcake Red Velvet', 'Clássico red velvet com cream cheese', 10.50, '/images/cupcakes/red-velvet.png'),
('Cupcake de Limão', 'Cupcake cítrico com cobertura de limão siciliano', 8.00, '/images/cupcakes/limao.png'),
('Cupcake de Cenoura', 'Cupcake de cenoura com cobertura de chocolate', 7.50, '/images/cupcakes/cenoura.png'),
('Cupcake de Coco', 'Cupcake de coco com cobertura cremosa e coco ralado', 9.50, '/images/cupcakes/coco.png')
ON CONFLICT DO NOTHING;

-- Comentários nas tabelas
COMMENT ON DATABASE cupcake_delivery IS 'Banco de dados do sistema de delivery de cupcakes (implementação final)';
COMMENT ON TABLE users IS 'Usuários do sistema (clientes, entregadores e administradores)';
COMMENT ON TABLE products IS 'Catálogo de cupcakes disponíveis para venda';
COMMENT ON TABLE orders IS 'Pedidos realizados pelos clientes';
COMMENT ON TABLE order_items IS 'Itens individuais de cada pedido';
COMMENT ON TABLE notifications IS 'Sistema de notificações para os usuários';

-- Comentários nas colunas principais
COMMENT ON COLUMN users.type IS 'Tipo de usuário: customer, delivery ou admin';
COMMENT ON COLUMN users.vehicle IS 'Informações do veículo (obrigatório para entregadores)';
COMMENT ON COLUMN orders.status IS 'Status atual do pedido';
COMMENT ON COLUMN order_items.price IS 'Preço do produto no momento da compra';
COMMENT ON COLUMN notifications.is_read IS 'Indica se a notificação foi lida pelo usuário';

-- Verificação final
SELECT 'Banco de dados criado com sucesso!' as status;
