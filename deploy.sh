#!/bin/bash

# Script de Deploy para Produção
# Sistema de Delivery de Cupcakes

set -e

echo "🚀 Iniciando deploy do Sistema de Delivery de Cupcakes..."

# Verificar se Docker está rodando
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker não está rodando. Por favor, inicie o Docker primeiro."
    exit 1
fi

# Verificar se docker-compose está instalado
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose não encontrado. Por favor, instale o docker-compose."
    exit 1
fi

# Criar arquivo .env se não existir
if [ ! -f .env ]; then
    echo "📝 Criando arquivo .env a partir do exemplo..."
    cp .env.example .env
    echo "⚠️  IMPORTANTE: Edite o arquivo .env com suas configurações antes de continuar!"
    echo "   Especialmente: DB_PASSWORD e JWT_SECRET"
    read -p "Pressione Enter após editar o arquivo .env..."
fi

# Parar containers existentes
echo "🛑 Parando containers existentes..."
docker-compose down

# Construir e iniciar containers
echo "🔨 Construindo e iniciando containers..."
docker-compose up -d --build

# Aguardar serviços ficarem prontos
echo "⏳ Aguardando serviços ficarem prontos..."

# Aguardar PostgreSQL
echo "📊 Aguardando PostgreSQL..."
until docker-compose exec postgres pg_isready -U cupcake_user; do
    echo "PostgreSQL ainda não está pronto..."
    sleep 2
done

# Aguardar Backend
echo "🔧 Aguardando Backend API..."
until curl -f http://localhost:8080/health 2>/dev/null; do
    echo "Backend ainda não está pronto..."
    sleep 2
done

# Aguardar Frontend
echo "🌐 Aguardando Frontend..."
until curl -f http://localhost:80 2>/dev/null; do
    echo "Frontend ainda não está pronto..."
    sleep 2
done

echo "✅ Deploy concluído com sucesso!"
echo ""
echo "🌐 Aplicação disponível em:"
echo "   Frontend: http://localhost"
echo "   Backend API: http://localhost:8080"
echo "   PostgreSQL: localhost:5432"
echo ""
echo "📊 Status dos serviços:"
docker-compose ps

echo ""
echo "📝 Logs em tempo real:"
echo "   docker-compose logs -f"
echo ""
echo "🛑 Para parar a aplicação:"
echo "   docker-compose down"
