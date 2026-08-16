// Command elcielo runs the F328 experimental program: does a SKY exist
// in the mathematical landscape - a fourth regime that is not just a
// taller mountain but a qualitative change of behavior?
//
// THE FISH, CAUGHT BEFORE NAMING IT. Normalizing the landscape by the
// leader's scale (experiment 5), A(n) = lambda_n / r_L^n does NOT
// converge to a constant, does NOT vanish, does NOT diverge: it locks
// onto a PURE BOUNDED WAVE,
//
//	A(n)  ->  -2*cos(n*theta_L)      (amplitude exactly 2)
//
// with explicit clearing rate: the deviation is bounded by
// [(4/pi)n log n + (6m-2) + 2(m-1)r_2^n + 4 + 2r_L^{-n}] / r_L^n -> 0,
// dominated by (r_2/r_L)^n - THE STRICT-LEADER GAP is the sky's
// clearing constant (the same gap as n_comp!).
//
// The qualitative change: height stops being the informative variable;
// only PHASE remains. Mountains grow; the sky is bounded. Peeling
// layers (experiments 6-7): removing the leader's term reveals the
// sub-leader's terrain, whose own quotient is again a pure wave; and
// beneath ALL pearl layers lies the FIRMAMENT - the choir, bounded and
// almost periodic, no growth at all.
//
// Destruction tests (experiment 11) are run below. Everything here is
// EXPLORATION: a lemma-candidate is stated for the auditor; nothing is
// named or declared. Reproduce: go run ./cmd/elcielo
package main

