-- Script de Criação do Banco de Dados Físico
-- Sistema de Delivery de Cupcakes
-- PostgreSQL 14+

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

-- Habilitar extensão UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Criar tipos ENUM para maior consistência
CREATE TYPE usuario_tipo AS ENUM ('CLIENTE', 'ADMIN', 'ENTREGADOR');
CREATE TYPE produto_status AS ENUM ('ATIVO', 'INATIVO', 'ESGOTADO');
CREATE TYPE produto_categoria AS ENUM ('TRADICIONAL', 'ESPECIAL');
CREATE TYPE pedido_status AS ENUM (
    'PENDENTE',
    'CONFIRMADO',
    'PREPARANDO',
    'ENVIADO',
    'ENTREGUE',
    'CANCELADO'
);
CREATE TYPE pagamento_status AS ENUM ('PENDENTE', 'APROVADO', 'REJEITADO');
CREATE TYPE pagamento_metodo AS ENUM ('CARTAO', 'PIX');
CREATE TYPE notificacao_tipo AS ENUM (
    'PEDIDO_CONFIRMADO',
    'PEDIDO_PREPARANDO',
    'SAIU_ENTREGA',
    'ENTREGUE',
    'CANCELADO'
);

-- Criação das tabelas
CREATE TABLE usuarios (
    id UUID DEFAULT uuid_generate_v4(),
    nome VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    senha VARCHAR(255) NOT NULL,
    tipo usuario_tipo NOT NULL,
    telefone VARCHAR(20),
    placa_veiculo VARCHAR(10),
    data_criacao TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    ultima_atualizacao TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_usuarios PRIMARY KEY (id),
    CONSTRAINT chk_entregador_placa CHECK (
        (tipo = 'ENTREGADOR' AND placa_veiculo IS NOT NULL) OR
        (tipo != 'ENTREGADOR' AND placa_veiculo IS NULL)
    )
);

CREATE TABLE enderecos (
    id UUID DEFAULT uuid_generate_v4(),
    usuario_id UUID NOT NULL,
    cep VARCHAR(9) NOT NULL,
    rua VARCHAR(100) NOT NULL,
    numero VARCHAR(10) NOT NULL,
    complemento VARCHAR(100),
    bairro VARCHAR(100) NOT NULL,
    cidade VARCHAR(100) NOT NULL,
    estado CHAR(2) NOT NULL,
    padrao BOOLEAN DEFAULT false,
    data_criacao TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_enderecos PRIMARY KEY (id),
    CONSTRAINT fk_enderecos_usuario FOREIGN KEY (usuario_id)
        REFERENCES usuarios (id) ON DELETE CASCADE
);

CREATE TABLE produtos (
    id UUID DEFAULT uuid_generate_v4(),
    nome VARCHAR(100) NOT NULL,
    descricao TEXT,
    preco DECIMAL(10,2) NOT NULL,
    imagem_url VARCHAR(255),
    categoria produto_categoria DEFAULT 'TRADICIONAL',
    status produto_status NOT NULL DEFAULT 'ATIVO',
    data_criacao TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    ultima_atualizacao TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_produtos PRIMARY KEY (id),
    CONSTRAINT chk_preco_positivo CHECK (preco > 0)
);

CREATE TABLE carrinhos (
    id UUID DEFAULT uuid_generate_v4(),
    usuario_id UUID NOT NULL,
    data_criacao TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    ultima_atualizacao TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_carrinhos PRIMARY KEY (id),
    CONSTRAINT fk_carrinhos_usuario FOREIGN KEY (usuario_id)
        REFERENCES usuarios (id) ON DELETE CASCADE
);

