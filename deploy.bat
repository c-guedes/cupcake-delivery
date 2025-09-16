@echo off
setlocal

REM Script de Deploy para Windows
REM Sistema de Delivery de Cupcakes

echo 🚀 Iniciando deploy do Sistema de Delivery de Cupcakes...

REM Verificar se Docker está rodando
docker info >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Docker não está rodando. Por favor, inicie o Docker primeiro.
    exit /b 1
)

REM Verificar se docker-compose está instalado
docker-compose --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ docker-compose não encontrado. Por favor, instale o docker-compose.
    exit /b 1
)

REM Criar arquivo .env se não existir
if not exist .env (
    echo 📝 Criando arquivo .env a partir do exemplo...
    copy .env.example .env
    echo ⚠️  IMPORTANTE: Edite o arquivo .env com suas configurações antes de continuar!
    echo    Especialmente: DB_PASSWORD e JWT_SECRET
    pause
)

REM Parar containers existentes
echo 🛑 Parando containers existentes...
docker-compose down

REM Construir e iniciar containers
echo 🔨 Construindo e iniciando containers...
docker-compose up -d --build

REM Aguardar serviços ficarem prontos
echo ⏳ Aguardando serviços ficarem prontos...

REM Aguardar um tempo para os serviços iniciarem
timeout /t 10 /nobreak >nul

echo ✅ Deploy concluído com sucesso!
echo.
echo 🌐 Aplicação disponível em:
echo    Frontend: http://localhost
echo    Backend API: http://localhost:8080
echo    PostgreSQL: localhost:5432
echo.
echo 📊 Status dos serviços:
docker-compose ps

echo.
echo 📝 Logs em tempo real:
echo    docker-compose logs -f
echo.
echo 🛑 Para parar a aplicação:
echo    docker-compose down

pause
