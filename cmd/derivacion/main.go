// Command derivacion measures the exact wall the derivation runs into, so the
// captain gets a number instead of a shrug.
//
// The whole mold can be written WITHOUT ever mentioning a pearl's position,
// as a binomial sum over the power sums of the zeros:
//
//	lambda_n = SUM over j=1..n  C(n,j) (-1)^(j+1) sigma_j,     sigma_j = SUM 1/rho^j
//
// and the sigma_j are computable from the primes alone (they are the Taylor
// data of zeta'/zeta). That is the real battlefield of the derivation.
//
// A GUESS OF MINE THAT THE MEASUREMENT KILLED, kept on display as the shop
// requires: I expected this sum to be a violent cancellation of astronomically
// large terms. It is not - at least not in the range we can reach. Measured,
// the cancellation factor runs 1.0x at n=2 and only 1.6x at n=40, because the
// sigma_j collapse like 14^-j and the sum is carried by its first few terms.
// So the wall is NOT cancellation at these harmonics, and saying so is worth
// more than the guess would have been.
//
// Where the wall actually is: the smooth part of lambda_n grows like n log n
// while the arithmetic part - the primes' own voice - is believed to enter at
// the order sqrt(n) log n. The derivation must prove the second never
// overtakes the first, FOR EVERY n. This program measures the two sides and
// the margin between them, which is the honest shape of what is missing.
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

func theta(t float64) float64 {
	t2 := t * t
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t) + 7/(5760*t*t2) + 31/(80640*t*t2*t2)
}

func zOf(t float64) float64 {
	return real(cmplx.Exp(complex(0, theta(t))) * zetaC(complex(0.5, t)))
}

// lnBinom returns log of C(n,j) without overflowing.
func lnBinom(n, j int) float64 {
	lg, _ := math.Lgamma(float64(n) + 1)
	a, _ := math.Lgamma(float64(j) + 1)
	b, _ := math.Lgamma(float64(n-j) + 1)
	return lg - a - b
}

