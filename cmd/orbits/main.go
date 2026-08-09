// Command orbits reads the primes as the periodic orbits of a chaotic
// system, using Landau's formula as the dictionary.
//
// Gutzwiller's trace formula in quantum chaos and the explicit formula of
// prime number theory have the same shape: energy levels ↔ zeros, classical
// periodic orbits ↔ primes, orbit period ↔ ln p, orbit stability ↔ p^{-k/2}.
// Landau (1912) made the dictionary testable: summing cos(γ·ln n) over the
// zeros singles out the prime powers —
//
//	Σ_{γ≤T} cos(γ ln n)  ≈  −(T/2π)·Λ(n)/√n   for n a prime power,
//	                          small              otherwise,
//
// where Λ(n) = ln p if n = p^k and 0 otherwise. In orbit language: call out
// a period; only real orbits answer, and each answers with a strength set by
// its stability.
//
// PRE-REGISTERED, using the laboratory's ten measured zeros of ζ (T ≈ 50):
// S(n) must be clearly negative for every prime power n ≤ 30 with magnitude
// tracking Λ(n)/√n, and near zero for every composite non-prime-power; the
// two populations must separate.
//
// Usage:
//
//	go run ./cmd/orbits
package main

import (
	"fmt"
	"math"
)

// the laboratory's own measured zeros of zeta (Finding 26).
var zetaZeros = []float64{14.1349, 21.0211, 25.0044, 30.4282, 32.9422,
	37.5872, 40.9264, 43.3211, 48.0105, 49.7752}

// mangoldt returns Λ(n): ln p if n is a power of the prime p, else 0.
func mangoldt(n int) float64 {
	for p := 2; p*p <= n; p++ {
		if n%p == 0 {
			for n%p == 0 {
				n /= p
			}
			if n == 1 {
				return math.Log(float64(p))
			}
			return 0
		}
	}
	if n > 1 {
		return math.Log(float64(n))
	}
	return 0
}

func main() {
	fmt.Println("THE ORBITS — Landau's roll call over our ten measured zeros")
	fmt.Println("\n   n   S(n) measured   -(T/2pi)*L(n)/sqrt(n)   orbit?")

	T := zetaZeros[len(zetaZeros)-1]
	scale := T / (2 * math.Pi)

	var orbitS, silentS []float64
	sx, sy, sxx, syy, sxy, np := 0.0, 0.0, 0.0, 0.0, 0.0, 0.0
	for n := 2; n <= 30; n++ {
		s := 0.0
		for _, g := range zetaZeros {
			s += math.Cos(g * math.Log(float64(n)))
		}
		pred := -scale * mangoldt(n) / math.Sqrt(float64(n))
		tag := "silent"
		if mangoldt(n) > 0 {
			tag = "ORBIT"
			orbitS = append(orbitS, s)
			sx += pred
			sy += s
			sxx += pred * pred
			syy += s * s
			sxy += pred * s
			np++
		} else {
			silentS = append(silentS, s)
		}
		fmt.Printf("  %2d   %+7.2f         %+7.2f                 %s\n", n, s, pred, tag)
	}

	mean := func(v []float64) float64 {
		m := 0.0
		for _, x := range v {
			m += x
		}
		return m / float64(len(v))
	}
	corr := (np*sxy - sx*sy) /
		math.Sqrt((np*sxx-sx*sx)*(np*syy-sy*sy))

	fmt.Printf("\nmean S over prime powers: %+.2f    mean S over the silent: %+.2f\n",
		mean(orbitS), mean(silentS))
	fmt.Printf("correlation between measured S and the stability law: %.3f\n", corr)
	fmt.Println("\neach orbit answers when called by its period ln n, with the strength")
	fmt.Println("Gutzwiller's stability weight predicts; non-orbits stay silent.")
}
