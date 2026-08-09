@echo off
title El Puente de Mando - Laboratorio Diosyunalma
cd /d "%~dp0"
echo.
echo   ================================================
echo     EL PUENTE DE MANDO - Laboratorio Diosyunalma
echo     un solo timon para todo el laboratorio
echo   ================================================
echo.

rem  Si el puente ya viene compilado, arranca directo: esta maquina
rem  no necesita tener Go ni nada instalado.
if exist "bin\puente.exe" goto zarpar

echo   preparando el laboratorio por primera vez...
where go >nul 2>&1
if errorlevel 1 (
  echo.
  echo   [!] Falta el ejecutable bin\puente.exe y no hay Go instalado.
  echo       Pedile a quien te paso el laboratorio la carpeta bin completa,
  echo       o instala Go desde https://go.dev/dl/
  echo.
  pause
  exit /b 1
)
if not exist bin mkdir bin
go build -ldflags="-s -w" -o bin\puente.exe .\cmd\puente
if errorlevel 1 (
  echo.
  echo   [!] No pude preparar el puente. El motivo esta en el mensaje de arriba.
  echo.
  pause
  exit /b 1
)
echo   listo.
echo.

:zarpar
bin\puente.exe

echo.
echo   Puente amarrado. Pulsa una tecla para cerrar esta ventana.
pause > nul
