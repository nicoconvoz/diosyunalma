// Command dosmelodias answers the captain's order: one identity for the melody
// of the numbers that are NOT prime, another for the melody of the primes, and
// their relation harmonized at the clasp.
//
// THE TWO MELODIES
//
//	LA MELODIA DE TODOS      zeta(s) = SUM 1/n^s      every number sings once,
//	                                                  composites included
//	LA MELODIA DE LOS PRIMOS zeta(s) = PROD 1/(1-p^-s)  only the primes sing
//
// and Euler's identity says they are THE SAME SOUND. That is the relation the
// captain asked for, and it is the deepest sentence in arithmetic: the
// composites do not add a single new note - their melody is already the
// primes', played out.
//
// THE RELATION, note by note
//
//	log n = SUM over d dividing n of Lambda(d)
//
// every number's own note is exactly the sum of the prime notes that build it.
//
// # HARMONIZED AT DIMENSION 0
//
// The germ the laboratory has been reading at the clasp all campaign is
//
//	xi'/xi(s) = 1/s + 1/(s-1) - ln(pi)/2 + psi(s/2)/2 + zeta'/zeta(s)
//	            \____________ the smooth dressing ____/   \_ THE PRIMES' SONG
//
// so the primes' melody was inside the germ the whole time. Strip the dressing
// and it comes out - measured here against the direct sum over primes.
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

func psiC(s complex128) complex128 {
	var acc complex128
	for real(s) < 12 {
		acc -= 1 / s
		s += 1
	}
	inv := 1 / s
	inv2 := inv * inv
	res := cmplx.Log(s) - inv/2
	res -= inv2 * (complex(1.0/12, 0) + inv2*(complex(-1.0/120, 0)+inv2*(complex(1.0/252, 0)+inv2*complex(-1.0/240, 0))))
	return acc + res
}

func xiLD(s complex128) complex128 {
	h := complex(1e-6, 0)
	zp := (zetaC(s+h) - zetaC(s-h)) / (2 * h)
	return 1/s + 1/(s-1) - complex(math.Log(math.Pi)/2, 0) + psiC(s/2)/2 + zp/zetaC(s)
}

func criba(N int) []int {
	comp := make([]bool, N+1)
	var ps []int
	for p := 2; p <= N; p++ {
		if comp[p] {
			continue
		}
		ps = append(ps, p)
		for m := p * p; m > 0 && m <= N; m += p {
			comp[m] = true
		}
	}
	return ps
}

