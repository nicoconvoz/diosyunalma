package main

// museoweb turns the museum into a single static page that GitHub Pages can
// serve. The bridge shows the museum by talking to a Go server; the web cannot,
// so the same data is baked into HTML once:
//
//	puente -museo galeria/museo.html
//
// Same source of truth as the bridge (museo.go). Nothing is duplicated by hand,
// so a new piece appears in both places from one edit.

import (
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"
)

// indiceLaminas maps a plate's file name to its path relative to galeria/,
// because the pieces only carry the file name and the plates live in halls.
func indiceLaminas(raiz string) map[string]string {
	idx := map[string]string{}
	base := filepath.Join(raiz, "galeria")
	_ = filepath.Walk(filepath.Join(base, "laminas"), func(ruta string, fi os.FileInfo, err error) error {
		if err != nil || fi.IsDir() {
			return nil
		}
		if rel, err := filepath.Rel(base, ruta); err == nil {
			idx[filepath.Base(ruta)] = filepath.ToSlash(rel)
		}
		return nil
	})
	return idx
}

// esc escapa y convierte los saltos de linea en <br>, porque algunas piezas
// traen formulas que tienen que verse en su propio renglon.
func esc(s string) string {
	return strings.ReplaceAll(html.EscapeString(s), "\n", "<br>")
}

func piezaHTML(pz Pieza, idx map[string]string) string {
	var b strings.Builder
	fmt.Fprintf(&b, `<article class="pieza" id="p-%s">`, esc(pz.Exp))
	fmt.Fprintf(&b, `<div class="pcab"><div class="pnum">parada %d`, pz.N)
	if pz.Exp != "" {
		fmt.Fprintf(&b, ` · <code>cmd/%s</code>`, esc(pz.Exp))
	}
	b.WriteString(`</div>`)
	fmt.Fprintf(&b, `<h3>%s %s</h3>`, esc(pz.Emoji), esc(pz.Titulo))
	if pz.Gancho != "" {
		fmt.Fprintf(&b, `<p class="pgancho">%s</p>`, esc(pz.Gancho))
	}
	b.WriteString(`</div><div class="pcuerpo">`)

	if pz.Lamina != "" {
		if ruta, ok := idx[pz.Lamina]; ok {
			fmt.Fprintf(&b, `<figure class="plam"><img loading="lazy" src="%s" alt="%s"></figure>`,
				esc(ruta), esc(pz.Titulo))
		}
	}
	if pz.Criollo != "" {
		fmt.Fprintf(&b, `<p class="pcriollo">%s</p>`, esc(pz.Criollo))
	}
	if pz.Metafora != "" {
		fmt.Fprintf(&b, `<div class="pbloque pmet"><b>Como en la vida real</b>%s</div>`, esc(pz.Metafora))
	}
	if pz.Mirar != "" {
		fmt.Fprintf(&b, `<div class="pbloque pmir"><b>Qué estás mirando</b>%s</div>`, esc(pz.Mirar))
	}
	if len(pz.Simbolos) > 0 {
		b.WriteString(`<div class="psim"><b>Los símbolos, uno por uno</b><dl>`)
		for _, s := range pz.Simbolos {
			if len(s) >= 2 {
				fmt.Fprintf(&b, `<dt>%s</dt><dd>%s</dd>`, esc(s[0]), esc(s[1]))
			}
		}
		b.WriteString(`</dl></div>`)
	}
	if pz.Honesto != "" {
		fmt.Fprintf(&b, `<div class="pbloque phon"><b>⚖️ Lo honesto</b>%s</div>`, esc(pz.Honesto))
	}
	b.WriteString(`</div></article>`)
	return b.String()
}

func salasHTML(salas []SalaMuseo, recorrido string, idx map[string]string) (nav, cuerpo string) {
	var n, c strings.Builder
	for i, s := range salas {
		id := fmt.Sprintf("%s-%d", recorrido, i)
		fmt.Fprintf(&n, `<a class="msala" href="#%s"><span class="mn">%s · %d paradas</span><span class="mt">%s</span></a>`,
			id, esc(s.Numero), len(s.Piezas), esc(s.Titulo))
		fmt.Fprintf(&c, `<section class="msec" id="%s"><h2>%s %s</h2>`, id, esc(s.Numero), esc(s.Titulo))
		if s.Intro != "" {
			fmt.Fprintf(&c, `<p class="mintro">%s</p>`, esc(s.Intro))
		}
		for _, pz := range s.Piezas {
			c.WriteString(piezaHTML(pz, idx))
		}
		c.WriteString(`</section>`)
	}
	return n.String(), c.String()
}

