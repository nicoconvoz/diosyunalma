// Command eldibujodelcapitan renders the diagram the captain drew on paper, as
// HE described it after correcting the first reading:
//
//	"son puntos que se cruzan en R5 y que forman como una especie de reloj de
//	 arena a medio terminar. R6 tiene la direccion opuesta pero es como un
//	 reflejo, y R1 tiene la mayor distancia de separacion y va disminuyendo
//	 hasta R5, que R5 es el mismo punto de partida el mismo punto de fin, y ahi
//	 se vuelve a abrir en R6 y toma un poquito de distancia espejado a R4."
//
// So it is not a star and not a chain: it is a CROSSING. The separation shrinks
// to exactly zero at R5 - start point and end point are the same - and then
// reopens on the other side. The reversal at R6 is not a chosen sign: it is the
// consequence of having crossed. The roles swapped.
//
// LEFT panel: his figure, drawn as he described it, plus the dashed continuation
// the shape implies and he has not drawn yet (the second half of the hourglass).
// RIGHT panel: how this workshop reads it, marked as a reading, not as his claim.
//
// Reproduce: go run ./cmd/eldibujodelcapitan
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	bg    = "#0b1526"
	panel = "#0d1830"
	ink   = "#dce8f7"
	dim   = "#8fb4d9"
	gold  = "#ffd98a"
	green = "#7ee0c0"
	blue  = "#7fb2ff"
	rose  = "#ff9aa8"
	grid  = "#1d3a63"
)

type hoja struct{ b strings.Builder }

func (l *hoja) raw(s string) { l.b.WriteString(s + "\n") }

func esc(s string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;").Replace(s)
}

