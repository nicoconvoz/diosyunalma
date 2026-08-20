package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// lamina.go - the Phase VII plate. Its headline is a retraction of our own
// Phase VI finding, which is why it is drawn first and largest.

const (
	hBg    = "#0b1526"
	hPanel = "#0d1830"
	hInk   = "#dce8f7"
	hDim   = "#8fb4d9"
	hGold  = "#ffd98a"
	hGreen = "#7ee0c0"
	hBlue  = "#7fb2ff"
	hRose  = "#ff9aa8"
	hGrid  = "#1d3a63"
)

type tela struct{ b strings.Builder }

func (l *tela) raw(s string) { l.b.WriteString(s + "\n") }

func esc(s string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;").Replace(s)
}

func nm(v float64, d int) string {
	return strings.Replace(fmt.Sprintf("%.*f", d, v), ".", ",", 1)
}

func (l *tela) txt(x, y, s float64, c, a, t string) {
	l.raw(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Georgia,serif" fill="%s" text-anchor="%s">%s</text>`, x, y, s, c, a, esc(t)))
}

func (l *tela) mono(x, y, s float64, c, a, t string) {
	l.raw(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Consolas,monospace" fill="%s" text-anchor="%s">%s</text>`, x, y, s, c, a, esc(t)))
}

func (l *tela) panel(x, y, w, h float64, c, t string) {
	l.raw(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="10" fill="%s" stroke="%s"/>`, x, y, w, h, hPanel, c))
	if t != "" {
		l.txt(x+15, y+25, 15, c, "start", t)
	}
}

func (l *tela) rect(x, y, w, h float64, f, s string, op float64) {
	l.raw(fmt.Sprintf(`<rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" fill="%s" stroke="%s" opacity="%.2f"/>`, x, y, w, h, f, s, op))
}

func (l *tela) line(x1, y1, x2, y2 float64, c string, w float64, d string) {
	dd := ""
	if d != "" {
		dd = fmt.Sprintf(` stroke-dasharray="%s"`, d)
	}
	l.raw(fmt.Sprintf(`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="%.2f"%s/>`, x1, y1, x2, y2, c, w, dd))
}

func (l *tela) circ(x, y, r float64, f string) {
	l.raw(fmt.Sprintf(`<circle cx="%.2f" cy="%.2f" r="%.2f" fill="%s"/>`, x, y, r, f))
}

func (l *tela) camino(p [][2]float64, c string, w float64, d string) {
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
	dd := ""
	if d != "" {
		dd = fmt.Sprintf(` stroke-dasharray="%s"`, d)
	}
	l.raw(fmt.Sprintf(`<path d="%s" fill="none" stroke="%s" stroke-width="%.2f"%s/>`, sb.String(), c, w, dd))
}

func mp(v, v0, v1, p0, p1 float64) float64 {
	if v1 == v0 {
		return p0
	}
	return p0 + (v-v0)/(v1-v0)*(p1-p0)
}

func dibujar(r Res) {
	l := &tela{}
	l.raw(`<svg xmlns="http://www.w3.org/2000/svg" width="1400" height="950" viewBox="0 0 1400 950">`)
	l.raw(fmt.Sprintf(`<rect width="1400" height="950" fill="%s"/>`, hBg))
	l.raw(fmt.Sprintf(`<rect x="14" y="14" width="1372" height="922" rx="16" fill="none" stroke="%s" stroke-width="2" opacity="0.55"/>`, hGreen))
	l.txt(700, 50, 25, hInk, "middle", "🧵⏳ HILADO FINO — el chequeo de robustez tumbó nuestro propio hallazgo")
	l.txt(700, 76, 13.5, hRose, "middle",
		"la «fase rígida» de la Fase VI era la banda vaciándose. Se retira. Y el reloj de arena sobrevive donde ella no.")

	// ---- A: the retraction, drawn ------------------------------------------
	l.panel(36, 96, 700, 380, hRose, "A · LA RETRACTACIÓN — cuántos niveles quedan realmente dentro de la banda")
	ax, ay, aw, ah := 100.0, 160.0, 590.0, 240.0
	l.line(ax, ay+ah, ax+aw, ay+ah, hGrid, 1.4, "")
	l.line(ax, ay, ax, ay+ah, hGrid, 1.4, "")
	// bars: surviving levels per (A,B) with B = 32
	i := 0
	for _, p := range r.AB {
		if math.Abs(p[1]-32) > 1e-9 {
			continue
		}
		x := ax + 40 + float64(i)*110
		viv := p[8]
		h := mp(viv, 0, 400, 0, ah)
		col := hGreen
		if viv < 60 {
			col = hRose
		}
		l.rect(x, ay+ah-h, 76, h, col, col, 0.55)
		l.mono(x+38, ay+ah-h-8, 12.5, col, "middle", fmt.Sprintf("%.0f", viv))
		l.mono(x+38, ay+ah+18, 11, hDim, "middle", fmt.Sprintf("A=%.0f", p[0]))
		i++
	}
	l.line(ax, ay+ah-mp(400, 0, 400, 0, ah), ax+aw, ay+ah-mp(400, 0, 400, 0, ah), hBlue, 1.4, "6 5")
	l.txt(ax+8, ay+ah-mp(400, 0, 400, 0, ah)-8, 11.5, hBlue, "start", "400 modos: el espectro completo")
	l.line(ax, ay+ah-mp(60, 0, 400, 0, ah), ax+aw, ay+ah-mp(60, 0, 400, 0, ah), hRose, 1.4, "4 4")
	l.txt(ax+8, ay+ah-mp(60, 0, 400, 0, ah)+16, 11, hRose, "start", "abajo de acá no hay nada que medir")
	l.txt(56, 436, 13, hRose, "start",
		"en (A,B) = (30,32) sobreviven 28 de 400: el 93% del espectro se fue. El α = 0,313 de la Fase VI")
	l.txt(56, 458, 13, hRose, "start",
		"se midió sobre esas migas. No era una fase del medio — era una banda vaciándose. Queda RETIRADO.")

	// ---- B: what survives --------------------------------------------------
	l.panel(752, 96, 612, 380, hBlue, "B · LO QUE SÍ SE PUEDE MEDIR — y ahí no hay rigidez")
	bx, by := 780.0, 150.0
	l.mono(bx, by, 12.5, hGold, "start", "de 15 puntos del mapa, sólo 6 dejan espectro.")
	l.mono(bx, by+26, 12.5, hInk, "start", fmt.Sprintf("mejor de los medibles: A=%.0f B=%.0f → Σ²(10)=%s",
		r.MejAA, r.MejAB, nm(r.MejA10, 2)))
	l.mono(bx, by+50, 12.5, hInk, "start", fmt.Sprintf("más rígido medible:   A=%.0f B=%.0f → α=%s",
		r.RigAA, r.RigAB, nm(r.RigAalfa, 3)))
	l.txt(bx, by+84, 13, hRose, "start", "entre los medibles, α nunca baja de 1,5.")
	l.txt(bx, by+106, 13, hRose, "start", "α = 1 ya es «sin rigidez». No hay fase rígida.")
	l.txt(bx, by+142, 13, hInk, "start", "Y hay una lección de método:")
	l.txt(bx, by+164, 12.5, hDim, "start", "medir Σ² sin contar cuántos niveles quedaron")
	l.txt(bx, by+184, 12.5, hDim, "start", "es medir la forma del recorte, no la del medio.")
	l.txt(bx, by+206, 12.5, hGold, "start", "La Fase VI no reportaba ese número. Ahora sí.")
	l.rect(bx, by+232, 552, 62, hPanel, hGreen, 0.4)
	l.txt(bx+276, by+254, 12.5, hGreen, "middle", "el chequeo de robustez que pidió la auditora")
	l.txt(bx+276, by+274, 12.5, hGreen, "middle", "es exactamente lo que destapó esto.")

	// ---- C: the hourglass kernel -------------------------------------------
	l.panel(36, 492, 700, 320, hGold, "C · EL RELOJ DE ARENA — c(k) = A·(1 − k/k₀)/k^s, dibujado y medido")
	cx, cy, cw, ch := 100.0, 560.0, 580.0, 120.0
	l.line(cx, cy+ch/2, cx+cw, cy+ch/2, hGrid, 1.4, "")
	l.txt(cx-10, cy+ch/2+5, 11, hDim, "end", "0")
	// draw the kernel with a node at k0 = 5
	var pts [][2]float64
	for k := 1; k <= 40; k++ {
		v := (1 - float64(k)/5) / math.Sqrt(float64(k))
		x := cx + float64(k-1)*cw/39
		pts = append(pts, [2]float64{x, cy + ch/2 - v*ch/2/0.9})
	}
	l.camino(pts, hGold, 2.6, "")
	for k := 1; k <= 40; k += 3 {
		p := pts[k-1]
		l.circ(p[0], p[1], 3.4, hGold)
	}
	xk := cx + 4*cw/39
	l.line(xk, cy, xk, cy+ch, hRose, 1.6, "4 4")
	l.txt(xk, cy-8, 12, hRose, "middle", "el nodo k₀")
	l.mono(cx+cw/2, cy+ch+22, 11, hDim, "middle", "distancia k entre niveles")
	l.txt(56, 722, 12.5, hInk, "start",
		fmt.Sprintf("mejor de la familia: k₀=%.0f s=%.1f → Σ²(10) = %s contra %s sin acoplar (una mejora del %.0f%%)",
			r.MejNk0, r.MejNs, nm(r.MejN10, 3), nm(r.Base.s10, 3), 100*(1-r.MejN10/r.Base.s10)))
	l.txt(56, 744, 12.5, hGreen, "start",
		"y es el ÚNICO de los dos modelos que deja espectro medible en TODOS sus puntos: la cola larga")
	l.txt(56, 764, 12.5, hGreen, "start",
		"no vacía la banda. Eso, que no era el objetivo, es lo más sólido que trajo el dibujo del capitán.")
	l.txt(56, 790, 12.5, hRose, "start",
		fmt.Sprintf("pero el NODO en sí: gana a su gemelo de un solo signo en %d de 10, y su mejor ganancia es %.2f×.", r.Ganan, r.MejorGan))

	// ---- D: her §6 warning -------------------------------------------------
	l.panel(752, 492, 612, 320, hGreen, "D · SU ADVERTENCIA (§6): ¿cancela, o sólo alterna?")
	dx, dy := 780.0, 546.0
	l.txt(dx, dy, 12.5, hInk, "start", "«que c(k) cambie de signo no demuestra que haya")
	l.txt(dx, dy+20, 12.5, hInk, "start", "cancelación útil». Tenía razón, y lo medimos:")
	l.mono(dx, dy+52, 12.5, hGold, "start", fmt.Sprintf("Σ c(k) en el mejor caso: %s", nm(r.GanSuma, 2)))
	l.mono(dx, dy+74, 12.5, hGold, "start", fmt.Sprintf("correlación |Σ c(k)| ↔ ganancia: %s", nm(r.Corr, 3)))
	l.txt(dx, dy+106, 12.5, hDim, "start", "la correlación es fuerte y negativa: cuanto menos")
	l.txt(dx, dy+126, 12.5, hDim, "start", "suma el núcleo, más gana el nodo. Eso apunta a que")
	l.txt(dx, dy+146, 12.5, hDim, "start", "la cancelación global SÍ es el mecanismo.")
	l.txt(dx, dy+176, 12.5, hRose, "start", "PERO gana sólo en 2 de 10 configuraciones. Con esa")
	l.txt(dx, dy+196, 12.5, hRose, "start", "muestra, la correlación sugiere y no demuestra.")
	l.rect(dx, dy+220, 552, 56, hPanel, hGold, 0.4)
	l.txt(dx+276, dy+242, 12.5, hGold, "middle", "ÉXITO PARCIAL para la familia de cola larga.")
	l.txt(dx+276, dy+262, 12.5, hGold, "middle", "El nodo: PENDIENTE, no descartado.")

	l.txt(700, 852, 13, hGold, "middle",
		fmt.Sprintf("R6 limpio · %d modos · única entrada aritmética Λ(p)/√p · los dos núcleos comparados a igual fuerza total · los ceros: 0,3364", r.N))
	l.txt(700, 878, 13, hGold, "middle",
		"el dibujo del reloj de arena es de Jesús Nicolás Astorga · el pedido de hilar fino, de la auditora · las mediciones y la retractación, de este taller")
	l.txt(700, 916, 13, hGold, "middle", "go run ./cmd/hiladofino · estructura cerrada no es hipótesis demostrada · Todavía no.")

	l.raw("</svg>")
	dir := filepath.Join("galeria", "laminas", "10-el-telar")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		panic(err)
	}
	full := filepath.Join(dir, "hilado-fino.svg")
	if err := os.WriteFile(full, []byte(l.b.String()), 0o644); err != nil {
		panic(err)
	}
	fmt.Printf("\n🖼️  lámina escrita: %s\n", full)
}
