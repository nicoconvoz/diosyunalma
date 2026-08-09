# ===================================================================
#  DESINSTALADOR DEL LABORATORIO DIOSYUNALMA
#
#  Quita los accesos directos y el registro. NO borra ni una linea del
#  trabajo: docs/, ckpt/, luz/, galeria/ y cmd/ quedan intactos. El
#  registro del laboratorio no se desinstala (Ley del Registro).
# ===================================================================

$ErrorActionPreference = 'Continue'
$raiz = Split-Path -Parent $PSScriptRoot
$bin  = Join-Path $raiz 'bin'

function Bien($t)  { Write-Host "      OK  $t" -ForegroundColor Green }
function Aviso($t) { Write-Host "      -   $t" -ForegroundColor DarkGray }

Clear-Host
Write-Host ""
Write-Host "  ==============================================================" -ForegroundColor DarkCyan
Write-Host "    DESINSTALAR EL LABORATORIO DIOSYUNALMA" -ForegroundColor Cyan
Write-Host "  ==============================================================" -ForegroundColor DarkCyan
Write-Host ""
Write-Host "  Esto quita los accesos directos y el registro de Windows." -ForegroundColor Gray
Write-Host "  NO se borra nada del trabajo: la bitacora, los hallazgos, las" -ForegroundColor Gray
Write-Host "  laminas, los checkpoints y el codigo quedan donde estan." -ForegroundColor Gray
Write-Host ""
$r = Read-Host "  Sigo? (S/N)"
if ($r -notmatch '^[sSyY]') { Write-Host "  Cancelado." -ForegroundColor Yellow; Start-Sleep 2; exit 0 }

# ---------- accesos directos ----------
Write-Host ""
Write-Host "  Quitando accesos directos" -ForegroundColor Yellow
$escritorio = Join-Path ([Environment]::GetFolderPath('Desktop')) 'El Puente de Mando.lnk'
if (Test-Path $escritorio) { Remove-Item $escritorio -Force; Bien "escritorio" } else { Aviso "no habia acceso en el escritorio" }

$menu = Join-Path ([Environment]::GetFolderPath('Programs')) 'Diosyunalma'
if (Test-Path $menu) { Remove-Item $menu -Recurse -Force; Bien "menu inicio" } else { Aviso "no habia carpeta en el menu inicio" }

# ---------- registro ----------
Write-Host ""
Write-Host "  Quitando el registro" -ForegroundColor Yellow
$clave = 'HKCU:\Software\Microsoft\Windows\CurrentVersion\Uninstall\Diosyunalma'
if (Test-Path $clave) { Remove-Item $clave -Recurse -Force; Bien "Agregar o quitar programas" } else { Aviso "no estaba registrado" }

# ---------- los ejecutables, solo si el capitan quiere ----------
Write-Host ""
if (Test-Path $bin) {
  $n = (Get-ChildItem -Path $bin -Filter *.exe -ErrorAction SilentlyContinue).Count
  Write-Host "  Quedan $n ejecutables en bin\ (se rehacen con INSTALAR.cmd)." -ForegroundColor Gray
  $r2 = Read-Host "  Los borro tambien? (S/N)"
  if ($r2 -match '^[sSyY]') {
    Get-ChildItem -Path $bin -Filter *.exe | Remove-Item -Force
    Bien "$n ejecutables borrados"
  } else {
    Aviso "los ejecutables se quedan"
  }
}

Write-Host ""
Write-Host "  ==============================================================" -ForegroundColor DarkCyan
Write-Host "    Desinstalado. El trabajo sigue entero en:" -ForegroundColor Green
Write-Host "    $raiz" -ForegroundColor Gray
Write-Host "  ==============================================================" -ForegroundColor DarkCyan
Write-Host ""
Write-Host "    Nada de lo escrito se perdio. Para volver a instalar," -ForegroundColor Gray
Write-Host "    doble clic en INSTALAR.cmd." -ForegroundColor Gray
Write-Host ""
Read-Host "  Enter para cerrar"
