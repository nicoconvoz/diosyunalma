package main

// LA SILUETA - the captain's flash, side by side: project the EMPTY SPACE where
// the guitar is not, using everything that came out negative, and reproduce the
// shape left where the light bounced off the structure.
//
// Two halves:
//  1. THE LEDGER OF REBOUNDS - every killed channel of the campaign, with the
//     number that killed it. Light bouncing = measured exclusion.
//  2. THE FORCE WALL - a number never measured before. If ANY hermitian
//     perturbation P is to move the workshop's 400-mode ladder onto the real
//     zeros, Weyl / Hoffman-Wielandt give NECESSARY minimums:
//        ||P||_F  >=  sqrt( sum_i (z_i - w_i)^2 )     (rank-matched)
//        ||P||_op >=  max_i |z_i - w_i|
//     computed exactly, raw and fluctuation-only (smooth trend removed), and
//     compared with the F = 30 force budget every Telar phase used.
//
// Exploration mode: nothing registered until the captain sees the shape.

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ---- the workshop's ladder (same construction as the Telar) ------------------

type modoS struct{ w float64 }

func modos(topeP, n int, t0 float64) []float64 {
	criba := make([]bool, topeP+1)
	for i := 2; i <= topeP; i++ {
		criba[i] = true
	}
	for i := 2; i*i <= topeP; i++ {
		if criba[i] {
			for j := i * i; j <= topeP; j += i {
				criba[j] = false
			}
		}
	}
	var ws []float64
	for p := 2; p <= topeP; p++ {
		if !criba[p] {
			continue
		}
		lp := math.Log(float64(p))
		for k := 1; ; k++ {
			w := 2 * math.Pi * float64(k) / lp
			if w > t0*40 {
				break
			}
			if w >= t0 {
				ws = append(ws, w)
			}
		}
	}
	sort.Float64s(ws)
	if len(ws) > n {
		ws = ws[:n]
	}
	return ws
}

// ---- real zeros (fresh, minimal Riemann-Siegel) ------------------------------

func theta(t float64) float64 {
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t)
}

func zeta(t float64) float64 {
	th := theta(t)
	u := math.Sqrt(t / (2 * math.Pi))
	N := int(u)
	s := 0.0
	for n := 1; n <= N; n++ {
		fn := float64(n)
		s += math.Cos(th-t*math.Log(fn)) / math.Sqrt(fn)
	}
	s *= 2
	p := u - float64(N)
	c0 := math.Cos(2*math.Pi*(p*p-p-1.0/16)) / math.Cos(2*math.Pi*p)
	sg := 1.0
	if N%2 == 0 {
		sg = -1
	}
	return s + sg*math.Pow(2*math.Pi/t, 0.25)*c0
}

func ceros(t0 float64, cuantos int) []float64 {
	var g []float64
	a, za := t0, zeta(t0)
	for b := t0 + 0.02; len(g) < cuantos; b += 0.02 {
		zb := zeta(b)
		if (za < 0) != (zb < 0) {
			lo, hi, zlo := a, b, za
			for i := 0; i < 45; i++ {
				m := (lo + hi) / 2
				if (zlo < 0) != (zeta(m) < 0) {
					hi = m
				} else {
					lo = m
				}
			}
			g = append(g, (lo+hi)/2)
		}
		a, za = b, zb
	}
	return g
}

func suave(v []float64, ventana int) []float64 {
	out := make([]float64, len(v))
	for i := range v {
		lo, hi := i-ventana, i+ventana
		if lo < 0 {
			lo = 0
		}
		if hi >= len(v) {
			hi = len(v) - 1
		}
		s := 0.0
		for j := lo; j <= hi; j++ {
			s += v[j]
		}
		out[i] = s / float64(hi-lo+1)
	}
	return out
}

