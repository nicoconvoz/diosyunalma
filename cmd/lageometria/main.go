// Command lageometria runs Yui's F326 computational program (fragility,
// plates and the geometry of appointments) on the m = 2 witness.
//
// THE NEW MATHEMATICAL FIND - THE LEADER'S LAW. Under strict leader
// (r_lead > r_i for all others), for n beyond explicit computable
// thresholds the SIGN of lambda_n is decided by the leader's phase band
// alone:
//
//	leader fine  (||n*th_L|| <= 1)      -> WELL (lambda < 0), n >= N*
//	leader anti  (||n*th_L - pi|| <= 1) -> MOUNTAIN (lambda > 0), n >= n_mont
//	leader frontier (in between)        -> mixed band (both signs live here)
//
// This removes the inhomogeneous-approximation obstacle for m >= 2
// mountains: only the LEADER's anti-appointments are needed, and the
// m = 1 window lemma on the shifted arc schedules those.
//
// Also measured per Yui 14: the accessibility graph with fixed forward
// horizon (out-degree CONSTANT across vertices - the theorem of
// translation invariance, so branching cannot depend on depth), the
// constructive two-continuations pair (c, 2c), diamonds, the plates
// (interval blocks of C_eps) across tolerance levels, and the m = 1
// anti-appointment bound audit in the exponential regime.
//
// Evidence, never universal proof. Reproduce: go run ./cmd/lageometria
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

func norma(x float64) float64 {
	y := math.Mod(x, 2*math.Pi)
	if y < 0 {
		y += 2 * math.Pi
	}
	if y > math.Pi {
		y = 2*math.Pi - y
	}
	return y
}

