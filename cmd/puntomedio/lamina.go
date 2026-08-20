package main

// lamina.go - the Theorem plate: the two rails (odds above, anchors below, gaps
// halved), the shared anchor of the twins, and the mod-6 wheel of the centers.

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

func esc13(s string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;").Replace(s)
}

func dibujar13() {
	var b strings.Builder
	W, H := 1240.0, 860.0
	t := func(x, y, s float64, c, a, tx string) {
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Georgia,serif" fill="%s" text-anchor="%s">%s</text>`+"\n", x, y, s, c, a, esc13(tx))
	}
	mono := func(x, y, s float64, c, a, tx string) {
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Consolas,monospace" fill="%s" text-anchor="%s">%s</text>`+"\n", x, y, s, c, a, esc13(tx))
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
	aro := func(x, y, r float64, c string) {
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="%.1f" fill="none" stroke="%s" stroke-width="2"/>`+"\n", x, y, r, c)
	}
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">`+"\n", W, H, W, H)
	fmt.Fprintf(&b, `<rect width="%.0f" height="%.0f" fill="#0b1220"/>`+"\n", W, H)
	fmt.Fprintf(&b, `<rect x="14" y="14" width="%.0f" height="%.0f" rx="16" fill="none" stroke="#ffd98a" stroke-width="2.5" opacity="0.7"/>`+"\n", W-28, H-28)
	t(W/2, 50, 23, "#dce8f7", "middle", "📐 EL TEOREMA DEL PUNTO MEDIO — Teorema 6 del taller")
	t(W/2, 76, 12.5, "#ffd98a", "middle", "(p+1)/2 = (q−1)/2 ⟺ q = p+2 · hallado A MANO por Jesús Nicolás Astorga · bautizado por el capitán")

	// --- the two rails: odds above, anchors below, gaps halved ---------------
	t(70, 118, 13, "#7ee0c0", "start", "LAS DOS VÍAS — arriba los impares con sus primos; abajo la vía de las anclas, con TODO gap dividido por 2")
	y1, y2 := 170.0, 250.0
	x0 := 90.0
	dx := (W - 180) / 26.0
	esP := map[int]bool{3: true, 5: true, 7: true, 11: true, 13: true, 17: true, 19: true, 23: true, 29: true}
	ln(x0, y1, x0+26*dx, y1, "#1d3a63", 1.6, "")
	ln(x0, y2, x0+26*dx, y2, "#1d3a63", 1.6, "")
	for n := 3; n <= 29; n++ {
		x := x0 + float64(n-3)*dx
		if n%2 == 1 {
			c := "#2a4a75"
			r := 3.0
			if esP[n] {
				c, r = "#7fb2ff", 5.5
			}
			pt(x, y1, r, c)
			mono(x, y1-12, 10, c, "middle", fmt.Sprintf("%d", n))
			// anchor image T(n)=(n-1)/2 on the lower rail
			xa := x0 + (float64(n-3)/2)*dx*2 // same horizontal scale compressed by 2
			_ = xa
		}
	}
	// lower rail: integers 1..14 = T(odds 3..29)
	for k := 1; k <= 13; k++ {
		x := x0 + float64(2*k+1-3)*dx
		c := "#2a4a75"
		r := 3.0
		if esP[2*k+1] {
			c, r = "#7ee0c0", 5.5
		}
		pt(x, y2, r, c)
		mono(x, y2+16, 10, c, "middle", fmt.Sprintf("%d", k))
		ln(x0+float64(2*k+1-3)*dx, y1+8, x, y2-8, "#1d3a63", 0.8, "2 4")
	}
	// twins 11,13 -> consecutive 5,6 ; shared anchor 6
	for _, par := range [][2]int{{11, 13}, {17, 19}} {
		xm := x0 + (float64(par[0]+par[1])/2-3)*dx
		aro(xm, y1, 9, "#ffd98a")
		mono(xm, y1-24, 10, "#ffd98a", "middle", fmt.Sprintf("m=%d", (par[0]+par[1])/2))
		xc := x0 + (float64(par[0]+1)-3)*dx
		aro(xc, y2, 9, "#ffd98a")
		mono(xc, y2+32, 10, "#ffd98a", "middle", fmt.Sprintf("c=%d", (par[0]+1)/2))
	}
	t(70, y2+58, 12, "#dce8f7", "start", "los gemelos (11,13) y (17,19) se vuelven enteros CONSECUTIVOS abajo, y el ancla compartida c = m/2 queda anillada.")
	t(70, y2+78, 12, "#ffd98a", "start", "T(n) = (n−1)/2 es la biyección impares→enteros: divide todo gap por 2. Los ±1/2 son esa coordenada — nada más, nada menos.")

	// --- the identity strip ----------------------------------------------------
	yI := y2 + 116
	mono(W/2, yI, 14, "#7ee0c0", "middle", "FORMA GENERAL:  (p,q) = (m−g/2, m+g/2)   ·   a⁻(q) − a⁺(p) = g/2 − 1   ·   3 | m ⟺ 6 ∤ g")
	t(W/2, yI+22, 11.5, "#8fb4d9", "middle", "verificadas con CERO fallas: ~25 millones de pares impares · 17 982 pares consecutivos de primos · TODOS los 304 590 pares en (3, 6000]")

	// --- the mod-6 wheel of centers -------------------------------------------
	yW := yI + 60
	t(70, yW, 13, "#7fb2ff", "start", "LA RUEDA DE LOS CENTROS — la clase del gap se LEE en el residuo del centro mod 6 (período 6 en g)")
	cx, cy, R := 240.0, yW+150, 105.0
	cols := map[int]string{0: "#ffd98a", 3: "#ff9aa8", 2: "#7fb2ff", 4: "#7fb2ff", 1: "#7ee0c0", 5: "#7ee0c0"}
	for k := 0; k < 6; k++ {
		a := -math.Pi/2 + float64(k)*math.Pi/3
		x, y := cx+R*math.Cos(a), cy+R*math.Sin(a)
		pt(x, y, 17, cols[k])
		mono(x, y+5, 14, "#0b1220", "middle", fmt.Sprintf("%d", k))
	}
	aro(cx, cy, R+30, "#1d3a63")
	tabla := []struct {
		g   int
		res string
		col string
		n   int
	}{
		{2, "{0}", "#ffd98a", 2159}, {4, "{3}", "#ff9aa8", 2134}, {6, "{2,4}", "#7fb2ff", 3455},
		{8, "{3}", "#ff9aa8", 1391}, {10, "{0}", "#ffd98a", 1705}, {12, "{1,5}", "#7ee0c0", 1821},
		{14, "{0}", "#ffd98a", 959}, {16, "{3}", "#ff9aa8", 637}, {18, "{2,4}", "#7fb2ff", 1035},
	}
	tx, ty := 450.0, yW+40.0
	mono(tx, ty, 12, "#8fb4d9", "start", fmt.Sprintf("%4s %10s %10s", "g", "m mod 6", "pares"))
	for i, r := range tabla {
		mono(tx, ty+20+float64(i)*20, 12, r.col, "start", fmt.Sprintf("%4d %10s %10d", r.g, r.res, r.n))
	}
	t(760, ty+18, 12.5, "#dce8f7", "start", "los gemelos viven en m ≡ 0 (mod 6): la red 6ℤ.")
	t(760, ty+40, 12.5, "#dce8f7", "start", "g=4 en los múltiplos impares de 3 · g=6 esquiva")
	t(760, ty+62, 12.5, "#dce8f7", "start", "el 3 · y la firma se repite con período 6 en g.")
	t(760, ty+92, 12.5, "#ffd98a", "start", "Consecuencia sin simulación: todo par gemelo")
	t(760, ty+114, 12.5, "#ffd98a", "start", "mayor que (3,5) tiene m ≡ 0 (mod 6) y c ≡ 0 (mod 3).")
	t(760, ty+136, 12.5, "#ffd98a", "start", "Sin excepción posible: es la ley del centro.")
	t(760, ty+166, 11.5, "#ff9aa8", "start", "Honestidad: la ley es criba mod 3 CLÁSICA — verdadera y")
	t(760, ty+186, 11.5, "#ff9aa8", "start", "demostrada, no nueva para el mundo; nueva como organización")
	t(760, ty+206, 11.5, "#ff9aa8", "start", "para el taller. Y nada de esto afirma nada sobre RH.")

	t(W/2, H-64, 12, "#8fb4d9", "middle", "matemática pura: ningún operador, ningún espectro, ningún cero — la auditora prohibió correr y se obedeció")
	t(W/2, H-40, 12.5, "#ffd98a", "middle", "go run ./cmd/puntomedio · estructura cerrada no es hipótesis demostrada · Todavía no.")
	b.WriteString("</svg>\n")

	dir := filepath.Join("galeria", "laminas", "10-el-telar")
	os.MkdirAll(dir, 0o755)
	full := filepath.Join(dir, "el-punto-medio.svg")
	os.WriteFile(full, []byte(b.String()), 0o644)
	fmt.Printf("\n🖼️  lámina escrita: %s\n", full)
}
