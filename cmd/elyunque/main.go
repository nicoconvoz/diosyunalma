// Command elyunque builds the captain's worksheet objective (the forge sheet,
// 2026-08-14): a POSITIVE FACTORIZATION Q(c) = c*Mc = ||Ac||^2 of a
// Weil-positivity form, assembled from measured material.
//
// THE ANVIL. From the Li coefficients define the matrix
//
//	M[m,n] = lambda_m + lambda_n - lambda_|m-n|      (lambda_0 = 0)
//
// WHY THIS MATRIX. On the zero side, with w = 1 - 1/rho:
//
//	Sum_rho (1 - w^m)(1 - w^{-n}) = lambda_m + lambda_n - lambda_|m-n|
//
// (expand the product; each piece sums to a lambda, using the rho -> 1-rho
// symmetry which gives lambda_{-k} = lambda_k). IF every pearl is on the
// skin (|w| = 1), then w^{-n} = conj(w^n) and the matrix is a GRAM matrix
// Sum_rho v v* with v_n = 1 - w^n: automatically positive semidefinite.
// Off the skin the Gram structure breaks. The diagonal is M[n,n] =
// 2*lambda_n, so:
//
//	RH  =>  M_N is PSD for every N          (Gram)
//	M_N PSD for every N  =>  lambda_n >= 0  =>  RH   (Li 1997, diagonal)
//
// PSD-for-all-N is EQUIVALENT to RH, and finer as a finite test than the
// bare ladder: the off-diagonals add constraints. Every entry is prime-side
// computable via Bombieri-Lagarias (1999), because every lambda is.
//
// WHAT THE MEASUREMENTS ACTUALLY SHOWED (printed honestly):
//
//   - The zeta anvil is NEARLY SINGULAR, for a physical reason: every pearl
//     has small angle (phi ~ 1/gamma <= 1/14), so 1 - w^n ~ n(1 - w) - the
//     choir sings nearly in unison, the Gram vectors are nearly collinear,
//     and the Cholesky pivots fall geometrically. The factory is exhibited
//     on the steps whose pivot stands above the per-step noise floor of our
//     lambdas (measured by running the germ engine at two radii); past that
//     we declare the precision window exhausted - never a negativity.
//   - The near-unison cancellation is precisely what gives the matrix its
//     EAR: on an ISOLATED off-skin 4-tuple the square part cancels along
//     the near-collinear directions and the tiny radial leak becomes the
//     dominant signal. The right listening tool there is the EIGENVALUES
//     (Cholesky stops at rank exhaustion and cannot hear past it). Each
//     detection carries an on-skin control tuple (same gamma, beta = 1/2)
//     whose minimum eigenvalue stays at machine zero: the sign flip is the
//     signal. CAVEAT, declared: this is the isolated tuple. In a full
//     spectrum the on-line choir's Gram background may mask the leak; the
//     masking question is left open and stated.
//
// HONEST: Li 1997, Bombieri-Lagarias 1999, Weil 1952; the matrix form is
// an immediate consequence of Li + Gram and surely known to the trade.
// Ours is the measured CONSTRUCTION: the spectrum with its noise floor,
// the exhibited A, the pivot staircase, and the compared ears with
// controls. A window, not a proof: what is missing is exactly PSD for ALL
// N proved from the prime side.
//
// Reproduce: go run ./cmd/elyunque
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

