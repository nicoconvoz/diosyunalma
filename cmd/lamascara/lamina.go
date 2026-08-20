package main

// lamina.go - the Phase XII plate: the captain's identity drawn as anchors, the
// exhaustive verification, and the honest verdict board of the sparse mask.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func e12(s string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;").Replace(s)
}

func dibujar12() {
	var b strings.Builder
	W, H := 1200.0, 820.0
	t := func(x, y, s float64, c, a, tx string) {
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Georgia,serif" fill="%s" text-anchor="%s">%s</text>`+"\n", x, y, s, c, a, e12(tx))
	}
	mono := func(x, y, s float64, c, a, tx string) {
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Consolas,monospace" fill="%s" text-anchor="%s">%s</text>`+"\n", x, y, s, c, a, e12(tx))
	}
	ln := func(x1, y1, x2, y2 float64, c string, w float64) {
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="%.1f"/>`+"\n", x1, y1, x2, y2, c, w)
	}
	pt := func(x, y, r float64, c string) {
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="%.1f" fill="%s"/>`+"\n", x, y, r, c)
	}
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">`+"\n", W, H, W, H)
	fmt.Fprintf(&b, `<rect width="%.0f" height="%.0f" fill="#0b1220"/>`+"\n", W, H)
	fmt.Fprintf(&b, `<rect x="14" y="14" width="%.0f" height="%.0f" rx="16" fill="none" stroke="#7ee0c0" stroke-width="2" opacity="0.55"/>`+"\n", W-28, H-28)
	t(W/2, 50, 22, "#dce8f7", "middle", "🧵🌌 LA MÁSCARA — la geometría de los gemelos, del papel al operador")
	t(W/2, 76, 12, "#ffd98a", "middle", "la relación (p+1)/2 = (q−1)/2 la encontró Jesús Nicolás Astorga A MANO, haciendo cuentas con primos")

	// the anchors, drawn: 11-12-13 shared anchor; 13..17 gap example
	y0 := 150.0
	t(90, y0-28, 13, "#7ee0c0", "start", "LAS ANCLAS: a±(p) = (p±1)/2 — y el diccionario entero se vuelve identidad")
	esc := func(v float64) float64 { return 80 + (v-4)*34 }
	ln(esc(4), y0, esc(17), y0, "#1d3a63", 1.4)
	for v := 4; v <= 17; v++ {
		pt(esc(float64(v)), y0, 2.2, "#1d3a63")
		mono(esc(float64(v)), y0+18, 9.5, "#8fb4d9", "middle", fmt.Sprintf("%d", v))
	}
	// p=11: anchors 5,6 ; q=13: anchors 6,7 -> shared 6
	for _, a := range []struct {
		v float64
		c string
		n string
	}{{5, "#7fb2ff", "a−(11)"}, {6, "#ffd98a", "a+(11) = a−(13)"}, {7, "#7fb2ff", "a+(13)"}} {
		pt(esc(a.v), y0, 5.5, a.c)
		t(esc(a.v), y0-12, 10, a.c, "middle", a.n)
	}
	// p=13,q=17 (g=4): a+(13)=7, a-(17)=8 adjacent
	for _, a := range []struct {
		v float64
		c string
		n string
	}{{8, "#ff9aa8", "a−(17)"}, {9, "#7fb2ff", "a+(17)"}} {
		pt(esc(a.v), y0, 5.5, a.c)
		t(esc(a.v), y0+34, 10, a.c, "middle", a.n)
	}
	t(90, y0+62, 12, "#dce8f7", "start", "g=2 comparte el ancla (oro) · g=4 anclas vecinas que se tocan · g>4 deja (g−4)/2 enteros de hueco")
	mono(90, y0+86, 12.5, "#7ee0c0", "start", "VERIFICADO sobre los 9590 pares de primos consecutivos ≤ 100000: 1224/1224 · 1215 · 7151 · CERO fallas")
	t(90, y0+108, 12.5, "#ffd98a", "start", "⟹ la igualdad del capitán es un TEOREMA chico y verdadero: identidad aritmética, no coincidencia visual.")

	// verdict board
	y := y0 + 160
	t(90, y, 13, "#7fb2ff", "start", "LA MÁSCARA RALA (regla congelada antes de todo espectro): signo −1 donde |p_i − p_j| = g")
	y += 26
	mono(90, y, 12, "#8fb4d9", "start", fmt.Sprintf("%-22s %8s %8s %10s %12s %14s", "clase", "enlaces", "vivos", "Σ²(10)", "vs azar", "vs permutada"))
	filas := []string{
		fmt.Sprintf("%-22s %8d %8d %10.3f %12s %14s", "gemelos g=2", 62, 181, 15.497, "+3,50σ", "+0,24σ"),
		fmt.Sprintf("%-22s %8d %8d %10.3f %12s %14s", "g=4", 66, 185, 15.491, "+0,03σ", "+0,80σ"),
		fmt.Sprintf("%-22s %8d %8d %10.3f %12s %14s", "g=6", 114, 185, 15.738, "+0,06σ", "+0,15σ"),
		fmt.Sprintf("%-22s %8d %8d %10.3f", "S0 sin máscara", 0, 187, 18.335),
	}
	for _, f := range filas {
		y += 22
		mono(90, y, 12, "#dce8f7", "start", f)
	}
	y += 34
	t(90, y, 12.5, "#ffd98a", "start", "El matiz que la auditora anticipó textualmente: los gemelos le ganan al azar crudo (+3,5σ) pero NO a su")
	y += 20
	t(90, y, 12.5, "#ffd98a", "start", "permutada por distancia (+0,24σ). Lo que aporta la máscara es EN QUÉ DISTANCIAS caen los gemelos —")
	y += 20
	t(90, y, 12.5, "#ffd98a", "start", "log p casi iguales ⟹ armónicos vecinos, k ∈ [3,80] — y esa distribución, copiada sin aritmética, rinde igual.")
	y += 30
	t(90, y, 12.5, "#ff9aa8", "start", "VEREDICTO: nulo limpio A ESTE TAMAÑO. Con 400 modos la máscara marca 62 de 40740 enlaces (0,15%):")
	y += 20
	t(90, y, 12.5, "#ff9aa8", "start", "existe pero casi no toca al operador. La pregunta no muere: queda «NO DECIDIBLE a N=400» — decidirla")
	y += 20
	t(90, y, 12.5, "#ff9aa8", "start", "pide un medio más grande (más gemelos adentro), no otra regla. El eco no se corrió, por regla.")

	t(W/2, H-64, 12, "#8fb4d9", "middle", "el hallazgo es del capitán, a mano · la formalización por anclas y los controles emparejados, del taller · R6 absoluto")
	t(W/2, H-40, 12.5, "#ffd98a", "middle", "go run ./cmd/lamascara · estructura cerrada no es hipótesis demostrada · Todavía no.")
	b.WriteString("</svg>\n")

	dir := filepath.Join("galeria", "laminas", "10-el-telar")
	os.MkdirAll(dir, 0o755)
	full := filepath.Join(dir, "la-mascara.svg")
	os.WriteFile(full, []byte(b.String()), 0o644)
	fmt.Printf("\n🖼️  lámina escrita: %s\n", full)
}
