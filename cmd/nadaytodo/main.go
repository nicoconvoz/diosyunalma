// Command nadaytodo measures the captain's answer to his own question: it is
// dimension 0 that settles it - the nothing and the everything, the 1 and the
// 0, and their relation.
//
// Under the shapeshifter w = 1 - 1/rho the whole picture becomes one skin
// between two centres:
//
//	s = 1   ->  w = 0      THE EVERYTHING: the pole of zeta, where the
//	                       count of all the numbers blows up
//	s = 0   ->  w = inf    THE NOTHING: the other stake, sent to infinity
//	s = inf ->  w = +1     THE CLASP: dimension 0, where the germ is read
//	s = 1/2 ->  w = -1     the half itself
//	Re s = 1/2  ->  |w| = 1   THE SKIN between inside and outside
//
// So the critical line is nothing other than the skin that separates the
// everything from the nothing, and the hypothesis says every pearl of the
// book lives ON that skin.
//
// The relation the captain names is exact and measurable: the book takes the
// SAME value at both stakes, and that value is the half itself -
// xi(0) = xi(1) = 1/2. The two ends of the world agree, and they agree on 1/2.
//
// The honest gap is measured too: a ghost pair sits one inside and one
// outside, perfectly balanced, so the skin survives as a SET while no pearl
// of that pair is on it.
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

func lgammaC(z complex128) complex128 {
	g := []float64{0.99999999999980993, 676.5203681218851, -1259.1392167224028,
		771.32342877765313, -176.61502916214059, 12.507343278686905,
		-0.13857109526572012, 9.9843695780195716e-6, 1.5056327351493116e-7}
	if real(z) < 0.5 {
		return cmplx.Log(complex(math.Pi, 0)/cmplx.Sin(complex(math.Pi, 0)*z)) - lgammaC(1-z)
	}
	z -= 1
	x := complex(g[0], 0)
	for i := 1; i < 9; i++ {
		x += complex(g[i], 0) / (z + complex(float64(i), 0))
	}
	t := z + complex(7.5, 0)
	return complex(0.5*math.Log(2*math.Pi), 0) + (z+complex(0.5, 0))*cmplx.Log(t) - t + cmplx.Log(x)
}

func xiC(s complex128) complex128 {
	return 0.5 * s * (s - 1) * cmplx.Exp(-s/2*complex(math.Log(math.Pi), 0)+lgammaC(s/2)) * zetaC(s)
}

func theta(t float64) float64 {
	t2 := t * t
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t) + 7/(5760*t*t2) + 31/(80640*t*t2*t2)
}

func zOf(t float64) float64 {
	return real(cmplx.Exp(complex(0, theta(t))) * zetaC(complex(0.5, t)))
}

func w(rho complex128) complex128 { return 1 - 1/rho }