func (l *hoja) txt(x, y, s float64, c, a, t string) {
	l.raw(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Georgia,serif" fill="%s" text-anchor="%s">%s</text>`, x, y, s, c, a, esc(t)))
}

func (l *hoja) ital(x, y, s float64, c, a, t string) {
	l.raw(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Georgia,serif" font-style="italic" font-weight="bold" fill="%s" text-anchor="%s">%s</text>`, x, y, s, c, a, esc(t)))
}

func (l *hoja) mono(x, y, s float64, c, a, t string) {
	l.raw(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Consolas,monospace" fill="%s" text-anchor="%s">%s</text>`, x, y, s, c, a, esc(t)))
}

func (l *hoja) panelBox(x, y, w, h float64, c, t string) {
	l.raw(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="10" fill="%s" stroke="%s"/>`, x, y, w, h, panel, c))
	if t != "" {
		l.txt(x+16, y+26, 15, c, "start", t)
	}
}

func (l *hoja) circ(x, y, r float64, f, st string, sw float64) {
	s := ""
	if st != "" {
		s = fmt.Sprintf(` stroke="%s" stroke-width="%.2f"`, st, sw)
	}
	l.raw(fmt.Sprintf(`<circle cx="%.2f" cy="%.2f" r="%.2f" fill="%s"%s/>`, x, y, r, f, s))
}

func (l *hoja) line(x1, y1, x2, y2 float64, c string, w float64, d string) {
	dd := ""
	if d != "" {
		dd = fmt.Sprintf(` stroke-dasharray="%s"`, d)
	}
	l.raw(fmt.Sprintf(`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="%.2f"%s/>`, x1, y1, x2, y2, c, w, dd))
}

func (l *hoja) camino(pts [][2]float64, c string, w float64, d string) {
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
	dd := ""
	if d != "" {
		dd = fmt.Sprintf(` stroke-dasharray="%s"`, d)
	}
	l.raw(fmt.Sprintf(`<path d="%s" fill="none" stroke="%s" stroke-width="%.2f"%s/>`, sb.String(), c, w, dd))
}

// flecha draws an arrow from x1 to x2 at height y; it points where x2 is.
func (l *hoja) flecha(x1, x2, y float64, c string, w float64) {
	if x1 == x2 {
		return
	}
	dir := 1.0
	if x2 < x1 {
		dir = -1
	}
	hl := 10.0
	l.line(x1, y, x2-dir*hl*0.7, y, c, w, "")
	l.raw(fmt.Sprintf(`<path d="M%.2f %.2f L%.2f %.2f L%.2f %.2f Z" fill="%s"/>`,
		x2, y, x2-dir*hl, y-5, x2-dir*hl, y+5, c))
}

func main() {
	l := &hoja{}
	const W, H = 1400.0, 950.0
	l.raw(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">`, W, H, W, H))
	l.raw(fmt.Sprintf(`<rect width="%.0f" height="%.0f" fill="%s"/>`, W, H, bg))
	l.raw(fmt.Sprintf(`<rect x="14" y="14" width="%.0f" height="%.0f" rx="16" fill="none" stroke="%s" stroke-width="2" opacity="0.55"/>`, W-28, H-28, green))
	l.txt(700, 50, 26, ink, "middle", "⏳ EL RELOJ DE ARENA DEL CAPITÁN — la separación se cierra en R5 y se abre espejada")
	l.txt(700, 78, 13.5, gold, "middle",
		"«R5 es el mismo punto de partida y el mismo punto de fin, y ahí se vuelve a abrir en R6, espejado a R4» — Jesús Nicolás Astorga")

	// =======================================================================
	// LEFT — the figure as he described it
	// =======================================================================
	l.panelBox(36, 104, 700, 706, gold, "SU FIGURA — el cruce, y el reloj de arena a medio terminar")

	eje := 400.0 // the axis: where the two points meet
	yTop, paso := 190.0, 74.0
	// separations, in his order and with his description: R1 the widest, down to
	// exactly zero at R5, then reopening mirrored at R6 (the mirror of R4)
	filas := []struct {
		nom  string
		sep  float64 // signed: positive before the crossing, negative after
		suyo bool    // drawn by him, or the continuation his shape implies
	}{
		{"R1", 260, true},
		{"R2", 196, true},
		{"R3", 134, true},
		{"R4", 74, true},
		{"R5", 0, true},
		{"R6", -74, true},
		{"R7", -134, false},
		{"R8", -196, false},
		{"R9", -260, false},
	}

	// the two paths, so the crossing is visible as a shape and not as a claim
	var izq, der [][2]float64
	for i, f := range filas {
		y := yTop + float64(i)*paso
		izq = append(izq, [2]float64{eje - f.sep/2, y})
		der = append(der, [2]float64{eje + f.sep/2, y})
	}
	// his half solid, the implied half dashed
	l.camino(izq[:6], blue, 2.2, "")
	l.camino(der[:6], green, 2.2, "")
	l.camino(izq[5:], blue, 1.6, "7 6")
	l.camino(der[5:], green, 1.6, "7 6")
	l.line(eje, yTop-26, eje, yTop+8*paso+26, grid, 1.3, "4 7")

	for i, f := range filas {
		y := yTop + float64(i)*paso
		xa, xb := eje-f.sep/2, eje+f.sep/2
		colFl := blue
		if f.sep < 0 {
			colFl = rose
		}
		op := 1.0
		if !f.suyo {
			op = 0.45
		}
		l.raw(fmt.Sprintf(`<g opacity="%.2f">`, op))
		if f.sep == 0 {
			// R5: one single point - start and end are the same
			l.circ(eje, y, 13, gold, ink, 2.4)
			l.txt(eje+34, y+6, 14.5, gold, "start", "R5 · el mismo punto de partida y de fin")
			l.txt(eje-34, y+6, 13, gold, "end", "SEPARACIÓN CERO")
		} else {
			l.circ(xa, y, 8, ink, "", 0)
			l.circ(xb, y, 8, ink, "", 0)
			l.flecha(xa+12, xb-12, y, colFl, 2.4)
			et := f.nom
			if !f.suyo {
				et = f.nom + " (lo que la forma implica)"
			}
			l.ital(eje+150, y+6, 15.5, ink, "start", et)
		}
		l.mono(56, y+5, 12, dim, "start", fmt.Sprintf("sep %+d", int(f.sep)))
		l.raw(`</g>`)
	}
	l.txt(56, 768, 13, blue, "start", "azul y verde: los dos puntos. Se acercan, se TOCAN en R5, y siguen — ahora cada uno del otro lado.")
	l.txt(56, 790, 13, rose, "start", "la flecha de R6 apunta al revés porque los roles se INTERCAMBIARON, no porque se haya elegido un signo.")

	// =======================================================================
	// RIGHT — the reading
	// =======================================================================
	l.panelBox(756, 104, 608, 706, green, "CÓMO LO LEEMOS ACÁ (lectura nuestra, no afirmación suya)")

	l.txt(778, 150, 14.5, ink, "start", "1 · Eso es un NODO.")
	l.txt(778, 174, 13, dim, "start", "Una cantidad que se hace exactamente cero y CAMBIA DE SIGNO.")
	l.txt(778, 194, 13, dim, "start", "El vuelco de R6 no es una elección: es lo que pasa después de un cero.")

	// draw the signed quantity as a curve through zero
	gx, gy, gw, gh := 800.0, 232.0, 520.0, 150.0
	l.line(gx, gy+gh/2, gx+gw, gy+gh/2, grid, 1.4, "")
	l.txt(gx-8, gy+gh/2+5, 12, dim, "end", "0")
	var curva [][2]float64
	for i, f := range filas {
		x := gx + float64(i)*gw/8
		curva = append(curva, [2]float64{x, gy + gh/2 - f.sep/260*(gh/2-8)})
	}
	l.camino(curva[:6], gold, 2.6, "")
	l.camino(curva[5:], gold, 1.8, "7 6")
	for i, f := range filas {
		p := curva[i]
		c := gold
		if !f.suyo {
			c = dim
		}
		l.circ(p[0], p[1], 4.6, c, "", 0)
		l.mono(p[0], gy+gh+18, 10.5, dim, "middle", f.nom)
	}
	l.circ(curva[4][0], curva[4][1], 8, "none", rose, 2.4)
	l.txt(curva[4][0], gy-6, 12, rose, "middle", "el nodo")

	l.txt(778, 428, 14.5, ink, "start", "2 · Y su forma hace una PREDICCIÓN.")
	l.txt(778, 452, 13, dim, "start", "Si de verdad es un reloj de arena, la mitad que falta no es libre:")
	l.mono(798, 478, 13, gold, "start", "R7 espeja a R3   ·   R8 espeja a R2   ·   R9 espeja a R1")
	l.txt(778, 504, 13, dim, "start", "Eso se puede dibujar, y se puede desmentir. Un dibujo que predice")
	l.txt(778, 524, 13, dim, "start", "es mucho más que un dibujo.")

	l.txt(778, 566, 14.5, ink, "start", "3 · Por qué nos importa acá.")
	l.txt(778, 590, 13, dim, "start", "La Fase VI dejó pidiendo un núcleo que dé ALCANCE sin diluir la fuerza.")
	l.txt(778, 610, 13, dim, "start", "Un núcleo con un NODO hace justo eso: las contribuciones lejanas se")
	l.txt(778, 630, 13, dim, "start", "CANCELAN en vez de sumarse, así que el alcance no cuesta empuje.")

	l.raw(fmt.Sprintf(`<rect x="778" y="654" width="564" height="132" rx="9" fill="%s" stroke="%s" opacity="0.45"/>`, panel, blue))
	l.mono(1060, 682, 13, gold, "middle", "c_k  con  c_5 = 0  y  c_k>5 de signo opuesto")
	l.txt(1060, 708, 12.5, ink, "middle", "un núcleo de alcance con un cero en el medio: eso es su reloj de arena")
	l.txt(1060, 730, 12.5, ink, "middle", "traducido a la única cosa que la Fase VI todavía no probó.")
	l.txt(1060, 756, 12.5, green, "middle", "No está medido. Es lo próximo que hay que medir, y sale de su dibujo.")
	l.txt(1060, 776, 11.5, dim, "middle", "(y el nodo en el MEDIO, no en el primer vecino: eso era mi lectura equivocada)")

	l.txt(700, 848, 13.5, gold, "middle",
		"el dibujo y su descripción son de Jesús Nicolás Astorga · la copia limpia y la lectura, de este taller · la lectura queda como pregunta, no como hallazgo")
	l.txt(700, 892, 13, gold, "middle", "go run ./cmd/eldibujodelcapitan · Todavía no.")

	l.raw("</svg>")
	dir := filepath.Join("galeria", "laminas", "10-el-telar")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		panic(err)
	}
	full := filepath.Join(dir, "el-reloj-de-arena.svg")
	if err := os.WriteFile(full, []byte(l.b.String()), 0o644); err != nil {
		panic(err)
	}
	fmt.Printf("🖼️  lámina escrita: %s\n", full)
}
