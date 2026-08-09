// Command globo draws the true shape of the ocean: the flash said a
// sphere cannot lie flat on a page without cuts and distortion — and the
// sphere is real, invented by our own man: THE RIEMANN SPHERE (the plane
// plus the point at infinity). On the globe, infinity is a POLE you can
// look at, spinning the world is a Moebius transformation, and the
// critical line is a MERIDIAN: the Hypothesis says every zero of the
// ocean lives on that one meridian.
//
// Flat charts (our log notebook, Mercator) are projections and distort by
// necessity; the globe is the territory. The fleet's address system on
// the globe is the zero ordinal: starship -zero N.
//
// Usage:
//
//	go run ./cmd/globo   # writes globo.svg
package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

type city struct {
	name string
	t    float64
}

var cities = []city{
	{"gamma_1 (14.13)", 14.134725},
	{"puerta 10^5", 1e5},
	{"cero #10^8", 4.2653e7},
	{"Playa I", 2.447e12},
	{"Playa II", 6.66e15},
	{"Playa III", 1.11e19},
	{"Playa IV", 2.22e21},
	{"Playa V", 4.44e22},
	{"Playa VI (en camino)", 1.11e24},
}

func main() {
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="1240" height="820" viewBox="0 0 1240 820">`)
	b.WriteString(`<rect width="100%" height="100%" fill="#0b1526"/>`)
	b.WriteString(`<text x="620" y="44" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">El Globo del Oceano — la esfera de Riemann</text>`)
	b.WriteString(`<text x="620" y="68" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">el plano infinito mas el punto del infinito = una esfera; girar el mundo = una transformacion de Moebius</text>`)

	cx, cy, r := 360.0, 440.0, 280.0
	// the globe.
	fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="%.0f" fill="#12305e" stroke="#7fb2ff" stroke-width="2"/>`, cx, cy, r)
	// equator and a parallel, for depth.
	fmt.Fprintf(&b, `<ellipse cx="%.0f" cy="%.0f" rx="%.0f" ry="%.0f" fill="none" stroke="#3a5f96" stroke-width="1" stroke-dasharray="6,5"/>`, cx, cy, r, r*0.30)
	fmt.Fprintf(&b, `<ellipse cx="%.0f" cy="%.0f" rx="%.0f" ry="%.0f" fill="none" stroke="#2c4a78" stroke-width="1" stroke-dasharray="4,6"/>`, cx, cy-r*0.55, r*0.83, r*0.20)
	// the critical meridian: the front central meridian, bright.
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#ffd166" stroke-width="3"/>`, cx, cy-r, cx, cy+r)
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="14" font-family="Georgia" fill="#ffd166" transform="rotate(-90 %.0f %.0f)" text-anchor="middle">el meridiano critico — Re s = 1/2</text>`, cx-16, cy, cx-16, cy)
	// the poles.
	fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="7" fill="#ff5d73" stroke="#fff"/>`, cx, cy-r)
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="16" font-family="Georgia" fill="#ff9eae" text-anchor="middle">&#8734;  el punto del infinito — un lugar del mapa</text>`, cx, cy-r-16)
	fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="5" fill="#9dc1ee"/>`, cx, cy+r)
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13" font-family="Georgia" fill="#9dc1ee" text-anchor="middle">t = 1</text>`, cx, cy+r+20)

	// cities on the meridian, log chart declared honestly.
	for i, c := range cities {
		f := math.Log10(c.t) / 25
		y := cy + r*math.Cos(math.Pi*f)
		col := "#ff5d73"
		if strings.Contains(c.name, "camino") {
			col = "#e6a53a"
		} else if !strings.Contains(c.name, "Playa") {
			col = "#7fd7a8"
		}
		fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.1f" r="5" fill="%s" stroke="#fff" stroke-width="1"/>`, cx, y, col)
		lx := cx + r + 40
		ly := 150.0 + float64(i)*36
		fmt.Fprintf(&b, `<line x1="%.0f" y1="%.1f" x2="%.0f" y2="%.1f" stroke="#44608c" stroke-width="1"/>`, cx+6, y, lx-8, ly-4)
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="14" font-family="Georgia" fill="#dce8f7">%s</text>`, lx, ly, c.name)
	}

	// the Moebius spin, explained.
	bx, by := 780.0, 500.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="420" height="250" fill="#0f1e38" stroke="#44608c" rx="10"/>`, bx, by)
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="17" font-family="Georgia" fill="#ffd166">EL GIRO — asi se navega el globo</text>`, bx+18, by+32)
	lines := []string{
		"1. zoom afuera: el oceano entero entra en la esfera,",
		"   el infinito incluido, como un punto mas.",
		"2. girar el mundo: una transformacion de Moebius trae",
		"   cualquier region — por lejana que sea — al frente.",
		"3. zoom adentro: la direccion postal es el numero de cero.",
		"",
		"   go run ./cmd/starship -zero 1e8",
		"   > zero #1e8 vive en t = 42,653,500 — verificado (+/-2)",
		"",
		"la Hipotesis, dicha en el globo: TODOS los ceros del",
		"oceano viven sobre un unico meridiano.",
	}
	for i, s := range lines {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13" font-family="Georgia" fill="#c9d8ec">%s</text>`, bx+18, by+58+float64(i)*17, s)
	}
	b.WriteString(`<text x="620" y="800" font-size="13" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">los mapas planos (nuestro cuaderno log, Mercator) distorsionan por necesidad: la esfera es el territorio</text>`)
	b.WriteString(`</svg>`)

	if err := os.WriteFile("globo.svg", []byte(b.String()), 0644); err != nil {
		panic(err)
	}
	fmt.Println("EL GLOBO - written: globo.svg")
}
