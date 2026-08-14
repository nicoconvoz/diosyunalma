// Command elinsignia takes the flagship number lambda_1 = 0.023096 and finds
// the captain's perfect 1/2 relation, harmonised at dimension 0.
//
// HIS ORDER: "take the 5, put our flagship number 0.023096, and find the
// perfect 1/2 relation harmonising it to leave it perfect in dimension 0."
//
// THE FLAGSHIP HAS THREE FACES, AND ALL THREE ARE THE HALF:
//
// FACE 1 - ARITHMETIC (the primes' side):
//
//		lambda_1 = 1 + (gamma - ln 4pi)/2
//	  The half is the harmoniser, literally: one plus HALF of (gamma - ln 4pi).
//	  And gamma comes from the sieve alone (F285).
//
// FACE 2 - GEOMETRIC (dimension 0): each pearl pair contributes
//
//		2*Re(1/rho) = 2*beta/(beta^2+gamma^2)
//	  and ON the cable beta = 1/2 makes the numerator EXACTLY 1:
//		term = 1/((1/2)^2 + gamma^2)
//	  Unity above, the half SQUARED below. That is the perfect form.
//
// FACE 3 - HARMONIC (the perfect square): on the cable, the same term equals
//
//		4*sin^2(phi/2)
//	  the HALF-ANGLE square of F280. A square cannot be negative - "perfecto"
//	  in both senses of the word.
//
// AND THE THREE FACES ARE ONE: on the cable, 1/(1/4+gamma^2) = 4 sin^2(phi/2)
// exactly (verified per pearl below), and their infinite sum is Face 1's
// arithmetic value. RH <=> the three faces coincide for the WHOLE necklace.
//
// Reproduce: go run ./cmd/elinsignia
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

const gammaE = 0.5772156649015329

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

func perlas(hasta float64) []float64 {
	var ps []float64
	prevT, prevZ := 12.0, zOf(12.0)
	for t := 12.02; t <= hasta; t += 0.02 {
		z := zOf(t)
		if z*prevZ < 0 {
			a, c := prevT, t
			for i := 0; i < 60; i++ {
				m := (a + c) / 2
				if zOf(m)*prevZ < 0 {
					c = m
				} else {
					a = m
				}
			}
			ps = append(ps, (a+c)/2)
		}
		prevT, prevZ = t, z
	}
	return ps
}

func main() {
	fmt.Println("🏅 EL INSIGNIA — λ₁ = 0,023096 y su ½ perfecto, en la dimensión 0")
	fmt.Println("\n   El número insignia tiene TRES caras, y las tres son el ½.")

	exacto := 1 + (gammaE-math.Log(4*math.Pi))/2

	// ---- CARA 1 ----
	fmt.Println("\nCARA 1 · ARITMÉTICA — el ½ como ARMONIZADOR, del lado de los primos")
	fmt.Println("\n        λ₁ = 1 + (γ − ln 4π)/2")
	fmt.Printf("\n        γ − ln 4π = %.9f   (un número negativo, feo)\n", gammaE-math.Log(4*math.Pi))
	fmt.Printf("        SU MITAD  = %.9f\n", (gammaE-math.Log(4*math.Pi))/2)
	fmt.Printf("        1 + mitad = %.9f  ← EL INSIGNIA\n", exacto)
	fmt.Println("\n   ⟹ El insignia es, literalmente, **uno más LA MITAD de (γ − ln 4π)**. El ½")
	fmt.Println("   armoniza el número feo con el 1 — y la γ sale de la criba sola (F285).")

	// ---- CARA 2 y 3 ----
	fmt.Println("\nCARAS 2 y 3 · GEOMÉTRICA Y ARMÓNICA — el ½ al cuadrado, y el cuadrado perfecto")
	fmt.Println("\n   Cada par de perlas aporta 2·Re(1/ρ) = 2β/(β²+γ²). Y SOBRE EL CABLE, β = ½")
	fmt.Println("   hace el numerador EXACTAMENTE 1:")
	fmt.Println("\n        aporte = 1 / (½² + γ²)         ← unidad arriba, el ½ AL CUADRADO abajo")
	fmt.Println("\n   Y por F280, ese mismo aporte es el CUADRADO PERFECTO del ángulo mitad:")
	fmt.Println("\n        aporte = 4·sen²(φ/2)           ← un cuadrado: no puede ser negativo")
	fmt.Println("\n   Verificado perla por perla que LAS DOS FORMAS SON UNA:")
	ps := perlas(120)
	fmt.Printf("\nperlas: %d\n", len(ps))
	fmt.Println("\n        γ            1/(¼+γ²)          4·sen²(φ/2)       diferencia")
	peor := 0.0
	var suma float64
	for i, g := range ps {
		f1 := 1 / (0.25 + g*g)
		w := 1 - 1/complex(0.5, g)
		phi := cmplx.Phase(w)
		f2 := 4 * math.Pow(math.Sin(phi/2), 2)
		d := math.Abs(f1 - f2)
		if d > peor {
			peor = d
		}
		suma += f1
		if i < 4 {
			fmt.Printf("   %12.6f %17.12f %17.12f %13.1e\n", g, f1, f2, d)
		}
	}
	fmt.Printf("\n        peor diferencia entre las dos formas ...... %.1e\n", peor)
	fmt.Println("        ⟹ **el ½ al cuadrado en el denominador Y el ángulo mitad al cuadrado")
	fmt.Println("        son EL MISMO NÚMERO. Las dos caras del ½, soldadas.**")

	// ---- la sintesis ----
	fmt.Println("\nLA SÍNTESIS · EL INSIGNIA ARMONIZADO, PERFECTO EN LA DIMENSIÓN 0")
	fmt.Printf("\n        suma de los cuadrados perfectos (38 perlas) ... %.9f\n", suma)
	fmt.Printf("        el insignia entero (cara aritmética) .......... %.9f\n", exacto)
	fmt.Printf("        lo que aportan las perlas que faltan .......... %.9f\n", exacto-suma)
	fmt.Println("\n        y cada perla faltante aporta OTRO cuadrado perfecto: positivo seguro")

	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("⚡ **EL ½ PERFECTO DEL INSIGNIA, ENCONTRADO — Y SON TRES CARAS DE UNO:**")
	fmt.Println("\n  · como ARMONIZADOR:   λ₁ = 1 + (γ − ln 4π)/2 — la mitad que cose el 1")
	fmt.Println("    con el número feo, y la γ viene de los primos solos")
	fmt.Println("  · como CUADRADO en el denominador:  1/(½² + γ²) — la unidad arriba,")
	fmt.Println("    porque 2β = 2·½ = 1 exacto sobre el cable")
	fmt.Printf("  · como ÁNGULO MITAD:  4·sen²(φ/2) — un cuadrado perfecto, idéntico al\n")
	fmt.Printf("    anterior a %.0e, que no puede ser negativo\n", peor)
	fmt.Println("\n⟹ RH, dicho con el insignia: **el número 0,023096 es una suma de infinitos")
	fmt.Println("  cuadrados perfectos del ½ — todos, sin que ninguno se rompa jamás.**")
	fmt.Println("\n⚖️ Honesto: las tres caras son álgebra clásica (la fórmula de λ₁, el ángulo")
	fmt.Println("  mitad, Mertens). La soldadura de las tres en un solo número, medida con")
	fmt.Println("  nuestros instrumentos, es el cierre — no un avance. Todavía no.")

	escribirLamina(exacto, suma, peor, len(ps))
}

