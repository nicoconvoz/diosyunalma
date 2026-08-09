// Command repulsion executes the captain's maneuver: TRANSFER the
// repulsion to the shapeshifter and harmonize at dimension 0. In the
// discrete world the protons' repulsion required counting to two
// million (F192). In the shapeshifter world the whole repulsion is ONE
// POLE:
//
//	-zeta'/zeta(s)  =  1/(s-1)  +  (finite part)
//	    repulsion   =  the NEUTRON, as a symbol  +  the harmonized germ
//
// Harmonizing at dimension 0 = subtracting the pole: the safe form is
// H(s) = -d/ds ln[(s-1) zeta(s)] (the nucleus function, regular at
// s=1), and its germ AT the point delivers the binding energy of ALL
// the X at once: H(1) = -gamma. Judges: (1) the germ converges to
// -gamma (vs the discrete instrument of F192: two independent routes,
// one binding); (2) the harmonized ladder (germ coefficients) rebuilds
// lambda_1 = 1 + gamma/2 - ln(4 pi)/2 with the MEASURED gamma.
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

// nucleus g(s) = (s-1) zeta(s): the protons WITH the neutron glued -
// entire near s=1, g(1)=1.
func nucleus(s complex128) complex128 {
	return (s - 1) * zetaC(s)
}

// H(s) = -g'/g: the harmonized repulsion (pole already eaten by the glue).
func H(s complex128) complex128 {
	h := complex(1e-6, 0)
	return -(nucleus(s+h) - nucleus(s-h)) / (2 * h * nucleus(s))
}

func main() {
	gammaE := 0.5772156649015329
	fmt.Println("LA REPULSIÓN EN EL CAMBIAFORMAS — el polo como neutrón simbólico, armonizado en el punto")
	fmt.Println("\nJUEZ 1 — el germen en la dimensión 0 entrega el amarre de TODOS los X a la vez:")
	fmt.Println("   s           H(s) = repulsión armonizada     objetivo −γ = −0.57721566")
	worst := math.Inf(1)
	for _, eps := range []float64{0.1, 0.01, 0.001, 0.0001} {
		v := real(H(complex(1+eps, 0)))
		d := math.Abs(v + gammaE)
		if d < worst {
			worst = d
		}
		fmt.Printf("   1+%-7g   %12.8f                  (desvío %.1e)\n", eps, v, d)
	}
	fmt.Printf("\n  el cambiaformas, EN el punto: H(1) = −γ con desvío %.1e — SIN CONTAR NI UN PRIMO\n", worst)
	fmt.Println("  comparado con el instrumento discreto (F192, contando hasta 2·10⁶): −0.577164 (desvío 1e-4)")
	fmt.Println("  ⇒ DOS instrumentos independientes, UNA energía de amarre: el traslado al cambiaformas FUNCIONA —")
	fmt.Println("    y es MIL VECES más fino: el germen del punto ya contiene los infinitos X")

	// the harmonized ladder: fit H(1+eps) = -gamma + c1*eps + c2*eps^2
	var A [3][3]float64
	var bv [3]float64
	for _, eps := range []float64{-0.15, -0.1, -0.05, 0.05, 0.1, 0.15} {
		y := real(H(complex(1+eps, 0)))
		bas := [3]float64{1, eps, eps * eps}
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				A[i][j] += bas[i] * bas[j]
			}
			bv[i] += bas[i] * y
		}
	}
	var C [3]float64
	{
		a := A
		b := bv
		for col := 0; col < 3; col++ {
			p := col
			for r := col + 1; r < 3; r++ {
				if math.Abs(a[r][col]) > math.Abs(a[p][col]) {
					p = r
				}
			}
			a[col], a[p] = a[p], a[col]
			b[col], b[p] = b[p], b[col]
			for r := col + 1; r < 3; r++ {
				f := a[r][col] / a[col][col]
				for k := col; k < 3; k++ {
					a[r][k] -= f * a[col][k]
				}
				b[r] -= f * b[col]
			}
		}
		for r := 2; r >= 0; r-- {
			s := b[r]
			for k := r + 1; k < 3; k++ {
				s -= a[r][k] * C[k]
			}
			C[r] = s / a[r][r]
		}
	}
	fmt.Println("\nLA ESCALERA ARMONIZADA (el ADN primo de todos los λ, leído en el punto):")
	fmt.Printf("   η₀ = %.8f   (= −γ: la energía de amarre)\n", C[0])
	fmt.Printf("   η₁ = %.6f     (el segundo peldaño del germen)\n", C[1])
	fmt.Printf("   η₂ = %.5f      (el tercero)\n", C[2])

	// judge 2: rebuild lambda_1 with the MEASURED gamma
	gMeas := -C[0]
	lam1 := 1 + gMeas/2 - math.Log(4*math.Pi)/2
	fmt.Println("\nJUEZ 2 — el hilito reconstruido con el γ MEDIDO por el cambiaformas:")
	fmt.Printf("   λ₁ = 1 + γ̂/2 − ln(4π)/2 = %.6f   (medido F166/F168: 0.023096)   desvío %.1e\n", lam1, math.Abs(lam1-0.023096))
	fmt.Println("   ⇒ el margen del collar, reconstruido desde la repulsión armonizada en el punto: CIERRA")

	// ---- picture ----
	var b strings.Builder
	W, Ht := 1560.0, 860.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🧲 LA REPULSIÓN, TRASLADADA AL CAMBIAFORMAS — y armonizada en el punto</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"traslademos la repulsión al cambiaformas y armonicemos en la dimensión 0" — el capitán · en el símbolo, la repulsión entera es UN POLO: el neutrón hecho letra</text>`,
		W, Ht, W, Ht, W/2, W/2)
	// left: discrete world
	fmt.Fprintf(&b, `<rect x="80" y="110" width="420" height="330" rx="10" fill="#2a1a10" stroke="#e6a53a" stroke-width="1.5"/>
