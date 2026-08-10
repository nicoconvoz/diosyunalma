; ===================================================================
;  LABORATORIO DIOSYUNALMA - guion de Inno Setup
;
;  Arma un setup.exe autocontenido: la maquina de destino no necesita
;  Go ni nada instalado. Se compila con:
;
;      EMPAQUETAR-SETUP.cmd     (desde la raiz del laboratorio)
;
;  UNA DECISION DE INGENIERIA, DECLARADA: la instalacion es POR USUARIO
;  y NO en Archivos de programa. Motivo real, no comodidad: el
;  laboratorio ESCRIBE dentro de su propia carpeta (luz/, ckpt/, las
;  laminas nuevas de galeria/, la bitacora). En Archivos de programa
;  esas escrituras fallarian sin permisos de administrador y el barco
;  quedaria mudo. Por eso PrivilegesRequired=lowest, sin excepcion.
;
;  Y LA LEY DEL REGISTRO TAMBIEN RIGE ACA: ckpt/ y luz/ se copian solo
;  si no existen (onlyifdoesntexist), asi una reinstalacion JAMAS pisa
;  el trabajo vivo. Y el desinstalador borra solo lo que instalo: todo
;  lo que el laboratorio escriba despues queda intacto.
; ===================================================================

#define Nombre    "Laboratorio Diosyunalma"
#define Version   GetDateTimeString('yyyy.mm.dd', '', '')
#define Autor     "el capitan y el Doc"
#define ExePuente "puente.exe"

