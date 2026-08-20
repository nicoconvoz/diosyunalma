@echo off
title Empaquetar el Laboratorio - Diosyunalma
cd /d "%~dp0"
echo.
echo   ================================================
echo     EMPAQUETAR EL LABORATORIO
echo     deja una carpeta lista para cualquier Windows,
echo     sin necesidad de instalar Go ni nada mas
echo   ================================================
echo.

where go >nul 2>&1
if errorlevel 1 (
  echo   [!] Para empaquetar hace falta Go en ESTA maquina
  echo       ^(la maquina de destino no necesita nada^).
  echo       Instalalo desde https://go.dev/dl/
  echo.
  pause
  exit /b 1
)

if not exist bin mkdir bin

echo   [1/3] compilando el puente...
go build -ldflags="-s -w" -o bin\puente.exe .\cmd\puente
if errorlevel 1 goto fallo

echo   [2/3] compilando los experimentos ^(tarda unos minutos^)...
for /d %%d in (cmd\*) do (
  go build -ldflags="-s -w" -o "bin\%%~nxd.exe" ".\cmd\%%~nxd" 2>nul
  if errorlevel 1 echo         - no compilo: %%~nxd
)

echo   [3/3] armando la carpeta de reparto...
if exist reparto rmdir /s /q reparto
mkdir reparto
copy /y PUENTE.cmd reparto\ >nul
copy /y README.md reparto\ >nul
copy /y go.mod reparto\ >nul
copy /y docs\guias\LEEME-PAQUETE.txt reparto\LEEME.txt >nul
robocopy bin     reparto\bin     /E /NFL /NDL /NJH /NJS /nc /ns /np >nul
robocopy cmd     reparto\cmd     /E /NFL /NDL /NJH /NJS /nc /ns /np >nul
robocopy galeria reparto\galeria /E /NFL /NDL /NJH /NJS /nc /ns /np >nul
robocopy docs    reparto\docs    /E /NFL /NDL /NJH /NJS /nc /ns /np >nul
robocopy luz     reparto\luz     /E /NFL /NDL /NJH /NJS /nc /ns /np >nul
if errorlevel 8 goto fallo

echo.
echo   LISTO. La carpeta "reparto" se copia entera a cualquier Windows
echo   y arranca con doble clic en PUENTE.cmd - sin instalar nada.
echo.
dir reparto\bin\*.exe /b | find /c ".exe" > "%TEMP%\_n.txt"
set /p NEXE=<"%TEMP%\_n.txt"
del "%TEMP%\_n.txt"
echo   experimentos compilados en el paquete: %NEXE%
echo.
pause
exit /b 0

:fallo
echo.
echo   [!] Algo fallo al empaquetar. El motivo esta arriba.
echo.
pause
exit /b 1
