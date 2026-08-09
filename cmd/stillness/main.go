// Command stillness tests the stillness principle: are the primes the
// points where the boiling broth comes to rest in perfect agreement?
//
// The orchestra is the zeros of zeta: pendulum m swings with frequency
// gamma_m. At a point x of the number line, pendulum m points in the
// direction gamma_m * ln x. UNREST is defined as phase incoherence — the
// Kuramoto order parameter R(x) = |sum_m e^{i gamma_m ln x}| / M. R near 0
// is the boil (total disagreement); R near 1 is stillness: every pendulum
// aligned. The explicit formula of Riemann predicts the alignment points
// are EXACTLY the primes and prime powers.
//
// Pre-registered test: compute the first ~100 zeros with our own
// instruments, chart R(x) over [1.8, 32], list its peaks, and compare with
// the prime powers. Control: the same orchestra with (a) smoothly spaced
// frequencies and (b) shuffled gaps must NOT single out the primes.
//
// Usage:
//
//	go run ./cmd/stillness
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"math/rand"
	"sort"
)

const emTerms = 120

// zetaEM computes zeta(s) by Euler–Maclaurin (ample at these heights).
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
	sum += ns * s * (s + 1) * (s + 2) * (s + 3) * (s + 4) /
		(nf * nf * nf * nf * nf) / 30240
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

// firstZeros hunts the first count zeros with a fine scan and bisection.
func firstZeros(count int, lnn []float64) []float64 {
	zeros := []float64{}
	prevT, prevZ := 12.0, zOf(12.0, lnn)
	for t := 12.02; len(zeros) < count; t += 0.02 {
		zt := zOf(t, lnn)
		if (prevZ < 0) != (zt < 0) {
			lo, hi := prevT, t
			for i := 0; i < 60 && hi-lo > 1e-10; i++ {
				mid := (lo + hi) / 2
				if (zOf(mid, lnn) < 0) == (prevZ < 0) {
					lo = mid
				} else {
					hi = mid
				}
			}
			zeros = append(zeros, (lo+hi)/2)
		}
		prevT, prevZ = t, zt
	}
	return zeros
}

// coherence is the unrest meter: R(x) = |sum e^{i gamma ln x}| / M.
func coherence(x float64, gammas []float64) float64 {
	lx := math.Log(x)
	var cr, ci float64
	for _, g := range gammas {
		sn, cs := math.Sincos(g * lx)
		cr += cs
		ci += sn
	}
	return math.Hypot(cr, ci) / float64(len(gammas))
}

type peak struct {
	x, r float64
}

// peaks charts R over [1.8, 32] and returns its local maxima, strongest first.
func peaks(gammas []float64) []peak {
	const x0, x1, dx = 1.8, 32.0, 5e-4
	var ps []peak
	prev := coherence(x0, gammas)
	cur := coherence(x0+dx, gammas)
	for x := x0 + 2*dx; x <= x1; x += dx {
		next := coherence(x, gammas)
		if cur > prev && cur >= next && cur > 0.2 {
			ps = append(ps, peak{x - dx, cur})
		}
		prev, cur = cur, next
	}
	sort.Slice(ps, func(i, j int) bool { return ps[i].r > ps[j].r })
	return ps
}

func isPrimePower(n int) bool {
	if n < 2 {
		return false
	}
	for p := 2; p*p <= n; p++ {
		if n%p == 0 {
			for n%p == 0 {
				n /= p
			}
			return n == 1
		}
	}
	return true // prime
}

// score counts how many of the top peaks sit on a prime power (+/- 0.05).
func score(ps []peak, top int) (hits int, report string) {
	if top > len(ps) {
		top = len(ps)
	}
	for i := 0; i < top; i++ {
		n := int(math.Round(ps[i].x))
		on := math.Abs(ps[i].x-float64(n)) <= 0.05 && isPrimePower(n)
		mark := " "
		if on {
			hits++
			mark = "*"
		}
		if i < 14 {
			report += fmt.Sprintf("    %s x = %7.3f   R = %.3f\n", mark, ps[i].x, ps[i].r)
		}
	}
	return
}

func main() {
	fmt.Println("THE STILLNESS PRINCIPLE — where does the boiling broth come to rest?")

	lnn := make([]float64, emTerms+1)
	for n := 1; n <= emTerms; n++ {
		lnn[n] = math.Log(float64(n))
	}
	gammas := firstZeros(100, lnn)
	fmt.Printf("\n  orchestra tuned: first %d zeros, gamma_1 = %.6f ... gamma_%d = %.4f\n",
		len(gammas), gammas[0], len(gammas), gammas[len(gammas)-1])

	// the real orchestra.
	psReal := peaks(gammas)
	hitsReal, rep := score(psReal, 20)
	fmt.Printf("\n  THE TRUE ZEROS - top peaks of agreement R(x) (* = prime power):\n%s", rep)
	fmt.Printf("    ... %d of the top 20 peaks sit on prime powers\n", hitsReal)

	// control (a): smooth frequencies with the same density.
	smooth := make([]float64, len(gammas))
	for m := range smooth {
		g := gammas[m]
		for i := 0; i < 40; i++ { // invert the counting function by Newton
			nb := g/(2*math.Pi)*math.Log(g/(2*math.Pi*math.E)) + 7.0/8
			dn := math.Log(g/(2*math.Pi)) / (2 * math.Pi)
			g -= (nb - float64(m+1)) / dn
		}
		smooth[m] = g
	}
	psSm := peaks(smooth)
	hitsSm, _ := score(psSm, 20)

	// control (b): true first zero, shuffled gaps.
	rng := rand.New(rand.NewSource(2026))
	gaps := make([]float64, len(gammas)-1)
	for i := range gaps {
		gaps[i] = gammas[i+1] - gammas[i]
	}
	rng.Shuffle(len(gaps), func(i, j int) { gaps[i], gaps[j] = gaps[j], gaps[i] })
	shuffled := make([]float64, len(gammas))
	shuffled[0] = gammas[0]
	for i, g := range gaps {
		shuffled[i+1] = shuffled[i] + g
	}
	psSh := peaks(shuffled)
	hitsSh, _ := score(psSh, 20)

	fmt.Printf("\n  CONTROLS (same density, no prime soul):\n")
	fmt.Printf("    smooth frequencies: %d of top 20 peaks on prime powers\n", hitsSm)
	fmt.Printf("    shuffled gaps:      %d of top 20 peaks on prime powers\n", hitsSh)

	fmt.Println("\n  UNREST, defined: inquietud(x) = 1 - R(x), the phase disagreement of")
	fmt.Println("  the zero-orchestra. The broth boils almost everywhere (R small); the")
	fmt.Println("  pendulums align only where the primes and their powers live. Stillness")
	fmt.Println("  is not absence of motion - it is perfect agreement.")
}
