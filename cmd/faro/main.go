// Command faro is the admiral's lighthouse: a live local dashboard.
// Land in sight (live voyage progress from the flight photographs), the
// map with every flag planted, and the last 50 log entries of the night
// log - refreshed automatically while the fleet sails.
//
// Usage:
//
//	go run ./cmd/faro    # then open http://localhost:8117
package main

import (
	"encoding/gob"
	"fmt"
	"html"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ckptFacets struct {
	T0, Span float64
	S, Shift int
	Engine   string
}

type ckptFold struct {
	T0, Span   float64
	S          int
	Engine     string
	Edge       float64
	K0         int64
	LnHi, LnLo float64
	Ph         float64
	NBlocks    int64
}

type beach struct {
	name, ordinal, status string
	t                     float64
	zeros                 int
}

var beaches = []beach{
	{"Playa I", "~10^13 (borde del mapa)", "bandera plantada", 2.447e12, 31},
	{"Playa II", "~3.56e16", "bandera plantada", 6.66e15, 8},
	{"Playa III", "~7.25e19", "bandera plantada", 1.11e19, 5},
	{"Playa IV", "~1.64e22 (doble firma)", "bandera plantada", 2.22e21, 4},
	{"Playa V", "~3.48e23 (la más honda)", "bandera plantada", 4.44e22, 6},
	{"Playa VI", "~9.28e24", "LA NAVE ESTÁ EN CAMINO", 1.11e24, 0},
}

func loadGob(path string, v interface{}) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	return gob.NewDecoder(f).Decode(v) == nil
}

func landInSight() string {
	var b strings.Builder
	found := false
	fac, _ := filepath.Glob("ckpt/facets-*.gob")
	for _, p := range fac {
		var ck ckptFacets
		if !loadGob(p, &ck) {
			continue
		}
		found = true
		st, _ := os.Stat(p)
		age := time.Since(st.ModTime()).Round(time.Second)
		pct := 100 * float64(ck.Shift+1) / 64
		fmt.Fprintf(&b, `<div class="card live"><b>⛵ t = %.3g</b> — nivel de remo: turno %d/64 (<b>%.0f%%</b>)`,
			ck.T0, ck.Shift+1, pct)
		var cf ckptFold
		if loadGob(strings.Replace(p, "facets", "fold", 1), &cf) {
			kB := 256 / math.Cbrt(9.0e-3/cf.T0)
			nTop := math.Sqrt(cf.T0 / (2 * math.Pi))
			fp := 100 * (float64(cf.K0) - kB) / (nTop - kB)
			if fp < 0 {
				fp = 0
			}
			fmt.Fprintf(&b, ` · plegado: k = %.3g (<b>%.0f%%</b>, %d bloques)`, float64(cf.K0), fp, cf.NBlocks)
		}
		fmt.Fprintf(&b, `<br><span class="dim">última fotografía hace %v — motor %s</span></div>`, age, ck.Engine)
	}
	if !found {
		b.WriteString(`<div class="card">🌙 Sin vela activa con fotografías — la nave está en puerto o acaba de zarpar (las primeras fotos tardan un turno).</div>`)
	}
	tiles, _ := filepath.Glob("luz/luz-*.gob")
	if len(tiles) > 0 {
		fmt.Fprintf(&b, `<div class="card">🗺️ Atlas de luz: <b>%d teselas</b> fotografiadas (`, len(tiles))
		names := []string{}
		for _, t := range tiles {
			names = append(names, strings.TrimSuffix(strings.TrimPrefix(filepath.Base(t), "luz-"), ".gob"))
		}
		b.WriteString(html.EscapeString(strings.Join(names, " · ")) + `) — replay en 0 ms</div>`)
	}
	return b.String()
}

type shout struct {
	t     float64
	swing float64
}

// the Carajo's shouts (expanded map, 24 strongest), strongest first.
var shouts = []shout{
	{4.78036e21, 1.74}, {1.21441e19, 1.59}, {4.62975e21, 1.56},
	{1.44879e21, 1.55}, {6.52069e20, 1.54}, {4.3533e21, 1.53},
	{1.50777e19, 1.53}, {1.14794e21, 1.53}, {1.3743e20, 1.53},
	{1.56934e19, 1.52}, {1.41961e19, 1.51}, {1.3756e19, 1.51},
	{3.15462e19, 1.51}, {6.26912e20, 1.50}, {2.80606e19, 1.50},
	{9.43626e20, 1.50}, {2.8834e19, 1.50}, {8.00594e20, 1.49},
	{5.18283e20, 1.49}, {1.59179e20, 1.49}, {1.67803e19, 1.49},
	{1.33284e22, 1.48}, {2.60923e19, 1.48}, {1.69668e19, 1.48},
}

// the lookout's LAND HO (F119): the twelve stillest harbors — pure
// stillness is the sign of land. swing holds the stillness residue.
var harbors = []shout{
	{1.02707e23, 0.012}, {2.71936e19, 0.015}, {8.46853e21, 0.016},
	{1.37247e22, 0.016}, {1.72453e23, 0.016}, {3.51433e24, 0.018},
	{7.06303e21, 0.018}, {1.89559e23, 0.018}, {1.58279e24, 0.019},
	{9.25757e23, 0.019}, {2.36578e22, 0.020}, {5.88459e23, 0.020},
}

