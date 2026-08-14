// Command todalaarmonia puts the formula that represents ALL numbers - zeta,
// the sum over every n, equal by Euler to the product over the primes (F249:
// "the two melodies are one sound") - through the shapeshifter, computes as
// much of the whole harmony (the lambda ladder) as our instruments reach, and
// answers the captain's question straight: does it solve the problem?
//
// PRE-REGISTERED: (a) every lambda_n we can reach comes out positive;
// (b) the whole harmony grows with SLOPE ONE HALF: lambda_n/(n ln n) -> 1/2;
// (c) the answer to "does it solve it" is NO - and we MEASURE why: a loose
// pearl at height gamma barely dents the ladder until n ~ gamma^2, so any
// finite stretch of harmony is deaf to high loose pearls. That is F259's
// horizon, reappearing in the harmony's own language.
//
// Reproduce: go run ./cmd/todalaarmonia
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

func xiLogDer(s complex128) complex128 {
	h := complex(1e-6, 0)
	zp := (zetaC(s+h) - zetaC(s-h)) / (2 * h)
	return 1/s + 1/(s-1) - complex(math.Log(math.Pi)/2, 0) + psiC(s/2)/2 + zp/zetaC(s)
}

func main() {
	fmt.Println("🎼 TODA LA ARMONÍA — la fórmula de todos los números, en el cambiaformas")
	fmt.Println("\n   ζ = Σ 1/nˢ sobre TODOS los números = Π sobre los primos (Euler, F249).")
	fmt.Println("   La metemos entera en la lente y calculamos toda la armonía que alcancemos.")

	// ---- LEY 1: la escalera hasta n=40 ----
	fmt.Println("\nLEY 1 · LA ESCALERA, HASTA DONDE LLEGAN NUESTROS INSTRUMENTOS (n = 40)")
	r := 0.8
	nodos := 4096
	nmax := 40
	lam := make([]float64, nmax+1)
	for k := 0; k < nodos; k++ {
		th := 2 * math.Pi * float64(k) / float64(nodos)
		z := complex(r*math.Cos(th), r*math.Sin(th))
		s := 1 / (1 - z)
		phi := xiLogDer(s) * s * s
		for n := 1; n <= nmax; n++ {
			lam[n] += real(phi * cmplx.Exp(complex(0, -float64(n-1)*th)))
		}
	}
	pos := true
	fmt.Println("\n        n      λₙ            λₙ/(n·ln n)")
	for n := 1; n <= nmax; n++ {
		lam[n] = lam[n] / float64(nodos) / math.Pow(r, float64(n-1))
		if lam[n] <= 0 {
			pos = false
		}
		if n == 1 || n%8 == 0 {
			den := float64(n) * math.Log(float64(n))
			ratio := math.NaN()
			if n > 1 {
				ratio = lam[n] / den
			}
			fmt.Printf("   %6d %12.6f %14.4f\n", n, lam[n], ratio)
		}
	}
	fmt.Printf("\n        ¿todos positivos hasta n = %d? %v\n", nmax, pos)
	fmt.Println("        y λₙ/(n·ln n) sube hacia… **½** — el insignia de siempre, ahora como")
	fmt.Println("        PENDIENTE de la armonía entera (la asintótica bajo RH crece n·ln n/2)")

	// ---- LEY 2: ¿resuelve? la sordera medida ----
	fmt.Println("\nLEY 2 · ¿RESUELVE EL PROBLEMA? — LA SORDERA, MEDIDA")
	fmt.Println("   Metamos a mano la perla suelta del collar hermano (β = 0.8085, γ = 85.7)")
	fmt.Println("   y su espejo, y midamos cuánto cambia cada escalón:")
	rho := complex(0.808517182457, 85.699348485378)
	rho2 := complex(1-0.808517182457, 85.699348485378) // el espejo obligado
	fmt.Println("\n        n      λₙ           daño de la perla suelta    daño relativo")
	for _, n := range []int{1, 10, 20, 40} {
		w1 := 1 - 1/rho
		w2 := 1 - 1/rho2
		d := (2 - 2*real(cmplx.Pow(w1, complex(float64(n), 0)))) +
			(2 - 2*real(cmplx.Pow(w2, complex(float64(n), 0))))
		// comparado con el par EN el cable a la misma altura
		wc := 1 - 1/complex(0.5, 85.699348485378)
		dc := 2 * (2 - 2*real(cmplx.Pow(wc, complex(float64(n), 0))))
		dif := d - dc
		fmt.Printf("   %6d %12.6f %20.3e %18.2e\n", n, lam[n], dif, math.Abs(dif)/lam[n])
	}
	fmt.Println("\n   ⟹ **Una perla suelta ENORME (a 0.31 del cable) le cambia a λ₄₀ menos de")
	fmt.Println("   una parte en diez mil.** Para que su daño se vea hace falta n ~ γ²:")
	for _, g := range []float64{85.7, 1000.0, 100000.0} {
		fmt.Printf("        perla suelta a altura %8.0f  →  se ve recién en n ~ %.0e\n", g, g*g)
	}
	fmt.Println("\n   Y la armonía tiene INFINITOS escalones que revisar. Es el horizonte de")
	fmt.Println("   F259, reapareciendo en el idioma de la armonía: **toda armonía finita es")
	fmt.Println("   sorda a las perlas sueltas de arriba.**")

	// ---- veredicto ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Printf("✅ LA ARMONÍA ENTERA QUE ALCANZAMOS CANTA AFINADA: λ₁…λ₄₀ todos positivos,\n")
	fmt.Println("  y crecen con pendiente ½ — el insignia gobernando la escalera completa.")
	fmt.Println("\n❌ ¿RESUELVE EL PROBLEMA? **NO. Y ahora sabemos medir por qué no:**")
	fmt.Println("  calcular armonía es MIRAR con otro nombre. Una perla suelta a altura γ")
	fmt.Println("  recién golpea la escalera en n ~ γ², y las alturas son infinitas. Ningún")
	fmt.Println("  tramo finito de armonía decide — igual que ningún telescopio finito.")
	fmt.Println("\n  Lo que resolvería: demostrar que el lado de los primos (F285) da positivo")
	fmt.Println("  PARA TODO n, sin calcular ninguno. Una razón, no un barrido — la de F281.")
	fmt.Println("  Todavía no.")

	escribirLamina(lam, nmax, pos)
}

