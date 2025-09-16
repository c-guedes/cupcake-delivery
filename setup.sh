#!/bin/bash

echo "🚀 Iniciando Sistema de Delivery de Cupcakes..."
echo

# Verificar se Go está instalado
if ! command -v go &> /dev/null; then
    echo "❌ Go não está instalado. Por favor, instale Go primeiro."
    echo "   Download: https://golang.org/dl/"
    exit 1
fi

# Verificar se Node.js está instalado
if ! command -v node &> /dev/null; then
    echo "❌ Node.js não está instalado. Por favor, instale Node.js primeiro."
    echo "   Download: https://nodejs.org/"
    exit 1
fi

echo "✅ Verificações iniciais concluídas!"
echo

# Instalar dependências do backend
echo "📦 Instalando dependências do backend..."
cd backend
go mod tidy
if [ $? -ne 0 ]; then
    echo "❌ Erro ao instalar dependências do backend"
    exit 1
fi

# Instalar dependências do frontend
echo "📦 Instalando dependências do frontend..."
cd ../frontend
npm install
if [ $? -ne 0 ]; then
    echo "❌ Erro ao instalar dependências do frontend"
    exit 1
fi

echo "✅ Dependências instaladas com sucesso!"
echo

# Voltar para o diretório raiz
cd ..

echo "🎯 Para executar o projeto:"
echo
echo "1. Backend (API):"
echo "   cd backend && go run cmd/api/main.go"
echo "   Será executado em: http://localhost:8080"
echo
echo "2. Frontend (React):"
echo "   cd frontend && npm run dev"
echo "   Será executado em: http://localhost:5173"
echo
echo "📝 Notas importantes:"
echo "   - Execute o backend primeiro (terminal 1)"
echo "   - Depois execute o frontend (terminal 2)"
echo "   - Certifique-se que as portas 8080 e 5173 estão livres"
echo

read -p "Pressione Enter para continuar..."
