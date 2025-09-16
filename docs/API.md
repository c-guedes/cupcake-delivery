# 🔗 API Documentation

Documentação completa da API REST do Sistema de Delivery de Cupcakes.

## Base URL
```
http://localhost:8080
```

## Autenticação

A API utiliza JWT (JSON Web Tokens) para autenticação. O token deve ser incluído no header Authorization:

```
Authorization: Bearer <seu_jwt_token>
```

## Estrutura de Resposta

### Sucesso
```json
{
  "data": { /* dados solicitados */ },
  "message": "Operação realizada com sucesso"
}
```

### Erro
```json
{
  "error": "error_type",
  "message": "Mensagem de erro legível",
  "validations": [
    {
      "field": "campo_com_erro",
      "message": "Descrição do erro específico"
    }
  ]
}
```

## Endpoints

### 🔐 Autenticação

#### Registrar Usuário
```http
POST /register
```

**Body:**
```json
{
  "name": "João Silva",
  "email": "joao@email.com",
  "password": "senha123",
  "type": "customer" // customer, delivery, admin
}
```

**Resposta (201):**
```json
{
  "message": "Usuário criado com sucesso",
  "user": {
    "id": 1,
    "name": "João Silva",
    "email": "joao@email.com",
    "type": "customer"
  }
}
```

#### Login
```http
POST /login
```

**Body:**
```json
{
  "email": "joao@email.com",
  "password": "senha123"
}
```

**Resposta (200):**
```json
{
  "message": "Login realizado com sucesso",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "João Silva",
    "email": "joao@email.com",
    "type": "customer"
  }
}
```

### 🧁 Produtos

#### Listar Produtos
```http
GET /products
```

**Resposta (200):**
```json
[
  {
    "id": 1,
    "name": "Cupcake de Chocolate",
    "description": "Delicioso cupcake com cobertura de chocolate",
    "price": 8.50,
    "imageUrl": "https://example.com/chocolate.jpg",
    "createdAt": "2024-01-15T10:30:00Z"
  }
]
```

#### Criar Produto
```http
POST /products
```
**Autenticação:** Requerida (Admin)

**Body:**
```json
{
  "name": "Cupcake de Morango",
  "description": "Cupcake sabor morango com chantilly",
  "price": 7.50,
  "imageUrl": "https://example.com/morango.jpg"
}
```

**Resposta (201):**
```json
{
  "message": "Produto criado com sucesso",
  "product": {
    "id": 2,
    "name": "Cupcake de Morango",
    "description": "Cupcake sabor morango com chantilly",
    "price": 7.50,
    "imageUrl": "https://example.com/morango.jpg"
  }
}
```

#### Atualizar Produto
```http
PUT /products/:id
```
**Autenticação:** Requerida (Admin)

**Body:**
```json
{
  "name": "Cupcake de Morango Premium",
  "price": 9.50
}
```

#### Deletar Produto
```http
DELETE /products/:id
```
**Autenticação:** Requerida (Admin)

### 📦 Pedidos

#### Criar Pedido
```http
POST /orders
```
**Autenticação:** Requerida

**Body:**
```json
{
  "items": [
    {
      "productId": 1,
      "quantity": 2
    },
    {
      "productId": 2,
      "quantity": 1
    }
  ],
  "deliveryAddress": "Rua das Flores, 123, Centro",
  "paymentMethod": "credit_card"
}
```

**Resposta (201):**
```json
{
  "message": "Pedido criado com sucesso",
  "order": {
    "id": 15,
    "userId": 1,
    "status": "pending",
    "total": 24.50,
    "deliveryAddress": "Rua das Flores, 123, Centro",
    "paymentMethod": "credit_card",
    "items": [
      {
        "id": 1,
        "productId": 1,
        "quantity": 2,
        "price": 8.50,
        "product": {
          "name": "Cupcake de Chocolate"
        }
      }
    ],
    "createdAt": "2024-01-15T14:30:00Z"
  }
}
```

#### Listar Pedidos
```http
GET /orders
```
**Autenticação:** Requerida

**Query Parameters:**
- `status`: Filtrar por status (pending, preparing, ready, delivering, delivered)
- `page`: Página (padrão: 1)
- `limit`: Itens por página (padrão: 10)