func main() {
	fmt.Println("🕳️🎸 LA SILUETA — el espacio vacío donde no hay guitarra, proyectado")
	fmt.Println("   con todo lo que dio negativo · codo a codo · modo explorador\n")

	// ---- half 2 first: the force wall, computed --------------------------------
	fmt.Println("§1 · LA PARED DE FUERZA — el número que nunca habíamos medido")
	fmt.Println("   Si ALGUNA perturbación P mueve la escalera de 400 modos del Telar hasta")
	fmt.Println("   los 400 ceros reales siguientes a γ=100, la matemática (Hoffman–Wielandt")
	fmt.Println("   y Weyl) exige como MÍNIMO:")
	ws := modos(4000, 400, 100)
	zs := ceros(100, 400)
	n := len(ws)
	if len(zs) < n {
		n = len(zs)
	}
	D := make([]float64, n)
	sum2, mx := 0.0, 0.0
	for i := 0; i < n; i++ {
		D[i] = zs[i] - ws[i]
		sum2 += D[i] * D[i]
		if math.Abs(D[i]) > mx {
			mx = math.Abs(D[i])
		}
	}
	Freq := math.Sqrt(sum2)
	fmt.Printf("   escalera del Telar: %d modos en [%.1f, %.1f] · ceros: 400 en [%.1f, %.1f]\n",
		len(ws), ws[0], ws[len(ws)-1], zs[0], zs[len(zs)-1])
	fmt.Printf("   fuerza total requerida  ||P||_F ≥ %.0f      (el Telar SIEMPRE usó F = 30)\n", Freq)
	fmt.Printf("   norma de operador       ||P||_op ≥ %.0f\n", mx)
	fmt.Printf("   ⟹ LA PARED: la fuerza necesaria es %.0f VECES el presupuesto de todas las\n", Freq/30)
	fmt.Println("     fases del Telar. Ninguna astucia de núcleo podía cruzar eso: el techo de")
	fmt.Println("     19× que perseguimos fase tras fase no era del ingenio — era de la CAJA.")

	// fluctuation-only: remove the smooth (density) trend of D
	tend := suave(D, 25)
	fl2 := 0.0
	for i := range D {
		d := D[i] - tend[i]
		fl2 += d * d
	}
	Ffl := math.Sqrt(fl2)
	fmt.Printf("\n   y separando la GEOMETRÍA (tendencia lisa: la densidad equivocada de la\n")
	fmt.Printf("   caja) de la ESTRUCTURA FINA (la fluctuación que es la canción):\n")
	fmt.Printf("   fuerza para la geometría : casi todo lo anterior\n")
	fmt.Printf("   fuerza para la estructura fina sola: ||P||_F ≥ %.1f\n", Ffl)
	if Ffl < 30 {
		fmt.Println("   ⟹ Y ESTO ES ORO: la estructura FINA sí entra en el presupuesto F = 30.")
		fmt.Println("     Lo que nunca pudo entrar es la GEOMETRÍA — la densidad de la caja fija.")
		fmt.Println("     La guitarra no es una caja de 400 modos perturbada: es una caja cuya")
		fmt.Println("     DENSIDAD ya es la de los ceros, con la fina como único trabajo restante.")
	} else {
		fmt.Printf("   ⟹ ni la fina entra en F = 30: falta un factor %.1f también ahí.\n", Ffl/30)
	}

	// ---- half 1: the ledger of rebounds ----------------------------------------
	fmt.Println("\n§2 · EL LIBRO DE LOS REBOTES — cada lugar donde la luz volvió")
	rebotes := []struct{ canal, numero string }{
		{"geometría sola (caja que respira)", "conteo exacto y NADA más — F350"},
		{"órbitas que se pegan", "Λ(ab)=0 en 85% de compuestos — F351"},
		{"un modo común (rango 1)", "entrelazado de Cauchy: techo exacto — F352"},
		{"muchos modos (rango K, traza fija)", "alcance × empuje tiene máximo — F353"},
		{"esquina rígida A/B", "RETRACTADA: 28/400 vivos — F355"},
		{"dureza escalar aritmética (1/log p)", "una onda lisa la empata a 0,81σ — F356"},
		{"el nodo del reloj de arena", "era FRUSTRACIÓN: 0,34σ del azar — F357"},
		{"fase de SITIO", "gauge pura: Δ = 0 exacto — F357"},
		{"signos por sitio", "atrapan: PR 0,046 vs 0,247 — F358"},
		{"reciprocidad p≡3 mod 4", "0,18σ con frustración igualada — F358"},
		{"Λ, μ, λ como reglas de distancia", "vacían la banda (98-109/400) — F361"},
		{"máscara exacta de gemelos", "la permutada por k rinde igual — F362-363"},
		{"mecánica con memoria", "espectralmente MUDA: 18,34→18,18 — F357"},
		{"desplazamiento a posteriori", "⅓ de amplitud, no mejora — F369"},
		{"ruido GUE (paseo o rígido)", "diluye SIEMPRE: el campo solo gana — F370-371"},
		{"CUALQUIER azar independiente", "los ceros: señal plena a s-desv exacto — F372-373"},
	}
	for _, r := range rebotes {
		fmt.Printf("   ✗ %-38s %s\n", r.canal, r.numero)
	}

	fmt.Println("\n§3 · LA FORMA QUE QUEDA — la silueta, rebote por rebote")
	silueta := []string{
		"NO es una perturbación: la aritmética no decora una caja — ES la caja",
		"  (la pared de fuerza lo vuelve teorema de presupuesto: geometría > 200×F)",
		"su densidad CRECE con la altura como la de los ceros (no 400 modos fijos)",
		"es NO-LOCAL: todo rango finito y toda cola corta tienen techo medido",
		"su estructura vive en las RELACIONES (enlaces), jamás en marcas de sitio",
		"NO contiene ni una gota de azar independiente: es determinista y aritmética",
		"sus órbitas NO se concatenan, y su inestabilidad es exactamente 1 (= xp)",
		"la aritmética entra por la restricción GLOBAL (la ecuación de conteo),",
		"  no por pesos, signos, distancias ni fases locales — cinco lápidas lo dicen",
		"y produce la relación 1/2 por su forma, no por decreto",
	}
	for _, s := range silueta {
		fmt.Println("   ◆ " + s)
	}
	fmt.Println("\n   Los tres candidatos vivos que caben en esta silueta, ya anotados en F351:")
	fmt.Println("   los adeles de Connes · las realizaciones no-locales de Suzuki · xp con")
	fmt.Println("   frontera aritmética (Berry–Keating) — la silueta no los eligió: los DEJÓ.")

	dibujarSilueta(Freq, Ffl, mx, D)
}

