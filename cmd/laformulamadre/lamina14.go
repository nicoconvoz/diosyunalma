package main

// lamina14.go - the Phase XIV plate: THE CURVE. Real signal against the
// empirical null, three depths, the sign crossing marked. Drawn from the very
// rows the run just measured.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func e14(s string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;").Replace(s)
}

// F rows: [tope, sLo, sHi, real, nullMean, nullSd, n]
func dibujar14(F [][7]float64) {
	var b strings.Builder
	W, H := 1240.0, 900.0
	t := func(x, y, s float64, c, a, tx string) {
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Georgia,serif" fill="%s" text-anchor="%s">%s</text>`+"\n", x, y, s, c, a, e14(tx))
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
	fmt.Fprintf(&b, `<rect x="14" y="14" width="%.0f" height="%.0f" rx="16" fill="none" stroke="#7fb2ff" stroke-width="2" opacity="0.6"/>`+"\n", W-28, H-28)
	t(W/2, 50, 22, "#dce8f7", "middle", "🪞📈 EL BARRIDO — la curva que los datos quisieron mostrar")
	t(W/2, 76, 12, "#ffd98a", "middle", "Fase XIV · señal(s) con s = gap/espaciado local · nulo empírico de 200 barajados · tres profundidades · malla congelada antes de correr")

	gx, gy, gw, gh := 120.0, 120.0, 760.0, 520.0
	vMin, vMax := -0.55, 0.85
	mapy := func(v float64) float64 { return gy + gh - (v-vMin)/(vMax-vMin)*gh }
	mapx := func(s float64) float64 {
		if s > 2.4 {
			s = 2.4
		}
		return gx + s/2.4*gw
	}
	ln(gx, gy, gx, gy+gh, "#1d3a63", 1.4, "")
	ln(gx, mapy(0), gx+gw, mapy(0), "#8fb4d9", 1.6, "")
	t(gx+gw+8, mapy(0)+4, 11, "#8fb4d9", "start", "cero")
	for _, v := range []float64{-0.4, -0.2, 0.2, 0.4, 0.6, 0.8} {
		ln(gx, mapy(v), gx+gw, mapy(v), "#1d3a63", 0.6, "3 5")
		t(gx-8, mapy(v)+4, 10, "#8fb4d9", "end", fmt.Sprintf("%+.1f", v))
	}
	for _, s := range []float64{0.5, 1.0, 1.5, 2.0} {
		ln(mapx(s), gy, mapx(s), gy+gh, "#1d3a63", 0.6, "3 5")
		t(mapx(s), gy+gh+18, 10.5, "#8fb4d9", "middle", fmt.Sprintf("%.1f", s))
	}
	t(gx+gw/2, gy+gh+38, 11.5, "#8fb4d9", "middle", "s = gap del par / espaciado medio LOCAL  →")
	t(gx-64, gy-16, 11.5, "#8fb4d9", "start", "−E medio en los 23 períodos log p (absorción hacia arriba)")

	cols := map[float64]string{1000: "#8fb4d9", 2000: "#7fb2ff", 4000: "#7ee0c0"}
	// null band of the deepest run
	for _, f := range F {
		if f[0] != 4000 {
			continue
		}
		xm := mapx((f[1] + f[2]) / 2)
		ln(xm, mapy(f[4]-2*f[5]), xm, mapy(f[4]+2*f[5]), "#5c4a2b", 7, "")
	}
	// curves per depth
	for _, tope := range []float64{1000, 2000, 4000} {
		var prev [2]float64
		primero := true
		for _, f := range F {
			if f[0] != tope {
				continue
			}
			xm, ym := mapx((f[1]+f[2])/2), mapy(f[3])
			if !primero {
				ln(prev[0], prev[1], xm, ym, cols[tope], 2.2, "")
			}
			pt(xm, ym, 4.4, cols[tope])
			prev, primero = [2]float64{xm, ym}, false
		}
	}
	// null mean of deepest as dashed line
	var prev [2]float64
	primero := true
	for _, f := range F {
		if f[0] != 4000 {
			continue
		}
		xm, ym := mapx((f[1]+f[2])/2), mapy(f[4])
		if !primero {
			ln(prev[0], prev[1], xm, ym, "#ffd98a", 1.6, "6 5")
		}
		prev, primero = [2]float64{xm, ym}, false
	}
	t(gx+14, gy+18, 11.5, "#7ee0c0", "start", "verde: 3474 ceros · azul: 1517 · celeste: 649")
	t(gx+14, gy+36, 11.5, "#ffd98a", "start", "punteada dorada: el NULO (gaps barajados, media de 200) — la trigonometría sola")
	t(gx+14, gy+54, 11.5, "#ff9aa8", "start", "la banda parda: ±2σ del nulo — casi invisible al lado de la señal")

	tx := 920.0
	t(tx, 130, 13.5, "#ffd98a", "start", "LO QUE LA CURVA DICE")
	for i, s := range []string{
		"· MONÓTONA en s: cuanto más",
		"  apretado el par, más absorción;",
		"  cuanto más ancho, más emisión.",
		"· CRUCE DE SIGNO estable en",
		"  s* ≈ 0,8–0,9, presente en las",
		"  tres profundidades.",
		"· MÍNIMO en s ∈ [1,3, 1,6].",
		"· El EXCEDENTE real−nulo tiene",
		"  magnitud ESTABLE con M",
		"  (Δ ≈ +0,43 arriba, −0,3 abajo):",
		"  crece la σ por datos, no la",
		"  física — la distinción de su §14.",
		"· Todos los bins del fondo a",
		"  |z| ≥ 6,3, hasta 34σ.",
	} {
		t(tx, 158+float64(i)*20, 12, "#dce8f7", "start", s)
	}
	t(tx, 470, 12.5, "#ff9aa8", "start", "HONESTIDAD: la física de fondo")
	t(tx, 490, 12, "#dce8f7", "start", "es Montgomery/Bogomolny–Keating")
	t(tx, 508, 12, "#dce8f7", "start", "en teoría; la FORMA —esta curva")
	t(tx, 526, 12, "#dce8f7", "start", "con este control— es del taller,")
	t(tx, 544, 12, "#dce8f7", "start", "sin búsqueda bibliográfica aún.")
	t(tx, 572, 12, "#ffd98a", "start", "s* NO se declara constante:")
	t(tx, 590, 12, "#ffd98a", "start", "deriva suave con M, pendiente")
	t(tx, 608, 12, "#ffd98a", "start", "de localización fina (su §13).")

	t(W/2, H-140, 12.5, "#7ee0c0", "middle", "el desdoblado LOCAL (s por espaciado local, no global) quedó declarado como mejora sobre F365 antes de correr:")
	t(W/2, H-120, 12.5, "#7ee0c0", "middle", "los ceros bajos tienen gaps anchos por ALTURA, no por soltura — sin corregirlo, la clase «ancho» era también la clase «bajo»")
	t(W/2, H-92, 12, "#8fb4d9", "middle", "malla, períodos, nulos y reglas de decisión congelados antes de la primera medición · semilla 20260821 · su regla de oro cumplida:")
	t(W/2, H-72, 12, "#8fb4d9", "middle", "no se buscó el resultado querido — se midió la curva que los datos quisieron mostrar")
	t(W/2, H-44, 12.5, "#ffd98a", "middle", "go run ./cmd/laformulamadre fase14 · estructura cerrada no es hipótesis demostrada · Todavía no.")
	b.WriteString("</svg>\n")

	dir := filepath.Join("galeria", "laminas", "10-el-telar")
	os.MkdirAll(dir, 0o755)
	full := filepath.Join(dir, "el-barrido.svg")
	os.WriteFile(full, []byte(b.String()), 0o644)
	fmt.Printf("\n🖼️  lámina escrita: %s\n", full)
}