[Setup]
AppId={{7D3A1C64-9B2E-4F58-A0D7-DIOSYUNALMA01}
AppName={#Nombre}
AppVersion={#Version}
AppVerName={#Nombre} {#Version}
AppPublisher={#Autor}
AppComments=Un laboratorio numerico completo de la funcion zeta, en Go puro.
VersionInfoDescription=Puente de mando del laboratorio Diosyunalma
VersionInfoProductName={#Nombre}
VersionInfoVersion=1.1.0.0

; por usuario: el laboratorio escribe en su propia carpeta
PrivilegesRequired=lowest
DefaultDirName={autopf}\Diosyunalma
DefaultGroupName=Diosyunalma
DisableProgramGroupPage=yes
DisableDirPage=no
AllowNoIcons=yes

OutputDir=salida
OutputBaseFilename=DiosyunalmaSetup
SetupIconFile=diosyunalma.ico
UninstallDisplayIcon={app}\instalador\diosyunalma.ico
UninstallDisplayName={#Nombre}

Compression=lzma2/max
SolidCompression=yes
ArchitecturesInstallIn64BitMode=x64compatible
WizardStyle=modern
ShowLanguageDialog=no
LicenseFile=..\LICENSE
InfoBeforeFile=LEEME-ANTES.txt
InfoAfterFile=LEEME-DESPUES.txt

[Languages]
Name: "es"; MessagesFile: "compiler:Languages\Spanish.isl"

[Tasks]
Name: "escritorio"; Description: "Crear un acceso directo en el escritorio"; GroupDescription: "Accesos directos:"

[Files]
; --- los ejecutables: el laboratorio entero, ya compilado ---
Source: "..\bin\*.exe";        DestDir: "{app}\bin";     Flags: ignoreversion

; --- la obra: laminas, documentos y codigo fuente ---
Source: "..\galeria\*";        DestDir: "{app}\galeria"; Flags: ignoreversion recursesubdirs createallsubdirs
Source: "..\docs\*";           DestDir: "{app}\docs";    Flags: ignoreversion recursesubdirs createallsubdirs
Source: "..\cmd\*";            DestDir: "{app}\cmd";     Flags: ignoreversion recursesubdirs createallsubdirs
Source: "..\go.mod";           DestDir: "{app}";         Flags: ignoreversion
Source: "..\README.md";        DestDir: "{app}";         Flags: ignoreversion
; las licencias viajan siempre: la AGPL las exige junto al programa
Source: "..\LICENSE";                DestDir: "{app}"; Flags: ignoreversion
Source: "..\NOTICE";                 DestDir: "{app}"; Flags: ignoreversion
Source: "..\LICENSE-CONTENIDO.txt";  DestDir: "{app}"; Flags: ignoreversion
Source: "..\LICENCIA-COMERCIAL.md";  DestDir: "{app}"; Flags: ignoreversion
Source: "..\LICENCIAS.md";           DestDir: "{app}"; Flags: ignoreversion
Source: "..\PUENTE.cmd";       DestDir: "{app}";         Flags: ignoreversion
Source: "diosyunalma.ico";     DestDir: "{app}\instalador"; Flags: ignoreversion

; --- tablas y datos chicos que algunos experimentos leen ---
Source: "..\control\*";        DestDir: "{app}\control";     Flags: ignoreversion recursesubdirs createallsubdirs skipifsourcedoesntexist
Source: "..\information\*";    DestDir: "{app}\information"; Flags: ignoreversion recursesubdirs createallsubdirs skipifsourcedoesntexist
Source: "..\pattern\*";        DestDir: "{app}\pattern";     Flags: ignoreversion recursesubdirs createallsubdirs skipifsourcedoesntexist
Source: "..\primes\*";         DestDir: "{app}\primes";      Flags: ignoreversion recursesubdirs createallsubdirs skipifsourcedoesntexist
Source: "..\riemann\*";        DestDir: "{app}\riemann";     Flags: ignoreversion recursesubdirs createallsubdirs skipifsourcedoesntexist
Source: "..\spectral\*";       DestDir: "{app}\spectral";    Flags: ignoreversion recursesubdirs createallsubdirs skipifsourcedoesntexist

; --- el trabajo vivo: SOLO si no existe, jamas pisar ---
Source: "..\ckpt\*";           DestDir: "{app}\ckpt"; Flags: onlyifdoesntexist recursesubdirs createallsubdirs skipifsourcedoesntexist
Source: "..\luz\*";            DestDir: "{app}\luz";  Flags: onlyifdoesntexist recursesubdirs createallsubdirs skipifsourcedoesntexist

[Dirs]
; que existan aunque vengan vacias: el barco escribe ahi
Name: "{app}\ckpt"
Name: "{app}\luz"

[Icons]
Name: "{group}\El Puente de Mando";     Filename: "{app}\bin\{#ExePuente}"; WorkingDir: "{app}"; IconFilename: "{app}\instalador\diosyunalma.ico"; Comment: "Un solo timon para todo el laboratorio"
Name: "{group}\El Faro del Almirante";  Filename: "{app}\bin\faro.exe";     WorkingDir: "{app}"; IconFilename: "{app}\instalador\diosyunalma.ico"; Comment: "Tablero en vivo de la flota"
Name: "{group}\La Galeria";             Filename: "{app}\galeria\index.html"; WorkingDir: "{app}"; IconFilename: "{app}\instalador\diosyunalma.ico"; Comment: "Las laminas del laboratorio"
Name: "{group}\{cm:UninstallProgram,{#Nombre}}"; Filename: "{uninstallexe}"
Name: "{autodesktop}\El Puente de Mando"; Filename: "{app}\bin\{#ExePuente}"; WorkingDir: "{app}"; IconFilename: "{app}\instalador\diosyunalma.ico"; Tasks: escritorio

[Run]
Filename: "{app}\bin\{#ExePuente}"; WorkingDir: "{app}"; Description: "Abrir el puente de mando ahora"; Flags: nowait postinstall skipifsilent

[Code]
// Al desinstalar, avisar en criollo que el trabajo NO se toca.
procedure CurUninstallStepChanged(CurStep: TUninstallStep);
begin
  if CurStep = usUninstall then
    MsgBox('Se van a quitar los programas y los accesos directos.' + #13#10 + #13#10 +
           'El trabajo NO se borra: la bitacora, los hallazgos, las laminas' + #13#10 +
           'nuevas y los checkpoints que hayas generado quedan donde estan.' + #13#10 + #13#10 +
           'La Ley del Registro tambien rige aca.', mbInformation, MB_OK);
end;
