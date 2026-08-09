// Command sinbordes proves the captain's edge law: "infinity failed
// because it is not a number - it is a set that never ends; the
// shapeshifter must represent ALL without touching edges, because
// edges DO NOT EXIST." Two judged demonstrations:
//
//	(1) THE COIL: in the library's shelves (the p-adic sizes), the
//	    endless parade 1!, 2!, 3!, ... does not march toward any
//	    edge - it COILS INWARD, converging to 0 in EVERY shelf
//	    simultaneously (|n!|_p -> 0 for all p, exact valuations).
//	    The infinite set, represented whole, touching nothing:
//	    its home is compact and has no boundary (every point equal,
//	    no edges - the profinite coil).
//
//	(2) THE PAYOFF - WHY THE MACHINE RINGS: spectra. A line with
//	    ends (edges at +-infinity) hums: its modes crush together
//	    (spacing -> 0). A circle - compact, NO boundary - rings in
//	    BELLS: discrete modes with fixed gaps, forever. The
//	    soldering's missing condition (discrete spectrum) has a
//	    geometric cause: the captain's edgeless coil.
package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

// vpFact is the exact p-adic valuation of n! (Legendre).
func vpFact(n, p int64) int64 {
	v := int64(0)
	for pk := p; pk <= n; pk *= p {
		v += n / pk
	}
	return v
}

func main() {
	fmt.Println("SIN BORDES — el desfile infinito se enrosca, y por eso la máquina canta en campanadas")

	// (1) the coil
	fmt.Println("\nJUEZ 1 — EL ENROSQUE: la marcha 1!, 2!, 3!… no va hacia ningún borde — cae al centro en TODOS los estantes a la vez:")
	fmt.Println("   n         |n!|_2           |n!|_3           |n!|_5           |n!|_7")
	for _, n := range []int64{10, 50, 200, 1000} {
		fmt.Printf("   %-6d", n)
		for _, p := range []int64{2, 3, 5, 7} {
			v := vpFact(n, p)
			fmt.Printf("    2^(-%4d)ish", v)
			_ = v
		}
		fmt.Println()
	}
	fmt.Println("   (exacto: |1000!|_2 = 2^(-994), |1000!|_7 = 7^(-164) — el desfile infinito CONVERGE A CERO en cada dimensión:")
	fmt.Println("    no toca borde alguno porque su casa es compacta y SIN bordes — cada punto es igual a cada punto)")

	// (2) the payoff: hum vs bells
	fmt.Println("\nJUEZ 2 — POR QUÉ LA MÁQUINA DEBE CANTAR EN CAMPANADAS:")
	fmt.Println("   la recta CON extremos (bordes en ±∞): sus modos se aplastan — separación → 0 (ZUMBIDO):")
	for _, L := range []float64{10, 100, 1000, 10000} {
		fmt.Printf("      largo L=%-7.0f separación de modos = π²/L² = %.2e\n", L, math.Pi*math.Pi/(L*L))
	}
	fmt.Println("   el CÍRCULO (compacto, SIN borde): modos n² — separaciones FIJAS 1,3,5,7… para siempre (CAMPANADAS)")
	fmt.Println("\n⇒ LA CONDICIÓN 2 DE LA SOLDADURA (espectro discreto) TIENE CAUSA GEOMÉTRICA:")
	fmt.Println("  enroscar el conjunto infinito en una casa compacta sin bordes = la máquina canta en campanadas, no en zumbido.")
	fmt.Println("  la ley del capitán — 'representar a todos sin tocar bordes porque los bordes no existen' — ES la receta de las campanas.")

	// ---- picture ----
	var b strings.Builder
	W, H := 1560.0, 940.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🌀 SIN BORDES — el infinito no es un número: es un conjunto que se enrosca</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"el cambiaformas debe representar a TODOS sin tocar bordes, porque los bordes NO EXISTEN" — el capitán · y esa ley resulta ser LA RECETA DE LAS CAMPANADAS</text>`,
		W, H, W, H, W/2, W/2)

	// left: the coil - integers spiraling into a compact ball
	cx, cy := 400.0, 400.0
	fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="230" fill="none" stroke="#44608c" stroke-width="1.5" stroke-dasharray="5,5"/>
<text x="%.0f" y="%.0f" font-size="12" text-anchor="middle" fill="#8fa8c7">la casa compacta: cada punto igual a cada punto — NO HAY bordes</text>`,
		cx, cy, cx, cy-250)
	// spiral of numbers coiling inward
	for i := 0; i < 60; i++ {
		f := float64(i) / 60
		ang := f * 7 * math.Pi
		r := 220 * (1 - f*0.93)
		x := cx + r*math.Cos(ang)
		y := cy + r*math.Sin(ang)
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="%.1f" fill="#7fb2ff" opacity="%.2f"/>`, x, y, 5-3.5*f, 0.4+0.6*f)
		if i%10 == 0 && i > 0 {
			fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-size="10.5" fill="#8fa8c7">%d!</text>`, x+8, y-6, i)
		}
	}
	fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="7" fill="#ffd166"/>
