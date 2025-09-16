@echo off
setlocal

echo 🚀 Iniciando Sistema de Delivery de Cupcakes...
echo.

REM Verificar se Go está instalado
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Go não está instalado. Por favor, instale Go primeiro.
    echo    Download: https://golang.org/dl/
    pause
    exit /b 1
)

REM Verificar se Node.js está instalado
node --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Node.js não está instalado. Por favor, instale Node.js primeiro.
    echo    Download: https://nodejs.org/
    pause
    exit /b 1
)

echo ✅ Verificações iniciais concluídas!
echo.

REM Instalar dependências do backend
echo 📦 Instalando dependências do backend...
cd backend
go mod tidy
if %errorlevel% neq 0 (
    echo ❌ Erro ao instalar dependências do backend
    pause
    exit /b 1
)

REM Instalar dependências do frontend
echo 📦 Instalando dependências do frontend...
cd ..\frontend
call npm install
if %errorlevel% neq 0 (
    echo ❌ Erro ao instalar dependências do frontend
    pause
    exit /b 1
)

echo ✅ Dependências instaladas com sucesso!
echo.

REM Voltar para o diretório raiz
cd ..

echo 🎯 Para executar o projeto:
echo.
echo 1. Backend (API):
echo    cd backend && go run cmd/api/main.go
echo    Será executado em: http://localhost:8080
echo.
echo 2. Frontend (React):
echo    cd frontend && npm run dev
echo    Será executado em: http://localhost:5173
echo.
echo 📝 Notas importantes:
echo    - Execute o backend primeiro (terminal 1)
echo    - Depois execute o frontend (terminal 2)
echo    - Certifique-se que as portas 8080 e 5173 estão livres
echo.

pause
