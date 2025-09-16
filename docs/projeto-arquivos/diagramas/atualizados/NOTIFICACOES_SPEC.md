# Sistema de Notificações - Cupcake Delivery

## Visão Geral
O sistema de notificações foi implementado para manter todos os usuários (clientes, entregadores e administradores) informados sobre mudanças de status dos pedidos em tempo real.

## Arquitetura

### Backend (Go)

#### Modelo de Dados
```go
type Notification struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    UserID    uint      `gorm:"not null" json:"user_id"`
    OrderID   *uint     `json:"order_id,omitempty"`
    Type      string    `gorm:"not null" json:"type"`
    Title     string    `gorm:"not null" json:"title"`
    Message   string    `gorm:"not null" json:"message"`
    IsRead    bool      `gorm:"default:false" json:"is_read"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

#### API Endpoints
- `GET /notifications` - Buscar notificações do usuário
- `GET /notifications/unread-count` - Contador de não lidas
- `PUT /notifications/:id/read` - Marcar como lida
- `PUT /notifications/mark-all-read` - Marcar todas como lidas
- `POST /notifications/test` - Criar notificação de teste (dev)

#### Tipos de Notificação
- `order_confirmed` - Pedido confirmado
- `order_preparing` - Pedido em preparação
- `order_ready` - Pedido pronto
- `order_delivering` - Pedido saiu para entrega
- `order_delivered` - Pedido entregue
- `order_cancelled` - Pedido cancelado
- `test` - Notificação de teste

#### Triggers Automáticos
As notificações são disparadas automaticamente quando:
1. Status do pedido é atualizado via `PUT /orders/:id/status`
2. `NotificationService.NotifyOrderStatusChange()` é chamado
3. Mensagens personalizadas são geradas para cada tipo de usuário

### Frontend (React)

#### Componentes
- `NotificationDropdown` - Dropdown com lista de notificações
- `BellIcon` - Ícone do sino com indicador de notificações
- `Toast` - Sistema de toast para feedback

#### Hooks
- `useNotifications` - Gerencia estado das notificações
  - Carregamento automático
  - Polling a cada 30 segundos
  - Sincronização com backend
  - Tratamento de erros

#### Serviços
- `NotificationService` - Abstração para APIs de notificação
  - Autenticação automática
  - Tratamento de erros
  - TypeScript typing

## Fluxos de Uso

### 1. Notificação Automática
1. Admin atualiza status do pedido
2. Backend identifica mudança de status
3. Sistema gera notificações para usuários relevantes
4. Frontend recebe via polling ou ao abrir dropdown

### 2. Visualização de Notificações
1. Usuário vê badge com contador no ícone do sino
2. Clica para abrir dropdown
3. Visualiza lista de notificações ordenadas por data
4. Pode marcar individualmente ou todas como lidas

### 3. Indicadores Visuais
- Badge laranja no sino (não lidas)
- Ponto pulsante na aba "Ativos" 
- Destaque visual em notificações não lidas
- Animações suaves de entrada/saída

## Regras de Negócio

### Destinatários por Status
- **Confirmed**: Cliente + Admin
- **Preparing**: Cliente
- **Ready**: Cliente + Todos os Entregadores
- **Delivering**: Cliente
- **Delivered**: Cliente + Admin
- **Cancelled**: Cliente + Admin

### Controle de Acesso
- Usuários só veem suas próprias notificações
- Autenticação obrigatória via JWT
- Filtragem automática por `user_id`

### Performance
- Polling limitado a 30 segundos
- Limite de 20 notificações por consulta
- Cache local no frontend
- Verificação de autenticação antes de requisições

## Melhorias Futuras
- [ ] WebSockets para notificações em tempo real
- [ ] Push notifications no mobile
- [ ] Persistência de preferências de notificação
- [ ] Templates customizáveis de mensagens
- [ ] Analytics de engajamento
- [ ] Notificações por email
- [ ] Categorização avançada de notificações

## Testes
- Endpoint de teste: `POST /notifications/test`
- Permite criar notificações para desenvolvimento
- Requer autenticação
- Aceita title, message e type customizados

## Monitoramento
- Logs de criação de notificações
- Métricas de taxa de leitura
- Contadores por tipo de notificação
- Performance de consultas no banco

## Database Schema
```sql
CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    order_id INTEGER REFERENCES orders(id),
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);
```
