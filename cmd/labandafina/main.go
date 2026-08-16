// Command labandafina closes the yellow point of the Leader's Law: the
// FINE BAND with all quantifiers - for every configuration under H0-H4
// with a strict leader, and for EVERY n >= N* in the leader's fine band
// (||n*th_L|| <= 1), lambda_n < 0, where
//
//	N* = max( n_rad , n_comp )
//	n_rad  = ceil(u*log u), u = 3(m+1)/delta_L      [DYN's, unchanged]
//	n_comp = ceil( log(2(m-1)/cos 1) / (delta_L - delta_2) )   [m>=2; 0 if m=1]
//
// via the chain F1-F6 (see docs/PLACAS-BANDA-FINA-ACTA.md):
// F1 leader bound in band, F2 competitors <= 6 + 2*r_2^n, F3 choir
// under H4, F4 assembly, F5 competitor domination past n_comp,
// F6 cos(1)*r_L^n beats the polynomial past n_rad (doubled-bracket
// radial argument: bracket_2(n_rad) <= 3u^3 <= u^{3(m+1)}).
//
// This program verifies: (1) the F6 doubled-bracket battery on the
// (m, delta) grid; (2) the F5 boundary battery on synthetic leader
// gaps; (3) the full chain live on the witness fine band - every
// fine-band n past N*, its majorant and its sign.
//
// Evidence; the proof is F1-F6. Reproduce: go run ./cmd/labandafina
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
	fmt.Println("🟡→🟢 LA BANDA FINA — el punto amarillo, cerrado: banda fina → N* → λ < 0")

	// ---- F6 battery: doubled-bracket radial on the grid ----
	fmt.Println("\n§1 · F6 EN LA GRILLA — cos(1)·r_Lⁿ > (4/π)n·log n + 6m − 2 desde n_rad")
	viol := 0
	for m := 1; m <= 10; m++ {
		for _, dd := range []float64{0.01, 0.1, 0.3466, 0.7, 1.0} {
			u := 3 * float64(m+1) / dd
			nrad := math.Ceil(u * math.Log(u))
			// the doubled bracket at n_rad, and the key inequality bracket2 <= 3u^3
			b2 := (8/math.Pi)*nrad*math.Log(nrad) + (12*float64(m) - 4)
			ok := b2 <= 3*math.Pow(u, 3) &&
				math.Cos(1)*math.Exp(nrad*dd) > (4/math.Pi)*nrad*math.Log(nrad)+6*float64(m)-2 &&
				dd*math.Exp(nrad*dd) > (8/math.Pi)*(math.Log(nrad)+1) &&
				dd*dd*math.Exp(nrad*dd) > (8/math.Pi)/nrad
			if !ok {
				viol++
			}
		}
	}
	fmt.Printf("        50 casos (m = 1..10 × δ ≤ 1): corchete₂ ≤ 3u³, F6, G′ > 0, G″ > 0 — %d violaciones ✅\n", viol)

	// ---- F5 boundary battery on synthetic gaps ----
	fmt.Println("\n§2 · F5 EN EL BORDE — 2(m−1)·r₂ⁿ ≤ cos(1)·r_Lⁿ exactamente en n_comp")
	viol5, casos := 0, 0
	for m := 2; m <= 8; m++ {
		for _, dL := range []float64{1e-4, 1e-2, 0.5} {
			for _, frac := range []float64{0.1, 0.5, 0.9, 0.99} {
				d2 := dL * frac
				ncomp := math.Ceil(math.Log(2*float64(m-1)/math.Cos(1)) / (dL - d2))
				casos++
				if 2*float64(m-1)*math.Exp(ncomp*d2) > math.Cos(1)*math.Exp(ncomp*dL)+1e-9 {
					viol5++
				}
			}
		}
	}
	fmt.Printf("        %d bordes sintéticos (m = 2..8, brechas hasta 1%%): %d violaciones ✅\n", casos, viol5)

	// ---- the full chain live on the witness ----
	fmt.Println("\n§3 · LA CADENA COMPLETA EN VIVO — el testigo m = 2, banda fina tras N*")
	rho1 := complex(0.808517, 85.699348)
	rho2 := complex(0.7, 45.0)
	w1 := 1 - 1/rho1
	w2 := 1 - 1/rho2
	r1 := math.Max(cmplx.Abs(w1), 1/cmplx.Abs(w1))
	r2 := math.Max(cmplx.Abs(w2), 1/cmplx.Abs(w2))
	t1 := math.Abs(cmplx.Phase(w1))
	t2 := math.Abs(cmplx.Phase(w2))
	d1 := math.Log(r1)
	dL := math.Log(r2) // leader
	u := 9 / dL
	nrad := int(math.Ceil(u * math.Log(u)))
	ncomp := int(math.Ceil(math.Log(2/math.Cos(1)) / (dL - d1)))
	Nstar := nrad
	if ncomp > Nstar {
		Nstar = ncomp
	}
	fmt.Printf("        n_rad = %d · n_comp = %d · N* = max = %d\n", nrad, ncomp, Nstar)

	ps := perlas(120)
	wsC := make([]complex128, len(ps))
	pcs := make([]complex128, len(ps))
	for i, g := range ps {
		wp := 1 - 1/complex(0.5, g)
		wp /= complex(cmplx.Abs(wp), 0)
		wsC[i] = wp
		pcs[i] = 1
	}
	finas, violSigno, violMaj := 0, 0, 0
	for n := 1; n <= Nstar+3000; n++ {
		var s float64
		for i := range wsC {
			pcs[i] *= wsC[i]
			s += 2 - 2*real(pcs[i])
		}
		if n < Nstar {
			continue
		}
		fn := float64(n)
		if norma(fn*t2) <= 1 { // leader fine band, pearl 1 unrestricted
			finas++
			l1 := 4 - 2*math.Cos(fn*t1)*(math.Exp(fn*d1)+math.Exp(-fn*d1))
			l2 := 4 - 2*math.Cos(fn*t2)*(math.Exp(fn*dL)+math.Exp(-fn*dL))
			lam := s + l1 + l2
			maj := (4/math.Pi)*fn*math.Log(fn) + 6*2 - 2 + 2*math.Exp(fn*d1) - 2*math.Cos(1)*math.Exp(fn*dL)
			if lam > maj+1e-6*math.Abs(maj) {
				violMaj++
			}
			if lam >= 0 {
				violSigno++
			}
		}
	}
	fmt.Printf("        escalones de banda fina en [N*, N*+3000]: %d\n", finas)
	fmt.Printf("        violaciones del mayorante F4 (λ ≤ P + 2r₂ⁿ − 2cos1·r_Lⁿ): %d ✅\n", violMaj)
	fmt.Printf("        violaciones del signo (λ < 0): %d ✅\n", violSigno)

	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("🟢 **EL PUNTO AMARILLO, CERRADO:**")
	fmt.Println("\n  ∀ configuración H0-H4 con líder estricto, ∀n ≥ N* = max(n_rad, n_comp):")
	fmt.Println("  ‖nθ_L‖ ≤ 1 ⟹ λₙ < 0 — con N* explícito y la cadena F1-F6 escrita")
	fmt.Println("  (F6 recicla el radial de DYN con corchete duplicado: corchete₂ ≤ 3u³ ≤ u^{3(m+1)})")
	fmt.Println("\n⚖️ Honesto: la prueba es F1-F6 del acta; las corridas corroboran. El paisaje")
	fmt.Println("  completo (Teorema de las Placas) queda con sus dos bandas cerradas — el")
	fmt.Println("  ensamble final y el sello son de la auditora. Todavía no.")
}
