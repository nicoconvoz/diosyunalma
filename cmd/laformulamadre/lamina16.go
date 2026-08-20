package main

// lamina16.go - the Phase XVI plate: observed against the two known models
// (pure GUE flat, GUE+explicit-formula), the selection terms, the residual.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func e16(s string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;").Replace(s)
}

func dibujar16(Eobs, Em, Eg, T2, T4 []float64, crObs float64, crM []float64) {
	var b strings.Builder
	W, H := 1240.0, 830.0
	t := func(x, y, s float64, c, a, tx string) {
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Georgia,serif" fill="%s" text-anchor="%s">%s</text>`+"\n", x, y, s, c, a, e16(tx))
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
	t(W/2, 48, 22, "#dce8f7", "middle", "🪞🧠 ¿POR QUÉ CRUZA EL ESPEJO? — derivar primero, comparar después")
	t(W/2, 74, 12, "#ffd98a", "middle", "Fase XVI · la identidad de cuatro términos · GUE puro · GUE + fórmula explícita · el residuo")

	gx, gy, gw, gh := 110.0, 112.0, 660.0, 440.0
	vMin, vMax := -0.35, 0.45
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
			pt(x, y, 3.6, col)
			prev = [2]float64{x, y}
		}
	}
	dib(Eg, "#8fb4d9", 1.4)   // pure GUE, ~flat
	dib(Em, "#ffd98a", 2.2)   // mechanism model
	dib(Eobs, "#7ee0c0", 2.6) // observed
	ln(mapx(crObs), gy, mapx(crObs), gy+gh, "#7ee0c0", 1.6, "6 4")
	t(gx+14, gy+18, 11.5, "#7ee0c0", "start", "verde: E(s) OBSERVADO en los 3474 ceros reales")
	t(gx+14, gy+36, 11.5, "#ffd98a", "start", "dorado: GUE + FÓRMULA EXPLÍCITA (densidad modulada por los primos, gap Wigner)")
	t(gx+14, gy+54, 11.5, "#8fb4d9", "start", "celeste: GUE PURO — plano: sin primos no hay curva")
	t(gx+gw/2, gy+gh+38, 11.5, "#8fb4d9", "middle", "s = gap / espaciado medio local →")

	tx := 800.0
	t(tx, 120, 13.5, "#ffd98a", "start", "LA IDENTIDAD (su §2–4)")
	for i, s := range []string{
		"E(s) = T1 + T2 + T3 + T4, exacta:",
		"cos(mT) = cosγT·cos(gT/2)",
		"        − sinγT·sin(gT/2)",
		"T2+T4 = SELECCIÓN × factor:",
		"«el cero de un par apretado",
		"vive en la cresta de la onda",
		"de los primos» — y domina.",
		"T1+T3 = covarianzas internas,",
		"chicas (ver corrida).",
	} {
		t(tx, 148+float64(i)*20, 12, "#dce8f7", "start", s)
	}
	t(tx, 340, 13.5, "#7ee0c0", "start", "EL VEREDICTO (su §10), honesto")
	crm := "—"
	if len(crM) > 0 {
		crm = fmt.Sprintf("%.2f/%.2f/%.2f", crM[0], crM[1], crM[2])
	}
	for i, s := range []string{
		"GUE puro: PLANO — sin primos",
		"no hay curva. Descartado.",
		"GUE + fórmula (densidad",
		"modulada): produce el SIGNO y",
		"la caída, pero SE QUEDA CORTO:",
		"amplitud ⅓ de la real, cruce en",
		fmt.Sprintf("%s (real: %.3f),", crm, crObs),
		"pendiente −0,6 (real: −0,92).",
		"El residuo conserva la mayor",
		"parte de la estructura.",
		"⟹ Y LA IDENTIDAD NOMBRA LO QUE",
		"FALTA: domina T4, el acople",
		"seno·seno — los ceros no sólo se",
		"APRIETAN donde la onda encresta:",
		"se CORREN con ella. El modelo",
		"movía gaps; falta mover las",
		"POSICIONES coherentemente.",
	} {
		t(tx, 368+float64(i)*20, 12, "#dce8f7", "start", s)
	}

	t(W/2, H-96, 12, "#ff9aa8", "middle", "«y si una explicación conocida mata la novedad, decirlo» — se dice lo contrario, con números: la explicación conocida NO mata la curva todavía —")
	t(W/2, H-76, 12, "#ff9aa8", "middle", "explica un tercio; los dos tercios restantes tienen nombre y domicilio: el término T4, el corrimiento coherente de los ceros con la onda de los primos")
	t(W/2, H-44, 12.5, "#ffd98a", "middle", "go run ./cmd/laformulamadre fase16 · estructura cerrada no es hipótesis demostrada · Todavía no.")
	b.WriteString("</svg>\n")

	dir := filepath.Join("galeria", "laminas", "10-el-telar")
	os.MkdirAll(dir, 0o755)
	full := filepath.Join(dir, "el-porque.svg")
	os.WriteFile(full, []byte(b.String()), 0o644)
	fmt.Printf("\n🖼️  lámina escrita: %s\n", full)
}