// the melody's twelve crests (compression breeds the record pairs).
var crests = []shout{
	{7.52248e23, 1.31}, {6.86484e23, 1.27}, {2.83697e19, 1.25},
	{6.05661e19, 1.24}, {5.51263e22, 1.23}, {5.77186e20, 1.23},
	{9.5309e20, 1.21}, {2.464e23, 1.21}, {2.42855e20, 1.21},
	{1.34007e22, 1.20}, {3.14376e24, 1.20}, {1.02771e23, 1.20},
}

func exists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

func theMap() string {
	var b strings.Builder
	b.WriteString(`<svg viewBox="0 0 1240 410" style="width:100%"><defs>
<linearGradient id="mar" x1="0" y1="0" x2="1" y2="0">
<stop offset="0" stop-color="#cfe8ff"/><stop offset="0.5" stop-color="#4a90d9"/><stop offset="1" stop-color="#032b5c"/></linearGradient>
<linearGradient id="cielo" x1="0" y1="0" x2="0" y2="1">
<stop offset="0" stop-color="#0b1526"/><stop offset="1" stop-color="#16294a"/></linearGradient></defs>
<rect x="40" y="18" width="1160" height="92" fill="url(#cielo)" stroke="#2c4a78"/>
<text x="52" y="36" font-size="12" fill="#ffd166">⛈ EL CIELO DE LAS TORMENTAS (los gritos del Carajo)</text>
<rect x="40" y="110" width="1160" height="120" fill="url(#mar)" stroke="#7fb2ff"/>
<text x="52" y="128" font-size="12" fill="#eaf4ff">🌊 EL OCÉANO</text>`)
	// the map now lives where the fleet hunts: 10^18.7 .. 10^25. The
	// human map and the shallow beaches compress into a thin left band.
	const xlo, xhi = 18.7, 25.0
	px := func(t float64) float64 {
		l := math.Log10(t)
		if l < xlo {
			return 40 + 30*l/xlo // the compressed old world, 30px wide
		}
		return 70 + 1130*(l-xlo)/(xhi-xlo)
	}
	fmt.Fprintf(&b, `<rect x="40" y="110" width="30" height="120" fill="#0a1e3d"/><text x="55" y="250" font-size="9" fill="#5f7ba0" text-anchor="middle">10^0..18</text>`)
	for L := 19; L <= 25; L++ {
		x := 70 + 1130*(float64(L)-xlo)/(xhi-xlo)
		fmt.Fprintf(&b, `<line x1="%.0f" y1="230" x2="%.0f" y2="238" stroke="#8fa8c7"/><text x="%.0f" y="254" font-size="12" text-anchor="middle" fill="#8fa8c7">10^%d</text>`, x, x, x, L)
	}

	// the storms, live from the atlas: caught (whale), swept, or pending.
	for i, s := range shouts {
		x := px(s.t)
		tile := fmt.Sprintf("luz/luz-%.6g.gob", s.t)
		duel := fmt.Sprintf("luz/luz-%.6g-colosal.gob", s.t)
		icon, col, state := "🌀", "#8fa8c7", "en la mira"
		spin := true
		switch {
		case exists(duel):
			icon, col, state, spin = "🐋", "#ffd166", "TIFÓN CONFIRMADO, doble firma", false
		case exists(tile):
			icon, col, state, spin = "⛈", "#7fd7a8", "barrido", false
		}
		size := 14 + int(12*(s.swing-1.5)/0.25)
		y := 62 + float64(i%3)*16
		if spin {
			fmt.Fprintf(&b, `<g transform="translate(%.0f,%.0f)"><g><text x="0" y="0" font-size="%d" text-anchor="middle"><title>t = %.6g — marejada prevista %.2f — %s</title>%s</text><animateTransform attributeName="transform" type="rotate" from="0 0 -5" to="360 0 -5" dur="6s" repeatCount="indefinite"/></g></g>`,
				x, y, size, s.t, s.swing, state, icon)
		} else {
			fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" text-anchor="middle"><title>t = %.6g — marejada prevista %.2f — %s</title>%s</text>`,
				x, y, size, s.t, s.swing, state, icon)
		}
		fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="110" stroke="%s" stroke-width="0.7" stroke-dasharray="2,3"/>`, x, y+4, x, col)
	}

	// the melody's crests, upper sky: compression breeds the record pairs.
	for i, s := range crests {
		x := px(s.t)
		tile := fmt.Sprintf("luz/luz-%.6g.gob", s.t)
		duel := fmt.Sprintf("luz/luz-%.6g-colosal.gob", s.t)
		icon, state := "💠", "cresta en la mira"
		switch {
		case exists(duel):
			icon, state = "🐋", "monstruo doble-firmado en la cresta"
		case exists(tile):
			icon, state = "👑", "cresta cosechada"
		}
		y := 34.0 + float64(i%2)*15
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle"><title>CRESTA t = %.6g — fuerza prevista %.2f — %s</title>%s</text>`,
			x, y, s.t, s.swing, state, icon)
	}

	// the beaches in the ocean lane.
	for _, be := range beaches {
		x := px(be.t)
		ic := "🚩"
		if be.zeros == 0 {
			ic = "⭐"
		}
		fmt.Fprintf(&b, `<text x="%.0f" y="185" font-size="20" text-anchor="middle"><title>%s — t=%.3g — %s</title>%s</text>`, x, be.name, be.t, be.status, ic)
		fmt.Fprintf(&b, `<text x="%.0f" y="205" font-size="10" fill="#eaf4ff" text-anchor="middle">%s</text>`,
			x, strings.Replace(be.name, "Playa ", "", 1))
	}

	// the ship, LIVE: drawn at the freshest photograph; stale sails as
	// pale anchors. The freshest photo IS where the ship rows right now.
	type sail struct {
		t, pct float64
		age    time.Duration
	}
	var sails []sail
	fac, _ := filepath.Glob("ckpt/facets-*.gob")
	for _, p := range fac {
		var ck ckptFacets
		if !loadGob(p, &ck) {
			continue
		}
		st, _ := os.Stat(p)
		pct := 100 * float64(ck.Shift+1) / 64
		// the segmented fold (F120) photographs per segment: fold0-, fold1-,
		// ... as well as the legacy fold-. The freshest of ANY of them is
		// the ship's true last photo; their K0s give a rough fold percent.
		segs, _ := filepath.Glob(strings.Replace(p, "facets-", "fold*-", 1))
		var bestK, nSegs float64
		for _, sp := range segs {
			var cf ckptFold
			if !loadGob(sp, &cf) || cf.T0 != ck.T0 {
				continue
			}
			nSegs++
			if fs, err := os.Stat(sp); err == nil && fs.ModTime().After(st.ModTime()) {
				st = fs
			}
			if float64(cf.K0) > bestK {
				bestK = float64(cf.K0)
			}
		}
		if nSegs > 0 {
			kB := 256 / math.Cbrt(9.0e-3/ck.T0)
			nT := math.Sqrt(ck.T0 / (2 * math.Pi))
			if fp := 100 * (bestK - kB) / (nT - kB); fp > 0 && fp < 100 {
				pct = (pct + fp) / 2
			}
		}
		sails = append(sails, sail{ck.T0, pct, time.Since(st.ModTime())})
	}
	// paired sweeps put TWO live rockets on the map: give each its own
	// vertical lane (one high, one low) so they can never overlap.
	rocketLane := 0
	anchorLane := 0
	for _, s := range sails {
		x := px(s.t)
		if s.age < 3*time.Minute {
			y := 148.0
			py := 141.0
			if rocketLane%2 == 1 {
				y, py = 210.0, 224.0
			}
			rocketLane++
			fmt.Fprintf(&b, `<g><text x="%.0f" y="%.0f" font-size="38" text-anchor="middle"><title>EL DELOREAN ESTÁ AQUÍ — t=%.3g, %.0f%%, foto hace %v</title>🏎️<animateTransform attributeName="transform" type="translate" values="0 0;0 -6;0 0;0 3;0 0" dur="1.6s" repeatCount="indefinite"/></text></g>`,
				x, y, s.t, s.pct, s.age.Round(time.Second))
			fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="44" height="18" rx="9" fill="#ffd166" stroke="#5c4a00"/><text x="%.0f" y="%.0f" font-size="13" font-weight="bold" fill="#000000" text-anchor="middle">%.0f%%</text>`,
				x-22, py-13, x, py, s.pct)
			fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="20" fill="none" stroke="#ffd166" stroke-width="2"><animate attributeName="r" values="14;30;14" dur="1.8s" repeatCount="indefinite"/><animate attributeName="opacity" values="0.9;0;0.9" dur="1.8s" repeatCount="indefinite"/></circle>`, x, y-10)
			fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="60" height="6" rx="3" fill="#1c3357" stroke="#44608c"/><rect x="%.0f" y="%.0f" width="%.0f" height="6" rx="3" fill="#ffd166"><animate attributeName="opacity" values="1;0.5;1" dur="1s" repeatCount="indefinite"/></rect>`,
				x-30, y+8, x-30, y+8, 60*s.pct/100)
		} else {
			y := 168.0
			if anchorLane%2 == 1 {
				y = 218.0
			}
			anchorLane++
			fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="15" text-anchor="middle" opacity="0.45"><title>vela en pausa — t=%.3g, %.0f%%, foto hace %v</title>⚓</text>`,
				x, y, s.t, s.pct, s.age.Round(time.Minute))
		}
	}
	// LA TIERRA FIRME (F119): the lookout's stillest harbors, live from
	// the atlas — pure stillness is the sign of land.
	b.WriteString(`<rect x="40" y="286" width="1160" height="74" fill="#2a2416" stroke="#e6a53a"/>
<text x="52" y="304" font-size="12" fill="#e6a53a">🏝 LA TIERRA FIRME (los puertos quietos del vigía — F119: quietud pura = tierra)</text>`)
	for _, h := range harbors {
		x := px(h.t)
		tile := fmt.Sprintf("luz/luz-%.6g.gob", h.t)
		if exists(tile) {
			fmt.Fprintf(&b, `<text x="%.0f" y="340" font-size="20" text-anchor="middle"><title>t = %.6g — quietud prevista %.3f — TIERRA PISADA (luz archivada)</title>🏝</text><text x="%.0f" y="356" font-size="10" fill="#7fd7a8" text-anchor="middle">pisada</text>`,
				x, h.t, h.swing, x)
		} else {
			fmt.Fprintf(&b, `<text x="%.0f" y="340" font-size="17" text-anchor="middle" opacity="0.4"><title>t = %.6g — quietud prevista %.3f — por pisar</title>🏝</text><text x="%.0f" y="356" font-size="10" fill="#8fa8c7" text-anchor="middle" opacity="0.7">%.3f</text>`,
				x, h.t, h.swing, x, h.swing)
		}
		fmt.Fprintf(&b, `<line x1="%.0f" y1="230" x2="%.0f" y2="322" stroke="#e6a53a" stroke-width="0.6" stroke-dasharray="2,4" opacity="0.5"/>`, x, x)
	}
	b.WriteString(`<text x="620" y="384" font-size="12" text-anchor="middle" fill="#8fa8c7">🐋 doble firma · ⛈ barrido · 🌀 en la mira · 👑 cresta cosechada · 💠 cresta pendiente · 🚩 bandera · 🏎️ EL DELOREAN · ⚓ vela en pausa · 🏝 puerto (nítido = pisado) — mouse sobre cada símbolo</text></svg>`)
	return b.String()
}

