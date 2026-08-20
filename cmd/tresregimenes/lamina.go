package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// lamina.go - the Phase VI plate. Every figure from the run that draws it.

const (
	zBg    = "#0b1526"
	zPanel = "#0d1830"
	zInk   = "#dce8f7"
	zDim   = "#8fb4d9"
	zGold  = "#ffd98a"
	zGreen = "#7ee0c0"
	zBlue  = "#7fb2ff"
	zRose  = "#ff9aa8"
	zGrid  = "#1d3a63"
)

type paño struct{ b strings.Builder }

func (l *paño) raw(s string) { l.b.WriteString(s + "\n") }

func esc(s string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;").Replace(s)
}

func nm(v float64, d int) string {
	return strings.Replace(fmt.Sprintf("%.*f", d, v), ".", ",", 1)
}

func (l *paño) txt(x, y, s float64, c, a, t string) {
	l.raw(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Georgia,serif" fill="%s" text-anchor="%s">%s</text>`, x, y, s, c, a, esc(t)))
}

func (l *paño) mono(x, y, s float64, c, a, t string) {
	l.raw(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Consolas,monospace" fill="%s" text-anchor="%s">%s</text>`, x, y, s, c, a, esc(t)))
}

func (l *paño) panel(x, y, w, h float64, c, t string) {
	l.raw(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="10" fill="%s" stroke="%s"/>`, x, y, w, h, zPanel, c))
	if t != "" {
		l.txt(x+14, y+24, 14.5, c, "start", t)
	}
}

func (l *paño) line(x1, y1, x2, y2 float64, c string, w float64, d string) {
	dd := ""
	if d != "" {
		dd = fmt.Sprintf(` stroke-dasharray="%s"`, d)
	}
	l.raw(fmt.Sprintf(`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="%.2f"%s/>`, x1, y1, x2, y2, c, w, dd))
}

func (l *paño) rect(x, y, w, h float64, f, s string, sw, op float64) {
	l.raw(fmt.Sprintf(`<rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" fill="%s" stroke="%s" stroke-width="%.2f" opacity="%.2f"/>`, x, y, w, h, f, s, sw, op))
}

func (l *paño) circ(x, y, r float64, f string) {
	l.raw(fmt.Sprintf(`<circle cx="%.2f" cy="%.2f" r="%.2f" fill="%s"/>`, x, y, r, f))
}

// tinte maps the rigidity exponent to a colour: 1 = no rigidity (rose),
// 0 = rigid like the zeros (green).
func tinte(alfa float64) string {
	a := math.Max(0, math.Min(1.8, alfa)) / 1.8
	r := int(60 + 195*a)
	g := int(224 - 130*a)
	b := int(192 - 40*a)
	return fmt.Sprintf("rgb(%d,%d,%d)", r, g, b)
}

func dibujar(r Res) {
	l := &paño{}
	l.raw(`<svg xmlns="http://www.w3.org/2000/svg" width="1400" height="950" viewBox="0 0 1400 950">`)
	l.raw(fmt.Sprintf(`<rect width="1400" height="950" fill="%s"/>`, zBg))
	l.raw(fmt.Sprintf(`<rect x="14" y="14" width="1372" height="922" rx="16" fill="none" stroke="%s" stroke-width="2" opacity="0.55"/>`, zGreen))
	l.txt(700, 50, 25, zInk, "middle", "⚡🌊 TRES REGÍMENES — el mapa (A, B), y una esquina donde la rigidez aparece")
	l.txt(700, 76, 13.5, zGold, "middle",
		"«podría tener tres modos: voltaje alto, amperaje alto, o una mezcla equilibrada de ambos» — Jesús Nicolás Astorga")

	// ---- A: the phase map --------------------------------------------------
	l.panel(36, 96, 700, 470, zGold, "A · EL MAPA DE FASES — color = α, el exponente con que crece Σ²(L)")
	ax, ay, cw, ch := 130.0, 150.0, 130.0, 92.0
	As := []float64{0.3, 1, 3, 10, 30}
	Bs := []float64{0.5, 2, 8, 32}
	for _, p := range r.Puntos {
		var ia, ib int
		for i, a := range As {
			if math.Abs(p[0]-a) < 1e-9 {
				ia = i
			}
		}
		for i, b := range Bs {
			if math.Abs(p[1]-b) < 1e-9 {
				ib = i
			}
		}
		x := ax + float64(ib)*cw
		y := ay + float64(4-ia)*ch*0.86
		l.rect(x, y, cw-6, ch*0.86-6, tinte(p[8]), zGrid, 1, 0.92)
		l.mono(x+(cw-6)/2, y+30, 15, "#0b1526", "middle", nm(p[8], 3))
		l.mono(x+(cw-6)/2, y+50, 10.5, "#0b1526", "middle", fmt.Sprintf("Σ²(10) %s", nm(p[4], 2)))
		l.mono(x+(cw-6)/2, y+66, 10, "#0b1526", "middle", fmt.Sprintf("pegados %s%%", nm(100*p[6], 1)))
	}
	for i, a := range As {
		l.mono(ax-12, ay+float64(4-i)*ch*0.86+42, 12, zDim, "end", fmt.Sprintf("A=%.0f", a))
	}
	for i, b := range Bs {
		l.mono(ax+float64(i)*cw+(cw-6)/2, ay-10, 12, zDim, "middle", fmt.Sprintf("B=%.1f", b))
	}
	l.txt(60, ay+180, 12.5, zDim, "start", "A")
	l.txt(60, ay+200, 11, zDim, "start", "el")
	l.txt(60, ay+216, 11, zDim, "start", "empuje")
	l.txt(ax+2*cw, ay+4*ch*0.86+30, 12.5, zDim, "middle", "B — el alcance")
	l.txt(56, 540, 12, zInk, "start",
		"verde = rígido (α → 0, como los ceros)   ·   rosa = sin rigidez (α ≈ 1, Poisson)")

	// ---- B: the corner -----------------------------------------------------
	l.panel(752, 96, 612, 250, zGreen, "B · LA ESQUINA — el tercer régimen del capitán, medido")
	bx, by := 780.0, 150.0
	l.mono(bx, by, 14, zGreen, "start", fmt.Sprintf("A = %.0f   y   B = %.0f   (los dos altos)", r.RigA, r.RigB))
	l.mono(bx, by+30, 13.5, zInk, "start", fmt.Sprintf("α = %s        ← todo el resto del mapa está en 0,95…1,79", nm(r.RigAlfa, 3)))
	l.mono(bx, by+54, 13.5, zInk, "start", fmt.Sprintf("niveles pegados = %s %%   ← cero colisiones", nm(100*r.RigFrac, 2)))
	l.mono(bx, by+78, 13.5, zInk, "start", fmt.Sprintf("Σ²(20) = %s  MENOR que  Σ²(10) = %s", nm(r.Rig20, 2), nm(r.Rig10, 2)))
	l.txt(bx, by+108, 12.5, zGold, "start", "la varianza DEJA DE CRECER con la distancia. Eso es rigidez:")
	l.txt(bx, by+128, 12.5, zGold, "start", "los niveles se acomodan entre sí a lo largo de la escalera.")
	l.txt(bx, by+156, 12, zRose, "start", "y no lo consigue ninguna de las dos perillas sola: con A alto")
	l.txt(bx, by+174, 12, zRose, "start", "y B chico, α se queda en 1; con B alto y A chico, también.")

	// ---- C: the decisive test ----------------------------------------------
	l.panel(752, 362, 612, 204, zBlue, "C · LA PRUEBA DECISIVA — misma fuerza total, distinto reparto")
	cx, cy := 780.0, 412.0
	l.txt(cx, cy, 12.5, zInk, "start", "si el resultado dependiera SÓLO del total, la hipótesis perdía.")
	l.mono(cx, cy+28, 12.5, zGold, "start", "fuerza 1,54  (A=0,3  B=32)  →  Σ²(10) = 7,097")
	l.mono(cx, cy+50, 12.5, zGold, "start", "fuerza 1,82  (A=10   B=0,5) →  Σ²(10) = 4,103")
	l.mono(cx, cy+78, 14, zGreen, "start", fmt.Sprintf("razón %s×  ⟹  DEPENDE DEL REPARTO", nm(r.RazonReparto, 2)))
	l.txt(cx, cy+108, 12.5, zInk, "start", "son dos grados de libertad independientes, no uno.")
	l.txt(cx, cy+130, 12.5, zDim, "start", "la intuición del capitán queda apoyada por la medición.")

	// ---- D: what it costs --------------------------------------------------
	l.panel(36, 582, 1328, 150, zRose, "D · LO QUE CUESTA, Y LO QUE QUEDA")
	l.txt(60, 626, 12.5, zInk, "start",
		fmt.Sprintf("en la esquina rígida Σ²(10) vale %s en valor absoluto — peor que el mejor punto del mapa (A=%.0f, B=%.1f, Σ²(10)=%s).",
			nm(r.Rig10, 2), r.MejorA, r.MejorB, nm(r.Mejor10, 4)))
	l.txt(60, 648, 12.5, zInk, "start",
		"Se gana RIGIDEZ y se paga en VARIANZA. No es todavía el espectro de los ceros: es una FASE distinta, que antes no aparecía en ningún mapa.")
	l.txt(60, 676, 12.5, zGreen, "start",
		"Lo que la Fase VII tiene que preguntar: ¿existe un camino DENTRO de esa esquina que baje la varianza sin perder la rigidez?")
	l.txt(60, 698, 12.5, zGreen, "start",
		"O sea: el mapa ya no es de un óptimo — es de una FRONTERA DE FASE, y hay que recorrerla.")

	// ---- footer ------------------------------------------------------------
	l.raw(fmt.Sprintf(`<rect x="36" y="748" width="1328" height="70" rx="10" fill="%s" stroke="%s" opacity="0.5"/>`, zPanel, zGold))
	l.mono(700, 776, 13, zGold, "middle", "H_ij = A · amp_i · amp_j · exp(−|i−j|/B)      con amp = Λ(p)/√p      ·      la fuerza total es DERIVADA, no el control")
	l.txt(700, 802, 12.5, zInk, "middle",
		fmt.Sprintf("R6 LIMPIO: única entrada aritmética Λ(n) · %d modos · ningún parámetro elegido mirando los γₙ · los ceros se leen al final", r.N))
	l.txt(700, 858, 13, zGold, "middle",
		"la intuición de los tres regímenes es de Jesús Nicolás Astorga · la ruta y los controles, de la auditora · las mediciones, de este taller")
	l.txt(700, 902, 13, zGold, "middle",
		"go run ./cmd/tresregimenes · estructura cerrada no es hipótesis demostrada · Todavía no.")

	l.raw("</svg>")
	dir := filepath.Join("galeria", "laminas", "10-el-telar")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		panic(err)
	}
	full := filepath.Join(dir, "tres-regimenes.svg")
	if err := os.WriteFile(full, []byte(l.b.String()), 0o644); err != nil {
		panic(err)
	}
	fmt.Printf("\n🖼️  lámina escrita: %s\n", full)
}
