# Dicionário de Dados - Sistema de Delivery de Cupcakes

## Tabela: usuarios
| Campo | Tipo | Descrição | Restrições |
|-------|------|-----------|------------|
| id | UUID | Identificador único do usuário | Chave primária |
| nome | VARCHAR(100) | Nome completo do usuário | NOT NULL |
| email | VARCHAR(100) | Email do usuário (usado para login) | NOT NULL, UNIQUE |
| senha | VARCHAR(255) | Senha criptografada | NOT NULL |
| tipo | VARCHAR(20) | Tipo do usuário (CLIENTE, ADMIN, ENTREGADOR) | NOT NULL |
| telefone | VARCHAR(20) | Número de telefone | Opcional |
| placa_veiculo | VARCHAR(10) | Placa do veículo (apenas para entregadores) | Opcional |
| data_criacao | TIMESTAMP | Data e hora de criação do registro | NOT NULL |
| ultima_atualizacao | TIMESTAMP | Data e hora da última atualização | NOT NULL |

## Tabela: enderecos
| Campo | Tipo | Descrição | Restrições |
|-------|------|-----------|------------|
| id | UUID | Identificador único do endereço | Chave primária |
| usuario_id | UUID | ID do usuário dono do endereço | FK usuarios(id) |
| cep | VARCHAR(9) | CEP do endereço | NOT NULL |
| rua | VARCHAR(100) | Nome da rua | NOT NULL |
| numero | VARCHAR(10) | Número do endereço | NOT NULL |
| complemento | VARCHAR(100) | Complemento do endereço | Opcional |
| bairro | VARCHAR(100) | Bairro | NOT NULL |
| cidade | VARCHAR(100) | Cidade | NOT NULL |
| estado | CHAR(2) | Sigla do estado | NOT NULL |
| padrao | BOOLEAN | Indica se é o endereço padrão | DEFAULT false |
| data_criacao | TIMESTAMP | Data e hora de criação do registro | NOT NULL |

## Tabela: produtos
| Campo | Tipo | Descrição | Restrições |
|-------|------|-----------|------------|
| id | UUID | Identificador único do produto | Chave primária |
| nome | VARCHAR(100) | Nome do cupcake | NOT NULL |
| descricao | TEXT | Descrição detalhada do cupcake | Opcional |
| preco | DECIMAL(10,2) | Preço unitário | NOT NULL |
| imagem_url | VARCHAR(255) | URL da imagem do produto | Opcional |
| categoria | VARCHAR(50) | Categoria (TRADICIONAL, ESPECIAL) | Opcional |
| status | VARCHAR(20) | Status do produto | NOT NULL |
| data_criacao | TIMESTAMP | Data e hora de criação do registro | NOT NULL |
| ultima_atualizacao | TIMESTAMP | Data e hora da última atualização | NOT NULL |

## Tabela: carrinhos
| Campo | Tipo | Descrição | Restrições |
|-------|------|-----------|------------|
| id | UUID | Identificador único do carrinho | Chave primária |
| usuario_id | UUID | ID do usuário dono do carrinho | FK usuarios(id) |
| data_criacao | TIMESTAMP | Data e hora de criação do registro | NOT NULL |
| ultima_atualizacao | TIMESTAMP | Data e hora da última atualização | NOT NULL |

## Tabela: itens_carrinho
| Campo | Tipo | Descrição | Restrições |
|-------|------|-----------|------------|
| id | UUID | Identificador único do item | Chave primária |
| carrinho_id | UUID | ID do carrinho | FK carrinhos(id) |
| produto_id | UUID | ID do produto | FK produtos(id) |
| quantidade | INTEGER | Quantidade do item | NOT NULL |
| preco_unitario | DECIMAL(10,2) | Preço unitário no momento da adição | NOT NULL |
| subtotal | DECIMAL(10,2) | Subtotal (quantidade * preço) | NOT NULL |
| data_criacao | TIMESTAMP | Data e hora de criação do registro | NOT NULL |

## Tabela: pedidos
| Campo | Tipo | Descrição | Restrições |
|-------|------|-----------|------------|
| id | UUID | Identificador único do pedido | Chave primária |
| usuario_id | UUID | ID do cliente | FK usuarios(id) |
| entregador_id | UUID | ID do entregador designado | FK usuarios(id) |
| endereco_id | UUID | ID do endereço de entrega | FK enderecos(id) |
| status | VARCHAR(20) | Status do pedido | NOT NULL |
| valor_produtos | DECIMAL(10,2) | Valor total dos produtos | NOT NULL |
| valor_frete | DECIMAL(10,2) | Valor do frete | NOT NULL |
| valor_total | DECIMAL(10,2) | Valor total (produtos + frete) | NOT NULL |
| data_criacao | TIMESTAMP | Data e hora de criação do registro | NOT NULL |
| ultima_atualizacao | TIMESTAMP | Data e hora da última atualização | NOT NULL |

## Tabela: itens_pedido
| Campo | Tipo | Descrição | Restrições |
|-------|------|-----------|------------|
| id | UUID | Identificador único do item | Chave primária |
| pedido_id | UUID | ID do pedido | FK pedidos(id) |
| produto_id | UUID | ID do produto | FK produtos(id) |
| quantidade | INTEGER | Quantidade do item | NOT NULL |
| preco_unitario | DECIMAL(10,2) | Preço unitário no momento da compra | NOT NULL |
| subtotal | DECIMAL(10,2) | Subtotal (quantidade * preço) | NOT NULL |

## Tabela: pagamentos
| Campo | Tipo | Descrição | Restrições |
|-------|------|-----------|------------|
| id | UUID | Identificador único do pagamento | Chave primária |
| pedido_id | UUID | ID do pedido | FK pedidos(id) |
| metodo | VARCHAR(50) | Método de pagamento | NOT NULL |
| status | VARCHAR(20) | Status do pagamento | NOT NULL |
| valor | DECIMAL(10,2) | Valor do pagamento | NOT NULL |
| transacao_id | VARCHAR(100) | ID da transação (gateway) | Opcional |
| data_criacao | TIMESTAMP | Data e hora de criação do registro | NOT NULL |
| ultima_atualizacao | TIMESTAMP | Data e hora da última atualização | NOT NULL |

## Tabela: provas_entrega
| Campo | Tipo | Descrição | Restrições |
|-------|------|-----------|------------|
| id | UUID | Identificador único da prova | Chave primária |
| pedido_id | UUID | ID do pedido | FK pedidos(id) |
| foto_url | VARCHAR(255) | URL da foto de confirmação | NOT NULL |
| data_hora | TIMESTAMP | Data e hora do registro | NOT NULL |
| latitude | DECIMAL(10,8) | Latitude da localização | Opcional |
| longitude | DECIMAL(11,8) | Longitude da localização | Opcional |

## Tabela: notificacoes
| Campo | Tipo | Descrição | Restrições |
|-------|------|-----------|------------|
| id | UUID | Identificador único da notificação | Chave primária |
| usuario_id | UUID | ID do usuário destinatário | FK usuarios(id) |
| pedido_id | UUID | ID do pedido relacionado | FK pedidos(id) |
| tipo | VARCHAR(50) | Tipo da notificação | NOT NULL |
| titulo | VARCHAR(100) | Título da notificação | NOT NULL |
| mensagem | TEXT | Conteúdo da notificação | NOT NULL |
| lida | BOOLEAN | Indica se foi lida | DEFAULT false |
| data_criacao | TIMESTAMP | Data e hora de criação | NOT NULL |