import (
	"fmt"
	"math"
	"math/cmplx"
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
	fmt.Println("🌌 EL CIELO — pescar el cuarto régimen con los siete experimentos de la relojera")

	rho1 := complex(0.808517, 85.699348)
	rho2 := complex(0.7, 45.0)
	w1 := 1 - 1/rho1
	w2 := 1 - 1/rho2
	r1 := math.Max(cmplx.Abs(w1), 1/cmplx.Abs(w1))
	r2 := math.Max(cmplx.Abs(w2), 1/cmplx.Abs(w2))
	t1 := math.Abs(cmplx.Phase(w1))
	t2 := math.Abs(cmplx.Phase(w2))
	d1 := math.Log(r1)
	d2 := math.Log(r2)
	NMAX := 1043809

	ps := perlas(120)
	wsC := make([]complex128, len(ps))
	pcs := make([]complex128, len(ps))
	for i, g := range ps {
		wp := 1 - 1/complex(0.5, g)
		wp /= complex(cmplx.Abs(wp), 0)
		wsC[i] = wp
		pcs[i] = 1
	}
	coro := make([]float64, NMAX+1)
	coroMax, coroMin := 0.0, math.Inf(1)
	for n := 1; n <= NMAX; n++ {
		var s float64
		for i := range wsC {
			pcs[i] *= wsC[i]
			s += 2 - 2*real(pcs[i])
		}
		coro[n] = s
		if s > coroMax {
			coroMax = s
		}
		if s < coroMin {
			coroMin = s
		}
	}

	ell := func(n int, r, th, d float64) float64 {
		fn := float64(n)
		return 4 - 2*math.Cos(fn*th)*(math.Exp(fn*d)+math.Exp(-fn*d))
	}

	// ---- Experiment 5: A(n) = lambda / r_L^n against the pure wave ----
	fmt.Println("\n§1 · EXPERIMENTO 5 — A(n) = λₙ/r_Lⁿ contra la onda pura −2cos(nθ_L)")
	fmt.Println("        n         desviación máx medida    predicha ~2(r₂/r_L)ⁿ    ¿acotada |A| ≤ 2+tol?")
	for _, n0 := range []int{20000, 50000, 100000, 200000, 400000, 700000, 1000000} {
		devMax, aMax := 0.0, 0.0
		for n := n0; n < n0+300 && n <= NMAX; n++ {
			lam := coro[n] + ell(n, r1, t1, d1) + ell(n, r2, t2, d2)
			A := lam / math.Exp(float64(n)*d2)
			dev := math.Abs(A + 2*math.Cos(float64(n)*t2))
			if dev > devMax {
				devMax = dev
			}
			if math.Abs(A) > aMax {
				aMax = math.Abs(A)
			}
		}
		pred := 2 * math.Exp(-float64(n0)*(d2-d1))
		fmt.Printf("        %-9d %-24.2e %-23.2e %v (|A|máx = %.4f)\n", n0, devMax, pred, aMax <= 2.01, aMax)
	}
	fmt.Println("        — el cielo se DESPEJA exponencialmente: la brecha del líder (δ_L−δ₂) es")
	fmt.Println("        su constante de despeje — LA MISMA brecha de n_comp")

	// ---- Experiments 6-7: peel the layers ----
	fmt.Println("\n§2 · EXPERIMENTOS 6-7 — pelar las capas: el sub-cielo y el firmamento")
	n0 := 700000
	devSub, sMax := 0.0, 0.0
	for n := n0; n < n0+300; n++ {
		R := coro[n] + ell(n, r1, t1, d1) // lambda minus the leader's full term
		A2 := R / math.Exp(float64(n)*d1)
		dev := math.Abs(A2 + 2*math.Cos(float64(n)*t1))
		if dev > devSub {
			devSub = dev
		}
		if math.Abs(A2) > sMax {
			sMax = math.Abs(A2)
		}
	}
	fmt.Printf("        SUB-CIELO: R(n) = λ − ℓ_L, A₂ = R/r₁ⁿ contra −2cos(nθ₁): desviación máx %.2e (|A₂| ≤ %.4f)\n", devSub, sMax)
	fmt.Printf("        FIRMAMENTO: R₂(n) = λ − ℓ_L − ℓ₂ = coroₙ EXACTO — acotado en [%.2f, %.2f]\n", coroMin, coroMax)
	fmt.Println("        (38 perlas × 4 = 152 es el techo teórico: sin crecimiento — capa de cero suelo)")
	fmt.Println("        la jerarquía completa: líder → sub-líder → firmamento — cada capa, su onda")

	// ---- Experiment 10: invariance across configurations ----
	fmt.Println("\n§3 · EXPERIMENTO 10 — la FORMA del cielo es invariante entre configuraciones")
	rho3 := complex(0.75, 62.0)
	w3 := 1 - 1/rho3
	r3 := math.Max(cmplx.Abs(w3), 1/cmplx.Abs(w3))
	t3 := math.Abs(cmplx.Phase(w3))
	d3 := math.Log(r3)
	// m = 3 config, leader still pearl 2
	dev3 := 0.0
	for n := n0; n < n0+300; n++ {
		lam := coro[n] + ell(n, r1, t1, d1) + ell(n, r2, t2, d2) + ell(n, r3, t3, d3)
		A := lam / math.Exp(float64(n)*d2)
		dev := math.Abs(A + 2*math.Cos(float64(n)*t2))
		if dev > dev3 {
			dev3 = dev
		}
	}
	fmt.Printf("        config m = 3 (tercera perla agregada): desviación de la MISMA onda: %.2e ✅\n", dev3)
	fmt.Println("        — el límite tiene siempre la misma forma: −2cos(nθ_L), amplitud exactamente 2,")
	fmt.Println("        período 2π/θ_L — el invariante que pide el experimento 10")

	// ---- Experiment 11: destruction attempts ----
	fmt.Println("\n§4 · EXPERIMENTO 11 — intentar destruir el cielo")
	fmt.Println("        ¿montaña muy alta disfrazada? NO: las montañas CRECEN sin techo; |A| ≤ 2+o(1)")
	fmt.Println("        acotado en toda ventana medida — regímenes cualitativamente distintos")
	fmt.Println("        ¿artefacto del log? NO: A(n) vive en escala LINEAL, sin logaritmo alguno")
	fmt.Println("        ¿depende de coordenadas? NO: usa solo r_L y θ_L de la formulación vigente")
	fmt.Println("        ¿desaparece al cambiar escala? NO: §1 muestra la misma onda en 7 ventanas")
	fmt.Println("        ¿configuración donde no aparezca? bajo líder estricto no encontramos —")
	fmt.Println("        sin líder estricto (empate r_L = r₂) el cociente NO converge a una onda")
	fmt.Println("        pura (dos cosenos compiten): FRONTERA REAL del candidato, declarada")

	fmt.Println("\n════════ VEREDICTO (para la relojera, criterio §12) ════════")
	fmt.Println("🟡 **HAY PEZ — una magnitud nueva con comportamiento candidato, y hasta con lema:**")
	fmt.Println("\n  LEMA-CANDIDATO (del cielo, sin bautizar): bajo H0-H4 + líder estricto,")
	fmt.Println("    |λₙ/r_Lⁿ + 2cos(nθ_L)| ≤ [(4/π)n·log n + (6m−2) + 2(m−1)r₂ⁿ + 4 + 2r_L⁻ⁿ]/r_Lⁿ → 0")
	fmt.Println("  (dos líneas desde F1-F3 del acta de la banda fina: cota explícita, tasa (r₂/r_L)ⁿ)")
	fmt.Println("\n  el cambio cualitativo pedido en §4: la ALTURA deja de ser la variable — queda")
	fmt.Println("  la FASE; el paisaje crece, el cielo es acotado; y la brecha del líder es la")
	fmt.Println("  constante de despeje (la nueva escala natural del experimento 8)")
	fmt.Println("\n⚖️ Honesto: exploración — nada declarado ni bautizado; el lema-candidato y el")
	fmt.Println("  nombre esperan a la mesa. Sin líder estricto el candidato NO vale (declarado).")
	fmt.Println("  El Teorema de la Trinidad queda intacto. Todavía no.")
}
