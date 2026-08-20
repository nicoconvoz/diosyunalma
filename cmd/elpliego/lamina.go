// lamina.go - the museum plate for EL PLIEGO.
//
// Everything drawn here comes from the run that draws it: the panels receive
// the Resultado the measurements produced, and the two spacing histograms that
// Resultado does not carry are recomputed here with the very same generators
// main.go uses. No figure appears on this plate that this program did not
// measure.
//
// House rule learned the hard way: EVERY piece of text goes through pPxml, so a
// raw < or > can never break the file.
//
// The drawing bench below is deliberately private to this package (prefix pP):
// cmd/losplanos has its own, and copying the style is cheaper than coupling.
package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// ---------------------------------------------------------------------------
// the house palette
// ---------------------------------------------------------------------------

const (
	pPbg    = "#0b1526"
	pPpanel = "#0d1830"
	pPink   = "#dce8f7"
	pPdim   = "#8fb4d9"
	pPgold  = "#ffd98a"
	pPgreen = "#7ee0c0"
	pPblue  = "#7fb2ff"
	pProse  = "#ff9aa8"
	pPgrid  = "#1d3a63"
)

// ---------------------------------------------------------------------------
// the bench
// ---------------------------------------------------------------------------

type pPlienzo struct {
	W, H float64
	b    strings.Builder
}

func pPnuevo(w, h float64) *pPlienzo {
	l := &pPlienzo{W: w, H: h}
	l.raw(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">`, w, h, w, h))
	l.raw(fmt.Sprintf(`<rect width="%.0f" height="%.0f" fill="%s"/>`, w, h, pPbg))
	l.raw(fmt.Sprintf(`<rect x="14" y="14" width="%.0f" height="%.0f" rx="16" fill="none" stroke="%s" stroke-width="2" opacity="0.55"/>`,
		w-28, h-28, pPgreen))
	return l
}

func (l *pPlienzo) raw(s string) { l.b.WriteString(s + "\n") }

// pPxml makes any caption safe as XML text content.
func pPxml(s string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;").Replace(s)
}

// pPmap sends a measured value onto a pixel range.
func pPmap(v, v0, v1, p0, p1 float64) float64 {
	if v1 == v0 {
		return p0
	}
	return p0 + (v-v0)/(v1-v0)*(p1-p0)
}

// pPnum writes a number the way the house speaks it, with a decimal comma.
func pPnum(v float64, dec int) string {
	return strings.Replace(fmt.Sprintf("%.*f", dec, v), ".", ",", 1)
}

func pPmaxf(vs ...float64) float64 {
	m := vs[0]
	for _, v := range vs {
		if v > m {
			m = v
		}
	}
	return m
}

// pPhist returns the normalised density histogram of a spacing sample, binned
// exactly the way the workshop's canto exam bins it.
func pPhist(sp []float64, nb int, sm float64) []float64 {
	h := make([]float64, nb)
	for _, s := range sp {
		i := int(s / sm * float64(nb))
		if i >= 0 && i < nb {
			h[i]++
		}
	}
	for i := range h {
		h[i] /= float64(len(sp)) * (sm / float64(nb))
	}
	return h
}

func (l *pPlienzo) txt(x, y, size float64, color, anchor, s string) {
	l.raw(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Georgia,serif" fill="%s" text-anchor="%s">%s</text>`,
		x, y, size, color, anchor, pPxml(s)))
}

func (l *pPlienzo) mono(x, y, size float64, color, anchor, s string) {
	l.raw(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Consolas,monospace" fill="%s" text-anchor="%s">%s</text>`,
		x, y, size, color, anchor, pPxml(s)))
}

