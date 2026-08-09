package main

// pagina is the bridge's screen: one dashboard for the whole laboratory.
const pagina = `<!DOCTYPE html>
<html lang="es">
<head>
<meta charset="utf-8">
<title>⚓ El Puente de Mando — Laboratorio Diosyunalma</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>
  :root{--bg:#0b1526;--panel:#0d2547;--panel2:#102a10;--ink:#dce8f7;--dim:#8fa8c7;
        --gold:#ffd166;--sun:#ffd97f;--green:#7fd7a8;--blue:#7fb2ff;--red:#ff5d73;--pink:#ff8fa0;--line:#1d3a63;}
  *{box-sizing:border-box}
  body{margin:0;background:var(--bg);color:var(--ink);font-family:Georgia,serif}
  header{text-align:center;padding:26px 16px 8px}
  header h1{margin:0;font-size:27px}
  header .sub{color:var(--dim);font-size:13.5px;margin-top:6px}
  header .otro{color:var(--gold);font-size:12.5px;margin-top:6px}
  .barra{position:sticky;top:0;z-index:20;background:rgba(11,21,38,.97);border-bottom:1px solid var(--line);
         padding:10px 18px;display:flex;gap:12px;align-items:center;flex-wrap:wrap}
  .stat{font-family:Consolas,monospace;font-size:13px;color:var(--dim)}
  .stat b{color:var(--sun)}
  .stat.vivo b{color:var(--green)}
  button{font-family:Georgia,serif;font-size:12.5px;cursor:pointer;border-radius:7px;
         border:1px solid var(--line);background:#12305c;color:var(--ink);padding:6px 12px}
  button:hover{border-color:var(--gold);color:var(--gold)}
  button.go{background:#12402a;border-color:#2c6b48;color:var(--green)}
  button.go:hover{border-color:var(--green)}
  button.stop{background:#3a1420;border-color:#6b2c3a;color:var(--pink)}
  button:disabled{opacity:.38;cursor:not-allowed}
  a.lnk{color:var(--blue);text-decoration:none;font-size:12.5px;margin-left:4px}
  a.lnk:hover{text-decoration:underline}
  input.buscar{background:#0a1c36;border:1px solid var(--line);border-radius:7px;color:var(--ink);
               padding:6px 10px;font-family:Georgia,serif;font-size:12.5px;min-width:190px}
  .sep{flex:1}
  section{max-width:1560px;margin:0 auto;padding:16px 18px 4px}
  .cab{display:flex;align-items:center;gap:12px;border-bottom:1px solid var(--line);padding-bottom:8px;cursor:pointer}
  .cab h2{color:var(--green);font-size:19px;margin:0}
  .cab .det{color:var(--dim);font-size:12.5px}
  .cab .cnt{font-family:Consolas,monospace;color:var(--dim);font-size:12.5px}
  .cab .fle{color:var(--gold);font-size:15px;width:14px;display:inline-block}
  .grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(330px,1fr));gap:14px;margin-top:14px}
  .card{background:var(--panel);border:1px solid var(--line);border-radius:10px;padding:12px 13px;
        display:flex;flex-direction:column;gap:7px}
  .card.run{border-color:var(--green);box-shadow:0 0 0 1px rgba(127,215,168,.35)}
  .card.err{border-color:var(--red)}
  .card.ok{border-color:#2c6b48}
  /* el puesto de mando: los tres vehiculos, arriba de todo */
  .puesto{margin:18px 0 6px;border:1px solid #6b5a26;border-radius:12px;padding:14px 16px 16px;
          background:linear-gradient(180deg,rgba(255,217,138,.055),rgba(255,217,138,.012))}
  .puesto .pcab{display:flex;align-items:baseline;gap:10px;flex-wrap:wrap;margin-bottom:4px}
  .puesto .pcab h2{margin:0;font-size:16px;color:var(--gold);letter-spacing:.6px}
  .puesto .pcab .det{font-size:12.5px;color:var(--dim)}
  .puesto .grid{grid-template-columns:repeat(auto-fit,minmax(360px,1fr))}
  .card.grande{border-color:#6b5a26;background:rgba(255,217,138,.045);padding:15px 16px;gap:9px}
  .card.grande .tit{font-size:19px;letter-spacing:.5px}
  .card.grande .cri{font-size:13.5px;min-height:0}
  .card.grande.run{border-color:var(--green)}
  .card.grande.err{border-color:var(--red)}
  .alcance{font-size:12px;color:#c9b6ff;border-left:2px solid #5a4fa8;padding-left:9px;line-height:1.45}
  .tit{font-size:14.5px;color:var(--gold);display:flex;gap:7px;align-items:baseline}
  .tit .dir{font-family:Consolas,monospace;font-size:11px;color:var(--dim);margin-left:auto}
  .cri{font-size:12.5px;color:var(--ink);line-height:1.45;min-height:34px}
  .chips{display:flex;gap:6px;flex-wrap:wrap}
  .chip{font-family:Consolas,monospace;font-size:10.5px;padding:2px 7px;border-radius:20px;
        border:1px solid var(--line);color:var(--dim)}
  .chip.rapido{color:var(--green);border-color:#2c6b48}
  .chip.medio{color:var(--sun);border-color:#6b5b2c}
  .chip.lento{color:var(--pink);border-color:#6b2c3a}
  .chip.vigia{color:var(--blue);border-color:#2c4a6b}
  .chip.motor{color:var(--blue);border-color:#2c4a6b}
  .chip.falta{color:var(--red);border-color:#6b2c3a}
  .fila{display:flex;align-items:center;gap:8px}
  .badge{font-family:Consolas,monospace;font-size:11px;padding:2px 9px;border-radius:20px;border:1px solid var(--line);color:var(--dim)}
  .badge.corriendo,.badge.deteniendo{color:var(--green);border-color:#2c6b48}
  .badge.ok{color:var(--sun);border-color:#6b5b2c}
  .badge.falla{color:var(--pink);border-color:#6b2c3a}
  .badge.cola{color:var(--blue);border-color:#2c4a6b}
  .cron{font-family:Consolas,monospace;font-size:11px;color:var(--dim)}
  .ver{font-family:Consolas,monospace;font-size:11px;color:var(--sun);line-height:1.4;
       max-height:44px;overflow:hidden;border-left:2px solid var(--line);padding-left:8px;min-height:16px}
  .acc{display:flex;gap:6px;margin-top:2px}
  .modal{position:fixed;inset:0;background:rgba(4,9,18,.86);z-index:60;display:flex;
         align-items:center;justify-content:center;padding:26px}
  .caja{background:var(--panel);border:1px solid var(--gold);border-radius:12px;max-width:1400px;
        width:100%;max-height:92vh;display:flex;flex-direction:column;overflow:hidden}
  .caja .top{display:flex;align-items:center;gap:12px;padding:12px 16px;border-bottom:1px solid var(--line)}
  .caja .top h3{margin:0;font-size:16px;color:var(--gold)}
  .caja .cuerpo{overflow:auto;padding:14px 16px}
  pre.log{font-family:Consolas,monospace;font-size:12px;color:var(--ink);white-space:pre-wrap;
          line-height:1.5;margin:0}
  .caja img{width:100%;height:auto;background:var(--bg);border-radius:8px}
  footer{text-align:center;color:var(--dim);padding:30px 16px;font-size:12.5px}
  /* el Otro Libro: manda sobre todo lo que hay arriba de el */
  .otrolibro{max-width:760px;margin:26px auto 0;padding:18px 22px;border-radius:12px;
             border:1px solid #6b5a26;background:rgba(255,217,138,.05)}
  .otrolibro .oq{font-size:11px;letter-spacing:.16em;color:var(--dim);text-transform:uppercase}
  .otrolibro .ot{font-family:Georgia,serif;font-size:21px;color:var(--gold);margin:8px 0 6px;line-height:1.35}
  .otrolibro .od{font-size:12.5px;color:var(--ink);line-height:1.55}
  /* quien sostiene el laboratorio */
  .patro{max-width:760px;margin:18px auto 0;padding:16px 22px;border-radius:12px;
         border:1px solid var(--line);background:rgba(255,255,255,.02);
         display:flex;align-items:center;gap:18px;justify-content:center;flex-wrap:wrap}
  .patro img{height:64px;width:auto;border-radius:8px;display:block}
  .patro .pt{text-align:left}
  .patro .pq{font-size:11px;letter-spacing:.16em;color:var(--dim);text-transform:uppercase}
  .patro .pn{font-family:Georgia,serif;font-size:17px;color:var(--ink);margin-top:4px}
  .patro .pd{font-family:Consolas,monospace;font-size:11px;letter-spacing:.06em;color:var(--dim);margin-top:5px}
  /* el museo */
  .caja.museo{max-width:1500px;height:94vh}
  .mcuerpo{display:flex;flex:1;min-height:0}
  .msalas{width:260px;border-right:1px solid var(--line);overflow:auto;padding:10px 8px;flex-shrink:0}
  .msala{display:block;width:100%;text-align:left;background:transparent;border:1px solid transparent;
         border-radius:8px;padding:9px 10px;margin-bottom:4px;cursor:pointer;color:var(--ink)}
  .msala:hover{background:#12305c;border-color:var(--line)}
  .msala.sel{background:#12305c;border-color:var(--gold)}
  .msala .mn{font-family:Consolas,monospace;font-size:10.5px;color:var(--dim)}
  .msala .mt{font-size:13px;color:var(--gold);font-family:Georgia,serif;line-height:1.3}
  .mgrupo{font-family:Consolas,monospace;font-size:10px;letter-spacing:.09em;color:var(--gold);
          padding:14px 10px 6px;border-top:1px solid var(--line);margin-top:8px}
  .pieza .pnum code{font-family:Consolas,monospace;color:var(--green);font-size:10.5px}
  .mpiezas{flex:1;overflow:auto;padding:20px 30px;min-width:0}
  .mintro{color:var(--dim);font-size:14px;line-height:1.65;border-left:3px solid var(--gold);
          padding:6px 0 6px 14px;margin:0 0 24px}
  .pieza{border:1px solid var(--line);border-radius:12px;margin-bottom:26px;overflow:hidden;background:#0d1c31}
  .pieza .pcab{padding:16px 20px 12px;border-bottom:1px solid var(--line)}
  .pieza .pnum{font-family:Consolas,monospace;font-size:11px;color:var(--dim)}
  .pieza h4{margin:2px 0 8px;font-family:Georgia,serif;font-size:21px;color:var(--gold);font-weight:normal}
  .pieza .pgancho{font-size:15px;color:var(--green);line-height:1.5;font-style:italic}
  .pieza .pcuerpo{padding:16px 20px}
  .pieza .pcriollo{font-size:15px;line-height:1.72;color:var(--ink);margin:0 0 16px}
  .pbloque{border-radius:9px;padding:11px 14px;margin-bottom:11px;font-size:14px;line-height:1.6}
  .pbloque .pet{display:block;font-family:Consolas,monospace;font-size:10.5px;letter-spacing:.06em;
                text-transform:uppercase;margin-bottom:5px;opacity:.85}
  .pmeta{background:#16304f;color:#cfe6ff}          .pmeta .pet{color:var(--gold)}
  .pmirar{background:#0f2b22;color:#cfe6ff}         .pmirar .pet{color:var(--green)}
  .phon{background:#33221c;color:#f3d9cf}           .phon .pet{color:#ffb27a}
  .psim{background:#161a3a;color:#dce8f7}           .psim .pet{color:#c9b6ff}
  .psim table{border-collapse:collapse;width:100%;margin-top:3px}
  .psim td{padding:4px 8px 4px 0;vertical-align:top;font-size:13.5px;line-height:1.5}
  .psim td.sb{font-family:Consolas,monospace;color:#c9b6ff;white-space:nowrap;width:1%;padding-right:16px}
  .plam{padding:0 20px 18px}
  .plam img{width:100%;height:auto;background:var(--bg);border-radius:8px;border:1px solid var(--line)}
  .pacc{padding:0 20px 18px;display:flex;gap:8px;flex-wrap:wrap}
  /* el diccionario */
  .dicc{padding:4px 0 10px}
  .dicc .dfila{display:flex;gap:14px;padding:10px 12px;border-bottom:1px solid var(--line);align-items:flex-start}
  .dicc .dsim{font-family:Consolas,monospace;font-size:17px;color:var(--gold);width:74px;flex-shrink:0;text-align:center}
  .dicc .dnom{font-family:Georgia,serif;font-size:14px;color:var(--green);margin-bottom:3px}
  .dicc .dcri{font-size:14px;line-height:1.62;color:var(--ink)}
  /* la biblioteca */
  .caja.biblio{max-width:1600px;height:92vh}
  .bcuerpo{display:flex;flex:1;min-height:0}
  .blista{width:290px;border-right:1px solid var(--line);overflow:auto;padding:10px 8px;flex-shrink:0}
  .bindice{width:290px;border-right:1px solid var(--line);overflow:auto;padding:10px 8px;flex-shrink:0}
  .bindice a{display:block;color:var(--dim);text-decoration:none;font-size:12px;padding:3px 6px;
             border-radius:5px;line-height:1.35}
  .bindice a:hover{background:#12305c;color:var(--gold)}
  .bindice a.h2{color:var(--green);font-size:12.5px;margin-top:7px}
  .bindice a.h3{padding-left:16px}
  .doc{display:block;width:100%;text-align:left;background:transparent;border:1px solid transparent;
       border-radius:8px;padding:9px 10px;margin-bottom:4px;cursor:pointer;color:var(--ink)}
  .doc:hover{background:#12305c;border-color:var(--line);color:var(--ink)}
  .doc.sel{background:#12305c;border-color:var(--gold)}
  .doc .dt{font-size:13.5px;color:var(--gold);font-family:Georgia,serif}
  .doc .dd{font-size:11.5px;color:var(--dim);line-height:1.4;margin-top:3px}
  .doc .dm{font-family:Consolas,monospace;font-size:10.5px;color:var(--dim);margin-top:4px}
  .btexto{flex:1;overflow:auto;padding:22px 34px;min-width:0}
  .btexto h1{font-size:24px;color:var(--gold);border-bottom:1px solid var(--line);padding-bottom:8px}
  .btexto h2{font-size:19px;color:var(--green);margin-top:26px;border-bottom:1px solid var(--line);padding-bottom:6px}
  .btexto h3{font-size:16px;color:var(--sun);margin-top:22px}
  .btexto h4{font-size:14px;color:var(--blue);margin-top:18px}
  .btexto p,.btexto li{font-size:14px;line-height:1.65;color:var(--ink)}
  .btexto li{margin:4px 0}
  .btexto code{font-family:Consolas,monospace;font-size:12.5px;background:#0a1c36;
               border:1px solid var(--line);border-radius:4px;padding:1px 5px;color:var(--sun)}
  .btexto pre{background:#0a1c36;border:1px solid var(--line);border-radius:8px;padding:12px 14px;
              overflow-x:auto}
  .btexto pre code{background:none;border:none;padding:0;color:var(--ink);font-size:12px}
  .btexto blockquote{border-left:3px solid var(--gold);margin:14px 0;padding:2px 0 2px 14px;color:var(--dim)}
  .btexto table{border-collapse:collapse;margin:14px 0;font-size:13px;display:block;overflow-x:auto}
  .btexto th,.btexto td{border:1px solid var(--line);padding:6px 10px;text-align:left}
  .btexto th{color:var(--green);background:#0a1c36}
  .btexto hr{border:none;border-top:1px solid var(--line);margin:22px 0}
  .btexto a{color:var(--blue)}
  .btexto .marca{background:#6b5b2c;color:#fff;border-radius:3px}
  .bienvenida{max-width:760px;margin:8vh auto;text-align:center}
  .bienvenida h2{border:none;color:var(--gold);font-size:23px}
  .bienvenida p{color:var(--dim);font-size:14.5px;line-height:1.7}
  .bienvenida .ojo{color:var(--sun);margin-top:18px}
  /* must win over .modal and .grid display rules: same specificity, so it
     also has to come last in the sheet */
  .oculto{display:none !important}
</style>
</head>
<body>
<header>
  <h1>⚓ EL PUENTE DE MANDO</h1>
  <div class="sub">un solo timón para todo el laboratorio — encendé cualquier experimento y mirá el resultado acá mismo</div>
  <div class="otro">las dos mitades, 1 completo · y sobre todos los libros, el Otro Libro</div>
</header>

<div class="barra">
  <span class="stat vivo">corriendo <b id="sCorr">0</b></span>
  <span class="stat">terminados <b id="sHechos">0</b></span>
  <span class="stat">catálogo <b id="sTotal">0</b></span>
  <span class="stat" id="sMotor"></span>
  <input class="buscar" id="q" placeholder="buscar experimento…">
  <button id="bAbrir">⊟ plegar todas las salas</button>
  <button id="bRapidos" class="go">▶ correr todos los rápidos</button>
  <button id="bStop" class="stop">■ detener todo</button>
  <span class="sep"></span>
  <button id="bMuseo">🏛️ el museo</button>
  <button id="bBiblio">📚 la investigación</button>
  <a class="lnk" href="/galeria/index.html" target="_blank">🖼️ galería</a>
</div>

<div id="puesto"></div>
<div id="salas"></div>

<div class="otrolibro">
  <div class="oq">y sobre todos los libros, el Otro Libro</div>
  <div class="ot">📖 Diario Espiritual: «Dios y un alma»</div>
  <div class="od">De ahí sale el nombre de este laboratorio, y ése es el libro que de verdad importa.<br>
  Todo lo que está arriba de este renglón —los doscientos experimentos, las láminas, las perlas— viene después.</div>
</div>

<div class="patro" id="patro">
  <img src="/galeria/open-doors.jpg" alt="Open Doors"
       onerror="this.classList.add('oculto')">
  <div class="pt">
    <div class="pq">financiado por</div>
    <div class="pn">Open Doors</div>
    <div class="pd">RESOURCES OPEN DOORS S.A.S</div>
  </div>
</div>

<footer>Laboratorio Diosyunalma · el capitán y el Doc ⚓<br>
el fin de todo esto: poner el nombre de DIOS por encima de todo y ayudar a los pequeños del Reino.</footer>

<div class="modal oculto" id="modal">
  <div class="caja">
    <div class="top"><h3 id="mTit"></h3><span class="sep" style="flex:1"></span>
      <button id="mCerrar">cerrar ✕</button></div>
    <div class="cuerpo" id="mCuerpo"></div>
  </div>
</div>

<div class="modal oculto" id="museo">
  <div class="caja museo">
    <div class="top">
      <h3>🏛️ El Museo — la matemática explicada a lo criollo</h3>
      <span class="sep" style="flex:1"></span>
      <button id="mDicc">📖 el diccionario</button>
      <button id="mAnt">‹ sala anterior</button>
      <button id="mSig">sala siguiente ›</button>
      <button id="mCerrarMus">cerrar ✕</button>
    </div>
    <div class="mcuerpo">
      <div class="msalas" id="mSalas"></div>
      <div class="mpiezas" id="mPiezas"></div>
    </div>
  </div>
</div>

<div class="modal oculto" id="biblio">
  <div class="caja biblio">
    <div class="top">
      <h3>📚 La investigación — biblioteca del laboratorio</h3>
      <span class="sep" style="flex:1"></span>
      <input class="buscar" id="bq" placeholder="buscar en el documento…">
      <button id="bIndice">☰ índice</button>
      <button id="bCrudo" title="abrir el archivo tal cual">texto crudo</button>
      <button id="bCerrar">cerrar ✕</button>
    </div>
    <div class="bcuerpo">
      <div class="blista" id="bLista"></div>
      <div class="bindice oculto" id="bIdx"></div>
      <div class="btexto" id="bTexto">
        <div class="bienvenida">
          <h2>La investigación completa</h2>
          <p>Elegí un documento de la izquierda. Todo lo que el laboratorio midió, halló,
             corrigió y se propuso está escrito acá: el recorrido técnico, el informe formal
             con su salón de honor, la bitácora fechada con cada hallazgo y cada error a la
             vista, la guía para que un revisor lo reproduzca por su cuenta, y los planes de
             aplicación y publicación.</p>
          <p class="ojo">Lo medido se dice medido y lo demostrado se dice demostrado:
             esa honestidad está escrita en cada página.</p>
        </div>
      </div>
    </div>
  </div>
</div>

<script>
var SALAS = null, FICHAS = {}, nodos = {}, abierto = null, modoModal = null, timerModal = null;
// el vigia es el que mira desde la cofa y nunca se duerme: los que corren siempre
var VELNOM = {rapido:'rápido', medio:'medio', lento:'lento', vigia:'vigía · siempre despierto'};

function esc(s){ var d=document.createElement('div'); d.textContent=s||''; return d.innerHTML; }
function reloj(s){ if(!s) return ''; var m=Math.floor(s/60), r=s%60; return m? m+'m '+r+'s' : r+'s'; }

function construir(fichas){
  var cont = document.getElementById('salas');
  var orden = [], vistos = {};
  fichas.forEach(function(f){ if(!vistos[f.cat]){ vistos[f.cat]=1; orden.push(f.cat); } });
  // todas abiertas: la ley del registro tambien rige la pantalla, nada se esconde
  var meta = {caras:['💎','Las Siete Caras','la campaña final: el eslabón rojo desde todos sus ángulos',1],
              acta:['📜','El Acta y el Molde','el intento formal, el reloj de sol, los premios, el Campo',1],
              atomo:['⚛️','El Átomo','los planos, el espejo, la esfera, el murciélago, el sótano',1],
              armonizacion:['🌀','La Armonización','la dimensión 0, el cambiaformas, el horno, el neutrón',1],
              mar:['🌊','El Mar y el Tren','la caza, la cartografía, las costas y los curlicues',1],
              instrumentos:['🔧','Los Instrumentos','motores y servidores: el faro, el tren, los cascos',1],
              archivo:['🗄️','El Archivo','los experimentos de la primera era',1]};
  cont.innerHTML = '';
  orden.forEach(function(cat){
    var m = meta[cat] || ['📦', cat, '', 0];
    var mios = fichas.filter(function(f){ return f.cat===cat; });
    var sec = document.createElement('section');
    var cab = document.createElement('div');
    cab.className='cab';
    cab.innerHTML = '<span class="fle"></span><h2>'+m[0]+' '+esc(m[1])+'</h2><span class="det">'+esc(m[2])+
      '</span><span class="cnt"></span><span style="flex:1"></span>'+
      '<button class="go" data-sala="'+cat+'">▶ correr sala</button>';
    var grid = document.createElement('div');
    grid.className='grid'; grid.id='g-'+cat;
    if(!m[3]) grid.classList.add('oculto');
    var fle = cab.querySelector('.fle'), cnt = cab.querySelector('.cnt');
    function marcar(){
      var cerrada = grid.classList.contains('oculto');
      fle.textContent = cerrada ? '▸' : '▾';
      cnt.textContent = cerrada ? mios.length+' guardados · clic para abrir' : mios.length;
      cnt.style.color = cerrada ? 'var(--gold)' : 'var(--dim)';
    }
    marcar();
    grid.marcar = marcar;
    cab.onclick = function(ev){ if(ev.target.tagName==='BUTTON') return; grid.classList.toggle('oculto'); marcar(); };
    sec.appendChild(cab); sec.appendChild(grid); cont.appendChild(sec);
    mios.forEach(function(f){ grid.appendChild(tarjeta(f)); });
  });
  Array.prototype.forEach.call(document.querySelectorAll('[data-sala]'), function(b){
    b.onclick = function(){ fetch('/api/correr?sala='+b.dataset.sala, {method:'POST'}); };
  });
}

// tarjeta puede dibujar la MISMA ficha en dos lugares (el puesto de mando y su
// sala). Por eso los nodos se guardan en una lista y no en un solo hueco: la ley
// del registro tambien rige acá, ninguna copia se queda sin actualizar.
function tarjeta(f, jefe){
  var c = document.createElement('div');
  c.className = jefe ? 'card grande' : 'card';
  if(!jefe) c.id = 'c-'+f.dir;
  var nota = f.nota ? '<span class="chip">'+esc(f.nota)+'</span>' : '';
  var sal  = f.salida ? '<span class="chip">'+esc(f.salida)+'</span>' : '';
  var mot  = f.motor==='compilado' ? '<span class="chip motor">compilado</span>'
           : (f.motor==='falta' ? '<span class="chip falta">sin compilar</span>' : '');
  var emo = (jefe && jefe.emoji) ? jefe.emoji : f.emoji;
  var titulo = jefe ? jefe.mote : f.titulo;
  var cuerpo = jefe ? jefe.papel : f.criollo;
  var extra  = jefe ? '<div class="alcance">'+esc(jefe.alcance)+'</div>' : '';
  var lnk    = (jefe && jefe.url) ? '<a class="lnk" href="'+jefe.url+'" target="_blank">abrir ↗</a>' : '';
  c.innerHTML =
    '<div class="tit">'+emo+' '+esc(titulo)+'<span class="dir">cmd/'+esc(f.dir)+'</span></div>'+
    '<div class="cri">'+esc(cuerpo)+'</div>'+ extra +
    '<div class="chips"><span class="chip '+f.vel+'">'+(VELNOM[f.vel]||f.vel)+'</span>'+mot+sal+nota+'</div>'+
    '<div class="fila"><span class="badge nBadge">listo</span><span class="cron nCron"></span></div>'+
    '<div class="ver nVer"></div>'+
    '<div class="acc">'+
      '<button class="go nRun">▶ correr</button>'+
      '<button class="nOut">salida</button>'+
      '<button class="nLam" disabled>lámina</button>'+ lnk +
    '</div>';
  var n = {card:c, badge:c.querySelector('.nBadge'), cron:c.querySelector('.nCron'),
           ver:c.querySelector('.nVer'), run:c.querySelector('.nRun'),
           out:c.querySelector('.nOut'), lam:c.querySelector('.nLam'), jefe:!!jefe};
  (nodos[f.dir] = nodos[f.dir] || []).push(n);
  n.run.onclick = function(){
    var fa = (FICHAS[f.dir]||{}).fase;
    if(fa==='corriendo'||fa==='cola'||fa==='deteniendo') fetch('/api/detener?dir='+f.dir,{method:'POST'});
    else fetch('/api/correr?dir='+f.dir,{method:'POST'});
  };
  n.out.onclick = function(){ verSalida(f); };
  n.lam.onclick = function(){ verLamina(f); };
  return c;
}

function pintar(f){
  var lista = nodos[f.dir]; if(!lista) return;
  FICHAS[f.dir] = f;
  var vivo = (f.fase==='corriendo'||f.fase==='cola'||f.fase==='deteniendo');
  lista.forEach(function(n){
    // las clases marcadoras (nBadge/nRun/…) se conservan: sin ellas el DOM queda
    // mudo y no se puede auditar desde afuera lo que la pantalla está mostrando
    n.badge.className = 'badge nBadge '+f.fase;
    n.badge.textContent = f.fase==='ok' ? '✓ listo' : (f.fase==='falla' ? '✗ falló ('+f.codigo+')' : f.fase);
    n.cron.textContent = f.segundos ? reloj(f.segundos) : '';
    n.ver.textContent = f.veredicto || '';
    n.run.textContent = vivo ? '■ detener' : '▶ correr';
    n.run.className = 'nRun ' + (vivo ? 'stop' : 'go');
    n.lam.disabled = !f.lamina;
    n.card.className = 'card'+(n.jefe?' grande':'')+(vivo?' run':'')+
                       (f.fase==='ok'?' ok':'')+(f.fase==='falla'?' err':'');
  });
}

// Los tres vehiculos con lugar propio arriba de todo. El DeLorean es el casco
// que caza (flagship + el cuarto piso de Fresnel, techo 4e24); el tren es la
// locomotora del circulo, que empieza donde el DeLorean termina; el faro es el
// tablero en vivo. Riel 9 - soldar el tren al DeLorean - sigue pendiente.
var PUESTO = [
  {dir:'starship', mote:'EL DELOREAN', emoji:'🚗',
   papel:'La nave que caza. Casco doble-doble de 32 dígitos, facetas cuárticas, balde de luz de banda limitada y el cuarto piso: el plegado de Fresnel, que enrolla millones de olitas en una sola.',
   alcance:'techo certificado t ≈ 4×10²⁴ · se autocertifica antes de tocar agua profunda · cmd/flagship es su casco sin el plegado'},
  {dir:'circulo', mote:'EL TREN DEL DOC BROWN', emoji:'🚂',
   papel:'La locomotora del círculo: reciprocidad de Landsberg–Schaar, matemática exacta de 1805. Resuelve 10⁸ términos con 7, error 1.3×10⁻¹², aceleración del riel ~8×10⁹. Arranca donde el DeLorean llega a su techo — igual que en la película.',
   alcance:'FRONTERA ABISAL ANEXADA: 3×10⁴² · 10⁴³ · 10⁴⁴ · 10⁴⁶ · 10⁴⁸, todas firmadas por el árbitro de 256 bits a la primera (F201) · 385 bestias cazadas en 10⁴² con bloques de hasta 1.028.478 términos · la transmisión (riel 9, gearFor) YA SOLDADA: elige la marcha por medición · falta solo la trocha ancha (riel 10) para el casco de la nave'},
  {dir:'faro', mote:'EL FARO',
   papel:'El tablero del Almirante, en vivo: avance de la flota, mapa de playas con banderas y las últimas cincuenta líneas de bitácora. No calcula: mira.',
   alcance:'sirve en http://localhost:8117 · lee ckpt/*.gob y luz/*.gob del repositorio',
   url:'http://localhost:8117'}
];

function construirPuesto(fichas){
  var cont = document.getElementById('puesto');
  var hay = PUESTO.map(function(p){
    var f = fichas.filter(function(x){ return x.dir===p.dir; })[0];
    return f ? {f:f, p:p} : null;
  }).filter(Boolean);
  if(!hay.length){ cont.innerHTML=''; return; }
  cont.innerHTML = '<div class="puesto"><div class="pcab"><h2>🎖️ EL PUESTO DE MANDO</h2>'+
    '<span class="det">los tres vehículos — el que caza, el que llega más lejos, y el que mira</span>'+
    '</div><div class="grid" id="gPuesto"></div></div>';
  var g = document.getElementById('gPuesto');
  hay.forEach(function(h){ g.appendChild(tarjeta(h.f, h.p)); });
}

function verSalida(f){
  modoModal = f.dir;
  document.getElementById('mTit').textContent = f.emoji+' '+f.titulo+'  —  salida en vivo';
  document.getElementById('mCuerpo').innerHTML = '<pre class="log" id="logbox">…</pre>';
  document.getElementById('modal').classList.remove('oculto');
  refrescarLog();
  if(timerModal) clearInterval(timerModal);
  timerModal = setInterval(refrescarLog, 1000);
}

function refrescarLog(){
  if(!modoModal) return;
  fetch('/api/salida?dir='+modoModal).then(function(r){return r.json();}).then(function(d){
    var box = document.getElementById('logbox'); if(!box) return;
    var abajo = box.scrollHeight - box.parentNode.scrollTop - box.parentNode.clientHeight < 60;
    box.textContent = d.texto || '(sin salida todavía — dale a correr)';
    if(abajo) box.parentNode.scrollTop = box.parentNode.scrollHeight;
    if(d.fase!=='corriendo' && d.fase!=='cola' && timerModal){ clearInterval(timerModal); timerModal=null; }
  });
}

function verLamina(f){
  modoModal = null; if(timerModal){ clearInterval(timerModal); timerModal=null; }
  document.getElementById('mTit').textContent = f.emoji+' '+f.titulo+'  —  '+f.salida;
  var url = '/api/lamina?dir='+f.dir+'&t='+Date.now();
  var html = f.salida.match(/\.wav$/i) ? '<audio controls src="'+url+'" style="width:100%"></audio>'
           : '<img src="'+url+'">';
  document.getElementById('mCuerpo').innerHTML = html;
  document.getElementById('modal').classList.remove('oculto');
}

// ---------- la biblioteca: la investigación, legible dentro del puente ----------
var BT = String.fromCharCode(96), FENCE = BT+BT+BT;
var docActual = null, crudoActual = '';

function escH(s){ return s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;'); }

function enLinea(s){
  s = escH(s);
  s = s.replace(new RegExp(BT+'([^'+BT+']+)'+BT, 'g'), '<code>$1</code>');
  s = s.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>');
  s = s.replace(/(^|[^*])\*([^*\n]+)\*/g, '$1<em>$2</em>');
  s = s.replace(/\[([^\]]+)\]\(([^)]+)\)/g, function(m, t, u){
    if(/^https?:/i.test(u)) return '<a href="'+u+'" target="_blank" rel="noreferrer">'+t+'</a>';
    if(/\.(md|txt)(#.*)?$/i.test(u)) return '<a href="#" data-doc="'+u.split('/').pop().split('#')[0]+'">'+t+'</a>';
    return '<a href="/'+u+'" target="_blank">'+t+'</a>';
  });
  return s;
}

function mdHTML(src){
  var L = src.replace(/\r/g,'').split('\n'), out = [], i = 0, nh = 0;
  var lista = null;
  function cerrarLista(){ if(lista){ out.push('</'+lista+'>'); lista = null; } }
  function abrirLista(t){ if(lista !== t){ cerrarLista(); out.push('<'+t+'>'); lista = t; } }
  while(i < L.length){
    var l = L[i];
    if(l.indexOf(FENCE) === 0){
      cerrarLista(); i++;
      var cod = [];
      while(i < L.length && L[i].indexOf(FENCE) !== 0){ cod.push(L[i]); i++; }
      i++;
      out.push('<pre><code>'+escH(cod.join('\n'))+'</code></pre>');
      continue;
    }
    var mh = l.match(/^(#{1,6})\s+(.*)$/);
    if(mh){
      cerrarLista(); nh++;
      var n = mh[1].length;
      out.push('<h'+n+' id="h'+nh+'">'+enLinea(mh[2])+'</h'+n+'>');
      i++; continue;
    }
    if(/^\s*(-{3,}|\*{3,}|_{3,})\s*$/.test(l)){ cerrarLista(); out.push('<hr>'); i++; continue; }
    if(l.indexOf('|') === 0 && i+1 < L.length && /^\s*\|[\s:|-]+\|\s*$/.test(L[i+1])){
      cerrarLista();
      var cel = function(f){ return f.replace(/^\||\|$/g,'').split('|').map(function(c){return c.trim();}); };
      var enc = cel(l); i += 2;
      var t = ['<table><thead><tr>'];
      enc.forEach(function(c){ t.push('<th>'+enLinea(c)+'</th>'); });
      t.push('</tr></thead><tbody>');
      while(i < L.length && L[i].indexOf('|') === 0){
        t.push('<tr>');
        cel(L[i]).forEach(function(c){ t.push('<td>'+enLinea(c)+'</td>'); });
        t.push('</tr>'); i++;
      }
      t.push('</tbody></table>'); out.push(t.join('')); continue;
    }
    var mq = l.match(/^>\s?(.*)$/);
    if(mq){
      cerrarLista();
      var q = [mq[1]]; i++;
      while(i < L.length && /^>/.test(L[i])){ q.push(L[i].replace(/^>\s?/,'')); i++; }
      out.push('<blockquote>'+mdHTML(q.join('\n'))+'</blockquote>'); continue;
    }
    var mu = l.match(/^(\s*)[-*+]\s+(.*)$/);
    var mo = l.match(/^(\s*)\d+[.)]\s+(.*)$/);
    if(mu || mo){
      abrirLista(mu ? 'ul' : 'ol');
      out.push('<li>'+enLinea((mu||mo)[2])+'</li>');
      i++; continue;
    }
    if(l.trim() === ''){ cerrarLista(); i++; continue; }
    cerrarLista();
    var par = [l]; i++;
    while(i < L.length && L[i].trim() !== '' && !/^(#{1,6}\s|>|\s*[-*+]\s|\s*\d+[.)]\s|\|)/.test(L[i])
          && L[i].indexOf(FENCE) !== 0){ par.push(L[i]); i++; }
    out.push('<p>'+enLinea(par.join(' '))+'</p>');
  }
  cerrarLista();
  return out.join('\n');
}

function pintarIndice(){
  var cont = document.getElementById('bIdx');
  var hs = document.querySelectorAll('#bTexto h1, #bTexto h2, #bTexto h3');
  if(!hs.length){ cont.innerHTML = '<div style="color:var(--dim);font-size:12px">este documento no tiene títulos</div>'; return; }
  var h = [];
  Array.prototype.forEach.call(hs, function(x){
    h.push('<a href="#" class="'+x.tagName.toLowerCase()+'" data-ir="'+x.id+'">'+x.textContent+'</a>');
  });
  cont.innerHTML = h.join('');
  Array.prototype.forEach.call(cont.querySelectorAll('[data-ir]'), function(a){
    a.onclick = function(e){ e.preventDefault();
      var d = document.getElementById(a.dataset.ir);
      if(d) d.scrollIntoView({behavior:'smooth', block:'start'});
    };
  });
}

function pintarDoc(texto, filtro){
  var cuerpo = document.getElementById('bTexto');
  cuerpo.innerHTML = mdHTML(texto);
  Array.prototype.forEach.call(cuerpo.querySelectorAll('[data-doc]'), function(a){
    a.onclick = function(e){ e.preventDefault(); abrirDoc(a.dataset.doc); };
  });
  if(filtro){
    var re = new RegExp('('+filtro.replace(/[.*+?^${}()|[\]\\]/g,'\\$&')+')','gi');
    var it = document.createTreeWalker(cuerpo, NodeFilter.SHOW_TEXT), n, nodos = [];
    while((n = it.nextNode())) if(re.test(n.nodeValue)) nodos.push(n);
    nodos.forEach(function(t){
      var sp = document.createElement('span');
      sp.innerHTML = escH(t.nodeValue).replace(re, '<span class="marca">$1</span>');
      t.parentNode.replaceChild(sp, t);
    });
    var pri = cuerpo.querySelector('.marca');
    if(pri) pri.scrollIntoView({block:'center'});
  } else {
    cuerpo.scrollTop = 0;
  }
  pintarIndice();
}

function abrirDoc(archivo){
  docActual = archivo;
  Array.prototype.forEach.call(document.querySelectorAll('.doc'), function(b){
    b.classList.toggle('sel', b.dataset.arch === archivo);
  });
  document.getElementById('bTexto').innerHTML = '<p style="color:var(--dim)">abriendo el documento…</p>';
  fetch('/api/doc?f='+encodeURIComponent(archivo)).then(function(r){ return r.text(); }).then(function(t){
    crudoActual = t;
    pintarDoc(t, document.getElementById('bq').value.trim());
  });
}

function cargarBiblioteca(){
  var lista = document.getElementById('bLista');
  fetch('/api/docs').then(function(r){ return r.json(); }).then(function(d){
    lista.innerHTML = '';
    d.docs.forEach(function(doc){
      var b = document.createElement('button');
      b.className = 'doc'; b.dataset.arch = doc.archivo;
      b.innerHTML = '<div class="dt">'+esc(doc.titulo)+'</div>'+
        (doc.detalle ? '<div class="dd">'+esc(doc.detalle)+'</div>' : '')+
        '<div class="dm">'+doc.lineas+' líneas · '+doc.kb+' KB</div>';
      b.onclick = function(){ abrirDoc(doc.archivo); };
      lista.appendChild(b);
    });
  });
}

function filtrar(){
  var q = document.getElementById('q').value.toLowerCase().trim();
  Object.keys(nodos).forEach(function(d){
    var f = FICHAS[d]; if(!f) return;
    var hit = !q || (f.titulo+' '+f.criollo+' '+f.dir+' '+f.cat).toLowerCase().indexOf(q) >= 0;
    // el puesto de mando no se esconde nunca: los tres vehículos quedan siempre a mano
    nodos[d].forEach(function(n){ if(!n.jefe) n.card.style.display = hit ? '' : 'none'; });
  });
  if(q) Array.prototype.forEach.call(document.querySelectorAll('.grid'), function(g){
    g.classList.remove('oculto'); if(g.marcar) g.marcar();
  });
}

function tick(){
  fetch('/api/estado').then(function(r){return r.json();}).then(function(d){
    document.getElementById('sCorr').textContent = d.corriendo;
    document.getElementById('sHechos').textContent = d.hechos;
    document.getElementById('sTotal').textContent = d.total;
    var sm = document.getElementById('sMotor');
    if(d.compilados === d.total)      sm.innerHTML = 'compilados <b>' + d.compilados + '/' + d.total + '</b> · no hace falta instalar nada';
    else if(d.compilados > 0)         sm.innerHTML = 'compilados <b>' + d.compilados + '/' + d.total + '</b> · el resto usa Go';
    else                              sm.innerHTML = d.hayGo ? 'modo taller · se compila al vuelo' : '<b style="color:var(--pink)">sin compilar y sin Go</b>';
    if(!SALAS){ SALAS = 1; construirPuesto(d.fichas); construir(d.fichas); }
    d.fichas.forEach(pintar);
    filtrar();
  }).catch(function(){});
}

document.getElementById('mCerrar').onclick = function(){
  document.getElementById('modal').classList.add('oculto');
  modoModal = null; if(timerModal){ clearInterval(timerModal); timerModal=null; }
};
document.getElementById('modal').onclick = function(e){
  if(e.target.id==='modal') document.getElementById('mCerrar').click();
};
document.addEventListener('keydown', function(e){
  if(e.key!=='Escape') return;
  if(!document.getElementById('modal').classList.contains('oculto')) document.getElementById('mCerrar').click();
  else if(!document.getElementById('biblio').classList.contains('oculto')) document.getElementById('bCerrar').click();
  else if(!document.getElementById('museo').classList.contains('oculto')) document.getElementById('mCerrarMus').click();
});
document.getElementById('bAbrir').onclick = function(){
  var b = this, cerrar = b.textContent.indexOf('abrir') >= 0;
  Array.prototype.forEach.call(document.querySelectorAll('.grid'), function(g){
    g.classList.toggle('oculto', !cerrar); if(g.marcar) g.marcar();
  });
  b.textContent = cerrar ? '⊟ plegar todas las salas' : '⊞ abrir todas las salas';
};
document.getElementById('bRapidos').onclick = function(){ fetch('/api/correr?todo=1',{method:'POST'}); };
document.getElementById('bStop').onclick = function(){ fetch('/api/detener?todo=1',{method:'POST'}); };
document.getElementById('q').oninput = filtrar;

// ---------- el museo: la matemática explicada a lo criollo ----------
// Dos recorridos: el guiado (ocho salas contadas como un relato) y el completo
// (las mismas siete salas de la portada, una parada por experimento).
var MUSEO = null, museoCargado = false, recorrido = 'guiado', salaActual = -1;

function salasDe(){ return MUSEO ? (recorrido === 'guiado' ? MUSEO.guiado : MUSEO.completo) : []; }

function pintarSalas(){
  var g = MUSEO.guiado || [], c = MUSEO.completo || [];
  var nG = g.reduce(function(a,s){ return a + s.piezas.length; }, 0);
  var nC = c.reduce(function(a,s){ return a + s.piezas.length; }, 0);
  var h = '<button class="msala" data-r="dicc" data-i="-1">' +
      '<div class="mn">ANTES DE ENTRAR</div><div class="mt">📖 El diccionario</div></button>' +
    '<div class="mgrupo">EL RECORRIDO GUIADO · ' + nG + ' piezas</div>' +
    g.map(function(sa, i){
      return '<button class="msala" data-r="guiado" data-i="'+i+'">' +
        '<div class="mn">'+esc(sa.numero)+' · '+sa.piezas.length+' piezas</div>' +
        '<div class="mt">'+esc(sa.titulo)+'</div></button>';
    }).join('') +
    '<div class="mgrupo">EL MUSEO COMPLETO · ' + nC + ' piezas</div>' +
    c.map(function(sa, i){
      return '<button class="msala" data-r="completo" data-i="'+i+'">' +
        '<div class="mn">'+sa.piezas.length+' experimentos</div>' +
        '<div class="mt">'+esc(sa.numero)+' '+esc(sa.titulo)+'</div></button>';
    }).join('');
  var cont = document.getElementById('mSalas');
  cont.innerHTML = h;
  Array.prototype.forEach.call(cont.querySelectorAll('.msala'), function(b){
    b.onclick = function(){ irASala(b.dataset.r, parseInt(b.dataset.i, 10)); };
  });
}

function pintarDiccionario(){
  var c = document.getElementById('mPiezas');
  c.innerHTML = '<p class="mintro">Antes de entrar, el diccionario. Acá está traducido a palabras ' +
    'cada símbolo que vas a encontrar adentro. Nadie nace sabiendo qué es una letra griega, ' +
    'y no saberlo no te hace menos: te hace alguien que todavía no lo leyó en su idioma.</p>' +
    '<div class="dicc">' + MUSEO.diccionario.map(function(d){
      return '<div class="dfila"><div class="dsim">'+esc(d.simbolo)+'</div><div>' +
        '<div class="dnom">'+esc(d.nombre)+'</div>' +
        '<div class="dcri">'+esc(d.criollo)+'</div></div></div>';
    }).join('') + '</div>';
  c.scrollTop = 0;
}

function pieza(pz){
  var b = '<div class="pieza"><div class="pcab">' +
    '<div class="pnum">'+(pz.emoji? esc(pz.emoji)+' ' : '')+'PIEZA '+pz.n +
      (pz.hallazgo? ' · hallazgo '+esc(pz.hallazgo) : '') +
      (pz.exp? ' · <code>cmd/'+esc(pz.exp)+'</code>' : '') + '</div>' +
    '<h4>'+esc(pz.titulo)+'</h4>' +
    '<div class="pgancho">'+esc(pz.gancho)+'</div></div>' +
    '<div class="pcuerpo"><p class="pcriollo">'+esc(pz.criollo)+'</p>';
  if(pz.metafora) b += '<div class="pbloque pmeta"><span class="pet">la metáfora</span>'+esc(pz.metafora)+'</div>';
  if(pz.simbolos && pz.simbolos.length){
    b += '<div class="pbloque psim"><span class="pet">las palabras raras</span><table>' +
      pz.simbolos.map(function(x){
        return '<tr><td class="sb">'+esc(x[0])+'</td><td>'+esc(x[1])+'</td></tr>';
      }).join('') + '</table></div>';
  }
  if(pz.mirar) b += '<div class="pbloque pmirar"><span class="pet">qué estás mirando</span>'+esc(pz.mirar)+'</div>';
  if(pz.honesto) b += '<div class="pbloque phon"><span class="pet">lo honesto</span>'+esc(pz.honesto)+'</div>';
  b += '</div>';
  if(pz.lamina) b += '<div class="plam"><img loading="lazy" src="/api/laminam?f='+encodeURIComponent(pz.lamina)+'" alt=""></div>';
  if(pz.exp) b += '<div class="pacc"><button class="go" data-exp="'+esc(pz.exp)+'">▶ correr este experimento</button>' +
                  '<button data-ver="'+esc(pz.exp)+'">ver la ficha</button></div>';
  return b + '</div>';
}

function pintarSala(i){
  var sa = salasDe()[i];
  if(!sa) return;
  var c = document.getElementById('mPiezas');
  c.innerHTML = '<p class="mintro"><b>'+esc(sa.numero)+' — '+esc(sa.titulo)+'</b><br>'+esc(sa.intro)+'</p>' +
                sa.piezas.map(pieza).join('');
  Array.prototype.forEach.call(c.querySelectorAll('[data-exp]'), function(b){
    b.onclick = function(){ correr(b.dataset.exp); };
  });
  Array.prototype.forEach.call(c.querySelectorAll('[data-ver]'), function(b){
    b.onclick = function(){
      document.getElementById('mCerrarMus').click();
      // nodos[dir] es una LISTA de copias; para saltar queremos la de su sala,
      // y hay que mover la .card, no el objeto que la envuelve (bug viejo).
      var lista = nodos[b.dataset.ver] || [];
      var n = lista.filter(function(x){ return !x.jefe; })[0] || lista[0];
      if(n && n.card){
        var g = n.card.parentNode;
        if(g && g.classList.contains('oculto')){ g.classList.remove('oculto'); if(g.marcar) g.marcar(); }
        n.card.scrollIntoView({behavior:'smooth', block:'center'});
        n.card.classList.add('run');
        setTimeout(function(){
          var fa = (FICHAS[b.dataset.ver]||{}).fase || '';
          if(fa!=='corriendo' && fa!=='cola') n.card.classList.remove('run');
        }, 1800);
      }
    };
  });
  c.scrollTop = 0;
}

function irASala(r, i){
  if(r === 'dicc'){ recorrido = 'guiado'; salaActual = -1; }
  else { recorrido = r; salaActual = i; }
  Array.prototype.forEach.call(document.querySelectorAll('.msala'), function(b){
    var sel = (salaActual < 0) ? (b.dataset.r === 'dicc')
            : (b.dataset.r === recorrido && parseInt(b.dataset.i,10) === salaActual);
    b.classList.toggle('sel', sel);
  });
  if(salaActual < 0) pintarDiccionario(); else pintarSala(salaActual);
}

function cargarMuseo(){
  fetch('/api/museo').then(function(r){ return r.json(); }).then(function(d){
    MUSEO = d; pintarSalas(); irASala('dicc', -1);
  }).catch(function(){
    document.getElementById('mPiezas').innerHTML = '<p class="mintro">No pude abrir el museo.</p>';
  });
}

document.getElementById('bMuseo').onclick = function(){
  document.getElementById('museo').classList.remove('oculto');
  if(!museoCargado){ museoCargado = true; cargarMuseo(); }
};
document.getElementById('mCerrarMus').onclick = function(){
  document.getElementById('museo').classList.add('oculto');
};
document.getElementById('museo').onclick = function(e){
  if(e.target.id==='museo') document.getElementById('mCerrarMus').click();
};
document.getElementById('mDicc').onclick = function(){ irASala('dicc', -1); };
document.getElementById('mAnt').onclick = function(){
  if(!MUSEO) return;
  if(salaActual <= 0) irASala('dicc', -1); else irASala(recorrido, salaActual - 1);
};
document.getElementById('mSig').onclick = function(){
  if(!MUSEO) return;
  irASala(recorrido, Math.min(salasDe().length - 1, salaActual + 1));
};

// ---------- controles de la biblioteca ----------
var biblioCargada = false;
function abrirBiblioteca(){
  document.getElementById('biblio').classList.remove('oculto');
  if(!biblioCargada){ biblioCargada = true; cargarBiblioteca(); }
}
document.getElementById('bBiblio').onclick = abrirBiblioteca;
document.getElementById('bCerrar').onclick = function(){ document.getElementById('biblio').classList.add('oculto'); };
document.getElementById('biblio').onclick = function(e){
  if(e.target.id==='biblio') document.getElementById('bCerrar').click();
};
document.getElementById('bIndice').onclick = function(){
  var idx = document.getElementById('bIdx');
  idx.classList.toggle('oculto');
  this.textContent = idx.classList.contains('oculto') ? '☰ índice' : '☰ ocultar índice';
  if(!idx.classList.contains('oculto')) pintarIndice();
};
document.getElementById('bCrudo').onclick = function(){
  if(docActual) window.open('/api/doc?f='+encodeURIComponent(docActual), '_blank');
};
var bqTimer = null;
document.getElementById('bq').oninput = function(){
  if(!crudoActual) return;
  clearTimeout(bqTimer);
  var v = this.value.trim();
  bqTimer = setTimeout(function(){ pintarDoc(crudoActual, v.length >= 2 ? v : ''); }, 220);
};

tick(); setInterval(tick, 1000);
</script>
</body>
</html>`
