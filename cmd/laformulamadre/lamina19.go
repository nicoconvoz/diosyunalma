package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func e19(s string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;").Replace(s)
}

func dibujar19(Ereal, Erig, Epas, Efas []float64, crReal float64, crRC []float64) {
	var b strings.Builder
	W, H := 1240.0, 800.0
	t := func(x, y, s float64, c, a, tx string) {
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Georgia,serif" fill="%s" text-anchor="%s">%s</text>`+"\n", x, y, s, c, a, e19(tx))
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
	t(W/2, 48, 22, "#dce8f7", "middle", "LA RIGIDEZ - el escalon con dientes")
	t(W/2, 74, 12, "#ffd98a", "middle", "Fase XIX - autovalores GUE genuinos (tridiagonal de Hermite) clavados por la misma ecuacion de conteo - cero perillas")

	gx, gy, gw, gh := 110.0, 112.0, 680.0, 470.0
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
			pt(x, y, 3.4, col)
			prev = [2]float64{x, y}
		}
	}
	dib(Efas, "#8fb4d9", 1.4)
	dib(Epas, "#ff9aa8", 1.8)
	dib(Erig, "#ffd98a", 2.4)
	dib(Ereal, "#7ee0c0", 2.6)
	ln(mapx(crReal), gy, mapx(crReal), gy+gh, "#7ee0c0", 1.4, "6 4")
	t(gx+14, gy+18, 11.5, "#7ee0c0", "start", "verde: REAL (3474 ceros)")
	t(gx+14, gy+36, 11.5, "#ffd98a", "start", "dorado: RIGIDO + campo (la hipotesis de esta fase)")
	t(gx+14, gy+54, 11.5, "#ff9aa8", "start", "rosa: paseo + campo (Fase XVIII)")
	t(gx+14, gy+72, 11.5, "#8fb4d9", "start", "celeste: rigido + campo con FASE ROTA (control asesino)")
	t(gx+gw/2, gy+gh+38, 11.5, "#8fb4d9", "middle", "s = gap / espaciado medio local")

	tx := 820.0
	t(tx, 130, 13.5, "#ffd98a", "start", "EL MATERIAL NUEVO")
	for i, s := range []string{
		"ensamble tridiagonal de Hermite",
		"(beta = 2): autovalores GUE",
		"EXACTOS, desplegados por la",
		"semicircular, solo el interior.",
		"Varianza de conteo en L=10:",
		"paseo ~1,8 - rigido ~0,6 - la",
		"rigidez de largo alcance que el",
		"paseo browniano no tenia.",
	} {
		t(tx, 158+float64(i)*20, 12, "#dce8f7", "start", s)
	}
	t(tx, 350, 13.5, "#7ee0c0", "start", "EL JUICIO")
	t(tx, 378, 12, "#dce8f7", "start", fmt.Sprintf("cruce rigido: %s", listaF(crRC)))
	t(tx, 398, 12, "#dce8f7", "start", fmt.Sprintf("cruce real: %.3f", crReal))
	t(tx, 418, 12, "#dce8f7", "start", "amplitudes, T4, espaciados,")
	t(tx, 438, 12, "#dce8f7", "start", "residuo y veredicto: en el acta.")
	t(tx, 470, 12, "#ffd98a", "start", "la brujula de la auditora:")
	t(tx, 490, 12, "#ffd98a", "start", "no aumentar la fuerza -")
	t(tx, 510, 12, "#ffd98a", "start", "cambiar la dinamica.")

	t(W/2, H-70, 12, "#8fb4d9", "middle", "XVI tamanos - XVII posiciones - XVIII juntos - XIX el material: la pregunta ya no es la fuerza, es que recibe la fuerza")
	t(W/2, H-44, 12.5, "#ffd98a", "middle", "go run ./cmd/laformulamadre fase19 - estructura cerrada no es hipotesis demostrada - Todavia no.")
	b.WriteString("</svg>\n")

	dir := filepath.Join("galeria", "laminas", "10-el-telar")
	os.MkdirAll(dir, 0o755)
	full := filepath.Join(dir, "la-rigidez.svg")
	os.WriteFile(full, []byte(b.String()), 0o644)
	fmt.Printf("\nlamina escrita: %s\n", full)
}
