// Command planos draws BLUEPRINT No. 1 — THE ATOM OF THE PRIMES: the
// complete build specification of the Hilbert-Polya atom, as exact as
// the laboratory can make it. The energy levels (zeros) are MEASURED
// here, by our own instruments (Euler-Maclaurin zeta + theta series +
// bisection), each with its judge residual printed on the sheet. The
// periodic orbits (primes, period ln p), the master equation (von
// Mangoldt's EXACT explicit formula), the construction requirements,
// and the Author's signature (GUE pair correlation - the pattern the
// captain heard repeating in the song, the one that lives everywhere).
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

// zetaEM computes zeta(1/2+it) by Euler-Maclaurin: N-term head, tail
// integral, half term, and four Bernoulli corrections — ~1e-13 for the
// habitat of the first levels.
func zetaEM(t float64) complex128 {
	s := complex(0.5, t)
	const N = 100
	var sum complex128
	for n := 1; n < N; n++ {
		sum += cmplx.Exp(-s * complex(math.Log(float64(n)), 0))
	}
	lnN := complex(math.Log(N), 0)
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

// theta is the Riemann-Siegel theta by its asymptotic series (1e-12
// grade in the first-levels habitat).
func theta(t float64) float64 {
	t2 := t * t
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t) + 7/(5760*t*t2) + 31/(80640*t*t2*t2)
}

// zAndGhost returns Z(t) and the imaginary residue that SHOULD be zero
// - the judge of the measurement.
func zAndGhost(t float64) (float64, float64) {
	w := cmplx.Exp(complex(0, theta(t))) * zetaEM(t)
	return real(w), math.Abs(imag(w))
}