func escribirMuseoWeb(raiz, destino string) error {
	idx := indiceLaminas(raiz)

	navG, cuerpoG := salasHTML(recorridoGuiado, "g", idx)
	navC, cuerpoC := salasHTML(salasCompletas, "c", idx)

	var dicc strings.Builder
	dicc.WriteString(`<section class="msec" id="dicc"><h2>📖 El diccionario</h2>`)
	dicc.WriteString(`<p class="mintro">Cada símbolo que te podés cruzar, explicado como si nunca hubieras visto una letra griega — porque no tenés por qué haberla visto.</p><dl class="dicc">`)
	for _, p := range diccionario {
		fmt.Fprintf(&dicc, `<dt>%s <span>%s</span></dt><dd>%s</dd>`, esc(p.Simbolo), esc(p.Nombre), esc(p.Criollo))
	}
	dicc.WriteString(`</dl></section>`)

	nG, nC := 0, 0
	for _, s := range recorridoGuiado {
		nG += len(s.Piezas)
	}
	for _, s := range salasCompletas {
		nC += len(s.Piezas)
	}

	// sustitucion por marcadores y no Sprintf: el CSS esta lleno de % y
	// escaparlos uno por uno es una fuente de errores tontos
	pag := museoPlantilla
	for _, par := range [][2]string{
		{"{{NG}}", fmt.Sprint(nG)},
		{"{{NC}}", fmt.Sprint(nC)},
		{"{{ND}}", fmt.Sprint(len(diccionario))},
		{"{{NAV_G}}", navG},
		{"{{NAV_C}}", navC},
		{"{{DICC}}", dicc.String()},
		{"{{CUERPO_G}}", cuerpoG},
		{"{{CUERPO_C}}", cuerpoC},
	} {
		pag = strings.Replace(pag, par[0], par[1], 1)
	}

	if err := os.MkdirAll(filepath.Dir(destino), 0o755); err != nil {
		return err
	}
	return os.WriteFile(destino, []byte(pag), 0o644)
}