func escribirLamina(exacto, suma, peor float64, n int) {
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="1400" height="640" viewBox="0 0 1400 640">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="700" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🏅 EL INSIGNIA — λ₁ = 0,023096 y su ½ perfecto en la dimensión 0</text>
<text x="700" y="76" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">tres caras, y las tres son el ½ — soldadas en un solo número, medido</text>
<rect x="50" y="110" width="420" height="250" rx="12" fill="#0f2b22" stroke="#2f7f63"/>
<text x="260" y="144" font-size="15" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">EL ½ ARMONIZADOR (los primos)</text>
<text x="260" y="196" font-size="19" text-anchor="middle" font-family="monospace" fill="#ffd98a">λ₁ = 1 + (γ − ln 4π)/2</text>
<text x="80" y="240" font-size="13.5" font-family="Georgia" fill="#cfe6ff">uno más LA MITAD del número feo —</text>
<text x="80" y="262" font-size="13.5" font-family="Georgia" fill="#cfe6ff">y la γ sale de la criba sola (F285)</text>
<text x="260" y="310" font-size="16" text-anchor="middle" font-family="monospace" fill="#7ee0c0">= %.9f</text>
<rect x="490" y="110" width="420" height="250" rx="12" fill="#101f36" stroke="#26456e"/>
<text x="700" y="144" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">EL ½ AL CUADRADO (dimensión 0)</text>
<text x="700" y="196" font-size="19" text-anchor="middle" font-family="monospace" fill="#ffd98a">aporte = 1 / (½² + γ²)</text>
<text x="520" y="240" font-size="13.5" font-family="Georgia" fill="#cfe6ff">la unidad arriba porque 2β = 2·½ = 1</text>
<text x="520" y="262" font-size="13.5" font-family="Georgia" fill="#cfe6ff">exacto sobre el cable — el ½ hace el 1</text>
<text x="700" y="310" font-size="14" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">una fracción limpia por perla</text>
<rect x="930" y="110" width="420" height="250" rx="12" fill="#161a3a" stroke="#5a4fa8"/>
<text x="1140" y="144" font-size="15" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">EL ÁNGULO MITAD (el cuadrado)</text>
<text x="1140" y="196" font-size="19" text-anchor="middle" font-family="monospace" fill="#ffd98a">aporte = 4·sen²(φ/2)</text>
<text x="960" y="240" font-size="13.5" font-family="Georgia" fill="#cfe6ff">un cuadrado perfecto: no puede ser</text>
<text x="960" y="262" font-size="13.5" font-family="Georgia" fill="#cfe6ff">negativo jamás</text>
<text x="1140" y="310" font-size="14" text-anchor="middle" font-family="monospace" fill="#7ee0c0">idéntico al anterior a %.0e</text>
<text x="700" y="430" font-size="22" text-anchor="middle" font-family="Georgia" fill="#ffd98a">RH, dicho con el insignia: 0,023096 es una suma de infinitos</text>
<text x="700" y="462" font-size="22" text-anchor="middle" font-family="Georgia" fill="#ffd98a">cuadrados perfectos del ½ — sin que ninguno se rompa jamás</text>
<text x="700" y="510" font-size="14" text-anchor="middle" font-family="monospace" fill="#cfe6ff">%d perlas nuestras suman %.9f · el resto del collar aporta %.9f — cuadrado tras cuadrado</text>
<text x="700" y="560" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">las tres caras son álgebra clásica; la soldadura en un solo número, medida acá, es el cierre — no un avance</text>
<text x="700" y="600" font-size="14" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Todavía no.</text>
</svg>
`, exacto, peor, n, suma, exacto-suma)
	os.WriteFile("el-insignia.svg", []byte(b.String()), 0o644)
	fmt.Println("\n🖼️  lámina escrita: el-insignia.svg")
}
