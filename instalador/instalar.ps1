# ===================================================================
#  INSTALADOR DEL LABORATORIO DIOSYUNALMA
#  Deja el puente de mando listo para usar con un doble clic.
#
#  No pide permisos de administrador: todo se escribe en el perfil del
#  usuario. Y NUNCA toca docs/, ckpt/, luz/ ni galeria/ — el registro
#  del laboratorio es intocable (Ley del Registro).
# ===================================================================

$ErrorActionPreference = 'Stop'
$raiz = Split-Path -Parent $PSScriptRoot
$bin  = Join-Path $raiz 'bin'
$ico  = Join-Path $PSScriptRoot 'diosyunalma.ico'

function Titulo($t) { Write-Host ""; Write-Host "  $t" -ForegroundColor Yellow }
function Bien($t)   { Write-Host "      OK  $t" -ForegroundColor Green }
function Aviso($t)  { Write-Host "      !   $t" -ForegroundColor DarkYellow }
function Mal($t)    { Write-Host "      X   $t" -ForegroundColor Red }

Clear-Host
Write-Host ""
Write-Host "  ==============================================================" -ForegroundColor DarkCyan
Write-Host "    LABORATORIO DIOSYUNALMA - instalador del puente de mando" -ForegroundColor Cyan
Write-Host "    un solo timon para todos los experimentos" -ForegroundColor DarkGray
Write-Host "  ==============================================================" -ForegroundColor DarkCyan
Write-Host ""
Write-Host "  carpeta del laboratorio: $raiz" -ForegroundColor DarkGray

if (-not (Test-Path (Join-Path $raiz 'go.mod'))) {
  Mal "Esta carpeta no parece el laboratorio (falta go.mod)."
  Write-Host ""; Read-Host "  Enter para cerrar"; exit 1
}

# ---------- 1. el motor ----------
Titulo "[1/5] Buscando el motor de compilacion"
$hayGo = $null -ne (Get-Command go -ErrorAction SilentlyContinue)
$hayPuente = Test-Path (Join-Path $bin 'puente.exe')
$fallados = @()

if ($hayGo) {
  $v = (& go version) 2>$null
  Bien "Go encontrado: $v"
} elseif ($hayPuente) {
  Aviso "No hay Go, pero el puente ya viene compilado: sigo sin recompilar."
} else {
  Mal "No hay Go instalado y tampoco hay bin\puente.exe."
  Write-Host ""
  Write-Host "      Dos salidas:" -ForegroundColor Gray
  Write-Host "        a) instalar Go desde https://go.dev/dl/  y volver a correr esto" -ForegroundColor Gray
  Write-Host "        b) pedir la carpeta bin ya compilada a quien te paso el laboratorio" -ForegroundColor Gray
  Write-Host ""; Read-Host "  Enter para cerrar"; exit 1
}

# ---------- 2. compilar ----------
Titulo "[2/5] Compilando el laboratorio"
$compilados = 0
if ($hayGo) {
  if (-not (Test-Path $bin)) { New-Item -ItemType Directory -Path $bin | Out-Null }
  Push-Location $raiz
  try {
    & go build -ldflags="-s -w" -o (Join-Path $bin 'puente.exe') ./cmd/puente
    if ($LASTEXITCODE -ne 0) { throw "no pude compilar el puente" }
    Bien "el puente de mando"

    $dirs = Get-ChildItem -Path (Join-Path $raiz 'cmd') -Directory | Sort-Object Name
    $i = 0
    foreach ($d in $dirs) {
      $i++
      Write-Progress -Activity "Compilando los experimentos" -Status "$($d.Name)   ($i de $($dirs.Count))" -PercentComplete ([int]($i * 100 / $dirs.Count))
      & go build -ldflags="-s -w" -o (Join-Path $bin "$($d.Name).exe") "./cmd/$($d.Name)" 2>$null
      if ($LASTEXITCODE -eq 0) { $compilados++ } else { $fallados += $d.Name }
    }
    Write-Progress -Activity "Compilando los experimentos" -Completed
    Bien "$compilados de $($dirs.Count) experimentos compilados"
    if ($fallados.Count -gt 0) {
      Aviso "$($fallados.Count) no compilaron: $($fallados -join ', ')"
      Aviso "el puente igual los puede correr con Go al vuelo"
    }
  } finally { Pop-Location }
} else {
  $compilados = (Get-ChildItem -Path $bin -Filter *.exe -ErrorAction SilentlyContinue).Count
  Bien "usando los $compilados ejecutables que ya estaban"
}

