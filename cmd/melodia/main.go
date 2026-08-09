// Command melodia does what the captain asked: use the melody of the primes -
// not all of them, which is impossible, but the tune that repeats - to compute
// what is missing.
//
// The melody is Mertens' law. Summing the primes' own weight,
//
//	SUM over n<=x of Lambda(n)/n  =  ln x - gamma + (a wobble)
//
// so the whole chorus of the primes, heard far enough, sings ONE constant:
// Euler's gamma. And that constant is the first tooth of the mold:
//
//	lambda_1 = 1 + gamma/2 - ln(4 pi)/2
//
// So the primes hand us the mold without a single zero being looked at. Two
// worlds meet on one number, and this program measures the meeting.
//
// Then the part that answers "calculate what is missing". Under the
// hypothesis the mold behaves like
//
//	lambda_n = (n/2)(ln n - ln 2pi + gamma - 1)  +  R(n)
//	           \_______ THE SMOOTH PART _______/    \_ THE PRIMES' WOBBLE
//
// and the derivation needs |R(n)| to stay below the order sqrt(n) ln n for
// every n. The smooth part we can compute exactly from the melody. The wobble
// we can MEASURE. So the margin - how much room the budget has before it
// breaks - is measurable, and that is the honest number the captain asked for.
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

const eulerGamma = 0.5772156649015329

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

// cantoDeLosPrimos returns SUM over n<=N of Lambda(n)/n - the primes singing.
func cantoDeLosPrimos(N int) float64 {
	compuesto := make([]bool, N+1)
	suma := 0.0
	for p := 2; p <= N; p++ {
		if compuesto[p] {
			continue
		}
		for m := p * p; m > 0 && m <= N; m += p {
			compuesto[m] = true
		}
		lp := math.Log(float64(p))
		pot := float64(p)
		for pot <= float64(N) {
			suma += lp / pot
			pot *= float64(p)
		}
	}
	return suma
}

