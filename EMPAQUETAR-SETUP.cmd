@echo off
rem  Arma un setup.exe con todo el laboratorio adentro.
rem  Requiere Go e Inno Setup EN ESTA maquina; la de destino no
rem  necesita nada. Tarda varios minutos: son 200 ejecutables.
title Empaquetar el Laboratorio en un setup.exe
cd /d "%~dp0"

if not exist "instalador\empaquetar.ps1" (
  echo.
  echo   [!] Falta instalador\empaquetar.ps1
  echo.
  pause
  exit /b 1
)

powershell -NoProfile -ExecutionPolicy Bypass -File "instalador\empaquetar.ps1"
if errorlevel 1 (
  echo.
  echo   El empaquetado termino con un problema. El motivo esta arriba.
  echo.
  pause
)
