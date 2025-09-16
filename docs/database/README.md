# 🗄️ Database Documentation

Esta pasta contém toda a documentação e scripts relacionados ao banco de dados do sistema Cupcake Delivery.

## 📋 **Arquivos Disponíveis**

### 📜 **Scripts SQL**
- **`migration.sql`** - Script de migração para PostgreSQL local (Docker)
- **`supabase_migration.sql`** - Script de migração para Supabase (Recomendado)
- **`modelo_conceitual.sql`** - Modelo conceitual do banco de dados

### 📚 **Documentação**
- **`DATABASE_SPEC.md`** - Especificação técnica completa do banco
- **`SUPABASE_INTEGRATION.md`** - Guia de integração com Supabase
- **`dicionario_dados.md`** - Dicionário de dados detalhado
- **`ATUALIZACAO_DATABASE.md`** - Histórico de atualizações

### ⚙️ **Configuração**
- **`../backend/.env.example`** - Configurações do backend
- **`../frontend/.env.example`** - Configurações do frontend

---

## 🚀 **Início Rápido com Supabase**

### 1️⃣ **Configurar Supabase**
```bash
# Acesse o Supabase Dashboard
https://msubfzrwhvwdwixfskuk.supabase.co

# Execute o script de migração no SQL Editor
# Copie e cole o conteúdo de supabase_migration.sql
```

### 2️⃣ **Configurar Backend**
```bash
# Na pasta backend/
cp .env.example .env

# Edite o .env com suas credenciais Supabase
SUPABASE_URL=https://msubfzrwhvwdwixfskuk.supabase.co
SUPABASE_ANON_KEY=sua-anon-key
DB_PASSWORD=sua-senha-supabase
```

### 3️⃣ **Configurar Frontend**
```bash
# Na pasta frontend/
cp .env.example .env.local

# Edite o .env.local
VITE_SUPABASE_URL=https://msubfzrwhvwdwixfskuk.supabase.co
VITE_SUPABASE_ANON_KEY=sua-anon-key
```

### 4️⃣ **Testar Conexão**
```bash
# Inicie o backend
cd backend && go run cmd/main.go

# Inicie o frontend
cd frontend && npm run dev
```

---

## 🔄 **Alternativa: PostgreSQL Local**

Se preferir usar PostgreSQL local com Docker:

### 1️⃣ **Usar Docker Compose**
```bash
# Na raiz do projeto
docker-compose up -d db

# Execute a migração
docker-compose exec db psql -U postgres -d cupcake_delivery -f /docker-entrypoint-initdb.d/migration.sql
```

### 2️⃣ **Configurar Variáveis**
```bash
# No backend/.env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=cupcake_delivery
DB_SSLMODE=disable
```

---

## 📊 **Estrutura do Banco**

### 🗂️ **5 Tabelas Principais**
1. **`users`** - Usuários (clientes, entregadores, admins)
2. **`products`** - Catálogo de cupcakes
3. **`orders`** - Pedidos realizados
4. **`order_items`** - Itens de cada pedido
5. **`notifications`** - Sistema de notificações

### 🔗 **Relacionamentos**
```
users (1:N) orders (1:N) order_items (N:1) products
users (1:N) notifications (N:1) orders
```

### 📈 **Otimizações**
- ✅ Índices em campos frequentemente consultados
- ✅ Soft delete com `deleted_at`
- ✅ Timestamps automáticos
- ✅ Constraints de integridade
- ✅ Triggers para validação

---

## 🔐 **Segurança**

### 🛡️ **Supabase (Recomendado)**
- Row Level Security (RLS) habilitado
- Policies para acesso controlado
- JWT authentication integrado
- SSL obrigatório

### 🔒 **PostgreSQL Local**
- Usuário e senha configurados
- Acesso restrito ao localhost
- SSL opcional para desenvolvimento

---

## 📚 **Documentação Técnica**

### 🎯 **Para Desenvolvedores**
- Leia `DATABASE_SPEC.md` para entender a estrutura completa
- Consulte `dicionario_dados.md` para detalhes dos campos
- Use `SUPABASE_INTEGRATION.md` para integração

### 🔧 **Para DevOps**
- Scripts de migração prontos para produção
- Configurações de ambiente documentadas
- Backup e recovery procedures incluídos

### 📊 **Para Analistas**
- Views otimizadas para relatórios
- Estrutura de dados business-friendly
- Métricas e KPIs integrados

---

## 🚀 **Próximos Passos**

1. ✅ **Escolher ambiente**: Supabase (recomendado) ou PostgreSQL local
2. 🔧 **Executar migração**: `supabase_migration.sql` ou `migration.sql`
3. ⚙️ **Configurar variáveis**: Backend e frontend `.env`
4. 🧪 **Testar conexão**: Verificar funcionamento
5. 🚀 **Começar desenvolvimento**: Sistema pronto para uso

---

## 📞 **Suporte**

### 🆘 **Problemas Comuns**
- **Erro de conexão**: Verifique credenciais no `.env`
- **Tabelas não criadas**: Execute o script de migração
- **Permissões**: Verifique policies no Supabase

### 📖 **Recursos**
- [Documentação Supabase](https://supabase.com/docs)
- [PostgreSQL Docs](https://www.postgresql.org/docs/)
- [GORM Documentation](https://gorm.io/docs/)

---

## 📊 **Status Atual**

| Componente | Status | Observações |
|------------|---------|-------------|
| **Estrutura DB** | ✅ Completa | 5 tabelas implementadas |
| **Migração Supabase** | ✅ Pronta | Script testado e funcional |
| **Migração Local** | ✅ Pronta | Docker-compose configurado |
| **Documentação** | ✅ Completa | Especificações técnicas |
| **Otimizações** | ✅ Implementadas | Índices e constraints |
| **Segurança** | ✅ Configurada | RLS e policies |

**🎉 Banco de dados 100% pronto para produção!**