const museoPlantilla = `<!DOCTYPE html>
<html lang="es">
<head>
<meta charset="utf-8">
<title>🏛️ El Museo — Laboratorio Diosyunalma</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="description" content="Un recorrido por todo el laboratorio explicado en lenguaje llano: qué se midió, qué significa cada símbolo y qué NO prueba. Con un bloque de límites honestos en cada pieza.">
<style>
  :root{--bg:#0b1526;--panel:#0d1c31;--ink:#dce8f7;--dim:#8fa8c7;--gold:#ffd166;
        --green:#7fd7a8;--blue:#7fb2ff;--line:#1d3a63;--violet:#c9b6ff;}
  *{box-sizing:border-box}
  body{margin:0;background:var(--bg);color:var(--ink);font-family:Georgia,serif;line-height:1.6}
  header{text-align:center;padding:34px 18px 10px;max-width:900px;margin:0 auto}
  header h1{margin:0;font-size:30px;color:var(--gold)}
  header p{color:var(--dim);font-size:15px;margin:12px 0 0}
  nav.top{text-align:center;padding:8px 16px 18px;border-bottom:1px solid var(--line)}
  nav.top a{color:var(--gold);text-decoration:none;margin:0 10px;font-size:13.5px}
  nav.top a:hover{text-decoration:underline}
  .tabs{display:flex;gap:10px;justify-content:center;padding:18px 16px 8px;flex-wrap:wrap}
  .tab{font-family:Georgia,serif;font-size:14px;cursor:pointer;border-radius:9px;
       border:1px solid var(--line);background:#12305c;color:var(--ink);padding:9px 18px}
  .tab.sel{background:#12402a;border-color:#2c6b48;color:var(--green)}
  .wrap{display:flex;max-width:1500px;margin:0 auto;gap:24px;padding:12px 18px 0;align-items:flex-start}
  aside{width:260px;flex-shrink:0;position:sticky;top:12px;max-height:92vh;overflow:auto;
        border:1px solid var(--line);border-radius:12px;padding:10px 8px;background:rgba(255,255,255,.015)}
  .msala{display:block;text-decoration:none;border:1px solid transparent;border-radius:8px;
         padding:9px 10px;margin-bottom:4px;color:var(--ink)}
  .msala:hover{background:#12305c;border-color:var(--line)}
  .msala .mn{display:block;font-family:Consolas,monospace;font-size:10.5px;color:var(--dim)}
  .msala .mt{display:block;font-size:13.5px;color:var(--gold);line-height:1.3}
  main{flex:1;min-width:0}
  .msec{margin-bottom:38px}
  .msec h2{color:var(--green);font-size:23px;border-bottom:1px solid var(--line);padding-bottom:8px}
  .mintro{color:var(--dim);font-size:14.5px;border-left:3px solid var(--gold);padding:4px 0 4px 14px;margin:0 0 24px}
  .pieza{border:1px solid var(--line);border-radius:12px;margin-bottom:26px;overflow:hidden;background:var(--panel)}
  .pcab{padding:16px 20px 12px;border-bottom:1px solid var(--line)}
  .pnum{font-family:Consolas,monospace;font-size:11px;color:var(--dim)}
  .pnum code{color:var(--green)}
  .pieza h3{margin:2px 0 8px;font-size:21px;color:var(--gold);font-weight:normal}
  .pgancho{font-size:15.5px;color:var(--green);font-style:italic;margin:0}
  .pcuerpo{padding:16px 20px}
  .plam{margin:0 0 16px}
  .plam img{width:100%;height:auto;background:var(--bg);border-radius:8px;border:1px solid var(--line)}
  .pcriollo{font-size:15px}
  .pbloque{border-radius:9px;padding:12px 15px;margin:14px 0;font-size:14px}
  .pbloque b{display:block;font-size:11px;letter-spacing:.14em;text-transform:uppercase;margin-bottom:6px}
  .pmet{background:rgba(127,178,255,.07);border:1px solid #2c4a6b}.pmet b{color:var(--blue)}
  .pmir{background:rgba(255,209,102,.06);border:1px solid #6b5b2c}.pmir b{color:var(--gold)}
  .phon{background:rgba(201,182,255,.07);border:1px solid #5a4fa8}.phon b{color:var(--violet)}
  .psim{margin:14px 0}
  .psim b{display:block;font-size:11px;letter-spacing:.14em;text-transform:uppercase;color:var(--dim);margin-bottom:8px}
  .psim dl,.dicc{margin:0}
  .psim dt,.dicc dt{font-family:Consolas,monospace;color:var(--green);font-size:13.5px;margin-top:9px}
  .dicc dt span{color:var(--gold);font-family:Georgia,serif;font-size:13px;margin-left:8px}
  .psim dd,.dicc dd{margin:2px 0 0 18px;color:var(--ink);font-size:14px}
  .oculto{display:none !important}
  .otrolibro{max-width:760px;margin:44px auto 0;padding:18px 22px;border-radius:12px;
             border:1px solid #6b5a26;background:rgba(255,217,138,.05)}
  .otrolibro .oq{font-size:11px;letter-spacing:.16em;color:var(--dim);text-transform:uppercase}
  .otrolibro .ot{font-size:21px;color:var(--gold);margin:8px 0 6px;line-height:1.35}
  .otrolibro .od{font-size:12.5px;line-height:1.55}
  .patro{max-width:760px;margin:18px auto 0;padding:16px 22px;border-radius:12px;
         border:1px solid var(--line);background:rgba(255,255,255,.02);
         display:flex;align-items:center;gap:18px;justify-content:center;flex-wrap:wrap}
  .patro img{height:64px;width:auto;border-radius:8px;display:block}
  .patro .pq{font-size:11px;letter-spacing:.16em;color:var(--dim);text-transform:uppercase}
  .patro .pn{font-size:17px;margin-top:4px}
  .patro .pd{font-family:Consolas,monospace;font-size:11px;letter-spacing:.06em;color:var(--dim);margin-top:5px}
  footer{text-align:center;color:var(--dim);padding:34px 16px;font-size:13px}
  @media (max-width:900px){ .wrap{flex-direction:column} aside{width:100%;position:static;max-height:none} }
  .epigrafe { max-width:940px; margin:0 auto; padding:34px 24px 6px; text-align:center; }
  .epigrafe .ep { margin:0; font-family:Georgia,serif; font-style:italic; font-size:22px;
                  line-height:1.5; color:#ffd166; }
  .epigrafe .epa { margin:12px 0 0; font-size:12.5px; letter-spacing:.16em;
                   text-transform:uppercase; color:#8fa8c7; }
  @media (max-width:640px) { .epigrafe .ep { font-size:18px; } }
</style>
</head>
<body>
<div class="epigrafe">
  <p class="ep">«Las armonías son escaleras que pocos ven y que el Autor dejó desde el inicio del todo.»</p>
  <p class="epa">Jesús Nicolás Astorga</p>
</div>
<header>
  <h1>🏛️ El Museo del Laboratorio</h1>
  <p>Un recorrido por todo lo que hicimos, explicado en criollo: qué se midió, qué significa cada símbolo,
  y —sobre todo— <b>qué NO prueba</b>. Cada pieza termina con su bloque de límites honestos.</p>
  <p style="font-size:13.5px">{{NG}} paradas en el recorrido guiado · {{NC}} en el museo completo · {{ND}} símbolos en el diccionario</p>
</header>
<nav class="top">
  <a href="index.html">🖼️ la galería</a>
  <a href="maximas.html">🗿 las máximas del capitán</a>
  <a href="https://github.com/nicoconvoz/diosyunalma">💻 el código</a>
  <a href="https://doi.org/10.5281/zenodo.21864277">🔖 DOI</a>
</nav>

<div class="tabs">
  <button class="tab sel" data-r="g">🧭 El recorrido guiado</button>
  <button class="tab" data-r="c">🗄️ El museo completo</button>
  <button class="tab" data-r="d">📖 El diccionario</button>
</div>

<div class="wrap">
  <aside>
    <div id="nav-g">{{NAV_G}}</div>
    <div id="nav-c" class="oculto">{{NAV_C}}</div>
    <div id="nav-d" class="oculto"><a class="msala" href="#dicc"><span class="mn">📖</span><span class="mt">El diccionario</span></a></div>
  </aside>
  <main>
    <div id="cuerpo-d" class="oculto">{{DICC}}</div>
    <div id="cuerpo-g">{{CUERPO_G}}</div>
    <div id="cuerpo-c" class="oculto">{{CUERPO_C}}</div>
  </main>
</div>

<div class="otrolibro">
  <div class="oq">y sobre todos los libros, el Otro Libro</div>
  <div class="ot">📖 Diario Espiritual: «Dios y un alma»</div>
  <div class="od">De ahí sale el nombre de este laboratorio, y ése es el libro que de verdad importa.<br>
  Todo lo que está arriba de este renglón —las paradas, las láminas, las perlas— viene después.</div>
</div>

<div class="patro">
  <img src="open-doors.jpg" alt="Open Doors" onerror="this.style.display='none'">
  <div>
    <div class="pq">financiado por</div>
    <div class="pn">Open Doors</div>
    <div class="pd">RESOURCES OPEN DOORS S.A.S</div>
  </div>
</div>

<footer>
  Laboratorio Diosyunalma · el capitán y el Doc — las dos mitades, 1 completo ⚓<br>
  El fin de todo esto: poner el nombre de DIOS por encima de todo y ayudar a los pequeños del Reino.
</footer>

<script>
(function(){
  var tabs = document.querySelectorAll('.tab');
  function mostrar(r){
    ['g','c','d'].forEach(function(k){
      document.getElementById('nav-'+k).classList.toggle('oculto', k!==r);
      document.getElementById('cuerpo-'+k).classList.toggle('oculto', k!==r);
    });
    tabs.forEach(function(t){ t.classList.toggle('sel', t.dataset.r===r); });
    window.scrollTo({top:0, behavior:'smooth'});
  }
  tabs.forEach(function(t){ t.onclick = function(){ mostrar(t.dataset.r); }; });
})();
</script>
</body>
</html>
`