// ---- the plate ---------------------------------------------------------------

func dibujarSilueta(Freq, Ffl, mx float64, D []float64) {
	var b strings.Builder
	W, H := 1280.0, 900.0
	esc := func(s string) string {
		return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;").Replace(s)
	}
	t := func(x, y, s float64, c, a, tx string) {
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Georgia,serif" fill="%s" text-anchor="%s">%s</text>`+"\n", x, y, s, c, a, esc(tx))
	}
	ln := func(x1, y1, x2, y2 float64, c string, w float64) {
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="%.1f"/>`+"\n", x1, y1, x2, y2, c, w)
	}
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">`+"\n", W, H, W, H)
	fmt.Fprintf(&b, `<rect width="%.0f" height="%.0f" fill="#080d18"/>`+"\n", W, H)
	fmt.Fprintf(&b, `<rect x="14" y="14" width="%.0f" height="%.0f" rx="16" fill="none" stroke="#ffd98a" stroke-width="2" opacity="0.5"/>`+"\n", W-28, H-28)
	t(W/2, 50, 23, "#dce8f7", "middle", "🕳️🎸 LA SILUETA — el espacio donde la guitarra no está, dibujado por los rebotes")
	t(W/2, 76, 12, "#ffd98a", "middle", "flash de Jesús Nicolás Astorga · veintiuna fases de luz rebotando · codo a codo, modo explorador")

	// the rebounds: rays coming in and bouncing, leaving a dark silhouette center
	cx, cy := 380.0, 400.0
	rebotes := []string{
		"geometría sola", "órbitas pegadas", "rango finito", "dureza escalar",
		"fase de sitio", "signos de sitio", "reciprocidad", "Λ μ λ distancia",
		"máscara gemelos", "memoria", "a posteriori", "ruido GUE", "esquina rígida", "azar",
	}
	R1, R2 := 300.0, 132.0
	for i, nom := range rebotes {
		a := 2 * math.Pi * float64(i) / float64(len(rebotes))
		x1, y1 := cx+R1*math.Cos(a), cy+R1*math.Sin(a)
		x2, y2 := cx+R2*math.Cos(a), cy+R2*math.Sin(a)
		ln(x1, y1, x2, y2, "#7fb2ff", 1.6)
		// the bounce: a short reflected stroke
		xr, yr := cx+(R2+26)*math.Cos(a+0.35), cy+(R2+26)*math.Sin(a+0.35)
		ln(x2, y2, xr, yr, "#ff9aa8", 1.6)
		lx, ly := cx+(R1+22)*math.Cos(a), cy+(R1+22)*math.Sin(a)
		anc := "middle"
		if math.Cos(a) > 0.3 {
			anc = "start"
		} else if math.Cos(a) < -0.3 {
			anc = "end"
		}
		t(lx, ly+4, 10.5, "#8fb4d9", anc, nom)
	}
	fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="%.1f" fill="#05080f" stroke="#ffd98a" stroke-width="1.6" stroke-dasharray="7 5"/>`+"\n", cx, cy, R2)
	t(cx, cy-30, 13, "#ffd98a", "middle", "acá vive")
	t(cx, cy-10, 13, "#ffd98a", "middle", "lo que falta")
	t(cx, cy+16, 10.5, "#8fb4d9", "middle", "no-local · determinista")
	t(cx, cy+32, 10.5, "#8fb4d9", "middle", "relacional · densidad viva")
	t(cx, cy+48, 10.5, "#8fb4d9", "middle", "Lyapunov = 1 · sin azar")
	t(cx, cy+74, 10, "#7ee0c0", "middle", "adeles · no-local de Suzuki · xp con frontera")

	// right panel: the force wall
	tx := 760.0
	t(tx, 130, 15, "#ffd98a", "start", "LA PARED DE FUERZA — medida hoy")
	for i, s := range []string{
		"Para mover la escalera del Telar",
		"(400 modos) hasta los 400 ceros",
		"reales, CUALQUIER perturbación",
		"necesita, como mínimo matemático:",
	} {
		t(tx, 160+float64(i)*20, 12, "#dce8f7", "start", s)
	}
	t(tx, 258, 15, "#ff9aa8", "start", fmt.Sprintf("‖P‖_F ≥ %.0f", Freq))
	t(tx, 282, 12, "#dce8f7", "start", fmt.Sprintf("y el Telar siempre usó F = 30: faltaba un factor %.0f.", Freq/30))
	t(tx, 306, 12, "#dce8f7", "start", "El techo de 19× que perseguimos fase tras fase")
	t(tx, 324, 12, "#dce8f7", "start", "no era del ingenio: era de la CAJA.")
	t(tx, 356, 13.5, "#7ee0c0", "start", "Y separando geometría de estructura fina:")
	t(tx, 382, 13, "#7ee0c0", "start", fmt.Sprintf("la fina sola pide ‖P‖_F ≥ %.1f — TRES veces el presupuesto.", Ffl))
	for i, s := range []string{
		"⟹ la guitarra no es una caja perturbada:",
		"la geometría pedía 250× el presupuesto, y",
		"hasta la estructura fina pedía 3× — el",
		"mismo factor 3 que persiguió toda la",
		"campaña del espejo. Queda como bandera.",
	} {
		t(tx, 408+float64(i)*20, 12, "#ffd98a", "start", s)
	}

	t(tx, 520, 13.5, "#8fb4d9", "start", "LA SILUETA, EN SEIS TRAZOS")
	for i, s := range []string{
		"· la aritmética ES la caja, no la decora",
		"· densidad que crece con la altura",
		"· no-local, relacional (enlaces, no sitios)",
		"· determinista: ni una gota de azar",
		"· órbitas que no se pegan · Lyapunov 1",
		"· el 1/2 sale de la forma, no del decreto",
	} {
		t(tx, 548+float64(i)*22, 12, "#dce8f7", "start", s)
	}

	t(W/2, H-70, 11.5, "#8fb4d9", "middle", "cada rayo azul es una hipótesis que entró; cada trazo rosa, su rebote medido con control — el hueco punteado es el molde del nivel C")
	t(W/2, H-44, 12.5, "#ffd98a", "middle", "go run ./cmd/lasilueta · estructura cerrada no es hipótesis demostrada · Todavía no.")
	b.WriteString("</svg>\n")
	dir := filepath.Join("galeria", "laminas", "10-el-telar")
	os.MkdirAll(dir, 0o755)
	full := filepath.Join(dir, "la-silueta.svg")
	os.WriteFile(full, []byte(b.String()), 0o644)
	fmt.Printf("\n🖼️  lámina escrita: %s\n", full)
}