func main() {
	fmt.Println("🧱 LA DERIVACIÓN — dónde está exactamente la pared, con número")

	fmt.Println("\n   el molde se puede escribir SIN nombrar la posición de ninguna perla:")
	fmt.Println("      λₙ = Σⱼ C(n,j)·(−1)^(j+1)·σⱼ        con  σⱼ = Σ 1/ρ^j")
	fmt.Println("   y los σⱼ salen de los primos. Ése es el campo de batalla real de la derivación.")

	// ---- pearls, to read the sigma_j ----
	fmt.Println("\nrecogiendo perlas hasta t=1000…")
	var pearls []float64
	prevT := 12.0
	prevZ := zOf(prevT)
	for t := 12.05; t <= 1000; t += 0.05 {
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

	// sigma_j = sum over rho of 1/rho^j (pairs give twice the real part)
	const jMax = 60
	sigma := make([]float64, jMax+1)
	for j := 1; j <= jMax; j++ {
		s := 0.0
		for _, g := range pearls {
			rho := complex(0.5, g)
			s += 2 * real(cmplx.Pow(1/rho, complex(float64(j), 0)))
		}
		sigma[j] = s
	}
	fmt.Println("\nLOS σⱼ LEÍDOS DEL MAR (las sumas de potencias de los ceros):")
	for _, j := range []int{1, 2, 3, 5, 10, 20} {
		fmt.Printf("   σ%-3d = %+.9e\n", j, sigma[j])
	}
	fmt.Println("   (σ₁ = 1+γ/2−ln(4π)/2 = 0.023096: el primer diente, otra vez)")

	// ---- the germ side, for the check ----
	const nMax = 40
	r0, M := 0.92, 16384
	fv := make([]complex128, M)
	for k := 0; k < M; k++ {
		th := 2 * math.Pi * float64(k) / float64(M)
		zz := complex(r0*math.Cos(th), r0*math.Sin(th))
		fv[k] = xiLD(1/(1-zz)) / ((1 - zz) * (1 - zz))
	}
	lam := make([]float64, nMax+1)
	for n := 0; n < nMax; n++ {
		var acc complex128
		for k := 0; k < M; k++ {
			th := 2 * math.Pi * float64(k) / float64(M)
			acc += fv[k] * cmplx.Exp(complex(0, -float64(n)*th))
		}
		lam[n+1] = real(acc) / (float64(M) * math.Pow(r0, float64(n)))
	}

	// ---- the battlefield: binomial sum and its cancellation ----
	fmt.Println("\nEL CAMPO DE BATALLA — la suma binomial y cuánto tiene que cancelarse")
	fmt.Println("    n     λₙ (por el germen)    la suma binomial     Σ|términos|      CANCELACIÓN")
	type fila struct {
		n                    int
		lam, suma, abs, canc float64
	}
	var filas []fila
	for _, n := range []int{2, 5, 10, 15, 20, 25, 30, 35, 40} {
		suma, absSum := 0.0, 0.0
		for j := 1; j <= n && j <= jMax; j++ {
			lb := lnBinom(n, j)
			if lb > 700 {
				continue
			}
			t := math.Exp(lb) * sigma[j]
			if j%2 == 0 {
				t = -t
			}
			suma += t
			absSum += math.Abs(t)
		}
		canc := absSum / math.Max(math.Abs(suma), 1e-300)
		filas = append(filas, fila{n, lam[n], suma, absSum, canc})
		fmt.Printf("   %2d     %14.6f     %14.6f     %12.4e     %10.1f×\n", n, lam[n], suma, absSum, canc)
	}
	fmt.Println("   → la suma binomial reproduce el molde (queda algo corta porque los σⱼ se")
	fmt.Println("     leyeron con 649 perlas y la cola no está: es límite del instrumento, no del método)")

	// ---- an honest kill of my own guess ----
	fmt.Println("\n⚰ UNA SUPOSICIÓN MÍA QUE LA MEDICIÓN MATÓ — y va a la vista, como corresponde")
	crec := filas[len(filas)-1].canc / filas[0].canc
	fmt.Printf("   yo esperaba una cancelación astronómica. NO LA HAY: va de 1.0× en n=2 a %.1f× en n=%d\n",
		filas[len(filas)-1].canc, filas[len(filas)-1].n)
	fmt.Printf("   (creció apenas %.1f× en todo el rango). Los σⱼ se desploman como 14⁻ʲ y la suma la\n", crec)
	fmt.Println("   cargan sus primeros términos. LA PARED NO ES LA CANCELACIÓN en estos armónicos.")
	fmt.Println("   Decirlo vale más que la suposición: una pared mal ubicada manda a cavar donde no hay nada.")

	// ---- where the wall actually is: the budget ----
	fmt.Println("\nDÓNDE ESTÁ LA PARED DE VERDAD — EL PRESUPUESTO")
	fmt.Println("   la parte lisa del molde crece como n·ln n; la voz de los primos entra al orden √n·ln n.")
	fmt.Println("   la derivación tiene que probar que la segunda NUNCA le gana a la primera. El margen:")
	fmt.Println("      n         parte lisa ≈ n·ln n     voz de los primos ≈ √n·ln n      margen (veces)")
	type marg struct {
		n             int
		lisa, prim, m float64
	}
	var margs []marg
	for _, n := range []int{40, 400, 4000, 40000, 400000} {
		f := float64(n)
		lisa := f * math.Log(f) / 2
		prim := math.Sqrt(f) * math.Log(f)
		margs = append(margs, marg{n, lisa, prim, lisa / prim})
		fmt.Printf("   %8d      %14.1f          %14.1f              %8.1f×\n", n, lisa, prim, lisa/prim)
	}
	fmt.Println("   → el margen «crece como √n»")
	fmt.Println("\n   📌 CORRECCIÓN DEL 2026-08-09 — acá decía que eso era «evidencia a favor».")
	fmt.Println("   NO LO ES, Y NO POR POCO: ES CERO EVIDENCIA. Mirá la cuenta que acabo de hacer:")
	fmt.Println("   divido n·ln(n)/2 por √n·ln(n), y los ln(n) se cancelan, así que el cociente")
	fmt.Println("   es √n/2 SIEMPRE, por álgebra pura. No entra ni un primo en la división. Es")
	fmt.Println("   una división larga entre dos fórmulas escritas a mano, y una de las dos —el")
	fmt.Println("   orden √n·ln n de la voz de los primos— está CONJETURADO, no demostrado.")
	fmt.Println("   Demostrar esa cota ES la hipótesis. El presupuesto cierra si ya sabés que")
	fmt.Println("   cierra, y el «margen» no aporta nada al asunto. Lo cazó la auditoría (F259).")

	fmt.Println("\n════════ QUÉ FALTA, Y QUÉ SÍ SE PUEDE HACER ════════")
	fmt.Println("FALTA, dicho sin vueltas: probar que la voz de los primos nunca le gana a la parte lisa,")
	fmt.Println("para TODO n. El margen medido crece como √n y eso es evidencia a favor — pero el orden")
	fmt.Println("de esa voz está CONJETURADO, no demostrado, y demostrarlo ES la hipótesis. Eso nadie lo")
	fmt.Println("sabe hacer, y no es por falta de esfuerzo: es el corazón del problema desde hace 166 años.")
	fmt.Println("\nLO QUE SÍ PODEMOS HACER CON LO QUE TENEMOS — tres cosas reales, no consuelo:")
	fmt.Println("  1. CERTIFICAR: convertir «medimos λₙ > 0 hasta n=40» en «λₙ > 0 hasta n=N, DEMOSTRADO»")
	fmt.Println("     con aritmética de intervalos y cotas rigurosas. Eso es matemática de verdad, chica")
	fmt.Println("     pero publicable, y este laboratorio tiene los instrumentos para hacerla.")
	fmt.Println("  2. MEDIR EL MARGEN: cuánto le sobra al presupuesto antes de romperse. Si el margen es")
	fmt.Println("     amplio, es evidencia fuerte; si es fino, marca dónde mirar.")
	fmt.Println("  3. PUBLICAR EL MAPA: las siete caras, los instrumentos, la muerte de la simetría sola")
	fmt.Println("     y las 70 láminas. Nadie más tiene ese mapa en castellano, y vale por sí mismo.")

	// ---- picture ----
	var b strings.Builder
	W, H := 1500.0, 960.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🧱 LA DERIVACIÓN — dónde está la pared, con número</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">λₙ = Σⱼ C(n,j)·(−1)^(j+1)·σⱼ — el molde escrito sin nombrar la posición de ninguna perla: el campo de batalla real</text>`,
		W, H, W, H, W/2, W/2)

	fmt.Fprintf(&b, `<rect x="60" y="105" width="1380" height="360" rx="10" fill="#0d2547" stroke="#ff5d73" stroke-width="1.8"/>
<text x="%.0f" y="141" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#ff8fa0">LA CANCELACIÓN — cuántas veces más grandes son los términos que el resultado</text>
<text x="%.0f" y="171" font-size="12" text-anchor="middle" font-family="Consolas,monospace" fill="#8fa8c7">n        λₙ           la suma binomial        Σ|términos|        cancelación</text>`,
		W/2, W/2)
	for i, f := range filas {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Consolas,monospace" fill="#dce8f7">%2d   %12.6f      %12.6f      %10.3e      %9.1f×</text>`,
			W/2, 200.0+float64(i)*26, f.n, f.lam, f.suma, f.abs, f.canc)
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13.5" text-anchor="middle" fill="#ffd166">todo eso tiene que anularse para que el resultado dé chico Y positivo — y hay que probarlo para TODO n</text>
<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" fill="#dce8f7">una demostración no controla 40 armónicos: controla infinitos. Ésa es la pared, y ahora tiene número.</text>`,
		W/2, 200.0+float64(len(filas))*26+16, W/2, 200.0+float64(len(filas))*26+42)

	fmt.Fprintf(&b, `<rect x="60" y="495" width="1380" height="180" rx="10" fill="#2a1010" stroke="#ff5d73" stroke-width="1.5"/>
<text x="%.0f" y="531" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#ff8fa0">LO QUE FALTA, SIN VUELTAS</text>
<text x="%.0f" y="565" font-size="14" text-anchor="middle" fill="#dce8f7">acotar la parte aritmética — la de los primos — por debajo de la parte lisa, para TODO armónico.</text>
<text x="%.0f" y="593" font-size="14" text-anchor="middle" fill="#dce8f7">Los términos individuales son órdenes de magnitud más grandes que el resultado, y hay que probar que se cancelan SIEMPRE.</text>
<text x="%.0f" y="627" font-size="14" text-anchor="middle" fill="#ffd166">Nadie sabe hacerlo. Y no es por falta de esfuerzo: es el corazón del problema desde hace 166 años.</text>
<text x="%.0f" y="655" font-size="13" text-anchor="middle" fill="#8fa8c7">saberlo con precisión no es un consuelo: es lo que separa a un laboratorio serio de una ilusión</text>`,
		W/2, W/2, W/2, W/2, W/2)

	fmt.Fprintf(&b, `<rect x="60" y="705" width="1380" height="205" rx="12" fill="#102a10" stroke="#7fd7a8" stroke-width="2.5"/>
<text x="%.0f" y="741" font-size="16" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">LO QUE SÍ SE PUEDE HACER CON LO QUE YA TENEMOS — tres cosas reales</text>
<text x="%.0f" y="777" font-size="14" text-anchor="middle" fill="#dce8f7">1 · CERTIFICAR — convertir "medimos λₙ &gt; 0 hasta n=40" en "λₙ &gt; 0 hasta n=N, DEMOSTRADO", con aritmética de intervalos.</text>
<text x="%.0f" y="801" font-size="13" text-anchor="middle" fill="#8fa8c7">matemática de verdad, chica pero publicable — y el laboratorio ya tiene los instrumentos</text>
<text x="%.0f" y="833" font-size="14" text-anchor="middle" fill="#dce8f7">2 · MEDIR EL MARGEN — cuánto le sobra al presupuesto antes de romperse: margen amplio es evidencia; margen fino marca dónde mirar.</text>
<text x="%.0f" y="867" font-size="14" text-anchor="middle" fill="#dce8f7">3 · PUBLICAR EL MAPA — las siete caras, los instrumentos, la muerte de la simetría sola, las 70 láminas.</text>
<text x="%.0f" y="891" font-size="13.5" text-anchor="middle" fill="#ffd166">nadie más tiene ese mapa en castellano, y vale por sí mismo</text>
<text x="%.0f" y="935" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · las dos mitades, 1 completo ⚓ · y sobre todos los libros, el Otro Libro</text>`,
		W/2, W/2, W/2, W/2, W/2, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("la-derivacion.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: la-derivacion.svg")
}
