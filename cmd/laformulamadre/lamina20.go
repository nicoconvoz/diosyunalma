package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func e20(s string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;").Replace(s)
}

func dibujar20(Ereal, E97, E9973, E99991 []float64, crReal float64) {
	var b strings.Builder
	W, H := 1240.0, 780.0
	t := func(x, y, s float64, c, a, tx string) {
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Georgia,serif" fill="%s" text-anchor="%s">%s</text>`+"\n", x, y, s, c, a, e20(tx))
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
	fmt.Fprintf(&b, `<rect x="14" y="14" width="%.0f" height="%.0f" rx="16" fill="none" stroke="#7ee0c0" stroke-width="2" opacity="0.6"/>`+"\n", W-28, H-28)
	t(W/2, 48, 22, "#dce8f7", "middle", "LA PRUEBA DE LA TRUNCACION - mas formula, sin mas ruido")
	t(W/2, 74, 12, "#ffd98a", "middle", "Fase XX - el campo puro N_liso + S = k-1/2, con la UNICA variable siendo el tope N de la suma - nada agregado")

	gx, gy, gw, gh := 110.0, 112.0, 680.0, 460.0
	vMin, vMax := -0.40, 0.45
	mapy := func(v float64) float64 { return gy + gh - (v-vMin)/(vMax-vMin)*gh }
	mapx := func(s float64) float64 { return gx + (s-zLo)/(zHi-zLo)*gw }
	ln(gx, mapy(0), gx+gw, mapy(0), "#8fb4d9", 1.6, "")
	for _, s := range []float64{0.6, 0.8, 1.0, 1.2} {
		ln(mapx(s), gy, mapx(s), gy+gh, "#1d3a63", 0.6, "3 5")
		t(mapx(s), gy+gh+16, 10, "#8fb4d9", "middle", fmt.Sprintf("%.1f", s))
	}
	dib := func(E []float64, col string, w float64) {
		var prev [2]float64
		for bi, v := range E {
			x, y := mapx(zLo+(float64(bi)+0.5)*zW), mapy(v)
			if bi > 0 {
				ln(prev[0], prev[1], x, y, col, w, "")
			}
			pt(x, y, 3.2, col)
			prev = [2]float64{x, y}
		}
	}
	dib(E97, "#8fb4d9", 1.6)
	dib(E9973, "#ff9aa8", 2.0)
	dib(E99991, "#ffd98a", 2.4)
	dib(Ereal, "#7ee0c0", 2.6)
	ln(mapx(crReal), gy, mapx(crReal), gy+gh, "#7ee0c0", 1.4, "6 4")
	t(gx+14, gy+18, 11.5, "#7ee0c0", "start", "verde: REAL (3474 ceros)")
	t(gx+14, gy+36, 11.5, "#ffd98a", "start", "dorado: N = 99991 (la formula mas completa)")
	t(gx+14, gy+54, 11.5, "#ff9aa8", "start", "rosa: N = 9973")
	t(gx+14, gy+72, 11.5, "#8fb4d9", "start", "celeste: N = 97 (donde arranco la campana)")
	t(gx+gw/2, gy+gh+38, 11.5, "#8fb4d9", "middle", "s = gap / espaciado medio local")

	tx := 820.0
	t(tx, 130, 13.5, "#ffd98a", "start", "LA ESCALERA DE LA CAMPANA")
	for i, s := range []string{
		"XVI: mover tamanos - un tercio",
		"XVII: mover posiciones - menos",
		"XVIII: juntos - el doble",
		"XIX: rigidez externa - muerta",
		"XX: quitar todo el ruido y",
		"completar la formula.",
	} {
		t(tx, 158+float64(i)*20, 12, "#dce8f7", "start", s)
	}
	t(tx, 300, 13.5, "#7ee0c0", "start", "LAS TRES CONVERGENCIAS BUSCADAS")
	for i, s := range []string{
		"amplitud: 0,285 hacia 0,346",
		"cruce: 0,94 hacia 0,862",
		"espaciados: quietos en 0,42",
		"mas T4 hacia 0,267.",
		"Los numeros por escalon estan",
		"en el acta. Y la regla del 16:",
		"si las tres convergen, FRENAR",
		"y AUDITAR - no interpretar.",
	} {
		t(tx, 328+float64(i)*20, 12, "#dce8f7", "start", s)
	}

	t(W/2, H-70, 12, "#8fb4d9", "middle", "esta fase no agrega una teoria: quita ruido y deja hablar a la formula - la auditora")
	t(W/2, H-44, 12.5, "#ffd98a", "middle", "go run ./cmd/laformulamadre fase20 - estructura cerrada no es hipotesis demostrada - Todavia no.")
	b.WriteString("</svg>\n")

	dir := filepath.Join("galeria", "laminas", "10-el-telar")
	os.MkdirAll(dir, 0o755)
	full := filepath.Join(dir, "la-truncacion.svg")
	os.WriteFile(full, []byte(b.String()), 0o644)
	fmt.Printf("\nlamina escrita: %s\n", full)
}