func (l *pPlienzo) panel(x, y, w, h float64, color, titulo string) {
	l.raw(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="10" fill="%s" stroke="%s" opacity="0.96"/>`,
		x, y, w, h, pPpanel, color))
	if titulo != "" {
		l.txt(x+14, y+24, 14.5, color, "start", titulo)
	}
}

func (l *pPlienzo) line(x1, y1, x2, y2 float64, color string, w float64, dash string) {
	d := ""
	if dash != "" {
		d = fmt.Sprintf(` stroke-dasharray="%s"`, dash)
	}
	l.raw(fmt.Sprintf(`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="%.2f"%s/>`,
		x1, y1, x2, y2, color, w, d))
}

func (l *pPlienzo) circ(x, y, r float64, fill, stroke string, sw float64) {
	s := ""
	if stroke != "" {
		s = fmt.Sprintf(` stroke="%s" stroke-width="%.2f"`, stroke, sw)
	}
	l.raw(fmt.Sprintf(`<circle cx="%.2f" cy="%.2f" r="%.2f" fill="%s"%s/>`, x, y, r, fill, s))
}

func (l *pPlienzo) rect(x, y, w, h float64, fill, stroke string, sw, op float64) {
	l.raw(fmt.Sprintf(`<rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" fill="%s" stroke="%s" stroke-width="%.2f" opacity="%.2f"/>`,
		x, y, w, h, fill, stroke, sw, op))
}

// camino draws a polyline through measured points - the workhorse of the plate.
func (l *pPlienzo) camino(pts [][2]float64, color string, w, op float64) {
	if len(pts) < 2 {
		return
	}
	var sb strings.Builder
	for i, p := range pts {
		if i == 0 {
			fmt.Fprintf(&sb, "M%.2f %.2f", p[0], p[1])
		} else {
			fmt.Fprintf(&sb, " L%.2f %.2f", p[0], p[1])
		}
	}
	l.raw(fmt.Sprintf(`<path d="%s" fill="none" stroke="%s" stroke-width="%.2f" opacity="%.2f" stroke-linejoin="round"/>`,
		sb.String(), color, w, op))
}

// poli draws a closed shape through corner points.
func (l *pPlienzo) poli(pts [][2]float64, fill, stroke string, w, op float64) {
	var sb strings.Builder
	for i, p := range pts {
		if i == 0 {
			fmt.Fprintf(&sb, "M%.2f %.2f", p[0], p[1])
		} else {
			fmt.Fprintf(&sb, " L%.2f %.2f", p[0], p[1])
		}
	}
	sb.WriteString(" Z")
	l.raw(fmt.Sprintf(`<path d="%s" fill="%s" stroke="%s" stroke-width="%.2f" opacity="%.2f"/>`,
		sb.String(), fill, stroke, w, op))
}

// llave draws a square brace spanning x1..x2, opening downwards from y.
func (l *pPlienzo) llave(x1, x2, y, d float64, color string) {
	l.raw(fmt.Sprintf(`<path d="M%.1f %.1f L%.1f %.1f L%.1f %.1f L%.1f %.1f" fill="none" stroke="%s" stroke-width="1.6"/>`,
		x1, y+d, x1, y, x2, y, x2, y+d, color))
	m := (x1 + x2) / 2
	l.line(m, y, m, y-6, color, 1.6, "")
}

// tilde draws the green check of a satisfied requirement.
func (l *pPlienzo) tilde(x, y float64, color string) {
	l.raw(fmt.Sprintf(`<path d="M%.1f %.1f L%.1f %.1f L%.1f %.1f" fill="none" stroke="%s" stroke-width="2.3" stroke-linecap="round" stroke-linejoin="round"/>`,
		x, y, x+4.5, y+5, x+11.5, y-6.5, color))
}

// nota writes the plain-language caption inside a soft box.
func (l *pPlienzo) nota(x, y, w float64, color string, lineas []string) {
	h := 16 + float64(len(lineas))*19
	l.raw(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="9" fill="%s" stroke="%s" opacity="0.35"/>`,
		x, y, w, h, pPpanel, color))
	for i, s := range lineas {
		l.txt(x+14, y+26+float64(i)*19, 13, pPink, "start", s)
	}
}

func (l *pPlienzo) formula(x, y, w float64, s string) {
	l.raw(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="34" rx="8" fill="%s" stroke="%s" opacity="0.5"/>`,
		x, y, w, pPpanel, pPgold))
	l.mono(x+w/2, y+23, 14, pPgold, "middle", s)
}

func (l *pPlienzo) guardar(dir, name string) string {
	l.raw("</svg>")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		panic(err)
	}
	full := filepath.Join(dir, name)
	if err := os.WriteFile(full, []byte(l.b.String()), 0o644); err != nil {
		panic(err)
	}
	return full
}