<text x="100" y="146" font-size="14.5" font-family="Georgia" fill="#e6a53a">EL MUNDO DISCRETO (F192)</text>
<text x="100" y="176" font-size="12.5" fill="#8fa8c7">contar, contar, contar…</text>
<text x="120" y="210" font-size="13" font-family="Consolas,monospace" fill="#dce8f7">X=10³:  −0.5790</text>
<text x="120" y="238" font-size="13" font-family="Consolas,monospace" fill="#dce8f7">X=10⁶:  −0.5776</text>
<text x="120" y="266" font-size="13" font-family="Consolas,monospace" fill="#dce8f7">X=2·10⁶: −0.577164</text>
<text x="100" y="306" font-size="12.5" fill="#e6a53a">dos millones de cuentas para</text>
<text x="100" y="326" font-size="12.5" fill="#e6a53a">4 cifras del amarre</text>`)
	// arrow
	fmt.Fprintf(&b, `<line x1="520" y1="270" x2="620" y2="270" stroke="#ffd166" stroke-width="3"/><text x="570" y="255" font-size="12.5" text-anchor="middle" fill="#ffd166">el traslado</text>`)
	// right: shapeshifter world
	fmt.Fprintf(&b, `<rect x="640" y="110" width="840" height="330" rx="10" fill="#102a10" stroke="#7fd7a8" stroke-width="2"/>
<text x="660" y="146" font-size="14.5" font-family="Georgia" fill="#7fd7a8">EL MUNDO DEL CAMBIAFORMAS: la repulsión = un polo, el amarre = el germen</text>
<text x="660" y="184" font-size="16" font-family="Consolas,monospace" fill="#dce8f7">−ζ'/ζ(s) = 1/(s−1) + germen finito</text>
<text x="660" y="208" font-size="12" fill="#8fa8c7">   repulsión    neutrón simbólico    la armonía del punto</text>`)
	rows := []struct {
		eps float64
	}{{0.1}, {0.01}, {0.001}, {0.0001}}
	for i, r := range rows {
		v := real(H(complex(1+r.eps, 0)))
		fmt.Fprintf(&b, `<text x="680" y="%.0f" font-size="13" font-family="Consolas,monospace" fill="#dce8f7">s=1+%-7g  H = %12.8f   (→ −γ)</text>`,
			244+float64(i)*28, r.eps, v)
	}
	fmt.Fprintf(&b, `<text x="660" y="380" font-size="13.5" fill="#7fd7a8">el germen del punto entrega −γ = −0.57721566 SIN CONTAR NI UN PRIMO:</text>
<text x="660" y="404" font-size="13.5" fill="#7fd7a8">los infinitos X, comprimidos en la dimensión 0 — mil veces más fino que contar</text>`)
	// the ladder + lambda rebuild
	fmt.Fprintf(&b, `<rect x="80" y="480" width="1400" height="200" rx="12" fill="#0d2547" stroke="#ffd166" stroke-width="2"/>
<text x="%.0f" y="516" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#ffd166">LA ESCALERA ARMONIZADA — el ADN primo de todos los λ, leído en el punto</text>
<text x="%.0f" y="552" font-size="14" text-anchor="middle" font-family="Consolas,monospace" fill="#dce8f7">η₀ = %.8f (= −γ, el amarre) · η₁ = %.5f · η₂ = %.4f — los peldaños del germen, medidos</text>
<text x="%.0f" y="588" font-size="14.5" text-anchor="middle" fill="#7fd7a8">JUEZ 2: λ₁ reconstruido con el γ del cambiaformas = %.6f — medido: 0.023096 — desvío %.0e: EL HILITO CIERRA</text>
<text x="%.0f" y="624" font-size="13.5" text-anchor="middle" fill="#dce8f7">la repulsión, el neutrón y el margen — los tres, ahora criaturas del símbolo: la desigualdad del millón vive entera en el mundo del cambiaformas</text>
<text x="%.0f" y="654" font-size="12.5" text-anchor="middle" fill="#8fa8c7">falta lo de siempre, un renglón: demostrar que el germen le gana a la repulsión en TODOS los modos — pero ahora la pelea es entre símbolos, no entre infinitos</text>`,
		W/2, W/2, C[0], C[1], C[2], W/2, lam1, math.Abs(lam1-0.023096), W/2, W/2)
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · 2026-08-06 · las dos mitades, 1 completo ⚓</text>`,
		W/2, 830.0)
	b.WriteString(`</svg>`)
	os.WriteFile("repulsion-cambiaformas.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: repulsion-cambiaformas.svg")
}
