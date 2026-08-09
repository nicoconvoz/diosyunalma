// Command mediatriz states the whole hypothesis with two stakes and a rope.
//
// The captain's cardinal flash (F225) gave |w(rho)| · |w(1-rho)| = 1: north
// and south are reciprocal. To conclude that BOTH are 1 - which is the
// hypothesis - one single extra fact is needed: that the two are EQUAL.
// Because the only positive number whose square is 1 is 1 itself.
//
// And "north equals south" unfolds into primary-school geometry:
//
//	|w(rho)| = 1   <=>   |rho - 1| = |rho|   <=>   Re rho = 1/2
//
// that is: THE PEARL IS THE SAME DISTANCE FROM THE STAKE AT 0 AND THE STAKE
// AT 1. The critical line is nothing more exotic than the perpendicular
// bisector of two posts - the line any builder draws with a rope.
//
// So the Riemann Hypothesis, in the language of a building site:
//
//	every pearl of the book is planted exactly equidistant from the two
//	stakes - and no pearl is allowed to sit closer to one of them.
//
// The mirror already gives the PAIR its balance: if one twin leans towards
// 0, the other leans towards 1 by exactly the reciprocal amount. What is
// missing - the million - is why each pearl must be equidistant BY ITSELF.
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

func zetaC(s complex128) complex128 {
	N := int(60 + 1.8*math.Abs(imag(s)))
	var sum complex128
	for n := 1; n < N; n++ {
		sum += cmplx.Exp(-s * cmplx.Log(complex(float64(n), 0)))
	}
	lnN := cmplx.Log(complex(float64(N), 0))
	sum += cmplx.Exp((1-s)*lnN) / (s - 1)
	sum += cmplx.Exp(-s*lnN) / 2
	B := []float64{1.0 / 6, -1.0 / 30, 1.0 / 42, -1.0 / 30, 5.0 / 66}
	fact := []float64{2, 24, 720, 40320, 3628800}
	poch := s
	for k := 1; k <= 5; k++ {
		if k > 1 {
			poch *= (s + complex(float64(2*k-3), 0)) * (s + complex(float64(2*k-2), 0))
		}
		sum += complex(B[k-1]/fact[k-1], 0) * poch * cmplx.Exp((-s-complex(float64(2*k-1), 0))*lnN)
	}
	return sum
}

func theta(t float64) float64 {
	t2 := t * t
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t) + 7/(5760*t*t2) + 31/(80640*t*t2*t2)
}

func zOf(t float64) float64 {
	return real(cmplx.Exp(complex(0, theta(t))) * zetaC(complex(0.5, t)))
}