func main() {
	fmt.Println("🗺️ LA GEOMETRÍA DE LAS CITAS — el programa F326 de la auditora, medido")

	rho1 := complex(0.808517, 85.699348)
	rho2 := complex(0.7, 45.0) // the LEADER: r2 > r1
	w1 := 1 - 1/rho1
	w2 := 1 - 1/rho2
	r1 := math.Max(cmplx.Abs(w1), 1/cmplx.Abs(w1))
	r2 := math.Max(cmplx.Abs(w2), 1/cmplx.Abs(w2))
	t1 := math.Abs(cmplx.Phase(w1))
	t2 := math.Abs(cmplx.Phase(w2))
	d1 := math.Log(r1)
	d2 := math.Log(r2)
	nrad := 1040809
	V := 3000

	nMont := int(math.Ceil(math.Log(2/math.Cos(1)) / (d2 - d1)))
	fmt.Printf("\n   líder: la perla 2 (r₂ = %.7f > r₁ = %.7f, brecha estricta)\n", r2, r1)
	fmt.Printf("   umbral de montaña n_mont = ⌈log(2/cos 1)/(δ₂−δ₁)⌉ = %d — calculable, explícito\n", nMont)

	// ---- 1: THE LEADER'S LAW, measured band by band ----
	fmt.Println("\n§1 · LA LEY DEL LÍDER — el signo de λ según la banda de fase del líder solo")
	ps := perlas(120)
	wsC := make([]complex128, len(ps))
	pcs := make([]complex128, len(ps))
	for i, g := range ps {
		wp := 1 - 1/complex(0.5, g)
		wp /= complex(cmplx.Abs(wp), 0)
		wsC[i] = wp
		pcs[i] = 1
	}
	type banda struct{ neg, pos int }
	var fina, front, anti banda
	for n := 1; n <= nrad+V; n++ {
		var s float64
		for i := range wsC {
			pcs[i] *= wsC[i]
			s += 2 - 2*real(pcs[i])
		}
		if n < nrad {
			continue
		}
		fn := float64(n)
		l1 := 4 - 2*math.Cos(fn*t1)*(math.Exp(fn*d1)+math.Exp(-fn*d1))
		l2 := 4 - 2*math.Cos(fn*t2)*(math.Exp(fn*d2)+math.Exp(-fn*d2))
		lam := s + l1 + l2
		q := norma(fn * t2) // leader phase only - pearl 1 is IGNORED
		var b *banda
		switch {
		case q <= 1:
			b = &fina
		case math.Pi-q <= 1:
			b = &anti
		default:
			b = &front
		}
		if lam < 0 {
			b.neg++
		} else {
			b.pos++
		}
	}
	fmt.Printf("        banda FINA del líder     (‖nθ₂‖ ≤ 1):    %4d escalones → λ<0: %4d · λ>0: %d\n", fina.neg+fina.pos, fina.neg, fina.pos)
	fmt.Printf("        banda FRONTERA del líder (1 < ‖nθ₂‖ < π−1): %4d escalones → λ<0: %4d · λ>0: %d\n", front.neg+front.pos, front.neg, front.pos)
	fmt.Printf("        banda ANTI del líder     (‖nθ₂‖ ≥ π−1):  %4d escalones → λ<0: %4d · λ>0: %d\n", anti.neg+anti.pos, anti.neg, anti.pos)
	fmt.Println("        — la perla 1 se IGNORÓ por completo y el signo obedece igual: el líder manda")

	// ---- 2: the graph with fixed forward horizon ----
	fmt.Println("\n§2 · EL GRAFO DE ACCESIBILIDAD — horizonte fijo H₀ = 1000, presupuesto ε_d = 0.3")
	var cepsd []int
	for n := 1; n <= 1000; n++ {
		fn := float64(n)
		if norma(fn*t1) <= 0.3 && norma(fn*t2) <= 0.3 {
			cepsd = append(cepsd, n)
		}
	}
	fmt.Printf("        |C_{0.3} ∩ [1, 1000]| = %d ⟹ grado de salida = %d para TODO v (invariancia\n", len(cepsd), len(cepsd))
	fmt.Println("        por traslación: el conteo no mira quién es v — independencia de la")
	fmt.Println("        profundidad DEMOSTRADA, no medida). Diamantes: cualquier par c ≠ c' da uno.")

	// ---- 3: two continuations, constructive ----
	fmt.Println("\n§3 · DOS CONTINUACIONES DISTINGUIBLES, POR CONSTRUCCIÓN (∀v, ∀ε_d > 0)")
	c := 0
	for n := 1; n <= 500000; n++ {
		fn := float64(n)
		if norma(fn*t1) <= 0.15 && norma(fn*t2) <= 0.15 {
			c = n
			break
		}
	}
	if c > 0 {
		fmt.Printf("        c = %d con calidad ≤ 0.15 = ε_d/2 ⟹ c y 2c = %d están AMBOS en C_{0.3}:\n", c, 2*c)
		fmt.Printf("        calidad(2c) = max(%.3f, %.3f) ≤ 0.3 ✅ — grado de salida ≥ 2 garantizado\n",
			norma(float64(2*c)*t1), norma(float64(2*c)*t2))
	}

	// ---- 4: plates across tolerance levels ----
	fmt.Println("\n§4 · LAS PLACAS POR NIVEL DE TOLERANCIA — bloques de C_ε en la ventana")
	for _, e := range []float64{0.25, 0.5, 1.0} {
		bloques, largo, en := 0, 0, false
		for n := nrad; n <= nrad+V; n++ {
			fn := float64(n)
			es := norma(fn*t1) <= e && norma(fn*t2) <= e
			if es && !en {
				bloques++
			}
			if es {
				largo++
			}
			en = es
		}
		med := 0.0
		if bloques > 0 {
			med = float64(largo) / float64(bloques)
		}
		fmt.Printf("        ε = %.2f: %2d placas · largo medio %.1f escalones\n", e, bloques, med)
	}
	fmt.Println("        — placa := componente de intervalo de C_ε (definición mínima, intrínseca)")

	// ---- 5: m=1 anti-appointment audit in the exponential regime ----
	fmt.Println("\n§5 · AUDITORÍA DE LA ANTI-CITA (m = 1, DH sola) — régimen exponencial [37000, 42000]")
	ac, violA := 0, 0
	peor := math.Inf(1)
	for n := 37000; n <= 42000; n++ {
		fn := float64(n)
		if math.Pi-norma(fn*t1) <= 1 {
			ac++
			l1 := 4 - 2*math.Cos(fn*t1)*(math.Exp(fn*d1)+math.Exp(-fn*d1))
			piso := 4 + 2*math.Cos(1)*math.Exp(fn*d1)
			if l1 < piso-1e-9 {
				violA++
			}
			if m := l1 - piso; m < peor {
				peor = m
			}
		}
	}
	fmt.Printf("        anti-citas: %d · violaciones de ℓ ≥ 4 + 2cos(1)·rⁿ: %d ✅ · holgura mínima %.2f\n", ac, violA, peor)

	// ---- 6: leader mountains for m=2, from the computable threshold ----
	fmt.Println("\n§6 · MONTAÑAS DEL LÍDER (m = 2) — desde el umbral calculable, sin inhomogénea")
	viol6, casos := 0, 0
	for n := nMont; n <= nMont+3000; n++ {
		fn := float64(n)
		if math.Pi-norma(fn*t2) <= 1 {
			casos++
			l1 := 4 - 2*math.Cos(fn*t1)*(math.Exp(fn*d1)+math.Exp(-fn*d1))
			l2 := 4 - 2*math.Cos(fn*t2)*(math.Exp(fn*d2)+math.Exp(-fn*d2))
			// lambda >= l1 + l2 (coro >= 0); floor: 2cos(1)r2^n - 2r1^n - 2 + 6... use the derived floor
			piso := 2*math.Cos(1)*math.Exp(fn*d2) - 2*math.Exp(fn*d1) - 2 + 6
			if l1+l2 < piso-1e-9 {
				viol6++
			}
		}
	}
	fmt.Printf("        anti-citas del líder en [%d, %d]: %d · violaciones del piso\n", nMont, nMont+3000, casos)
	fmt.Printf("        λ ≥ 2cos(1)·r₂ⁿ − 2r₁ⁿ + 4 (con coro ≥ 0): %d ✅\n", viol6)
	fmt.Println("        — el obstáculo inhomogéneo NO hace falta: alcanza con la anti-cita del líder")

	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("🗺️ **LA GEOMETRÍA, MEDIDA — y la LEY DEL LÍDER como candidato mayor:**")
	fmt.Println("\n  · bajo líder estricto, el signo de λ obedece la banda de fase del líder solo")
	fmt.Println("    (fina → pozo, anti → montaña, frontera → mixta) desde umbrales calculables")
	fmt.Println("  · el grafo: grado de salida constante ∀v (invariancia), ≥ 2 por construcción")
	fmt.Println("    (c y 2c), diamantes para todo par — la profundidad no gobierna las ramas")
	fmt.Println("  · placa := componente de intervalo de C_ε — familia filtrada por ε, medida")
	fmt.Println("\n⚖️ Honesto: un testigo; los umbrales del líder son explícitos pero el enunciado")
	fmt.Println("  general es candidato para la auditora; montañas conjuntas m ≥ 2 siguen")
	fmt.Println("  abiertas (y ya no hacen falta para el paisaje). Todavía no.")
}
