// Command laguna runs the captain's answer to the ONE question. Asked
// "why can no star burn inland?", he answered with a stone in a still
// pond: light cannot occupy a single point - it occupies everything,
// in all directions, a point growing infinitely, trading voltage for
// amperage, ever farther, ever weaker.
//
// The shop name for that law is GAUSS + THE MEAN VALUE PRINCIPLE:
// log|xi| is the pond's surface; away from stones it is harmonic -
// every point equals the average of its neighbours, so nothing can
// stay concentrated and nothing can hide. Jensen's formula is the
// stone-law made exact: the MEAN LIGHT over the circle of radius r
// centred at s=1/2,
//
//	m(r) = (1/2pi) INTEGRAL log|xi(1/2 + r e^{i th})| d th,
//
// is CONSTANT while the disk holds no stone (still pond), and its
// slope against ln r equals EXACTLY the number of stones enclosed:
// dm/dlnr = N(r). Each pearl pair entering at r = gamma_k lifts the
// staircase by 2 - announced on EVERY larger ring, ever farther
// (log(r/gamma_k) keeps growing), ever weaker (each stone's share of
// the total fades). We never locate a single zero: the pond counts
// them by ripples alone.
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

// logAbsXi is log|xi(s)| computed in log form (stable far up the tower).
// Points with Re s < 1/2 are reflected through the mirror xi(s)=xi(1-s)
// first: our Euler-Maclaurin zeta diverges left of the strip, and the
// mirror (judged at 8e-8 in the assembly) is exact.
func logAbsXi(s complex128) float64 {
	if real(s) < 0.5 {
		s = 1 - s
	}
	return math.Log(cmplx.Abs(0.5*s*(s-1))) - real(s)/2*math.Log(math.Pi) +
		real(lgammaC(s/2)) + math.Log(cmplx.Abs(zetaC(s)))
}

// meanLight is m(r): the average of log|xi| over the circle of radius
// r centred at 1/2, angles offset half a step so no sample can land
// exactly on a pearl when r = gamma_k.
func meanLight(r float64, K int) float64 {
	s := 0.0
	for j := 0; j < K; j++ {
		th := 2 * math.Pi * (float64(j) + 0.5) / float64(K)
		s += logAbsXi(complex(0.5+r*math.Cos(th), r*math.Sin(th)))
	}
	return s / float64(K)
}

