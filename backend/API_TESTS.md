# Cupcake Delivery API - Testes

## Base URL: http://localhost:8080

### 1. Health Check
```http
GET /health
```

### 2. Registro de Cliente
```http
POST /register
Content-Type: application/json

{
  "name": "João Silva",
  "email": "joao@email.com",
  "password": "123456",
  "type": "customer"
}
```

### 3. Registro de Entregador
```http
POST /register
Content-Type: application/json

{
  "name": "Maria Santos",
  "email": "maria@email.com", 
  "password": "123456",
  "type": "delivery",
  "vehicle": "bicicleta"
}
```

### 4. Login
```http
POST /login
Content-Type: application/json

{
  "email": "admin@cupcakedelivery.com",
  "password": "admin123"
}
```

### 5. Login de Cliente
```http
POST /login
Content-Type: application/json

{
  "email": "joao@email.com",
  "password": "123456"
}
```

### 6. Listar Produtos (público)
```http
GET /products
```

### 7. Ver Produto Específico (público)
```http
GET /products/1
```

### 8. Criar Produto (ADMIN APENAS)
```http
POST /products
Authorization: Bearer {TOKEN_ADMIN}
Content-Type: application/json

{
  "name": "Cupcake Especial",
  "description": "Cupcake especial do dia",
  "price": 10.50,
  "imageUrl": "https://via.placeholder.com/300x300/FF69B4/FFFFFF?text=Especial"
}
```

### 9. Atualizar Produto (ADMIN APENAS)
```http
PUT /products/1
Authorization: Bearer {TOKEN_ADMIN}
Content-Type: application/json

{
  "name": "Cupcake de Chocolate Premium",
  "description": "Delicioso cupcake de chocolate belga premium",
  "price": 12.00,
  "imageUrl": "https://via.placeholder.com/300x300/8B4513/FFFFFF?text=Premium"
}
```

### 10. Deletar Produto (ADMIN APENAS)
```http
DELETE /products/7
Authorization: Bearer {TOKEN_ADMIN}
```

### 11. Criar Pedido (CLIENTE AUTENTICADO)
```http
POST /orders
Authorization: Bearer {TOKEN_CLIENTE}
Content-Type: application/json

{
  "items": [
    {
      "product_id": 1,
      "quantity": 2
    },
    {
      "product_id": 3,
      "quantity": 1
    }
  ]
}
```

### 12. Listar Pedidos (AUTENTICADO)
```http
GET /orders
Authorization: Bearer {TOKEN}
```

### 13. Atualizar Status do Pedido (ENTREGADOR/ADMIN)
```http
PUT /orders/1/status
Authorization: Bearer {TOKEN_ENTREGADOR_OU_ADMIN}
Content-Type: application/x-www-form-urlencoded

status=delivering
```

### 14. Marcar Pedido como Entregue (ENTREGADOR)
```http
PUT /orders/1/status  
Authorization: Bearer {TOKEN_ENTREGADOR}
Content-Type: application/x-www-form-urlencoded

status=delivered
```

## Fluxo de Teste Completo:

1. **Verificar se o servidor está rodando**: Health Check
2. **Criar contas**: Registrar cliente e entregador
3. **Fazer login como admin**: Usar credenciais admin@cupcakedelivery.com / admin123
4. **Testar produtos**: Listar, ver, criar, atualizar, deletar (admin)
5. **Fazer login como cliente**: Usar credenciais do cliente criado
6. **Criar pedido**: Cliente faz pedido com produtos existentes
7. **Fazer login como entregador**: Usar credenciais do entregador
8. **Atualizar status**: Entregador pega e entrega pedido

## Observações:

- Substitua `{TOKEN}` pelo token JWT retornado no login
- Todos os endpoints protegidos precisam do header Authorization
- Admin pode fazer tudo
- Cliente pode apenas criar pedidos e ver seus pedidos
- Entregador pode ver pedidos disponíveis e atualizar status dos seus
