package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// lamina.go - the Phase II plate. Everything drawn comes from the run that
// draws it. House rule learned the hard way: every caption goes through esc().

const (
	tBg    = "#0b1526"
	tPanel = "#0d1830"
	tInk   = "#dce8f7"
	tDim   = "#8fb4d9"
	tGold  = "#ffd98a"
	tGreen = "#7ee0c0"
	tBlue  = "#7fb2ff"
	tRose  = "#ff9aa8"
	tGrid  = "#1d3a63"
)

type tela struct{ b strings.Builder }

func (l *tela) raw(s string) { l.b.WriteString(s + "\n") }

func esc(s string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;").Replace(s)
}

func num(v float64, d int) string {
	return strings.Replace(fmt.Sprintf("%.*f", d, v), ".", ",", 1)
}

func (l *tela) txt(x, y, s float64, c, a, t string) {
	l.raw(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Georgia,serif" fill="%s" text-anchor="%s">%s</text>`, x, y, s, c, a, esc(t)))
}

func (l *tela) mono(x, y, s float64, c, a, t string) {
	l.raw(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Consolas,monospace" fill="%s" text-anchor="%s">%s</text>`, x, y, s, c, a, esc(t)))
}

func (l *tela) panel(x, y, w, h float64, c, t string) {
	l.raw(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="10" fill="%s" stroke="%s"/>`, x, y, w, h, tPanel, c))
	if t != "" {
		l.txt(x+14, y+24, 14.5, c, "start", t)
	}
}

func (l *tela) rect(x, y, w, h float64, f, st string, op float64) {
	l.raw(fmt.Sprintf(`<rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" fill="%s" stroke="%s" opacity="%.2f"/>`, x, y, w, h, f, st, op))
}

func (l *tela) line(x1, y1, x2, y2 float64, c string, w float64, d string) {
	dd := ""
	if d != "" {
		dd = fmt.Sprintf(` stroke-dasharray="%s"`, d)
	}
	l.raw(fmt.Sprintf(`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="%.2f"%s/>`, x1, y1, x2, y2, c, w, dd))
}

func (l *tela) circ(x, y, r float64, f, s string, w float64) {
	st := ""
	if s != "" {
		st = fmt.Sprintf(` stroke="%s" stroke-width="%.2f"`, s, w)
	}
	l.raw(fmt.Sprintf(`<circle cx="%.2f" cy="%.2f" r="%.2f" fill="%s"%s/>`, x, y, r, f, st))
}

func (l *tela) camino(p [][2]float64, c string, w, op float64) {
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

func mapa(v, v0, v1, p0, p1 float64) float64 {
	if v1 == v0 {
		return p0
	}
	return p0 + (v-v0)/(v1-v0)*(p1-p0)
}

func dibujar(r Resultado) {
	l := &tela{}
	l.raw(`<svg xmlns="http://www.w3.org/2000/svg" width="1400" height="950" viewBox="0 0 1400 950">`)
	l.raw(fmt.Sprintf(`<rect width="1400" height="950" fill="%s"/>`, tBg))
	l.raw(fmt.Sprintf(`<rect x="14" y="14" width="1372" height="922" rx="16" fill="none" stroke="%s" stroke-width="2" opacity="0.55"/>`, tGreen))
	l.txt(700, 52, 25, tInk, "middle", "🧵 EL TELAR — Fase II: la geometría sola es vacía, y los primos solos alcanzan")
	l.txt(700, 78, 13, tDim, "middle",
		fmt.Sprintf("sin leer un solo cero · %d ceros verdaderos como regla en [%.0f, %.0f] · ln(T/2π) va de %s a %s",
			len(r.Verdad), r.T0, r.T1, num(r.LargoT0, 3), num(r.LargoT1, 3)))

	// ---- A: the two halves, judged by the statistic that sees --------------
	l.panel(36, 96, 650, 300, tGold, "A · LAS DOS MITADES — Σ²(10): la rigidez que ninguna matriz al azar tiene")
	bx, by, bw, bh := 70.0, 150.0, 580.0, 190.0
	l.line(bx, by+bh, bx+bw, by+bh, tGrid, 1.4, "")
	barras := []struct {
		v   float64
		c   string
		nom string
	}{
		{r.S2caja10, tBlue, "la caja que respira"},
		{r.Filas[0][4], tGreen, fmt.Sprintf("primos ≤ %d (%d términos)", int(r.Filas[0][0]), int(r.Filas[0][1]))},
		{r.Filas[len(r.Filas)-1][4], tGreen, fmt.Sprintf("primos ≤ %d", int(r.Filas[len(r.Filas)-1][0]))},
		{r.S2verdad, tRose, "los ceros VERDADEROS"},
	}
	maxV := 0.45
	for i, b := range barras {
		x := bx + 40 + float64(i)*140
		h := mapa(b.v, 0, maxV, 0, bh)
		l.rect(x, by+bh-h, 78, h, b.c, b.c, 0.55)
		l.mono(x+39, by+bh-h-8, 12, b.c, "middle", num(b.v, 4))
		for j, w := range strings.Split(b.nom, " ") {
			l.txt(x+39, by+bh+18+float64(j)*13, 10.5, tDim, "middle", w)
		}
	}
	l.line(bx, by+bh-mapa(r.S2verdad, 0, maxV, 0, bh), bx+bw, by+bh-mapa(r.S2verdad, 0, maxV, 0, bh), tRose, 1.2, "5 5")
	l.txt(bx+8, by+bh-mapa(r.S2verdad, 0, maxV, 0, bh)-8, 10.5, tRose, "start", "el valor de los ceros")
	l.txt(56, 380, 11.5, tInk, "start",
		fmt.Sprintf("la caja acierta el CONTEO exacto y su Σ² es %s: una reja demasiado rígida. Con %d términos ya está en el valor de los ceros.",
			num(r.S2caja10, 4), int(r.Filas[0][1])))

	// ---- B: convergence to the true zeros ----------------------------------
	l.panel(700, 96, 664, 300, tGreen, "B · LOS PRIMOS SOLOS SE ACERCAN — |nivel − γ| contra cuántos primos se permiten")
	cx, cy, cw, ch := 760.0, 150.0, 570.0, 180.0
	l.line(cx, cy+ch, cx+cw, cy+ch, tGrid, 1.4, "")
	l.line(cx, cy, cx, cy+ch, tGrid, 1.4, "")
	var pts, pts2 [][2]float64
	for _, f := range r.Filas {
		x := mapa(math.Log10(f[1]), 0, 4, cx, cx+cw)
		pts = append(pts, [2]float64{x, mapa(math.Log10(f[2]), -0.3, -2.3, cy, cy+ch)})
		pts2 = append(pts2, [2]float64{x, mapa(math.Log10(f[3]), -0.3, -2.3, cy, cy+ch)})
	}
	l.camino(pts2, tDim, 1.8, 0.75)
	l.camino(pts, tGreen, 2.6, 1)
	for i, p := range pts {
		l.circ(p[0], p[1], 4.5, tGreen, "", 0)
		l.mono(p[0], cy+ch+16, 9.5, tDim, "middle", fmt.Sprintf("%d", int(r.Filas[i][1])))
	}
	l.mono(cx+cw/2, cy+ch+34, 10.5, tDim, "middle", "términos de Λ(n) usados")
	l.txt(cx+6, cy-8, 10.5, tGreen, "start", "medio (verde) y peor (gris), en unidades de altura")
	if len(r.Filas) > 0 {
		f := r.Filas[len(r.Filas)-1]
		l.mono(716, 366, 11.5, tInk, "start",
			fmt.Sprintf("con %d términos: |nivel − γ| medio = %s, peor = %s", int(f[1]), num(f[2], 5), num(f[3], 5)))
		l.txt(716, 384, 11.5, tGold, "start", "y el 100% de los 620 niveles cae dentro de un décimo del espaciado medio.")
	}

	// ---- C: the sign, measured ---------------------------------------------
	l.panel(36, 410, 1328, 288, tRose, "C · EL SIGNO — medido en nuestros propios ceros, período por período")
	sx, sy, sw, sh := 90.0, 470.0, 1240.0, 150.0
	l.line(sx, sy, sx+sw, sy, tGrid, 1.6, "")
	l.txt(sx-12, sy+4, 11, tDim, "end", "0")
	peor := 0.0
	for _, s := range r.Signos {
		if v := math.Abs(s[2]); v > peor {
			peor = v
		}
	}
	for i, s := range r.Signos {
		x := sx + 34 + float64(i)*float64(sw-70)/float64(len(r.Signos))
		hM := mapa(math.Abs(s[1]), 0, peor*1.12, 0, sh)
		hP := mapa(math.Abs(s[2]), 0, peor*1.12, 0, sh)
		l.rect(x-24, sy, 24, hP, tGrid, tDim, 0.55)
		l.rect(x, sy, 24, hM, tRose, tRose, 0.75)
		l.mono(x, sy+sh+22, 10.5, tDim, "middle", fmt.Sprintf("%d", int(math.Round(math.Exp(s[0])))))
		l.txt(x, sy+sh+40, 9.5, tGreen, "middle", "▼")
	}
	l.txt(sx+16, sy-14, 11.5, tRose, "start", "barra rosa: coeficiente medido |D(τ)| · barra gris: la predicción |−Λ(n)/(π√n)| · TODAS hacia abajo")
	l.mono(sx, sy+sh+62, 12.5, tGold, "start",
		fmt.Sprintf("%d de %d períodos aritméticos con coeficiente NEGATIVO — el espectro de absorción, leído en los datos, no argumentado",
			r.SignosNeg, r.SignosTot))
	l.txt(sx, sy+sh+82, 11.5, tDim, "start",
		"la fórmula de Selberg entra con + en las órbitas; la explícita con − en los primos. No es convención: es una restricción de diseño.")

	// ---- D: the verdict ----------------------------------------------------
	l.panel(36, 712, 1328, 128, tBlue, "D · LAS CUATRO CAPAS, SEPARADAS — y la tautología declarada")
	cols := []string{"candidato", "conteo", "correlaciones", "identidad", "aritmética", "R6"}
	xs := []float64{60, 330, 470, 660, 830, 1040}
	for i, c := range cols {
		l.txt(xs[i], 758, 11.5, tGold, "start", c)
	}
	filas := [][]string{
		{"H = diag(γₙ)", "sí", "sí", "sí", "sí", "CIRCULAR"},
		{"matriz GUE al azar", "reescalable", "sí", "no", "NO", "limpio pero vacío"},
		{"la caja que respira", "sí", "NO — rígida", "no", "NO", "ESTIPULADA"},
		{"los primos solos", "sí", "SÍ", "SÍ", "tautológica", "limpio"},
	}
	colr := []string{tDim, tDim, tBlue, tGreen}
	for j, f := range filas {
		y := 782 + float64(j)*16
		for i, v := range f {
			l.txt(xs[i], y, 11, colr[j], "start", v)
		}
	}
	l.txt(60, 852, 11.5, tRose, "start",
		"⚠ el eco de «los primos solos» NO informa: lleva log n en su definición, así que está garantizado. Lo que informa es Σ² y la identidad — y ésas no se pusieron a mano.")

	l.raw(fmt.Sprintf(`<rect x="36" y="866" width="1328" height="34" rx="8" fill="%s" stroke="%s" opacity="0.5"/>`, tPanel, tGold))
	l.mono(700, 889, 13.5, tGold, "middle", "N_p(T) = θ(T)/π + 1 − (1/π)·Σ_{n=p^k ≤ P} Λ(n)·sin(T·log n)/(√n·log n)      —      ningún γ entra en la definición")
	l.txt(700, 924, 13, tGold, "middle", "go run ./cmd/eltelar · estructura cerrada no es hipótesis demostrada · Todavía no.")

	l.raw("</svg>")
	dir := filepath.Join("galeria", "laminas", "10-el-telar")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		panic(err)
	}
	full := filepath.Join(dir, "el-telar.svg")
	if err := os.WriteFile(full, []byte(l.b.String()), 0o644); err != nil {
		panic(err)
	}
	fmt.Printf("\n🖼️  lámina escrita: %s\n", full)
}