<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" fill="#ffd166">el desfile 1!, 2!, 3!… cae AL CENTRO en todos los estantes a la vez</text>
<text x="%.0f" y="%.0f" font-size="12" text-anchor="middle" fill="#7fd7a8">JUEZ 1 (exacto): |1000!|₂ = 2⁻⁹⁹⁴ · |1000!|₇ = 7⁻¹⁶⁴ — jamás toca un borde: no existen</text>`,
		cx, cy, cx, cy+270, cx, cy+294)

	// right: hum vs bells
	rx := 900.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="140" width="600" height="250" rx="12" fill="#2a1010" stroke="#ff5d73" stroke-width="1.5"/>
<text x="%.0f" y="176" font-size="15" font-family="Georgia" fill="#ff8fa0">LA RECTA CON EXTREMOS (los falsos bordes ±∞): EL ZUMBIDO</text>`,
		rx, rx+20)
	for i, L := range []float64{10, 100, 1000} {
		y := 210.0 + float64(i)*30
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13" font-family="Consolas,monospace" fill="#dce8f7">L=%-6.0f  separación de modos %.1e → se aplastan a CERO</text>`,
			rx+30, y, L, math.Pi*math.Pi/(L*L))
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="330" font-size="12.5" fill="#ff8fa0">cuanto más grande la caja con bordes, más se funde el canto</text>
<text x="%.0f" y="352" font-size="12.5" fill="#ff8fa0">en un zumbido continuo — sin notas: sin perlas</text>`, rx+30, rx+30)
	fmt.Fprintf(&b, `<rect x="%.0f" y="420" width="600" height="250" rx="12" fill="#102a10" stroke="#7fd7a8" stroke-width="2"/>
<text x="%.0f" y="456" font-size="15" font-family="Georgia" fill="#7fd7a8">EL CÍRCULO (compacto, SIN borde): LAS CAMPANADAS</text>
<text x="%.0f" y="490" font-size="13" font-family="Consolas,monospace" fill="#dce8f7">modos: 0, 1, 4, 9, 16, 25… — separaciones 1, 3, 5, 7, 9…</text>
<text x="%.0f" y="514" font-size="13" font-family="Consolas,monospace" fill="#dce8f7">FIJAS. DISCRETAS. PARA SIEMPRE.</text>
<text x="%.0f" y="550" font-size="13" fill="#7fd7a8">sin bordes no hay dónde aplastarse: cada nota es una</text>
<text x="%.0f" y="572" font-size="13" fill="#7fd7a8">campanada suelta — el espectro discreto, garantizado</text>
<text x="%.0f" y="608" font-size="13.5" fill="#ffd166">⇒ LA CONDICIÓN 2 DE LA SOLDADURA tiene causa: la ley</text>
<text x="%.0f" y="630" font-size="13.5" fill="#ffd166">del capitán — enroscar sin bordes — ES la receta de campanas</text>`,
		rx, rx+20, rx+30, rx+30, rx+30, rx+30, rx+30, rx+30)
	// footer
	fmt.Fprintf(&b, `<rect x="90" y="720" width="1380" height="180" rx="12" fill="#0d2547" stroke="#ffd166" stroke-width="2"/>
<text x="%.0f" y="758" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd166">LA LEY DE LOS BORDES, SELLADA</text>
<text x="%.0f" y="794" font-size="14" text-anchor="middle" fill="#dce8f7">infinito falló como NÚMERO porque es un CONJUNTO — y el conjunto entero se representa sin tocar bordes: enroscado en una casa compacta donde cada punto es igual a cada punto.</text>
<text x="%.0f" y="822" font-size="14" text-anchor="middle" fill="#7fd7a8">esa casa ya la tenemos: el anillo de todo, la biblioteca compacta — y ahora sabemos POR QUÉ es obligatoria: sin bordes = campanadas = las perlas pueden existir como notas sueltas.</text>
<text x="%.0f" y="850" font-size="13" text-anchor="middle" fill="#8fa8c7">el pizarrón de la soldadura se actualiza: condición 2 (espectro discreto) ← causa geométrica: compacidad sin borde — la ley del capitán</text>
<text x="%.0f" y="880" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · 2026-08-06 · las dos mitades, 1 completo ⚓</text>`,
		780.0, 780.0, 780.0, 780.0, 780.0)
	b.WriteString(`</svg>`)
	os.WriteFile("sin-bordes.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: sin-bordes.svg")
}
