-- Modelo Conceitual do Banco de Dados
-- Sistema de Delivery de Cupcakes

/*
ENTIDADES PRINCIPAIS:

1. USUÁRIO (Tabela base para autenticação)
   - Herança será implementada usando uma coluna "tipo" para identificar o perfil

2. PRODUTO (Cupcakes disponíveis)
   - Informações sobre os produtos, incluindo preço e status

3. PEDIDO (Compras realizadas)
   - Cabeçalho do pedido com informações gerais

4. CARRINHO (Carrinho de compras temporário)
   - Armazena itens ainda não convertidos em pedido

5. ENDEREÇO (Endereços de entrega)
   - Vinculado ao usuário e pode ser reutilizado em vários pedidos

RELACIONAMENTOS PRINCIPAIS:

1. Um USUÁRIO pode ter vários ENDEREÇOS
2. Um USUÁRIO pode ter vários PEDIDOS
3. Um PEDIDO tem vários ITENS_PEDIDO
4. Um CARRINHO pertence a um USUÁRIO
5. Um CARRINHO tem vários ITENS_CARRINHO
*/

-- Criação das tabelas conceituais (Primeira Forma Normal)

CREATE TABLE usuarios (
    id UUID PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    senha VARCHAR(255) NOT NULL,
    tipo VARCHAR(20) NOT NULL, -- 'CLIENTE', 'ADMIN', 'ENTREGADOR'
    telefone VARCHAR(20),
    placa_veiculo VARCHAR(10), -- Só para entregadores
    data_criacao TIMESTAMP NOT NULL,
    ultima_atualizacao TIMESTAMP NOT NULL
);

CREATE TABLE enderecos (
    id UUID PRIMARY KEY,
    usuario_id UUID NOT NULL,
    cep VARCHAR(9) NOT NULL,
    rua VARCHAR(100) NOT NULL,
    numero VARCHAR(10) NOT NULL,
    complemento VARCHAR(100),
    bairro VARCHAR(100) NOT NULL,
    cidade VARCHAR(100) NOT NULL,
    estado CHAR(2) NOT NULL,
    padrao BOOLEAN DEFAULT false,
    data_criacao TIMESTAMP NOT NULL
);

CREATE TABLE produtos (
    id UUID PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    descricao TEXT,
    preco DECIMAL(10,2) NOT NULL,
    imagem_url VARCHAR(255),
    categoria VARCHAR(50), -- 'TRADICIONAL', 'ESPECIAL'
    status VARCHAR(20) NOT NULL, -- 'ATIVO', 'INATIVO', 'ESGOTADO'
    data_criacao TIMESTAMP NOT NULL,
    ultima_atualizacao TIMESTAMP NOT NULL
);

CREATE TABLE carrinhos (
    id UUID PRIMARY KEY,
    usuario_id UUID NOT NULL,
    data_criacao TIMESTAMP NOT NULL,
    ultima_atualizacao TIMESTAMP NOT NULL
);

CREATE TABLE itens_carrinho (
    id UUID PRIMARY KEY,
    carrinho_id UUID NOT NULL,
    produto_id UUID NOT NULL,
    quantidade INTEGER NOT NULL,
    preco_unitario DECIMAL(10,2) NOT NULL,
    subtotal DECIMAL(10,2) NOT NULL,
    data_criacao TIMESTAMP NOT NULL
);

CREATE TABLE pedidos (
    id UUID PRIMARY KEY,
    usuario_id UUID NOT NULL,
    entregador_id UUID,
    endereco_id UUID NOT NULL,
    status VARCHAR(20) NOT NULL, -- 'PENDENTE', 'CONFIRMADO', 'PREPARANDO', 'ENVIADO', 'ENTREGUE', 'CANCELADO'
    valor_produtos DECIMAL(10,2) NOT NULL,
    valor_frete DECIMAL(10,2) NOT NULL,
    valor_total DECIMAL(10,2) NOT NULL,
    data_criacao TIMESTAMP NOT NULL,
    ultima_atualizacao TIMESTAMP NOT NULL
);

CREATE TABLE itens_pedido (
    id UUID PRIMARY KEY,
    pedido_id UUID NOT NULL,
    produto_id UUID NOT NULL,
    quantidade INTEGER NOT NULL,
    preco_unitario DECIMAL(10,2) NOT NULL,
    subtotal DECIMAL(10,2) NOT NULL
);

CREATE TABLE pagamentos (
    id UUID PRIMARY KEY,
    pedido_id UUID NOT NULL,
    metodo VARCHAR(50) NOT NULL, -- 'CARTAO', 'PIX'
    status VARCHAR(20) NOT NULL, -- 'PENDENTE', 'APROVADO', 'REJEITADO'
    valor DECIMAL(10,2) NOT NULL,
    transacao_id VARCHAR(100),
    data_criacao TIMESTAMP NOT NULL,
    ultima_atualizacao TIMESTAMP NOT NULL
);