func lastLogs() string {
	data, err := os.ReadFile("docs/BITACORA-NOCTURNA.md")
	if err != nil {
		return `<div class="card">sin bitácora</div>`
	}
	lines := []string{}
	for _, ln := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(ln) != "" {
			lines = append(lines, ln)
		}
	}
	if len(lines) > 50 {
		lines = lines[len(lines)-50:]
	}
	var b strings.Builder
	b.WriteString(`<div class="log">`)
	for _, ln := range lines {
		e := html.EscapeString(ln)
		cls := ""
		switch {
		case strings.Contains(ln, "CERTIFICADA") || strings.Contains(ln, "CERTIFICADO") || strings.Contains(ln, "esfera OK") || strings.Contains(ln, "QUIETUD"):
			cls = "ok"
		case strings.Contains(ln, "ROTA") || strings.Contains(ln, "FUGA") || strings.Contains(ln, "TIERRA"):
			cls = "bad"
		case strings.HasPrefix(ln, "##"):
			cls = "hdr"
		}
		fmt.Fprintf(&b, `<div class="%s">%s</div>`, cls, e)
	}
	b.WriteString(`</div>`)
	return b.String()
}

func flags() string {
	var b strings.Builder
	b.WriteString(`<table><tr><th>playa</th><th>altura t</th><th>ceros</th><th>primer cero (ordinal)</th><th>estado</th></tr>`)
	for _, be := range beaches {
		st := be.status
		cls := "ok"
		if be.zeros == 0 {
			cls = "warn"
		}
		fmt.Fprintf(&b, `<tr><td><b>%s</b></td><td>%.3g</td><td>%d</td><td>%s</td><td class="%s">%s</td></tr>`,
			be.name, be.t, be.zeros, be.ordinal, cls, st)
	}
	b.WriteString(`</table>`)
	return b.String()
}

