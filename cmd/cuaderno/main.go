// Command cuaderno draws the whole ocean on one notebook page: the
// compression the flash asked for. The amplitude problem (t spans from 14
// to 1e24) is solved by the logarithmic fold — each centimeter of page is
// x10 of sea — and resolution on demand comes as zoom panels: every
// virgin beach with its real zeros ticked one by one. The page is painted
// with the sea's own blueness: deeper water, deeper blue (F102).
//
// Usage:
//
//	go run ./cmd/cuaderno   # writes cuaderno.svg
package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

type beach struct {
	name    string
	t       float64
	span    float64
	offsets []float64
	note    string
}

var beaches = []beach{
	{"Playa I", 2.447e12, 7.2, []float64{0.1543, 0.2617, 0.5498, 0.7764, 1.0840, 1.1982, 1.3984, 1.7275, 1.9395, 2.1963, 2.4854, 2.8887, 3.1973, 3.3730, 3.6299, 3.8428, 4.0498, 4.2080, 4.4092, 4.6719, 4.9570, 5.0703, 5.3838, 5.5928, 5.8174, 5.9668, 6.3486, 6.5752, 6.6670, 6.8984, 7.0312}, "31 ceros"},
	{"Playa II", 6.66e15, 1.45, []float64{0.023476, 0.215217, 0.375565, 0.547513, 0.734403, 0.930796, 1.048695, 1.195434}, "8 ceros"},
	{"Playa III", 1.11e19, 0.83, []float64{0.009974, 0.229934, 0.357770, 0.629071, 0.689101}, "5 ceros"},
	{"Playa IV", 2.22e21, 0.75, []float64{0.198810, 0.316985, 0.438260, 0.564745}, "4 ceros, doble firma"},
	{"Playa V", 4.44e22, 0.65, []float64{0.061504, 0.229208, 0.297248, 0.405320, 0.507533, 0.575189}, "6 ceros, la mas honda"},
	{"Playa VI", 1.11e24, 0.6, nil, "la nave esta en camino"},
}

