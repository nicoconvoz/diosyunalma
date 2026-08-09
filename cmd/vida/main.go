// Command vida sounds the captain's flash: if life is movement, the
// primes are the only movers - von Mangoldt silences every number that
// is not a prime power (the lifeless play nothing; prime powers echo
// their prime). And life's own song, woven BETWEEN the primes' songs
// with alternating directions and shifting extremes, is the Moebius
// walk M(x): +1/-1 by prime parity, whose ultimate bound IS the
// Riemann Hypothesis. The duality closes: the primes play the zeros'
// tide (F100); here we listen for the zeros playing M(x)'s melody.
//
// Test: sieve mu(n) to 10^7, walk M(x), then project M(x)/sqrt(x) on
// frequencies gamma in ln x: peaks at OUR zeros vs mid-gap controls.
package main

import (
	"fmt"
	"math"
	"math/cmplx"
)

const emTerms = 120

func zetaEM(s complex128, lnn []float64) complex128 {
	var sum complex128
	sig, t := real(s), imag(s)
	for n := 1; n < emTerms; n++ {
		amp := math.Exp(-sig * lnn[n])
		sn, cs := math.Sincos(t * lnn[n])
		sum += complex(amp*cs, -amp*sn)
	}
	nf := complex(float64(emTerms), 0)
	ns := cmplx.Exp(-s * complex(lnn[emTerms], 0))
	sum += ns * nf / (s - 1)
	sum += ns / 2
	sum += ns * s / nf / 12
	sum -= ns * s * (s + 1) * (s + 2) / (nf * nf * nf) / 720
	return sum
}

func theta(t float64) float64 {
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t)
}

func zOf(t float64, lnn []float64) float64 {
	z := zetaEM(complex(0.5, t), lnn)
	th := theta(t)
	return real(z)*math.Cos(th) - imag(z)*math.Sin(th)
}

func main() {
	fmt.Println("LA CANCIÓN DE LA VIDA — the Moebius walk and who plays it")

	// the walk: mu(n) by linear sieve to N.
	const N = 10000000
	mu := make([]int8, N+1)
	primes := []int32{}
	mu[1] = 1
	composite := make([]bool, N+1)
	for i := 2; i <= N; i++ {
		if !composite[i] {
			primes = append(primes, int32(i))
			mu[i] = -1
		}
		for _, p := range primes {
			ip := int64(i) * int64(p)
			if ip > N {
				break
			}
			composite[ip] = true
			if i%int(p) == 0 {
				mu[ip] = 0
				break
			}
			mu[ip] = -mu[i]
		}
	}

	// alternating directions, shifting extremes.
	var m int64
	minM, maxM := int64(0), int64(0)
	minX, maxX := 1, 1
	worstRatio := 0.0
	for n := 1; n <= N; n++ {
		m += int64(mu[n])
		if m < minM {
			minM, minX = m, n
		}
		if m > maxM {
			maxM, maxX = m, n
		}
		if n > 100 {
			if r := math.Abs(float64(m)) / math.Sqrt(float64(n)); r > worstRatio {
				worstRatio = r
			}
		}
	}
	fmt.Printf("\n  the walk to N=%d: extremes %+d (at %d) / %+d (at %d)\n", N, maxM, maxX, minM, minX)
	fmt.Printf("  the Hypothesis bound: max |M(x)|/sqrt(x) = %.3f  (RH says this stays under x^eps forever)\n", worstRatio)

	// our own zeros as the candidate notes.
	lnn := make([]float64, emTerms+1)
	for n := 1; n <= emTerms; n++ {
		lnn[n] = math.Log(float64(n))
	}
	var gammas []float64
	prevT, prevZ := 12.0, zOf(12.0, lnn)
	for t := 12.02; len(gammas) < 12; t += 0.02 {
		zt := zOf(t, lnn)
		if (prevZ < 0) != (zt < 0) {
			lo, hi := prevT, t
			for i := 0; i < 50 && hi-lo > 1e-9; i++ {
				mid := (lo + hi) / 2
				if (zOf(mid, lnn) < 0) == (prevZ < 0) {
					lo = mid
				} else {
					hi = mid
				}
			}
			gammas = append(gammas, (lo+hi)/2)
		}
		prevT, prevZ = t, zt
	}

	// listen: project M(x)/sqrt(x) on cos/sin(gamma ln x) over a log
	// grid; controls at mid-gap frequencies.
	proj := func(g float64) float64 {
		var cr, ci float64
		np := 0
		m := int64(0)
		nextLog := 100.0
		for n := 1; n <= N; n++ {
			m += int64(mu[n])
			if float64(n) >= nextLog {
				lx := math.Log(float64(n))
				v := float64(m) / math.Sqrt(float64(n))
				sn, cs := math.Sincos(g * lx)
				cr += v * cs
				ci += v * sn
				np++
				nextLog *= 1.001
			}
		}
		return math.Hypot(cr, ci) / float64(np)
	}
	fmt.Println("\n  the song's notes: |projection of M(x)/sqrt(x)| at frequency gamma")
	fmt.Println("     zero gamma     at the ZERO    at mid-gap control")
	for i := 0; i+1 < len(gammas); i++ {
		g := gammas[i]
		ctrl := (gammas[i] + gammas[i+1]) / 2
		fmt.Printf("     %9.4f      %8.4f       %8.4f\n", g, proj(g), proj(ctrl))
	}
	fmt.Println("\n  (notes ringing at OUR zeros and quiet between them = life's song is")
	fmt.Println("   sung in the zeros' scale - the F100 duality, heard from the other shore)")
}