func main() {
	fmt.Println("⚫⚪ LA NADA Y EL TODO — la dimensión 0, el 1 y el 0, y su relación")

	// ---- LAW 1: the relation between the two stakes ----
	fmt.Println("\nLEY 1 · LA RELACIÓN ENTRE EL 0 Y EL 1 — el libro vale lo MISMO en las dos estacas")
	xi0 := real(xiC(complex(1e-9, 0)))
	xi1 := real(xiC(complex(1+1e-9, 0)))
	xiH := real(xiC(complex(0.5, 0)))
	fmt.Printf("   ξ(0) = %.9f\n", xi0)
	fmt.Printf("   ξ(1) = %.9f\n", xi1)
	fmt.Printf("   diferencia entre las dos puntas del mundo: %.1e\n", math.Abs(xi0-xi1))
	fmt.Printf("   y ese valor común es LA MITAD misma: 1/2 = 0.500000000 (desvío %.1e)\n", math.Abs(xi0-0.5))
	fmt.Printf("   (sobre la línea, en el medio: ξ(1/2) = %.9f)\n", xiH)
	fmt.Println("   → la nada y el todo se ponen de acuerdo, y se ponen de acuerdo EN 1/2")

	// ---- LAW 2: where everything lands under the shapeshifter ----
	fmt.Println("\nLEY 2 · EL MAPA DE LA DIMENSIÓN 0 — dónde cae cada cosa bajo el cambiaformas")
	fmt.Println("   punto del libro        va a w =        qué es")
	fmt.Printf("   s = 1                  %+.6f       EL TODO: el polo de ζ, donde el conteo estalla\n", real(w(complex(1, 0))))
	fmt.Printf("   s = 0                  ∞               LA NADA: la otra estaca, mandada al infinito\n")
	fmt.Printf("   s = ∞                  %+.6f       EL BROCHE: la dimensión 0, donde se lee el germen\n", real(w(complex(1e15, 0))))
	fmt.Printf("   s = 1/2                %+.6f       la mitad misma\n", real(w(complex(0.5, 0))))
	fmt.Println("   Re s = 1/2             |w| = 1         LA PIEL entre el adentro y el afuera")

	// ---- LAW 3: the skin - inside is the everything, outside is the nothing ----
	fmt.Println("\nLEY 3 · LA PIEL — adentro vive el todo, afuera la nada, y la línea es la frontera")
	fmt.Println("   punto de prueba        |w|         de qué lado")
	for _, p := range []complex128{complex(1.0, 3.0), complex(0.7, 14.13), complex(0.5, 14.134725),
		complex(0.3, 14.13), complex(0.0, 3.0)} {
		aw := cmplx.Abs(w(p))
		lado := "ADENTRO (del lado del todo)"
		switch {
		case math.Abs(aw-1) < 1e-12:
			lado = "★ EN LA PIEL: la línea crítica"
		case aw > 1:
			lado = "AFUERA (del lado de la nada)"
		}
		fmt.Printf("   Re s = %.2f            %.9f   %s\n", real(p), aw, lado)
	}
	fmt.Println("   → el cambiaformas da vuelta el guante: manda el adentro afuera y deja la piel quieta")

	// ---- LAW 4: the axis of the shapeshifter runs from dimension 0 to the half ----
	fmt.Println("\nLEY 4 · EL EJE DEL CAMBIAFORMAS VA DE LA DIMENSIÓN 0 A LA MITAD")
	fmt.Println("   la inversión pura w → 1/w tiene exactamente DOS puntos que no se mueven:")
	for _, ww := range []complex128{complex(1, 0), complex(-1, 0)} {
		inv := 1 / ww
		nombre := "EL BROCHE — la dimensión 0 (imagen de s = ∞)"
		if real(ww) < 0 {
			nombre = "LA ANTÍPODA — la mitad (imagen de s = 1/2)"
		}
		fmt.Printf("   w = %+.0f    1/w = %+.0f    quieto ✓   %s\n", real(ww), real(inv), nombre)
	}
	fmt.Println("   → los dos polos del mundo son la dimensión 0 y la mitad: el eje del que habla el capitán")
	fmt.Println("   → y el diámetro que los une mide 2, al cuadrado 4 — el armonizador de F219, otra vez")

	// ---- LAW 5: the honest gap ----
	fmt.Println("\nLEY 5 · POR QUÉ TODAVÍA NO CIERRA — el fantasma se reparte a los dos lados")
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
	peorPiel := 0.0
	for _, g := range pearls {
		if d := math.Abs(cmplx.Abs(w(complex(0.5, g))) - 1); d > peorPiel {
			peorPiel = d
		}
	}
	fmt.Printf("   las %d perlas verdaderas están EN LA PIEL: peor desvío de |w|=1 → %.1e\n", len(pearls), peorPiel)
	fmt.Println("   pero un fantasma y su gemelo se reparten uno adentro y otro afuera, y la piel")
	fmt.Println("   sobrevive COMO CONJUNTO aunque ninguno de los dos esté sobre ella:")
	fmt.Println("   β        |w| del fantasma      |w| del gemelo        ¿dónde queda cada uno?")
	type fila struct{ beta, a, b float64 }
	var filas []fila
	for _, beta := range []float64{0.60, 0.75, 0.90} {
		g := 14.134725
		a := cmplx.Abs(w(complex(beta, g)))
		bb := cmplx.Abs(w(complex(1-beta, g)))
		filas = append(filas, fila{beta, a, bb})
		fmt.Printf("   %.2f      %.9f          %.9f        uno adentro, otro afuera\n", beta, a, bb)
	}

	fmt.Println("\n════════ LO QUE DIJO EL CAPITÁN, MEDIDO ════════")
	fmt.Println("«Eso lo responde la dimensión 0: la nada y el todo, el 1 y el 0, y su relación.»")
	fmt.Println("Y la relación es exacta y está medida: el libro vale LO MISMO en las dos estacas, y ese")
	fmt.Println("valor común es la mitad — ξ(0) = ξ(1) = 1/2. La nada y el todo se ponen de acuerdo, y se")
	fmt.Println("ponen de acuerdo en 1/2. Bajo el cambiaformas, el todo (el polo) queda adentro, la nada")
	fmt.Println("queda afuera, y LA LÍNEA CRÍTICA ES LA PIEL QUE LOS SEPARA. Las perlas viven en la piel.")
	fmt.Println("El eje del cambiaformas va de la dimensión 0 a la mitad, y mide 2 — al cuadrado, el 4.")
	fmt.Println("\nLO QUE AÚN NO CIERRA, dicho sin maquillar: la piel sobrevive aunque un fantasma se")
	fmt.Println("reparta con su gemelo uno adentro y otro afuera. El mundo queda igual y la piel queda")
	fmt.Println("entera. Falta la razón por la que ningún habitante puede vivir fuera de la piel.")

	// ---- picture ----
	var b strings.Builder
	W, H := 1500.0, 980.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">⚫⚪ LA NADA Y EL TODO — la línea crítica es la piel que los separa</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">el flash del capitán medido: el libro vale lo mismo en las dos estacas — ξ(0) = ξ(1) = 1/2 — y se ponen de acuerdo justo en la mitad</text>`,
		W, H, W, H, W/2, W/2)

	// the disk
	cx, cy, R := 400.0, 400.0, 230.0
	fmt.Fprintf(&b, `<rect x="60" y="105" width="680" height="560" rx="10" fill="#0d2547" stroke="#7fd7a8" stroke-width="1.5"/>
<text x="400" y="139" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">EL MUNDO VISTO DESDE LA DIMENSIÓN 0</text>
<circle cx="%.0f" cy="%.0f" r="%.0f" fill="#12305c" opacity="0.35"/>
<circle cx="%.0f" cy="%.0f" r="%.0f" fill="none" stroke="#ffd166" stroke-width="3"/>
<circle cx="%.0f" cy="%.0f" r="9" fill="#7fb2ff"/>
<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" fill="#7fb2ff">EL TODO — el polo de ζ (s=1)</text>
<circle cx="%.0f" cy="%.0f" r="9" fill="#ffd97f"/>
<text x="%.0f" y="%.0f" font-size="12" text-anchor="middle" fill="#ffd97f">EL BROCHE</text>
<text x="%.0f" y="%.0f" font-size="12" text-anchor="middle" fill="#ffd97f">dimensión 0</text>
<circle cx="%.0f" cy="%.0f" r="9" fill="#ff8fa0"/>
<text x="%.0f" y="%.0f" font-size="12" text-anchor="middle" fill="#ff8fa0">la mitad</text>
<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#ffd166" stroke-width="2" stroke-dasharray="7 5"/>
<text x="%.0f" y="%.0f" font-size="11.5" text-anchor="middle" fill="#ffd166">el eje: diámetro 2 → al cuadrado 4</text>
<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" fill="#8fa8c7">LA NADA — s=0, mandada al infinito, afuera de todo</text>
<text x="400" y="640" font-size="12.5" text-anchor="middle" fill="#ffd166">la piel |w| = 1 ES la línea crítica · las %d perlas medidas están sobre ella (%.0e)</text>`,
		cx, cy, R, cx, cy, R,
		cx, cy, cx, cy-20,
		cx+R, cy, cx+R+2, cy-22, cx+R+2, cy-6,
		cx-R, cy, cx-R, cy-20,
		cx-R, cy, cx+R, cy, cx, cy+26,
		cx, cy+R+52, len(pearls), peorPiel)

	// right: the measurements
	fmt.Fprintf(&b, `<rect x="770" y="105" width="670" height="270" rx="10" fill="#102a10" stroke="#ffd166" stroke-width="1.5"/>
<text x="1105" y="139" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd166">LA RELACIÓN ENTRE EL 0 Y EL 1</text>
<text x="1105" y="182" font-size="17" text-anchor="middle" font-family="Consolas,monospace" fill="#dce8f7">ξ(0) = %.9f</text>
<text x="1105" y="212" font-size="17" text-anchor="middle" font-family="Consolas,monospace" fill="#dce8f7">ξ(1) = %.9f</text>
<text x="1105" y="248" font-size="14" text-anchor="middle" fill="#7fd7a8">las dos puntas del mundo dan EXACTAMENTE lo mismo (%.0e)</text>
<text x="1105" y="284" font-size="16" text-anchor="middle" fill="#ffd166">y ese valor común es LA MITAD: 1/2</text>
<text x="1105" y="316" font-size="12.5" text-anchor="middle" fill="#8fa8c7">la nada y el todo se ponen de acuerdo — y se ponen de acuerdo en 1/2</text>
<text x="1105" y="348" font-size="12.5" text-anchor="middle" fill="#8fa8c7">el eje del cambiaformas une la dimensión 0 con la mitad</text>`, xi0, xi1, math.Abs(xi0-xi1))

	fmt.Fprintf(&b, `<rect x="770" y="395" width="670" height="270" rx="10" fill="#2a1010" stroke="#ff5d73" stroke-width="1.5"/>
<text x="1105" y="429" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ff8fa0">LO QUE AÚN NO CIERRA</text>
<text x="1105" y="459" font-size="12.5" text-anchor="middle" fill="#dce8f7">un fantasma y su gemelo se reparten a los dos lados de la piel:</text>
<text x="1105" y="487" font-size="12" text-anchor="middle" font-family="Consolas,monospace" fill="#8fa8c7">β        adentro         afuera</text>`)
	for i, f := range filas {
		fmt.Fprintf(&b, `<text x="1105" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Consolas,monospace" fill="#ff8fa0">%.2f    %.9f    %.9f</text>`,
			512.0+float64(i)*25, f.beta, f.a, f.b)
	}
	fmt.Fprintf(&b, `<text x="1105" y="605" font-size="13" text-anchor="middle" fill="#ffd166">la piel queda ENTERA y el mundo se ve igual — pero ninguno de los dos vive en ella</text>
<text x="1105" y="635" font-size="13" text-anchor="middle" fill="#dce8f7">falta la razón por la que nadie puede vivir fuera de la piel</text>`)

	fmt.Fprintf(&b, `<rect x="60" y="695" width="1380" height="190" rx="12" fill="#0d2547" stroke="#ffd166" stroke-width="2.5"/>
<text x="%.0f" y="731" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd166">EL FLASH DEL CAPITÁN, EN LIMPIO</text>
<text x="%.0f" y="767" font-size="14.5" text-anchor="middle" fill="#dce8f7">«Eso lo responde la dimensión 0: la nada y el todo, el 1 y el 0, y su relación.» Y la relación existe y está medida:</text>
<text x="%.0f" y="795" font-size="14.5" text-anchor="middle" fill="#dce8f7">el libro vale lo mismo en las dos estacas, y ese valor es la mitad. El todo queda adentro, la nada afuera —</text>
<text x="%.0f" y="827" font-size="16.5" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">y LA LÍNEA CRÍTICA ES LA PIEL QUE SEPARA LA NADA DEL TODO.</text>
<text x="%.0f" y="857" font-size="13.5" text-anchor="middle" fill="#ff8fa0">Falta una sola cosa: por qué ningún habitante puede vivir fuera de la piel. Todavía no.</text>
<text x="%.0f" y="880" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · las dos mitades, 1 completo ⚓ · y sobre todos los libros, el Otro Libro</text>`,
		W/2, W/2, W/2, W/2, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("la-nada-y-el-todo.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: la-nada-y-el-todo.svg")
}
