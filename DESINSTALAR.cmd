@echo off
rem  Quita los accesos directos y el registro del laboratorio.
rem  NO borra el trabajo: bitacora, hallazgos, laminas y codigo quedan.
title Desinstalar el Laboratorio Diosyunalma
cd /d "%~dp0"

if not exist "instalador\desinstalar.ps1" (
  echo.
  echo   [!] Falta instalador\desinstalar.ps1
  echo.
  pause
  exit /b 1
)

powershell -NoProfile -ExecutionPolicy Bypass -File "instalador\desinstalar.ps1"
