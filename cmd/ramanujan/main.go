// Command ramanujan climbs to the second floor of the Langlands tower.
//
// Every dial so far was a GL(1) instrument — Dirichlet characters, degree
// one. The first citizen of the second floor is Ramanujan's Δ, the weight-12
// modular form whose coefficients τ(n) this laboratory already verified to
// be multiplicative with the Hecke recursion (Finding 51's neighbourhood).
// Its L-function has its own zeros — never measured by this laboratory —
// and its von Mangoldt weights are Λ_Δ(p^k) = (α^k + ᾱ^k)·ln p with
// α + ᾱ = τ(p)/p^{11/2}, |α| = 1 (Deligne).
//
// The radio: τ(n) computed exactly from η(z)^24 by 24 sparse pentagonal
// multiplications; the signal ψ_Δ(x)/√x has no smooth term (L(Δ) is
// entire), so its periodogram peaks ARE the second floor's stations.
//
// PRE-REGISTERED: distinct stations must appear below γ ≈ 23 and match the
// published zeros of L(s, Δ) on subsequent LMFDB verification. Instrument
// checks: τ(2) = −24 and τ(25) = −25499225 exactly, and every |a_p| ≤ 2
// (Deligne's bound) — a precision alarm for the float convolution.
//
// Usage:
//
//	go run ./cmd/ramanujan [-limit N]
package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"sort"

	"github.com/nicoconvoz/diosyunalma/primes"
	"github.com/nicoconvoz/diosyunalma/spectral"
)

func main() {
	limit := flag.Int("limit", 300_000, "compute tau(n) up to this value")
	flag.Parse()
	if *limit < 50_000 {
		fmt.Fprintln(os.Stderr, "limit must be at least 50000")
		os.Exit(1)
	}
	n := *limit

	// eta(z)^24: start from the series 1 and multiply 24 times by the
	// pentagonal series 1 + sum s_i q^(g_i).
	type pent struct {
		g int
		s float64
	}
	ps := []pent{}
	for k := 1; ; k++ {
		g1 := k * (3*k - 1) / 2
		g2 := k * (3*k + 1) / 2
		if g1 > n && g2 > n {
			break
		}
		s := 1.0
		if k%2 == 1 {
			s = -1.0
		}
		if g1 <= n {
			ps = append(ps, pent{g1, s})
		}
		if g2 <= n {
			ps = append(ps, pent{g2, s})
		}
	}
	a := make([]float64, n+1)
	a[0] = 1
	for rep := 0; rep < 24; rep++ {
		for i := n; i >= 1; i-- {
			for _, p := range ps {
				if p.g > i {
					break
				}
				a[i] += p.s * a[i-p.g]
			}
		}
	}
	tau := func(m int) float64 { return a[m-1] }

	// instrument checks.
	if tau(2) != -24 || tau(25) != -25499225 {
		fmt.Fprintf(os.Stderr, "precision alarm: tau(2)=%v tau(25)=%v\n", tau(2), tau(25))
		os.Exit(1)
	}
	pr := primes.Sieve(n)
	maxA := 0.0
	ap := map[int]float64{}
	for _, p := range pr {
		v := tau(p) / math.Pow(float64(p), 5.5)
		ap[p] = v
		if math.Abs(v) > maxA {
			maxA = math.Abs(v)
		}
	}
	fmt.Printf("THE RAMANUJAN RADIO — tau exact to %d, max |a_p| = %.6f (Deligne demands <= 2)\n", n, maxA)
	if maxA > 2 {
		fmt.Println("precision alarm: Deligne violated - do not trust the dial")
		os.Exit(1)
	}

	// the signal: psi_Delta(x)/sqrt(x) on the log grid.
	const du = 0.002
	var us []float64
	var xs []int
	for u := math.Log(1000); u <= math.Log(float64(n)); u += du {
		us = append(us, u)
		xs = append(xs, int(math.Round(math.Exp(u))))
	}
	type event struct {
		at int
		v  float64
	}
	evs := []event{}
	for _, p := range pr {
		lg := math.Log(float64(p))
		b1, b0 := ap[p], 2.0
		evs = append(evs, event{p, b1 * lg})
		bk1, bk2 := b1, b0
		for pk := p * p; pk <= n && pk > 0; pk *= p {
			bk := ap[p]*bk1 - bk2
			evs = append(evs, event{pk, bk * lg})
			bk2, bk1 = bk1, bk
		}
	}
	sort.Slice(evs, func(i, j int) bool { return evs[i].at < evs[j].at })
	es := make([]float64, len(xs))
	sum, ei := 0.0, 0
	for i, x := range xs {
		for ei < len(evs) && evs[ei].at <= x {
			sum += evs[ei].v
			ei++
		}
		es[i] = sum / math.Sqrt(float64(x))
	}
	m := float64(len(es) - 1)
	for i := range es {
		es[i] *= 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/m))
	}

	freqs := []float64{}
	for f := 3.0; f <= 26.0; f += 0.005 {
		freqs = append(freqs, f)
	}
	pw := spectral.Periodogram(us, es, freqs)
	type peak struct{ f, p float64 }
	cands := []peak{}
	for i := 1; i+1 < len(pw); i++ {
		if pw[i] > pw[i-1] && pw[i] > pw[i+1] {
			cands = append(cands, peak{freqs[i], pw[i]})
		}
	}
	sort.Slice(cands, func(i, j int) bool { return cands[i].p > cands[j].p })
	out := []peak{}
	for _, c := range cands {
		ok := true
		for _, k := range out {
			if math.Abs(k.f-c.f) < 1.2 {
				ok = false
				break
			}
		}
		if ok {
			out = append(out, c)
		}
		if len(out) == 5 {
			break
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].f < out[j].f })
	fmt.Print("\nsecond-floor stations heard:")
	for _, k := range out {
		fmt.Printf("  %.3f", k.f)
	}
	fmt.Println()
	fmt.Println("\nthe first GL(2) music this laboratory has ever tuned: the song of a")
	fmt.Println("modular form, whose whole melody is written by its primes alone.")
}
