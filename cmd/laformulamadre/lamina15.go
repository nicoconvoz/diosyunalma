package main

// lamina15.go - the Phase XV plate: the zoom curve at the deepest depth with
// its crossing, the null world's crossing histogram, and the s*(M) table.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func e15(s string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;").Replace(s)
}

// E: zoom curve (deepest). nulos: null-world crossings. res: rows [M, s*, lo, hi, slope].
func dibujar15(E []float64, nulos []float64, res [][5]float64) {
	var b strings.Builder
	W, H := 1200.0, 800.0
	t := func(x, y, s float64, c, a, tx string) {
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Georgia,serif" fill="%s" text-anchor="%s">%s</text>`+"\n", x, y, s, c, a, e15(tx))
	}
	mono := func(x, y, s float64, c, a, tx string) {
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Consolas,monospace" fill="%s" text-anchor="%s">%s</text>`+"\n", x, y, s, c, a, e15(tx))
	}
	ln := func(x1, y1, x2, y2 float64, c string, w float64, da string) {
		dd := ""
		if da != "" {
			dd = fmt.Sprintf(` stroke-dasharray="%s"`, da)
		}
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="%.1f"%s/>`+"\n", x1, y1, x2, y2, c, w, dd)
	}
	pt := func(x, y, r float64, c string) {
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="%.1f" fill="%s"/>`+"\n", x, y, r, c)
	}
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">`+"\n", W, H, W, H)
	fmt.Fprintf(&b, `<rect width="%.0f" height="%.0f" fill="#0b1220"/>`+"\n", W, H)
	fmt.Fprintf(&b, `<rect x="14" y="14" width="%.0f" height="%.0f" rx="16" fill="none" stroke="#ffd98a" stroke-width="2" opacity="0.6"/>`+"\n", W-28, H-28)
	t(W/2, 48, 22, "#dce8f7", "middle", "🪞🎯 LOCALIZAR EL CRUCE — dónde el espejo cambia de signo")
	t(W/2, 74, 12, "#ffd98a", "middle", "Fase XV · zoom s ∈ [0,50, 1,30) en bins de 0,05 · E(s) = real − nulo · bootstrap y mundo nulo de cruces")

	gx, gy, gw, gh := 110.0, 110.0, 640.0, 420.0
	vMin, vMax := -0.30, 0.40
	mapy := func(v float64) float64 { return gy + gh - (v-vMin)/(vMax-vMin)*gh }
	mapx := func(s float64) float64 { return gx + (s-zLo)/(zHi-zLo)*gw }
	ln(gx, mapy(0), gx+gw, mapy(0), "#8fb4d9", 1.6, "")
	ln(gx, gy, gx, gy+gh, "#1d3a63", 1.2, "")
	for _, s := range []float64{0.6, 0.7, 0.8, 0.9, 1.0, 1.1, 1.2} {
		ln(mapx(s), gy, mapx(s), gy+gh, "#1d3a63", 0.6, "3 5")
		mono(mapx(s), gy+gh+16, 10, "#8fb4d9", "middle", fmt.Sprintf("%.1f", s))
	}
	var prev [2]float64
	for bi, v := range E {
		x, y := mapx(zLo+(float64(bi)+0.5)*zW), mapy(v)
		if bi > 0 {
			ln(prev[0], prev[1], x, y, "#7ee0c0", 2.4, "")
		}
		pt(x, y, 4.2, "#7ee0c0")
		prev = [2]float64{x, y}
	}
	sR := res[len(res)-1][1]
	ln(mapx(sR), gy, mapx(sR), gy+gh, "#ffd98a", 2, "6 4")
	t(mapx(sR), gy-8, 12, "#ffd98a", "middle", fmt.Sprintf("s* = %.3f", sR))
	t(gx+gw/2, gy+gh+38, 11.5, "#8fb4d9", "middle", "s = gap / espaciado medio local → · verde: E(s) a fondo (3474 ceros)")

	// null crossing histogram (below)
	hy := gy + gh + 70
	t(gx, hy-8, 12, "#ff9aa8", "start", "EL MUNDO NULO (su §6): dónde cruzan los 200 barajados — histograma de sus cruces")
	bins := make([]int, 16)
	for _, x := range nulos {
		if x >= zLo && x < zHi {
			bins[int((x-zLo)/zW)]++
		}
	}
	mx := 1
	for _, c := range bins {
		if c > mx {
			mx = c
		}
	}
	for bi, c := range bins {
		x := mapx(zLo + float64(bi)*zW)
		h := float64(c) / float64(mx) * 70
		fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#ff9aa8" opacity="0.55"/>`+"\n",
			x, hy+80-h, gw/16-2, h)
	}
	ln(mapx(sR), hy, mapx(sR), hy+82, "#ffd98a", 2, "6 4")
	ln(gx, hy+80, gx+gw, hy+80, "#1d3a63", 1, "")

	// right panel: the table and verdicts
	tx := 790.0
	t(tx, 120, 13.5, "#ffd98a", "start", "s*(M) — ¿converge o se mueve?")
	mono(tx, 148, 11.5, "#8fb4d9", "start", fmt.Sprintf("%6s %8s %17s %8s", "M", "s*", "boot 95%", "dE/ds"))
	for i, r := range res {
		mono(tx, 170+float64(i)*22, 11.5, "#dce8f7", "start",
			fmt.Sprintf("%6.0f %8.3f [%6.3f, %6.3f] %8.2f", r[0], r[1], r[2], r[3], r[4]))
	}
	y := 260.0
	for _, s := range []string{
		"· el cruce EXISTE en las tres",
		"  profundidades, dentro del zoom.",
		"· la pendiente es NEGATIVA y",
		"  estable: cruce suave, no meseta",
		"  ni salto (su §7).",
		"· el segundo control (centros",
		"  intactos, etiquetas permutadas)",
		"  APLANA la curva entera: toda la",
		"  dependencia en s es del",
		"  emparejamiento, no de las",
		"  marginales (su §10).",
		"· batería de refutación: sobrevive",
		"  al desdoblado global, y a las",
		"  dos mitades independientes de",
		"  los ceros (ver corrida).",
	} {
		t(tx, y, 11.5, "#dce8f7", "start", s)
		y += 19
	}
	t(tx, y+8, 11.5, "#ffd98a", "start", "s* NO se declara constante:")
	t(tx, y+26, 11.5, "#ffd98a", "start", "se declara localizado, con su")
	t(tx, y+44, 11.5, "#ffd98a", "start", "intervalo y su mundo nulo.")

	t(W/2, H-70, 12, "#8fb4d9", "middle", "«no persigas el 34,3σ — perseguí la forma que produce ese 34,3σ» — la auditora · zoom, nulos y reglas congelados antes de medir")
	t(W/2, H-44, 12.5, "#ffd98a", "middle", "go run ./cmd/laformulamadre fase15 · estructura cerrada no es hipótesis demostrada · Todavía no.")
	b.WriteString("</svg>\n")

	dir := filepath.Join("galeria", "laminas", "10-el-telar")
	os.MkdirAll(dir, 0o755)
	full := filepath.Join(dir, "el-cruce.svg")
	os.WriteFile(full, []byte(b.String()), 0o644)
	fmt.Printf("\n🖼️  lámina escrita: %s\n", full)
}
