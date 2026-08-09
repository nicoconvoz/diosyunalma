@echo off
rem  Instalador del Laboratorio Diosyunalma.
rem  Doble clic y listo: compila todo, crea los accesos directos y
rem  registra el laboratorio. No pide permisos de administrador.
title Instalar el Laboratorio Diosyunalma
cd /d "%~dp0"

if not exist "instalador\instalar.ps1" (
  echo.
  echo   [!] Falta instalador\instalar.ps1
  echo       La carpeta del laboratorio esta incompleta.
  echo.
  pause
  exit /b 1
)

powershell -NoProfile -ExecutionPolicy Bypass -File "instalador\instalar.ps1"
if errorlevel 1 (
  echo.
  echo   El instalador termino con un problema. El motivo esta arriba.
  echo.
  pause
)
