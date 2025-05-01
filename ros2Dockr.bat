@echo off
setlocal ENABLEDELAYEDEXPANSION

:: Required files
set FILES=compose.yaml Dockerfile .env

echo 📂 Checking for required files...
for %%F in (%FILES%) do (
    if exist %%F (
        echo ✅ Found: %%F
    ) else (
        echo ❌ Missing: %%F
        exit /b 1
    )
)

:: Display usage help
if "%1"=="" goto help

:: Command dispatcher
if "%1"=="-start" (
    echo 🚀 Starting Docker Compose (no rebuild)...
    docker compose up -d
    goto end
)

if "%1"=="-clean" (
    echo ♻️ Rebuilding and cleaning volumes...
    docker rm -f ros2 2>nul
    docker compose down
    docker compose up --build -d
    goto end
)

if "%1"=="-stop" (
    echo 🛑 Stopping Docker services...
    docker compose down
    goto end
)

if "%1"=="-shell" (
    echo 🔧 Entering container shell...
    docker exec -it ros2 bash
    goto end
)

if "%1"=="-logs" (
    echo 📜 Streaming logs...
    docker compose logs -f
    goto end
)

:help
echo.
echo ℹ️ Usage: ros2dockr.bat [-clean | -start | -stop | -shell | -logs]

:end
pause