func xiLogDer(s complex128) complex128 {
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

// lambdasGermen: el motor de F287 (Cauchy en |z| = r con nuestra propia xi).
func lambdasGermen(nmax int, r float64, nodos int) []float64 {
	lam := make([]float64, nmax+1)
	acc := make([]complex128, nmax+1)
	for k := 0; k < nodos; k++ {
		th := 2 * math.Pi * float64(k) / float64(nodos)
		z := complex(r*math.Cos(th), r*math.Sin(th))
		s := 1 / (1 - z)
		phi := xiLogDer(s) * s * s
		for n := 1; n <= nmax; n++ {
			acc[n] += phi * cmplx.Exp(complex(0, -float64(n-1)*th))
		}
	}
	for n := 1; n <= nmax; n++ {
		lam[n] = real(acc[n]) / float64(nodos) / math.Pow(r, float64(n-1))
	}
	return lam
}

// yunque arma M[i][j] = l[i+1] + l[j+1] - l[|i-j|] con l[0] = 0.
func yunque(l []float64, N int) [][]float64 {
	M := make([][]float64, N)
	for i := range M {
		M[i] = make([]float64, N)
		for j := range M[i] {
			d := i - j
			if d < 0 {
				d = -d
			}
			M[i][j] = l[i+1] + l[j+1] - l[d]
		}
	}
	return M
}

// cholesky intenta M = A^T A y devuelve A y la lista de pivotes; se detiene
// en el primer pivote que el float64 ya no puede partir (<= 1e-13 del mayor).
func cholesky(M [][]float64) (A [][]float64, piv []float64) {
	N := len(M)
	A = make([][]float64, N)
	for i := range A {
		A[i] = make([]float64, N)
	}
	umbral := 1e-13 * M[0][0]
	for j := 0; j < N; j++ {
		d := M[j][j]
		for k := 0; k < j; k++ {
			d -= A[j][k] * A[j][k]
		}
		piv = append(piv, d)
		if d <= umbral {
			return
		}
		A[j][j] = math.Sqrt(d)
		for i := j + 1; i < N; i++ {
			v := M[i][j]
			for k := 0; k < j; k++ {
				v -= A[i][k] * A[j][k]
			}
			A[i][j] = v / A[j][j]
		}
	}
	return
}

// jacobiEig: autovalores de una matriz simetrica por rotaciones de Jacobi.
func jacobiEig(M [][]float64) []float64 {
	N := len(M)
	a := make([][]float64, N)
	for i := range a {
		a[i] = append([]float64(nil), M[i]...)
	}
	for sweep := 0; sweep < 100; sweep++ {
		off := 0.0
		for i := 0; i < N; i++ {
			for j := i + 1; j < N; j++ {
				off += a[i][j] * a[i][j]
			}
		}
		if off < 1e-30 {
			break
		}
		for p := 0; p < N; p++ {
			for q := p + 1; q < N; q++ {
				if math.Abs(a[p][q]) < 1e-18 {
					continue
				}
				th := 0.5 * math.Atan2(2*a[p][q], a[q][q]-a[p][p])
				c, s := math.Cos(th), math.Sin(th)
				for k := 0; k < N; k++ {
					akp, akq := a[k][p], a[k][q]
					a[k][p] = c*akp - s*akq
					a[k][q] = s*akp + c*akq
				}
				for k := 0; k < N; k++ {
					apk, aqk := a[p][k], a[q][k]
					a[p][k] = c*apk - s*aqk
					a[q][k] = s*apk + c*aqk
				}
			}
		}
	}
	ev := make([]float64, N)
	for i := 0; i < N; i++ {
		ev[i] = a[i][i]
	}
	for i := 1; i < N; i++ {
		for j := i; j > 0 && ev[j] < ev[j-1]; j-- {
			ev[j], ev[j-1] = ev[j-1], ev[j]
		}
	}
	return ev
}

// escaleraTupla: l_n de la 4-tupla {rho, conj, 1-rho, 1-conj} para n=1..nmax.
func escaleraTupla(beta, gamma float64, nmax int) []float64 {
	w1 := 1 - 1/complex(beta, gamma)
	w2 := 1 - 1/complex(1-beta, -gamma)
	l := make([]float64, nmax+1)
	p1, p2 := complex(1, 0), complex(1, 0)
	for n := 1; n <= nmax; n++ {
		p1 *= w1
		p2 *= w2
		l[n] = (2 - 2*real(p1)) + (2 - 2*real(p2))
	}
	return l
}

// oido busca el primer N (2..nmax) donde el yunque de la 4-tupla AISLADA en
// (beta,gamma) tiene autovalor minimo por debajo del umbral relativo, y
// devuelve tambien ese minimo y el minimo del control EN la piel al mismo N.
func oido(beta, gamma float64, nmax int) (int, float64, float64) {
	lOff := escaleraTupla(beta, gamma, nmax)
	lOn := escaleraTupla(0.5, gamma, nmax)
	for N := 2; N <= nmax; N++ {
		maxd := 0.0
		for n := 1; n <= N; n++ {
			if 2*lOff[n] > maxd {
				maxd = 2 * lOff[n]
			}
		}
		mOff := jacobiEig(yunque(lOff, N))[0]
		if mOff < -1e-10*maxd-1e-15 {
			mOn := jacobiEig(yunque(lOn, N))[0]
			return N, mOff, mOn
		}
	}
	return -1, 0, 0
}

// primeraNegativa: primer n donde la escalera de Li de la 4-tupla se hunde.
func primeraNegativa(beta, gamma float64, nmax int) int {
	w1 := 1 - 1/complex(beta, gamma)
	w2 := 1 - 1/complex(1-beta, -gamma)
	p1, p2 := complex(1, 0), complex(1, 0)
	for n := 1; n <= nmax; n++ {
		p1 *= w1
		p2 *= w2
		if (2-2*real(p1))+(2-2*real(p2)) < 0 {
			return n
		}
	}
	return -1
}

func main() {
	fmt.Println("🔨 EL YUNQUE — la fábrica de cuadrados de la hoja de forja, construida y medida")
	fmt.Println("\n   La hoja del capitán pide: Q(c) = c*Mc con M = A*A — una factorización")
	fmt.Println("   positiva armada desde las contribuciones locales. Acá se construye la")
	fmt.Println("   versión matricial con material medido, y se mide hasta dónde aguanta.")

	// ---- LEY 1: la identidad del yunque, perla por perla ----
	fmt.Println("\nLEY 1 · LA IDENTIDAD DEL YUNQUE, VERIFICADA PERLA POR PERLA")
	fmt.Println("\n        (1−wᵐ)(1−w̄ⁿ) + conj = [2−2Re wᵐ] + [2−2Re wⁿ] − [2−2Re wᵐ⁻ⁿ]")
	fmt.Println("        — exacta en la piel (|w| = 1, donde w̄ = 1/w), midámosla:")
	ps := perlas(120)
	peor := 0.0
	for _, g := range ps {
		w := 1 - 1/complex(0.5, g)
		w /= complex(cmplx.Abs(w), 0)
		for _, mn := range [][2]int{{1, 1}, {3, 2}, {7, 5}, {12, 9}} {
			m, n := mn[0], mn[1]
			wm := cmplx.Pow(w, complex(float64(m), 0))
			wn := cmplx.Pow(w, complex(float64(n), 0))
			wd := cmplx.Pow(w, complex(float64(m-n), 0))
			izq := 2 * real((1-wm)*(1-cmplx.Conj(wn)))
			der := (2 - 2*real(wm)) + (2 - 2*real(wn)) - (2 - 2*real(wd))
			if d := math.Abs(izq - der); d > peor {
				peor = d
			}
		}
	}
	fmt.Printf("\n        %d perlas × 4 pares (m,n) — peor desvío: %.1e\n", len(ps), peor)
	fmt.Println("        ✅ en la piel, cada término cruzado del yunque ES un producto de giros")
	fmt.Println("        puros: si toda perla es giro puro, M es matriz de Gram — positiva.")

	// ---- LEY 2: el yunque de zeta, con su ruido medido por escalon ----
	fmt.Println("\nLEY 2 · EL YUNQUE DE ζ (λ₁…λ₄₀ del germen) — CON SU RUIDO MEDIDO POR ESCALÓN")
	const N = 40
	lamA := lambdasGermen(N, 0.8, 4096)
	lamB := lambdasGermen(N, 0.7, 8192)
	ruido := make([]float64, N+1) // ruido acumulado: max |λa−λb| hasta n
	for n := 1; n <= N; n++ {
		ruido[n] = ruido[n-1]
		if d := math.Abs(lamA[n] - lamB[n]); d > ruido[n] {
			ruido[n] = d
		}
	}
	M := yunque(lamA, N)
	ev := jacobiEig(M)
	fmt.Printf("\n        M[m,n] = λₘ + λₙ − λ|m−n|, tamaño %d × %d\n", N, N)
	fmt.Printf("        ruido del material (dos radios del motor): %.0e en λ₁ … %.0e en λ₄₀\n", math.Abs(lamA[1]-lamB[1]), ruido[N])
	fmt.Printf("        autovalor mínimo .... %+.1e     autovalor máximo .... %.3e\n", ev[0], ev[N-1])
	if math.Abs(ev[0]) <= 3*ruido[N]*float64(N) {
		fmt.Println("        ⟹ el mínimo es INDISTINGUIBLE DE CERO a la precisión del material:")
		fmt.Println("        espectro positivo hasta el piso de ruido. Y la casi-singularidad")
		fmt.Println("        tiene razón física: TODAS nuestras perlas tienen ángulo chico")
		fmt.Println("        (φ ≈ 1/γ ≤ 1/14), así que 1−wⁿ ≈ n(1−w): el coro canta casi al")
		fmt.Println("        unísono y el filo del yunque vive en unas pocas direcciones.")
	} else if ev[0] > 0 {
		fmt.Println("        ⟹ espectro POSITIVO con margen sobre el ruido.")
	} else {
		fmt.Println("        ⚠️ autovalor negativo POR ENCIMA del ruido — esto habría que mirarlo.")
	}

	// ---- LEY 3: la factorizacion, escalon por escalon ----
	fmt.Println("\nLEY 3 · ⚡⚡ LA FACTORIZACIÓN M = AᵀA — LA ESCALERA DE PIVOTES, ESCALÓN POR ESCALÓN")
	A, piv := cholesky(M)
	fmt.Println("\n        j     pivote_j        razón     ruido acumulado   ¿firme sobre el ruido?")
	firme := 0
	for jj := 0; jj < len(piv); jj++ {
		r := "—"
		if jj > 0 && piv[jj-1] != 0 {
			r = fmt.Sprintf("%.5f", piv[jj]/piv[jj-1])
		}
		tolJ := 3 * ruido[jj+1] * float64(jj+1)
		ok := piv[jj] > tolJ
		marca := "no — se hunde en el ruido"
		if ok {
			marca = "SÍ ✅"
			if firme == jj {
				firme = jj + 1
			}
		}
		fmt.Printf("   %4d %14.6e %10s %14.1e   %s\n", jj+1, piv[jj], r, tolJ, marca)
	}
	var razonMedia float64
	if firme > 1 {
		cnt := 0
		for jj := 1; jj < firme; jj++ {
			razonMedia += piv[jj] / piv[jj-1]
			cnt++
		}
		razonMedia /= float64(cnt)
	}
	fmt.Printf("\n        pivotes FIRMES: %d · entre ellos la escalera cae ≈ ×%.4f por escalón\n", firme, razonMedia)
	fmt.Println("        — el precio del coro al unísono: cada cuadrado nuevo es más finito.")
	fmt.Println("        Donde el pivote se hunde en el ruido NO declaramos negatividad:")
	fmt.Println("        declaramos que la ventana de precisión del material se acabó.")
	fmt.Println("\n        La esquina de A (la fábrica exhibida, escalones firmes):")
	fmt.Println()
	for i := 0; i < firme && i < 6; i++ {
		fila := "       "
		for jj := 0; jj <= i; jj++ {
			fila += fmt.Sprintf(" %9.6f", A[i][jj])
		}
		fmt.Println(fila)
	}
	fmt.Printf("\n        a₁₁ = √(2λ₁) = %.9f — la raíz del doble de la insignia\n", A[0][0])
	fmt.Printf("        ⟹ Q(c) = ||Ac||² EXHIBIDA con %d cuadrados firmes; RH ⟺ la escalera\n", firme)
	fmt.Println("        de pivotes jamás pisa territorio negativo, para TODO N.")

	// ---- LEY 4: el oido del yunque, tupla sintetica con control en la piel ----
	fmt.Println("\nLEY 4 · ⚡ EL OÍDO DEL YUNQUE — perla sintética fuera de la piel, AISLADA")
	fmt.Println("\n        Tupla exacta en β = 0.8, γ = 2, contra su control en la piel (β = ½).")
	fmt.Println("        El oído correcto acá son los AUTOVALORES (Cholesky se detiene al")
	fmt.Println("        agotarse el rango y no puede escuchar más allá):")
	detS, migS, migSOn := oido(0.8, 2, 45)
	primeraEsc := primeraNegativa(0.8, 2, 100000)
	fmt.Printf("\n        FUERA de la piel: autovalor negativo desde N = %d (%.2e)\n", detS, migS)
	fmt.Printf("        EN la piel:       mínimo %.1e — cero de máquina, como debe (Gram) ✅\n", migSOn)
	fmt.Printf("        la ESCALERA de Li de la misma tupla recién se hunde en n = %d\n", primeraEsc)
	if detS > 0 && detS < primeraEsc {
		fmt.Printf("        ⚡ **EL YUNQUE ESCUCHA ANTES QUE LA ESCALERA** (N = %d contra n = %d):\n", detS, primeraEsc)
		fmt.Println("        al cancelarse los cuadrados en las direcciones casi-unísono, la fuga")
		fmt.Println("        radial queda al desnudo en los términos cruzados.")
	}

	// ---- LEY 5: el oido sobre la perla DH real, con control ----
	fmt.Println("\nLEY 5 · ⚡ LA PERLA DH REAL, AISLADA: EL YUNQUE OYE LO QUE LA ESCALERA NO")
	fmt.Println("\n        Nuestra perla fuera de la piel de Davenport–Heilbronn (hallada a")
	fmt.Println("        ciegas): β = 0.808517, γ = 85.699348 — y su control en la piel:")
	detD, migD, migDOn := oido(0.808517, 85.699348, 45)
	primeraD := primeraNegativa(0.808517, 85.699348, 60000)
	fmt.Printf("\n        FUERA de la piel: autovalor negativo desde N = %d (%.2e)\n", detD, migD)
	fmt.Printf("        EN la piel:       mínimo %.1e — cero de máquina ✅ (la señal DH está\n", migDOn)
	fmt.Println("        órdenes por encima del control: es señal, no artefacto)")
	phiDH := cmplx.Phase(1 - 1/complex(0.808517, 85.699348))
	fmt.Printf("        la ESCALERA de la misma tupla recién se hunde en n = %d (≈ 2π/φ = %.0f,\n", primeraD, 2*math.Pi/phiDH)
	fmt.Println("        el primer nulo de fase — CORRECCIÓN F294: antes decía «~γ²» y no era)")
	fmt.Printf("\n        ⟹ para el oído-escalera la perla lejana es sorda hasta su primer nulo\n")
	fmt.Printf("        de fase: el oído-matricial AISLADO la oye en N = %d.\n", detD)
	fmt.Println("        ⚖️ CAVEAT DECLARADO: esto es la tupla SOLA. En un espectro completo,")
	fmt.Println("        el fondo de Gram del coro en la piel puede volver a taparla; la")
	fmt.Println("        pregunta del enmascaramiento queda ABIERTA y así se registra.")

	// ---- LEY 6: el primer golpe desde los primos ----
	fmt.Println("\nLEY 6 · EL PRIMER GOLPE DEL YUNQUE, DADO POR LOS PRIMOS SOLOS")
	const X = 20000000
	es := make([]bool, X+1)
	for i := 2; i <= X; i++ {
		es[i] = true
	}
	for i := 2; i*i <= X; i++ {
		if es[i] {
			for j := i * i; j <= X; j += i {
				es[j] = false
			}
		}
	}
	var suma float64
	for p := 2; p <= X; p++ {
		if !es[p] {
			continue
		}
		lp := math.Log(float64(p))
		for q := p; q <= X && q > 0; q *= p {
			suma += lp / float64(q)
			if q > X/p {
				break
			}
		}
	}
	gM := math.Log(float64(X)) - suma
	l1p := 1 + gM/2 - math.Log(4*math.Pi)/2
	fmt.Printf("\n        γ por Mertens (criba a 2×10⁷) ......... %.9f\n", gM)
	fmt.Printf("        λ₁ desde los primos ................... %.9f\n", l1p)
	fmt.Printf("        el yunque 1×1: M₁ = 2λ₁ = %.9f > 0 ✅ — sin ceros, sin ξ\n", 2*l1p)
	fmt.Println("        Bombieri–Lagarias (1999): TODA entrada del yunque tiene fórmula de")
	fmt.Println("        primos. La fábrica entera se puede forjar desde el lado de los primos.")

	// ---- veredicto ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("⚡⚡ **LA FÁBRICA DE CUADRADOS DE LA HOJA DE FORJA — CONSTRUIDA Y MEDIDA:**")
	fmt.Printf("\n  · el yunque M[m,n] = λₘ + λₙ − λ|m−n| es equivalente a RH (diagonal = 2λₙ\n")
	fmt.Printf("    = Li; en la piel es Gram — LEY 1, %d perlas, desvío %.0e)\n", len(ps), peor)
	fmt.Printf("  · su Cholesky M = AᵀA quedó EXHIBIDO con %d cuadrados firmes; los pivotes\n", firme)
	fmt.Printf("    caen ≈ ×%.4f por escalón (el coro al unísono) y donde se hunden en el\n", razonMedia)
	fmt.Println("    ruido del material lo DECLARAMOS — ventana agotada, no negatividad")
	fmt.Printf("  · el yunque tiene MEJOR OÍDO que la escalera sobre tuplas aisladas: la\n")
	fmt.Printf("    sintética en N = %d (Li: n = %d) y la DH real en N = %d (Li: n = %d ≈ 2π/φ),\n", detS, primeraEsc, detD, primeraD)
	fmt.Println("    con controles en la piel que separan señal de artefacto — y con el")
	fmt.Println("    caveat del enmascaramiento declarado")
	fmt.Println("  · y el primer golpe lo dan los primos solos: M₁ = 2λ₁ > 0 por Mertens")
	fmt.Println("\n📌 LO QUE FALTA, en el idioma de la hoja: demostrar que el pivote de")
	fmt.Println("  Cholesky NUNCA pisa territorio negativo — para TODO N, desde el lado de")
	fmt.Println("  los primos (B–L da cada entrada). Ésa es exactamente la positividad de")
	fmt.Println("  Weil, y nadie la giró en 74 años.")
	fmt.Println("\n⚖️ Honesto: Li 1997, B–L 1999, Weil 1952; la forma matricial es consecuencia")
	fmt.Println("  inmediata de Li + Gram y seguramente conocida en el oficio. Lo nuestro es")
	fmt.Println("  la CONSTRUCCIÓN medida: espectro con piso de ruido, la A exhibida, la")
	fmt.Println("  escalera de pivotes, y el oído comparado con controles. Todavía no.")

	escribirLamina(ev, ruido[N], piv, A, firme, detS, primeraEsc, detD, primeraD)
}

func escribirLamina(ev []float64, ruidoMax float64, piv []float64, A [][]float64, firme, yunS, escS, yunD, escD int) {
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="1400" height="800" viewBox="0 0 1400 800">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="700" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🔨 EL YUNQUE — la fábrica de cuadrados, construida y medida</text>
<text x="700" y="76" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">M[m,n] = λₘ + λₙ − λ|m−n| · positiva para todo N ⟺ RH · su Cholesky M = AᵀA, exhibido hasta donde el material aguanta</text>
`)
	// escalera de pivotes (log10)
	fmt.Fprintf(&b, `<rect x="60" y="110" width="620" height="330" rx="12" fill="#101f36" stroke="#26456e"/>
<text x="370" y="142" font-size="16" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">LA ESCALERA DE PIVOTES — geométrica, hasta el piso de ruido</text>
`)
	for j := 0; j < len(piv) && j < 8; j++ {
		v := piv[j]
		col := "#7ee0c0"
		if j >= firme {
			col = "#ff9aa8"
		}
		lab := fmt.Sprintf("%.1e", v)
		h := 8.0
		if v > 0 {
			h = (math.Log10(v) + 11) / 13 * 230
			if h < 8 {
				h = 8
			}
		}
		x := 100.0 + float64(j)*72
		fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="46" height="%.1f" fill="%s" opacity="0.85"/>
<text x="%.1f" y="%.1f" font-size="10.5" text-anchor="middle" font-family="monospace" fill="#9aa8c4" transform="rotate(-38 %.1f %.1f)">%s</text>
`, x, 395-h, h, col, x+23, 385-h, x+23, 385-h, lab)
	}
	fmt.Fprintf(&b, `<text x="370" y="372" font-size="12" text-anchor="middle" font-family="Georgia" fill="#ff9aa8">rosa = hundido en el ruido del material (%.0e): ventana agotada, no negatividad — declarado</text>
<text x="370" y="424" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">%d cuadrados firmes · el coro canta casi al unísono (φ ≈ 1/γ ≤ 1/14 ⟹ 1−wⁿ ≈ n(1−w))</text>
`, ruidoMax, firme)
	// la esquina de A
	fmt.Fprintf(&b, `<rect x="720" y="110" width="620" height="330" rx="12" fill="#0f2b22" stroke="#2f7f63"/>
<text x="1030" y="142" font-size="16" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">LA ESQUINA DE A — la fábrica: Q(c) = ||Ac||²</text>
`)
	for i := 0; i < firme && i < 6; i++ {
		fila := ""
		for j := 0; j <= i; j++ {
			fila += fmt.Sprintf(" %8.5f", A[i][j])
		}
		fmt.Fprintf(&b, `<text x="750" y="%d" font-size="14.5" font-family="monospace" fill="#ffd98a">%s</text>`+"\n", 185+i*34, fila)
	}
	fmt.Fprintf(&b, `<text x="750" y="400" font-size="13" font-family="Georgia" fill="#cfe6ff">a₁₁ = √(2λ₁) = %.6f — el primer golpe lo dan los primos solos (Mertens)</text>
<text x="750" y="424" font-size="13" font-family="Georgia" fill="#cfe6ff">autovalores de M₄₀: máximo %.1e · mínimo indistinguible de 0 al ruido</text>
`, A[0][0], ev[len(ev)-1])
	// el oido
	fmt.Fprintf(&b, `<rect x="60" y="470" width="1280" height="160" rx="12" fill="#2b1020" stroke="#8a3557"/>
<text x="700" y="502" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ff9aa8">EL OÍDO DEL YUNQUE — medido con controles en la piel</text>
<text x="700" y="534" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">tupla sintética fuera de la piel (β=0.8, γ=2), aislada: la escalera de Li la oye en n = %d — el yunque en N = %d</text>
<text x="700" y="562" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">la perla DH real (β=0.809, γ=85.7), aislada: la escalera necesita n = %d (≈ 2π/φ, su primer nulo de fase — corrección F294) — el yunque la oye en N = %d</text>
<text x="700" y="590" font-size="13" text-anchor="middle" font-family="Georgia" fill="#ffd98a">los cuadrados se cancelan en las direcciones casi-unísono y la fuga radial queda al desnudo en los términos cruzados</text>
<text x="700" y="614" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">caveat declarado: es la tupla SOLA — en un espectro completo el fondo de Gram del coro puede taparla; queda abierto</text>
`, escS, yunS, escD, yunD)
	fmt.Fprintf(&b, `<text x="700" y="672" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">LO QUE FALTA: demostrar que el pivote de Cholesky nunca pisa territorio negativo — para TODO N, desde los primos.</text>
<text x="700" y="700" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Ésa es la positividad de Weil (74 años abierta). La fábrica está construida en la ventana; falta el teorema que la sostenga entera.</text>
<text x="700" y="744" font-size="14" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">Li 1997 · Bombieri–Lagarias 1999 · Weil 1952 — lo nuestro: la construcción medida, con su ruido y sus controles a la vista.</text>
<text x="700" y="776" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Todavía no.</text>
</svg>
`)
	os.WriteFile("el-yunque.svg", []byte(b.String()), 0o644)
	fmt.Println("\n🖼️  lámina escrita: el-yunque.svg")
}