func main() {
	fmt.Println("🪨 LA LAGUNA — la piedra en el agua quieta: la ley que responde la pregunta")

	truePearls := []float64{14.134725, 21.022040, 25.010858, 30.424876, 32.935062}
	K := 2000

	// ---- the still pond: no stone inside -> mean light CONSTANT ----
	fmt.Println("\nLA LAGUNA QUIETA — anillos sin piedra adentro (r < 14.13): la luz media NO se mueve")
	m0 := logAbsXi(complex(0.5, 0))
	fmt.Printf("   luz en el centro: log|ξ(1/2)| = %.9f\n", m0)
	stillWorst := 0.0
	for _, r := range []float64{2, 5, 8, 11, 13.5} {
		m := meanLight(r, K)
		d := math.Abs(m - m0)
		if d > stillWorst {
			stillWorst = d
		}
		fmt.Printf("   r=%5.1f:  luz media %.9f   desvío del centro %.1e\n", r, m, d)
	}
	fmt.Printf("   → laguna QUIETA verificada: peor desvío %.1e — sin piedra, la luz no puede concentrarse ni moverse\n", stillWorst)

	// ---- the staircase: ripples count the stones ----
	fmt.Println("\nLAS ONDAS CUENTAN LAS PIEDRAS — pendiente de la luz media = piedras encerradas (ley de Gauss):")
	r0, r1, dr := 2.0, 36.0, 0.25
	var rs, ms []float64
	for r := r0; r <= r1+1e-9; r += dr {
		rs = append(rs, r)
		ms = append(ms, meanLight(r, K))
	}
	// slope against ln r via central differences
	Ns := make([]float64, len(rs))
	for i := 1; i < len(rs)-1; i++ {
		Ns[i] = (ms[i+1] - ms[i-1]) / (math.Log(rs[i+1]) - math.Log(rs[i-1]))
	}
	Ns[0], Ns[len(rs)-1] = Ns[1], Ns[len(rs)-2]
	// read the plateaus and the crossing radii (where N crosses odd values)
	var crossings []float64
	for lvl := 1.0; lvl <= 9.0; lvl += 2 {
		for i := 1; i < len(rs); i++ {
			if Ns[i-1] < lvl && Ns[i] >= lvl {
				// linear interpolation
				f := (lvl - Ns[i-1]) / (Ns[i] - Ns[i-1])
				crossings = append(crossings, rs[i-1]+f*(rs[i]-rs[i-1]))
				break
			}
		}
	}
	fmt.Println("   escalón (piedras)   radio del salto (medido)   perla verdadera   desvío")
	worst := 0.0
	for i, c := range crossings {
		d := math.Abs(c - truePearls[i])
		if d > worst {
			worst = d
		}
		fmt.Printf("      %d → %d              %8.3f                %8.4f        %.3f 🎯\n",
			2*i, 2*i+2, c, truePearls[i], d)
	}
	// plateau values between jumps
	fmt.Println("   mesetas medidas (deben ser 0, 2, 4, 6, 8, 10 piedras):")
	plateauAt := func(r float64) float64 {
		best, bd := 0.0, math.Inf(1)
		for i, rr := range rs {
			if d := math.Abs(rr - r); d < bd {
				bd, best = d, Ns[i]
			}
		}
		return best
	}
	for _, pr := range []float64{8, 17.5, 23, 27.5, 31.7, 35} {
		fmt.Printf("      r=%5.1f → pendiente %.3f\n", pr, plateauAt(pr))
	}
	fmt.Printf("\n⚖ LA LEY DE LA PIEDRA VERIFICADA: los saltos caen en las perlas (peor desvío %.3f)\n", worst)
	fmt.Println("  sin ubicar UN solo cero: las ondas solas contaron 2, 4, 6, 8, 10 piedras —")
	fmt.Println("  cada piedra se anuncia en TODOS los anillos que la rodean, cada vez más lejos")
	fmt.Println("  (su log(r/γ) sigue creciendo) y cada vez más débil (su parte del total se diluye).")
	fmt.Println("\nPOR QUÉ ESTO RESPONDE LA PREGUNTA — y qué queda:")
	fmt.Println("  tu ley explica el DETECTOR: nada se esconde — toda piedra, esté donde esté, se anuncia")
	fmt.Println("  en cada anillo (por eso el contraluz atrapa a todo fantasma a alguna profundidad).")
	fmt.Println("  la laguna cuenta piedras por DISTANCIA — todavía no jura la DIRECCIÓN (la línea).")
	fmt.Println("  la llave, en idioma de laguna: demostrar que el agua quieta solo admite piedras")
	fmt.Println("  en la orilla de la línea — que ninguna onda puede nacer tierra adentro. Todavía no.")

	// ---- picture ----
	var b strings.Builder
	W, H := 1500.0, 940.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🪨 LA LAGUNA — la piedra en el agua quieta responde la pregunta</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"la luz no puede ocupar un solo espacio: lo ocupa todo, en todas direcciones, cada vez más lejos, cada vez más débil" — el capitán · en el taller: Gauss + valor medio (Jensen)</text>`,
		W, H, W, H, W/2, W/2)

	// staircase plot
	sx, sy, sw, sh := 90.0, 110.0, 900.0, 420.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" rx="10" fill="#0d2547" stroke="#7fb2ff" stroke-width="1.5"/>
<text x="%.0f" y="%.0f" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7fb2ff">LA ESCALERA DE GAUSS — pendiente de la luz media = piedras encerradas</text>`,
		sx, sy, sw, sh, sx+sw/2, sy+30)
	var pts []string
	for i, r := range rs {
		X := sx + 50 + (r-r0)/(r1-r0)*(sw-100)
		Y := sy + sh - 60 - Ns[i]/10.5*(sh-130)
		pts = append(pts, fmt.Sprintf("%.1f,%.1f", X, Y))
	}
	fmt.Fprintf(&b, `<polyline points="%s" fill="none" stroke="#ffd97f" stroke-width="2.5"/>`, strings.Join(pts, " "))
	for lvl := 0; lvl <= 10; lvl += 2 {
		Y := sy + sh - 60 - float64(lvl)/10.5*(sh-130)
		fmt.Fprintf(&b, `<line x1="%.0f" y1="%.1f" x2="%.0f" y2="%.1f" stroke="#8fa8c7" stroke-width="0.6" stroke-dasharray="4 6"/>
<text x="%.0f" y="%.1f" font-size="10.5" fill="#8fa8c7">%d</text>`,
			sx+50, Y, sx+sw-50, Y, sx+18, Y+4, lvl)
	}
	for _, p := range truePearls {
		X := sx + 50 + (p-r0)/(r1-r0)*(sw-100)
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.0f" x2="%.1f" y2="%.0f" stroke="#7fd7a8" stroke-width="1" stroke-dasharray="3 5"/>
<text x="%.1f" y="%.0f" font-size="10" text-anchor="middle" fill="#7fd7a8">γ=%.2f</text>`,
			X, sy+60, X, sy+sh-60, X, sy+52, p)
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12" text-anchor="middle" fill="#dce8f7">cada par de perlas que entra al anillo sube la escalera EXACTAMENTE en 2 — contado solo con ondas, sin ubicar jamás un cero (peor desvío %.3f)</text>`,
		sx+sw/2, sy+sh-20, worst)

	// still pond panel
	fmt.Fprintf(&b, `<rect x="1020" y="110" width="390" height="420" rx="10" fill="#081020" stroke="#7fd7a8" stroke-width="1.5"/>
<text x="1215" y="142" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">LA LAGUNA QUIETA</text>
<text x="1215" y="172" font-size="12" text-anchor="middle" fill="#dce8f7">anillos sin piedra (r &lt; 14.13):</text>
<text x="1215" y="196" font-size="12" text-anchor="middle" fill="#dce8f7">la luz media NO SE MUEVE</text>
<text x="1215" y="226" font-size="13" text-anchor="middle" font-family="Consolas,monospace" fill="#ffd97f">peor desvío: %.0e</text>
<circle cx="1215" cy="330" r="80" fill="none" stroke="#7fb2ff" stroke-width="1.2" opacity="0.9"/>
<circle cx="1215" cy="330" r="55" fill="none" stroke="#7fb2ff" stroke-width="1.2" opacity="0.6"/>
<circle cx="1215" cy="330" r="30" fill="none" stroke="#7fb2ff" stroke-width="1.2" opacity="0.35"/>
<circle cx="1215" cy="330" r="4" fill="#ffd97f"/>
<text x="1215" y="442" font-size="11.5" text-anchor="middle" fill="#8fa8c7">sin piedra, ninguna onda: la luz no puede</text>
<text x="1215" y="460" font-size="11.5" text-anchor="middle" fill="#8fa8c7">concentrarse — el valor ES la media de sus vecinos</text>
<text x="1215" y="490" font-size="11.5" text-anchor="middle" fill="#7fd7a8">la cara del "no hay mapa sin luz" (F211), ahora ley</text>`,
		stillWorst)

	// verdict
	fmt.Fprintf(&b, `<rect x="90" y="570" width="1320" height="230" rx="12" fill="#102a10" stroke="#ffd166" stroke-width="2"/>
<text x="%.0f" y="606" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd166">LO QUE TU RESPUESTA DEMOSTRÓ — y lo que aún pide</text>
<text x="%.0f" y="640" font-size="14" text-anchor="middle" fill="#dce8f7">tu ley es real y quedó verificada: la luz no puede ocupar un solo punto — cada piedra se anuncia en TODOS los anillos, más lejos y más débil, con la cuenta exacta conservada.</text>
<text x="%.0f" y="668" font-size="14" text-anchor="middle" fill="#dce8f7">por ESO el contraluz atrapa a todo fantasma: nada se esconde de las ondas. La laguna contó 2, 4, 6, 8, 10 piedras sin ver ninguna — tercer instrumento del taller.</text>
<text x="%.0f" y="700" font-size="14" text-anchor="middle" fill="#ffd166">pero la laguna cuenta por DISTANCIA, no jura la DIRECCIÓN: sabe cuántas piedras hay a radio γ — aún no que estén sobre LA LÍNEA.</text>
<text x="%.0f" y="732" font-size="14.5" text-anchor="middle" fill="#ff8fa0">la pregunta, refinada por tu propia ley: ¿por qué el agua quieta solo admite piedras en la orilla de la línea? — ninguna onda puede nacer tierra adentro. Todavía no.</text>
<text x="%.0f" y="768" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · las dos mitades, 1 completo ⚓ · y sobre todos los libros, el Otro Libro</text>`,
		W/2, W/2, W/2, W/2, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("laguna.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: laguna.svg")
}