// ---------------------------------------------------------------------------
// dibujar is the plate. Written after the measurements so it can only ever
// show numbers this run produced.
// ---------------------------------------------------------------------------

func dibujar(r Resultado) {
	const W, H = 1400.0, 950.0
	l := pPnuevo(W, H)

	l.txt(W/2, 52, 26, pPink, "middle",
		"⚛️ EL PLIEGO — el operador que piden R1-R5, y el requisito que lo prohíbe")
	l.txt(W/2, 80, 13.5, pPdim, "middle",
		fmt.Sprintf("no demostrar que el Átomo existe: intentar construirlo · %d ceros medidos en [%.0f, %.0f], ley suave %.1f",
			len(r.Ceros), r.T0, r.T1, r.Esperados))

	// =======================================================================
	// PANEL A — el examen no distingue
	// =======================================================================
	l.panel(30, 98, 655, 318, pPgreen, "A · EL EXAMEN NO DISTINGUE")

	mZ, sZ := media(r.CantosZ), desvio(r.CantosZ)
	mB, sB := media(r.CantosB), desvio(r.CantosB)
	mR, sR := media(r.CantosR), desvio(r.CantosR)
	cmax := pPmaxf(mZ+sZ, mB+sB, mR+sR) * 1.55

	const abx0, abx1 = 232.0, 656.0
	xa := func(v float64) float64 { return pPmap(v, 0, cmax, abx0, abx1) }

	l.txt(44, 140, 11, pPdim, "start", fmt.Sprintf("mismo examen, %d espaciados cada uno", r.N89))

	barras := []struct {
		y, m, s float64
		col     string
		nom     string
	}{
		{164, mZ, sZ, pPgreen, fmt.Sprintf("ceros de ζ · %d bloques", len(r.CantosZ))},
		{188, mB, sB, pPblue, fmt.Sprintf("matriz GUE · %d matrices", r.Ensamble)},
		{212, mR, sR, pProse, fmt.Sprintf("sorteo de Wigner · %d", len(r.CantosR))},
	}
	for _, b := range barras {
		l.txt(abx0-12, b.y+4, 11.5, pPdim, "end", b.nom)
		l.line(xa(b.m-b.s), b.y, xa(b.m+b.s), b.y, b.col, 2.4, "")
		l.line(xa(b.m-b.s), b.y-6, xa(b.m-b.s), b.y+6, b.col, 2.4, "")
		l.line(xa(b.m+b.s), b.y-6, xa(b.m+b.s), b.y+6, b.col, 2.4, "")
		l.circ(xa(b.m), b.y, 4.3, b.col, pPbg, 1.2)
		l.mono(xa(b.m+b.s)+10, b.y+4, 10.5, b.col, "start", pPnum(b.m, 4)+" ± "+pPnum(b.s, 4))
	}

	// the brace: the whole distance between the true zeros and pure noise
	l.line(xa(mZ), 148, xa(mZ), 158, pPgold, 1, "3 3")
	l.line(xa(mR), 148, xa(mR), 206, pPgold, 1, "3 3")
	l.llave(xa(mZ), xa(mR), 148, 6, pPgold)
	l.mono((xa(mZ)+xa(mR))/2, 136, 11.5, pPgold, "middle", pPnum(r.Sigmas, 2)+" σ")

	// axis of the canto statistic
	l.line(abx0, 234, abx1, 234, pPgrid, 1.4, "")
	for v := 0.0; v <= cmax; v += 0.05 {
		l.line(xa(v), 234, xa(v), 239, pPgrid, 1.2, "")
		l.mono(xa(v), 252, 10, pPdim, "middle", pPnum(v, 2))
	}
	l.mono(abx1, 252, 10, pPdim, "end", "canto")

	// the three histograms, recomputed HERE with the same generators main.go
	// uses, so the plate never draws a shape it did not measure itself.
	dh := &dado{s: 40208171}
	var poolG []float64
	for i := 0; i < 10; i++ {
		poolG = append(poolG, gueNiveles(dh, 60)...)
	}
	poolW := make([]float64, 4000)
	for i := range poolW {
		poolW[i] = dh.wignerSorteo()
	}

	const ahx0, ahx1 = 232.0, 656.0
	const ahy0, ahy1 = 272.0, 386.0
	const aymax = 1.25
	hx := func(s float64) float64 { return pPmap(s, 0, smax, ahx0, ahx1) }
	hy := func(v float64) float64 { return pPmap(v, 0, aymax, ahy1, ahy0) }

	l.line(ahx0, ahy1, ahx1, ahy1, pPgrid, 1.4, "")
	l.line(ahx0, ahy0-8, ahx0, ahy1, pPgrid, 1.4, "")

	escalon := func(h []float64, col string, op float64) {
		pts := [][2]float64{{hx(0), hy(0)}}
		for i, v := range h {
			pts = append(pts, [2]float64{hx(float64(i) * smax / nbins), hy(v)})
			pts = append(pts, [2]float64{hx(float64(i+1) * smax / nbins), hy(v)})
		}
		pts = append(pts, [2]float64{hx(smax), hy(0)})
		l.camino(pts, col, 1.5, op)
	}
	escalon(pPhist(r.SpZeros, nbins, smax), pPgreen, 0.55)
	escalon(pPhist(poolG, nbins, smax), pPblue, 0.55)
	escalon(pPhist(poolW, nbins, smax), pProse, 0.55)

	var wc [][2]float64
	for s := 0.0; s <= smax+1e-9; s += 0.01 {
		wc = append(wc, [2]float64{hx(s), hy(wigner(s))})
	}
	l.camino(wc, pPgold, 2.3, 0.95)

	l.mono(ahx0+6, ahy0-2, 10, pPdim, "start", "p(s)")
	for s := 0.0; s <= 2.0+1e-9; s += 1 {
		l.mono(hx(s), 402, 10, pPdim, "middle", pPnum(s, 0))
	}
	l.mono(ahx1, 402, 10, pPdim, "end", "espaciado s")

	// legend, in the empty column at the left of the histogram
	leyenda := []struct {
		y   float64
		col string
		nom string
		w   float64
	}{
		{298, pPgreen, "ceros de ζ", 1.5},
		{320, pPblue, "matriz GUE", 1.5},
		{342, pProse, "sorteo puro", 1.5},
		{364, pPgold, "ley de Wigner", 2.3},
	}
	for _, e := range leyenda {
		l.line(44, e.y-4, 68, e.y-4, e.col, e.w, "")
		l.txt(75, e.y, 11, pPdim, "start", e.nom)
	}
	l.txt(44, 390, 10.5, pPdim, "start", "las tres se apilan")

	// the statistic that is NOT blind: the number variance sees correlations
	if len(r.Sigma2) > 0 {
		sg := r.Sigma2[0]
		l.rect(38, 400, 640, 44, pPpanel, pProse, 1, 0.5)
		l.txt(52, 418, 11.5, pProse, "start", "pero el examen mira UN espaciado por vez. La varianza de numero mira CORRELACIONES:")
		l.mono(52, 436, 11, pPink, "start",
			fmt.Sprintf("Sigma2(%.0f): ceros %s · GUE %s · sorteo SIN memoria %s   ->   t = %s   (ya no es ciega)",
				sg[0], pPnum(sg[1], 3), pPnum(sg[2], 3), pPnum(sg[3], 3), pPnum(r.Tsigma2, 1)))
	}

	// =======================================================================
	// PANEL C — la caja que respira
	// =======================================================================
	l.panel(715, 98, 655, 318, pPblue, "C · LA CAJA QUE RESPIRA")

	l.txt(732, 148, 10.5, pPdim, "start", "a escala · ancho ∝ L, paso ∝ 2π/L")

	const cyT, cyB = 168.0, 372.0
	const kL, kg = 13.0, 62.0
	L0 := r.Ventanas[0][2]
	L1 := r.Ventanas[len(r.Ventanas)-1][2]

	// the fixed box: one width, one constant rung spacing
	cx1 := 780.0
	wF := kL * r.Lfijo
	l.rect(cx1-wF/2, cyT, wF, cyB-cyT, "none", pProse, 1.8, 0.9)
	for y := cyB - kg/r.Lfijo; y > cyT+2; y -= kg / r.Lfijo {
		l.line(cx1-wF/2+3, y, cx1+wF/2-3, y, pProse, 1.4, "")
	}
	l.txt(cx1, 388, 11, pProse, "middle", "caja FIJA")
	l.mono(cx1, 403, 10, pPdim, "middle", "L = "+pPnum(r.Lfijo, 3))

	// the box the zeros actually ask for: it has to grow with the height
	cx2 := 890.0
	w0, w1 := kL*L0, kL*L1
	l.poli([][2]float64{
		{cx2 - w0/2, cyB}, {cx2 + w0/2, cyB}, {cx2 + w1/2, cyT}, {cx2 - w1/2, cyT},
	}, "none", pPgreen, 1.8, 0.95)
	Lat := func(y float64) float64 { return pPmap(y, cyB, cyT, L0, L1) }
	for y := cyB; y > cyT+2; {
		y -= kg / Lat(y)
		if y <= cyT+2 {
			break
		}
		wy := kL * Lat(y)
		l.line(cx2-wy/2+3, y, cx2+wy/2-3, y, pPgreen, 1.4, "")
	}
	l.txt(cx2, 388, 11, pPgreen, "middle", "la caja real")
	l.mono(cx2, 403, 10, pPdim, "middle", "L "+pPnum(L0, 3)+"→"+pPnum(L1, 3))

	// the five measured windows against ln(T/2pi)
	l.txt(1010, 148, 10.5, pPdim, "start", "L medido = 2π/Δ  ·  curva ln(T/2π)")
	const csx0, csx1 = 1010.0, 1352.0
	const csy0, csy1 = 168.0, 362.0
	const cT0, cT1 = 120.0, 1030.0
	const cL0, cL1 = 2.9, 5.35
	px := func(T float64) float64 { return pPmap(T, cT0, cT1, csx0, csx1) }
	py := func(L float64) float64 { return pPmap(L, cL0, cL1, csy1, csy0) }

	l.line(csx0, csy1, csx1, csy1, pPgrid, 1.4, "")
	l.line(csx0, csy0-4, csx0, csy1, pPgrid, 1.4, "")
	for T := 200.0; T <= 1000.0+1e-9; T += 200 {
		l.line(px(T), csy1, px(T), csy1+5, pPgrid, 1.2, "")
		l.mono(px(T), 378, 9.5, pPdim, "middle", pPnum(T, 0))
	}
	for L := 3.0; L <= 5.0+1e-9; L += 1 {
		l.line(csx0-5, py(L), csx0, py(L), pPgrid, 1.2, "")
		l.mono(csx0-8, py(L)+4, 9.5, pPdim, "end", pPnum(L, 0))
	}

	var lnc [][2]float64
	for T := cT0; T <= cT1; T += 5 {
		lnc = append(lnc, [2]float64{px(T), py(math.Log(T / (2 * math.Pi)))})
	}
	l.camino(lnc, pPgold, 2.2, 0.95)
	l.mono(csx1, py(math.Log(cT1/(2*math.Pi)))-8, 10, pPgold, "end", "ln(T/2π)")

	l.line(csx0, py(r.Lfijo), csx1, py(r.Lfijo), pProse, 1.4, "5 5")
	l.mono(csx1, py(r.Lfijo)-6, 9.5, pProse, "end", "mejor caja FIJA · L = "+pPnum(r.Lfijo, 3))

	for _, v := range r.Ventanas {
		l.circ(px(v[0]), py(v[2]), 4.2, pPgreen, pPbg, 1.2)
	}
	l.txt(px(r.Ventanas[0][0])+10, py(r.Ventanas[0][2])+16, 10, pPgreen, "start", "medido")

	l.txt(1010, 398, 11, pProse, "start",
		"se equivoca hasta "+pPnum(r.PeorFijo, 1)+"% · L crece ×"+pPnum(r.Crecimiento, 3))

	// =======================================================================
	// PANEL B — el eco sí distingue (the heart of the plate)
	// =======================================================================
	l.panel(30, 428, 990, 330, pPgold, "B · EL ECO SÍ DISTINGUE — |E(T)| en los períodos k·log p")

	l.mono(44, 470, 10.5, pPgreen, "start",
		"ceros de ζ  · dentro "+pPnum(r.EcoZdentro, 5)+" · fuera "+pPnum(r.EcoZfuera, 5)+" · razón "+pPnum(r.EcoZrazon, 3))
	l.mono(44, 487, 10.5, pPblue, "start",
		"espectro GUE · dentro "+pPnum(r.EcoGdentro, 5)+" · fuera "+pPnum(r.EcoGfuera, 5)+" · razón "+pPnum(r.EcoGrazon, 3))

	const bex0, bex1 = 95.0, 1005.0
	const beyT, beyB = 524.0, 706.0
	const bT0 = 0.35
	ex := func(T float64) float64 { return pPmap(T, bT0, r.Tope, bex0, bex1) }

	const nE = 1400
	cz := make([][2]float64, 0, nE+1)
	cg := make([][2]float64, 0, nE+1)
	maxE := 0.0
	for i := 0; i <= nE; i++ {
		T := bT0 + (r.Tope-bT0)*float64(i)/float64(nE)
		vz := math.Abs(eco(r.VentanaCeros, T))
		vg := math.Abs(eco(r.NivelesGUE, T))
		if vz > maxE {
			maxE = vz
		}
		if vg > maxE {
			maxE = vg
		}
		cz = append(cz, [2]float64{ex(T), vz})
		cg = append(cg, [2]float64{ex(T), vg})
	}
	ey := func(v float64) float64 { return pPmap(v, 0, maxE*1.06, beyB, beyT) }
	for i := range cz {
		cz[i][1] = ey(cz[i][1])
		cg[i][1] = ey(cg[i][1])
	}

	// the arithmetic periods, drawn behind the curves
	for i, p := range r.Periodos {
		if p < bT0 || p > r.Tope {
			continue
		}
		x := ex(p)
		l.line(x, beyT-2, x, beyB, pPgrid, 1.1, "3 5")
		l.line(x, beyB, x, beyB+7, pPgold, 1.6, "")
		yl := beyT - 8
		if i%2 == 1 {
			yl = beyT - 23
		}
		l.mono(x, yl, 10.5, pPgold, "middle", fmt.Sprintf("log %d", int(math.Round(math.Exp(p)))))
	}

	l.line(bex0, beyB, bex1, beyB, pPgrid, 1.4, "")
	l.line(bex0, beyT-2, bex0, beyB, pPgrid, 1.4, "")
	l.camino(cg, pPblue, 1.5, 0.85)
	l.camino(cz, pPgreen, 1.8, 0.98)

	l.mono(bex0-8, beyT+10, 10, pPdim, "end", "|E|")
	for T := 0.5; T <= r.Tope+1e-9; T += 0.5 {
		l.line(ex(T), beyB, ex(T), beyB+5, pPgrid, 1.2, "")
		l.mono(ex(T), beyB+22, 10, pPdim, "middle", pPnum(T, 1))
	}
	l.mono(bex1, beyB+22, 10, pPdim, "end", "T")
	l.txt(bex0, beyB+44, 10.5, pPdim, "start",
		fmt.Sprintf("%d períodos aritméticos hasta T = %s · el GUE fue reescalado a la misma densidad media que los ceros",
			len(r.Periodos), pPnum(r.Tope, 1)))

	l.mono(738, 690, 11.5, pProse, "start",
		fmt.Sprintf("control: los MISMOS ceros medidos en %d periodos AL AZAR (no aritmeticos, %d sorteos) -> razon %s ± %s, maximo %s",
			len(r.Periodos), 120, pPnum(r.EcoAzar, 4), pPnum(r.EcoAzarDes, 4), pPnum(r.EcoAzarMax, 4)))
	// =======================================================================
	// PANEL D — el pliego
	// =======================================================================
	l.panel(1040, 428, 330, 330, pProse, "D · EL PLIEGO")

	l.txt(1054, 476, 11.5, pPdim, "start",
		fmt.Sprintf("H = diag(γ₁ … γ_%d) ya los cumple todos:", len(r.Ceros)))

	pedidos := []string{
		"R1 · autoadjunta",
		"R2 · estadística de los ceros",
		fmt.Sprintf("R3 · conteo N(T): %d vs %s", len(r.Ceros), pPnum(r.Esperados, 1)),
		"R4 · eco no trivial",
		"R5 · picos en k·log p — " + pPnum(r.EcoZrazon, 3),
	}
	for i, s := range pedidos {
		y := 502 + float64(i)*26
		l.tilde(1054, y-4, pPgreen)
		l.txt(1074, y, 11.5, pPink, "start", s)
	}

	// the cheat, drawn: a diagonal matrix that already knows the answer
	const mxa, mya, mw, mh = 1252.0, 496.0, 100.0, 100.0
	l.line(mxa, mya, mxa+10, mya, pPdim, 1.6, "")
	l.line(mxa, mya, mxa, mya+mh, pPdim, 1.6, "")
	l.line(mxa, mya+mh, mxa+10, mya+mh, pPdim, 1.6, "")
	l.line(mxa+mw, mya, mxa+mw-10, mya, pPdim, 1.6, "")
	l.line(mxa+mw, mya, mxa+mw, mya+mh, pPdim, 1.6, "")
	l.line(mxa+mw, mya+mh, mxa+mw-10, mya+mh, pPdim, 1.6, "")
	for i := 0; i < 5; i++ {
		f := (float64(i) + 0.5) / 5
		l.circ(mxa+16+f*(mw-32), mya+10+f*(mh-20), 3.4, pPgold, "", 0)
	}
	l.mono(mxa+mw/2, mya+mh+18, 10, pPdim, "middle", "diag(γ)")

	l.line(1054, 628, 1356, 628, pPgrid, 1.2, "")
	l.rect(1054, 646, 12, 12, "none", pProse, 1.8, 0.95)
	l.txt(1076, 657, 12.5, pProse, "start", "R6 · NO CIRCULARIDAD")
	l.txt(1054, 680, 11, pPink, "start", "el operador tiene que definirse SIN mirar")
	l.txt(1054, 698, 11, pPink, "start", "los ceros — y ahí choca con R1: la ley de")
	l.txt(1054, 716, 11, pPink, "start", "Weyl de una caja fija da densidad CONSTANTE.")
	l.txt(1054, 740, 11, pPgold, "start", "sin R6, el pliego se satisface con trampa.")

	// =======================================================================
	// formula strip, nota and footer
	// =======================================================================
	l.formula(30, 770, 1340,
		"N(T) = θ(T)/π + 1  y  densidad ~ (1/2π)·log(T/2π)    ·    caja fija: N(k) ~ (L/π)·k, densidad CONSTANTE")

	l.nota(30, 816, 1340, pPgreen, []string{
		"Le tomamos el mismo examen a tres candidatos: los ceros de verdad, una matriz de números al azar, y un sorteo que no sabe absolutamente nada.",
		"El examen del canto no los distingue: las tres marcas caen una encima de la otra, separadas por " + pPnum(r.Sigmas, 2) + " del ruido del propio examen. Cantar lindo no prueba nada.",
		"El eco sí los distingue: en los períodos de los números primos, los ceros pegan un salto " + pPnum(r.EcoZrazon, 1) + " veces más alto que en el resto; el azar queda plano en " + pPnum(r.EcoGrazon, 2) + ".",
		"Y queda el obstáculo medido: para seguir a los ceros, la caja tendría que ir creciendo con la altura (×" + pPnum(r.Crecimiento, 3) + " en este tramo) — y ninguna caja fija crece.",
	})

	l.txt(W/2, H-26, 13, pPgold, "middle",
		fmt.Sprintf("Reproducir: go run ./cmd/elpliego · %d ceros medidos en [%.0f, %.0f] · estructura cerrada no es hipótesis demostrada · Todavía no.",
			len(r.Ceros), r.T0, r.T1))

	full := l.guardar(filepath.Join("galeria", "laminas", "10-el-telar"), "el-pliego.svg")
	fmt.Printf("\n🖼️  lámina escrita en %s\n", full)
}