CREATE TABLE itens_carrinho (
    id UUID DEFAULT uuid_generate_v4(),
    carrinho_id UUID NOT NULL,
    produto_id UUID NOT NULL,
    quantidade INTEGER NOT NULL,
    preco_unitario DECIMAL(10,2) NOT NULL,
    subtotal DECIMAL(10,2) NOT NULL,
    data_criacao TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_itens_carrinho PRIMARY KEY (id),
    CONSTRAINT fk_itens_carrinho_carrinho FOREIGN KEY (carrinho_id)
        REFERENCES carrinhos (id) ON DELETE CASCADE,
    CONSTRAINT fk_itens_carrinho_produto FOREIGN KEY (produto_id)
        REFERENCES produtos (id) ON DELETE RESTRICT,
    CONSTRAINT chk_quantidade_positiva CHECK (quantidade > 0),
    CONSTRAINT chk_preco_positivo CHECK (preco_unitario > 0),
    CONSTRAINT chk_subtotal CHECK (subtotal = quantidade * preco_unitario)
);

CREATE TABLE pedidos (
    id UUID DEFAULT uuid_generate_v4(),
    usuario_id UUID NOT NULL,
    entregador_id UUID,
    endereco_id UUID NOT NULL,
    status pedido_status NOT NULL DEFAULT 'PENDENTE',
    valor_produtos DECIMAL(10,2) NOT NULL,
    valor_frete DECIMAL(10,2) NOT NULL,
    valor_total DECIMAL(10,2) NOT NULL,
    data_criacao TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    ultima_atualizacao TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_pedidos PRIMARY KEY (id),
    CONSTRAINT fk_pedidos_usuario FOREIGN KEY (usuario_id)
        REFERENCES usuarios (id) ON DELETE RESTRICT,
    CONSTRAINT fk_pedidos_entregador FOREIGN KEY (entregador_id)
        REFERENCES usuarios (id) ON DELETE RESTRICT,
    CONSTRAINT fk_pedidos_endereco FOREIGN KEY (endereco_id)
        REFERENCES enderecos (id) ON DELETE RESTRICT,
    CONSTRAINT chk_entregador_tipo CHECK (
        entregador_id IS NULL OR 
        EXISTS (
            SELECT 1 FROM usuarios 
            WHERE id = entregador_id AND tipo = 'ENTREGADOR'
        )
    ),
    CONSTRAINT chk_valores CHECK (
        valor_total = valor_produtos + valor_frete AND
        valor_produtos >= 0 AND
        valor_frete >= 0
    )
);

CREATE TABLE itens_pedido (
    id UUID DEFAULT uuid_generate_v4(),
    pedido_id UUID NOT NULL,
    produto_id UUID NOT NULL,
    quantidade INTEGER NOT NULL,
    preco_unitario DECIMAL(10,2) NOT NULL,
    subtotal DECIMAL(10,2) NOT NULL,
    CONSTRAINT pk_itens_pedido PRIMARY KEY (id),
    CONSTRAINT fk_itens_pedido_pedido FOREIGN KEY (pedido_id)
        REFERENCES pedidos (id) ON DELETE CASCADE,
    CONSTRAINT fk_itens_pedido_produto FOREIGN KEY (produto_id)
        REFERENCES produtos (id) ON DELETE RESTRICT,
    CONSTRAINT chk_quantidade_positiva CHECK (quantidade > 0),
    CONSTRAINT chk_preco_positivo CHECK (preco_unitario > 0),
    CONSTRAINT chk_subtotal CHECK (subtotal = quantidade * preco_unitario)
);

CREATE TABLE pagamentos (
    id UUID DEFAULT uuid_generate_v4(),
    pedido_id UUID NOT NULL,
    metodo pagamento_metodo NOT NULL,
    status pagamento_status NOT NULL DEFAULT 'PENDENTE',
    valor DECIMAL(10,2) NOT NULL,
    transacao_id VARCHAR(100),
    data_criacao TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    ultima_atualizacao TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_pagamentos PRIMARY KEY (id),
    CONSTRAINT fk_pagamentos_pedido FOREIGN KEY (pedido_id)
        REFERENCES pedidos (id) ON DELETE RESTRICT,
    CONSTRAINT chk_valor_positivo CHECK (valor > 0)
);

