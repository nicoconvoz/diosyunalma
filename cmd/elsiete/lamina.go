package main

import (
	"fmt"
	"math"
	"strings"
)

// resultado is what the experiment measured, for the plate.
type resultado struct {
	hasta, primos, maxSin, maxCon int
	escapan, bases                []int
	artin                         [][3]float64 // tope, primes with 2 as primitive root, odd primes
}

// The palette of the house's plates.
const (
	fondo   = "#0d1017"
	panel   = "#131826"
	borde   = "#2a3040"
	tinta   = "#e8e2d0"
	suave   = "#c9c3b0"
	tenue   = "#8a8676"
	oro     = "#f2c14e"
	verde   = "#6ee7a0"
	coral   = "#e07a5f"
	apagado = "#5a6070"
)

// miles writes 664579 as 664.579, and coma writes 0.37391 as 0,37391.
func miles(n int) string {
	s := fmt.Sprint(n)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "." + s[i:]
	}
	return s
}

func coma(x float64, d int) string { return strings.Replace(fmt.Sprintf("%.*f", d, x), ".", ",", 1) }

// lamina draws the plate: the clock of 9 that the 2 walks, the prison of the
// 3, the captain's branches, the theorem, and what the experiment measured.
func lamina(r resultado) string {
	var b strings.Builder
	w := func(f string, a ...any) { fmt.Fprintf(&b, f, a...) }
	txt := func(x, y float64, tam int, color, ancla, s string) {
		w(`<text x="%.1f" y="%.1f" font-family="Georgia,serif" font-size="%d" fill="%s" text-anchor="%s">%s</text>`, x, y, tam, color, ancla, s)
	}
	caja := func(x, y, ww, hh float64) {
		w(`<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="%s" stroke="%s"/>`, x, y, ww, hh, panel, borde)
	}
	w(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1060 740" width="1060" height="740"><rect width="1060" height="740" fill="%s"/>`, fondo)
	w(`<defs><marker id="f" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto"><path d="M0,0 L10,5 L0,10 z" fill="%s"/></marker>`, oro)
	w(`<marker id="c" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto"><path d="M0,0 L10,5 L0,10 z" fill="%s"/></marker></defs>`, coral)
	txt(530, 40, 24, tinta, "middle", "EL SIETE DEL CAPITÁN")
	txt(530, 66, 12, tenue, "middle", "«2 = 2+2 = 4, 4+4 = 8, 8+8 = 16, 1+6 = 7 · todos dan 7» — y al 3 que se escapaba, le abrió el 12 en dos ramas")
	txt(530, 92, 13, verde, "middle", "duplicar una cifra, sumar las de dos: ¿por qué todo primo cae en el 7?")

	// the clock of 9
	caja(20, 112, 500, 360)
	txt(270, 138, 13, oro, "middle", "1 · EL RELOJ DEL 9")
	cx, cy, rad := 270.0, 300.0, 120.0
	pos := func(d int) (float64, float64) {
		a := -math.Pi/2 + 2*math.Pi*float64(d%9)/9
		return cx + rad*math.Cos(a), cy + rad*math.Sin(a)
	}
	w(`<circle cx="%.0f" cy="%.0f" r="%.0f" fill="none" stroke="%s" stroke-dasharray="3 5"/>`, cx, cy, rad, borde)
	flecha := func(a, z int, color, marca string) {
		x1, y1 := pos(a)
		x2, y2 := pos(z)
		dx, dy := x2-x1, y2-y1
		l := math.Hypot(dx, dy)
		w(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="2" marker-end="url(#%s)"/>`,
			x1+dx/l*20, y1+dy/l*20, x2-dx/l*22, y2-dy/l*22, color, marca)
	}
	for _, p := range [][2]int{{1, 2}, {2, 4}, {4, 8}, {8, 7}, {7, 5}, {5, 1}} {
		flecha(p[0], p[1], oro, "f")
	}
	flecha(3, 6, coral, "c")
	flecha(6, 3, coral, "c")
	for d := 1; d <= 9; d++ {
		x, y := pos(d)
		color, relleno := tinta, panel
		switch {
		case d == 7:
			color, relleno = fondo, verde
		case d%3 == 0:
			color = coral
		}
		w(`<circle cx="%.1f" cy="%.1f" r="17" fill="%s" stroke="%s"/>`, x, y, relleno, map[bool]string{true: coral, false: oro}[d%3 == 0])
		txt(x, y+6, 17, color, "middle", fmt.Sprint(d))
	}
	txt(270, 158, 12, suave, "middle", "sumar cifras = el resto por 9 · duplicar = r → 2r")
	txt(270, 450, 12, oro, "middle", "el 2 recorre las seis caras: 1 → 2 → 4 → 8 → 7 → 5 → 1")
	txt(270, 466, 11, coral, "middle", "el 3, el 6 y el 9: la cárcel de los múltiplos de 3")

	// the theorem
	caja(540, 112, 500, 360)
	txt(790, 138, 13, oro, "middle", "2 · EL TEOREMA (demostrado para todo número)")
	lineas := []struct{ c, s string }{
		{tinta, "(1) Sin ramas, n llega al 7 ⟺ 3 no divide a n."},
		{tinta, "(2) Todo primo llega al 7, salvo el 3."},
		{tinta, "(3) El 7 no es especial: el 1, 2, 4, 5 y 8, igual."},
		{verde, "(4) Con las ramas del capitán, TODO n ≥ 1 llega."},
		{tinta, "(5) En base b, el reloj está completo ⟺"},
		{tinta, "      el 2 es raíz primitiva módulo b − 1."},
		{tenue, ""},
		{suave, "Prueba: 10 = 9 + 1, así que sumar cifras guarda"},
		{suave, "el resto por 9, y duplicar lo multiplica por 2."},
		{suave, "Una cifra ES su resto. El 2 genera el reloj del 9;"},
		{suave, "el {0, 3, 6} se cierra sobre sí mismo (2·3 = 6, 2·6 ≡ 3)."},
		{coral, "No es de los primos: lo cumple todo no múltiplo de 3."},
	}
	for i, l := range lineas {
		txt(560, 170+float64(i)*24, 14, l.c, "start", l.s)
	}

	// the captain's branches
	caja(20, 490, 500, 230)
	txt(270, 516, 13, verde, "middle", "3 · LA LLAVE DEL CAPITÁN: LAS RAMAS")
	txt(270, 552, 16, tinta, "middle", "3 → 6 → 12")
	txt(270, 576, 12, tenue, "middle", "el 12 se abre en dos ramas, una por cifra:")
	txt(270, 600, 15, verde, "middle", "rama 1: 1 → 2 → 4 → 8 → 16 → 7")
	txt(270, 624, 15, verde, "middle", "rama 2: 2 → 4 → 8 → 16 → 7")
	txt(270, 652, 12, suave, "middle", "sumar las cifras del 12 da 3: respeta el resto, no sale de la cárcel")
	txt(270, 672, 12, suave, "middle", "abrirlas en 1 y 2 NO lo respeta: esa es la llave, y es suya")
	txt(270, 700, 12, tenue, "middle", "el 6 abre el 12, el 9 abre el 18 (1 y 8): nadie queda preso")

	// the experiment and the bridge
	caja(540, 490, 500, 230)
	txt(790, 516, 13, oro, "middle", fmt.Sprintf("4 · EL EXPERIMENTO, LITERAL, DE 1 A %s", miles(r.hasta)))
	txt(560, 546, 13, tinta, "start", fmt.Sprintf("0 excepciones al teorema · primos: %s, escapa solo %v", miles(r.primos), r.escapan))
	txt(560, 568, 13, tinta, "start", fmt.Sprintf("con ramas llegan todos (a lo sumo %d pasos; sin ramas, %d)", r.maxCon, r.maxSin))
	txt(560, 590, 13, tinta, "start", fmt.Sprintf("reloj completo en las bases %v", r.bases))
	txt(560, 618, 12, verde, "start", "el puente: ¿en qué parte de los primos p el 2 es raíz primitiva?")
	txt(560, 634, 11, tenue, "start", "la conjetura de Artin (abierta; Hooley 1967, si vale Riemann generalizada)")
	const constante = 0.3739558136192
	y0, x0, ancho := 704.0, 600.0, 360.0
	escala := func(v float64) float64 { return y0 - v*120 }
	w(`<line x1="%.0f" y1="%.1f" x2="%.0f" y2="%.1f" stroke="%s" stroke-dasharray="4 3"/>`, x0-10, escala(constante), x0+ancho, escala(constante), oro)
	txt(x0+ancho+4, escala(constante)+4, 11, oro, "start", coma(constante, 5))
	for i, a := range r.artin {
		parte := a[1] / a[2]
		x := x0 + float64(i)*(ancho/float64(len(r.artin))) + 10
		w(`<rect x="%.1f" y="%.1f" width="60" height="%.1f" fill="%s" opacity="0.8"/>`, x, escala(parte), y0-escala(parte), verde)
		txt(x+30, escala(parte)-4, 11, tinta, "middle", coma(parte, 5))
		txt(x+30, y0+12, 10, tenue, "middle", "hasta "+miles(int(a[0])))
	}
	w(`</svg>`)
	return b.String()
}