CREATE TABLE provas_entrega (
    id UUID PRIMARY KEY,
    pedido_id UUID NOT NULL,
    foto_url VARCHAR(255) NOT NULL,
    data_hora TIMESTAMP NOT NULL,
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8)
);

CREATE TABLE notificacoes (
    id UUID PRIMARY KEY,
    usuario_id UUID NOT NULL,
    pedido_id UUID,
    tipo VARCHAR(50) NOT NULL, -- 'PEDIDO_CONFIRMADO', 'SAIU_ENTREGA', 'ENTREGUE', etc
    titulo VARCHAR(100) NOT NULL,
    mensagem TEXT NOT NULL,
    lida BOOLEAN DEFAULT false,
    data_criacao TIMESTAMP NOT NULL
);

-- Índices para otimização de consultas frequentes

CREATE INDEX idx_usuarios_email ON usuarios(email);
CREATE INDEX idx_usuarios_tipo ON usuarios(tipo);
CREATE INDEX idx_enderecos_usuario ON enderecos(usuario_id);
CREATE INDEX idx_produtos_status ON produtos(status);
CREATE INDEX idx_pedidos_usuario ON pedidos(usuario_id);
CREATE INDEX idx_pedidos_entregador ON pedidos(entregador_id);
CREATE INDEX idx_pedidos_status ON pedidos(status);
CREATE INDEX idx_itens_pedido_pedido ON itens_pedido(pedido_id);
CREATE INDEX idx_pagamentos_pedido ON pagamentos(pedido_id);
CREATE INDEX idx_notificacoes_usuario ON notificacoes(usuario_id);

-- Chaves estrangeiras para garantir integridade referencial

ALTER TABLE enderecos
    ADD CONSTRAINT fk_enderecos_usuario
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id);

ALTER TABLE carrinhos
    ADD CONSTRAINT fk_carrinhos_usuario
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id);

ALTER TABLE itens_carrinho
    ADD CONSTRAINT fk_itens_carrinho_carrinho
    FOREIGN KEY (carrinho_id) REFERENCES carrinhos(id),
    ADD CONSTRAINT fk_itens_carrinho_produto
    FOREIGN KEY (produto_id) REFERENCES produtos(id);

ALTER TABLE pedidos
    ADD CONSTRAINT fk_pedidos_usuario
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id),
    ADD CONSTRAINT fk_pedidos_entregador
    FOREIGN KEY (entregador_id) REFERENCES usuarios(id),
    ADD CONSTRAINT fk_pedidos_endereco
    FOREIGN KEY (endereco_id) REFERENCES enderecos(id);

ALTER TABLE itens_pedido
    ADD CONSTRAINT fk_itens_pedido_pedido
    FOREIGN KEY (pedido_id) REFERENCES pedidos(id),
    ADD CONSTRAINT fk_itens_pedido_produto
    FOREIGN KEY (produto_id) REFERENCES produtos(id);

ALTER TABLE pagamentos
    ADD CONSTRAINT fk_pagamentos_pedido
    FOREIGN KEY (pedido_id) REFERENCES pedidos(id);

ALTER TABLE provas_entrega
    ADD CONSTRAINT fk_provas_entrega_pedido
    FOREIGN KEY (pedido_id) REFERENCES pedidos(id);

ALTER TABLE notificacoes
    ADD CONSTRAINT fk_notificacoes_usuario
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id),
    ADD CONSTRAINT fk_notificacoes_pedido
    FOREIGN KEY (pedido_id) REFERENCES pedidos(id);

-- Comentários das tabelas para documentação

COMMENT ON TABLE usuarios IS 'Armazena informações de todos os usuários do sistema (clientes, administradores e entregadores)';
COMMENT ON TABLE enderecos IS 'Endereços de entrega cadastrados pelos usuários';
COMMENT ON TABLE produtos IS 'Catálogo de cupcakes disponíveis para venda';
COMMENT ON TABLE carrinhos IS 'Carrinhos de compras ativos dos usuários';
COMMENT ON TABLE itens_carrinho IS 'Itens adicionados aos carrinhos de compras';
COMMENT ON TABLE pedidos IS 'Pedidos realizados pelos clientes';
COMMENT ON TABLE itens_pedido IS 'Itens incluídos em cada pedido';
COMMENT ON TABLE pagamentos IS 'Registro de pagamentos dos pedidos';
COMMENT ON TABLE provas_entrega IS 'Fotos e dados de confirmação de entrega';
COMMENT ON TABLE notificacoes IS 'Notificações enviadas aos usuários';
