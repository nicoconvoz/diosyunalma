package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// lamina.go - the Phase IV plate. Every figure from the run that draws it.
// House rule learned the hard way: all text through esc().

const (
	fBg    = "#0b1526"
	fPanel = "#0d1830"
	fInk   = "#dce8f7"
	fDim   = "#8fb4d9"
	fGold  = "#ffd98a"
	fGreen = "#7ee0c0"
	fBlue  = "#7fb2ff"
	fRose  = "#ff9aa8"
	fGrid  = "#1d3a63"
)

type ola struct{ b strings.Builder }

func (l *ola) raw(s string) { l.b.WriteString(s + "\n") }

func esc(s string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;").Replace(s)
}

func nm(v float64, d int) string {
	return strings.Replace(fmt.Sprintf("%.*f", d, v), ".", ",", 1)
}

func (l *ola) txt(x, y, s float64, c, a, t string) {
	l.raw(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Georgia,serif" fill="%s" text-anchor="%s">%s</text>`, x, y, s, c, a, esc(t)))
}

func (l *ola) mono(x, y, s float64, c, a, t string) {
	l.raw(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Consolas,monospace" fill="%s" text-anchor="%s">%s</text>`, x, y, s, c, a, esc(t)))
}

func (l *ola) panel(x, y, w, h float64, c, t string) {
	l.raw(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="10" fill="%s" stroke="%s"/>`, x, y, w, h, fPanel, c))
	if t != "" {
		l.txt(x+14, y+24, 14.5, c, "start", t)
	}
}

func (l *ola) line(x1, y1, x2, y2 float64, c string, w float64, d string) {
	dd := ""
	if d != "" {
		dd = fmt.Sprintf(` stroke-dasharray="%s"`, d)
	}
	l.raw(fmt.Sprintf(`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="%.2f"%s/>`, x1, y1, x2, y2, c, w, dd))
}

func (l *ola) circ(x, y, r float64, f string) {
	l.raw(fmt.Sprintf(`<circle cx="%.2f" cy="%.2f" r="%.2f" fill="%s"/>`, x, y, r, f))
}

func (l *ola) rect(x, y, w, h float64, f, s string, op float64) {
	l.raw(fmt.Sprintf(`<rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" fill="%s" stroke="%s" opacity="%.2f"/>`, x, y, w, h, f, s, op))
}

func (l *ola) camino(p [][2]float64, c string, w, op float64) {
	if len(p) < 2 {
		return
	}
	var sb strings.Builder
	for i, q := range p {
		if i == 0 {
			fmt.Fprintf(&sb, "M%.2f %.2f", q[0], q[1])
		} else {
			fmt.Fprintf(&sb, " L%.2f %.2f", q[0], q[1])
		}
	}
	l.raw(fmt.Sprintf(`<path d="%s" fill="none" stroke="%s" stroke-width="%.2f" opacity="%.2f"/>`, sb.String(), c, w, op))
}

func mp(v, v0, v1, p0, p1 float64) float64 {
	if v1 == v0 {
		return p0
	}
	return p0 + (v-v0)/(v1-v0)*(p1-p0)
}

func dibujar(r R) {
	l := &ola{}
	l.raw(`<svg xmlns="http://www.w3.org/2000/svg" width="1400" height="950" viewBox="0 0 1400 950">`)
	l.raw(fmt.Sprintf(`<rect width="1400" height="950" fill="%s"/>`, fBg))
	l.raw(fmt.Sprintf(`<rect x="14" y="14" width="1372" height="922" rx="16" fill="none" stroke="%s" stroke-width="2" opacity="0.55"/>`, fGreen))
	l.txt(700, 50, 25, fInk, "middle", "🌊 EL FLUIDO — la intuición del capitán, medida: las ondas SÍ se sienten, pero sólo a la vecina")
	l.txt(700, 76, 13.5, fGold, "middle",
		"«algunas ondas grandes, otras pequeñas, otras medianas, todas resonando y creando una melodía única» — Jesús Nicolás Astorga")

	// ---- A: the medium, drawn ---------------------------------------------
	l.panel(36, 96, 650, 286, fBlue, "A · EL MEDIO — cada primo, una excitación de escala log p, todas al mismo campo")
	ax, ay := 70.0, 150.0
	for i, p := range []int{2, 3, 5, 7, 11} {
		y := ay + float64(i)*30
		lp := math.Log(float64(p))
		l.mono(ax-8, y+4, 11, fDim, "end", fmt.Sprintf("p=%d", p))
		var pts [][2]float64
		for k := 0; k <= 240; k++ {
			x := ax + float64(k)*2.1
			pts = append(pts, [2]float64{x, y - 9*math.Sin(2*math.Pi*float64(k)*2.1/(30+lp*34))})
		}
		l.camino(pts, fBlue, 1.5, 0.85)
	}
	// the common field they all couple to
	l.line(ax, ay+168, ax+520, ay+168, fGold, 2.6, "")
	for i := 0; i < 5; i++ {
		y := ay + float64(i)*30
		l.line(ax+520, y, ax+560, ay+168, fGold, 1, "3 4")
	}
	l.circ(ax+560, ay+168, 7, fGold)
	l.txt(ax+574, ay+172, 12.5, fGold, "start", "el fluido")
	l.txt(ax, ay+196, 12, fDim, "start", "cada onda aporta su escalera 2πk/log p y se acopla al campo común con la fuerza")
	l.mono(ax, ay+216, 12, fGold, "start", "Λ(p)/√p = log p · p^(−1/2)   ← la que le da la propia fórmula explícita")

	// ---- B: what the coupling does ----------------------------------------
	l.panel(700, 96, 664, 286, fRose, "B · LO QUE HACE EL ACOPLAMIENTO — y dónde satura")
	bx, by, bw, bh := 760.0, 150.0, 560.0, 150.0
	l.line(bx, by+bh, bx+bw, by+bh, fGrid, 1.4, "")
	l.line(bx, by, bx, by+bh, fGrid, 1.4, "")
	var pMin, pFrac [][2]float64
	for i, f := range r.Filas {
		x := bx + float64(i)*bw/float64(len(r.Filas)-1)
		pMin = append(pMin, [2]float64{x, mp(math.Log10(f[2]), -5.2, -2.4, by+bh, by)})
		pFrac = append(pFrac, [2]float64{x, mp(f[3], 0.10, 0, by, by+bh)})
		l.mono(x, by+bh+16, 10, fDim, "middle", nm(f[0], 2))
	}
	l.camino(pMin, fGreen, 2.6, 1)
	l.camino(pFrac, fRose, 2.6, 1)
	for _, p := range pMin {
		l.circ(p[0], p[1], 3.6, fGreen)
	}
	for _, p := range pFrac {
		l.circ(p[0], p[1], 3.6, fRose)
	}
	l.txt(bx+bw/2, by+bh+34, 11, fDim, "middle", "fuerza del acoplamiento g")
	l.txt(bx+8, by-6, 11, fGreen, "start", "verde: espaciado mínimo (sube = hay repulsión)")
	l.txt(bx+8, by+12, 11, fRose, "start", "rosa: espaciados pegados (baja = se separan)")
	l.mono(716, 356, 12, fInk, "start",
		fmt.Sprintf("mínimo: %.1e → %.1e (×%.0f)   ·   pegados: %s%% → %s%%",
			r.Filas[0][2], r.Filas[len(r.Filas)-1][2], r.Filas[len(r.Filas)-1][2]/r.Filas[0][2],
			nm(100*r.Filas[0][3], 2), nm(100*r.Filas[len(r.Filas)-1][3], 2)))

	// ---- C: the ceiling, and why -------------------------------------------
	l.panel(36, 396, 650, 300, fGold, "C · EL TECHO — el entrelazado de Cauchy, dibujado")
	cx, cy := 90.0, 460.0
	for i := 0; i < 6; i++ {
		x := cx + float64(i)*95
		l.line(x, cy, x, cy+52, fBlue, 3, "")
	}
	l.txt(cx-14, cy-10, 11.5, fBlue, "start", "niveles sin acoplar")
	for i := 0; i < 5; i++ {
		x := cx + float64(i)*95 + 60
		l.line(x, cy+70, x, cy+122, fGold, 3, "")
		l.line(cx+float64(i)*95, cy+56, x, cy+66, fDim, 1, "3 3")
		l.line(cx+float64(i+1)*95, cy+56, x, cy+66, fDim, 1, "3 3")
	}
	l.txt(cx-14, cy+142, 11.5, fGold, "start", "acoplados: EXACTAMENTE uno entre cada par de vecinos")
	l.txt(56, cy+176, 12.5, fInk, "start", "un acoplamiento de rango uno intercala: empuja cada nivel, pero")
	l.txt(56, cy+196, 12.5, fInk, "start", "no puede empujarlo más allá de sus vecinos inmediatos.")
	l.txt(56, cy+220, 13, fRose, "start", "por eso satura: con g = 1, 10 o 100 da lo mismo.")
	l.mono(56, cy+244, 11.5, fDim, "start", "Σ²(10): 7,58 sin acoplar → 6,87 acoplado → los ceros 0,336")

	// ---- D: the one-by-one experiment --------------------------------------
	l.panel(700, 396, 664, 300, fGreen, "D · EL EXPERIMENTO QUE PIDIÓ LA AUDITORA — excitaciones de a una")
	dx, dy := 730.0, 440.0
	l.mono(dx, dy, 11.5, fGold, "start", "primos   modos   Σ² sin acoplar   Σ² acoplado    mejora")
	for i, e := range r.Evol {
		y := dy + 22 + float64(i)*22
		l.mono(dx, y, 11.5, fInk, "start",
			fmt.Sprintf("%5.0f  %6.0f   %12s   %11s   %8s", e[0], e[1], nm(e[2], 4), nm(e[3], 4), nm(e[4], 4)))
	}
	l.txt(dx, dy+180, 12, fGreen, "start", "la mejora CRECE con la cantidad de ondas (0,12 → 0,86): el efecto ES colectivo.")
	l.txt(dx, dy+200, 12, fRose, "start", "pero Σ² crece mucho más rápido (0,22 → 7,58): el acoplamiento nunca la alcanza.")
	l.txt(dx, dy+226, 12.5, fGold, "start", "⟹ la intuición del fluido NO muere: muere la versión de UN SOLO modo.")
	l.txt(dx, dy+246, 12.5, fGold, "start", "el medio necesita MUCHOS modos comunes — y eso es medible.")

	// ---- footer ------------------------------------------------------------
	l.raw(fmt.Sprintf(`<rect x="36" y="712" width="1328" height="90" rx="10" fill="%s" stroke="%s" opacity="0.5"/>`, fPanel, fGold))
	l.mono(700, 740, 13, fGold, "middle", "1 = g · Σᵢ vᵢ² / (E − ωᵢ)      la secular del medio: una raíz exacta entre cada par de niveles")
	l.txt(700, 766, 12.5, fInk, "middle",
		"R6 LIMPIO: la única entrada es Λ(n). Ningún parámetro se tocó mirando los γₙ — el espectro se midió antes de compararlo.")
	l.txt(700, 788, 12.5, fRose, "middle",
		"el eco aritmético sale 103 pero es TAUTOLÓGICO: las escalas log p están puestas en la definición. Lo declaramos nosotros.")

	l.txt(700, 838, 13.5, fGreen, "middle",
		"RESULTADO NEGATIVO Y ÚTIL: el fluido de un solo modo queda descartado, con la razón exacta.")
	l.txt(700, 862, 13.5, fGreen, "middle",
		"Las ondas del capitán sí se sienten entre sí — pero por ahora sólo con la de al lado.")
	l.txt(700, 906, 13, fGold, "middle",
		"go run ./cmd/elfluido · idea de Jesús Nicolás Astorga, formalizada por la auditora · estructura cerrada no es hipótesis demostrada · Todavía no.")

	l.raw("</svg>")
	dir := filepath.Join("galeria", "laminas", "10-el-telar")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		panic(err)
	}
	full := filepath.Join(dir, "el-fluido.svg")
	if err := os.WriteFile(full, []byte(l.b.String()), 0o644); err != nil {
		panic(err)
	}
	fmt.Printf("\n🖼️  lámina escrita: %s\n", full)
}