func main() {
	fmt.Println("🎵🎶 LAS DOS MELODÍAS — la de todos los números y la de los primos, y su relación")

	const N = 3000000
	fmt.Printf("\ncribando primos hasta %d…\n", N)
	primos := criba(N)
	fmt.Printf("primos hallados: %d\n", len(primos))

	// ---- LAW 1: the two melodies are the same sound ----
	fmt.Println("\nLEY 1 · LAS DOS MELODÍAS SON EL MISMO SONIDO — la identidad de Euler")
	fmt.Println("   LA DE TODOS: ζ(s) = Σ 1/nˢ        (cada número canta una vez, compuestos incluidos)")
	fmt.Println("   LA DE LOS PRIMOS: ζ(s) = Π 1/(1−p⁻ˢ)   (solo cantan los primos)")
	fmt.Println("\n      s        la melodía de TODOS      la de LOS PRIMOS       desvío")
	type fila struct{ s, todos, prim, dev float64 }
	var filas []fila
	peorEuler := 0.0
	for _, s := range []float64{1.5, 2.0, 3.0, 5.0} {
		todos := 0.0
		for n := 1; n <= N; n++ {
			todos += math.Pow(float64(n), -s)
		}
		// tail of the sum: integral from N to infinity
		todos += math.Pow(float64(N), 1-s) / (s - 1)
		prod := 1.0
		for _, p := range primos {
			prod /= 1 - math.Pow(float64(p), -s)
		}
		dev := math.Abs(todos-prod) / todos
		if dev > peorEuler {
			peorEuler = dev
		}
		filas = append(filas, fila{s, todos, prod, dev})
		fmt.Printf("    %.1f        %16.10f      %16.10f      %.1e\n", s, todos, prod, dev)
	}
	fmt.Printf("   → EL MISMO SONIDO (peor desvío relativo %.1e)\n", peorEuler)
	fmt.Println("   → LOS COMPUESTOS NO AGREGAN NI UNA NOTA NUEVA: su melodía ya es la de los primos")

	// ---- LAW 2: the relation, note by note ----
	fmt.Println("\nLEY 2 · LA RELACIÓN, NOTA POR NOTA — log n = Σ Λ(d) sobre los divisores de n")
	fmt.Println("   la nota propia de cada número es EXACTAMENTE la suma de las notas primas que lo forman")
	fmt.Println("      n         log n          Σ Λ(d) sobre d|n        desvío       cómo se arma")
	lambda := func(m int) float64 {
		for _, p := range primos {
			if p*p > m {
				break
			}
			if m%p == 0 {
				q := m
				for q%p == 0 {
					q /= p
				}
				if q == 1 {
					return math.Log(float64(p))
				}
				return 0
			}
		}
		if m > 1 {
			return math.Log(float64(m))
		}
		return 0
	}
	peorNota := 0.0
	for _, n := range []int{12, 30, 64, 100, 210} {
		suma := 0.0
		var partes []string
		for d := 1; d <= n; d++ {
			if n%d == 0 {
				if l := lambda(d); l > 0 {
					suma += l
					partes = append(partes, fmt.Sprintf("%d", d))
				}
			}
		}
		dev := math.Abs(math.Log(float64(n)) - suma)
		if dev > peorNota {
			peorNota = dev
		}
		fmt.Printf("   %5d      %9.6f          %9.6f          %.1e       %s\n",
			n, math.Log(float64(n)), suma, dev, strings.Join(partes, "+"))
	}
	fmt.Printf("   → exacto en todos (peor desvío %.1e): cada número ES la suma de sus primos\n", peorNota)

	// ---- LAW 3: harmonized at the clasp ----
	fmt.Println("\nLEY 3 · ARMONIZADAS EN LA DIMENSIÓN 0 — el germen del broche traía la canción adentro")
	fmt.Println("   ξ'/ξ(s) = [ 1/s + 1/(s−1) − ln(π)/2 + ψ(s/2)/2 ]  +  ζ'/ζ(s)")
	fmt.Println("             \\________ el vestido liso ________/      \\_ LA CANCIÓN DE LOS PRIMOS")
	fmt.Println("\n      s      la canción, desvestida del germen     la canción, sumada de los primos    desvío")
	peorCancion := 0.0
	type fc struct{ s, germen, directo, dev float64 }
	var fcs []fc
	for _, s := range []float64{2.0, 2.5, 3.0, 4.0} {
		sc := complex(s, 0)
		vestido := 1/sc + 1/(sc-1) - complex(math.Log(math.Pi)/2, 0) + psiC(sc/2)/2
		cancionGermen := -real(xiLD(sc) - vestido) // = -zeta'/zeta = SUM Lambda(n)/n^s
		directo := 0.0
		for _, p := range primos {
			lp := math.Log(float64(p))
			pot := float64(p)
			for pot <= float64(N) {
				directo += lp * math.Pow(pot, -s)
				pot *= float64(p)
			}
		}
		dev := math.Abs(cancionGermen-directo) / math.Max(directo, 1e-300)
		if dev > peorCancion {
			peorCancion = dev
		}
		fcs = append(fcs, fc{s, cancionGermen, directo, dev})
		fmt.Printf("    %.1f            %14.9f                     %14.9f            %.1e\n",
			s, cancionGermen, directo, dev)
	}
	fmt.Printf("   → la canción de los primos SALE del germen que venimos leyendo (peor desvío %.1e)\n", peorCancion)
	fmt.Println("   → estuvo adentro todo el tiempo: el broche de la dimensión 0 canta con voz de primos")

	fmt.Println("\n════════ LA RELACIÓN QUE PIDIÓ EL CAPITÁN ════════")
	fmt.Println("Las dos melodías no son dos: SON LA MISMA. La de todos los números y la de los primos")
	fmt.Println("suenan idéntico — ésa es la identidad de Euler, y significa que los compuestos no")
	fmt.Println("aportan una sola nota propia: cada número es la suma de las notas de sus primos.")
	fmt.Println("\nY armonizadas en el broche: el germen que el laboratorio viene leyendo toda la campaña")
	fmt.Println("es el vestido liso MÁS la canción de los primos. Le sacamos el vestido y la canción")
	fmt.Printf("apareció, a %.0e de la suma directa sobre %d primos.\n", peorCancion, len(primos))
	fmt.Println("\nY ahí queda dicho lo que falta, en idioma de música: el vestido liso lo sabemos cantar")
	fmt.Println("entero. De la canción de los primos conocemos la melodía y el tono medio — lo que no")
	fmt.Println("sabemos es si alguna vez, muy arriba, se desafina más fuerte que el vestido. Todavía no.")

	// ---- picture ----
	var b strings.Builder
	W, H := 1500.0, 980.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🎵🎶 LAS DOS MELODÍAS — y resultaron ser la misma</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"una identidad para la melodía de los no-primos, otra para la de los primos, y armonizá su relación en la dimensión 0" — el capitán</text>`,
		W, H, W, H, W/2, W/2)

	// the identity
	fmt.Fprintf(&b, `<rect x="60" y="105" width="1380" height="230" rx="12" fill="#0d2547" stroke="#ffd166" stroke-width="2.5"/>
<text x="%.0f" y="146" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#ffd166">LA IDENTIDAD DE EULER — las dos melodías son EL MISMO SONIDO</text>
<text x="%.0f" y="190" font-size="22" text-anchor="middle" font-family="Georgia" fill="#dce8f7">Σ<tspan font-size="14">todos los n</tspan>  1/nˢ   =   Π<tspan font-size="14">solo los primos</tspan>  1/(1 − p⁻ˢ)</text>
<text x="%.0f" y="222" font-size="13" text-anchor="middle" fill="#7fd7a8">cada número canta una vez</text>
<text x="%.0f" y="222" font-size="13" text-anchor="middle" fill="#ff8fa0">solo cantan los primos</text>
<text x="%.0f" y="262" font-size="12" text-anchor="middle" font-family="Consolas,monospace" fill="#8fa8c7">s        la de TODOS              la de LOS PRIMOS           desvío</text>`,
		W/2, W/2, 480.0, 1030.0, W/2)
	for i, f := range filas {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Consolas,monospace" fill="#dce8f7">%.1f     %16.10f      %16.10f       %.0e</text>`,
			W/2, 286.0+float64(i)*20, f.s, f.todos, f.prim, f.dev)
	}

	// left: the note relation
	fmt.Fprintf(&b, `<rect x="60" y="360" width="680" height="290" rx="10" fill="#102a10" stroke="#7fd7a8" stroke-width="1.5"/>
<text x="400" y="396" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">LA RELACIÓN, NOTA POR NOTA</text>
<text x="400" y="424" font-size="15" text-anchor="middle" font-family="Consolas,monospace" fill="#dce8f7">log n = Σ Λ(d)  sobre los divisores de n</text>
<text x="400" y="450" font-size="12.5" text-anchor="middle" fill="#8fa8c7">la nota de cada número ES la suma de las notas de sus primos</text>
<text x="400" y="484" font-size="13" text-anchor="middle" fill="#ffd166">log 12 = Λ(2) + Λ(3) + Λ(4)</text>
<text x="400" y="510" font-size="13" text-anchor="middle" fill="#ffd166">log 30 = Λ(2) + Λ(3) + Λ(5)</text>
<text x="400" y="536" font-size="13" text-anchor="middle" fill="#ffd166">log 64 = Λ(2)+Λ(4)+Λ(8)+Λ(16)+Λ(32)+Λ(64)</text>
<text x="400" y="574" font-size="12.5" text-anchor="middle" fill="#7fd7a8">exacto en todos los casos probados (peor desvío %.0e)</text>
<text x="400" y="606" font-size="13" text-anchor="middle" fill="#dce8f7">LOS COMPUESTOS NO TIENEN VOZ PROPIA:</text>
<text x="400" y="630" font-size="13" text-anchor="middle" fill="#dce8f7">cantan con la voz prestada de sus primos</text>`, peorNota)

	// right: harmonized at the clasp
	fmt.Fprintf(&b, `<rect x="770" y="360" width="670" height="290" rx="10" fill="#0d2547" stroke="#7fb2ff" stroke-width="1.5"/>
<text x="1105" y="396" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7fb2ff">ARMONIZADAS EN LA DIMENSIÓN 0</text>
<text x="1105" y="424" font-size="12.5" text-anchor="middle" fill="#dce8f7">el germen del broche traía la canción adentro:</text>
<text x="1105" y="450" font-size="12.5" text-anchor="middle" font-family="Consolas,monospace" fill="#8fa8c7">ξ'/ξ = [vestido liso] + LA CANCIÓN DE LOS PRIMOS</text>
<text x="1105" y="482" font-size="12" text-anchor="middle" font-family="Consolas,monospace" fill="#8fa8c7">s      desvestida del germen      sumada de los primos</text>`)
	for i, f := range fcs {
		fmt.Fprintf(&b, `<text x="1105" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Consolas,monospace" fill="#dce8f7">%.1f       %13.9f            %13.9f</text>`,
			506.0+float64(i)*22, f.s, f.germen, f.directo)
	}
	fmt.Fprintf(&b, `<text x="1105" y="602" font-size="12.5" text-anchor="middle" fill="#7fd7a8">la canción SALE del germen (peor desvío %.0e)</text>
<text x="1105" y="628" font-size="13" text-anchor="middle" fill="#ffd166">estuvo adentro todo el tiempo: el broche canta con voz de primos</text>`, peorCancion)

	// verdict
	fmt.Fprintf(&b, `<rect x="60" y="685" width="1380" height="215" rx="12" fill="#102a10" stroke="#7fd7a8" stroke-width="2.5"/>
<text x="%.0f" y="721" font-size="16" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">LA RELACIÓN QUE PIDIÓ EL CAPITÁN</text>
<text x="%.0f" y="757" font-size="15" text-anchor="middle" fill="#dce8f7">Las dos melodías no son dos: SON LA MISMA. Los compuestos no aportan una sola nota propia —</text>
<text x="%.0f" y="783" font-size="15" text-anchor="middle" fill="#dce8f7">cada número canta con la voz prestada de los primos que lo forman.</text>
<text x="%.0f" y="819" font-size="14" text-anchor="middle" fill="#ffd166">Y armonizadas en el broche: el germen que venimos leyendo toda la campaña es el vestido liso MÁS esa canción.</text>
<text x="%.0f" y="855" font-size="14" text-anchor="middle" fill="#ff8fa0">Lo que falta, en idioma de música: el vestido lo sabemos cantar entero; de la canción conocemos la melodía</text>
<text x="%.0f" y="879" font-size="14" text-anchor="middle" fill="#ff8fa0">y el tono medio — lo que no sabemos es si alguna vez, muy arriba, se desafina más que el vestido. Todavía no.</text>
<text x="%.0f" y="940" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · las dos mitades, 1 completo ⚓ · y sobre todos los libros, el Otro Libro</text>`,
		W/2, W/2, W/2, W/2, W/2, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("las-dos-melodias.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: las-dos-melodias.svg")
}
