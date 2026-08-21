# 🧁 Cupcake Delivery System

Sistema completo de delivery de cupcakes desenvolvido como monorepo com React, Go, PostgreSQL e arquitetura moderna.

## 📋 Visão Geral

Este projeto implementa um sistema de delivery de cupcakes com três tipos de usuários:
- **Clientes**: Fazem pedidos e acompanham entregas
- **Administradores**: Gerenciam produtos, pedidos e usuários  
- **Entregadores**: Recebem e entregam pedidos

## 🏗️ Arquitetura

### Backend (Go)
- **Framework**: Gin Web Framework
- **ORM**: GORM 
- **Banco**: PostgreSQL
- **Autenticação**: JWT
- **Validação**: Validadores customizados
- **Testes**: Go testing + testify + mocks

### Frontend (React)
- **Framework**: React 18 + TypeScript
- **Build**: Vite
- **Styling**: TailwindCSS
- **Roteamento**: React Router v7
- **Estado**: Context API
- **Testes**: Jest + Testing Library

### Funcionalidades Principais
- ✅ Autenticação JWT multi-perfil
- ✅ CRUD completo de produtos
- ✅ Sistema de pedidos com status
- ✅ Carrinho de compras persistente
- ✅ Dashboard específico por tipo de usuário
- ✅ Sistema de notificações
- ✅ Dark mode
- ✅ Tratamento robusto de erros
- ✅ Testes automatizados

## 🚀 Como Executar

### Pré-requisitos
- Go 1.19+
- Node.js 18+
- PostgreSQL 13+

### 1. Backend
```bash
cd backend
cp .env.example .env
# Configure as variáveis no .env
go mod download
go run cmd/api/main.go
```

### 2. Frontend  
```bash
cd frontend
npm install
npm run dev
```

### 3. Banco de Dados
```sql
-- No PostgreSQL
CREATE DATABASE cupcake_delivery;
-- As migrações rodarão automaticamente
```

## 🧪 Testes

### Backend
```bash
cd backend
go test ./...                    # Todos os testes
go test ./internal/validators/   # Testes de validação
go test ./internal/handlers/     # Testes de handlers
```

### Frontend
```bash
cd frontend
npm test                    # Todos os testes
npm run test:coverage      # Com cobertura
npm run test:watch         # Watch mode
```

## 📚 Documentação para entrega — PIT II

- [Diagramas UML finais validados](docs/diagramas/README.md): classes, casos de uso, sequência e banco de dados.
- [Mockups e IHC validados](docs/mockups/README.md).
- [Modelo físico](docs/database/modelo_fisico.md) e [migração SQL](docs/database/migration.sql).
- [Dicionário de dados](docs/database/dicionario_dados.md).
- [Relatório de validação técnica](docs/VALIDACAO_TECNICA_PIT_II.md).
- [Capturas técnicas da aplicação em execução](docs/pit_ii/evidencias/README.md).
- [Checklist de entrega do PIT II](docs/PIT_II_CHECKLIST.md).
- [Auditoria de entregáveis e pendências externas](docs/ENTREGAVEIS_PIT_II.md).

> O repositório comprova a implementação e a validação técnica. As cinco fichas de avaliação precisam ser preenchidas por participantes reais antes do envio acadêmico.

## 📁 Estrutura do Projeto

```
cupcake-delivery/
├── backend/
│   ├── cmd/api/                 # Ponto de entrada
│   ├── internal/
│   │   ├── config/             # Configurações
│   │   ├── database/           # Conexão DB
│   │   ├── handlers/           # Controllers
│   │   ├── middleware/         # Middlewares
│   │   ├── models/             # Modelos GORM
│   │   ├── router/             # Rotas
│   │   ├── services/           # Lógica de negócio
│   │   ├── utils/              # Utilidades
│   │   └── validators/         # Validadores
│   └── scripts/                # Scripts auxiliares
├── frontend/
│   ├── src/
│   │   ├── components/         # Componentes React
│   │   ├── contexts/           # Contexts (Auth, Cart, Theme)
│   │   ├── hooks/              # Custom hooks
│   │   ├── pages/              # Páginas/rotas
│   │   ├── services/           # API client
│   │   └── types/              # TypeScript types
│   └── public/                 # Assets estáticos
└── docs/                       # Documentação
```

## 🔗 API Endpoints

### Autenticação
- `POST /register` - Cadastro de usuário
- `POST /login` - Login

### Produtos
- `GET /products` - Listar produtos
- `POST /products` - Criar produto (admin)
- `PUT /products/:id` - Atualizar produto (admin)
- `DELETE /products/:id` - Deletar produto (admin)

### Pedidos
- `POST /orders` - Criar pedido
- `GET /orders` - Listar pedidos
- `PUT /orders/:id/status` - Atualizar status

### Notificações
- `GET /notifications` - Listar notificações
- `PUT /notifications/:id/read` - Marcar como lida

## 🎨 UI/UX

### Dark Mode
Sistema completo de dark mode usando TailwindCSS com suporte a:
- Toggle automático
- Persistência no localStorage
- Transições suaves

### Componentes Principais
- **Navbar**: Navegação responsiva com notificações
- **Dashboard**: Específico por tipo de usuário
- **Cart**: Modal/dropdown com persistência
- **Toast**: Sistema de notificações temporárias
- **ErrorDisplay**: Tratamento visual de erros

## 🔒 Segurança

- JWT com expiração configurável
- Middleware de autenticação
- Validação de dados no backend
- Sanitização de inputs
- Proteção contra XSS/CSRF

## 📊 Status dos Pedidos

1. **Pendente** → Aguardando confirmação
2. **Preparando** → Em produção
3. **Pronto** → Aguardando entregador
4. **Entregando** → Em rota de entrega
5. **Entregue** → Concluído

## 🧑‍💻 Desenvolvimento

### Padrões Adotados
- Clean Architecture
- Repository Pattern
- Context Pattern (React)
- Error Boundary Pattern
- Custom Hooks Pattern

### Ferramentas de Qualidade
- ESLint + Prettier (Frontend)
- Go fmt + vet (Backend)
- Testes automatizados
- TypeScript strict mode

## 📄 Licença

Este projeto foi desenvolvido para fins acadêmicos como parte do curso de Análise e Desenvolvimento de Sistemas.

## 👥 Autor

**Caique Guedes de Almeida**
- Projeto Integrador de Competências (PIC)
- Análise e Desenvolvimento de Sistemas

---

Para mais detalhes, consulte a documentação específica em cada módulo ou entre em contato.