func page(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<!doctype html><html><head><meta charset="utf-8">
<meta http-equiv="refresh" content="30">
<title>EL FARO — Laboratorio Diosyunalma</title><style>
body{background:#0b1526;color:#dce8f7;font-family:Georgia,serif;margin:24px;max-width:1240px}
h1{color:#ffd166;font-size:26px;margin-bottom:2px} h2{color:#7fb2ff;font-size:18px;margin:26px 0 8px}
.sub{color:#8fa8c7;font-size:13px}
.card{background:#0f1e38;border:1px solid #44608c;border-radius:10px;padding:12px 16px;margin:8px 0;font-size:15px}
.live{border-color:#ffd166} .dim{color:#8fa8c7;font-size:12px}
.log{background:#081020;border:1px solid #2c4a78;border-radius:10px;padding:12px;font-family:Consolas,monospace;font-size:12.5px;max-height:420px;overflow-y:auto}
.log div{padding:1px 4px} .ok{color:#7fd7a8} .bad{color:#ff5d73} .hdr{color:#ffd166;margin-top:6px} .warn{color:#e6a53a}
table{border-collapse:collapse;width:100%%;font-size:14px} th,td{border:1px solid #2c4a78;padding:7px 10px;text-align:left} th{color:#7fb2ff}
@keyframes brillo{0%%,100%%{box-shadow:0 0 12px #ffd166,0 0 30px #ffd16644}50%%{box-shadow:0 0 40px #ffd166,0 0 80px #ffd166aa}}
@keyframes entrada{from{transform:translate(-50%%,-40px);opacity:0}to{transform:translate(-50%%,0);opacity:1}}
#pregonero{display:none;position:fixed;top:18px;left:50%%;transform:translateX(-50%%);z-index:99;max-width:88%%;
background:#1a2a10;border:3px solid #ffd166;border-radius:14px;padding:16px 26px;font-size:19px;color:#dce8f7;
animation:brillo 1s infinite,entrada 0.4s ease-out;text-align:center}
#pregonero b{color:#ffd166}
</style>
<script>
let visto = null;
async function pregonar(){
  try{
    const r = await fetch('/presa'); const t = await r.text();
    const i = t.indexOf('|'); const n = parseInt(t.slice(0,i)); const ln = t.slice(i+1);
    if(visto === null){ visto = n; return; }
    if(n > visto){
      visto = n;
      let ico = '🌊', que = 'OLA COHERENTE';
      if(ln.includes('MUDA')){ ico='🏝'; que='ISLA MUDA'; }
      if(ln.startsWith('CARDUMEN')){ ico='🐬'; que='¡CARDUMEN!'; }
      if(ln.startsWith('MARCHA')){ ico='🚂'; que='¡LA MARCHA AVANZA!'; }
      const p = document.getElementById('pregonero');
      p.innerHTML = ico+' <b>¡TESORO CAPTURADO!</b> '+que+'<br><span style="font-size:14px;font-family:Consolas,monospace">'+ln.replace(/</g,'&lt;')+'</span>';
      p.style.display = 'block';
      clearTimeout(p._t); p._t = setTimeout(()=>{p.style.display='none';}, 15000);
    }
  }catch(e){}
}
setInterval(pregonar, 5000); pregonar();
</script>
</head><body><div id="pregonero"></div>
<h1>⚓ EL FARO</h1><div class="sub">Laboratorio Diosyunalma — %s · se refresca solo cada 30 s</div>
<h2>🔭 ¿Tierra a la vista?</h2>%s
<h2>🗺️ El mapa y las banderas</h2><div class="card">%s</div>
<h2>🚂 EL CAZADERO DEL TREN — caza continua en el abismo (F152)</h2>%s
<h2>📖 El glosario del mar — cada palabra con su matemática de verdad</h2>%s
<h2>🏝 El Atlas de la Costa — los primos cartografiados (F119)</h2>%s
<h2>🚩 Banderas plantadas</h2>%s
<h2>📖 Bitácora — últimos 50 asientos</h2>%s
</body></html>`,
		time.Now().Format("2006-01-02 15:04:05"), landInSight(), theMap(), cazadero(), glosario(), costa(), flags(), lastLogs())
}

// cazadero shows the train's standing hunt: the ABYSS MAP (every beast
// plotted in its k-band, waves and islands, the train drawn where it
// hunts) plus the latest signed catches (luz/cazadero.log).
func cazadero() string {
	var b strings.Builder
	b.WriteString(`<div class="card">`)
	data, err := os.ReadFile("luz/cazadero.log")
	if err != nil || len(data) == 0 {
		b.WriteString(`<div class="dim">el tren aún no salió de caza — lanzar: go run ./cmd/circulo -cazar</div></div>`)
		return b.String()
	}
	lines := []string{}
	for _, ln := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(ln) != "" {
			lines = append(lines, ln)
		}
	}
	total := len(lines)

	// THE TREASURE MAP OF THE ABYSS: two water lanes, every beast
	// plotted by k-band, records crowned, the fresh catch flagged.
	type catchT struct {
		frac, coh, ola, juez float64
		muda                 bool
		lane                 string
	}
	lanes := map[string][]catchT{}
	type schoolT struct {
		frac   float64
		presas int
		lane   string
	}
	var schools []schoolT
	var last catchT
	var haveLast bool
	var maxWave, minIsle catchT
	minIsle.coh = 99
	for _, ln := range lines {
		if strings.HasPrefix(ln, "CARDUMEN") {
			var tt, kk float64
			var pr int
			if n, _ := fmt.Sscanf(ln, "CARDUMEN: t=%g k=%g presas=%d", &tt, &kk, &pr); n == 3 {
				nT := math.Sqrt(tt / (2 * math.Pi))
				schools = append(schools, schoolT{frac: kk / nT, presas: pr, lane: fmt.Sprintf("%.0e", tt)})
			}
			continue
		}
		var kind string
		var tt, kk, ola, coh, juez float64
		var LL int64
		n, _ := fmt.Sscanf(ln, "BESTIA %s t=%g k=%g L=%d |ola|=%g coh=%gσ juez=%g", &kind, &tt, &kk, &LL, &ola, &coh, &juez)
		if n < 6 {
			continue
		}
		kind = strings.TrimSuffix(kind, ":")
		nT := math.Sqrt(tt / (2 * math.Pi))
		key := fmt.Sprintf("%.0e", tt)
		c := catchT{frac: kk / nT, coh: coh, ola: ola, juez: juez, muda: kind == "MUDA", lane: key}
		lanes[key] = append(lanes[key], c)
		last, haveLast = c, true
		if !c.muda && c.coh > maxWave.coh {
			maxWave = c
		}
		if c.muda && c.coh < minIsle.coh {
			minIsle = c
		}
	}
	// keep the chart readable: records span the whole book, but only
	// the freshest catches are drawn (the sonar fills the sea fast)
	for key, cs := range lanes {
		if len(cs) > 750 {
			lanes[key] = cs[len(cs)-750:]
		}
	}
	// lanes are DYNAMIC: an annexed water gets its own lane the moment
	// its first beast lands in the book
	laneKeys := make([]string, 0, len(lanes))
	for key := range lanes {
		laneKeys = append(laneKeys, key)
	}
	sort.Slice(laneKeys, func(i, j int) bool {
		a, _ := strconv.ParseFloat(laneKeys[i], 64)
		c, _ := strconv.ParseFloat(laneKeys[j], 64)
		return a < c
	})
	fmt.Fprintf(&b, `<div style="font-size:14px;color:#ffd166">🚩 aguas de caza: %s — 💀 4º fantasma CAZADO (F154): la pared cayó, bandas completas en toda agua</div>`,
		strings.Join(laneKeys, " · "))
	fresh := false
	if fi, err := os.Stat("luz/cazadero.log"); err == nil {
		fresh = time.Since(fi.ModTime()) < 2*time.Minute
	}
	// the sighting banner: a fresh target just entered the sonar
	if fresh && haveLast {
		kind, col := "OLA COHERENTE 🌊", "#7fb2ff"
		if last.muda {
			kind, col = "ISLA MUDA 🏝", "#7fd7a8"
		}
		fmt.Fprintf(&b, `<style>@keyframes cartel{0%%,100%%{box-shadow:0 0 6px #ffd166}50%%{box-shadow:0 0 26px #ffd166}}</style>
<div style="margin:10px 0;padding:10px 16px;border:2px solid #ffd166;border-radius:10px;background:#1a2a10;animation:cartel 1.2s infinite;font-size:16px">
⚑ <b style="color:#ffd166">¡OBJETIVO A LA VISTA!</b> <span style="color:%s">%s</span> en t=%s — banda k/nTop=%.3f · |ola|=%.1f · %.3fσ · juez %.1e — <b>marcada en el mapa</b></div>`,
			col, kind, last.lane, last.frac, last.ola, last.coh, last.juez)
	}
	nLanes := len(laneKeys)
	if nLanes == 0 {
		nLanes = 1
	}
	base := 78.0
	afterLanes := base + 110*float64(nLanes-1) + 28
	tickY := afterLanes + 22
	axisY := tickY + 16
	legY := axisY + 14
	svgH := legY + 76
	fmt.Fprintf(&b, `<svg viewBox="0 0 1160 %.0f" style="width:100%%">
<defs><linearGradient id="mar" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="#123564"/><stop offset="1" stop-color="#0a1f42"/></linearGradient>
<radialGradient id="ola" cx="0.5" cy="0.5" r="0.5"><stop offset="0" stop-color="#cfe6ff"/><stop offset="1" stop-color="#5a94e0"/></radialGradient></defs>
<text x="60" y="26" font-size="17" fill="#ffd166" font-family="Georgia">🧭 CARTA DEL ABISMO — el cazadero del tren</text>
<text x="1120" y="26" font-size="13" text-anchor="end" fill="#8fa8c7">N↑ · rumbo áureo</text>`, svgH)
	laneY := map[string]float64{}
	for i, key := range laneKeys {
		laneY[key] = base + 110*float64(i)
	}
	for _, key := range laneKeys {
		y := laneY[key]
		fmt.Fprintf(&b, `<rect x="60" y="%.0f" width="1060" height="56" rx="6" fill="url(#mar)" stroke="#44608c"/><text x="8" y="%.0f" font-size="13" fill="#ffd166">t=%s</text>`, y-28, y+4, key)
		// nautical grid: a line each tenth of the band axis
		for g := 1; g < 10; g++ {
			gx := 60 + 1060*float64(g)/10
			fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#2c4a78" stroke-dasharray="2,4" opacity="0.6"/>`, gx, y-28, gx, y+28)
		}
		nW, nI := 0, 0
		for _, c := range lanes[key] {
			x := 60 + 1060*c.frac
			if c.muda {
				nI++
				sz := 12 + 8*(0.05-c.coh)/0.05
				fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%.0f" text-anchor="middle"><title>ISLA — silencio %.3f&#963; · |ola|=%.1f · juez %.1e · k/nTop=%.3f</title>🏝</text>`, x, y+6, sz, c.coh, c.ola, c.juez, c.frac)
			} else {
				nW++
				r := math.Max(2.5, 2+2.5*(c.coh-2.2))
				fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="%.1f" fill="url(#ola)" opacity="0.9"><title>OLA — %.2f&#963; · |ola|=%.1f · juez %.1e · k/nTop=%.3f</title></circle>`, x, y, r, c.coh, c.ola, c.juez, c.frac)
			}
		}
		fmt.Fprintf(&b, `<text x="1124" y="%.0f" font-size="12" fill="#8fa8c7">🌊%d 🏝%d</text>`, y+4, nW, nI)
	}
	// axis ticks
	for g := 0; g <= 10; g++ {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="10" text-anchor="middle" fill="#8fa8c7">%.1f</text>`, 60+1060*float64(g)/10, tickY, float64(g)/10)
	}
	fmt.Fprintf(&b, `<text x="590" y="%.0f" font-size="11" text-anchor="middle" fill="#8fa8c7">eje: banda k/nTop (0 → 1)</text>`, axisY)
	// the dolphin's schools: prey found swimming together
	if len(schools) > 40 {
		schools = schools[len(schools)-40:]
	}
	for _, s := range schools {
		if y, ok := laneY[s.lane]; ok {
			fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="17" text-anchor="middle"><title>CARDUMEN de %d presas · k/nTop=%.3f</title>🐬</text>`, 60+1060*s.frac, y-16, s.presas, s.frac)
		}
	}
	// the crowns: the strongest wave and the deepest silence ever signed
	if y, ok := laneY[maxWave.lane]; ok {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="16" text-anchor="middle"><title>👑 LA OLA RÉCORD: %.3f&#963; · |ola|=%.1f</title>👑</text>`, 60+1060*maxWave.frac, y-32, maxWave.coh, maxWave.ola)
	}
	if y, ok := laneY[minIsle.lane]; ok {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="16" text-anchor="middle"><title>💰 EL SILENCIO RÉCORD: %.3f&#963;</title>💰</text>`, 60+1060*minIsle.frac, y-32, minIsle.coh)
	}
	// the fresh catch, flagged and ringing on the map
	if y, ok := laneY[last.lane]; fresh && haveLast && ok {
		fx := 60 + 1060*last.frac
		fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="10" fill="none" stroke="#ffd166" stroke-width="2.5"><animate attributeName="r" values="6;20;6" dur="1.2s" repeatCount="indefinite"/><animate attributeName="opacity" values="1;0.15;1" dur="1.2s" repeatCount="indefinite"/></circle><text x="%.0f" y="%.0f" font-size="18" text-anchor="middle">⚑</text>`, fx, y, fx, y-38)
	}
	// the train, hunting on the last-visited band
	if y, ok := laneY[last.lane]; haveLast && ok {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="26" text-anchor="middle">🚂<animateTransform attributeName="transform" type="translate" values="0 0;0 -4;0 0" dur="1.4s" repeatCount="indefinite"/></text>`,
			60+1060*last.frac, y+52)
	}
	// the cartouche: the legend of the chart
	fmt.Fprintf(&b, `<rect x="60" y="%.0f" width="1060" height="64" rx="10" fill="#0f1e38" stroke="#ffd166" opacity="0.95"/>
<text x="80" y="%.0f" font-size="13" fill="#dce8f7">🌊 círculo = OLA (tamaño = coherencia) · 🏝 = ISLA (tamaño = hondura del silencio) · 🐬 = CARDUMEN (presas juntas) · 👑 = ola récord · 💰 = silencio récord</text>
<text x="80" y="%.0f" font-size="13" fill="#dce8f7">⚑ + anillo dorado = presa FRESCA (recién firmada) · 🚂 = el tren, donde caza AHORA · pasar el mouse sobre cualquier presa = su ficha completa</text>`,
		legY, legY+24, legY+48)
	b.WriteString(`</svg>`)

	if len(lines) > 8 {
		lines = lines[len(lines)-8:]
	}
	fmt.Fprintf(&b, `<div class="dim">bestias firmadas: %d — últimas capturas:</div><div class="log" style="max-height:150px">`, total)
	for _, ln := range lines {
		cls := "ok"
		if strings.Contains(ln, "MUDA") {
			cls = "warn"
		}
		if strings.HasPrefix(ln, "MARCHA") || strings.HasPrefix(ln, "CARDUMEN") {
			cls = "hdr"
		}
		fmt.Fprintf(&b, `<div class="%s">%s</div>`, cls, html.EscapeString(ln))
	}
	b.WriteString(`</div></div>`)
	return b.String()
}

// glosario is the captain's demand: every sea-word of the lab defined
// in plain criollo AND with the actual mathematics — no metaphor left
// unbacked.
func glosario() string {
	type entry struct{ term, formula, criollo string }
	entries := []entry{
		{"LA SUMA MADRE (el mar)",
			"Z(t) = 2·Σ_{k=1}^{nTop} cos(θ(t) − t·ln k)/√k + resto,  nTop = ⌊√(t/2π)⌋",
			"Todo el cazadero vive adentro de esta suma: la fórmula de Riemann-Siegel para la zeta en la línea crítica. Cada término es una flecha (un vector unitario girado t·ln k) achicada por 1/√k. Un BLOQUE es un tramo de L flechas consecutivas arrancando en k₀. El mar es literalmente esta suma; navegar es recorrerla."},
		{"σ — LA COHERENCIA",
			"S = Σ_{j=0}^{L−1} peso·e^{−2πi·(t/2π)·ln(k₀+j)},  σ = |S|/√L",
			"Si las L fases fueran ruido independiente, las flechas se pisan al azar y |S| ≈ √L (paseo del borracho) → σ ≈ 1. σ dice cuántas veces MÁS (o MENOS) que el azar mide la resultante. Es el número que decide todo lo demás."},
		{"OLA (bestia coherente)",
			"σ ≥ 2.4 — probabilidad bajo azar puro: P(σ≥s) = e^{−s²} ≈ 0.3%",
			"Decenas de miles de flechas que en vez de pisarse CONSPIRAN: sus fases t·ln k caen alineadas y la resultante mide cientos. Interferencia constructiva masiva donde el azar no la explica. Las montañas de 3.9σ del libro: P ~ 2 en un millón por bloque."},
		{"ISLA (bestia muda)",
			"σ ≤ 0.05 — probabilidad bajo azar: P(σ≤s) = 1−e^{−s²} ≈ s² ≈ 0.25%",
			"Lo contrario y más raro de lo que parece: hasta 100.000 flechas que se comen ENTRE TODAS hasta casi cero (|S| < 1 con L=9.000 se ha visto). Cancelación destructiva casi perfecta. El silencio también es estructura — por eso las cazamos igual."},
		{"CARDUMEN",
			"≥ 2 bloques CONTIGUOS (k₀, k₀±L, k₀±2L…) todos anómalos y firmados",
			"La anomalía no cabe en un bloque: la estructura de fase está correlacionada a escala MAYOR que L. El delfín los encuentra persiguiendo bloques vecinos mientras sigan sonando. Récord: 16 bloques seguidos. Que las islas naden en manada es un hallazgo del laboratorio: la mudez se extiende."},
		{"EL JUEZ (y por qué 'bestia FIRMADA')",
			"e = |S_tren − S_juez| / |S_juez| < 0.05,  juez = suma directa en doble-doble (~32 dígitos)",
			"Nada entra al libro sin recálculo INDEPENDIENTE a fuerza bruta, término a término, en aritmética dd. El 'juez=1.4e-02' de cada línea es ese acuerdo: cuanto más chico, más fina la firma. Sin firma no hay bestia — es la honestidad del registro hecha número."},
		{"k/nTop — LA BANDA (el eje del mapa)",
			"posición del bloque dentro de la suma madre: 0 = primeros términos, 1 = término nTop",
			"La coordenada horizontal de la carta. Cada agua t tiene su suma madre de nTop términos; la banda dice en qué tramo de esa suma está fondeado el bloque."},
		{"EL SONAR",
			"escucha del prefijo: L_s = 1500 → 6000 → 24000 → L; se descarta si 0.35 < σ_prefijo < 1.8",
			"Antes de remar las 100.000 flechas de un bloque, se escuchan las primeras 1.500: si ese prefijo suena a mar ordinario (σ≈1), se cuelga y se sigue — el 88% del mar se descarta así. La onda se propaga por etapas y la respuesta crece con el rango, como el canto de la ballena."},
		{"LA MARCHA (aguas anexadas)",
			"sonda: mismo bloque por 2 motores independientes (cascada vs directa dd), e < 0.02, en 3 bandas distintas",
			"Cómo crece el mar navegable: un agua se declara ANEXADA cuando tres sondas en bandas distintas firman el acuerdo de los dos motores. Así cayó la frontera entera hasta 10³⁶ el día que murió el 4º fantasma (F154)."},
		{"TORMENTA",
			"marea S(t) = ceros hallados − exigidos por la ley suave N(T) = θ(T)/π + 1, en una ventana",
			"Otra caza distinta: acá no se miden bloques sino CEROS. La ley suave dicta cuántos ceros debe haber en cada tramo; se cuentan los reales; la diferencia es la marea S. Tormenta = |S| extremo. Récord del laboratorio: S = −3.00 en t≈4.78×10²¹ (faltaban 3 ceros donde la esfera exigía 5)."},
		{"CRESTA",
			"separación entre ceros vecinos ≪ espaciamiento medio (récord: 0.369 espaciamientos)",
			"Lo contrario de la tormenta: ceros que se APRIETAN. Las crestas extremas son candidatas a pares de Lehmer — dos ceros tan pegados que |Z| apenas despega entre ellos: los casi-besos del átomo, sus momentos más frágiles. Cada uno que aguanta es un test de estabilidad aprobado."},
		{"BALLENA",
			"tormenta-récord con DOBLE FIRMA: dos motores independientes, acuerdo ~10⁻¹¹",
			"La clase máxima del catálogo de tormentas: no alcanza con que un motor la vea — la ballena exige que el motor clásico y el colosal la midan por separado y coincidan a 11 dígitos. Sin doble firma, no es ballena."},
		{"BANDERA / AGUA VERIFICADA",
			"profundidad t donde tren y juez acuerdan sobre bloques reales — hoy: 7 aguas hasta 10³⁶",
			"El mapa de dónde la aritmética as-built dice la verdad. Cada bandera costó una expedición; las últimas cinco cayeron juntas cuando el forense de 256 bits mató al 4º fantasma."},
		{"¿Y DÓNDE ESTÁN LOS PRIMOS? (pregunta del capitán)",
			"la suma madre corre sobre TODOS los enteros k — los primos aparecen del otro lado: ψ(x) = x − Σ_ρ x^ρ/ρ − …",
			"En ESTE mapa (olas, islas, cardúmenes) los primos NO son los tesoros: las bestias son el clima del mar de Z, interferencias de flechas de todos los enteros. Los primos viven en los otros mapas: en el ATLAS DE LA COSTA cada isla ES un primo (altura ln p, censo 181/181); en el MAPA DEL ÁTOMO son los lazos dorados por el núcleo; en el ECO DEL MURCIÉLAGO, cada valle. La cadena completa: mar de enteros → ceros de Z → (fórmula explícita) → primos. Dos orillas de un puente — acá se caza en una, los primos viven en la otra, y el murciélago demostró que el puente anda en ambos sentidos."},
	}
	var b strings.Builder
	b.WriteString(`<div class="card">`)
	for _, e := range entries {
		fmt.Fprintf(&b, `<div style="margin:10px 0;padding:10px 14px;border-left:3px solid #ffd166;background:#0d2244;border-radius:0 8px 8px 0">
<div style="font-size:15px;color:#ffd166;font-weight:bold">%s</div>
<div style="font-family:Consolas,monospace;font-size:12.5px;color:#7fb2ff;margin:5px 0">%s</div>
<div style="font-size:13.5px;color:#dce8f7">%s</div></div>`, e.term, e.formula, e.criollo)
	}
	b.WriteString(`</div>`)
	return b.String()
}

// costa embeds the self-drawn prime atlas: the sonar travels with the
// ship, the li ruler steals no visibility, island height = ln(prime).
func costa() string {
	if !exists("galeria/laminas/05-mar-y-tren/costa-atlas.svg") {
		return `<div class="card dim">el atlas aún no está dibujado — corré: go run ./cmd/costa</div>`
	}
	return `<div class="card"><img src="/costa-atlas.svg" style="width:100%"/>
<div class="dim">sonar a bordo (T=4x, nitidez constante) · eje = regla li(x) (F114) · altura de isla = ln(primo) · censo 181/181, 0 falsas — 🏝 verde = isla verificada · celeste = banco de arena (potencia) · ▼ = potencia profunda hundida bajo el ruido (física, no falla)</div></div>`
}

// presa answers the town crier's poll: total book lines and the last
// entry, so the page can flash a glowing banner the moment a treasure
// lands — without waiting for the 30 s refresh.
func presa(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("luz/cazadero.log")
	if err != nil {
		fmt.Fprint(w, "0|")
		return
	}
	lines := []string{}
	for _, ln := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(ln) != "" {
			lines = append(lines, ln)
		}
	}
	lastLn := ""
	if len(lines) > 0 {
		lastLn = lines[len(lines)-1]
	}
	fmt.Fprintf(w, "%d|%s", len(lines), lastLn)
}

func main() {
	http.HandleFunc("/", page)
	http.HandleFunc("/presa", presa)
	http.HandleFunc("/costa-atlas.svg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		http.ServeFile(w, r, "galeria/laminas/05-mar-y-tren/costa-atlas.svg")
	})
	fmt.Println("EL FARO encendido: http://localhost:8117")
	http.ListenAndServe(":8117", nil)
}