CREATE TABLE provas_entrega (
    id UUID DEFAULT uuid_generate_v4(),
    pedido_id UUID NOT NULL,
    foto_url VARCHAR(255) NOT NULL,
    data_hora TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    CONSTRAINT pk_provas_entrega PRIMARY KEY (id),
    CONSTRAINT fk_provas_entrega_pedido FOREIGN KEY (pedido_id)
        REFERENCES pedidos (id) ON DELETE CASCADE,
    CONSTRAINT chk_coordenadas CHECK (
        (latitude IS NULL AND longitude IS NULL) OR
        (latitude IS NOT NULL AND longitude IS NOT NULL AND
         latitude BETWEEN -90 AND 90 AND
         longitude BETWEEN -180 AND 180)
    )
);

CREATE TABLE notificacoes (
    id UUID DEFAULT uuid_generate_v4(),
    usuario_id UUID NOT NULL,
    pedido_id UUID,
    tipo notificacao_tipo NOT NULL,
    titulo VARCHAR(100) NOT NULL,
    mensagem TEXT NOT NULL,
    lida BOOLEAN DEFAULT false,
    data_criacao TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_notificacoes PRIMARY KEY (id),
    CONSTRAINT fk_notificacoes_usuario FOREIGN KEY (usuario_id)
        REFERENCES usuarios (id) ON DELETE CASCADE,
    CONSTRAINT fk_notificacoes_pedido FOREIGN KEY (pedido_id)
        REFERENCES pedidos (id) ON DELETE CASCADE
);

-- Índices
CREATE INDEX idx_usuarios_email ON usuarios(email);
CREATE INDEX idx_usuarios_tipo ON usuarios(tipo);
CREATE INDEX idx_enderecos_usuario ON enderecos(usuario_id);
CREATE INDEX idx_produtos_status ON produtos(status);
CREATE INDEX idx_pedidos_usuario ON pedidos(usuario_id);
CREATE INDEX idx_pedidos_entregador ON pedidos(entregador_id);
CREATE INDEX idx_pedidos_status ON pedidos(status);
CREATE INDEX idx_pagamentos_pedido ON pagamentos(pedido_id);
CREATE INDEX idx_pagamentos_status ON pagamentos(status);
CREATE INDEX idx_notificacoes_usuario ON notificacoes(usuario_id);
CREATE INDEX idx_notificacoes_lida ON notificacoes(lida);

-- Triggers para atualização automática de timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.ultima_atualizacao = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_usuarios_modtime
    BEFORE UPDATE ON usuarios
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_produtos_modtime
    BEFORE UPDATE ON produtos
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_carrinhos_modtime
    BEFORE UPDATE ON carrinhos
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_pedidos_modtime
    BEFORE UPDATE ON pedidos
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_pagamentos_modtime
    BEFORE UPDATE ON pagamentos
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Função para calcular o total do pedido
CREATE OR REPLACE FUNCTION calcular_total_pedido()
RETURNS TRIGGER AS $$
BEGIN
    NEW.valor_produtos = (
        SELECT COALESCE(SUM(subtotal), 0)
        FROM itens_pedido
        WHERE pedido_id = NEW.id
    );
    NEW.valor_total = NEW.valor_produtos + NEW.valor_frete;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER atualizar_total_pedido
    BEFORE INSERT OR UPDATE ON pedidos
    FOR EACH ROW
    EXECUTE FUNCTION calcular_total_pedido();

-- Comentários
COMMENT ON DATABASE cupcake_delivery IS 'Banco de dados do sistema de delivery de cupcakes';
COMMENT ON TABLE usuarios IS 'Armazena informações de todos os usuários do sistema';
COMMENT ON TABLE enderecos IS 'Endereços de entrega cadastrados pelos usuários';
COMMENT ON TABLE produtos IS 'Catálogo de cupcakes disponíveis para venda';
COMMENT ON TABLE carrinhos IS 'Carrinhos de compras ativos dos usuários';
COMMENT ON TABLE itens_carrinho IS 'Itens adicionados aos carrinhos de compras';
COMMENT ON TABLE pedidos IS 'Pedidos realizados pelos clientes';
COMMENT ON TABLE itens_pedido IS 'Itens incluídos em cada pedido';
COMMENT ON TABLE pagamentos IS 'Registro de pagamentos dos pedidos';
COMMENT ON TABLE provas_entrega IS 'Fotos e dados de confirmação de entrega';
COMMENT ON TABLE notificacoes IS 'Notificações enviadas aos usuários';
