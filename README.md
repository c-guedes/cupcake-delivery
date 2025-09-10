# Cupcake Delivery

Sistema de delivery de cupcakes desenvolvido com React e Go.

## Estrutura do Projeto

```
cupcake-delivery/
├── docs/               # Documentação
│   ├── diagramas/     # Diagramas UML
│   ├── mockups/       # Protótipos
│   └── database/      # Scripts e docs do banco
├── frontend/          # Aplicação React
├── backend/           # API em Go
└── README.md
```

## Tecnologias Utilizadas

### Frontend
- React
- TypeScript
- Vite
- Tailwind CSS
- React Router
- React Query

### Backend
- Go
- Gin (Web Framework)
- GORM (ORM)
- JWT (Autenticação)
- PostgreSQL

### Banco de Dados
- PostgreSQL 14+

## Pré-requisitos

- Node.js 18+
- Go 1.20+
- PostgreSQL 14+
- Git

## Configuração do Ambiente

### Backend (Go)

1. Instalar Go 1.20 ou superior
2. Configurar variáveis de ambiente:
```bash
export POSTGRES_URL="postgres://user:password@localhost:5432/cupcake_delivery"
export JWT_SECRET="seu_secret_aqui"
```

### Frontend (React)

1. Instalar Node.js 18 ou superior
2. Instalar dependências:
```bash
cd frontend
npm install
```

### Banco de Dados

1. Criar banco de dados:
```sql
CREATE DATABASE cupcake_delivery;
```

2. Executar scripts de migração:
```bash
cd docs/database
psql -U postgres -d cupcake_delivery -f create_database.sql
```

## Executando o Projeto

### Backend
```bash
cd backend
go run main.go
```

### Frontend
```bash
cd frontend
npm run dev
```

## Funcionalidades

### Cliente
- Cadastro e login de usuários
- Visualização do catálogo de cupcakes
- Gerenciamento do carrinho de compras
- Checkout com múltiplas formas de pagamento
- Acompanhamento de pedidos
- Histórico de pedidos

### Entregador
- Login no sistema
- Visualização de pedidos disponíveis
- Atualização de status de entrega
- Upload de foto de confirmação

### Administrador
- Gerenciamento de produtos
- Atribuição de pedidos
- Visualização de relatórios
- Gestão de usuários

## Documentação

- [Diagramas UML](./docs/diagramas/)
- [Mockups](./docs/mockups/)
- [Documentação do Banco](./docs/database/)

## Time

- Caique Guedes de Almeida

## Licença

Este projeto está sob a licença MIT.
