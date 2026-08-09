# ===================================================================
#  EMPAQUETAR EL LABORATORIO EN UN setup.exe
#
#  Compila todo el laboratorio y lo mete en un unico instalador de
#  Windows con Inno Setup. La maquina de destino no necesita nada.
#
#  Requiere Go (para compilar) e Inno Setup (para empaquetar) EN ESTA
#  maquina. La de destino no necesita ninguno de los dos.
# ===================================================================

$ErrorActionPreference = 'Stop'
$raiz = Split-Path -Parent $PSScriptRoot
$bin  = Join-Path $raiz 'bin'
$iss  = Join-Path $PSScriptRoot 'diosyunalma.iss'
$ico  = Join-Path $PSScriptRoot 'diosyunalma.ico'
$out  = Join-Path $PSScriptRoot 'salida'

function Titulo($t) { Write-Host ""; Write-Host "  $t" -ForegroundColor Yellow }
function Bien($t)   { Write-Host "      OK  $t" -ForegroundColor Green }
function Aviso($t)  { Write-Host "      !   $t" -ForegroundColor DarkYellow }
function Mal($t)    { Write-Host "      X   $t" -ForegroundColor Red }

Clear-Host
Write-Host ""
Write-Host "  ==============================================================" -ForegroundColor DarkCyan
Write-Host "    EMPAQUETAR EL LABORATORIO EN UN setup.exe" -ForegroundColor Cyan
Write-Host "    un solo archivo, y la maquina de destino no necesita nada" -ForegroundColor DarkGray
Write-Host "  ==============================================================" -ForegroundColor DarkCyan

# ---------- 1. las herramientas ----------
Titulo "[1/5] Buscando las herramientas"

if ($null -eq (Get-Command go -ErrorAction SilentlyContinue)) {
  Mal "Falta Go en esta maquina. Se necesita para compilar el laboratorio."
  Write-Host "      https://go.dev/dl/" -ForegroundColor Gray
  Write-Host ""; Read-Host "  Enter para cerrar"; exit 1
}
Bien "Go: $((& go version) 2>$null)"

$iscc = $null
foreach ($c in @(
  "$env:ProgramFiles\Inno Setup 7\ISCC.exe",
  "${env:ProgramFiles(x86)}\Inno Setup 6\ISCC.exe",
  "$env:ProgramFiles\Inno Setup 6\ISCC.exe",
  "${env:ProgramFiles(x86)}\Inno Setup 5\ISCC.exe")) {
  if (Test-Path $c) { $iscc = $c; break }
}
if (-not $iscc) {
  $g = Get-Command ISCC.exe -ErrorAction SilentlyContinue
  if ($g) { $iscc = $g.Source }
}
if (-not $iscc) {
  Mal "No encuentro ISCC.exe (el compilador de Inno Setup)."
  Write-Host "      https://jrsoftware.org/isdl.php" -ForegroundColor Gray
  Write-Host ""; Read-Host "  Enter para cerrar"; exit 1
}
Bien "Inno Setup: $iscc"

# ---------- 2. compilar el laboratorio ----------
Titulo "[2/5] Compilando el laboratorio entero"
if (-not (Test-Path $bin)) { New-Item -ItemType Directory -Path $bin | Out-Null }
$compilados = 0; $fallados = @()
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
  Bien "$compilados de $($dirs.Count) experimentos"
  if ($fallados.Count -gt 0) { Aviso "no compilaron: $($fallados -join ', ')" }
} finally { Pop-Location }

# ---------- 3. el icono ----------
Titulo "[3/5] Dibujando el sello"
Push-Location $raiz
try { & go run ./cmd/icono $ico | Out-Null } catch { } finally { Pop-Location }
if (Test-Path $ico) { Bien "icono listo" } else { Mal "sin icono"; exit 1 }

# ---------- 4. el peso ----------
Titulo "[4/5] Pesando la carga"
$mb = [math]::Round((Get-ChildItem $bin -Filter *.exe | Measure-Object Length -Sum).Sum / 1MB, 1)
Bien "$((Get-ChildItem $bin -Filter *.exe).Count) ejecutables · $mb MB sin comprimir"
Write-Host "      La compresion LZMA2 al maximo tarda varios minutos." -ForegroundColor DarkGray
Write-Host "      Es una sola vez, y baja el paquete a menos de la mitad." -ForegroundColor DarkGray

# ---------- 5. empaquetar ----------
Titulo "[5/5] Armando el setup.exe"
if (-not (Test-Path $out)) { New-Item -ItemType Directory -Path $out | Out-Null }
$reloj = [Diagnostics.Stopwatch]::StartNew()
& $iscc /Q $iss
if ($LASTEXITCODE -ne 0) {
  Mal "Inno Setup fallo. El motivo esta arriba."
  Write-Host ""; Read-Host "  Enter para cerrar"; exit 1
}
$reloj.Stop()

$paq = Get-ChildItem $out -Filter *.exe | Sort-Object LastWriteTime -Descending | Select-Object -First 1
$pmb = [math]::Round($paq.Length / 1MB, 1)
Bien "$($paq.Name) · $pmb MB · $([int]$reloj.Elapsed.TotalSeconds)s"

Write-Host ""
Write-Host "  ==============================================================" -ForegroundColor DarkCyan
Write-Host "    LISTO. El instalador esta armado." -ForegroundColor Green
Write-Host "  ==============================================================" -ForegroundColor DarkCyan
Write-Host ""
Write-Host "    $($paq.FullName)" -ForegroundColor Gray
Write-Host ""
Write-Host "    Ese solo archivo instala el laboratorio completo en" -ForegroundColor Gray
Write-Host "    cualquier Windows. No hace falta Go ni nada mas del otro lado." -ForegroundColor Gray
Write-Host ""
if ($fallados.Count -gt 0) {
  Write-Host "    Honestidad: $($fallados.Count) experimento(s) quedaron afuera por no compilar." -ForegroundColor DarkYellow
  Write-Host "    El puente los marca 'sin compilar'. Nada se oculta." -ForegroundColor DarkYellow
  Write-Host ""
}
$r = Read-Host "  Abro la carpeta del instalador? (S/N)"
if ($r -match '^[sSyY]') { Start-Process explorer.exe $out }