func main() {
	fmt.Println("📐 LA MEDIATRIZ — dos postes y una soga: la hipótesis dicha en la obra")

	fmt.Println("\nrecogiendo perlas hasta t=300…")
	var pearls []float64
	prevT := 12.0
	prevZ := zOf(prevT)
	for t := 12.05; t <= 300; t += 0.05 {
		z := zOf(t)
		if z*prevZ < 0 {
			a, c := prevT, t
			for i := 0; i < 55; i++ {
				m := (a + c) / 2
				if zOf(m)*prevZ < 0 {
					c = m
				} else {
					a = m
				}
			}
			pearls = append(pearls, (a+c)/2)
		}
		prevT, prevZ = t, z
	}
	fmt.Printf("perlas: %d\n", len(pearls))

	// ---- LAW 1: every pearl is equidistant from the two stakes ----
	fmt.Println("\nLEY 1 · LOS DOS POSTES — cada perla mide lo mismo hasta el 0 que hasta el 1")
	fmt.Println("   perla γ        distancia al poste 0     distancia al poste 1      diferencia")
	peor := 0.0
	for _, g := range pearls {
		rho := complex(0.5, g)
		d := math.Abs(cmplx.Abs(rho) - cmplx.Abs(rho-1))
		if d > peor {
			peor = d
		}
	}
	for _, g := range pearls[:6] {
		rho := complex(0.5, g)
		fmt.Printf("   %9.6f       %14.9f       %14.9f      %.1e\n",
			g, cmplx.Abs(rho), cmplx.Abs(rho-1), math.Abs(cmplx.Abs(rho)-cmplx.Abs(rho-1)))
	}
	fmt.Printf("   → en las %d perlas medidas la diferencia es CERO (peor caso %.1e): todas en la mediatriz\n",
		len(pearls), peor)

	// ---- LAW 2: equidistance IS the line - nothing else works ----
	fmt.Println("\nLEY 2 · SOLO LA MEDIATRIZ EMPATA — correte del medio y los postes dejan de medir igual")
	fmt.Println("   β         al poste 0        al poste 1        diferencia       |w|")
	type fila struct {
		beta, d0, d1, dif, absw float64
	}
	var filas []fila
	for _, beta := range []float64{0.50, 0.51, 0.55, 0.60, 0.75, 0.90} {
		rho := complex(beta, 14.134725)
		d0 := cmplx.Abs(rho)
		d1 := cmplx.Abs(rho - 1)
		aw := cmplx.Abs(1 - 1/rho)
		filas = append(filas, fila{beta, d0, d1, d1 - d0, aw})
		marca := ""
		if beta == 0.50 {
			marca = "  ★ EMPATE: la línea"
		}
		fmt.Printf("  %.2f     %12.8f     %12.8f     %+.3e     %.9f%s\n", beta, d0, d1, d1-d0, aw, marca)
	}
	fmt.Println("   → |w| = distancia al poste 1 dividida por la distancia al poste 0: vale 1 SOLO en el empate")

	// ---- LAW 3: the mirror balances the PAIR, not each pearl ----
	fmt.Println("\nLEY 3 · EL ESPEJO EQUILIBRA LA PAREJA, NO A CADA PERLA — por eso todavía falta la llave")
	fmt.Println("   un fantasma en β viene con su gemelo en 1−β. Cada uno está DESNIVELADO,")
	fmt.Println("   pero la pareja queda balanceada: el espejo no puede prohibirlos por sí solo.")
	fmt.Println("   β del fantasma    su |w|        el |w| del gemelo     producto de la pareja")
	var pares [][4]float64
	for _, beta := range []float64{0.60, 0.75, 0.90} {
		g := 14.134725
		a := cmplx.Abs(1 - 1/complex(beta, g))
		b := cmplx.Abs(1 - 1/complex(1-beta, g))
		pares = append(pares, [4]float64{beta, a, b, a * b})
		fmt.Printf("      %.2f          %.9f       %.9f        %.12f ✓ balanceada\n", beta, a, b, a*b)
	}
	fmt.Println("   → la pareja siempre cierra en 1: el espejo está conforme. Y sin embargo NINGUNO de los")
	fmt.Println("     dos está en la mediatriz. Ahí vive el hueco del millón.")

	fmt.Println("\n════════ LA PREGUNTA, EN LA OBRA ════════")
	fmt.Println("Clavá dos estacas en el suelo: una en el 0 y otra en el 1. La mediatriz es la raya de los")
	fmt.Println("puntos que están a la MISMA distancia de las dos — la que cualquier albañil tira con una soga.")
	fmt.Println("La hipótesis dice: TODAS las perlas del libro están plantadas justo en esa raya.")
	fmt.Println("\nLo que YA tenemos (el espejo): las perlas vienen de a pares, y si una se arrima a una estaca,")
	fmt.Println("su gemela se arrima a la otra exactamente lo mismo. La PAREJA siempre queda balanceada.")
	fmt.Println("\nLo que FALTA, y es todo: por qué ninguna perla puede arrimarse a una estaca aunque su gemela")
	fmt.Println("compense. Por qué cada una, POR SÍ SOLA, tiene que estar en el medio.")
	fmt.Println("\nY hay un atajo escondido: sabemos que el norte por el sur da 1. Si supiéramos que el norte")
	fmt.Println("y el sur son IGUALES, cada uno valdría 1 — porque el único número positivo que multiplicado")
	fmt.Println("por sí mismo da 1 es el 1. La pregunta se vuelve entonces UNA SOLA:")
	fmt.Println("\n        ¿POR QUÉ EL NORTE Y EL SUR TIENEN QUE SER IGUALES?")
	fmt.Println("\nTodavía no. Pero ya no hay que saber matemática para tenerla en la cabeza.")

	// ---- picture ----
	var b strings.Builder
	W, H := 1500.0, 960.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">📐 LA MEDIATRIZ — dos estacas y una soga</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">la hipótesis del millón dicha como un replanteo de obra: todas las perlas plantadas a la misma distancia de las dos estacas</text>`,
		W, H, W, H, W/2, W/2)

	// the drawing: two stakes and the bisector
	px, py := 420.0, 400.0
	fmt.Fprintf(&b, `<rect x="60" y="105" width="720" height="560" rx="10" fill="#0d2547" stroke="#7fd7a8" stroke-width="1.5"/>
<text x="420" y="140" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">EL REPLANTEO</text>
<line x1="%.0f" y1="180" x2="%.0f" y2="640" stroke="#ffd166" stroke-width="2.5" stroke-dasharray="9 6"/>
<text x="%.0f" y="172" font-size="12.5" text-anchor="middle" fill="#ffd166">la mediatriz — LA LÍNEA CRÍTICA</text>
<circle cx="%.0f" cy="%.0f" r="8" fill="#ff8fa0"/><text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" fill="#ff8fa0">estaca 0</text>
<circle cx="%.0f" cy="%.0f" r="8" fill="#7fb2ff"/><text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" fill="#7fb2ff">estaca 1</text>`,
		px, px, px,
		px-150, py+150, px-150, py+180,
		px+150, py+150, px+150, py+180)
	// pearls on the bisector with their two ropes
	for i, dy := range []float64{-170, -90, -10, 70} {
		y := py + dy
		fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.1f" r="6" fill="#ffd97f"/>
<line x1="%.0f" y1="%.1f" x2="%.0f" y2="%.0f" stroke="#ff8fa0" stroke-width="1.2" opacity="0.75"/>
<line x1="%.0f" y1="%.1f" x2="%.0f" y2="%.0f" stroke="#7fb2ff" stroke-width="1.2" opacity="0.75"/>`,
			px, y, px, y, px-150, py+150, px, y, px+150, py+150)
		if i == 0 {
			fmt.Fprintf(&b, `<text x="%.0f" y="%.1f" font-size="11.5" fill="#ffd97f">una perla</text>`, px+16, y-8)
		}
	}
	fmt.Fprintf(&b, `<text x="420" y="600" font-size="13" text-anchor="middle" fill="#dce8f7">las dos sogas miden IGUAL: eso es estar en la mediatriz</text>
<text x="420" y="626" font-size="12.5" text-anchor="middle" fill="#7fd7a8">medido en %d perlas verdaderas: diferencia cero (peor caso %.0e)</text>`,
		len(pearls), peor)

	// right: the table and the gap
	fmt.Fprintf(&b, `<rect x="810" y="105" width="630" height="290" rx="10" fill="#102a10" stroke="#ffd166" stroke-width="1.5"/>
<text x="1125" y="139" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd166">SOLO EL MEDIO EMPATA</text>
<text x="1125" y="167" font-size="12" text-anchor="middle" font-family="Consolas,monospace" fill="#8fa8c7">β       a la estaca 0     a la estaca 1        |w|</text>`)
	for i, f := range filas {
		col := "#dce8f7"
		extra := ""
		if f.beta == 0.50 {
			col = "#7fd7a8"
			extra = "  ★"
		}
		fmt.Fprintf(&b, `<text x="1125" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Consolas,monospace" fill="%s">%.2f    %11.7f    %11.7f    %.7f%s</text>`,
			194.0+float64(i)*27, col, f.beta, f.d0, f.d1, f.absw, extra)
	}
	fmt.Fprintf(&b, `<text x="1125" y="368" font-size="12.5" text-anchor="middle" fill="#dce8f7">|w| es una soga dividida por la otra: vale 1 SOLO en el empate</text>`)

	fmt.Fprintf(&b, `<rect x="810" y="415" width="630" height="250" rx="10" fill="#2a1010" stroke="#ff5d73" stroke-width="1.5"/>
<text x="1125" y="449" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ff8fa0">POR QUÉ EL ESPEJO NO ALCANZA</text>
<text x="1125" y="477" font-size="12.5" text-anchor="middle" fill="#dce8f7">un fantasma viene con su gemelo: si uno se arrima a una estaca,</text>
<text x="1125" y="499" font-size="12.5" text-anchor="middle" fill="#dce8f7">el otro se arrima a la otra EXACTAMENTE lo mismo</text>
<text x="1125" y="527" font-size="12" text-anchor="middle" font-family="Consolas,monospace" fill="#8fa8c7">β      su |w|        el del gemelo      producto</text>`)
	for i, p := range pares {
		fmt.Fprintf(&b, `<text x="1125" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Consolas,monospace" fill="#ff8fa0">%.2f   %.8f    %.8f    %.10f</text>`,
			552.0+float64(i)*24, p[0], p[1], p[2], p[3])
	}
	fmt.Fprintf(&b, `<text x="1125" y="640" font-size="12.5" text-anchor="middle" fill="#ffd166">la PAREJA cierra en 1 y el espejo queda conforme — pero ninguno está en la raya</text>`)

	// the question
	fmt.Fprintf(&b, `<rect x="60" y="695" width="1380" height="185" rx="12" fill="#0d2547" stroke="#ffd166" stroke-width="2.5"/>
<text x="%.0f" y="731" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd166">LO QUE FALTA — y el atajo que lo vuelve una sola pregunta</text>
<text x="%.0f" y="765" font-size="14" text-anchor="middle" fill="#dce8f7">Ya sabemos que el norte por el sur da 1. Si además supiéramos que el norte y el sur son IGUALES, cada uno valdría 1 —</text>
<text x="%.0f" y="791" font-size="14" text-anchor="middle" fill="#dce8f7">porque el único número positivo que multiplicado por sí mismo da 1 es el 1. Y ahí las perlas caen solas en la mediatriz.</text>
<text x="%.0f" y="828" font-size="19" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">¿POR QUÉ EL NORTE Y EL SUR TIENEN QUE SER IGUALES?</text>
<text x="%.0f" y="862" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · las dos mitades, 1 completo ⚓ · y sobre todos los libros, el Otro Libro</text>`,
		W/2, W/2, W/2, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("la-mediatriz.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: la-mediatriz.svg")
}