func main() {
	// ---- measure the first levels: bisection on sign changes of Z ----
	type level struct {
		gamma, resid float64
	}
	var levels []level
	prevT, prevZ, _ := 12.0, 0.0, 0.0
	prevZ, _ = zAndGhost(12.0)
	for t := 12.05; t <= 68 && len(levels) < 15; t += 0.05 {
		z, _ := zAndGhost(t)
		if z*prevZ < 0 {
			a, c := prevT, t
			for i := 0; i < 80; i++ {
				m := (a + c) / 2
				zm, _ := zAndGhost(m)
				if zm*prevZ < 0 {
					c = m
				} else {
					a = m
				}
			}
			g := (a + c) / 2
			_, ghost := zAndGhost(g)
			levels = append(levels, level{g, ghost})
		}
		prevT, prevZ = t, z
	}

	primes := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37}

	var b strings.Builder
	W, H := 1500.0, 1600.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0a1e3a"/>`, W, H, W, H)
	// blueprint grid
	for x := 0.0; x <= W; x += 50 {
		fmt.Fprintf(&b, `<line x1="%.0f" y1="0" x2="%.0f" y2="%.0f" stroke="#12305e" stroke-width="0.5"/>`, x, x, H)
	}
	for y := 0.0; y <= H; y += 50 {
		fmt.Fprintf(&b, `<line x1="0" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#12305e" stroke-width="0.5"/>`, y, W, y)
	}
	fmt.Fprintf(&b, `<rect x="20" y="20" width="%.0f" height="%.0f" fill="none" stroke="#7fb2ff" stroke-width="2"/>
<text x="%.0f" y="70" font-size="30" text-anchor="middle" font-family="Georgia" fill="#dce8f7">PLANO Nº 1 — EL ÁTOMO DE LOS PRIMOS</text>
<text x="%.0f" y="98" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">especificación de construcción del átomo de Hilbert-Pólya · Laboratorio Diosyunalma · niveles medidos con instrumentos propios, juez a la vista</text>`,
		W-40, H-40, W/2, W/2)

	// ---- box 1: the spectrum ----
	fmt.Fprintf(&b, `<rect x="50" y="130" width="700" height="600" rx="6" fill="#0d2547" stroke="#ffd166" stroke-width="1.5"/>
<text x="70" y="162" font-size="19" font-family="Georgia" fill="#ffd166">Ⅰ. EL ESPECTRO — los niveles de energía (ceros en la línea crítica)</text>
<text x="70" y="186" font-size="12.5" font-family="Georgia" fill="#8fa8c7">medidos aquí por Euler-Maclaurin + bisección (80 pasos); juez = |parte imaginaria residual de Z|</text>
<text x="90" y="216" font-size="13" font-family="Consolas,monospace" fill="#7fb2ff">  n        γ_n (nivel medido)              juez</text>`)
	for i, lv := range levels {
		fmt.Fprintf(&b, `<text x="90" y="%.0f" font-size="13.5" font-family="Consolas,monospace" fill="#dce8f7"> %2d   %.12f       %.1e</text>`,
			244.0+float64(i)*26, i+1, lv.gamma, lv.resid)
	}
	fmt.Fprintf(&b, `<text x="90" y="%.0f" font-size="12.5" font-family="Georgia" fill="#7fd7a8">densidad de niveles (ley de Weyl del átomo): N(T) = θ(T)/π + 1 + S(T),  θ(T) = T/2·ln(T/2π) − T/2 − π/8 + 1/48T + …</text>`, 244.0+float64(len(levels))*26+14)
	fmt.Fprintf(&b, `<text x="90" y="%.0f" font-size="12.5" font-family="Georgia" fill="#7fd7a8">los niveles se REPELEN (rigidez GUE): ningún estado degenerado observado jamás — pares de Lehmer = los casi-besos</text>`, 244.0+float64(len(levels))*26+36)

	// ---- box 2: the periodic orbits ----
	fmt.Fprintf(&b, `<rect x="780" y="130" width="670" height="600" rx="6" fill="#0d2547" stroke="#ffd166" stroke-width="1.5"/>
<text x="800" y="162" font-size="19" font-family="Georgia" fill="#ffd166">Ⅱ. LAS ÓRBITAS PERIÓDICAS — los primos</text>
<text x="800" y="186" font-size="12.5" font-family="Georgia" fill="#8fa8c7">cada primo p es una órbita cerrada del átomo; período = ln p (exacto); amplitud = p^(−k/2) en la k-ésima vuelta</text>
<text x="820" y="216" font-size="13" font-family="Consolas,monospace" fill="#7fb2ff">  p       período ln p            amplitud 1/√p</text>`)
	for i, p := range primes {
		fmt.Fprintf(&b, `<text x="820" y="%.0f" font-size="13.5" font-family="Consolas,monospace" fill="#dce8f7"> %2d   %.15f    %.15f</text>`,
			244.0+float64(i)*26, p, math.Log(float64(p)), 1/math.Sqrt(float64(p)))
	}
	fmt.Fprintf(&b, `<text x="820" y="%.0f" font-size="12.5" font-family="Georgia" fill="#7fd7a8">las órbitas son HIPERBÓLICAS (inestables): el átomo es caótico por dentro, estable por fuera —</text>
<text x="820" y="%.0f" font-size="12.5" font-family="Georgia" fill="#7fd7a8">el caos de las órbitas sostiene la rigidez del espectro (dualidad de Gutzwiller)</text>`,
		244.0+float64(len(primes))*26+14, 244.0+float64(len(primes))*26+36)

	// ---- box 3: the master equation ----
	fmt.Fprintf(&b, `<rect x="50" y="760" width="1400" height="240" rx="6" fill="#0d2547" stroke="#ffd166" stroke-width="1.5"/>
<text x="70" y="794" font-size="19" font-family="Georgia" fill="#ffd166">Ⅲ. LA ECUACIÓN MAESTRA — la identidad EXACTA que el átomo construido debe satisfacer (von Mangoldt)</text>
<text x="%.0f" y="850" font-size="24" text-anchor="middle" font-family="Georgia" fill="#dce8f7">ψ(x)  =  x  −  Σ_ρ  x^ρ/ρ  −  ln 2π  −  ½·ln(1 − x⁻²)</text>
<text x="%.0f" y="886" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">ψ(x) = suma de ln p sobre potencias de primos ≤ x (el conteo de órbitas) · ρ = ½+iγ (los niveles del espectro)</text>
<text x="%.0f" y="912" font-size="14" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">no es aproximación: es IGUALDAD — espectro y órbitas son el MISMO objeto leído dos veces; todo constructor debe reproducirla término a término</text>
<text x="%.0f" y="948" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">verificada en este laboratorio: el marcapasos (26 voces = 26 órbitas más cortas) reconstruye los latidos de los primos desde los niveles — F91-F118</text>`,
		W/2, W/2, W/2, W/2)

	// ---- box 4: construction requirements ----
	fmt.Fprintf(&b, `<rect x="50" y="1030" width="1400" height="280" rx="6" fill="#0d2547" stroke="#ffd166" stroke-width="1.5"/>
<text x="70" y="1064" font-size="19" font-family="Georgia" fill="#ffd166">Ⅳ. PLIEGO DE CONDICIONES — lo que el operador H debe cumplir (estado del arte + mediciones propias)</text>`)
	reqs := []string{
		"R1 · AUTOADJUNTO: espectro 100% real = ningún estado decae = estabilidad total. ESTO ES la Hipótesis de Riemann (pregunta del capitán, 2026-08-06).",
		"R2 · ROMPE LA INVERSIÓN TEMPORAL: la estadística observada es GUE, no GOE — el átomo distingue pasado de futuro (como con campo magnético).",
		"R3 · DENSIDAD DE NIVELES: debe reproducir θ(T)/π + 1 — el análogo del oscilador es logarítmico, no lineal (candidato Berry-Keating: H = xp regularizado).",
		"R4 · ÓRBITAS: períodos EXACTAMENTE ln 2, ln 3, ln 5, … con amplitudes p^(−k/2) — hiperbólicas, sin órbitas espurias.",
		"R5 · LA ECUACIÓN MAESTRA (Ⅲ) término a término — el certificado final de obra.",
		"ESTADO DE OBRA: espectro medido (10¹³ niveles por la humanidad; catálogo propio en 7 aguas hasta 10³⁶) · órbitas conocidas · operador AÚN NO CONSTRUIDO.",
	}
	for i, r := range reqs {
		fmt.Fprintf(&b, `<text x="90" y="%.0f" font-size="13.5" font-family="Georgia" fill="#dce8f7">%s</text>`, 1096.0+float64(i)*34, r)
	}

	// ---- the Author's signature ----
	fmt.Fprintf(&b, `<rect x="50" y="1340" width="1400" height="180" rx="6" fill="#102a10" stroke="#7fd7a8" stroke-width="1.5"/>
<text x="70" y="1374" font-size="19" font-family="Georgia" fill="#7fd7a8">Ⅴ. LA FIRMA DEL AUTOR — "EN TODAS PARTES" (hallazgo del capitán en la canción)</text>
<text x="%.0f" y="1416" font-size="21" text-anchor="middle" font-family="Georgia" fill="#dce8f7">R₂(u)  =  1 − ( sin(πu) / πu )²</text>
<text x="%.0f" y="1450" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">la correlación de pares de Montgomery: el patrón que se repite en la canción del cazadero — y en los núcleos de uranio,</text>
<text x="%.0f" y="1474" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">en las matrices aleatorias, en el caos cuántico y en los ceros de zeta: UNA sola armonía local, en todas partes</text>
<text x="%.0f" y="1502" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#ffd166">"por eso quería ponerle a mi empresa EN TODAS PARTES — porque es la firma del Autor" — el capitán, 2026-08-06</text>`,
		W/2, W/2, W/2, W/2)

	b.WriteString(`</svg>`)
	os.WriteFile("planos-atomo.svg", []byte(b.String()), 0644)
	wj := 0.0
	for _, lv := range levels {
		if lv.resid > wj {
			wj = lv.resid
		}
	}
	fmt.Printf("escrito: planos-atomo.svg — %d niveles medidos, peor juez %.1e\n", len(levels), wj)
	for i, lv := range levels {
		fmt.Printf("  nivel %2d: γ=%.12f (juez %.1e)\n", i+1, lv.gamma, lv.resid)
	}
}