**Resposta (200):**
```json
{
  "orders": [
    {
      "id": 15,
      "userId": 1,
      "status": "preparing",
      "total": 24.50,
      "deliveryAddress": "Rua das Flores, 123, Centro",
      "items": [
        {
          "productId": 1,
          "quantity": 2,
          "price": 8.50,
          "product": {
            "name": "Cupcake de Chocolate"
          }
        }
      ],
      "createdAt": "2024-01-15T14:30:00Z",
      "updatedAt": "2024-01-15T15:00:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "limit": 10
}
```

#### Atualizar Status do Pedido
```http
PUT /orders/:id/status
```
**Autenticação:** Requerida (Admin ou Delivery)

**Body:**
```json
{
  "status": "ready",
  "notes": "Pedido pronto para entrega"
}
```

### 🔔 Notificações

#### Listar Notificações
```http
GET /notifications
```
**Autenticação:** Requerida

**Query Parameters:**
- `unread`: `true` para apenas não lidas
- `page`: Página (padrão: 1)
- `limit`: Itens por página (padrão: 20)

**Resposta (200):**
```json
{
  "notifications": [
    {
      "id": 1,
      "userId": 1,
      "type": "order_status",
      "title": "Pedido Atualizado",
      "message": "Seu pedido #15 está sendo preparado",
      "data": {
        "orderId": 15,
        "status": "preparing"
      },
      "isRead": false,
      "createdAt": "2024-01-15T15:00:00Z"
    }
  ],
  "unreadCount": 3
}
```

#### Marcar Notificação como Lida
```http
PUT /notifications/:id/read
```
**Autenticação:** Requerida

## Status dos Pedidos

| Status | Descrição |
|--------|-----------|
| `pending` | Aguardando confirmação |
| `preparing` | Em preparação |
| `ready` | Pronto para entrega |
| `delivering` | Em rota de entrega |
| `delivered` | Entregue |

## Códigos de Erro

| Código | Descrição |
|--------|-----------|
| 400 | Bad Request - Dados inválidos |
| 401 | Unauthorized - Token inválido/ausente |
| 403 | Forbidden - Sem permissão |
| 404 | Not Found - Recurso não encontrado |
| 422 | Unprocessable Entity - Erro de validação |
| 500 | Internal Server Error - Erro interno |

## Tipos de Erro

### `validation_error`
Erro de validação de dados de entrada.

### `authentication_error`
Erro de autenticação (credenciais inválidas).

### `authorization_error`
Erro de autorização (sem permissão).

### `not_found_error`
Recurso solicitado não encontrado.

### `server_error`
Erro interno do servidor.

## Middleware de Autenticação

Rotas protegidas requerem header de autorização:
```
Authorization: Bearer <jwt_token>
```

O token JWT contém:
```json
{
  "user_id": 1,
  "user_type": "customer",
  "exp": 1642012345
}
```

## Rate Limiting

- **Limite**: 100 requests por minuto por IP
- **Header de resposta**: `X-RateLimit-*`

## CORS

A API aceita requests de:
- `http://localhost:5173` (desenvolvimento)
- Configurável via variável de ambiente `FRONTEND_URL`

## Variáveis de Ambiente

```bash
# Banco de dados
DATABASE_URL="postgres://user:pass@localhost/cupcake_delivery"

# JWT
JWT_SECRET="seu_secret_super_seguro"
JWT_EXPIRY="24h"

# Servidor
PORT="8080"
GIN_MODE="release" # development/release

# CORS
FRONTEND_URL="http://localhost:5173"
```

## Exemplos de Uso

### Fluxo Completo de Pedido

1. **Registrar/Login**
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"name":"João","email":"joao@email.com","password":"123456","type":"customer"}'
```

2. **Listar Produtos**
```bash
curl -X GET http://localhost:8080/products
```

3. **Criar Pedido**
```bash
curl -X POST http://localhost:8080/orders \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"items":[{"productId":1,"quantity":2}],"deliveryAddress":"Rua A, 123"}'
```

4. **Acompanhar Pedido**
```bash
curl -X GET http://localhost:8080/orders \
  -H "Authorization: Bearer <token>"
```