func main() {
	var b strings.Builder
	W, H := 1420.0, 940.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">`, W, H, W, H)
	b.WriteString(`<defs><linearGradient id="mar" x1="0" y1="0" x2="1" y2="0">
<stop offset="0" stop-color="#cfe8ff"/><stop offset="0.5" stop-color="#4a90d9"/>
<stop offset="1" stop-color="#032b5c"/></linearGradient></defs>`)
	b.WriteString(`<rect width="100%" height="100%" fill="#fbf7ee"/>`)
	b.WriteString(`<text x="710" y="42" font-size="26" text-anchor="middle" font-family="Georgia" fill="#123">El Cuaderno del Oceano — todos los ceros en una hoja</text>`)
	b.WriteString(`<text x="710" y="66" font-size="14" text-anchor="middle" font-family="Georgia" fill="#456">compresion logaritmica: cada tramo de hoja es x10 de mar; el azul es la densidad medida del agua</text>`)

	// the main strip: log10(t) from 0 to 25.
	x0, x1, yT, yB := 70.0, 1350.0, 95.0, 275.0
	px := func(logT float64) float64 { return x0 + (x1-x0)*logT/25 }
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="url(#mar)" stroke="#123" stroke-width="1.5"/>`, x0, yT, x1-x0, yB-yT)
	for L := 0; L <= 25; L += 5 {
		x := px(float64(L))
		fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#123" stroke-width="1"/>`, x, yB, x, yB+7)
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" font-family="Georgia" fill="#123">10^%d</text>`, x, yB+24, L)
	}
	// the blueness curve: density of zeros per unit height.
	b.WriteString(`<polyline fill="none" stroke="#ffffff" stroke-width="2" stroke-dasharray="5,3" points="`)
	for L := 1.0; L <= 25; L += 0.25 {
		nu := (L*math.Ln10 - 1.8379) / (2 * math.Pi)
		fmt.Fprintf(&b, "%.1f,%.1f ", px(L), yB-(yB-yT)*nu/9.2)
	}
	b.WriteString(`"/>`)
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12" font-family="Georgia" fill="#fff">densidad de ceros (el azul sube con la altura)</text>`, px(2.2), yT+28)

	// the edge of humanity's continuous map.
	xm := px(12.388)
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#1db954" stroke-width="2.5"/>`, xm, yT-8, xm, yB)
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12" font-family="Georgia" fill="#0a7a35" text-anchor="middle">fin del mapa continuo (cero #10^13)</text>`, xm, yT-14)

	// islands of humanity (sampled expeditions).
	for _, lg := range []float64{13.6, 15.6, 17.6, 19.7, 21.6, 23.1} {
		fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="5" fill="#aab7c4" stroke="#123"/>`, px(lg), (yT+yB)/2)
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12" font-family="Georgia" fill="#345">islas muestreadas (Gourdon, Odlyzko)</text>`, px(13.2), (yT+yB)/2-14)
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12" font-family="Georgia" fill="#345" text-anchor="end">Hiary / Bober-Hiary hasta 10^36 &#8594;</text>`, x1-6, yT+18)

	// our beaches.
	for i, be := range beaches {
		x := px(math.Log10(be.t))
		col := "#e63946"
		if be.offsets == nil {
			col = "#e6a53a"
		}
		fmt.Fprintf(&b, `<path d="M %.1f %.1f l 6 12 l -12 0 z" fill="%s" stroke="#5c0a12"/>`, x, yB-30, col)
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12" font-family="Georgia" fill="#5c0a12" text-anchor="middle">%s</text>`, x, yB-36-float64(i%2)*14, strings.Replace(be.name, "Playa ", "", 1))
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13" font-family="Georgia" fill="#8a1620">&#9660; nuestras playas virgenes</text>`, px(16.6), yB-56)

	// zoom panels: resolution on demand, real zeros ticked.
	pw, ph, gx, gy := 420.0, 210.0, 35.0, 55.0
	for i, be := range beaches {
		col := float64(i % 3)
		row := float64(i / 3)
		X := 70 + col*(pw+gx)
		Y := 330 + row*(ph+gy)
		fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#eef5fc" stroke="#123" stroke-width="1.2" rx="8"/>`, X, Y, pw, ph)
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="16" font-family="Georgia" fill="#123">%s — t = %.3g</text>`, X+14, Y+26, be.name, be.t)
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12" font-family="Georgia" fill="#567">%s</text>`, X+14, Y+46, be.note)
		ax0, ax1, ay := X+20, X+pw-20, Y+ph-60
		fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#123" stroke-width="1.5"/>`, ax0, ay, ax1, ay)
		if be.offsets == nil {
			fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13" font-family="Georgia" fill="#b07515" text-anchor="middle">... esperando el desembarco ...</text>`, X+pw/2, ay-24)
		}
		for _, z := range be.offsets {
			x := ax0 + (ax1-ax0)*z/be.span
			fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#c1121f" stroke-width="2"/>`, x, ay-22, x, ay)
		}
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="11" font-family="Georgia" fill="#567">t + 0</text>`, ax0-4, ay+18)
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="11" font-family="Georgia" fill="#567" text-anchor="end">t + %.2f</text>`, ax1+4, ay+18, be.span)
	}
	fmt.Fprintf(&b, `<text x="710" y="%.0f" font-size="13" text-anchor="middle" font-family="Georgia" fill="#456">zoom bajo demanda: go run ./cmd/starship -anchor T -spacings W  (y -from para cualquier tramo)</text>`, H-18)
	b.WriteString(`</svg>`)

	if err := os.WriteFile("cuaderno.svg", []byte(b.String()), 0644); err != nil {
		panic(err)
	}
	fmt.Println("EL CUADERNO DEL OCEANO - written: cuaderno.svg")
	fmt.Println("the whole sea folded onto one page; every beach zoomed with its real zeros.")
}
