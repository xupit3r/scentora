@echo off
REM Scentora Application Launcher for Windows
REM Manages the full stack: CouchDB, Backend API, and Frontend

setlocal EnableDelayedExpansion

echo.
echo ========================================
echo   Scentora Application Manager
echo ========================================
echo.

REM Check for Node.js
where node >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Node.js is not installed
    echo Please install Node.js from https://nodejs.org
    pause
    exit /b 1
)

REM Check for npm
where npm >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] npm is not installed
    pause
    exit /b 1
)

REM Check for Docker
where docker >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo [WARNING] Docker is not installed
    echo CouchDB may not start automatically
    echo.
)

set SCRIPT_DIR=%~dp0
set BACKEND_DIR=%SCRIPT_DIR%backend
set FRONTEND_DIR=%SCRIPT_DIR%frontend

REM Parse command
set CMD=%1
if "%CMD%"=="" set CMD=start

if /I "%CMD%"=="start" goto START
if /I "%CMD%"=="stop" goto STOP
if /I "%CMD%"=="status" goto STATUS
if /I "%CMD%"=="help" goto HELP

echo [ERROR] Unknown command: %CMD%
echo.
goto HELP

:START
echo [INFO] Starting all services...
echo.

REM Start CouchDB
echo [INFO] Starting CouchDB...
cd "%SCRIPT_DIR%"
start /B docker-compose up >nul 2>&1
timeout /t 3 >nul

REM Check if CouchDB is running
curl -s http://localhost:5984 >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    echo [OK] CouchDB is ready
) else (
    echo [WARNING] CouchDB may not be running
)
echo.

REM Install backend dependencies if needed
if not exist "%BACKEND_DIR%\node_modules" (
    echo [INFO] Installing backend dependencies...
    cd "%BACKEND_DIR%"
    call npm install
    echo.
)

REM Install frontend dependencies if needed
if not exist "%FRONTEND_DIR%\node_modules" (
    echo [INFO] Installing frontend dependencies...
    cd "%FRONTEND_DIR%"
    call npm install
    echo.
)

REM Start Backend
echo [INFO] Starting Backend API...
cd "%BACKEND_DIR%"
start "Scentora Backend" cmd /c "npm run dev"
timeout /t 5 >nul

REM Start Frontend
echo [INFO] Starting Frontend...
cd "%FRONTEND_DIR%"
start "Scentora Frontend" cmd /c "npm run dev"
timeout /t 5 >nul

echo.
echo ========================================
echo   All services started!
echo ========================================
echo.
echo URLs:
echo   Frontend:  http://localhost:5173
echo   Backend:   http://localhost:3000
echo   CouchDB:   http://localhost:5984/_utils
echo.
echo To stop services, run: scentora.bat stop
echo Or close the terminal windows manually
echo.
pause
goto END

:STOP
echo [INFO] Stopping all services...
echo.
echo Please close the following windows:
echo   - Scentora Backend
echo   - Scentora Frontend
echo.
echo To stop CouchDB:
cd "%SCRIPT_DIR%"
docker-compose down
echo.
echo [OK] Services stopped
pause
goto END

:STATUS
echo [INFO] Checking service status...
echo.

curl -s http://localhost:5984 >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    echo [RUNNING] CouchDB - http://localhost:5984
) else (
    echo [STOPPED] CouchDB
)

curl -s http://localhost:3000/api/health >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    echo [RUNNING] Backend - http://localhost:3000
) else (
    echo [STOPPED] Backend
)

curl -s http://localhost:5173 >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    echo [RUNNING] Frontend - http://localhost:5173
) else (
    echo [STOPPED] Frontend
)

echo.
pause
goto END

:HELP
echo Usage: scentora.bat [COMMAND]
echo.
echo Commands:
echo   start       Start all services (default)
echo   stop        Stop all services
echo   status      Show service status
echo   help        Show this help message
echo.
echo Examples:
echo   scentora.bat              Start all services
echo   scentora.bat start        Start all services
echo   scentora.bat status       Check status
echo   scentora.bat stop         Stop all services
echo.
pause
goto END

:END
endlocal