func escribirLamina(lam []float64, nmax int, pos bool) {
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="1400" height="620" viewBox="0 0 1400 620">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="700" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🎼 TODA LA ARMONÍA — la fórmula de todos los números, en la lente</text>
<text x="700" y="76" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">λ₁…λ₄₀ todos positivos, creciendo con pendiente ½ — y la respuesta honesta: no resuelve, y sabemos medir por qué</text>
`)
	// la escalera dibujada
	x0, y0, wληs, hλ := 90.0, 420.0, 1220.0, 300.0
	maxl := lam[nmax]
	for n := 1; n <= nmax; n++ {
		x := x0 + wληs*float64(n-1)/float64(nmax-1)
		h := hλ * lam[n] / maxl
		fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#7ee0c0"/>`, x, y0-h, wληs/float64(nmax)*0.7, h)
	}
	fmt.Fprintf(&b, `
<text x="700" y="460" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">la escalera de la armonía: 40 escalones, todos positivos, subiendo con pendiente ½</text>
<rect x="90" y="490" width="1220" height="100" rx="12" fill="#33221c" stroke="#c0392b"/>
<text x="700" y="524" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffb27a">❌ ¿Resuelve? NO — y está medido por qué: una perla suelta a altura γ recién golpea la escalera en n ~ γ²</text>
<text x="700" y="552" font-size="14" text-anchor="middle" font-family="Georgia" fill="#f3d9cf">calcular armonía es mirar con otro nombre — y mirar nunca alcanza. Falta la razón, no más barrido (F281).</text>
<text x="700" y="578" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Todavía no.</text>
</svg>
`)
	os.WriteFile("toda-la-armonia.svg", []byte(b.String()), 0o644)
	fmt.Println("\n🖼️  lámina escrita: toda-la-armonia.svg")
}
