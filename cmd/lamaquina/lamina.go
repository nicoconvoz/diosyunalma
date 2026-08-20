package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// lamina.go - the Phase III plate. Every figure comes from the run that draws
// it. House rule learned the hard way: all text through esc().

const (
	mBg    = "#0b1526"
	mPanel = "#0d1830"
	mInk   = "#dce8f7"
	mDim   = "#8fb4d9"
	mGold  = "#ffd98a"
	mGreen = "#7ee0c0"
	mBlue  = "#7fb2ff"
	mRose  = "#ff9aa8"
	mGrid  = "#1d3a63"
)

type paño struct{ b strings.Builder }

func (l *paño) raw(s string) { l.b.WriteString(s + "\n") }

func esc(s string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;").Replace(s)
}

func nm(v float64, d int) string {
	return strings.Replace(fmt.Sprintf("%.*f", d, v), ".", ",", 1)
}

func (l *paño) txt(x, y, s float64, c, a, t string) {
	l.raw(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Georgia,serif" fill="%s" text-anchor="%s">%s</text>`, x, y, s, c, a, esc(t)))
}

func (l *paño) mono(x, y, s float64, c, a, t string) {
	l.raw(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Consolas,monospace" fill="%s" text-anchor="%s">%s</text>`, x, y, s, c, a, esc(t)))
}

func (l *paño) panel(x, y, w, h float64, c, t string) {
	l.raw(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="10" fill="%s" stroke="%s"/>`, x, y, w, h, mPanel, c))
	if t != "" {
		l.txt(x+14, y+24, 14.5, c, "start", t)
	}
}

func (l *paño) line(x1, y1, x2, y2 float64, c string, w float64, d string) {
	dd := ""
	if d != "" {
		dd = fmt.Sprintf(` stroke-dasharray="%s"`, d)
	}
	l.raw(fmt.Sprintf(`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="%.2f"%s/>`, x1, y1, x2, y2, c, w, dd))
}

func (l *paño) circ(x, y, r float64, f string) {
	l.raw(fmt.Sprintf(`<circle cx="%.2f" cy="%.2f" r="%.2f" fill="%s"/>`, x, y, r, f))
}

func (l *paño) rect(x, y, w, h float64, f, s string, op float64) {
	l.raw(fmt.Sprintf(`<rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" fill="%s" stroke="%s" opacity="%.2f"/>`, x, y, w, h, f, s, op))
}

func mp(v, v0, v1, p0, p1 float64) float64 {
	if v1 == v0 {
		return p0
	}
	return p0 + (v-v0)/(v1-v0)*(p1-p0)
}

func dibujar(r Res) {
	l := &paño{}
	l.raw(`<svg xmlns="http://www.w3.org/2000/svg" width="1400" height="950" viewBox="0 0 1400 950">`)
	l.raw(fmt.Sprintf(`<rect width="1400" height="950" fill="%s"/>`, mBg))
	l.raw(fmt.Sprintf(`<rect x="14" y="14" width="1372" height="922" rx="16" fill="none" stroke="%s" stroke-width="2" opacity="0.55"/>`, mGreen))
	l.txt(700, 52, 25, mInk, "middle", "🔥 LA MÁQUINA — Fase III: los primos no pueden ser órbitas de un sistema que concatena")
	l.txt(700, 78, 13, mDim, "middle",
		"de primo a órbita · el peso derivado, no impuesto · y un no-go nuevo con sus hipótesis exactas")

	// ---- A: the concatenation obstruction ---------------------------------
	l.panel(36, 96, 1328, 246, mRose, "A · LA OBSTRUCCIÓN DE CONCATENACIÓN — las longitudes que un sistema normal está OBLIGADO a tener")
	ax, ay, aw := 70.0, 178.0, 1260.0
	l.line(ax, ay, ax+aw, ay, mGrid, 1.6, "")
	const NV = 110
	for n := 2; n <= NV; n++ {
		x := mp(math.Log(float64(n)), math.Log(2), math.Log(NV), ax, ax+aw)
		if lambda(n) > 0 {
			l.line(x, ay, x, ay-34, mGold, 2.4, "")
			l.circ(x, ay-38, 3.4, mGold)
		} else {
			l.line(x, ay, x, ay+26, mRose, 1.6, "")
			l.line(x-4, ay+26+4, x+4, ay+26-4, mRose, 1.8, "") // struck out
			l.line(x-4, ay+26-4, x+4, ay+26+4, mRose, 1.8, "")
		}
	}
	l.txt(ax, ay-52, 12.5, mGold, "start", "hacia arriba: potencias de primo — la fórmula explícita SÍ les da peso Λ(n)")
	l.txt(ax, ay+58, 12.5, mRose, "start", "hacia abajo y tachadas: los compuestos — la fórmula explícita les da peso EXACTAMENTE CERO")
	l.mono(ax, ay+80, 11.5, mDim, "start", "eje: log n, de log 2 a log 110")
	l.rect(890, 210, 440, 106, mPanel, mRose, 0.5)
	l.mono(1110, 236, 13, mInk, "middle", fmt.Sprintf("hasta n = %d", r.TopeN))
	l.mono(1110, 258, 12.5, mDim, "middle", fmt.Sprintf("obligadas por concatenación: %d", r.Forzadas))
	l.mono(1110, 278, 12.5, mGold, "middle", fmt.Sprintf("permitidas por Weil:        %d", r.Permitidas))
	l.mono(1110, 298, 13, mRose, "middle", fmt.Sprintf("PROHIBIDAS: %d  (%s %%)", r.Prohibidas,
		nm(100*float64(r.Prohibidas)/float64(r.Forzadas), 1)))

	// ---- B: the disjoint ladders ------------------------------------------
	l.panel(36, 356, 650, 300, mBlue, "B · LA ÚNICA SALIDA: órbitas que no se tocan — y lo que cuesta")
	bx, by, bw := 70.0, 420.0, 580.0
	for i, p := range []int{2, 3, 5, 7, 11, 13} {
		y := by + float64(i)*26
		lp := math.Log(float64(p))
		l.mono(bx-10, y+4, 11, mDim, "end", fmt.Sprintf("p=%d", p))
		for m := 1; ; m++ {
			e := 2 * math.Pi * float64(m) / lp
			if e > 160 {
				break
			}
			if e < 100 {
				continue
			}
			x := mp(e, 100, 160, bx, bx+bw)
			l.line(x, y-8, x, y+8, mBlue, 1.6, "")
		}
	}
	l.txt(bx, by+178, 12, mDim, "start", "cada escalera es un primo: 2πm/log p. Sin tocarse, no hay concatenación — y no hay repulsión.")
	l.mono(bx, by+202, 12, mRose, "start", fmt.Sprintf("Σ²(10) = %s   contra   %s de los ceros", nm(r.S2a, 3), nm(r.S2v, 3)))
	l.mono(bx, by+222, 12, mRose, "start", fmt.Sprintf("densidad CONSTANTE %s/unidad · espaciado mínimo %.1e", nm(r.DensA, 2), r.MinA))
	l.mono(bx, by+242, 12, mGreen, "start", fmt.Sprintf("eco en k·log p = %s  ← la aritmética SÍ está", nm(r.EcoA, 1)))

	// ---- C: the weight derived --------------------------------------------
	l.panel(700, 356, 664, 300, mGold, "C · EL PESO, DERIVADO — la tasa de decaimiento fuerza λ = log p")
	cx, cy := 730.0, 420.0
	l.mono(cx, cy, 12.5, mInk, "start", "órbita estable (círculo):   peso ℓ, sin amortiguación")
	l.mono(cx, cy+22, 12.5, mInk, "start", "órbita hiperbólica:         ℓ / (2·senh(mλ/2))  → tasa λ")
	l.mono(cx, cy+44, 12.5, mGold, "start", "la fórmula explícita pide:  2·log p · p^(−m/2)   → tasa log p")
	l.line(cx, cy+62, cx+600, cy+62, mGrid, 1.2, "")
	l.txt(cx, cy+86, 13.5, mGreen, "start", "las dos tasas coinciden ⟺ λ = log p")
	l.txt(cx, cy+110, 12.5, mInk, "start", "el exponente de inestabilidad de la órbita ES su propia longitud:")
	l.txt(cx, cy+132, 13, mGold, "start", "exponente de Lyapunov exactamente 1 — el flujo de xp (ẋ = x, ṗ = −x)")
	l.rect(cx, cy+152, 600, 74, mPanel, mRose, 0.45)
	l.txt(cx+14, cy+176, 12, mRose, "start", "y lo que sobra después de emparejar amplitudes es −2·(1 − p^(−m)):")
	l.txt(cx+14, cy+196, 12, mRose, "start", "el SIGNO no se va con ningún ajuste. Hay que derivarlo de otra cosa.")
	l.txt(cx+14, cy+216, 11.5, mDim, "start", "(y el factor (1 − p^(−m)) es el «problema asintótico»: coinciden sólo cuando m → ∞)")

	// ---- D: where the hole is ---------------------------------------------
	l.panel(36, 670, 1328, 174, mGreen, "D · LA MATRIZ DE EVALUACIÓN — y el hueco que queda")
	cols := []string{"candidato", "R6", "es operador", "conteo", "correlaciones", "eco", "muere en"}
	xs := []float64{60, 300, 380, 500, 600, 740, 880}
	for i, c := range cols {
		l.txt(xs[i], 716, 11.5, mGold, "start", c)
	}
	filas := [][]string{
		{"inversión aritmética (F350)", "limpio", "NO", "sí", "sí", "tautológico", "no es un operador"},
		{"grafo o flujo que concatena", "limpio", "sí", "—", "—", "IMPOSIBLE", "§A: longitudes prohibidas"},
		{"escaleras disjuntas", "limpio", "sí", "NO", "NO", "sí", "conteo y correlaciones"},
		{"Berry-Keating 2011", "limpio", "sí", "sí", "?", "NO", "no tiene aritmética"},
	}
	colr := []string{mDim, mRose, mBlue, mDim}
	for j, f := range filas {
		y := 740 + float64(j)*17
		for i, v := range f {
			l.txt(xs[i], y, 11, colr[j], "start", v)
		}
	}
	l.rect(60, 812, 1280, 24, mPanel, mGreen, 0.4)
	l.txt(700, 829, 12.5, mGreen, "middle",
		"el hueco es exactamente uno: un operador con la inestabilidad de xp cuyas órbitas cerradas sean los primos SIN poder concatenarse")

	l.raw(fmt.Sprintf(`<rect x="36" y="856" width="1328" height="34" rx="8" fill="%s" stroke="%s" opacity="0.5"/>`, mPanel, mGold))
	l.mono(700, 879, 13, mGold, "middle",
		"Λ(n) = 0 para todo n compuesto  ·  pero  log a + log b = log(ab)  ·  las dos cosas no pueden convivir en un sistema que concatena")
	l.txt(700, 918, 13, mGold, "middle",
		"go run ./cmd/lamaquina · única entrada aritmética: Λ(n) · los γₙ sólo como regla · Todavía no.")

	l.raw("</svg>")
	dir := filepath.Join("galeria", "laminas", "10-el-telar")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		panic(err)
	}
	full := filepath.Join(dir, "la-maquina.svg")
	if err := os.WriteFile(full, []byte(l.b.String()), 0o644); err != nil {
		panic(err)
	}
	fmt.Printf("\n🖼️  lámina escrita: %s\n", full)
}
