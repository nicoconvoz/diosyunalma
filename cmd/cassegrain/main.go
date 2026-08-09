// Command cassegrain is the engine's crankshaft: the exact Gauss-sum
// folder, certified on its own test bench.
//
// The first build bounced two same-sized mirrors against each other
// forever — the bench caught the infinite ping-pong, and the fix is the
// classical one: the bounce cascade must terminate in GAUSS'S CLOSED FORM
// (1805). Every complete quadratic sum collapses to an eighth root of
// unity times √N, with the root decided by the Jacobi symbol — whose
// computation is itself a Euclid-like cascade of log-many bounces. One
// Landsberg–Schaar bounce handles the odd-numerator case; Gauss does the
// rest. A billion-term wave becomes a phase and a square root, exactly.
//
// The gearbox joining this to incomplete, amplitude-weighted sums (the
// Fresnel seams of the full t^(1/3) engine) is the shipyard's registered
// next stage — the honest distance still separating the fingers.
//
// Usage:
//
//	go run ./cmd/cassegrain
package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

var bounceCount int

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// jacobi computes the Jacobi symbol (a|n) for odd n>0 - the bounce cascade.
func jacobi(a, n int64) int64 {
	a %= n
	if a < 0 {
		a += n
	}
	result := int64(1)
	for a != 0 {
		for a%2 == 0 {
			a /= 2
			if r := n % 8; r == 3 || r == 5 {
				result = -result
			}
			bounceCount++
		}
		a, n = n, a
		bounceCount++
		if a%4 == 3 && n%4 == 3 {
			result = -result
		}
		a %= n
	}
	if n == 1 {
		return result
	}
	return 0
}

// s2 computes S(a,N) = sum_{n<N} e^{2 pi i a n^2 / N} exactly (Gauss).
func s2(a, N int64) (float64, float64) {
	a %= N
	if a < 0 {
		a += N
	}
	if N == 1 {
		return 1, 0
	}
	if a == 0 {
		return float64(N), 0
	}
	if g := gcd(a, N); g > 1 {
		r, i := s2(a/g, N/g)
		return float64(g) * r, float64(g) * i
	}
	rt := math.Sqrt(float64(N))
	switch N % 4 {
	case 1:
		j := float64(jacobi(a, N))
		return j * rt, 0
	case 3:
		j := float64(jacobi(a, N))
		return 0, j * rt
	case 2:
		return 0, 0
	default: // N = 0 mod 4, a odd
		// S = (1+i) * conj(eps_a) * (N|a) * sqrt(N), eps_a = 1 or i.
		j := float64(jacobi(N, a))
		var er, ei float64 // conj(eps_a)
		if a%4 == 1 {
			er, ei = 1, 0
		} else {
			er, ei = 0, -1
		}
		// (1+i)*(er+ei i) = (er-ei) + (er+ei) i
		cr, ci := er-ei, er+ei
		return j * cr * rt, j * ci * rt
	}
}

// gp computes G(p,q) = sum_{n<q} e^{i pi p n^2 / q}, pq even.
func gp(p, q int64) (float64, float64) {
	p %= 2 * q
	if p < 0 {
		p += 2 * q
	}
	if p%2 == 0 {
		return s2(p/2, q)
	}
	// p odd (so q even): one Landsberg-Schaar bounce onto an odd modulus.
	// G(p,q) = e^{i pi/4} sqrt(q/p) conj(G(q mod 2p, p)).
	bounceCount++
	gr, gi := gp(q%(2*p), p)
	gi = -gi
	s := math.Sqrt(float64(q) / float64(p))
	c := math.Sqrt2 / 2
	return s * c * (gr - gi), s * c * (gr + gi)
}

// direct computes G(p,q) by brute force with exact integer phases.
func direct(p, q int64) (float64, float64) {
	var sr, si float64
	for n := int64(0); n < q; n++ {
		ph := float64((p*n%(2*q))*n%(2*q)) * math.Pi / float64(q)
		sr += math.Cos(ph)
		si += math.Sin(ph)
	}
	return sr, si
}

func main() {
	fmt.Println("THE CASSEGRAIN CRANKSHAFT — Gauss's closed form at the end of the bounces")

	rng := rand.New(rand.NewSource(2026))
	const trials = 500
	worst := 0.0
	worstP, worstQ := int64(0), int64(0)
	bounceCount = 0
	for i := 0; i < trials; i++ {
		q := int64(rng.Intn(9000) + 1000)
		p := int64(rng.Intn(int(2*q-2)) + 1)
		if (p*q)%2 == 1 {
			p++
		}
		dr, di := direct(p, q)
		fr, fi := gp(p, q)
		err := math.Hypot(fr-dr, fi-di) / math.Sqrt(float64(q))
		if err > worst {
			worst, worstP, worstQ = err, p, q
		}
	}
	fmt.Printf("\n  bench: %d random complete Gauss sums, q up to 10000\n", trials)
	fmt.Printf("  worst relative error vs direct summation: %.2e", worst)
	if worst > 1e-9 {
		fmt.Printf("   (worst case p=%d q=%d)", worstP, worstQ)
	}
	fmt.Println()
	fmt.Printf("  total mirror bounces across the bench: %d (~%.1f per fold)\n",
		bounceCount, float64(bounceCount)/trials)

	p, q := int64(246913578), int64(999999937)
	t0 := time.Now()
	bounceCount = 0
	fr, fi := gp(p, q)
	el := time.Since(t0)
	fmt.Printf("\n  showpiece: G(p=%d, q=%d)\n", p, q)
	fmt.Printf("  a %d-term wave folded through %d bounces in %s\n", q, bounceCount, el)
	fmt.Printf("  value = %.6f %+.6fi   |G|/sqrt(q) = %.6f\n",
		fr, fi, math.Hypot(fr, fi)/math.Sqrt(float64(q)))
	fmt.Println("\na billion-scale wave collapsed to a phase and a square root, exactly:")
	fmt.Println("the crankshaft turns. the Fresnel gearbox remains the registered next")
	fmt.Println("stage - the honest distance between the fingers.")
}