# ---------- 3. el icono ----------
Titulo "[3/5] Preparando el icono"
if ($hayGo) {
  Push-Location $raiz
  try { & go run ./cmd/icono $ico | Out-Null } catch { } finally { Pop-Location }
}
if (Test-Path $ico) { Bien "icono listo" } else { Aviso "sin icono propio; Windows usara el suyo" }

# ---------- 4. accesos directos ----------
Titulo "[4/5] Creando los accesos directos"
$ws = New-Object -ComObject WScript.Shell

function Atajo($ruta, $destino, $desc) {
  $sc = $ws.CreateShortcut($ruta)
  $sc.TargetPath = $destino
  $sc.WorkingDirectory = $raiz
  $sc.Description = $desc
  if (Test-Path $ico) { $sc.IconLocation = "$ico,0" }
  $sc.Save()
}

$escritorio = [Environment]::GetFolderPath('Desktop')
$menu = Join-Path ([Environment]::GetFolderPath('Programs')) 'Diosyunalma'
if (-not (Test-Path $menu)) { New-Item -ItemType Directory -Path $menu | Out-Null }

$puenteExe = Join-Path $bin 'puente.exe'
Atajo (Join-Path $escritorio 'El Puente de Mando.lnk') $puenteExe 'El puente de mando del laboratorio Diosyunalma'
Bien "escritorio: El Puente de Mando"

Atajo (Join-Path $menu 'El Puente de Mando.lnk') $puenteExe 'El puente de mando'
$faroExe = Join-Path $bin 'faro.exe'
if (Test-Path $faroExe) { Atajo (Join-Path $menu 'El Faro del Almirante.lnk') $faroExe 'Tablero en vivo de la flota' }
$galeria = Join-Path $raiz 'galeria\index.html'
if (Test-Path $galeria) { Atajo (Join-Path $menu 'La Galeria.lnk') $galeria 'Las laminas del laboratorio' }
Atajo (Join-Path $menu 'Desinstalar Diosyunalma.lnk') (Join-Path $raiz 'DESINSTALAR.cmd') 'Quitar accesos directos y registro'
Bien "menu inicio: carpeta Diosyunalma"

# ---------- 5. registro ----------
Titulo "[5/5] Registrando en Programas y caracteristicas"
$clave = 'HKCU:\Software\Microsoft\Windows\CurrentVersion\Uninstall\Diosyunalma'
if (-not (Test-Path $clave)) { New-Item -Path $clave -Force | Out-Null }
Set-ItemProperty $clave DisplayName     'Laboratorio Diosyunalma - Puente de Mando'
Set-ItemProperty $clave DisplayVersion  (Get-Date -Format 'yyyy.MM.dd')
Set-ItemProperty $clave Publisher       'el capitan y el Doc'
Set-ItemProperty $clave InstallLocation $raiz
Set-ItemProperty $clave UninstallString ('"' + (Join-Path $raiz 'DESINSTALAR.cmd') + '"')
if (Test-Path $ico) { Set-ItemProperty $clave DisplayIcon $ico }
Set-ItemProperty $clave NoModify 1 -Type DWord
Set-ItemProperty $clave NoRepair 1 -Type DWord
Bien "aparece en Agregar o quitar programas"

# ---------- final ----------
Write-Host ""
Write-Host "  ==============================================================" -ForegroundColor DarkCyan
Write-Host "    LISTO. El laboratorio quedo instalado." -ForegroundColor Green
Write-Host "  ==============================================================" -ForegroundColor DarkCyan
Write-Host ""
Write-Host "    Doble clic en 'El Puente de Mando' del escritorio y zarpa." -ForegroundColor Gray
Write-Host "    El puente abre solo el navegador en http://localhost:8118" -ForegroundColor DarkGray
Write-Host ""
Write-Host "    Arriba de todo esta EL PUESTO DE MANDO con los tres vehiculos:" -ForegroundColor Gray
Write-Host "    el DeLorean, el Tren del Doc Brown y el Faro." -ForegroundColor Gray
Write-Host ""
if ($fallados.Count -gt 0) {
  Write-Host "    Honestidad: $($fallados.Count) experimento(s) no compilaron." -ForegroundColor DarkYellow
  Write-Host "    Nada se oculta - el puente los marca 'sin compilar'." -ForegroundColor DarkYellow
  Write-Host ""
}
$r = Read-Host "  Abro el puente ahora? (S/N)"
if ($r -match '^[sSyY]') { Start-Process -FilePath $puenteExe -WorkingDirectory $raiz }