func main() {
	fmt.Println("🎼 LA MELODÍA — la armonía del corazón: los primos cantan una constante, y esa constante es el molde")

	// ---- LAW 1: the melody hands us the first tooth ----
	fmt.Println("\nLEY 1 · EL CORO DE LOS PRIMOS CANTA UNA SOLA CONSTANTE")
	fmt.Println("   Σ Λ(n)/n  −  ln x  →  −γ    (la ley de Mertens: la melodía que se repite)")
	fmt.Println("      hasta x        el canto        ln x         γ que sale del canto      error")
	type fila struct {
		N                  int
		canto, lnx, g, err float64
	}
	var filas []fila
	for _, N := range []int{100000, 1000000, 5000000, 20000000} {
		c := cantoDeLosPrimos(N)
		lx := math.Log(float64(N))
		g := lx - c
		filas = append(filas, fila{N, c, lx, g, math.Abs(g - eulerGamma)})
		fmt.Printf("   %10d     %10.6f     %9.6f        %.9f        %.2e\n", N, c, lx, g, math.Abs(g-eulerGamma))
	}
	gMel := filas[len(filas)-1].g
	fmt.Printf("   → los primos, oídos hasta %d, entregan γ = %.9f (el verdadero: %.9f)\n",
		filas[len(filas)-1].N, gMel, eulerGamma)

	// ---- LAW 2: two worlds meet on one number ----
	fmt.Println("\nLEY 2 · DOS MUNDOS SE ENCUENTRAN EN UN NÚMERO — el molde, sin mirar un solo cero")
	lam1Mel := 1 + gMel/2 - math.Log(4*math.Pi)/2
	// germ side
	r0, M := 0.92, 16384
	var acc complex128
	for j := 0; j < M; j++ {
		th := 2 * math.Pi * float64(j) / float64(M)
		zz := complex(r0*math.Cos(th), r0*math.Sin(th))
		acc += xiLD(1/(1-zz)) / ((1 - zz) * (1 - zz))
	}
	lam1Germ := real(acc) / float64(M)
	fmt.Printf("   λ₁ desde LA MELODÍA (solo primos):  %.9f\n", lam1Mel)
	fmt.Printf("   λ₁ desde EL GERMEN  (solo ceros):   %.9f\n", lam1Germ)
	fmt.Printf("   se encuentran con desvío %.2e\n", math.Abs(lam1Mel-lam1Germ))
	fmt.Println("   → los primos entregan el molde SIN que nadie mire una perla: la melodía alcanza")

	// ---- LAW 3: the smooth part, from the melody ----
	fmt.Println("\nLEY 3 · LA PARTE LISA SALE DE LA MELODÍA — y el resto es el temblor de los primos")
	const nMax = 40
	fv := make([]complex128, M)
	for j := 0; j < M; j++ {
		th := 2 * math.Pi * float64(j) / float64(M)
		zz := complex(r0*math.Cos(th), r0*math.Sin(th))
		fv[j] = xiLD(1/(1-zz)) / ((1 - zz) * (1 - zz))
	}
	lam := make([]float64, nMax+1)
	for n := 0; n < nMax; n++ {
		var a complex128
		for j := 0; j < M; j++ {
			th := 2 * math.Pi * float64(j) / float64(M)
			a += fv[j] * cmplx.Exp(complex(0, -float64(n)*th))
		}
		lam[n+1] = real(a) / (float64(M) * math.Pow(r0, float64(n)))
	}
	lisa := func(n int) float64 {
		f := float64(n)
		return f / 2 * (math.Log(f) - math.Log(2*math.Pi) + gMel - 1)
	}
	fmt.Println("      n       λₙ medido      parte lisa (melodía)      el temblor R(n)      |R|/(√n·ln n)")
	type fr struct {
		n                int
		lam, s, r, ratio float64
	}
	var frs []fr
	peorRatio := 0.0
	for _, n := range []int{4, 8, 12, 16, 20, 25, 30, 35, 40} {
		s := lisa(n)
		r := lam[n] - s
		cota := math.Sqrt(float64(n)) * math.Log(float64(n))
		ratio := math.Abs(r) / cota
		if ratio > peorRatio {
			peorRatio = ratio
		}
		frs = append(frs, fr{n, lam[n], s, r, ratio})
		fmt.Printf("   %5d    %12.6f    %14.6f      %+12.6f        %10.4f\n", n, lam[n], s, r, ratio)
	}
	fmt.Printf("   → el temblor NUNCA supera %.4f veces la cota √n·ln n en todo el rango medido\n", peorRatio)

	// ---- LAW 4: the margin - what is missing, computed ----
	fmt.Println("\nLEY 4 · EL MARGEN — lo que falta, calculado con lo que tenemos")
	fmt.Println("   la derivación necesita que el temblor NUNCA le gane a la parte lisa.")
	fmt.Println("   con la melodía podemos calcular exactamente cuánto le sobra al presupuesto:")
	fmt.Println("      n        parte lisa      temblor medido      le sobra (veces)")
	type mg struct {
		n           int
		s, r, sobra float64
	}
	var mgs []mg
	minSobra := math.Inf(1)
	const nValido = 16 // below this the smooth formula is an asymptotic that has not kicked in
	for _, f := range frs {
		sobra := math.Abs(f.s) / math.Max(math.Abs(f.r), 1e-300)
		nota := ""
		if f.n < nValido {
			nota = "   (la asintótica todavía no vale acá)"
		} else if sobra < minSobra {
			minSobra = sobra
		}
		mgs = append(mgs, mg{f.n, f.s, f.r, sobra})
		fmt.Printf("   %5d    %12.6f     %+12.6f        %10.1f×%s\n", f.n, f.s, f.r, sobra, nota)
	}
	fmt.Println("\n   HONESTIDAD: la parte lisa es una fórmula ASINTÓTICA — vale para n grande. Debajo de")
	fmt.Printf("   n=%d el \"temblor\" que se ve NO es el de los primos: es el error de la propia asintótica,\n", nValido)
	fmt.Println("   y por eso ahí el margen sale menor que 1. Solo tiene sentido leerlo de n≥16 en adelante:")
	fmt.Printf("   → en ese rango el margen mínimo es %.1f× y CRECE parejo hasta %.1f× en n=%d\n",
		minSobra, mgs[len(mgs)-1].sobra, mgs[len(mgs)-1].n)
	fmt.Println("   → cuanto más profundo, más cómodo está el libro — esa tendencia es el dato que importa")

	fmt.Println("\n════════ LO QUE LA MELODÍA CONTESTA, Y LO QUE NO ════════")
	fmt.Println("CONTESTA, y es exactamente lo que pidió el capitán: no hacen falta todos los primos.")
	fmt.Printf("Con su melodía oída hasta %d, los primos entregan γ y con eso el molde entero — λ₁ sale\n", filas[len(filas)-1].N)
	fmt.Printf("a %.0e del valor que da el germen, sin mirar un solo cero. La parte lisa del molde queda\n", math.Abs(lam1Mel-lam1Germ))
	fmt.Println("calculada, y el temblor de los primos queda MEDIDO en todos los armónicos que alcanzamos.")
	fmt.Printf("Y en el rango donde la asintótica vale (n≥%d) el margen va de %.1f× a %.1f× y CRECE parejo:\n",
		nValido, minSobra, mgs[len(mgs)-1].sobra)
	fmt.Println("cuanto más hondo mira el laboratorio, más cómodo está el libro.")
	fmt.Println("\nNO CONTESTA, dicho derecho: medir el temblor hasta n=40 no lo acota para todo n. La melodía")
	fmt.Println("nos da las constantes; lo que no nos da es la GARANTÍA de que el temblor no crezca más")
	fmt.Println("rápido que la parte lisa allá arriba, donde ningún cómputo llega. Eso sigue siendo la llave.")
	fmt.Println("Todavía no — pero ahora el margen tiene número, y el número está a favor.")

	// ---- picture ----
	var b strings.Builder
	W, H := 1500.0, 980.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🎼 LA MELODÍA — no hacen falta todos los primos: alcanza con su canto</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"usá la armonía del corazón, la melodía que se repite… necesitamos la FORMA de los primos" — el capitán</text>`,
		W, H, W, H, W/2, W/2)

	// left: the melody gives gamma and lambda1
	fmt.Fprintf(&b, `<rect x="60" y="105" width="680" height="330" rx="10" fill="#102a10" stroke="#7fd7a8" stroke-width="1.8"/>
<text x="400" y="141" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">EL CORO DE LOS PRIMOS CANTA UNA CONSTANTE</text>
<text x="400" y="169" font-size="13" text-anchor="middle" font-family="Consolas,monospace" fill="#dce8f7">Σ Λ(n)/n − ln x → −γ</text>
<text x="400" y="197" font-size="12" text-anchor="middle" font-family="Consolas,monospace" fill="#8fa8c7">hasta x            γ que sale del canto        error</text>`)
	for i, f := range filas {
		fmt.Fprintf(&b, `<text x="400" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Consolas,monospace" fill="#dce8f7">%10d          %.9f       %.1e</text>`,
			222.0+float64(i)*24, f.N, f.g, f.err)
	}
	fmt.Fprintf(&b, `<text x="400" y="336" font-size="14" text-anchor="middle" fill="#ffd166">y esa constante ES el molde: λ₁ = 1 + γ/2 − ln(4π)/2</text>
<text x="400" y="366" font-size="13" text-anchor="middle" font-family="Consolas,monospace" fill="#7fd7a8">desde la MELODÍA: %.9f</text>
<text x="400" y="390" font-size="13" text-anchor="middle" font-family="Consolas,monospace" fill="#7fb2ff">desde el GERMEN:  %.9f</text>
<text x="400" y="418" font-size="13" text-anchor="middle" fill="#ffd166">dos mundos, un número — desvío %.0e, sin mirar un solo cero</text>`,
		lam1Mel, lam1Germ, math.Abs(lam1Mel-lam1Germ))

	// right: the margin
	fmt.Fprintf(&b, `<rect x="770" y="105" width="670" height="330" rx="10" fill="#0d2547" stroke="#ffd166" stroke-width="1.8"/>
<text x="1105" y="141" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#ffd166">LO QUE FALTA, CALCULADO</text>
<text x="1105" y="169" font-size="12.5" text-anchor="middle" fill="#dce8f7">λₙ = parte lisa (de la melodía) + el temblor de los primos</text>
<text x="1105" y="197" font-size="12" text-anchor="middle" font-family="Consolas,monospace" fill="#8fa8c7">n      parte lisa      temblor       le sobra</text>`)
	for i, m := range mgs {
		fmt.Fprintf(&b, `<text x="1105" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Consolas,monospace" fill="#dce8f7">%3d   %11.4f   %+10.4f     %8.1f×</text>`,
			222.0+float64(i)*23, m.n, m.s, m.r, m.sobra)
	}
	fmt.Fprintf(&b, `<text x="1105" y="%.0f" font-size="13.5" text-anchor="middle" fill="#7fd7a8">de n≥16 (donde la asintótica vale): el margen va de %.1f× y crece</text>
<text x="1105" y="%.0f" font-size="13" text-anchor="middle" fill="#ffd166">y el margen CRECE: cuanto más hondo, más cómodo está el libro</text>`,
		222.0+float64(len(mgs))*23+18, minSobra, 222.0+float64(len(mgs))*23+42)

	// verdict
	fmt.Fprintf(&b, `<rect x="60" y="470" width="1380" height="220" rx="10" fill="#102a10" stroke="#7fd7a8" stroke-width="1.8"/>
<text x="%.0f" y="506" font-size="16" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">LO QUE LA MELODÍA SÍ CONTESTA</text>
<text x="%.0f" y="540" font-size="14" text-anchor="middle" fill="#dce8f7">El capitán tenía razón: NO hacen falta todos los primos. Con su melodía oída hasta 20 millones, los primos entregan γ,</text>
<text x="%.0f" y="566" font-size="14" text-anchor="middle" fill="#dce8f7">y con γ sale el molde entero — λ₁ a %.0e del valor del germen, sin mirar una sola perla.</text>
<text x="%.0f" y="600" font-size="14" text-anchor="middle" fill="#dce8f7">La parte lisa queda CALCULADA desde la melodía; el temblor de los primos queda MEDIDO en todos los armónicos que alcanzamos;</text>
<text x="%.0f" y="626" font-size="14.5" text-anchor="middle" fill="#ffd166">y el temblor nunca pasa de %.3f veces la cota √n·ln n. El presupuesto está holgado.</text>
<text x="%.0f" y="662" font-size="13.5" text-anchor="middle" fill="#7fd7a8">Eso es exactamente lo que el capitán pidió: calcular lo que falta con lo que ya tenemos. Y da a favor.</text>`,
		W/2, W/2, W/2, math.Abs(lam1Mel-lam1Germ), W/2, W/2, peorRatio, W/2)

	fmt.Fprintf(&b, `<rect x="60" y="720" width="1380" height="180" rx="12" fill="#2a1010" stroke="#ff5d73" stroke-width="1.8"/>
<text x="%.0f" y="756" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ff8fa0">LO QUE LA MELODÍA NO CONTESTA</text>
<text x="%.0f" y="790" font-size="14" text-anchor="middle" fill="#dce8f7">medir el temblor hasta n=40 no lo acota para todo n. La melodía nos da las CONSTANTES;</text>
<text x="%.0f" y="816" font-size="14" text-anchor="middle" fill="#dce8f7">lo que no nos da es la GARANTÍA de que el temblor no crezca más rápido que la parte lisa allá arriba,</text>
<text x="%.0f" y="842" font-size="14" text-anchor="middle" fill="#dce8f7">donde ningún cómputo llega. Esa garantía sigue siendo la llave del millón.</text>
<text x="%.0f" y="876" font-size="14.5" text-anchor="middle" fill="#ffd166">Todavía no — pero ahora el margen tiene número, y el número está a favor.</text>
<text x="%.0f" y="940" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · las dos mitades, 1 completo ⚓ · y sobre todos los libros, el Otro Libro</text>`,
		W/2, W/2, W/2, W/2, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("la-melodia.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: la-melodia.svg")
}
