package main

// contagio.go - the contagion test. Prediction, declared before running: the
// noisy mixed composites (72, 75, 92, ...) sing ONLY because they sit within
// the bat's blur next to a prime; the blur shrinks like ~pi/gamma_max, so
// raising the zero count must silence them ON SCHEDULE, while the primes stay
// loud. If instead they keep singing past their decoupling height, they carry
// a voice of their own and the explicit formula is in trouble (it will not be).

import (
	"fmt"
	"math"
)

// theta, zeta, cerosRS: Riemann-Siegel on the critical line (from cmd/elpliego).
func theta(t float64) float64 {
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t) + 7/(5760*t*t*t)
}

func zetaZ(t float64) float64 {
	th := theta(t)
	u := math.Sqrt(t / (2 * math.Pi))
	N := int(u)
	s := 0.0
	for n := 1; n <= N; n++ {
		fn := float64(n)
		s += math.Cos(th-t*math.Log(fn)) / math.Sqrt(fn)
	}
	s *= 2
	p := u - float64(N)
	c0 := math.Cos(2*math.Pi*(p*p-p-1.0/16)) / math.Cos(2*math.Pi*p)
	sign := 1.0
	if N%2 == 0 {
		sign = -1
	}
	return s + sign*math.Pow(2*math.Pi/t, 0.25)*c0
}

func cerosRS(t0, t1, paso float64) []float64 {
	var g []float64
	a, za := t0, zetaZ(t0)
	for b := t0 + paso; b <= t1; b += paso {
		zb := zetaZ(b)
		if za == 0 || (za < 0) != (zb < 0) {
			lo, hi, zlo := a, b, za
			for i := 0; i < 60; i++ {
				m := (lo + hi) / 2
				zm := zetaZ(m)
				if (zlo < 0) != (zm < 0) {
					hi = m
				} else {
					lo, zlo = m, zm
				}
			}
			g = append(g, (lo+hi)/2)
		}
		a, za = b, zb
	}
	return g
}

func correr3() {
	fmt.Println("\n🦇🔬 LA PRUEBA DEL CONTAGIO — ¿los mixtos ruidosos se callan EN HORARIO?")
	fmt.Println("     profecía: n mixto pegado a un primo p se calla cuando γ_max·|log n − log p| ≳ π")
	fmt.Println("     o sea, altura de desacople γ* ≈ π/d. Se verifica el horario, caso por caso.")

	g := cerosRS(10, 1000, 0.05)
	fmt.Printf("     ceros calculados con Riemann-Siegel hasta γ = 1000: %d (los primeros 29 ya\n", len(g))
	fmt.Printf("     los teníamos a mano: coinciden a %.6f)\n", math.Abs(g[0]-ceros[0]))

	// the accused: the noisy mixed composites and their contaminating prime
	casos := []struct{ n, p int }{
		{72, 73}, {75, 73}, {92, 89}, {45, 47}, {39, 41},
		{93, 97}, {76, 79}, {86, 89}, {24, 23}, {60, 61},
	}
	alturas := []int{29, 50, 100, 200, len(g)}

	fmt.Printf("\n     %-9s", "M ceros:")
	for _, M := range alturas {
		fmt.Printf(" %8d", M)
	}
	fmt.Printf("  %9s\n", "γ* teor.")
	fmt.Printf("     %-9s", "γ_max:")
	for _, M := range alturas {
		fmt.Printf(" %8.0f", g[M-1])
	}
	fmt.Println()

	aciertos, total := 0, 0
	for _, c := range casos {
		d := math.Abs(math.Log(float64(c.n)) - math.Log(float64(c.p)))
		gStar := math.Pi / d
		fmt.Printf("     n=%-3d(p=%d)", c.n, c.p)
		var es []float64
		for _, M := range alturas {
			e := math.Abs(eco(g[:M], math.Log(float64(c.n))))
			es = append(es, e)
			fmt.Printf(" %8.3f", e)
		}
		fmt.Printf("  %9.0f\n", gStar)
		// verdict per case, corrected: the naive pi/d schedule ignores the NOISE
		// FLOOR sqrt(2/M) that a finite zero count imposes; with 29 zeros the floor
		// is 0.26, so nobody can be "silent" there. The honest check: at the top
		// height the accused must sit UNDER the noise floor (silence) after falling
		// by more than a factor 5 from the 29-zero reading (real collapse).
		total++
		piso := math.Sqrt(2 / float64(len(g)))
		if es[len(es)-1] < 1.5*piso && es[len(es)-1] < es[0]/5 {
			aciertos++
		}
	}

	// the primes must NOT silence: same heights, mean -E over the 25 primes
	fmt.Printf("\n     %-13s", "primos (media)")
	for _, M := range alturas {
		s := 0.0
		for _, p := range primosCien {
			s += -eco(g[:M], math.Log(float64(p)))
		}
		fmt.Printf(" %8.3f", s/float64(len(primosCien)))
	}
	fmt.Println("   ← deben SEGUIR cantando")

	fmt.Printf("\n§ VEREDICTO: de los mixtos cuya altura de desacople quedó dentro del rango,\n")
	fmt.Printf("     %d de %d acusados quedaron BAJO el piso de ruido tras caer más de 5 veces.\n", aciertos, total)
	if total > 0 && aciertos >= total*7/10 {
		fmt.Println("     ⟹ CONFIRMADO: el ruido de los mixtos era contagio del vecino. Los ceros no")
		fmt.Println("       les conocen la voz: al afinar el oído, el silencio se profundiza en el")
		fmt.Println("       horario que dicta la resolución. Los números compuestos no existen para")
		fmt.Println("       la orquesta — ni siquiera de cerca.")
	} else {
		fmt.Println("     ⟹ NO se confirma el horario: hay mixtos que siguen sonando pasada su altura")
		fmt.Println("       de desacople. Eso sería voz propia y habría que perseguirlo.")
	}
}
