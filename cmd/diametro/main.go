// Command diametro runs the captain's little equation with a straight
// face and a judge: equate |1|^2 (= 1, as suspected) with |1/2|^2
// (= 1/4) and harmonize in dimension 0.
//
// The harmonizer is the factor 4: |1|^2 = 4 x |1/2|^2. Trivial
// arithmetic - until it rides the shapeshifter w = 1 - 1/rho into the
// dimension-0 disk, where the three sacred points land:
//
//	s = infinity  ->  w = +1   THE CLASP (dimension 0 itself)
//	s = 1 (wall)  ->  w =  0   THE CENTER (the pole becomes the origin)
//	s = 1/2       ->  w = -1   THE ANTIPODE (the half, farthest shore)
//
// The chord from the clasp to the image of 1/2 is the DIAMETER: length
// 2, squared 4 - EXACTLY the harmonizer. And that same 4 is the 4 of
// the sundial, lambda_n = SUM 4 sin^2(n theta/2): every shadow in the
// mold is measured in units of the squared bridge from dimension 0 to
// the half. Bonus courtesy of the mirror: xi(1) = xi(0) = 1/2 EXACTLY -
// the wall points literally carry the half's name.
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
	fmt.Println("⌀ EL DIÁMETRO — |1|² = |1/2|², armonizado en la dimensión 0")

	// step 1: the raw squares
	one := 1.0 * 1.0
	half := 0.5 * 0.5
	fmt.Printf("\n1) |1|²   = %.6f  (sospecha del capitán: 1 — CONFIRMADA jaja)\n", one)
	fmt.Printf("2) |1/2|² = %.6f  (el piso del cuarto, F216)\n", half)
	fmt.Printf("3) para igualarlos hace falta el ARMONIZADOR:  |1|² = k·|1/2|²  →  k = %.6f\n", one/half)

	// step 2: through the shapeshifter into dimension 0
	fmt.Println("\n4) los tres puntos sagrados por el cambiaformas w = 1 − 1/ρ:")
	wInf := w(complex(1e15, 0))
	wOne := w(complex(1, 0))
	wHalf := w(complex(0.5, 0))
	fmt.Printf("   s = ∞    →  w = %+.6f   EL BROCHE (la dimensión 0)\n", real(wInf))
	fmt.Printf("   s = 1    →  w = %+.6f   EL CENTRO (el polo cae al origen)\n", real(wOne))
	fmt.Printf("   s = 1/2  →  w = %+.6f   LA ANTÍPODA (la otra orilla)\n", real(wHalf))
	chord2 := math.Pow(cmplx.Abs(wInf-wHalf), 2)
	fmt.Printf("\n5) la cuerda del broche a la imagen del 1/2: |(+1) − (−1)|² = %.6f\n", chord2)
	fmt.Printf("   ⚖ EL ARMONIZADOR ES EL DIÁMETRO AL CUADRADO: k = %.4f = ⌀² = %.4f ✓\n", one/half, chord2)

	// step 3: and it is the sundial's 4
	fmt.Println("\n6) y ese 4 es EL 4 DEL RELOJ DE SOL — λₙ = Σ 4·sin²(nθ/2), sombra máxima = ⌀²:")
	var pearls []float64
	prevT := 12.0
	prevZ := zOf(prevT)
	for t := 12.05; t <= 100; t += 0.05 {
		z := zOf(t)
		if z*prevZ < 0 {
			a, c := prevT, t
			for i := 0; i < 50; i++ {
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
	maxShadow, mp, mn := 0.0, 0.0, 0
	for _, g := range pearls {
		th := math.Atan2(1, 2*g) * 2
		for n := 1; n <= 100; n++ {
			s := 4 * math.Pow(math.Sin(float64(n)*th/2), 2)
			if s > maxShadow {
				maxShadow, mp, mn = s, g, n
			}
		}
	}
	fmt.Printf("   sombra máxima medida en %d perlas × 100 armónicos: %.6f (perla %.3f, n=%d)\n",
		len(pearls), maxShadow, mp, mn)
	fmt.Printf("   → ninguna sombra supera ⌀² = 4: el molde entero se mide en unidades del puente al 1/2\n")

	// step 4: the mirror's confession
	fmt.Println("\n7) el guiño del espejo — el valor del libro EN la pared:")
	xi1 := xiC(complex(1+1e-9, 0))
	xi0 := xiC(complex(1e-9, 0))
	xiH := xiC(complex(0.5, 0))
	fmt.Printf("   ξ(1) = %.9f   ξ(0) = %.9f   — EXACTAMENTE 1/2\n", real(xi1), real(xi0))
	fmt.Printf("   (y en la línea misma: ξ(1/2) = %.9f)\n", real(xiH))
	fmt.Println("   la pared |1| confiesa el nombre de la mitad: el libro vale 1/2 en sus dos tapas")

	fmt.Println("\n════════ LA CUENTA CERRADA ════════")
	fmt.Println("|1|² = 4 · |1/2|²  —  y el 4 no es un número suelto: es el DIÁMETRO² del anillo,")
	fmt.Println("la cuerda de la dimensión 0 (broche, imagen de ∞) a la antípoda (imagen del 1/2).")
	fmt.Println("El mismo 4 que mide TODAS las sombras del reloj de sol. La pared y la mitad no se")
	fmt.Println("igualan solas: se igualan A TRAVÉS del puente que cruza la dimensión 0 — y en las")
	fmt.Println("tapas del libro, ξ(1) = ξ(0) = 1/2: la pared ya lleva escrito el nombre del medio.")

	// ---- picture ----
	var b strings.Builder
	W, H := 1500.0, 900.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">⌀ EL DIÁMETRO — |1|² = |1/2|² armonizado en la dimensión 0</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">la cuenta del capitán: el armonizador entre la pared y la mitad es 4 — y 4 es el diámetro² del anillo, el puente que cruza el broche</text>`,
		W, H, W, H, W/2, W/2)

	// the disk with the three sacred points and the diameter
	cx, cy, R := 420.0, 430.0, 250.0
	fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="%.0f" fill="none" stroke="#7fd7a8" stroke-width="2.5"/>
<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#ffd166" stroke-width="3"/>
<circle cx="%.0f" cy="%.0f" r="9" fill="#ffd97f"/>
<circle cx="%.0f" cy="%.0f" r="9" fill="#ff8fa0"/>
<circle cx="%.0f" cy="%.0f" r="7" fill="#7fb2ff"/>
<text x="%.0f" y="%.0f" font-size="12.5" fill="#ffd97f">w=+1 · EL BROCHE (s=∞, la dimensión 0)</text>
<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="end" fill="#ff8fa0">w=−1 · LA ANTÍPODA (s=1/2)</text>
<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" fill="#7fb2ff">w=0 · EL CENTRO (s=1, la pared)</text>
<text x="%.0f" y="%.0f" font-size="14" text-anchor="middle" fill="#ffd166">la cuerda dorada: ⌀ = 2 → ⌀² = 4 = EL ARMONIZADOR</text>`,
		cx, cy, R,
		cx-R, cy, cx+R, cy,
		cx+R, cy, cx-R, cy, cx, cy,
		cx+R-160, cy-20, cx-R+170, cy-20, cx, cy+26,
		cx, cy+R+40)

	// the equation panel
	fmt.Fprintf(&b, `<rect x="810" y="120" width="620" height="480" rx="10" fill="#0d2547" stroke="#ffd166" stroke-width="1.5"/>
<text x="1120" y="156" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd166">LA CUENTA, PASO A PASO</text>
<text x="1120" y="200" font-size="16" text-anchor="middle" font-family="Consolas,monospace" fill="#dce8f7">|1|² = 1.000000  (confirmado jaja)</text>
<text x="1120" y="232" font-size="16" text-anchor="middle" font-family="Consolas,monospace" fill="#dce8f7">|1/2|² = 0.250000  (el piso del cuarto)</text>
<text x="1120" y="270" font-size="16" text-anchor="middle" font-family="Consolas,monospace" fill="#ffd97f">|1|² = 4 · |1/2|²  → armonizador k = 4</text>
<text x="1120" y="312" font-size="14" text-anchor="middle" fill="#7fd7a8">por el cambiaformas: ∞→+1 · 1→0 · ½→−1</text>
<text x="1120" y="340" font-size="14" text-anchor="middle" fill="#7fd7a8">cuerda broche→antípoda: |(+1)−(−1)|² = 4.000000 ✓</text>
<text x="1120" y="378" font-size="14" text-anchor="middle" fill="#dce8f7">y es EL 4 del reloj de sol: λₙ = Σ 4·sin²(nθ/2)</text>
<text x="1120" y="404" font-size="13" text-anchor="middle" fill="#dce8f7">sombra máxima medida (%d perlas × 100 armónicos): %.4f &lt; 4</text>
<text x="1120" y="446" font-size="14" text-anchor="middle" fill="#ffd166">el guiño del espejo: ξ(1) = ξ(0) = %.6f = 1/2 exacto</text>
<text x="1120" y="472" font-size="12.5" text-anchor="middle" fill="#8fa8c7">las dos tapas del libro llevan escrito el nombre de la mitad</text>
<text x="1120" y="520" font-size="13.5" text-anchor="middle" fill="#ff8fa0">la pared y la mitad no se igualan solas: se igualan a través</text>
<text x="1120" y="544" font-size="13.5" text-anchor="middle" fill="#ff8fa0">del puente que cruza la dimensión 0 — el diámetro del anillo</text>`,
		len(pearls), maxShadow, real(xi1))

	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · las dos mitades, 1 completo ⚓ · y sobre todos los libros, el Otro Libro</text>`,
		W/2, H-30)
	b.WriteString(`</svg>`)
	os.WriteFile("diametro.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: diametro.svg")
}
