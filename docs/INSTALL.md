# 🚀 Guia de Instalação e Deploy

Este guia detalha como instalar, configurar e fazer deploy do Sistema de Delivery de Cupcakes.

## 📋 Pré-requisitos

### Desenvolvimento Local
- **Go** 1.19 ou superior
- **Node.js** 18 ou superior  
- **PostgreSQL** 13 ou superior
- **Git**

### Produção
- **Servidor Linux** (Ubuntu 20.04+ recomendado)
- **Docker** e **Docker Compose** (opcional)
- **Nginx** (para proxy reverso)
- **SSL Certificate** (Let's Encrypt recomendado)

## 🔧 Instalação Local

### 1. Clone do Repositório
```bash
git clone https://github.com/usuario/cupcake-delivery.git
cd cupcake-delivery
```

### 2. Configuração do Banco de Dados

#### Instalação PostgreSQL (Ubuntu/Debian)
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

#### Criação do Banco
```bash
sudo -u postgres psql
```

```sql
CREATE DATABASE cupcake_delivery;
CREATE USER cupcake_user WITH PASSWORD 'senha_segura';
GRANT ALL PRIVILEGES ON DATABASE cupcake_delivery TO cupcake_user;
\q
```

### 3. Configuração do Backend

#### Instalar Go (se necessário)
```bash
# Ubuntu/Debian
sudo apt install golang-go

# Ou baixar de https://golang.org/dl/
```

#### Configurar Backend
```bash
cd backend

# Criar arquivo de configuração
cp .env.example .env

# Editar variáveis de ambiente
nano .env
```

#### Arquivo .env
```bash
# Banco de dados
DATABASE_URL="postgres://cupcake_user:senha_segura@localhost:5432/cupcake_delivery?sslmode=disable"

# JWT
JWT_SECRET="seu_secret_super_seguro_aqui_128_bits_minimo"
JWT_EXPIRY="24h"

# Servidor
PORT="8080"
GIN_MODE="development"  # ou "release" para produção

# CORS
FRONTEND_URL="http://localhost:5173"

# Upload de arquivos (opcional)
UPLOAD_PATH="./uploads"
MAX_FILE_SIZE="10MB"
```

#### Instalar Dependências e Executar
```bash
# Baixar dependências
go mod download

# Executar testes
go test ./...

# Executar aplicação
go run cmd/api/main.go
```

### 4. Configuração do Frontend

#### Instalar Node.js (se necessário)
```bash
# Ubuntu/Debian via NodeSource
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# Verificar instalação
node --version
npm --version
```

#### Configurar Frontend
```bash
cd frontend

# Instalar dependências
npm install

# Executar testes
npm test

# Executar em desenvolvimento
npm run dev
```

#### Configuração de Ambiente (opcional)
```bash
# Criar .env.local se necessário
echo "VITE_API_URL=http://localhost:8080" > .env.local
```

## 🐳 Deploy com Docker

### 1. Docker Compose
```yaml
# docker-compose.yml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: cupcake_delivery
      POSTGRES_USER: cupcake_user
      POSTGRES_PASSWORD: senha_segura
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backend/scripts/init.sql:/docker-entrypoint-initdb.d/init.sql
    ports:
      - "5432:5432"
    restart: unless-stopped

  backend:
    build: 
      context: ./backend
      dockerfile: Dockerfile
    environment:
      DATABASE_URL: "postgres://cupcake_user:senha_segura@postgres:5432/cupcake_delivery?sslmode=disable"
      JWT_SECRET: "seu_secret_super_seguro_aqui"
      GIN_MODE: "release"
      PORT: "8080"
    ports:
      - "8080:8080"
    depends_on:
      - postgres
    restart: unless-stopped

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    ports:
      - "3000:80"
    restart: unless-stopped

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./nginx/ssl:/etc/nginx/ssl
    depends_on:
      - backend
      - frontend
    restart: unless-stopped

volumes:
  postgres_data:
```

### 2. Dockerfile Backend
```dockerfile
# backend/Dockerfile
FROM golang:1.19-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080
CMD ["./main"]
```

### 3. Dockerfile Frontend
```dockerfile
# frontend/Dockerfile
FROM node:18-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

### 4. Executar com Docker
```bash
# Build e start
docker-compose up -d --build

# Ver logs
docker-compose logs -f

# Parar
docker-compose down
```

## 🌐 Deploy em Produção

### 1. Configuração do Servidor

#### Atualizar Sistema
```bash
sudo apt update && sudo apt upgrade -y
sudo apt install -y curl wget git nginx certbot python3-certbot-nginx
```

#### Instalar Docker
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
sudo systemctl enable docker
```

### 2. Configuração Nginx

#### /etc/nginx/sites-available/cupcake-delivery
```nginx
server {
    listen 80;
    server_name seu-dominio.com www.seu-dominio.com;

    # Redirect HTTP to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name seu-dominio.com www.seu-dominio.com;

    # SSL Configuration
    ssl_certificate /etc/letsencrypt/live/seu-dominio.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/seu-dominio.com/privkey.pem;
    
    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # Frontend
    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }

    # API Backend
    location /api/ {
        proxy_pass http://localhost:8080/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
```

#### Ativar Site
```bash
sudo ln -s /etc/nginx/sites-available/cupcake-delivery /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 3. SSL com Let's Encrypt
```bash
sudo certbot --nginx -d seu-dominio.com -d www.seu-dominio.com
```

### 4. Deploy da Aplicação

#### Usando Git (recomendado)
```bash
# Clone no servidor
git clone https://github.com/usuario/cupcake-delivery.git /opt/cupcake-delivery
cd /opt/cupcake-delivery

# Configurar ambiente de produção
cp backend/.env.example backend/.env
# Editar backend/.env com dados de produção

# Build e deploy
docker-compose -f docker-compose.prod.yml up -d --build
```

### 5. Monitoramento e Logs

#### Systemd Service (alternativa ao Docker)
```ini
# /etc/systemd/system/cupcake-backend.service
[Unit]
Description=Cupcake Delivery Backend
After=network.target postgresql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/cupcake-delivery/backend
ExecStart=/opt/cupcake-delivery/backend/main
Restart=always
RestartSec=10

Environment=GIN_MODE=release
EnvironmentFile=/opt/cupcake-delivery/backend/.env

[Install]
WantedBy=multi-user.target
```

#### Ativar Serviço
```bash
sudo systemctl daemon-reload
sudo systemctl enable cupcake-backend
sudo systemctl start cupcake-backend
sudo systemctl status cupcake-backend
```

### 6. Backup e Manutenção

#### Script de Backup
```bash
#!/bin/bash
# backup.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/opt/backups"
DB_NAME="cupcake_delivery"

# Criar diretório se não existe
mkdir -p $BACKUP_DIR

# Backup do banco
pg_dump $DB_NAME > $BACKUP_DIR/db_backup_$DATE.sql

# Backup dos uploads (se existir)
tar -czf $BACKUP_DIR/uploads_backup_$DATE.tar.gz /opt/cupcake-delivery/uploads

# Manter apenas backups dos últimos 7 dias
find $BACKUP_DIR -name "*.sql" -mtime +7 -delete
find $BACKUP_DIR -name "*.tar.gz" -mtime +7 -delete
```

#### Cron para Backup Automático
```bash
# Editar crontab
sudo crontab -e

# Adicionar linha para backup diário às 2h
0 2 * * * /opt/scripts/backup.sh
```

## 🔍 Troubleshooting

### Problemas Comuns

#### Backend não conecta ao banco
```bash
# Verificar se PostgreSQL está rodando
sudo systemctl status postgresql

# Verificar logs
sudo journalctl -u postgresql

# Testar conexão
psql -h localhost -U cupcake_user -d cupcake_delivery
```

#### Frontend não carrega
```bash
# Verificar se Node.js está instalado
node --version

# Limpar cache e reinstalar
rm -rf node_modules package-lock.json
npm install
```

#### Erro de CORS
- Verificar `FRONTEND_URL` no .env do backend
- Confirmar que as portas estão corretas

#### SSL não funciona
```bash
# Verificar certificados
sudo certbot certificates

# Renovar se necessário
sudo certbot renew --dry-run
```

### Logs Úteis

```bash
# Backend logs
docker-compose logs backend

# Frontend logs  
docker-compose logs frontend

# Nginx logs
sudo tail -f /var/log/nginx/error.log
sudo tail -f /var/log/nginx/access.log

# PostgreSQL logs
sudo tail -f /var/log/postgresql/postgresql-15-main.log
```

## 📊 Monitoramento

### Health Check Endpoints

#### Backend
```bash
curl http://localhost:8080/health
```

#### Banco de Dados
```bash
curl http://localhost:8080/health/db
```

### Ferramentas Recomendadas
- **Logs**: ELK Stack, Grafana Loki
- **Métricas**: Prometheus + Grafana
- **Uptime**: UptimeRobot, Pingdom
- **Errors**: Sentry

## 🔐 Segurança

### Checklist de Segurança
- [ ] JWT Secret com 256+ bits
- [ ] HTTPS ativado
- [ ] Firewall configurado (UFW)
- [ ] Atualizações automáticas
- [ ] Backup regular
- [ ] Rate limiting ativo
- [ ] Headers de segurança
- [ ] Logs de auditoria

### Hardening do Servidor
```bash
# Firewall básico
sudo ufw allow ssh
sudo ufw allow 80
sudo ufw allow 443
sudo ufw enable

# Fail2ban para SSH
sudo apt install fail2ban
sudo systemctl enable fail2ban

# Atualizações automáticas
sudo apt install unattended-upgrades
sudo dpkg-reconfigure unattended-upgrades
```

---

Para suporte adicional, consulte a documentação específica de cada componente ou entre em contato com a equipe de desenvolvimento.
